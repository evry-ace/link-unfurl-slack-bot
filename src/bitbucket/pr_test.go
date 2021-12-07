package bitbucket

import (
	"testing"

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
