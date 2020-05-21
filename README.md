# organic - Nature inspired resilient systems

Organic tries to bridge the gap between the original object-oriented programming
where messages and late-binding where the building blocks and network systems
where all messages are async.

Similar to the original OO concept, the building blocks (cells), take inspiration
from biology.

The network is based around the actor model where messages (signals) happens
by addressing a given actor (cell) and a given port (receptor). State internal
to the cells are opaque to the outside world and can only be mutated by the
cell itself.

Each cell is addressable by a unique 256 bit identifier which is the sha256 of
a unique string (the cell "name"). Cells exchange messages via their "mebrane"
which is an abstract medium on which CBOR/JSON encoded messages are transfered.

The membrane provides the given guarantees:

- signals are delivered in a best-effort fashion
- if the same signal is sent more than once it might be delivered more than once
- signals might arrive out of order
- there is no acknowledge when a signal is delivered to the cell
- signals "might" be queued up
- if a signal is sent and the cell doesn't have a receptor, the signal might be lost
- if multiple cells have the same ID the signal might be detected by more than one

# Direct inspirations

- Alan Kay concept of object-oriented programming on which the most important aspects
are: isolated state, message passing and late-binding.
- Carl Heweit concept of actors communicating using async messages over an opaque medium
- Joe Armstrong for creating Erlang and showing that such resilient and complex
systems are acheivable using simple and unreliable processes
- Biology itself (although I'm not a Biologist)
