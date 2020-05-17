package taigo

import "time"

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
	ProjectBlockedCode  interface{}         `json:"project_blocked_code,omitempty"`
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

// User represents User detail | https://taigaio.github.io/taiga-doc/dist/api.html#object-user-detail
type User struct {
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
	ProjectBlockedCode  interface{}         `json:"project_blocked_code"`
	ProjectIsPrivate    bool                `json:"project_is_private"`
	ProjectName         string              `json:"project_name"`
	ProjectSlug         string              `json:"project_slug"`
	Ref                 int                 `json:"ref"`
	Slug                string              `json:"slug"`
	Status              interface{}         `json:"status"`
	StatusColor         interface{}         `json:"status_color"`
	Subject             string              `json:"subject"`
	TagsColors          []TagsColors        `json:"tags_colors,omitempty"`
	TotalFans           int                 `json:"total_fans"`
	TotalWatchers       int                 `json:"total_watchers"`
	Type                string              `json:"type"`
}

// UserAuthenticationDetail represents all details of User logging in
type UserAuthenticationDetail struct {
	ReadNewTerms                  bool      `json:"read_new_terms"`
	MaxPrivateProjects            int       `json:"max_private_projects"`
	Roles                         []string  `json:"roles"`
	Email                         string    `json:"email"`
	BigPhoto                      string    `json:"big_photo"`
	MaxMembershipsPrivateProjects int       `json:"max_memberships_private_projects"`
	Username                      string    `json:"username"`
	AuthToken                     string    `json:"auth_token"`
	TotalPublicProjects           int       `json:"total_public_projects"`
	DateJoined                    time.Time `json:"date_joined"`
	GravatarID                    string    `json:"gravatar_id"`
	FullName                      string    `json:"full_name"`
	AcceptedTerms                 bool      `json:"accepted_terms"`
	Timezone                      string    `json:"timezone"`
	IsActive                      bool      `json:"is_active"`
	Bio                           string    `json:"bio"`
	UUID                          string    `json:"uuid"`
	MaxMembershipsPublicProjects  int       `json:"max_memberships_public_projects"`
	TotalPrivateProjects          int       `json:"total_private_projects"`
	Photo                         string    `json:"photo"`
	FullNameDisplay               string    `json:"full_name_display"`
	Theme                         string    `json:"theme"`
	Color                         string    `json:"color"`
	MaxPublicProjects             int       `json:"max_public_projects"`
	Lang                          string    `json:"lang"`
	ID                            int       `json:"id"`
}
