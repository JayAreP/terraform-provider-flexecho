package flexecho

import sdk "github.com/silk-us/flexecho-go-sdk/flexecho"

// per-provider config, server + bearer token
type Config struct {
	Server string
	Token  string
}

// used to validate the server looked like an ip/host here before connecting,
// dropped it - the sdk surfaces a clean enough error on a bad host anyway
// func (c *Config) Client() (*sdk.Credentials, error) {
// 	if c.Server == "" {
// 		return nil, fmt.Errorf("server is required")
// 	}
// 	return sdk.Connect(c.Server, c.Token), nil
// }

// build a configured echo sdk client. NOthing fancy, just wraps Connect
func (c *Config) Client() (*sdk.Credentials, error) {
	return sdk.Connect(c.Server, c.Token), nil
}
