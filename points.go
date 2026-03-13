package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// Point -> https://docs.taiga.io/api.html#points
type Point struct {
	ID        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Order     int      `json:"order,omitempty"`
	ProjectID int      `json:"project_id,omitempty"`
	Value     *float64 `json:"value,omitempty"`
}

// PointCreateRequest represents payload for creating points.
type PointCreateRequest struct {
	Name    string   `json:"name"`
	Order   int      `json:"order,omitempty"`
	Project int      `json:"project"`
	Value   *float64 `json:"value,omitempty"`
}

// PointEditRequest represents sparse non-destructive updates for points.
type PointEditRequest struct {
	Name    string   `json:"name,omitempty"`
	Order   int      `json:"order,omitempty"`
	Project int      `json:"project,omitempty"`
	Value   *float64 `json:"value,omitempty"`
}

// PointPatch represents explicit PATCH payload for points.
type PointPatch struct {
	Name    *string  `json:"name,omitempty"`
	Order   *int     `json:"order,omitempty"`
	Project *int     `json:"project,omitempty"`
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
func (s *PointService) Create(request *PointCreateRequest) (*Point, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responsePoint Point
	projectID, err := resolveProjectID(request.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(request.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *request
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responsePoint)
	if err != nil {
		return nil, err
	}
	return &responsePoint, nil
}

// Edit -> https://docs.taiga.io/api.html#points-edit
func (s *PointService) Edit(pointID int, request *PointEditRequest) (*Point, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("pointID", pointID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(pointID))
	var responsePoint Point
	payload, err := sparsePatchMapFromStruct(request)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responsePoint)
	if err != nil {
		return nil, err
	}
	return &responsePoint, nil
}

// Update is an alias for Edit.
func (s *PointService) Update(pointID int, request *PointEditRequest) (*Point, error) {
	return s.Edit(pointID, request)
}

// Patch sends an explicit PATCH payload to edit a point.
func (s *PointService) Patch(pointID int, patch *PointPatch) (*Point, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("pointID", pointID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(pointID))
	var responsePoint Point
	_, err := s.client.Request.Patch(url, patch, &responsePoint)
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
