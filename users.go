package taigo

import (
	"fmt"
)

const endpointUsers = "/users"

// UserService is a handle to actions related to Users
//
// https://taigaio.github.io/taiga-doc/dist/api.html#users
type UserService struct {
	client *Client
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#users-list
func (s *UserService) List() ([]User, error) {
	url := s.client.APIURL + endpointUsers
	var users []User

	err := s.client.Request.GetRequest(url, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#users-get
func (s *UserService) Get(user User) (*User, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUsers, user.ID)
	var respUser User
	err := s.client.Request.GetRequest(url, &respUser)
	if err != nil {
		return nil, err
	}
	return &respUser, nil
}

// Me => https://taigaio.github.io/taiga-doc/dist/api.html#users-me
func (s *UserService) Me() (*User, error) {
	var user User
	url := s.client.APIURL + endpointUsers + "/me"

	err := s.client.Request.GetRequest(url, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetStats => https://taigaio.github.io/taiga-doc/dist/api.html#users-stats
func (s *UserService) GetStats(user *User) (*UserStatsDetail, error) {
	url := s.client.APIURL + endpointUsers + fmt.Sprintf("/%d/stats", user.ID)
	var userStatsDetail UserStatsDetail

	err := s.client.Request.GetRequest(url, &userStatsDetail)
	if err != nil {
		return nil, err
	}

	return &userStatsDetail, nil
}

// GetWatchedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-watched
//
// TODO: Implement query param filtering
func (s *UserService) GetWatchedContent(user *User) (*UserWatched, error) {
	url := s.client.APIURL + endpointUsers + fmt.Sprintf("/%d/watched", user.ID)
	var userWatched UserWatched

	err := s.client.Request.GetRequest(url, &userWatched)
	if err != nil {
		return nil, err
	}

	return &userWatched, nil
}

// GetLikedContent => https://taigaio.github.io/taiga-doc/dist/api.html#users-liked
//
// TODO: Implement query param filtering
func (s *UserService) GetLikedContent(user *User) (*UserLiked, error) {
	url := s.client.APIURL + endpointUsers + fmt.Sprintf("/%d/liked", user.ID)
	var userLiked UserLiked

	err := s.client.Request.GetRequest(url, &userLiked)
	if err != nil {
		return nil, err
	}

	return &userLiked, nil
}

// https://taigaio.github.io/taiga-doc/dist/api.html#users-voted
// func GetVotedContent(s *UserService) (User, error) {}

// Edit sends a PATCH request to edit a user
// https://taigaio.github.io/taiga-doc/dist/api.html#users-edit
func (s *UserService) Edit(userID int, patchedUser *User) (*User, error) {
	url := s.client.APIURL + endpointUsers + fmt.Sprintf("/%d", userID)
	var responseUser User

	err := s.client.Request.PatchRequest(url, &patchedUser, &responseUser)
	if err != nil {
		return nil, err
	}
	return &responseUser, nil
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#users-delete
func (s *UserService) Delete(user *User) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUsers, user.ID)
	err := s.client.Request.DeleteRequest(url)
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
