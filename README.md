# Link Unfurl Slack-bot

Slack-bot to do link unfurl for private endpoints.

## Features

* [x] Atlassian Bitbucket Server
* [ ] Atlassian Confluence Server
* [ ] Atlassian JIRA Server

## Configuration

| Environment Variable | Description | Required | Default |
|----------------------|-------------|----------|---------|
| `LOGLEVEL`           | Logrus log level | `false` | `debug` |
| `LOGFORMAT`          | Logrus log format | `false` | `text` |
| `BITBUCKET_PAT`      | Bitbucket Personal Access Token | `true` | `""` |
| `BITBUCKET_SERVER`   | Bitbucket Server Hostname | `true` | `""` |
| `SLACK_APP_TOKEN`    | Slack App Token | `true` | `""` |
| `SLACK_BOT_TOKEN`    | Slack Bot Token | `true` | `""` |
| `CHANNEL_REGEX`      | Enabled channels for link unfurling | `false` | `"^devops-([a-zA-Z0-9_]+)$"` |
