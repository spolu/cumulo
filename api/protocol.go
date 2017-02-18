package api

const (
	// DefaultPort is the default port on which the API runs.
	DefaultPort int = 9000
	// Version is the current protocol version.
	Version string = "0"
	// TimeResolutionNs is the resolution of our time variables in nanoseconds
	// (aka resolution in milliseconds).
	TimeResolutionNs int64 = 1000 * 1000
)

// UserResource is the representation of a user in the Cumulo API.
type UserResource struct {
	Token   string `json:"token"`
	Created int64  `json:"created"`

	Phone  string  `json:"phone"`
	PubKey *string `json:"pubkey"`
}
