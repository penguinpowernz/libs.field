package parsers

import (
	"context"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
)

var Parsers = map[string]ParserFunc{
	"commits":  Commits,
	"repos":    Repos,
	"releases": Releases,
	"tags":     Tags,
	"search":   Search,
	"commit":   Commit,
}

type ParserFunc func(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error
