writeas-cli
===========
![MIT license](https://img.shields.io/github/license/writeas/writeas-cli.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/writeas/writeas-cli)](https://goreportcard.com/report/github.com/writeas/writeas-cli) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Public Slack discussion](http://slack.write.as/badge.svg)](http://slack.write.as/)

Command line interface for [Write.as](https://write.as) and [Write.as on Tor](http://writeas7pm7rcdqg.onion/). Works on Windows, OS X, and Linux.

Like the [Android app](https://play.google.com/store/apps/details?id=com.abunchtell.writeas), the command line client keeps track of the posts you make, so future editing / deleting is easier than [doing it with cURL](http://cmd.write.as/). The goal is for this to serve as the backend for any future GUI app we build for the desktop.

It is currently **alpha**, so a) functionality is basic and b) everything is subject to change — i.e., watch the [changelog](https://write.as/changelog-cli.html).

## Download
[![Latest release](https://img.shields.io/github/release/writeas/writeas-cli.svg)](https://github.com/writeas/writeas-cli/releases/latest) ![Total downloads](https://img.shields.io/github/downloads/writeas/writeas-cli/total.svg) 

Get the latest version for your operating system as a standalone executable.

**Windows**: [64-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_windows_amd64.zip) – [32-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_windows_386.zip)

**OS X**: [64-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_darwin_amd64.zip) – [32-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_darwin_386.zip)

**Linux**: [64-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_linux_amd64.tar.gz) – [32-bit](https://github.com/writeas/writeas-cli/releases/download/v0.4/writeas_0.4_linux_386.tar.gz)

### Go get it
`go get github.com/writeas/writeas-cli/writeas`

## Usage

See full usage documentation on our [User Guide](GUIDE.md).

```
writeas [global options] command [command options] [arguments...]

COMMANDS:
   post     Alias for default action: create post from stdin
   new      Compose a new post from the command-line and publish
   delete   Delete a post
   update   Update (overwrite) a post
   get      Read a raw post
   add      Add an existing post locally
   list     List local posts
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --tor, -t		 Perform action on Tor hidden service
   --tor-port "9150" Use a different port to connect to Tor
   --help, -h		 show help
   --version, -v	 print the version
```
