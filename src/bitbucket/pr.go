package bitbucket

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/xeonx/timeago"
)

// PullRequests is a list of Pull Requests
type PullRequests struct {
	Size       int           `json:"size"`
	Limit      int           `json:"limit"`
	IsLastPage bool          `json:"isLastPage"`
	List       []PullRequest `json:"values"`
}

// PullRequest is a single Pull Request
type PullRequest struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state"`
	IsOpen      bool     `json:"open"`
	IsClosed    bool     `json:"closed"`
	CreatedDate int64    `json:"createdDate"`
	UpdatedDate int64    `json:"updatedDate"`
	FromRef     GitRef   `json:"fromRef"`
	ToRef       GitRef   `json:"toRef"`
	Author      Author   `json:"author"`
	Reviewers   []Author `json:"reviewers"`
	Links       struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

// ReviewBy returns a list of reviewers as a string
func (pr PullRequest) ReviewBy() string {
	var reviewers []string
	for _, reviewer := range pr.Reviewers {
		if reviewer.Role == PullRequestUserRoleReviewer && reviewer.Status != PullRequestReviewStatusUnapproved {
			reviewers = append(reviewers, fmt.Sprintf("%s (%s)", reviewer.User.DisplayName, reviewer.Status))
		}
	}

	if len(reviewers) == 0 {
		return "No reviews :sob:"
	}

	return strings.Join(reviewers, ", ")
}

// OpenSince returns how long the Pull Request has been open in human readable
// format.
func (pr PullRequest) OpenSince() string {
	t := time.Unix(pr.CreatedDate/1000, 0)
	s := timeago.NoMax(timeago.English).Format(t)

	return s
}

// RepoSlug returns the combind string of repo and project
func (pr PullRequest) RepoSlug() string {
	return fmt.Sprintf(
		"%s/%s#%d",
		pr.ToRef.Repository.Slug,
		pr.ToRef.Repository.Project.Key,
		pr.ID,
	)
}

// ToString returns the Pull Request as a string
func (pr PullRequest) ToString() string {
	return fmt.Sprintf(
		"_%s_ (%s) by %s opened %s",
		pr.Title,
		pr.ApprovalStatus(true),
		pr.Author.User.DisplayName,
		pr.OpenSince(),
	)
}

// IsWorkInProgress returns true if the Pull Request is Work In Progress
func (pr PullRequest) IsWorkInProgress() bool {
	var wipTitle = regexp.MustCompile("WIP")

	return wipTitle.MatchString(pr.Title)
}

// IsApproved return if the pull request is approved by reviewers
func (pr PullRequest) IsApproved() bool {
	approved := false
	for _, u := range pr.Reviewers {
		switch status := u.Status; status {
		// case "UNAPPROVED"
		case "APPROVED":
			approved = true
		case "NEEDS_WORK":
			return false
		}
	}

	return approved
}

// ApprovalStatus returns a human readable approval status
func (pr PullRequest) ApprovalStatus(showEmojis bool) string {
	if pr.IsApproved() {
		if showEmojis {
			return ":white_check_mark: Approved"
		}

		return "Approved"
	}

	if showEmojis {
		return ":disappointed: Unapproved"
	}

	return "Unapproved"
}

// GitRef is the git reference of a Pull Request
type GitRef struct {
	ID           string     `json:"id"`
	DisplayID    string     `json:"displayId"`
	LatestCommit string     `json:"latestCommit"`
	Repository   Repository `json:"repository"`
}
