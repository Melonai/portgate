package portgate

import (
	"net/url"
	"path"
	"strconv"
	"strings"
)

// Path represents a routing destination the user gave.
type Path struct {
	// Identifier to which port of the destination host the path points to and to which the
	// user's request will be proxied to.
	// If there was no identifier given it's -1.
	DestinationIdentifier int
	// The path without the port identifier.
	// This is the path which will be requested from the destination.
	ResourcePath string
}

// ParsePath creates a Path from the requested URL.
func ParsePath(p string) Path {
	p = path.Clean("/" + p)

	// Get first path part, which is the destinationIdentifier

	destinationIdentifierEnd := strings.Index(p[1:], "/") + 1

	// If there is no '/' at in the path, apart from the root, the first part is the entire path
	if destinationIdentifierEnd == 0 {
		destinationIdentifierEnd = len(p)
	}

	destinationIdentifier := p[1:destinationIdentifierEnd]

	// We have to to check that destinationIdentifier is a port
	port, err := strconv.Atoi(destinationIdentifier)
	if err == nil {
		// We got an identifier and can split the path

		resourcePath := path.Clean("/" + p[destinationIdentifierEnd:])
		return Path{
			DestinationIdentifier: port,
			ResourcePath:          resourcePath,
		}
	} else {
		// We got some other path without an identifier

		return Path{
			DestinationIdentifier: -1,
			ResourcePath:          p,
		}
	}
}

// ParsePathFromReferer tries to create a Path from the Referer header of the request.
func ParsePathFromReferer(p Path, r string) (Path, error) {
	u, err := url.Parse(r)
	if err != nil {
		return Path{}, err
	}

	// p has the correct resource path but the wrong port, so we create a new path
	// with the correct data from both.
	rp := ParsePath(u.Path)
	return Path{
		DestinationIdentifier: rp.DestinationIdentifier,
		ResourcePath:          p.ResourcePath,
	}, nil
}
