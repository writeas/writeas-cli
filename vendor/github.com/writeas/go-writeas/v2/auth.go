package writeas

import (
	"fmt"
	"net/http"
)

// LogIn authenticates a user with Write.as.
// See https://developer.write.as/docs/api/#authenticate-a-user
func (c *Client) LogIn(username, pass string) (*AuthUser, error) {
	u := &AuthUser{}
	up := struct {
		Alias string `json:"alias"`
		Pass  string `json:"pass"`
	}{
		Alias: username,
		Pass:  pass,
	}

	env, err := c.post("/auth/login", up, u)
	if err != nil {
		return nil, err
	}

	var ok bool
	if u, ok = env.Data.(*AuthUser); !ok {
		return nil, fmt.Errorf("Wrong data returned from API.")
	}

	status := env.Code
	if status != http.StatusOK {
		if status == http.StatusBadRequest {
			return nil, fmt.Errorf("Bad request: %s", env.ErrorMessage)
		} else if status == http.StatusUnauthorized {
			return nil, fmt.Errorf("Incorrect password.")
		} else if status == http.StatusNotFound {
			return nil, fmt.Errorf("User does not exist.")
		} else if status == http.StatusTooManyRequests {
			return nil, fmt.Errorf("Too many log in attempts in a short period of time.")
		}
		return nil, fmt.Errorf("Problem authenticating: %d. %v\n", status, err)
	}

	c.SetToken(u.AccessToken)
	return u, nil
}

// LogOut logs the current user out, making the Client's current access token
// invalid.
func (c *Client) LogOut() error {
	env, err := c.delete("/auth/me", nil)
	if err != nil {
		return err
	}

	status := env.Code
	if status != http.StatusNoContent {
		if status == http.StatusNotFound {
			return fmt.Errorf("Access token is invalid or doesn't exist")
		}
		return fmt.Errorf("Unable to log out: %v", env.ErrorMessage)
	}

	// Logout successful, so update the Client
	c.token = ""

	return nil
}

func (c *Client) isNotLoggedIn(code int) bool {
	if c.token == "" {
		return false
	}
	return code == http.StatusUnauthorized
}
