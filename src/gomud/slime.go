package gomud

import (
	"encoding/json"
)

type Slime struct {
	ctrl  Controller
	id    int
	state SlimeState
	place *Place
	edge  *Edge
}

func (s *Slime) Controller() Controller {
	return s.ctrl
}
func (s *Slime) SetController(c Controller) {
	if s.ctrl != nil {
		s.ctrl.Stop()
	}
	s.ctrl = c
}
func (s *Slime) Edge() *Edge {
	return s.edge
}
func (s *Slime) SetEdge(e *Edge) {
	s.place = nil
	s.edge = e
}
func (s *Slime) Place() *Place {
	return s.place
}
func (s *Slime) SetPlace(p *Place) {
	s.edge = nil
	s.place = p
}
func (s *Slime) Notify(n Notification) {
	if s.ctrl != nil {
		s.ctrl.Notify(n)
	}
}
func (s *Slime) State() DynamicObjectState {
	return &s.state
}
func NewSlime(id int, ss SlimeState) *Slime {
	return &Slime{
		id:    id,
		state: ss,
	}
}

type SlimeState struct {
	Size int
}

func (ss *SlimeState) Serialize() ([]byte, error) {
	return json.Marshal(ss)
}
func (ss *SlimeState) Deserialize(b []byte) error {
	return json.Unmarshal(b, ss)
}
