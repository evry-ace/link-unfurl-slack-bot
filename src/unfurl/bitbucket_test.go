package unfurl

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/evry-ace/link-unfurl-slack-bot/src/bitbucket"
	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"
)

const (
	testdataDir = "../../testdata"

	server  = "bitbucket.corp.org"
	project = "MY-PROJ"
	repo    = "my-repo"
	pr      = "297"
)

var (
	client = bitbucket.Client{
		Server: server,
		PAT:    "my-token",
	}
)

func TestBitbucketLinkType(t *testing.T) {
	typeLinks := map[string][]string{
		BitbucketURLPullRequestType: {
			"/projects/MY-PRO/repos/my-repo/pull-requests/123",
			"/projects/MY-PRO/repos/my-repo/pull-requests/123/overview",
			"/projects/MY-PRO/repos/my-repo/pull-requests/123/diff",
			"/projects/MY-PRO/repos/my-repo/pull-requests/123/commits",
		},
		BitbucketURLSourceCodeType: {
			"/projects/MY-PRO/repos/my-repo/browse/file",
			"/projects/MY-PRO/repos/my-repo/browse/file.ext",
			"/projects/MY-PRO/repos/my-repo/browse/some/file.ext",
		},
		BitbucketURLRepoType: {
			"/projects/MY-PRO/repos/my-repo/browse",
			"/projects/MY-PRO/repos/my-repo/commits",
			"/projects/MY-PRO/repos/my-repo/branches",
			"/projects/MY-PRO/repos/my-repo/settings",
		},
		BitbucketURLUnknownType: {
			"/dashboard",
			"/admin",
			"/profile",
			"/account",
		},
	}

	for shouldBeType, paths := range typeLinks {
		for _, path := range paths {
			u := url.URL{Path: path}
			isType, _ := bitbucketLinkType(&u)
			if isType != shouldBeType {
				t.Errorf("URL type should be %s but was %s for url %+v", shouldBeType, isType, u)
			}
		}
	}
}

func TestBitbucketLink(t *testing.T) {
	t.Run("PullRequest", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		linkPath := fmt.Sprintf(bitbucket.APIPaths["pullRequest"], project, repo, pr)
		linkUrl := fmt.Sprintf("https://%s/%s", server, linkPath)

		prJSON := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-pull-requests-297.json")
		prAPI := fmt.Sprintf(bitbucket.APIPaths["base"], server, linkPath)

		httpmock.RegisterResponder("GET", prAPI,
			httpmock.NewStringResponder(200, httpmock.File(prJSON).String()))

		statusJSON := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-build-status-654382.json")
		statusPath := fmt.Sprintf(bitbucket.StatusPaths["status"], "a68adcc6e8461db084acf7e76401d3c1542bb8ad")
		statusAPI := fmt.Sprintf(bitbucket.StatusPaths["base"], server, statusPath)

		httpmock.RegisterResponder("GET", statusAPI,
			httpmock.NewStringResponder(200, httpmock.File(statusJSON).String()))

		attachment, err := BitbucketLink(linkUrl, client)

		if err != nil {
			t.Errorf("Error should be nil but was %s", err)
		}

		assert.Equal(t, "#297 My new feature", attachment.Title)
		assert.Equal(t, "User D", attachment.AuthorName)
		assert.Equal(t, "My awesome description", attachment.Text)
		assert.Equal(t, 4, len(attachment.Fields))
	})
}
