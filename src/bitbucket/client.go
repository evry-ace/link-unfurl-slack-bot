package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client is a Bitbucket client that is used for storing server configuration,
// authentiation and to run the actual Bitbucket requests.
type Client struct {
	Server    string
	PAT       string
	timeout   int
	useragent string
}

// Timeout returns the configured connection timeout for the HTTP client.
func (c Client) Timeout() int {
	// @TODO find a better way for checking if timeout has been set or not.
	// c.timeout is 0 by default which means no timeout.
	if c.timeout == 0 {
		return 2
	}

	return c.timeout
}

// Useragent returns the configured client useragnt or a default one.
func (c Client) Useragent() string {
	if c.useragent == "" {
		return "bitbucket-go-sdk"
	}

	return c.useragent
}

// RawRequest does a API request and returns content as a string. This is just a
// helper method used by other Client functions.
// https://docs.atlassian.com/bitbucket-server/rest/5.16.0/bitbucket-rest.html
func (c Client) RawRequest(url string) ([]byte, int, error) {
	bearer := fmt.Sprintf("Bearer %s", c.PAT)

	httpClient := http.Client{
		Timeout: time.Second * time.Duration(c.Timeout()),
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return []byte{}, 0, reqErr
	}

	req.Header.Set("User-Agent", c.Useragent())
	req.Header.Add("Authorization", bearer)

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return []byte{}, 0, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return []byte{}, 0, readErr
	}

	return body, res.StatusCode, nil
}

// rawUrl returns a URL for a given API path and a list of path parameters.
func (c Client) rawUrl(apis map[string]string, path string, args ...interface{}) string {
	var url string

	if len(args) > 0 {
		url = fmt.Sprintf(apis[path], args...)
	} else {
		url = apis[path]
	}

	return fmt.Sprintf(apis["base"], c.Server, url)
}

// PullRequests returns a list of PullRequests for a given repo in a given
// project.
// @TODO add support for pagination when there are lots of pull requests?
// @TODO add support for ofset when there are lots of pull requests
func (c Client) PullRequests(project string, repo string) (PullRequests, error) {
	var prs PullRequests

	data, status, err := c.RawRequest(c.rawUrl(APIPaths, "pullRequests", project, repo))
	if err != nil {
		return prs, err
	}

	if status != 200 {
		return prs, fmt.Errorf("HTTP request failed with unexpected status code %d", status)
	}

	if err := json.Unmarshal(data, &prs); err != nil {
		return prs, err
	}

	return prs, nil
}

// PullRequest returns a single PullRequest in a given repo in a given project.
func (c Client) PullRequest(project string, repo string, id int) (PullRequest, error) {
	var pr PullRequest

	data, status, err := c.RawRequest(c.rawUrl(APIPaths, "pullRequest", project, repo, fmt.Sprint(id)))
	if err != nil {
		return pr, err
	}

	if status != 200 {
		return pr, fmt.Errorf("HTTP request failed with unexpected status code %d", status)
	}

	if err := json.Unmarshal(data, &pr); err != nil {
		return pr, err
	}

	return pr, nil
}

// Repository returns a single Repo in a given project.
func (c Client) Repository(project string, repo string) (Repository, error) {
	var repository Repository

	data, status, err := c.RawRequest(c.rawUrl(APIPaths, "repo", project, repo))
	if err != nil {
		return repository, err
	}

	if status != 200 {
		return repository, fmt.Errorf("HTTP request failed with unexpected status code %d", status)
	}

	if err := json.Unmarshal(data, &repository); err != nil {
		return repository, err
	}

	return repository, nil
}

// Commits returns a list of Commits for a given repo in a given project
func (c Client) Commits(project string, repo string, co CommitOptions) (CommitList, error) {
	var commits CommitList

	url := c.rawUrl(APIPaths, "repoCommits", project, repo)
	query := co.ToQueryString()

	data, status, err := c.RawRequest(fmt.Sprintf("%s?%s", url, query))
	if err != nil {
		return commits, err
	}

	if status != 200 {
		return commits, fmt.Errorf("HTTP request failed with unexpected status code %d", status)
	}

	if err := json.Unmarshal(data, &commits); err != nil {
		return commits, err
	}

	return commits, nil
}

// Status returns a status for a given commit
func (c Client) Status(sha string) (StatusList, error) {
	var s StatusList

	data, status, err := c.RawRequest(c.rawUrl(StatusPaths, "status", sha))
	if err != nil {
		return s, err
	}

	if status != 200 {
		return s, fmt.Errorf("HTTP request failed with unexpected status code %d", status)
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return s, err
	}

	return s, nil
}
