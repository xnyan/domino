import slackweb
import sys

with open('slack_webhook_url.txt', 'r') as file:
    webhook_url = file.read().strip()

slack = slackweb.Slack(url=webhook_url)
slack.notify(text=sys.argv[1])
