package bitbucket

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestPullRequestReviewBy(t *testing.T) {
	t.Run("no reviews when there are no reviewers", func(t *testing.T) {
		pr := PullRequest{}
		assert.Equal(t, "No reviews :sob:", pr.ReviewedBy())
	})

	t.Run("ignors reviewers without review", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusUnapproved,
				},
			},
		}
		assert.Equal(t, "No reviews :sob:", pr.ReviewedBy())
	})

	t.Run("review by single reviewer", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusApproved,
					User: User{
						DisplayName: "Test User",
					},
				},
			},
		}
		assert.Equal(t, "Test User (APPROVED)", pr.ReviewedBy())
	})

	t.Run("reviews by multiple reviewers", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusApproved,
					User: User{
						DisplayName: "Test User",
					},
				},
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusUnapproved,
				},
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusNeedsWork,
					User: User{
						DisplayName: "Test User 2",
					},
				},
			},
		}
		assert.Equal(t, "Test User (APPROVED), Test User 2 (NEEDS_WORK)", pr.ReviewedBy())
	})
}

func TestPullRequestRepoSlug(t *testing.T) {
	t.Run("returns repo slug", func(t *testing.T) {
		pr := PullRequest{
			ID: 123,
			ToRef: GitRef{
				Repository: Repository{
					Slug: "test-repo",
					Project: Project{
						Key: "test-project",
					},
				},
			},
		}
		assert.Equal(t, "test-project/test-repo#123", pr.RepoSlug())
	})
}

func TestPullRequestString(t *testing.T) {
	t.Run("returns string representation", func(t *testing.T) {
		pr := PullRequest{
			ID:    123,
			Title: "Test PR",
			State: "OPEN",
			Author: Author{
				User: User{
					DisplayName: "Test User",
				},
			},
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusApproved,
				},
			},
			ToRef: GitRef{
				Repository: Repository{
					Slug: "test-repo",
					Project: Project{
						Key: "test-project",
					},
				},
			},
			CreatedDate: time.Now().UnixMicro(),
		}
		assert.Equal(t, "_Test PR_ (:white_check_mark: Approved) by Test User opened in about a second", pr.String())
	})
}

func TestPullRequestIsWorkInProgress(t *testing.T) {
	t.Run("returns true when pull request is work in progress", func(t *testing.T) {
		pr := PullRequest{
			Title: "WIP: Test PR",
		}
		assert.Equal(t, true, pr.IsWorkInProgress())
	})

	t.Run("returns false when pull request is not work in progress", func(t *testing.T) {
		pr := PullRequest{
			Title: "Test PR",
		}
		assert.Equal(t, false, pr.IsWorkInProgress())
	})
}

func TestPullRequestIsApproved(t *testing.T) {
	t.Run("returns true when pull request is approved", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusApproved,
				},
			},
		}
		assert.Equal(t, true, pr.IsApproved())
	})

	t.Run("returns false when pull request is not approved", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusUnapproved,
				},
			},
		}
		assert.Equal(t, false, pr.IsApproved())
	})
}

func TestPullRequestAppovalStatus(t *testing.T) {
	t.Run("returns approved when pull request is approved", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusApproved,
				},
			},
		}
		assert.Equal(t, "Approved", pr.ApprovalStatus(false))
	})

	t.Run("returns unapproved when pull request is not approved", func(t *testing.T) {
		pr := PullRequest{
			Reviewers: []Author{
				{
					Role:   PullRequestUserRoleReviewer,
					Status: PullRequestReviewStatusUnapproved,
				},
			},
		}
		assert.Equal(t, "Unapproved", pr.ApprovalStatus(false))
	})
}
