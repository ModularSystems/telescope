![](./telescope.png)

Telescope
===

Telescope is used to monitor websites. Much like a telescope, it can crudely see things from far away. It can watch for new vulnerabilities, unexpected changes, outages, and more. It integrates with Sendgrid to send email alerts. Future releases aim to support integrations for:

- Slack
- Text via Email
- Twillio (Text)
- Twillio (Call)
- Configurable webhook

It works by fetching configured URIs at their configured time, storing the request output and other associated data, and running regexes against those configured attributes.
If a regex is matched against a particular attribute, the alert is triggered.

# Usage

Telescope is a binary designed to be run in a container ([modularsystems/telescope:latest](https://hub.docker.com/r/modularsystems/telescope)), so you can either run the container, or deploy it with your orchestrator of choice.

Here are some examples on how you can run it:

## Cron

Cron is a popular favorite for task automation. You can use this container with cron by creating your config locally. Then, pass in your configuration to the container by a volume mount.

Example `/etc/telescope/config.yaml`:
```
alerts:
  - name: Vulnerability Report

```



# Configuration:

Secrets data is passed via environment variables. Here is what can be configured by environment variables:

|Environment Variable|Description|
|-|-|
|SENDGRID_API_KEY|Configure your sendgrid token to receive email alerts|
|SENDGRID_SENDER_NAME|Your sendgrid verified user's name [See: Sendgrid sender identity](https://sendgrid.com/docs/for-developers/sending-email/sender-identity/)|
|SENDGRID_SENDER_EMAIL|Your sendgrid verified user's email address  [See also: Sendgrid sender identity](https://sendgrid.com/docs/for-developers/sending-email/sender-identity/)|
|WPVULNDB_API_KEY|Configure your wpvulndb token for richer wordpress vulnerability detection, you can get one free [here](https://wpvulndb.com/users/sign_up)|

The following flags exist to help you configure the application:

|flag|description|
|-|-|
|--debug|turn on debug logging|
|--config|path to configuration on the filesystem. default: /etc/telescope/config.yaml|

# Configuring Alerts

Alert configuration is done in YAML. An example alert might look like this:
```
alerts:
  - name: Non 200 Returned from mywebsite.com
    uris:
      - https://mywebsite.com
      - https://zombo.com
      - https://www.zombo.com
    scanner: HTMLScan
    attribute: "return code"
    regex: "^((?!.*200.*).)*$"
    every: 5m
    send: text
    page:
      - "555-555-5555"
      - "5555555555"

  - name: Plugin Vulnerability Discovered
    scanner: WPScan
    uris:
      - https://www.zombo.com
      - https://zombo.com
    send: email
    page:
      - "Dev Null dev@null"
      - "dead beef dead@beef"
    regex: "Plugin.Vulernability.Discovered"
    timed: "13:00 UTC"
```
