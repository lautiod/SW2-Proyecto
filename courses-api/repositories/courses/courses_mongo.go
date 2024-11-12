package courses

import (
	"context"
	"fmt"
	"log"

	coursesDAO "courses-api/dao/courses"
	inscriptionsDAO "courses-api/dao/inscriptions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host                   string
	Port                   string
	Username               string
	Password               string
	Database               string
	CoursesCollection      string
	InscriptionsCollection string
}

type Mongo struct {
	client                 *mongo.Client
	database               string
	coursesCollection      string
	inscriptionsCollection string
}

const (
	connectionURI = "mongodb://%s:%s"
)

func NewMongo(config MongoConfig) Mongo {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx := context.Background()
	uri := fmt.Sprintf(connectionURI, config.Host, config.Port)
	cfg := options.Client().ApplyURI(uri).SetAuth(credentials)

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}

	return Mongo{
		client:                 client,
		database:               config.Database,
		coursesCollection:      config.CoursesCollection,
		inscriptionsCollection: config.InscriptionsCollection,
	}
}

func (repository Mongo) GetCourses(ctx context.Context) (coursesDAO.Courses, error) {
	// Ejecutar la consulta Find para obtener todos los documentos
	cursor, err := repository.client.Database(repository.database).Collection(repository.coursesCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error obtaining course list: %w", err)
	}
	defer cursor.Close(ctx) // Asegurarse de cerrar el cursor cuando ya no se necesite

	// Crear un slice para almacenar los cursos
	var courses coursesDAO.Courses // Ahora, solo necesitamos un slice de coursesDAO.Courses

	// Iterar sobre el cursor y decodificar cada documento en el slice
	for cursor.Next(ctx) {
		var course coursesDAO.Course
		if err := cursor.Decode(&course); err != nil {
			return nil, fmt.Errorf("error decoding course document: %w", err)
		}
		courses = append(courses, course)
	}

	// Verificar si hubo un error durante la iteración del cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Retornar la lista de cursos
	return courses, nil
}

func (repository Mongo) GetCourseByID(ctx context.Context, id string) (coursesDAO.Course, error) {
	// Get from MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return coursesDAO.Course{}, fmt.Errorf("error converting id to mongo ID: %w", err)
	}
	result := repository.client.Database(repository.database).Collection(repository.coursesCollection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return coursesDAO.Course{}, fmt.Errorf("error finding document: %w", result.Err())
	}

	// Convert document to DAO
	var courseDAO coursesDAO.Course
	if err := result.Decode(&courseDAO); err != nil {
		return coursesDAO.Course{}, fmt.Errorf("error decoding result: %w", err)
	}
	return courseDAO, nil
}

func (repository Mongo) CreateCourse(ctx context.Context, course coursesDAO.Course) (string, error) {
	// Insert into mongo
	result, err := repository.client.Database(repository.database).Collection(repository.coursesCollection).InsertOne(ctx, course)
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}

	// Get inserted ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("error converting mongo ID to object ID")
	}
	return objectID.Hex(), nil
}

func (repository Mongo) UpdateCourse(ctx context.Context, course coursesDAO.Course) error {
	// Convert hotel ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(course.ID)
	if err != nil {
		return fmt.Errorf("error converting id to mongo ID: %w", err)
	}

	// Create an update document
	update := bson.M{}

	// Only set the fields that are not empty or their default value
	if course.Name != "" {
		update["name"] = course.Name
	}
	if course.Description != "" {
		update["description"] = course.Description
	}
	if course.Professor != "" {
		update["professor"] = course.Professor
	}
	if course.ImageURL != "" {
		update["image_url"] = course.ImageURL
	}
	if course.Duration != 0 { // Assuming 0 is for not modification
		update["duration"] = course.Duration
	}
	if course.Requirement != "" { // Assuming 0 is for not modification
		update["requirement"] = course.Requirement
	}
	if course.Availability != 0 { // Assuming 0 is for not modification
		update["availability"] = course.Availability
	}

	// Update the document in MongoDB
	if len(update) == 0 {
		return fmt.Errorf("no fields to update for course ID %s", course.ID)
	}

	filter := bson.M{"_id": objectID}
	result, err := repository.client.Database(repository.database).Collection(repository.coursesCollection).UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID %s", course.ID)
	}

	return nil
}

// ******************************** I N S C R I P T I O N S

func (repository Mongo) EnrollUser(ctx context.Context, inscription inscriptionsDAO.Inscription) (string, error) {
	// Insert into mongo
	result, err := repository.client.Database(repository.database).Collection(repository.inscriptionsCollection).InsertOne(ctx, inscription)
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}

	// Get inserted ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("error converting mongo ID to object ID")
	}
	return objectID.Hex(), nil
}

func (repository Mongo) ValidateEnrrol(ctx context.Context, inscription inscriptionsDAO.Inscription) error {
	result := repository.client.Database(repository.database).Collection(repository.inscriptionsCollection).FindOne(ctx, bson.M{"course_id": inscription.CourseID, "user_id": inscription.UserID})

	if result.Err() != nil {
		return fmt.Errorf("error finding document: %v", result.Err())
	}

	return nil
}

func (repository Mongo) GetInscriptionsByUserId(ctx context.Context, userID string) ([]inscriptionsDAO.Inscription, error) {
	cursor, err := repository.client.Database(repository.database).Collection(repository.inscriptionsCollection).Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error obtaining courseIDs list: %w", err)
	}
	defer cursor.Close(ctx)

	// Crear un slice para almacenar los cursos
	var inscriptions []inscriptionsDAO.Inscription // Ahora, solo necesitamos un slice de coursesDAO.Courses

	// Iterar sobre el cursor y decodificar cada documento en el slice
	for cursor.Next(ctx) {
		var inscription inscriptionsDAO.Inscription
		if err := cursor.Decode(&inscription); err != nil {
			return nil, fmt.Errorf("error decoding inscription document: %w", err)
		}
		inscriptions = append(inscriptions, inscription)
	}

	// Verificar si hubo un error durante la iteración del cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Retornar la lista de cursos
	return inscriptions, nil
}
