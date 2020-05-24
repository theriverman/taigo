package taigo

import "time"

// UsersQueryParams holds fields to be used as URL query parameters to filter the queried objects
type UsersQueryParams struct {
	Project int `url:"project,omitempty"`
}

// User represents User detail | https://taigaio.github.io/taiga-doc/dist/api.html#object-user-detail
type User struct {
	authToken                     string    // internal; non-standard
	AcceptedTerms                 bool      `json:"accepted_terms,omitempty"`
	BigPhoto                      string    `json:"big_photo,omitempty"`
	Bio                           string    `json:"bio,omitempty"`
	Color                         string    `json:"color,omitempty"`
	DateJoined                    time.Time `json:"date_joined,omitempty"`
	Email                         string    `json:"email,omitempty"`
	FullName                      string    `json:"full_name,omitempty"`
	FullNameDisplay               string    `json:"full_name_display,omitempty"`
	GravatarID                    string    `json:"gravatar_id,omitempty"`
	ID                            int       `json:"id,omitempty"`
	IsActive                      bool      `json:"is_active,omitempty"`
	Lang                          string    `json:"lang,omitempty"`
	MaxMembershipsPrivateProjects int       `json:"max_memberships_private_projects,omitempty"`
	MaxMembershipsPublicProjects  int       `json:"max_memberships_public_projects,omitempty"`
	MaxPrivateProjects            int       `json:"max_private_projects,omitempty"`
	MaxPublicProjects             int       `json:"max_public_projects,omitempty"`
	Photo                         string    `json:"photo,omitempty"`
	ReadNewTerms                  bool      `json:"read_new_terms,omitempty"`
	Roles                         []string  `json:"roles,omitempty"`
	Theme                         string    `json:"theme,omitempty"`
	Timezone                      string    `json:"timezone,omitempty"`
	TotalPrivateProjects          int       `json:"total_private_projects,omitempty"`
	TotalPublicProjects           int       `json:"total_public_projects,omitempty"`
	Username                      string    `json:"username,omitempty"`
	UUID                          string    `json:"uuid,omitempty"`
}

// GetToken returns teh token string embedded into User
func (u User) GetToken() string {
	return u.authToken
}

// UserAuthenticationDetail is a superset of User extended by an AuthToken field
type UserAuthenticationDetail struct {
	AuthToken string `json:"auth_token"`
	User             // Embedding type User struct
}

// AsUser returns a *User from *UserAuthenticationDetail
//
// The AuthToken can be accessed from *User via .
func (u *UserAuthenticationDetail) AsUser() *User {
	user := &User{}
	err := convertStructViaJSON(u, user)
	if err != nil {
		return nil
	}
	user.authToken = u.AuthToken
	return user
}

// Liked represents Liked | https://taigaio.github.io/taiga-doc/dist/api.html#object-liked-detail
type Liked struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	Description         string              `json:"description,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsFan               bool                `json:"is_fan,omitempty"`
	IsPrivate           bool                `json:"is_private,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	LogoSmallURL        string              `json:"logo_small_url,omitempty"`
	Name                string              `json:"name,omitempty"`
	Project             string              `json:"project,omitempty"`
	ProjectBlockedCode  string              `json:"project_blocked_code,omitempty"`
	ProjectIsPrivate    bool                `json:"project_is_private,omitempty"`
	ProjectName         string              `json:"project_name,omitempty"`
	ProjectSlug         string              `json:"project_slug,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Slug                string              `json:"slug,omitempty"`
	Status              int                 `json:"status,omitempty"`
	StatusColor         string              `json:"status_color,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalFans           int                 `json:"total_fans,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	Type                string              `json:"type,omitempty"`
}

// Voted represents Voted | https://taigaio.github.io/taiga-doc/dist/api.html#object-voted-detail
type Voted struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	Description         string              `json:"description,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsPrivate           bool                `json:"is_private,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	LogoSmallURL        string              `json:"logo_small_url,omitempty"`
	Name                string              `json:"name,omitempty"`
	Project             int                 `json:"project,omitempty"`
	ProjectBlockedCode  string              `json:"project_blocked_code,omitempty"`
	ProjectIsPrivate    bool                `json:"project_is_private,omitempty"`
	ProjectName         string              `json:"project_name,omitempty"`
	ProjectSlug         string              `json:"project_slug,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Slug                string              `json:"slug,omitempty"`
	Status              string              `json:"status,omitempty"`
	StatusColor         string              `json:"status_color,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	Type                string              `json:"type,omitempty"`
}

