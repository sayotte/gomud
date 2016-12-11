package model

import (
	"encoding/json"
	"sync"
)

type Slime struct {
	ctrl   ModelController
	id     ObjectID
	size   int
	place  *Place
	edge   Edge
	rwlock sync.RWMutex
}

func (s *Slime) ID() ObjectID {
	return s.id
}
func (s *Slime) Controller() ModelController {
	return s.ctrl
}
func (s *Slime) SetController(c ModelController) {
	if s.ctrl != nil {
		s.ctrl.Stop()
	}
	s.ctrl = c
}
func (s *Slime) Edge() Edge {
	s.rwlock.RLock()
	e := s.edge
	s.rwlock.Unlock()
	return e
}
func (s *Slime) setEdge(e Edge) {
	s.rwlock.Lock()
	if s.place != nil {
		s.place.RemoveObject(s)
	}
	s.place = nil
	s.edge = e
	e.AddObject(s)
	s.rwlock.Unlock()
}
func (s *Slime) Place() *Place {
	s.rwlock.RLock()
	p := s.place
	s.rwlock.RUnlock()
	return p
}
func (s *Slime) setPlace(p *Place) {
	s.rwlock.Lock()
	if s.edge != nil {
		s.edge.RemoveObject(s)
	}
	s.edge = nil
	s.place = p
	p.AddObject(s)
	s.rwlock.Unlock()
}
func (s *Slime) Notify(e Event) {
	if s.ctrl != nil {
		s.ctrl.Notify(e)
	}
}
func (s *Slime) serialize() ([]byte, error) {
	return json.Marshal(s)
}
func (s *Slime) deserialize(b []byte) error {
	return json.Unmarshal(b, s)
}
func NewSlime(id ObjectID, size int) *Slime {
	return &Slime{
		id:   id,
		size: size,
	}
}
