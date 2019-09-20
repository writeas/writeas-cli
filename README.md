writeas-cli / wf-cli
====================
![GPL](https://img.shields.io/github/license/writeas/writeas-cli.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/writeas/writeas-cli)](https://goreportcard.com/report/github.com/writeas/writeas-cli) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Discuss on our forum](https://img.shields.io/discourse/https/discuss.write.as/users.svg?label=forum)](https://discuss.write.as/c/development)

Command line utility for publishing to [Write.as](https://write.as) and any other [WriteFreely](https://writefreely.org) instance. Works on Windows, macOS, and Linux.

## Features

* Authenticate with a Write.as / WriteFreely account
* Publish anonymous posts or drafts to Write.as or WriteFreely, respectively
* A stable, easy back-end for your [GUI app](https://write.as/apps/desktop) or desktop-based workflow
* Compatible with the [Write.as Tor hidden service](http://writeas7pm7rcdqg.onion/)
* Update and delete posts
* Fetch any post by ID
* ...and more, depending on which client you're using (see respective READMEs for more)

## Installing
The easiest way to get the CLI is to download a pre-built executable for your OS.

### Download
[![Latest release](https://img.shields.io/github/release/writeas/writeas-cli.svg)](https://github.com/writeas/writeas-cli/releases/latest) ![Total downloads](https://img.shields.io/github/downloads/writeas/writeas-cli/total.svg) 

Get the latest version for your operating system as a standalone executable.

**Write.as CLI**<br />
See the [writeas-cli README](https://github.com/writeas/writeas-cli/cmd/writeas#readme)

**WriteFreely CLI**<br />
See the [wf-cli README](https://github.com/writeas/writeas-cli/cmd/wf#readme)

## Usage

**Write.as CLI**<br />
See full usage documentation on our [writeas-cli User Guide](https://github.com/writeas/writeas-cli/blob/master/cmd/writeas/GUIDE.md).

**WriteFreely CLI**<br />
See full usage documentation on our [wf-cli User Guide](https://github.com/writeas/writeas-cli/blob/master/cmd/wf/GUIDE.md).

## Contributing to the CLI

For a complete guide to contributing, see the [Contribution Guide](.github/CONTRIBUTING.md).

We welcome any kind of contributions including documentation, organizational improvements, tutorials, bug reports, feature requests, new features, answering questions, etc.

### Getting Support

We're available on [several channels](https://write.as/contact), and prefer our [forum](https://discuss.write.as) for project discussion. Please don't use the GitHub issue tracker to ask questions.

### Reporting Issues

If you believe you have found a bug in the CLI or its documentation, file an issue on this repo. If you're not sure if it's a bug or not, [reach out to us](https://write.as/contact) in one way or another. Be sure to provide the version of the CLI (with `writeas --version` or `wf --version`) in your report.
