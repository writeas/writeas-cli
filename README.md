writeas-cli
===========
![GPL](https://img.shields.io/github/license/writeas/writeas-cli.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/writeas/writeas-cli)](https://goreportcard.com/report/github.com/writeas/writeas-cli) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Discuss on our forum](https://img.shields.io/discourse/https/discuss.write.as/users.svg?label=forum)](https://discuss.write.as/c/development)

Command line interface for [Write.as](https://write.as). Works on Windows, macOS, and Linux.

## Features

* Publish anonymously to Write.as
* Authenticate with a Write.as account
* A stable, easy back-end for your GUI app or desktop-based workflow
* Compatible with our [Tor hidden service](http://writeas7pm7rcdqg.onion/)
* Locally keeps track of any posts you make
* Update and delete posts, anonymous and authenticated
* Fetch any post by ID
* Add anonymous post credentials (like for one published with the [Android app](https://play.google.com/store/apps/details?id=com.abunchtell.writeas)) for editing

## Installing
The easiest way to get the CLI is to download a pre-built executable for your OS.

### Download
[![Latest release](https://img.shields.io/github/release/writeas/writeas-cli.svg)](https://github.com/writeas/writeas-cli/releases/latest) ![Total downloads](https://img.shields.io/github/downloads/writeas/writeas-cli/total.svg) 

Get the latest version for your operating system as a standalone executable.

**Windows**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v1.1/writeas_1.1_windows_amd64.zip) or [32-bit](https://github.com/writeas/writeas-cli/releases/download/v1.1/writeas_1.1_windows_386.zip) executable and put it somewhere in your `%PATH%`.

**macOS**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v1.1/writeas_1.1_darwin_amd64.tar.gz) executable and put it somewhere in your `$PATH`, like `/usr/local/bin`.

**Linux**<br />
Download the [64-bit](https://github.com/writeas/writeas-cli/releases/download/v1.1/writeas_1.1_linux_amd64.tar.gz) or [32-bit](https://github.com/writeas/writeas-cli/releases/download/v1.1/writeas_1.1_linux_386.tar.gz) executable and put it somewhere in your `$PATH`, like `/usr/local/bin`.

### Go get it
```bash
go get github.com/writeas/writeas-cli/cmd/writeas
```

Once this finishes, you'll see `writeas` or `writeas.exe` inside `$GOPATH/bin/`.

## Upgrading

To upgrade the CLI, download and replace the executable you downloaded before.

If you previously installed with `go get`, run it again with the `-u` option.

```bash
go get -u github.com/writeas/writeas-cli/cmd/writeas
```

## Usage

See full usage documentation on our [User Guide](GUIDE.md).

```
   writeas [global options] command [command options] [arguments...]

COMMANDS:
     post     Alias for default action: create post from stdin
     new      Compose a new post from the command-line and publish
     publish  Publish a file to Write.as
     delete   Delete a post
     update   Update (overwrite) a post
     get      Read a raw post
     add      Add an existing post locally
     posts    List all of your posts
     blogs    List blogs
     auth     Authenticate with Write.as
     logout   Log out of Write.as
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -c value, -b value  Optional blog to post to
   --tor, -t           Perform action on Tor hidden service
   --tor-port value    Use a different port to connect to Tor (default: 9150)
   --code              Specifies this post is code
   --md                Returns post URL with Markdown enabled
   --verbose, -v       Make the operation more talkative
   --font value        Sets post font to given value (default: "mono")
   --lang value        Sets post language to given ISO 639-1 language code
   --user-agent value  Sets the User-Agent for API requests
   --help, -h          show help
   --version, -V       print the version
```

## Contributing to the CLI

For a complete guide to contributing, see the [Contribution Guide](.github/CONTRIBUTING.md).

We welcome any kind of contributions including documentation, organizational improvements, tutorials, bug reports, feature requests, new features, answering questions, etc.

### Getting Support

We're available on [several channels](https://write.as/contact), and prefer our [forum](https://discuss.write.as) for project discussion. Please don't use the GitHub issue tracker to ask questions.

### Reporting Issues

If you believe you have found a bug in the CLI or its documentation, file an issue on this repo. If you're not sure if it's a bug or not, [reach out to us](https://write.as/contact) in one way or another. Be sure to provide the version of the CLI (with `writeas --version`) in your report.
