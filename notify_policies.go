package taigo

// NotifyPolicy is a raw DTO for /notify-policies endpoints.
type NotifyPolicy = RawResource

// NotifyPolicyService is a handle to actions related to notification policies.
type NotifyPolicyService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#notify-policies-list
func (s *NotifyPolicyService) List() ([]NotifyPolicy, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, nil)
}

// Get -> https://docs.taiga.io/api.html#notify-policies-get
func (s *NotifyPolicyService) Get(policyID int) (*NotifyPolicy, error) {
	return getRawResource(s.client, s.Endpoint, policyID)
}

// Replace -> https://docs.taiga.io/api.html#notify-policies-edit
func (s *NotifyPolicyService) Replace(policyID int, payload any) (*NotifyPolicy, error) {
	return putRawResource(s.client, s.Endpoint, policyID, payload)
}

// Edit -> https://docs.taiga.io/api.html#notify-policies-edit
func (s *NotifyPolicyService) Edit(policyID int, payload any) (*NotifyPolicy, error) {
	return patchRawResource(s.client, s.Endpoint, policyID, payload)
}

// Update is an alias for Edit.
func (s *NotifyPolicyService) Update(policyID int, payload any) (*NotifyPolicy, error) {
	return s.Edit(policyID, payload)
}
