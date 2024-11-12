package courses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"search-api/dao/courses"

	"github.com/stevenferrer/solr-go"
)

type SolrConfig struct {
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	// Construct the BaseURL using the provided host and port
	baseURL := fmt.Sprintf("http://%s:%s", config.Host, config.Port)
	client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
	}
}

// Index adds a new course document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, course courses.Course) (string, error) {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"id":           course.ID,
		"name":         course.Name,
		"description":  course.Description,
		"professor":    course.Professor,
		"image_url":    course.ImageURL,
		"duration":     course.Duration,
		"requirement":  course.Requirement,
		"availability": course.Availability,
	}

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in Solr
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling course document: %w", err)
	}

	// Index the document in Solr
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error indexing course: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("failed to index course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return course.ID, nil
}

// Update modifies an existing course document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, course courses.Course) error {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"id":           course.ID,
		"name":         course.Name,
		"description":  course.Description,
		"professor":    course.Professor,
		"image_url":    course.ImageURL,
		"duration":     course.Duration,
		"requirement":  course.Requirement,
		"availability": course.Availability,
	}

	// Prepare the update request
	updateRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Update the document in Solr
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	// Execute the update request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating course: %v", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	// Prepare the delete request
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	// Update the document in Solr
	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	// Execute the delete request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting course: %v", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to index course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]courses.Course, error) {
	// Prepare the Solr query with limit and offset
	// solrQuery := fmt.Sprintf("q=(name:%s)&rows=%d&start=%d", query, limit, offset)
	solrQuery := fmt.Sprintf("q=(name:%s OR description:%s)&rows=%d&start=%d", query, query, limit, offset)

	// Execute the search request
	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %v", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to execute search query: %v", resp.Error)
	}

	// Parse the response and extract course documents
	var coursesList []courses.Course
	for _, doc := range resp.Response.Documents {

		// Safely extract course fields with type assertions
		course := courses.Course{
			ID:           getStringField(doc, "id"),
			Name:         getStringField(doc, "name"),
			Description:  getStringField(doc, "description"),
			Professor:    getStringField(doc, "professor"),
			ImageURL:     getStringField(doc, "image_url"),
			Duration:     getFloatField(doc, "duration"),
			Requirement:  getStringField(doc, "requirement"),
			Availability: getFloatField(doc, "availability"),
		}
		coursesList = append(coursesList, course)
	}

	return coursesList, nil
}

// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get float64 fields from the document
func getFloatField(doc map[string]interface{}, field string) float64 {
	if val, ok := doc[field].(float64); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if floatVal, ok := val[0].(float64); ok {
			return floatVal
		}
	}
	return 0.0
}
