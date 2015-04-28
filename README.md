writeas-cli
===========
Command line interface for [Write.as](https://write.as) and [Write.as on Tor](http://writeas7pm7rcdqg.onion/). Works on Windows, OS X, and Linux.

Like the [Android app](https://play.google.com/store/apps/details?id=com.abunchtell.writeas), the command line client keeps track of the posts you make, so future editing / deleting is easier than [doing it with cURL](http://cmd.write.as/). It is currently **ALPHA**, so only basic functionality is available. But the goal is for this to hold the logic behind any future GUI app we build for the desktop.

## Usage

```
writeas [global options] command [command options] [arguments...]

COMMANDS:
   post     Alias for default action: create post from stdin
   delete   Delete a post
   get      Read a raw post
   add      Add a post locally for easy modification
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