// UserStatsDetail represents User stats detail | https://taigaio.github.io/taiga-doc/dist/api.html#object-user-stats-detail
type UserStatsDetail struct {
	Roles                     []string `json:"roles,omitempty"`
	TotalNumClosedUserstories int      `json:"total_num_closed_userstories,omitempty"`
	TotalNumContacts          int      `json:"total_num_contacts,omitempty"`
	TotalNumProjects          int      `json:"total_num_projects,omitempty"`
}

// UserContactDetail represents User contact detail | https://taigaio.github.io/taiga-doc/dist/api.html#object-contact-detail
type UserContactDetail struct {
	BigPhoto        string   `json:"big_photo,omitempty"`
	Bio             string   `json:"bio,omitempty"`
	Color           string   `json:"color,omitempty"`
	FullName        string   `json:"full_name,omitempty"`
	FullNameDisplay string   `json:"full_name_display,omitempty"`
	GravatarID      string   `json:"gravatar_id,omitempty"`
	ID              int      `json:"id,omitempty"`
	IsActive        bool     `json:"is_active,omitempty"`
	Lang            string   `json:"lang,omitempty"`
	Photo           string   `json:"photo,omitempty"`
	Roles           []string `json:"roles,omitempty"`
	Theme           string   `json:"theme,omitempty"`
	Timezone        string   `json:"timezone,omitempty"`
	Username        string   `json:"username,omitempty"`
}

// Watched represents Watched | https://taigaio.github.io/taiga-doc/dist/api.html#object-watched-detail
type Watched struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	Description         string              `json:"description,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsPrivate           bool                `json:"is_private,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	LogoSmallURL        string              `json:"logo_small_url,omitempty"`
	Name                string              `json:"name,omitempty"`
	Project             int                 `json:"project,omitempty"`
	ProjectBlockedCode  string              `json:"project_blocked_code,omitempty"`
	ProjectIsPrivate    bool                `json:"project_is_private,omitempty"`
	ProjectName         string              `json:"project_name,omitempty"`
	ProjectSlug         string              `json:"project_slug,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Slug                string              `json:"slug,omitempty"`
	Status              string              `json:"status,omitempty"`
	StatusColor         string              `json:"status_color,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	Type                string              `json:"type,omitempty"`
}

// UserWatched => https://taigaio.github.io/taiga-doc/dist/api.html#object-watched-detail
type UserWatched struct {
	AssignedTo          int                 `json:"assigned_to"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info"`
	CreatedDate         time.Time           `json:"created_date"`
	Description         string              `json:"description"`
	ID                  int                 `json:"id"`
	IsPrivate           bool                `json:"is_private"`
	IsVoter             bool                `json:"is_voter"`
	IsWatcher           bool                `json:"is_watcher"`
	LogoSmallURL        string              `json:"logo_small_url"`
	Name                string              `json:"name"`
	Project             int                 `json:"project"`
	ProjectBlockedCode  string              `json:"project_blocked_code"`
	ProjectIsPrivate    bool                `json:"project_is_private"`
	ProjectName         string              `json:"project_name"`
	ProjectSlug         string              `json:"project_slug"`
	Ref                 int                 `json:"ref"`
	Slug                string              `json:"slug"`
	Status              string              `json:"status"`
	StatusColor         string              `json:"status_color"`
	Subject             string              `json:"subject"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalVoters         int                 `json:"total_voters"`
	TotalWatchers       int                 `json:"total_watchers"`
	Type                string              `json:"type"`
}

// UserLiked => https://taigaio.github.io/taiga-doc/dist/api.html#object-liked-detail
type UserLiked struct {
	AssignedTo          int                 `json:"assigned_to"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info"`
	CreatedDate         time.Time           `json:"created_date"`
	Description         string              `json:"description"`
	ID                  int                 `json:"id"`
	IsFan               bool                `json:"is_fan"`
	IsPrivate           bool                `json:"is_private"`
	IsWatcher           bool                `json:"is_watcher"`
	LogoSmallURL        string              `json:"logo_small_url"`
	Name                string              `json:"name"`
	Project             int                 `json:"project"`
	ProjectBlockedCode  string              `json:"project_blocked_code"`
	ProjectIsPrivate    bool                `json:"project_is_private"`
	ProjectName         string              `json:"project_name"`
	ProjectSlug         string              `json:"project_slug"`
	Ref                 int                 `json:"ref"`
	Slug                string              `json:"slug"`
	Status              int                 `json:"status"`
	StatusColor         string              `json:"status_color"`
	Subject             string              `json:"subject"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalFans           int                 `json:"total_fans"`
	TotalWatchers       int                 `json:"total_watchers"`
	Type                string              `json:"type"`
}
