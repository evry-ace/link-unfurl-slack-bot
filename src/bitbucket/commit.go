package bitbucket

import (
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/xeonx/timeago"
)

// CommitList is a list of commits
type CommitList struct {
	Size          int      `json:"size"`
	IsLastPage    bool     `json:"isLastPage"`
	Start         int      `json:"start"`
	Limit         int      `json:"limit"`
	NextPageStart int      `json:"nextPageStart"`
	Values        []Commit `json:"values"`
}

// CommitOptions is the options for a commit
type CommitOptions struct {
	FollowRenames bool `url:"followRenames,omitempty"`
	IngoreMissing bool `url:"ignoreMissing,omitempty"`

	// if present, controls how merge commits should be filtered. Can be either exclude,
	// to exclude merge commits, include, to include both merge commits and non-merge
	// commits or only, to only return merge commits.
	Merges string `url:"merges,omitempty"`

	Path       string `url:"path,omitempty"`
	Since      string `url:"since,omitempty"`
	Until      string `url:"until,omitempty"`
	WithCounts bool   `url:"withCounts,omitempty"`

	Limit int `url:"limit,omitempty"`
	Start int `url:"start,omitempty"`
}

// ToQueryString returns the CommitOptions as a query string
func (co CommitOptions) ToQueryString() string {
	v, _ := query.Values(co)
	return v.Encode()
}

// Commit is a commit in a branch
type Commit struct {
	ID                 string `json:"id"`
	DisplayID          string `json:"displayId"`
	Author             User   `json:"author"`
	AuthorTimestamp    int64  `json:"authorTimestamp"`
	Committer          User   `json:"committer"`
	CommitterTimestamp int64  `json:"committerTimestamp"`
	Message            string `json:"message"`
	Parents            []struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
	} `json:"parents"`
	Properties struct {
		JIRAIssueKeys []string `json:"jira-key"`
	} `json:"properties"`
}

// JIRAIssueKeys returns the JIRA issue key from the commit
func (c Commit) JIRAIssueKeys() []string {
	return c.Properties.JIRAIssueKeys
}

// String returns the Commit as a string
func (c Commit) String() string {
	return fmt.Sprintf(
		"%s (%s) by %s %s",
		c.Message,
		c.ID,
		c.Author.DisplayName,
		c.TimeAgo(),
	)
}

// TimeAgo returns the time since the commit was made in human readable format
func (c Commit) TimeAgo() string {
	t := time.Unix(c.AuthorTimestamp/1000, 0)
	s := timeago.NoMax(timeago.English).Format(t)

	return s
}
