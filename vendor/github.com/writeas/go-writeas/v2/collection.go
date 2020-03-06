package writeas

import (
	"fmt"
	"net/http"
)

type (
	// Collection represents a collection of posts. Blogs are a type of collection
	// on Write.as.
	Collection struct {
		Alias       string `json:"alias"`
		Title       string `json:"title"`
		Description string `json:"description"`
		StyleSheet  string `json:"style_sheet"`
		Private     bool   `json:"private"`
		Views       int64  `json:"views"`
		Domain      string `json:"domain,omitempty"`
		Email       string `json:"email,omitempty"`
		URL         string `json:"url,omitempty"`

		TotalPosts int `json:"total_posts"`

		Posts *[]Post `json:"posts,omitempty"`
	}

	// CollectionParams holds values for creating a collection.
	CollectionParams struct {
		Alias       string `json:"alias"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
	}
)

// CreateCollection creates a new collection, returning a user-friendly error
// if one comes up. Requires a Write.as subscription. See
// https://developer.write.as/docs/api/#create-a-collection
func (c *Client) CreateCollection(sp *CollectionParams) (*Collection, error) {
	p := &Collection{}
	env, err := c.post("/collections", sp, p)
	if err != nil {
		return nil, err
	}

	var ok bool
	if p, ok = env.Data.(*Collection); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	status := env.Code
	if status != http.StatusCreated {
		if status == http.StatusBadRequest {
			return nil, fmt.Errorf("Bad request: %s", env.ErrorMessage)
		} else if status == http.StatusForbidden {
			return nil, fmt.Errorf("Casual or Pro user required.")
		} else if status == http.StatusConflict {
			return nil, fmt.Errorf("Collection name is already taken.")
		} else if status == http.StatusPreconditionFailed {
			return nil, fmt.Errorf("Reached max collection quota.")
		}
		return nil, fmt.Errorf("Problem getting post: %d. %v\n", status, err)
	}
	return p, nil
}

// GetCollection retrieves a collection, returning the Collection and any error
// (in user-friendly form) that occurs. See
// https://developer.write.as/docs/api/#retrieve-a-collection
func (c *Client) GetCollection(alias string) (*Collection, error) {
	coll := &Collection{}
	env, err := c.get(fmt.Sprintf("/collections/%s", alias), coll)
	if err != nil {
		return nil, err
	}

	var ok bool
	if coll, ok = env.Data.(*Collection); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}
	status := env.Code

	if status == http.StatusOK {
		return coll, nil
	} else if status == http.StatusNotFound {
		return nil, fmt.Errorf("Collection not found.")
	} else {
		return nil, fmt.Errorf("Problem getting collection: %d. %v\n", status, err)
	}
}

// GetCollectionPosts retrieves a collection's posts, returning the Posts
// and any error (in user-friendly form) that occurs. See
// https://developer.write.as/docs/api/#retrieve-collection-posts
func (c *Client) GetCollectionPosts(alias string) (*[]Post, error) {
	coll := &Collection{}
	env, err := c.get(fmt.Sprintf("/collections/%s/posts", alias), coll)
	if err != nil {
		return nil, err
	}

	var ok bool
	if coll, ok = env.Data.(*Collection); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}
	status := env.Code

	if status == http.StatusOK {
		return coll.Posts, nil
	} else if status == http.StatusNotFound {
		return nil, fmt.Errorf("Collection not found.")
	} else {
		return nil, fmt.Errorf("Problem getting collection: %d. %v\n", status, err)
	}
}

// GetCollectionPost retrieves a post from a collection
// and any error (in user-friendly form) that occurs). See
// https://developers.write.as/docs/api/#retrieve-a-collection-post
func (c *Client) GetCollectionPost(alias, slug string) (*Post, error) {
	post := Post{}

	env, err := c.get(fmt.Sprintf("/collections/%s/posts/%s", alias, slug), &post)
	if err != nil {
		return nil, err
	}

	if _, ok := env.Data.(*Post); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	if env.Code == http.StatusOK {
		return &post, nil
	} else if env.Code == http.StatusNotFound {
		return nil, fmt.Errorf("Post %s not found in collection %s", slug, alias)
	}

	return nil, fmt.Errorf("Problem getting post %s from collection %s: %d. %v\n", slug, alias, env.Code, err)
}

// GetUserCollections retrieves the authenticated user's collections.
// See https://developers.write.as/docs/api/#retrieve-user-39-s-collections
func (c *Client) GetUserCollections() (*[]Collection, error) {
	colls := &[]Collection{}
	env, err := c.get("/me/collections", colls)
	if err != nil {
		return nil, err
	}

	var ok bool
	if colls, ok = env.Data.(*[]Collection); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}
	status := env.Code

	if status != http.StatusOK {
		if c.isNotLoggedIn(status) {
			return nil, fmt.Errorf("Not authenticated.")
		}
		return nil, fmt.Errorf("Problem getting collections: %d. %v\n", status, err)
	}
	return colls, nil
}

// DeleteCollection permanently deletes a collection and makes any posts on it
// anonymous.
//
// See https://developers.write.as/docs/api/#delete-a-collection.
func (c *Client) DeleteCollection(alias string) error {
	endpoint := "/collections/" + alias
	env, err := c.delete(endpoint, nil /* data */)
	if err != nil {
		return err
	}

	status := env.Code
	switch status {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("Not authenticated.")
	case http.StatusBadRequest:
		return fmt.Errorf("Bad request: %s", env.ErrorMessage)
	default:
		return fmt.Errorf("Problem deleting collection: %d. %s\n", status, env.ErrorMessage)
	}
}
