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
func (u *Unfurl) bitbucketLinkType(url *url.URL) (string, []string) {
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

	if isPullRequest.MatchString(url.Path) {
		return BitbucketURLPullRequestType, isPullRequest.FindStringSubmatch(url.Path)
	} else if isSourceCode.MatchString(url.Path) {
		return BitbucketURLSourceCodeType, isSourceCode.FindStringSubmatch(url.Path)
	} else if isRepo.MatchString(url.Path) {
		return BitbucketURLRepoType, isRepo.FindStringSubmatch(url.Path)
	}

	return BitbucketURLUnknownType, []string{}
}

// bitbucketLink returns a Slack Attachment for Bitbucket links
func (u *Unfurl) bitbucketLink(URL *url.URL) (slack.Attachment, error) {
	attachement := slack.Attachment{}

	// Parse what type of link this is
	linkType, matches := u.bitbucketLinkType(URL)
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
		return u.bitbucketPRLink(proj, repo, prid)

	case BitbucketURLSourceCodeType:
		// @TODO
		fmt.Println("BitbucketURLSourceCodeType is not implemented for BitbucketLink()")

	case BitbucketURLRepoType:
		proj := matches[1]
		repo := matches[2]

		fmt.Printf("project=%s, repo=%s", proj, repo)
		return u.bitbucketRepoLink(proj, repo)

	default:
		return slack.Attachment{}, errors.New("bitbucket link not supported")
	}

	return attachement, nil
}

// bitbucketPRLink returns a Slack Attachment for Bitbucket Pull Request links
func (u *Unfurl) bitbucketPRLink(proj string, repo string, prid int) (slack.Attachment, error) {
	attachement := slack.Attachment{}

	// Get the Pull Request
	pr, err := u.Bitbucket.PullRequest(proj, repo, prid)
	if err != nil {
		return attachement, err
	}

	// Get the Pull Request Status
	st, err := u.Bitbucket.Status(pr.FromRef.LatestCommit)
	if err != nil {
		return attachement, err
	}

	fields := []slack.AttachmentField{
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
	}
	if pr.ApprovalStatus(false) != "" {
		fields = append(
			fields,
			slack.AttachmentField{
				Title: "Reviewers",
				Value: pr.ReviewedBy(),
				Short: true,
			}, slack.AttachmentField{
				Title: "Review Status",
				Value: pr.ApprovalStatus(true),
				Short: true,
			},
		)
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
	attachement.Fields = fields

	return attachement, nil
}

// bitbucketRepoLink returns a Slack Attachment for a Bitbucket Repo links
func (u *Unfurl) bitbucketRepoLink(project string, repo string) (slack.Attachment, error) {
	var attachement slack.Attachment

	// Get repo info
	r, err := u.Bitbucket.Repository(project, repo)
	if err != nil {
		return attachement, err
	}

	// Get repo commits
	co, err := u.Bitbucket.Commits(project, repo, bitbucket.CommitOptions{})
	if err != nil {
		return attachement, err
	}

	// Get build status for latest commit
	st, err := u.Bitbucket.Status(co.Values[0].ID)
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
