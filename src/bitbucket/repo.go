package bitbucket

// Repository is the repository reference of a Pull Request
type Repository struct {
	Slug        string          `json:"slug"`
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Project     Project         `json:"project"`
	Links       RepositoryLinks `json:"links"`
}

// RepositoryLinks are links to the repository
type RepositoryLinks struct {
	Self  []Link `json:"self"`
	Clone []Link `json:"clone"`
}

// Link is a link to a resource
type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// Project is the project reference of a Pull Request
type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Author is the creator of a Pull Request
type Author struct {
	User        User   `json:"user"`
	Role        string `json:"role"`
	HasApproved bool   `json:"approved"`
	Status      string `json:"status"`
}

// User is the user object for an Author
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
	Links       struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type Links struct {
	Self []struct {
		Href string `json:"href"`
	} `json:"self"`
}
