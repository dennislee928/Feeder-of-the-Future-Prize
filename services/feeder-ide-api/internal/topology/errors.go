package topology

import "errors"

var (
	ErrTopologyNotFound = errors.New("topology not found")
	ErrInvalidTopology  = errors.New("invalid topology")
)

