package courses

import (
	"context"
	"fmt"
	"log"

	inscriptionsDAO "inscriptions-api/dao/inscriptions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
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
		client:     client,
		database:   config.Database,
		collection: config.Collection,
	}
}

func (repository Mongo) EnrollUser(ctx context.Context, inscription inscriptionsDAO.Inscription) (string, error) {
	// Insert into mongo
	result, err := repository.client.Database(repository.database).Collection(repository.collection).InsertOne(ctx, inscription)
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
	result := repository.client.Database(repository.database).Collection(repository.collection).FindOne(ctx, bson.M{"course_id": inscription.CourseID, "user_id": inscription.UserID})

	if result.Err() != nil {
		return fmt.Errorf("error finding document: %v", result.Err())
	}

	return nil
}
