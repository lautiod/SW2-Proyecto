package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	coursesDomain "search-api/domain/courses"
)

type HTTPConfig struct {
	Host string
	Port string
}

type HTTP struct {
	baseURL func(courselID string) string
}

func NewHTTP(config HTTPConfig) HTTP {
	return HTTP{
		baseURL: func(courseID string) string {
			return fmt.Sprintf("http://%s:%s/courses/%s", config.Host, config.Port, courseID)
		},
	}
}

func (repository HTTP) GetCourseByID(ctx context.Context, id string) (coursesDomain.Course, error) {
	resp, err := http.Get(repository.baseURL(id))
	if err != nil {
		return coursesDomain.Course{}, fmt.Errorf("error fetching course (%s): %v", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return coursesDomain.Course{}, fmt.Errorf("failed to fetch course (%s): received status code %d", id, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return coursesDomain.Course{}, fmt.Errorf("error reading response body for course (%s): %v", id, err)
	}

	// Unmarshal the course details into the course struct
	var course coursesDomain.Course
	if err := json.Unmarshal(body, &course); err != nil {
		return coursesDomain.Course{}, fmt.Errorf("error unmarshaling course data (%s): %v", id, err)
	}

	return course, nil
}
