writeas-cli
===========
Command line interface for [Write.as](https://write.as) and [Write.as on Tor](http://writeas7pm7rcdqg.onion/). Works on Windows, OS X, and Linux.

Like the [Android app](https://play.google.com/store/apps/details?id=com.abunchtell.writeas), the command line client keeps track of the posts you make, so future editing / deleting is easier than [doing it with cURL](http://cmd.write.as/). The goal is for this to serve as the backend for any future GUI app we build for the desktop.

It is currently **alpha**, so a) functionality is basic and b) everything is subject to change — i.e., watch the [changelog](https://write.as/changelog-cli.html).

## Usage

```
writeas [global options] command [command options] [arguments...]

COMMANDS:
   post     Alias for default action: create post from stdin
   delete   Delete a post
   update   Update (overwrite) a post
   get      Read a raw post
   add      Add a post locally
   list     List local posts
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --tor, -t		 Perform action on Tor hidden service
   --tor-port "9150" Use a different port to connect to Tor
   --help, -h		 show help
   --version, -v	 print the version
```

## Download

Get it on the [web](https://write.as/cli.html) or [hidden service](http://writeas7pm7rcdqg.onion/cli.html).

## Go get it
`go get github.com/writeas/writeas-cli`
