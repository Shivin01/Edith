# Edith slack bot
This Slack bot improves the workflow of development teams. Especially with focus on Jenkins and some other functionalities inluding add user, remove user, modify user, request for leaves, approval leaves etc.

# Installation
## 1st) Create+prepare the Slack App:
### Recommended way: Use Manifest file as App template
1. Create a [Slack App](https://api.slack.com/apps?new_app=1)
2. Select "From an app manifest"
3. Select your Workspace
4. Paste this Yaml code:

```yaml
_metadata:
  major_version: 1
  minor_version: 1
display_information:
  name: edith
features:
  app_home:
    messages_tab_enabled: true
    messages_tab_read_only_enabled: false
  bot_user:
    display_name: edith
    always_online: true
oauth_config:
  scopes:
    bot:
      - app_mentions:read
      - channels:read
      - chat:write
      - im:history
      - im:write
      - mpim:history
      - reactions:read
      - reactions:write
      - users:read
      - files:write
settings:
  event_subscriptions:
    bot_events:
      - app_mention
      - message.im
  interactivity:
    is_enabled: true
  org_deploy_enabled: false
  socket_mode_enabled: true
```

5. Create the App!
6. Go to "Basic Information"
7. -> in "App-Level Tokens", "Generate a Token" with the scope "connections:write"
8. You will see a App-Level Token (beginning with xapp-). Sse it in the config.yaml as slack.socket_token.
9. Go to "OAuth & Permissions":
10. -> "Install to Workspace"
11. -> you should see a "Bot User OAuth Access Token" (beginning with "xoxb-"). Use it as slack.token in the config.yaml
12. start the bot!

## 2nd) Run the bot

### Option 1: run via go
1. install go
2. create a config.yaml (at least a slack token is required) or take a look in config-example.yaml
3. go inside slack bot folder `cd slack-bot`
4. run `go run cmd/edith/main.go`


# Usage
**Note:** You have to invite the bot into the channel to be able to handle commands.
