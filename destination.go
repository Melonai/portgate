package portgate

import (
	"net/url"
	"path"
	"strconv"
	"strings"
)

// Destination represents a routing destination the user gave.
type Destination struct {
	// Identifier to which port of the destination host the path points to and to which the
	// user's request will be proxied to.
	Port int
	// The path without the port identifier.
	// This is the path which will be requested from the destination.
	Path string
	IsPortgatePath bool
}

// ParseDestinationFromURL creates a Path from the requested URL.
func DestinationFromURL(p string) Destination {
	p = path.Clean("/" + p)

	// Get first path part, which is the potential port identifier.
	destinationIdentifierEnd := strings.Index(p[1:], "/") + 1

	// If there is no '/' in the path, apart from the root, the first part is the entire path
	if destinationIdentifierEnd == 0 {
		destinationIdentifierEnd = len(p)
	}

	destinationIdentifier := p[1:destinationIdentifierEnd]

	// We have to to check that destinationIdentifier is a port
	port, err := strconv.Atoi(destinationIdentifier)
	if err == nil {
		// We got an identifier and can split the path
		resourcePath := path.Clean("/" + p[destinationIdentifierEnd:])

		return Destination{
			Port: port,
			Path: resourcePath,
		}
	} else {
		destination := Destination{
			Port: 0,
			Path: p,
		}

		if strings.HasPrefix(destination.Path, "/_portgate") {
			destination.IsPortgatePath = true
			return destination
		}

		return destination
	}
}

// DestinationFromReferer tries to create a Path from the Referer header of the request.
func (d Destination) AddReferer(referer string) Destination {
	u, err := url.Parse(referer)
	if err != nil {
		return Destination{}
	}

	// d has the correct resource path but the wrong port, so we create a new destination
	// with the correct data from both.
	newDestination := DestinationFromURL(u.Path)

	return Destination{
		Port: newDestination.Port,
		Path: d.Path,
	}
}