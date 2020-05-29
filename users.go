package taigo

import (
	"fmt"
	"strconv"

	"github.com/google/go-querystring/query"
)

// UserService is a handle to actions related to Users
//
// https://taigaio.github.io/taiga-doc/dist/api.html#users
type UserService struct {
	client   *Client
	Endpoint string
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#users-list
func (s *UserService) List(queryParams *UsersQueryParams) ([]User, error) {
	url := s.client.MakeURL(s.Endpoint)
	if queryParams != nil {
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	var users []User
	_, err := s.client.Request.Get(url, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#users-get
func (s *UserService) Get(userID int) (*User, error) {
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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "stats")
	var usd UserStatsDetail
	_, err := s.client.Request.Get(url, &usd)
	if err != nil {
		return nil, err
	}
	return &usd, nil
}

// GetWatchedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-watched
//
// TODO: Implement query param filtering
func (s *UserService) GetWatchedContent(userID int) (*UserWatched, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "watched")
	var uw UserWatched
	_, err := s.client.Request.Get(url, &uw)
	if err != nil {
		return nil, err
	}
	return &uw, nil
}

// GetLikedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-liked
//
// TODO: Implement query param filtering
func (s *UserService) GetLikedContent(userID int) (*UserLiked, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userID), "liked")
	var ul UserLiked
	_, err := s.client.Request.Get(url, &ul)
	if err != nil {
		return nil, err
	}
	return &ul, nil
}

// https://taigaio.github.io/taiga-doc/dist/api.html#users-voted
// func GetVotedContent(s *UserService) (User, error) {}

// Edit sends a PATCH request to edit a user
// https://taigaio.github.io/taiga-doc/dist/api.html#users-edit
func (s *UserService) Edit(user *User) (*User, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(user.ID))
	var u User
	_, err := s.client.Request.Patch(url, &user, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#users-delete
func (s *UserService) Delete(userID int) error {
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
