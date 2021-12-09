package unfurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/evry-ace/link-unfurl-slack-bot/src/bitbucket"
	"github.com/slack-go/slack"
)

const (
	BitbucketIcon               = "https://avatars.slack-edge.com/2021-06-20/2187759053413_fb4aad0a769aaadbdc62_72.png"
	BitbucketURLPullRequestType = "pull_request"
	BitbucketURLRepoType        = "repo"
	BitbucketURLSourceCodeType  = "source_code"
	BitbucketURLUnknownType     = "unknown"
)

// bitbucketLinkType returns the type of Bitbucket link and the matches
func bitbucketLinkType(u *url.URL) (string, []string) {
	var isPullRequest = regexp.MustCompile(
		"^/" + fmt.Sprintf(bitbucket.APIPaths["pullRequest"],
			"([^/]+)", "([^/]+)", "([^/]+)",
		),
	)

	var isSourceCode = regexp.MustCompile(
		"^/" + fmt.Sprintf(bitbucket.APIPaths["browse"],
			"([^/]+)", "([^/]+)", "(.+)",
		),
	)

	var isRepo = regexp.MustCompile(
		"^/" + fmt.Sprintf(bitbucket.APIPaths["repo"], "([^/]+)", "([^/]+)"),
	)

	if isPullRequest.MatchString(u.Path) {
		return BitbucketURLPullRequestType, isPullRequest.FindStringSubmatch(u.Path)
	} else if isSourceCode.MatchString(u.Path) {
		return BitbucketURLSourceCodeType, isSourceCode.FindStringSubmatch(u.Path)
	} else if isRepo.MatchString(u.Path) {
		return BitbucketURLRepoType, isRepo.FindStringSubmatch(u.Path)
	}

	return BitbucketURLUnknownType, []string{}
}

// BitbucketLink returns a Slack Attachment for Bitbucket links
func BitbucketLink(URL string, c bitbucket.Client) (slack.Attachment, error) {
	attachement := slack.Attachment{}

	u, err := url.Parse(URL)
	if err != nil {
		return slack.Attachment{}, err
	}

	// Parse what type of link this is
	linkType, matches := bitbucketLinkType(u)
	fmt.Printf("linkType=%s matches=%s", linkType, matches)
	switch linkType {
	case BitbucketURLPullRequestType:
		proj := matches[1]
		repo := matches[2]
		prid, err := strconv.Atoi(matches[3])
		if err != nil {
			return attachement, err
		}

		fmt.Printf("project=%s, repo=%s, prid=%d", proj, repo, prid)

		return bitbucketPRLink(proj, repo, prid, c)

	case BitbucketURLSourceCodeType:
		// @TODO
		fmt.Println("BitbucketURLSourceCodeType is not implemented for BitbucketLink()")
	case BitbucketURLRepoType:
		// @TODO
		fmt.Println("BitbucketURLRepoType is not implemented for BitbucketLink()")
	default:
		return slack.Attachment{}, errors.New("bitbucket link not supported")
	}

	return attachement, nil
}

// bitbucketPRLink returns a Slack Attachment for Bitbucket Pull Request links
func bitbucketPRLink(proj string, repo string, prid int, c bitbucket.Client) (slack.Attachment, error) {
	attachement := slack.Attachment{}

	// Get the Pull Request
	pr, err := c.PullRequest(proj, repo, prid)
	if err != nil {
		return slack.Attachment{}, err
	}

	// Get the Pull Request Status
	st, err := c.Status(pr.FromRef.LatestCommit)
	if err != nil {
		return slack.Attachment{}, err
	}

	attachement.Ts = json.Number(fmt.Sprint(pr.CreatedDate))
	attachement.FooterIcon = BitbucketIcon
	attachement.Footer = "Bitbucket"
	attachement.AuthorID = fmt.Sprintf("%d", pr.Author.User.ID)
	attachement.AuthorName = pr.Author.User.DisplayName
	attachement.AuthorLink = pr.Author.User.Links.Self[0].Href
	attachement.Title = fmt.Sprintf("#%d %s", pr.ID, pr.Title)
	attachement.TitleLink = pr.Links.Self[0].Href
	attachement.Text = pr.Description
	attachement.Fields = []slack.AttachmentField{
		{
			Title: "PR State",
			Value: pr.State,
			Short: true,
		},
		{
			Title: "Build Status",
			Value: st.State(),
			Short: true,
		},
		{
			Title: "Reviewers",
			Value: pr.ReviewedBy(),
			Short: true,
		},
		{
			Title: "Review Status",
			Value: pr.ApprovalStatus(true),
			Short: true,
		},
	}

	return attachement, nil
}

// BitbucketRepoLink returns a Slack Attachment for a Bitbucket Repo links
func BitbucketRepoLink(project string, repo string, c bitbucket.Client) (slack.Attachment, error) {
	var attachement slack.Attachment

	// Get repo info
	r, err := c.Repository(project, repo)
	if err != nil {
		return attachement, err
	}

	// Get repo commits
	co, err := c.Commits(project, repo, bitbucket.CommitOptions{})
	if err != nil {
		return attachement, err
	}

	// Get build status for latest commit
	st, err := c.Status(co.Values[0].ID)
	if err != nil {
		return attachement, err
	}

	attachement.Title = r.Name
	attachement.TitleLink = r.Links.Self[0].Href
	attachement.Text = r.Description
	attachement.Fields = []slack.AttachmentField{
		{
			Title: "Last Commit",
			Value: co.Values[0].String(),
			Short: true,
		},
		{
			Title: "Build Status",
			Value: st.State(),
			Short: true,
		},
	}

	return attachement, nil
}
