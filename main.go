package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/evry-ace/link-unfurl-slack-bot/src/bitbucket"
	"github.com/evry-ace/link-unfurl-slack-bot/src/unfurl"
	"github.com/evry-ace/link-unfurl-slack-bot/src/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	utils.SetupLogging()
}

func main() {
	c, configErr := utils.ConfigFromEnvironment(".env")
	if configErr != nil {
		logrus.Fatal(configErr.Error(), "config loading failed")
	}

	b := bitbucket.Client{Server: c.BitbucketServer, PAT: c.BitbucketPAT}

	if !strings.HasPrefix(c.SlackAppToken, "xapp-") {
		logrus.Fatal("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	if !strings.HasPrefix(c.SLackBotToken, "xoxb-") {
		logrus.Fatal("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	// Slack SDK
	api := slack.New(
		c.SLackBotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(c.SlackAppToken),
	)

	// Slack Events API
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Listen for events
	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				logrus.Info("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				logrus.Info("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				logrus.Info("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					logrus.WithField("event", evt).Warn("Event type is not EventsAPIEvent")
					continue
				}

				logrus.WithField("event", eventsAPIEvent).Info("Slack event received")

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					logrus.Debug("Callback event received")
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.LinkSharedEvent:
						logrus.WithField("event", ev).Debug("LinkSharedEvent received")

						// Create a new map to store link unfurled data as Slack attachments
						unfurls := make(map[string]slack.Attachment, len(ev.Links))

						// Get slack channel name
						channel, err := api.GetConversationInfo(ev.Channel, false)
						if err != nil {
							logrus.WithError(err).WithField("event", ev).Error("Failed to get Slack channel info")
							continue
						}

						// Only unfurl links from the channel matching regex
						re := regexp.MustCompile(c.ChannelRegex)
						if !re.MatchString(channel.Name) {
							logrus.WithFields(logrus.Fields{
								"event":   ev,
								"regex":   c.ChannelRegex,
								"channel": channel,
							}).Debug("Ignoring link unfurl from unsupported channel")
							continue
						}

						// Unfurl all the shared links
						for _, link := range ev.Links {
							// Check the link domain, discard if not supported
							// Currently only Bitbucket links are supported
							if link.Domain != c.BitbucketServer {
								continue
							}

							// Create the unfurl attachement
							attachement, linkErr := unfurl.BitbucketLink(link.URL, b)
							if linkErr != nil {
								logrus.WithError(linkErr).WithField("link", link).Error("Failed to unfurl Bitbucket link")
							} else {
								unfurls[link.URL] = attachement
							}
						}

						logrus.WithField("unfurls", unfurls).Debug("Unfurls")

						if len(unfurls) > 0 {
							_, _, err := api.PostMessage(
								ev.Channel,
								slack.MsgOptionUnfurl(ev.MessageTimeStamp, unfurls),
							)
							if err != nil {
								logrus.WithError(err).WithField("unfurls", unfurls).Error("Failed to post Slack message")
							}
						}
					default:
						client.Debugf("unsupported Callback API event received")
					}

				default:
					client.Debugf("unsupported Events API event received")
				}

			case socketmode.EventTypeHello:
				//numConnections := evt.Request.NumConnections
				logrus.WithField("evt", fmt.Sprintf("%v", evt)).Info("Hello event received")

			default:
				logrus.WithFields(logrus.Fields{
					"event": evt,
					"type":  evt.Type,
				}).Warn("Unhandled event received")
			}
		}
	}()

	client.Run()
}
