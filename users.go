package taigo

import (
	"errors"
	"strconv"
)

// UserService is a handle to actions related to Users
//
// https://taigaio.github.io/taiga-doc/dist/api.html#users
type UserService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// UserPatch represents an explicit PATCH payload for users.
// Pointer fields allow intentionally setting zero-values (false, 0, "").
type UserPatch struct {
	Bio          *string `json:"bio,omitempty"`
	Color        *string `json:"color,omitempty"`
	Email        *string `json:"email,omitempty"`
	FullName     *string `json:"full_name,omitempty"`
	Lang         *string `json:"lang,omitempty"`
	ReadNewTerms *bool   `json:"read_new_terms,omitempty"`
	Theme        *string `json:"theme,omitempty"`
	Timezone     *string `json:"timezone,omitempty"`
	Username     *string `json:"username,omitempty"`
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#users-list
func (s *UserService) List(queryParams *UsersQueryParams) ([]User, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var users []User
	_, err = s.client.Request.Get(url, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#users-get
func (s *UserService) Get(userID int) (*User, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID))
	var u User
	_, err := s.client.Request.Get(url, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Me => https://taigaio.github.io/taiga-doc/dist/api.html#users-me
func (s *UserService) Me() (*User, error) {
	var u User
	url := s.client.MakeURL(s.Endpoint, "me")
	_, err := s.client.Request.Get(url, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetStats => https://taigaio.github.io/taiga-doc/dist/api.html#users-stats
func (s *UserService) GetStats(userID int) (*UserStatsDetail, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "stats")
	var usd UserStatsDetail
	_, err := s.client.Request.Get(url, &usd)
	if err != nil {
		return nil, err
	}
	return &usd, nil
}

// GetWatchedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-watched
func (s *UserService) GetWatchedContent(userID int, queryParams *UsersHighlightedQueryParams) ([]UserWatched, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "watched")
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var watched []UserWatched
	_, err = s.client.Request.Get(url, &watched)
	if err != nil {
		return nil, err
	}
	return watched, nil
}

// GetLikedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-liked
func (s *UserService) GetLikedContent(userID int, queryParams *UsersHighlightedQueryParams) ([]UserLiked, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "liked")
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var liked []UserLiked
	_, err = s.client.Request.Get(url, &liked)
	if err != nil {
		return nil, err
	}
	return liked, nil
}

// https://taigaio.github.io/taiga-doc/dist/api.html#users-voted
func (s *UserService) GetVotedContent(userID int, queryParams *UsersHighlightedQueryParams) ([]Voted, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "voted")
	var err error
	url, err = appendQueryParams(url, queryParams)
	if err != nil {
		return nil, err
	}
	var voted []Voted
	_, err = s.client.Request.Get(url, &voted)
	if err != nil {
		return nil, err
	}
	return voted, nil
}

// Edit sends a PATCH request to edit a user
// https://taigaio.github.io/taiga-doc/dist/api.html#users-edit
func (s *UserService) Edit(user *User) (*User, error) {
	if err := requireNonNil("user", user); err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user ID is required")
	}
	patchPayload := map[string]any{}
	if user.Bio != "" {
		patchPayload["bio"] = user.Bio
	}
	if user.Color != "" {
		patchPayload["color"] = user.Color
	}
	if user.Email != "" {
		patchPayload["email"] = user.Email
	}
	if user.FullName != "" {
		patchPayload["full_name"] = user.FullName
	}
	if user.Lang != "" {
		patchPayload["lang"] = user.Lang
	}
	if user.ReadNewTerms {
		patchPayload["read_new_terms"] = user.ReadNewTerms
	}
	if user.Theme != "" {
		patchPayload["theme"] = user.Theme
	}
	if user.Timezone != "" {
		patchPayload["timezone"] = user.Timezone
	}
	if user.Username != "" {
		patchPayload["username"] = user.Username
	}
	if len(patchPayload) == 0 {
		return nil, errors.New("no updatable user fields were provided; use Patch for explicit zero-value updates")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(user.ID))
	var responseUser User
	_, err := s.client.Request.Patch(url, &patchPayload, &responseUser)
	if err != nil {
		return nil, err
	}
	return &responseUser, nil
}

// Patch sends an explicit PATCH payload to edit a user.
func (s *UserService) Patch(userID int, patch *UserPatch) (*User, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID))
	var u User
	_, err := s.client.Request.Patch(url, patch, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update is an alias for Edit.
func (s *UserService) Update(user *User) (*User, error) {
	return s.Edit(user)
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#users-delete
func (s *UserService) Delete(userID int) error {
	if err := requirePositiveID("userID", userID); err != nil {
		return err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID))
	_, err := s.client.Request.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

// GetContacts => https://taigaio.github.io/taiga-doc/dist/api.html#users-get-contacts
// func GetContacts(s *UserService) (User, error) {}

// CancelUserAccount => https://taigaio.github.io/taiga-doc/dist/api.html#users-cancel
// func CancelUserAccount(s *UserService) (User, error) {}

// ChangeAvatar => https://taigaio.github.io/taiga-doc/dist/api.html#users-change-avatar
// func ChangeAvatar(s *UserService) (User, error) {}

// RemoveAvatar => https://taigaio.github.io/taiga-doc/dist/api.html#users-remove-avatar
// func RemoveAvatar(s *UserService) (User, error) {}

// ChangeEmail => https://taigaio.github.io/taiga-doc/dist/api.html#users-change-email
// func ChangeEmail(s *UserService) (User, error) {}

// ChangePassword => https://taigaio.github.io/taiga-doc/dist/api.html#users-change-password
// func ChangePassword(s *UserService) (User, error) {}

// PasswordRecovery => https://taigaio.github.io/taiga-doc/dist/api.html#users-password-recovery
// func PasswordRecovery(s *UserService) (User, error) {}

// ChangePasswordFromRecovery => https://taigaio.github.io/taiga-doc/dist/api.html#users-change-password-from-recovery
// func ChangePasswordFromRecovery(s *UserService) (User, error) {}
