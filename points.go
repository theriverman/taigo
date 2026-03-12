package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// Point -> https://docs.taiga.io/api.html#points
type Point struct {
	ID      int      `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Order   int      `json:"order,omitempty"`
	Project int      `json:"project,omitempty"`
	Value   *float64 `json:"value,omitempty"`
}

// PointService is a handle to actions related to points.
type PointService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#points-list
func (s *PointService) List(queryParams *ProjectIDQueryParams) ([]Point, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var points []Point
	_, err = s.client.Request.Get(url, &points)
	if err != nil {
		return nil, err
	}
	return points, nil
}

// Get -> https://docs.taiga.io/api.html#points-get
func (s *PointService) Get(pointID int) (*Point, error) {
	if err := requirePositiveID("pointID", pointID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(pointID))
	var point Point
	_, err := s.client.Request.Get(url, &point)
	if err != nil {
		return nil, err
	}
	return &point, nil
}

// Create -> https://docs.taiga.io/api.html#points-create
func (s *PointService) Create(point *Point) (*Point, error) {
	if err := requireNonNil("point", point); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responsePoint Point
	projectID, err := resolveProjectID(point.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(point.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *point
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responsePoint)
	if err != nil {
		return nil, err
	}
	return &responsePoint, nil
}

// Edit -> https://docs.taiga.io/api.html#points-edit
func (s *PointService) Edit(point *Point) (*Point, error) {
	if err := requireNonNil("point", point); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(point.ID))
	var responsePoint Point
	if err := requirePositiveID("pointID", point.ID); err != nil {
		return nil, err
	}
	payload, err := sparsePatchMapFromStruct(point, "id")
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responsePoint)
	if err != nil {
		return nil, err
	}
	return &responsePoint, nil
}

// Delete -> https://docs.taiga.io/api.html#points-delete
func (s *PointService) Delete(pointID int) (*http.Response, error) {
	if err := requirePositiveID("pointID", pointID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(pointID))
	return s.client.Request.Delete(url)
}
