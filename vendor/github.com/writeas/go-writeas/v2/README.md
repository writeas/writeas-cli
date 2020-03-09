# go-writeas

[![godoc](https://godoc.org/go.code.as/writeas.v2?status.svg)](https://godoc.org/go.code.as/writeas.v2)

Official Write.as Go client library.

## Installation

**Warning**: the `v2` branch is under heavy development and its API will change without notice.

For a stable API, use `go.code.as/writeas.v1` and upgrade to `v2` once everything is merged into `master`.

```bash
go get go.code.as/writeas.v2
```

## Documentation

See all functionality and usages in the [API documentation](https://developer.write.as/docs/api/).

### Example usage

```go
import "go.code.as/writeas.v2"

func main() {
	// Create the client
	c := writeas.NewClient()

	// Publish a post
	p, err := c.CreatePost(&writeas.PostParams{
		Title:   "Title!",
		Content: "This is a post.",
		Font:    "sans",
	})
	if err != nil {
		// Perhaps show err.Error()
	}

	// Save token for later, since it won't ever be returned again
	token := p.Token

	// Update a published post
	p, err = c.UpdatePost(p.ID, token, &writeas.PostParams{
		Content: "Now it's been updated!",
	})
	if err != nil {
		// handle
	}

	// Get a published post
	p, err = c.GetPost(p.ID)
	if err != nil {
		// handle
	}

	// Delete a post
	err = c.DeletePost(p.ID, token)
}
```

## Contributing

The library covers our usage, but might not be comprehensive of the API. So we always welcome contributions and improvements from the community. Before sending pull requests, make sure you've done the following:

* Run `goimports` on all updated .go files.
* Document all exported structs and funcs.

## License

MIT
