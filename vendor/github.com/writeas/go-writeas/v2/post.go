package writeas

import (
	"fmt"
	"net/http"
	"time"
)

type (
	// Post represents a published Write.as post, whether anonymous, owned by a
	// user, or part of a collection.
	Post struct {
		ID        string    `json:"id"`
		Slug      string    `json:"slug"`
		Token     string    `json:"token"`
		Font      string    `json:"appearance"`
		Language  *string   `json:"language"`
		RTL       *bool     `json:"rtl"`
		Listed    bool      `json:"listed"`
		Created   time.Time `json:"created"`
		Updated   time.Time `json:"updated"`
		Title     string    `json:"title"`
		Content   string    `json:"body"`
		Views     int64     `json:"views"`
		Tags      []string  `json:"tags"`
		Images    []string  `json:"images"`
		OwnerName string    `json:"owner,omitempty"`

		Collection *Collection `json:"collection,omitempty"`
	}

	// OwnedPostParams are, together, fields only the original post author knows.
	OwnedPostParams struct {
		ID    string `json:"id"`
		Token string `json:"token,omitempty"`
	}

	// PostParams holds values for creating or updating a post.
	PostParams struct {
		// Parameters only for updating
		ID    string `json:"-"`
		Token string `json:"token,omitempty"`

		// Parameters for creating or updating
		Slug     string     `json:"slug"`
		Created  *time.Time `json:"created,omitempty"`
		Updated  *time.Time `json:"updated,omitempty"`
		Title    string     `json:"title,omitempty"`
		Content  string     `json:"body,omitempty"`
		Font     string     `json:"font,omitempty"`
		IsRTL    *bool      `json:"rtl,omitempty"`
		Language *string    `json:"lang,omitempty"`

		// Parameters only for creating
		Crosspost []map[string]string `json:"crosspost,omitempty"`

		// Parameters for collection posts
		Collection string `json:"-"`
	}

	// PinnedPostParams holds values for pinning a post
	PinnedPostParams struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}

	// BatchPostResult contains the post-specific result as part of a larger
	// batch operation.
	BatchPostResult struct {
		ID           string `json:"id,omitempty"`
		Code         int    `json:"code,omitempty"`
		ErrorMessage string `json:"error_msg,omitempty"`
	}

	// ClaimPostResult contains the post-specific result for a request to
	// associate a post to an account.
	ClaimPostResult struct {
		ID           string `json:"id,omitempty"`
		Code         int    `json:"code,omitempty"`
		ErrorMessage string `json:"error_msg,omitempty"`
		Post         *Post  `json:"post,omitempty"`
	}
)

// GetPost retrieves a published post, returning the Post and any error (in
// user-friendly form) that occurs. See
// https://developer.write.as/docs/api/#retrieve-a-post.
func (c *Client) GetPost(id string) (*Post, error) {
	p := &Post{}
	env, err := c.get(fmt.Sprintf("/posts/%s", id), p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*Post); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}
	status := env.Code

	if status == http.StatusOK {
		return p, nil
	} else if status == http.StatusNotFound {
		return nil, fmt.Errorf("Post not found.")
	} else if status == http.StatusGone {
		return nil, fmt.Errorf("Post unpublished.")
	}
	return nil, fmt.Errorf("Problem getting post: %d. %s\n", status, env.ErrorMessage)
}

// CreatePost publishes a new post, returning a user-friendly error if one comes
// up. See https://developer.write.as/docs/api/#publish-a-post.
func (c *Client) CreatePost(sp *PostParams) (*Post, error) {
	p := &Post{}
	endPre := ""
	if sp.Collection != "" {
		endPre = "/collections/" + sp.Collection
	}
	env, err := c.post(endPre+"/posts", sp, p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*Post); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	status := env.Code
	if status != http.StatusCreated {
		if status == http.StatusBadRequest {
			return nil, fmt.Errorf("Bad request: %s", env.ErrorMessage)
		}
		return nil, fmt.Errorf("Problem creating post: %d. %s\n", status, env.ErrorMessage)
	}
	return p, nil
}

// UpdatePost updates a published post with the given PostParams. See
// https://developer.write.as/docs/api/#update-a-post.
func (c *Client) UpdatePost(id, token string, sp *PostParams) (*Post, error) {
	return c.updatePost("", id, token, sp)
}

func (c *Client) updatePost(collection, identifier, token string, sp *PostParams) (*Post, error) {
	p := &Post{}
	endpoint := "/posts/" + identifier
	/*
		if collection != "" {
			endpoint = "/collections/" + collection + endpoint
		} else {
			sp.Token = token
		}
	*/
	sp.Token = token
	env, err := c.put(endpoint, sp, p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*Post); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	status := env.Code
	if status != http.StatusOK {
		if c.isNotLoggedIn(status) {
			return nil, fmt.Errorf("Not authenticated.")
		} else if status == http.StatusBadRequest {
			return nil, fmt.Errorf("Bad request: %s", env.ErrorMessage)
		}
		return nil, fmt.Errorf("Problem updating post: %d. %s\n", status, env.ErrorMessage)
	}
	return p, nil
}

