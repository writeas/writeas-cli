# WriteFreely CLI User Guide

The WriteFreely Command-Line Interface (CLI) is a cross-platform tool for publishing text to any [WriteFreely](https://writefreely.org) instance. It is designed to be simple, scriptable, do one job (publishing) well, and work as you'd expect with other command-line tools.

WriteFreely is the software behind [Write.as](https://write.as). While the WriteFreely CLI supports publishing to Write.as, we recommend using the dedicated [Write.as CLI](https://github.com/writeas/writeas-cli/tree/master/cmd/writeas#readme) to get the full features of the platform, including anonymous publishing.

**The WriteFreely CLI is compatible with WriteFreely v0.11 or later.**

## Uses

These are a few common uses for `wf`. If you get stuck or want to know more, run `wf [command] --help`. If you still have questions, [ask us](https://write.as/contact).

### Overview

```
   wf [global options] command [command options] [arguments...]

COMMANDS:
     post     Alias for default action: create post from stdin
     new      Compose a new post from the command-line and publish
     publish  Publish a file
     delete   Delete a post
     update   Update (overwrite) a post
     get      Read a raw post
     posts    List all of your posts
     blogs    List blogs
     auth     Authenticate with a WriteFreely instance
     logout   Log out of a WriteFreely instance
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -c value, -b value      Optional blog to post to
   --insecure              Send request insecurely.
   --tor, -t               Perform action on Tor hidden service
   --tor-port value        Use a different port to connect to Tor (default: 9150)
   --code                  Specifies this post is code
   --verbose, -v           Make the operation more talkative
   --font value            Sets post font to given value (default: "mono")
   --lang value            Sets post language to given ISO 639-1 language code
   --user-agent value      Sets the User-Agent for API requests
   --host value, -H value  Use the given WriteFreely instance hostname
   --user value, -u value  Use the given account username
   --help, -h              show help
   --version, -V           print the version
```

#### Authenticate

To use the WriteFreely CLI, you'll first need to authenticate with the WriteFreely instance you want to interact with.

You may authenticate with as many WriteFreely instances and accounts as you want. But the first account you authenticate with will automatically be set as the default instance to operate on, so you don't have to supply `--host` and `--user` with every command.

```bash
$ wf --host pencil.writefree.ly auth username
Password: ************
```

In this example, you'll be authenticated as the user **username** on the WriteFreely instance **https://pencil.writefree.ly**.

#### Choosing an account

To select the WriteFreely instance and account you want to interact with, supply the `--host` and `--user` flags at the beginning of your `wf` command, e.g.:

```
$ wf --host pencil.writefree.ly --user username <subcommand>
```

If you're authenticated with only one account on any given WriteFreely instance, you only need to supply the `--host`, and `wf` will automatically use the correct account. E.g.:

```
$ wf --host pencil.writefree.ly <subcommand>
```

If a default account is set in `~/.writefreely/config.ini` and you want to use it, you don't need to supply any additional arguments. E.g.:

```
$ wf <subcommand>
```

#### Share something

By default, `wf` creates a post with a `monospace` typeface that doesn't word wrap (scrolls horizontally). It will return a single line with a URL, and automatically copy that URL to the clipboard.

```bash
$ echo "Hello world!" | wf
https://pencil.writefree.ly/aaaaazzzzz
```

This is generally more useful for posting terminal output or code, like so (the `--code` flag turns on syntax highlighting):

macOS / Linux: `cat cmd/wf/cli.go | wf --code`

Windows: `type cmd/wf/cli.go | wf.exe --code`

#### Output a post

This outputs any WriteFreely post with the given ID.

```bash
$ wf get aaaaazzzzz
Hello world!
```

#### List all blogs

This will output a list of the authenticated user's blogs.
```bash
$ wf blogs
Alias    Title
user     An Example Blog
dev      My Dev Log
```

#### List posts

This lists all draft posts you've published.

Pass the `--url` flag to show the list with full post URLs.

```bash
$ wf posts
aaaaazzzzz

$ wf posts -url
https://pencil.writefree.ly/aaaaazzzzz

$ wf posts
ID
aaaaazzzzz
```

#### Delete a post

This permanently deletes a post with the given ID.

```bash
$ wf delete aaaaazzzzz
```

#### Update a post

This completely overwrites an existing post with the given ID.

```bash
$ echo "See you later!" | wf update aaaaazzzzz
```

### Composing posts

If you simply have a penchant for never leaving your keyboard, `wf` is great for composing new posts from the command-line. Just use the `new` subcommand.

`wf new` will open your favorite command-line editor, as specified by your `WRITEAS_EDITOR` or `EDITOR` environment variables (in that order), falling back to `vim` on OS X / *nix.

Customize your post's appearance with the `--font` flag:

| Argument | Appearance (Typeface) | Word Wrap? |
| -------- | --------------------- | ---------- |
| `sans` | Sans-serif (Open Sans) | Yes |
| `serif` | Serif (Lora) | Yes |
| `wrap` | Monospace | Yes |
| `mono` | Monospace | No |
| `code` | Syntax-highlighted monospace | No |

Put it all together, e.g. publish with a sans-serif font: `wf new --font sans`
