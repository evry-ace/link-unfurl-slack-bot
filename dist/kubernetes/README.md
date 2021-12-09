# Deploy to Kubernetes

This quickstart will walk you through deploying the bot to Kubernetes.

First you need to create a secret with the bot's token.

```bash
kubectl create secret generic link-unfurl-slack-bot \
  --from-literal=BITBUCKET_PAT=<your-bitbucket-pat> \
  --from-literal=BITBUCKET_SERVER=<your-bitbucket-server> \
  --from-literal=SLACK_APP_TOKEN=<your-slack-app-token> \
  --from-literal=SLACK_BOT_TOKEN=<your-slack-bot-token>
```

Then you can deploy the bot to Kubernetes.

```bash
kubectl apply -f deployment.yaml
```
