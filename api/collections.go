package api

// RemoteColl represents a collection of posts
// It is a reduced set of data from a go-writeas Collection
type RemoteColl struct {
	Alias string
	Title string
	URL   string
}
