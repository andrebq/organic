package medium

import (
	"os"

	"github.com/andrebq/organic/cell"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type (
	// Agar is the default medium which can be used by cells to
	// interact with the environment
	Agar struct {
		ctx context.Context
		cli *redis.Client
	}
)

// NewAgar connects to the given redis server to provide access
// to Membranes
func NewAgar(ctx context.Context, addr string) (*Agar, error) {
	if addr == "" {
		addr = os.Getenv("ORGANIC_AGAR_REDIS")
	}
	if addr == "" {
		addr = os.Getenv("REDIS_SERVER")
	}
	if addr == "" {
		addr = "localhost:6379"
	}
	r := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := r.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &Agar{ctx: ctx, cli: r}, nil
}

// Membrane returns a new membrane which can be used by the given cell
func (a *Agar) Membrane(id cell.ID) cell.Membrane {
	pubsub := a.cli.Subscribe(a.ctx, id.String())
	return newMembrane(a.ctx, pubsub, a.cli, id)
}

// Close disconnects the agar from redis
func (a *Agar) Close() error {
	return a.cli.Close()
}
