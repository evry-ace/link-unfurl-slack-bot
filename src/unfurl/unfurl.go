package unfurl

import (
	"net/url"

	"github.com/bndr/gojenkins"
	"github.com/evry-ace/link-unfurl-slack-bot/src/bitbucket"
	"github.com/evry-ace/link-unfurl-slack-bot/src/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Unfurl is an inverted control structure for the unfurl package
type Unfurl struct {
	Logger    *logrus.Logger
	Jenkins   *gojenkins.Jenkins
	Bitbucket *bitbucket.Client
	Config    *utils.Config
}

// Links unfurls a all links from a Slack LinkSharedEvent and returns a
// slack.Attachment for each link.
func (u *Unfurl) Links(event *slackevents.LinkSharedEvent) (map[string]slack.Attachment, error) {
	// Create a new map to store link unfurled data as Slack attachments
	unfurls := make(map[string]slack.Attachment, len(event.Links))

	// Unfurl all the shared links
	for _, link := range event.Links {
		var attachement slack.Attachment
		var err error

		// Parse the link
		URL, urlErr := url.Parse(link.URL)
		if urlErr != nil {
			u.Logger.Errorf("Error parsing url: %s", err)
			continue
		}

		u.Logger.Infof("Unfurling link: %s", URL.String())

		// Check the link domain, discard if not supported
		switch link.Domain {
		case u.Config.BitbucketServer:
			attachement, err = u.bitbucketLink(URL)

		case u.Config.JenkinsServer:
			attachement, err = u.jenkinsLink(URL)

		default:
			u.Logger.Debugf("Unsupported link domain: %s", link.Domain)
			continue
		}

		if err != nil {
			u.Logger.WithError(err).WithField("link", link).Error("Failed to unfurl link")
		} else {
			unfurls[link.URL] = attachement
		}
	}

	return unfurls, nil
}
