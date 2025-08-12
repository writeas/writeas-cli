# How to contribute

We're happy you're considering contributing to the Write.as command-line client!

It won't take long to get up to speed on this. Here are our development resources:

* We do our project management in [Phabricator](https://todo.musing.studio/tag/write.as_cli/).
* We accept and respond to bugs here on [GitHub](https://github.com/writeas/writeas-cli/issues).
* Ask any questions you have on [our forum](https://discuss.write.as).

## Testing

We try to write tests for all parts of the CLI, but aren't there yet. While not required, including tests with your new code will bring us closer to where we want to be and speed up our review.

## Submitting changes

Please send a [pull request](https://github.com/writeas/writeas-cli/compare) with a clear list of what you've done.

Please follow our coding conventions below and make sure all of your commits are atomic. Larger changes should have commits with more detailed information on what changed, any impact on existing code, rationales, etc.

## Coding conventions

We strive for consistency above all. Reading the small codebase should give you a good idea of the conventions we follow.

* We use `goimports` before committing anything
* We aim to document all exported entities
* Go files are broken up into logical functional components
* General functions are extracted into modules when possible

### Import Groups

We aim for two import groups:

* Standard library imports
* Everything else

`goimports` already does this for you along with running `go fmt`.

## Design conventions

We maintain a few high-level design principles in all decisions we make. Keep these in mind while devising new functionality:

* Updates should be backwards compatible or provide a seamless migration path from *any* previous version
* Each subcommand should perform one action and do it well
* Each subcommand will ideally work well in a script
* Avoid clever functionality and assume each function will be used in ways we didn't imagine
