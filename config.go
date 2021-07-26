package portgate

import "fmt"

// Config represents the global Portgate config.
type Config struct {
	// Where Portgate will be running at.
	portgatePort int
	portgateHost string

	// Where the requests will be proxied to.
	targetHost string

	allowedPorts   []int
	forbiddenPorts []int

	key string

	jwtSecret string
}

// GetConfig creates the Portgate config from outside sources such as
// the environment variables and the portgate.yml file.
func GetConfig() (Config, error) {
	// TODO: Read config from environment/file
	return Config{
		portgatePort: 8080,
		targetHost:   "localhost",

		allowedPorts:   []int{80},
		forbiddenPorts: []int{},

		key: "password",
	}, nil
}

// PortgateAddress is the address on which Portgate will run.
func (c *Config) PortgateAddress() string {
	return fmt.Sprintf("%s:%d", c.portgateHost, c.portgatePort)
}

// TargetAddress is the address of the destination server.
func (c *Config) TargetAddress(port int) string {
	return fmt.Sprintf("%s:%d", c.targetHost, port)
}

// MakeUrl creates the URL on the destination host that the user wants to access.
func (c *Config) MakeUrl(p Path) string {
	// TODO: Figure out what to do with TLS
	return fmt.Sprintf("http://%s:%d%s", c.targetHost, p.DestinationIdentifier, p.ResourcePath)
}

// CheckKey checks whether the givenKey matches the one in the config.
func (c *Config) CheckKey(givenKey string) bool {
	return c.key == givenKey
}
