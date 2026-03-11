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
	if isEmpty(point.Project) || isEmpty(point.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, &point, &responsePoint)
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
	if point.ID == 0 {
		return nil, errors.New("passed Point does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &point, &responsePoint)
	if err != nil {
		return nil, err
	}
	return &responsePoint, nil
}

// Delete -> https://docs.taiga.io/api.html#points-delete
func (s *PointService) Delete(pointID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(pointID))
	return s.client.Request.Delete(url)
}
