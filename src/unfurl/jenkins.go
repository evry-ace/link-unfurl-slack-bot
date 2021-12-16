package unfurl

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"github.com/xeonx/timeago"
)

const (
	JenkinsURLBuildType  = "build"
	JenkinsURLUknownType = "unknown"
)

// jenkinsLinkType returns the type of Bitbucket link and the matches
func (u *Unfurl) jenkinsLinkType(URL *url.URL) (string, []string) {
	// /job/k8s/job/tf-dockyard/job/master/821/
	var isBuild = regexp.MustCompile(`^/job/([^/]+)/job/([^/]+)/job/([^/]+)/([0-9]+)`)

	if isBuild.MatchString(URL.Path) {
		return JenkinsURLBuildType, isBuild.FindStringSubmatch(URL.Path)
	}

	return JenkinsURLUknownType, []string{}
}

// jenkinsLink returns a slack.Attachment for a Jenkins link
func (u *Unfurl) jenkinsLink(URL *url.URL) (slack.Attachment, error) {
	attachement := slack.Attachment{}

	// Parse what type of link this is
	linkType, matches := u.jenkinsLinkType(URL)
	fmt.Printf("linkType=%s matches=%s", linkType, matches)
	switch linkType {
	case JenkinsURLBuildType:
		project := matches[1]
		repo := matches[2]
		branch := matches[3]
		buildID, err := strconv.Atoi(matches[4])
		if err != nil {
			return attachement, err
		}

		fmt.Printf("project=%s repo=%s branch=%s buildID=%d", project, repo, branch, buildID)

		return u.jenkinsBuildLink(project, repo, branch, buildID)

	default:
		return attachement, errors.New("jenkins link not supported")
	}
}

func (u *Unfurl) jenkinsBuildLink(project, repo, branch string, buildNumber int) (slack.Attachment, error) {
	attachement := slack.Attachment{}
	ctx := context.Background()

	jobName := fmt.Sprintf("%s/job/%s/job/%s", project, repo, branch)
	fmt.Printf("jobName=%s", jobName)

	build, err := u.Jenkins.GetBuild(ctx, jobName, int64(buildNumber))
	if err != nil {
		return attachement, err
	}

	fmt.Printf("build=%+v", build.Raw)

	// Build result
	result := build.GetResult()
	if result == "" {
		result = "IN PROGRESS"
	}

	// Build start time
	started := timeago.NoMax(timeago.English).Format(time.Unix(build.Raw.Timestamp/1000, 0))

	// Build Duration
	var duration time.Duration

	if build.Raw.Duration > 0 {
		duration = time.Duration(build.Raw.Duration) * time.Second
	} else {
		duration = time.Duration(time.Now().Unix()-build.Raw.Timestamp/1000) * time.Second
	}

	attachement.Title = build.Raw.FullDisplayName
	attachement.TitleLink = build.GetUrl()
	attachement.AuthorName = build.Raw.FullDisplayName
	// attachement.Text = build.Raw.ChangeSet.Items[0].Msg

	// check if jenkins build is waiting for input
	if build.Raw.Building {
		attachement.Text = "Waiting for input"
	} else {
		attachement.Text = fmt.Sprintf("%s\n%s", result, duration.String())
	}

	attachement.Fields = []slack.AttachmentField{
		{
			Title: "Status",
			Value: result,
			Short: true,
		},
		{
			Title: "Duration",
			Value: duration.String(),
			Short: true,
		},
		{
			Title: "Started",
			Value: started,
			Short: true,
		},
	}

	attachement.CallbackID = "jenkins_build"
	attachement.Actions = []slack.AttachmentAction{
		{
			Name: "build log",
			Text: ":page_facing_up: Build Log",
			Type: "button",
			URL:  build.GetUrl() + "console",
		},
		{
			Name: "build changes",
			Text: ":compass: Change Log",
			Type: "button",
			URL:  build.GetUrl() + "changes",
		},
	}

	return attachement, nil
}
