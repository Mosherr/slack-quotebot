# Hosting a Slack Slash Command on Google App Engine

This is built as a slash command custom integration for Slack,
running on App Engine.

## Prerequisites

You need to be able to deploy Go applications to App Engine. I recommend
following the [guide for Go App Engine][go-appengine], first. At a minimum, you
should [install Go][go-install] and the [Go App Engine SDK][go-appengine-sdk].

[go-appengine]: https://cloud.google.com/appengine/docs/go/
[go-install]: https://golang.org/doc/install
[go-appengine-sdk]: https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go

## App Engine setup
- SLACKQUOTEBOT=~/src/slackquotebot/slackquotebot-2017-03-05-15-50
- git clone https://github.com/Mosherr/slack-quotebot.git $SLACKQUOTEBOT
- cd $SLACKQUOTEBOT
- git checkout master
- goapp serve -host 0.0.0.0 $PWD
- gcloud app create
- goapp deploy -application slackquotebot -version 0
- For Mongo to work you need to open a port
-- enable Compute Engine API
-- gcloud compute firewall-rules create mongodb --allow tcp:27017
-- where 27017 is the port in your mongo connection string

## Setup Slack

We need to create a [new slash command][new-slash-command] custom integration.

- Point it at `https://YOUR-PROJECT.appspot.com`, replace
  YOUR-PROJECT with your Google Cloud Project ID.
- Give it a name, like `/quotebot`.
- Write the token it gives you to `config.go`.

[new-slash-command]: https://my.slack.com/services/new/slash-commands

## Build and Deploy

Then, deploy. Replace your-project with your Google Cloud Project ID:

```
goapp deploy -application your-project app.yaml
```


## Try It Out

Run `/quotebot` in Slack
- /quotebot -g (get a random quote)
- /quotebot -g @user (get a random quote from specified user)
- /quotebot -a @user "new quote to add" (adds a quote for a user)
