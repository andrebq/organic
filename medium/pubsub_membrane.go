package medium

import (
	"context"

	"github.com/andrebq/organic/cell"
	"github.com/andrebq/organic/helper/signal"
	"github.com/go-redis/redis/v8"
)

type (
	pubsubMembrane struct {
		ps  *redis.PubSub
		cli *redis.Client
		ctx context.Context

		id    cell.ID
		idstr string
	}

	crossSignal struct {
		sig cell.Signal
		age int
	}
)

func newMembrane(ctx context.Context, ps *redis.PubSub, cli *redis.Client, id cell.ID) *pubsubMembrane {
	return &pubsubMembrane{
		ps:    ps,
		cli:   cli,
		ctx:   ctx,
		id:    id,
		idstr: id.String(),
	}
}

// ExchangeOut sends the given signal to a cell identified with the given CellID
// and to the provided Receptor
func (ps *pubsubMembrane) ExchangeOut(c cell.Carrier) {
	if c.From != ps.id {
		return
	}
	b := signal.Build().Bytes(c.From.Bytes()).Bytes(c.To.Bytes()).Str(string(c.Receptor)).Bytes(c.Signal.Content()).Signal()
	ps.cli.Publish(ps.ctx, c.To.String(), b.Content())
}

// Receives a signal one one of the given receptors
func (ps *pubsubMembrane) ExchangeIn(c *cell.Carrier) {
	msg, err := ps.ps.ReceiveMessage(ps.ctx)
	if err != nil {
		// TODO think about how to handle errors here...
		return
	}
	if msg.Channel != ps.idstr {
		return
	}
	sig, err := cell.RawSignalStr(msg.Payload)
	var buf []byte
	m := signal.Match(sig, signal.ID(&c.From), signal.ID(&c.To), signal.Receptor(&c.Receptor), signal.Bytes(&buf))
	if !m {
		*c = cell.Carrier{}
		return
	}
	sig, err = cell.RawSignal(buf)
	if err != nil {
		*c = cell.Carrier{}
	}
	c.Signal = sig
	return
}
