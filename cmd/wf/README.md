wf-cli
======
![GPL](https://img.shields.io/github/license/writeas/writeas-cli.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/writeas/writeas-cli)](https://goreportcard.com/report/github.com/writeas/writeas-cli) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Discuss on our forum](https://img.shields.io/discourse/https/discuss.write.as/users.svg?label=forum)](https://discuss.write.as/c/development)

Command line utility for publishing to any [WriteFreely](https://writefreely.org) instance. Works on Windows, macOS, and Linux.

**The WriteFreely CLI is compatible with WriteFreely v0.11 or later.**

## Features

* Authenticate with any WriteFreely instance
* Publish drafts
* Manage multiple WriteFreely accounts on multiple instances
* A stable, easy back-end for your GUI app or desktop-based workflow
* Locally keeps track of any posts you make
* Update and delete posts
* Fetch any post by ID

## Installing
The easiest way to get the CLI is to download a pre-built executable for your OS.

### Download
[![Latest release](https://img.shields.io/github/release/writeas/writeas-cli.svg)](https://github.com/writeas/writeas-cli/releases/latest) ![Total downloads](https://img.shields.io/github/downloads/writeas/writeas-cli/total.svg) 

Get the latest version for your operating system as a standalone executable.

**Windows**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v2.0.0/wf_2.0.0_windows_amd64.zip) or [32-bit](https://github.com/writeas/writeas-cli/releases/download/v2.0.0/wf_2.0.0_windows_386.zip) executable and put it somewhere in your `%PATH%`.

**macOS**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v2.0.0/wf_2.0.0_darwin_amd64.zip) executable and put it somewhere in your `$PATH`, like `/usr/local/bin`.

**Debian-based Linux**<br />
```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys DBE07445
sudo add-apt-repository "deb http://updates.writeas.org xenial main"
sudo apt-get update && sudo apt-get install wf-cli
```

**Linux (other)**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v2.0.0/wf_2.0.0_linux_amd64.tar.gz) or [32-bit](https://github.com/writeas/writeas-cli/releases/download/v2.0.0/wf_2.0.0_linux_386.tar.gz) executable and put it somewhere in your `$PATH`, like `/usr/local/bin`.

### Go get it
```bash
go get github.com/writeas/writeas-cli/cmd/wf
```

Once this finishes, you'll see `wf` or `wf.exe` inside `$GOPATH/bin/`.

## Upgrading

To upgrade the CLI, download and replace the executable you downloaded before.

If you previously installed with `go get`, run it again with the `-u` option.

```bash
go get -u github.com/writeas/writeas-cli/cmd/wf
```

## Usage

See full usage documentation on our [User Guide](https://github.com/writeas/writeas-cli/blob/master/cmd/wf/GUIDE.md).

```
   wf [global options] command [command options] [arguments...]

COMMANDS:
     post      Alias for default action: create post from stdin
     new       Compose a new post from the command-line and publish
     publish   Publish a file
     delete    Delete a post
     update    Update (overwrite) a post
     get       Read a raw post
     posts     List draft posts
     blogs     List blogs
     accounts  List all currently logged in accounts
     auth      Authenticate with a WriteFreely instance
     logout    Log out of a WriteFreely instance
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -c value, -b value      Optional blog to post to
   --tor, -t               Perform action on Tor hidden service
   --tor-port value        Use a different port to connect to Tor (default: 9150)
   --code                  Specifies this post is code
   --md                    Returns post URL with Markdown enabled
   --verbose, -v           Make the operation more talkative
   --font value            Sets post font to given value (default: "mono")
   --lang value            Sets post language to given ISO 639-1 language code
   --user-agent value      Sets the User-Agent for API requests
   --host value, -H value  Operate against a custom hostname
   --user value, -u value  Use authenticated user, other than default
   --help, -h              show help
   --version, -V           print the version
```
