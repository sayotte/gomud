package gomud

import (
	"errors"
	"fmt"
)

type Place struct {
	desc          string
	X             int64
	Y             int64
	Z             int64
	outgoingEdges []Edge
}

func NewPlace(desc string) *Place {
	return &Place{desc: desc}
}
func (p *Place) Less(op *Place) bool {
	return (p.X + p.Y + p.Z) < (op.X + op.Y + op.Z)
}
func (p *Place) Desc() string {
	return p.desc
}
func (p *Place) Edges() []Edge {
	return p.outgoingEdges
}

// An Edge is a possibly bidirectional edge in the graph of Places
type Edge interface {
	OutgoingFromPlaces() []*Place
}
type edge struct {
	a     *Place
	fromA bool
	b     *Place
	fromB bool
}

func (e edge) OutgoingFromPlaces() []*Place {
	var places []*Place
	if e.fromA {
		places = append(places, e.a)
	}
	if e.fromB {
		places = append(places, e.b)
	}
	return places
}

func NewEdge(a, b *Place, fromA, fromB bool) (Edge, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("Neither a nor b may be nil, got %p / %p", a, b)
	}
	if !fromA && !fromB {
		return nil, errors.New("One of fromA or fromB must be true.")
	}
	e := &edge{
		a:     a,
		fromA: fromA,
		b:     b,
		fromB: fromB,
	}
	if fromA {
		a.outgoingEdges = append(a.outgoingEdges, e)
	}
	if fromB {
		b.outgoingEdges = append(b.outgoingEdges, e)
	}
	return e, nil
}
