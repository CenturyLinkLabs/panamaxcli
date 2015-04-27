# pmxcli

[![Circle CI](https://circleci.com/gh/CenturyLinkLabs/panamaxcli.svg?style=svg)](https://circleci.com/gh/CenturyLinkLabs/panamaxcli)

`pmxcli` allows you to manage deployments to Panamax Remote Agents, including
creation, redeployment, and deletion. It is installed alongside Panamax, but
can be run separately without Panamax or its dependencies.

## Installation

The utility is installed when you follow [the normal Panamax installation
process](http://panamax.io/get-panamax/), or it can be downloaded separately
for 64-bit [OSX](http://download.panamax.io/panamaxcli/panamaxcli-darwin) or
        [Linux](http://download.panamax.io/panamaxcli/panamaxcli-linux).

## Usage

There are two resources that can be managed in `pmxcli`: remotes and
deployments. Remotes are the Panamax Agents you have installed, and deployments
are the applications that are currently deployed on any one of those agents.
You can get an exhaustive list of commands and help by running `pmxcli` with no
arguments, but here are the basics to get you started.

First, you should add a remote. The token file should contain the text of the
token, as provided by the Panamax web UI:

```bash
% pmxcli remote add demo /path/to/tokenfile.txt
Successfully added! 'demo' is your active remote.

% pmxcli remote list
ACTIVE  NAME    ENDPOINT
*       demo    https://192.168.1.1:3001
```

Your first remote is automatically made active. The active remote will be the
one whose deployments you'll be interacting with when you run any `pmxcli
deployment` commands.

You can deploy any Panamax template, both existing ones you've downloaded from
[the public templates
repository](https://github.com/CenturyLinkLabs/panamax-public-templates), or
those you create yourself:

```bash
% pmxcli deployment create wordpress.pmx
Template successfully deployed as '1'

% pmxcli deployment describe 1
ID              1
Name            Wordpress with MySQL
Redeployable    true

SERVICES
ID              STATE
db.service      load_state: loaded; active_state: activating; sub_state: start-pre
wp.service      load_state: loaded; active_state: activating; sub_state: start-pre
```

Run `pmxcli deployment help` for a list of commands to interact with
deployments.

## Gotchas

#### SSL Warnings

Communication with the remotes happens over SSL, and `pmxcli` verifies of the
SSL certification from the agent to ensure that it is communicating with the
same server that originally generated the token. There may be edge cases where
you want to disable that verification, and so an `--insecure` global flag has
been included to skip that step.

If you've set up a remote agent using Panamax Remote Agent Installer 0.1.3 or
below, you may see a specific warning in cases where your remote uses an IP
address and not a hostname. You'll see something like:

```
x509: cannot validate certificate for X.X.X.X because it doesn't contain any IP SANs
```

You can use the `--insecure` flag as directed by the warning, but we recommend
upgrading the installer and reinstalling the agent.

#### Debugging

If you see unexpected results, there is a `--debug` global flag that will log
out all the requests and responses from the remote agent. This can be useful
for troubleshooting.
