package gomud

import (
	"errors"
	"fmt"
	"sync"
)

type PlaceID uint64
type Place struct {
	ID   PlaceID
	desc string
	X    int64
	Y    int64
	Z    int64
	// This is a convenience / performance lookup list.
	// The authoritative incidence of each edge is
	// stored within the edge itself.
	outgoingEdges []*Edge
	// This is a convenience / performance lookup list.
	// The authoritative location for every object is
	// stored within the object itself.
	objects []DynamicObject
	rwlock  sync.RWMutex
}

func NewPlace(id PlaceID, desc string) *Place {
	return &Place{
		ID:   id,
		desc: desc,
	}
}
func (p *Place) Less(op *Place) bool {
	return (p.X + p.Y + p.Z) < (op.X + op.Y + op.Z)
}
func (p *Place) Desc() string {
	return p.desc
}
func (p *Place) Edges() []*Edge {
	p.rwlock.RLock()
	eList := make([]*Edge, len(p.outgoingEdges))
	copy(eList, p.outgoingEdges)
	p.rwlock.RUnlock()
	return eList
}
func (p *Place) AddObject(o DynamicObject) {
	p.rwlock.Lock()
	p.objects = append(p.objects, o)
	p.rwlock.Unlock()
}
func (p *Place) RemoveObject(o DynamicObject) {
	p.rwlock.Lock()
	i := indexOfObjectInSlice(o, p.objects)
	if i != -1 {
		p.objects = append(p.objects[:i], p.objects[i+1:]...)
	}
	p.rwlock.Unlock()
}
func (p *Place) Objects() []DynamicObject {
	p.rwlock.RLock()
	oList := make([]DynamicObject, len(p.objects))
	copy(oList, p.objects)
	p.rwlock.RUnlock()
	return oList
}
func indexOfObjectInSlice(o DynamicObject, oList []DynamicObject) int {
	for i, oli := range oList {
		if oli == o {
			return i
		}
	}
	return -1
}

type EdgeID uint64

// An Edge is a possibly bidirectional edge in the graph of Places
type Edge struct {
	ID      EdgeID
	a       *Place
	fromA   bool
	b       *Place
	fromB   bool
	objects []DynamicObject
	rwlock  sync.RWMutex
}

func (e *Edge) OutgoingFromPlaces() []*Place {
	var places []*Place
	if e.fromA {
		places = append(places, e.a)
	}
	if e.fromB {
		places = append(places, e.b)
	}
	return places
}
func (e *Edge) AddObject(o DynamicObject) {
	e.rwlock.Lock()
	e.objects = append(e.objects, o)
	e.rwlock.Unlock()
}
func (e *Edge) RemoveObject(o DynamicObject) {
	e.rwlock.Lock()
	i := indexOfObjectInSlice(o, e.objects)
	if i != -1 {
		e.objects = append(e.objects[:i], e.objects[i+1:]...)
	}
	e.rwlock.Unlock()
}

func NewEdge(a, b *Place, fromA, fromB bool) (*Edge, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("Neither a nor b may be nil, got %p / %p", a, b)
	}
	if !fromA && !fromB {
		return nil, errors.New("One of fromA or fromB must be true.")
	}
	e := &Edge{
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
