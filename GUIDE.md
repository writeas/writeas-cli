# Write.as CLI User Guide

The Write.as Command-Line Interface (CLI) is a cross-platform tool for publishing text to [Write.as](https://write.as) and its other sites, like [Paste.as](https://paste.as). It is designed to be simple, scriptable, do one job (publishing) well, and work as you'd expect with other command-line tools.

Write.as is a text-publishing service that protects your privacy. There's no sign up required to publish, but if you do sign up, you can access posts across devices and compile collections of them in what most people would call a "blog".

**Note** accounts are not supported in CLI v1.0. They'll be available in [v2.0](https://github.com/writeas/writeas-cli/milestone/4).

## Uses

These are a few common uses for `writeas`. If you get stuck or want to know more, run `writeas [command] --help`. If you still have questions, [ask us](https://write.as/contact).

### Overview

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
   --code            Specifies this post is code
   --verbose, -v     Make the operation more talkative
   --font value      Sets post font to given value (default: "mono")
   --help, -h		 show help
   --version, -v	 print the version
```

#### Share something

By default, `writeas` creates a post with a `monospace` typeface that doesn't word wrap (scrolls horizontally). It will return a single line with a URL, and automatically copy that URL to the clipboard:

```bash
$ echo "Hello world!" | writeas
https://write.as/aaaaaaaaaaaa
```

This is generally more useful for posting terminal output or code, like so (the `--code` flag turns on syntax highlighting):

macOS / Linux: `cat writeas/cli.go | writeas --code`

Windows: `type writeas/cli.go | writeas.exe --code`

#### Output a post

This outputs any Write.as post with the given ID.

```bash
$ writeas get aaaaaaaaaaaa
Hello world!
```

#### List all published posts

This lists all posts you've published from your device. Pass the `--url` flag to show the list with full URLs.

```bash
$ writeas list
aaaaaaaaaaaa
```

#### Delete a post

This permanently deletes a post you own.

```bash
$ writeas delete aaaaaaaaaaaa
```

#### Update a post

This completely overwrites an existing post you own.

```bash
$ echo "See you later!" | writeas update aaaaaaaaaaaa
```

### Composing posts

If you simply have a penchant for never leaving your keyboard, `writeas` is great for composing new posts from the command-line. Just use the `new` subcommand.

`writeas new` will open your favorite command-line editor, as specified by your `WRITEAS_EDITOR` or `EDITOR` environment variables (in that order), falling back to `vim` on OS X / *nix.

Customize your post's appearance with the `--font` flag:

| Argument | Appearance (Typeface) | Word Wrap? |
| -------- | --------------------- | ---------- |
| `sans` | Sans-serif (Open Sans) | Yes |
| `serif` | Serif (Lora) | Yes |
| `wrap` | Monospace | Yes |
| `mono` | Monospace | No |
| `code` | Syntax-highlighted monospace | No |

Put it all together, e.g. publish with a sans-serif font: `writeas new --font sans`