// DeletePost permanently deletes a published post. See
// https://developer.write.as/docs/api/#delete-a-post.
func (c *Client) DeletePost(id, token string) error {
	return c.deletePost("", id, token)
}

func (c *Client) deletePost(collection, identifier, token string) error {
	p := map[string]string{}
	endpoint := "/posts/" + identifier
	/*
		if collection != "" {
			endpoint = "/collections/" + collection + endpoint
		} else {
			p["token"] = token
		}
	*/
	p["token"] = token
	env, err := c.delete(endpoint, p)
	if err != nil {
		return err
	}

	status := env.Code
	if status == http.StatusNoContent {
		return nil
	} else if c.isNotLoggedIn(status) {
		return fmt.Errorf("Not authenticated.")
	} else if status == http.StatusBadRequest {
		return fmt.Errorf("Bad request: %s", env.ErrorMessage)
	}
	return fmt.Errorf("Problem deleting post: %d. %s\n", status, env.ErrorMessage)
}

// ClaimPosts associates anonymous posts with a user / account.
// https://developer.write.as/docs/api/#claim-posts.
func (c *Client) ClaimPosts(sp *[]OwnedPostParams) (*[]ClaimPostResult, error) {
	p := &[]ClaimPostResult{}
	env, err := c.post("/posts/claim", sp, p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*[]ClaimPostResult); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	status := env.Code
	if status == http.StatusOK {
		return p, nil
	} else if c.isNotLoggedIn(status) {
		return nil, fmt.Errorf("Not authenticated.")
	} else if status == http.StatusBadRequest {
		return nil, fmt.Errorf("Bad request: %s", env.ErrorMessage)
	} else {
		return nil, fmt.Errorf("Problem claiming post: %d. %s\n", status, env.ErrorMessage)
	}
	// TODO: does this also happen with moving posts?
}

// GetUserPosts retrieves the authenticated user's posts.
// See https://developers.write.as/docs/api/#retrieve-user-39-s-posts
func (c *Client) GetUserPosts() (*[]Post, error) {
	p := &[]Post{}
	env, err := c.get("/me/posts", p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*[]Post); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}
	status := env.Code

	if status != http.StatusOK {
		if c.isNotLoggedIn(status) {
			return nil, fmt.Errorf("Not authenticated.")
		}
		return nil, fmt.Errorf("Problem getting user posts: %d. %s\n", status, env.ErrorMessage)
	}
	return p, nil
}

// PinPost pins a post in the given collection.
// See https://developers.write.as/docs/api/#pin-a-post-to-a-collection
func (c *Client) PinPost(alias string, pp *PinnedPostParams) error {
	res := &[]BatchPostResult{}
	env, err := c.post(fmt.Sprintf("/collections/%s/pin", alias), []*PinnedPostParams{pp}, res)
	if err != nil {
		return err
	}

	var ok bool
	if res, ok = env.Data.(*[]BatchPostResult); !ok {
		return fmt.Errorf("Wrong data returned from API.")
	}

	// Check for basic request errors on top level response
	status := env.Code
	if status != http.StatusOK {
		if c.isNotLoggedIn(status) {
			return fmt.Errorf("Not authenticated.")
		}
		return fmt.Errorf("Problem pinning post: %d. %s\n", status, env.ErrorMessage)
	}

	// Check the individual post result
	if len(*res) == 0 || len(*res) > 1 {
		return fmt.Errorf("Wrong data returned from API.")
	}
	if (*res)[0].Code != http.StatusOK {
		return fmt.Errorf("Problem pinning post: %d", (*res)[0].Code)
		// TODO: return ErrorMessage (right now it'll be empty)
		// return fmt.Errorf("Problem pinning post: %s", res[0].ErrorMessage)
	}
	return nil
}

// UnpinPost unpins a post from the given collection.
// See https://developers.write.as/docs/api/#unpin-a-post-from-a-collection
func (c *Client) UnpinPost(alias string, pp *PinnedPostParams) error {
	res := &[]BatchPostResult{}
	env, err := c.post(fmt.Sprintf("/collections/%s/unpin", alias), []*PinnedPostParams{pp}, res)
	if err != nil {
		return err
	}

	var ok bool
	if res, ok = env.Data.(*[]BatchPostResult); !ok {
		return fmt.Errorf("Wrong data returned from API.")
	}

	// Check for basic request errors on top level response
	status := env.Code
	if status != http.StatusOK {
		if c.isNotLoggedIn(status) {
			return fmt.Errorf("Not authenticated.")
		}
		return fmt.Errorf("Problem unpinning post: %d. %s\n", status, env.ErrorMessage)
	}

	// Check the individual post result
	if len(*res) == 0 || len(*res) > 1 {
		return fmt.Errorf("Wrong data returned from API.")
	}
	if (*res)[0].Code != http.StatusOK {
		return fmt.Errorf("Problem unpinning post: %d", (*res)[0].Code)
		// TODO: return ErrorMessage (right now it'll be empty)
		// return fmt.Errorf("Problem unpinning post: %s", res[0].ErrorMessage)
	}
	return nil
}
