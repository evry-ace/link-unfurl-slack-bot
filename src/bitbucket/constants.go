package bitbucket

var APIPaths = map[string]string{
	"base":         "https://%s/rest/api/1.0/%s",
	"repo":         "projects/%s/repos/%s",
	"repoCommits":  "projects/%s/repos/%s/commits",
	"browse":       "projects/%s/repos/%s/browse/%s",
	"pullRequests": "projects/%s/repos/%s/pull-requests",
	"pullRequest":  "projects/%s/repos/%s/pull-requests/%s",
}

var StatusPaths = map[string]string{
	"base":   "https://%s/rest/build-status/1.0/%s",
	"status": "commits/%s",
}

const (
	// PullRequestReviewStatusNeedsWork is the status for a pull request review
	// when the pull request needs more work before it can be merged.
	PullRequestReviewStatusNeedsWork = "NEEDS_WORK"

	// PullRequestReviewStatusApproved is the status for an approved pull request review
	PullRequestReviewStatusApproved = "APPROVED"

	// PullRequestReviewStatusUnapproved is the status for an unapproved pull request review
	PullRequestReviewStatusUnapproved = "UNAPPROVED"

	PullRequestUserRoleAuthor   = "AUTHOR"
	PullRequestUserRoleReviewer = "REVIEWER"
)
