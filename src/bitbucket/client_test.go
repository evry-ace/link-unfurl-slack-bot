package bitbucket

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"
)

const (
	bitbucketServer    = "bitbucket.corp.org"
	bitbucketProject   = "MY-PROJ"
	bitbucketRepo      = "my-repo"
	bitbucketPAT       = "my-fake-token"
	bitbucketRRID      = 297
	bitbucketCommitSHA = "65438285b7e2433f7857e53603089d6a63ebf0bb"

	testdataDir = "../../testdata"
)

func TestBitbucketClientPullRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	prsJSONFile := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-pull-requests.json")
	prsReqPath := fmt.Sprintf(
		APIPaths["base"],
		bitbucketServer,
		fmt.Sprintf(APIPaths["pullRequests"], bitbucketProject, bitbucketRepo),
	)

	// Set up mock Bitbucket Server
	httpmock.RegisterResponder("GET", prsReqPath,
		httpmock.NewStringResponder(200, httpmock.File(prsJSONFile).String()))

	// Get Pull Requests using Client
	client := Client{Server: bitbucketServer, PAT: bitbucketPAT}
	prs, err := client.PullRequests(bitbucketProject, bitbucketRepo)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, prs.Size, 6)
	assert.Equal(t, len(prs.List), 6)
}

func TestBitbucketClientPullRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	prJSONFile := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-pull-requests-297.json")
	prReqPath := fmt.Sprintf(
		APIPaths["base"],
		bitbucketServer,
		fmt.Sprintf(
			APIPaths["pullRequest"],
			bitbucketProject,
			bitbucketRepo,
			fmt.Sprint(bitbucketRRID),
		),
	)

	// Set up mock Bitbucket Server
	httpmock.RegisterResponder("GET", prReqPath,
		httpmock.NewStringResponder(200, httpmock.File(prJSONFile).String()))

	// Get Pull Requests using Client
	client := Client{Server: bitbucketServer, PAT: bitbucketPAT}
	pr, err := client.PullRequest(bitbucketProject, bitbucketRepo, bitbucketRRID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, pr.ID, bitbucketRRID)
	assert.Equal(t, pr.Title, "My new feature")
}

func TestBitbucketClientRepository(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	repoJSONFile := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-repo.json")
	repoReqPath := fmt.Sprintf(
		APIPaths["base"],
		bitbucketServer,
		fmt.Sprintf(APIPaths["repo"], bitbucketProject, bitbucketRepo),
	)

	// Set up mock Bitbucket Server
	httpmock.RegisterResponder("GET", repoReqPath,
		httpmock.NewStringResponder(200, httpmock.File(repoJSONFile).String()))

	// Get Repository using Client
	client := Client{Server: bitbucketServer, PAT: bitbucketPAT}
	repo, err := client.Repository(bitbucketProject, bitbucketRepo)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, repo.Name, bitbucketRepo)
}

func TestBitbucketClientCommits(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	jsonFilePath := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-repo-commits.json")
	reqPath := fmt.Sprintf(
		APIPaths["base"],
		bitbucketServer,
		fmt.Sprintf(APIPaths["repoCommits"], bitbucketProject, bitbucketRepo),
	)

	// Set up mock Bitbucket Server
	httpmock.RegisterResponder("GET", reqPath,
		httpmock.NewStringResponder(200, httpmock.File(jsonFilePath).String()))

	// Get Pull Requests using Client
	client := Client{Server: bitbucketServer, PAT: bitbucketPAT}
	c, err := client.Commits(bitbucketProject, bitbucketRepo, CommitOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, c.Size, 25)
	assert.Equal(t, len(c.Values), 25)
}

func TestBitbucketClientStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	repoJSONFile := fmt.Sprintf("%s/%s", testdataDir, "bitbucket-build-status-654382.json")
	repoReqPath := fmt.Sprintf(
		StatusPaths["base"],
		bitbucketServer,
		fmt.Sprintf(StatusPaths["status"], bitbucketCommitSHA),
	)

	// Set up mock Bitbucket Server
	httpmock.RegisterResponder("GET", repoReqPath,
		httpmock.NewStringResponder(200, httpmock.File(repoJSONFile).String()))

	// Get Repository using Client
	client := Client{Server: bitbucketServer, PAT: bitbucketPAT}
	status, err := client.Status(bitbucketCommitSHA)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(status.Values))
	assert.Equal(t, StatusInProgress, status.Values[0].State)
}
