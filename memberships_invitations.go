package taigo

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Membership is a raw DTO for /memberships endpoints.
type Membership = RawResource

// MembershipInvitation is a raw DTO for membership invitation records.
type MembershipInvitation = RawResource

// MembershipsQueryParams holds optional list filters for memberships.
type MembershipsQueryParams struct {
	Project int `url:"project,omitempty"`
	Role    int `url:"role,omitempty"`
}

// MembershipInvitationsQueryParams holds optional list filters for invitations.
type MembershipInvitationsQueryParams struct {
	Project int `url:"project,omitempty"`
}

// MembershipInvitationService is a handle to actions related to memberships and invitations.
type MembershipInvitationService struct {
	client               *Client
	defaultProjectID     int
	Endpoint             string
	PublicInvitationPath string
}

// ListMemberships -> https://docs.taiga.io/api.html#memberships-list
func (s *MembershipInvitationService) ListMemberships(queryParams *MembershipsQueryParams) ([]Membership, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, queryParams)
}

// GetMembership -> https://docs.taiga.io/api.html#memberships-get
func (s *MembershipInvitationService) GetMembership(membershipID int) (*Membership, error) {
	return getRawResource(s.client, s.Endpoint, membershipID)
}

// EditMembership -> https://docs.taiga.io/api.html#memberships-edit
func (s *MembershipInvitationService) EditMembership(membershipID int, payload any) (*Membership, error) {
	return patchRawResource(s.client, s.Endpoint, membershipID, payload)
}

// UpdateMembership is an alias for EditMembership.
func (s *MembershipInvitationService) UpdateMembership(membershipID int, payload any) (*Membership, error) {
	return s.EditMembership(membershipID, payload)
}

// DeleteMembership -> https://docs.taiga.io/api.html#memberships-delete
func (s *MembershipInvitationService) DeleteMembership(membershipID int) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, membershipID)
}

// BulkCreateMemberships -> https://docs.taiga.io/api.html#memberships-bulk-create
func (s *MembershipInvitationService) BulkCreateMemberships(payload any) ([]Membership, error) {
	return postRawResourceListAtPath(s.client, payload, s.Endpoint, "bulk_create")
}

// CreateMembership -> https://docs.taiga.io/api.html#memberships-create
func (s *MembershipInvitationService) CreateMembership(payload any) (*Membership, error) {
	return createRawResource(s.client, s.Endpoint, payload)
}

// ListInvitations -> https://docs.taiga.io/api.html#memberships-invitations-list
// Invitation resources are represented by memberships in Taiga.
func (s *MembershipInvitationService) ListInvitations(queryParams *MembershipInvitationsQueryParams) ([]MembershipInvitation, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, queryParams)
}

// CreateInvitation -> https://docs.taiga.io/api.html#memberships-invitations-create
// Invitation resources are created through memberships.
func (s *MembershipInvitationService) CreateInvitation(payload any) (*MembershipInvitation, error) {
	return s.CreateMembership(payload)
}

// GetInvitation -> https://docs.taiga.io/api.html#memberships-invitations-get
// Invitation resources are represented by memberships in Taiga.
func (s *MembershipInvitationService) GetInvitation(invitationID int) (*MembershipInvitation, error) {
	return getRawResource(s.client, s.Endpoint, invitationID)
}

// ResendInvitation -> https://docs.taiga.io/api.html#memberships-invitations-resend
// https://docs.taiga.io/api.html#memberships-resend-invitation
func (s *MembershipInvitationService) ResendInvitation(invitationID int) (*MembershipInvitation, error) {
	return postRawResourceAtPath(s.client, nil, s.Endpoint, strconv.Itoa(invitationID), "resend_invitation")
}

// DeleteInvitation -> https://docs.taiga.io/api.html#memberships-invitations-delete
// Invitation resources are represented by memberships in Taiga.
func (s *MembershipInvitationService) DeleteInvitation(invitationID int) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, invitationID)
}

// GetInvitationByToken resolves invitation details from a public invitation token.
// https://docs.taiga.io/api.html#invitations
func (s *MembershipInvitationService) GetInvitationByToken(invitationUUID string) (*MembershipInvitation, error) {
	invitationUUID = strings.TrimSpace(invitationUUID)
	if invitationUUID == "" {
		return nil, errors.New("invitationUUID is required")
	}
	invitation, err := getRawResourceAtPath(s.client, s.PublicInvitationPath, invitationUUID)
	if err == nil {
		return invitation, nil
	}
	// Older Taiga versions may expose this endpoint as POST-only.
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusMethodNotAllowed {
		return postRawResourceAtPath(s.client, nil, s.PublicInvitationPath, invitationUUID)
	}
	return nil, err
}

// ApplyInvitationByToken applies invitation data for a public invitation token.
func (s *MembershipInvitationService) ApplyInvitationByToken(invitationUUID string, payload any) (*MembershipInvitation, error) {
	invitationUUID = strings.TrimSpace(invitationUUID)
	if invitationUUID == "" {
		return nil, errors.New("invitationUUID is required")
	}
	if err := requireNonNil("payload", payload); err != nil {
		return nil, err
	}
	return postRawResourceAtPath(s.client, payload, s.PublicInvitationPath, invitationUUID)
}
