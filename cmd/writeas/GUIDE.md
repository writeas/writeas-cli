# Write.as CLI User Guide

The Write.as Command-Line Interface (CLI) is a cross-platform tool for publishing text to [Write.as](https://write.as) and its other sites, like [Paste.as](https://paste.as). It is designed to be simple, scriptable, do one job (publishing) well, and work as you'd expect with other command-line tools.

Write.as is a text-publishing service that protects your privacy. There's no sign up required to publish, but if you do sign up, you can access posts across devices and compile collections of them in what most people would call a "blog".

## Uses

These are a few common uses for `writeas`. If you get stuck or want to know more, run `writeas [command] --help`. If you still have questions, [ask us](https://write.as/contact).

### Overview

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
     claim    Claim local unsynced posts
     blogs    List blogs
     claim    Claim local unsynced posts
     auth     Authenticate with Write.as
     logout   Log out of Write.as
     help, h  Shows a list of commands or help for one command

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
   --help, -h              show help
   --version, -V           print the version
```

#### Share something

By default, `writeas` creates a post with a `monospace` typeface that doesn't word wrap (scrolls horizontally). It will return a single line with a URL, and automatically copy that URL to the clipboard:

```bash
$ echo "Hello world!" | writeas
https://write.as/aaaazzzzzzzza
```

This is generally more useful for posting terminal output or code, like so (the `--code` flag turns on syntax highlighting):

macOS / Linux: `cat writeas/cli.go | writeas --code`

Windows: `type writeas/cli.go | writeas.exe --code`

#### Output a post

This outputs any Write.as post with the given ID.

```bash
$ writeas get aaaazzzzzzzza
Hello world!
```

#### Authenticate

This will authenticate with write.as and store the user access token locally, until you explicitly logout.
```bash
$ writeas auth username
Password: ************
```

#### List all blogs

This will output a list of the authenticated user's blogs.
```bash
$ writeas blogs
Alias    Title
user     An Example Blog
dev      My Dev Log
```

#### List posts

This lists all anonymous posts you've published. If authenticated, it will include posts on your account as well as any local / unclaimed posts.

Pass the `--url` flag to show the list with full post URLs, and the `--md` flag to return URLs with Markdown enabled.

To see post IDs with their Edit Tokens pass the `--v` flag.

```bash
$ writeas posts
aaaazzzzzzzza

$ writeas posts -url
https://write.as/aaaazzzzzzzza

$ writeas posts -v
ID              Token
aaaazzzzzzzza   dhuieoj23894jhf984hdfs9834hdf84j
```

#### Delete a post

This permanently deletes a post you own.

```bash
$ writeas delete aaaazzzzzzzza
```

#### Update a post

This completely overwrites an existing post you own.

```bash
$ echo "See you later!" | writeas update aaaazzzzzzzza
```

#### Claim a post

This moves an unsynced local post to a draft on your account. You will need to authenticate first.
```bash
$ writeas claim aaaazzzzzzzza
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

If you're publishing Markdown, supply the `--md` flag to get a URL back that will render Markdown, e.g.: `writeas new --font sans --md`
