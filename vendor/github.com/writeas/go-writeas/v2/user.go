package writeas

import "time"

type (
	// AuthUser represents a just-authenticated user. It contains information
	// that'll only be returned once (now) per user session.
	AuthUser struct {
		AccessToken string `json:"access_token,omitempty"`
		Password    string `json:"password,omitempty"`
		User        *User  `json:"user"`
	}

	// User represents a registered Write.as user.
	User struct {
		Username string    `json:"username"`
		Email    string    `json:"email"`
		Created  time.Time `json:"created"`

		// Optional properties
		Subscription *UserSubscription `json:"subscription"`
	}

	// UserSubscription contains information about a user's Write.as
	// subscription.
	UserSubscription struct {
		Name       string    `json:"name"`
		Begin      time.Time `json:"begin"`
		End        time.Time `json:"end"`
		AutoRenew  bool      `json:"auto_renew"`
		Active     bool      `json:"is_active"`
		Delinquent bool      `json:"is_delinquent"`
	}
)
