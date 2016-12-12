package gomud

import (
	"encoding/json"
)

type Slime struct {
	ctrl  Controller
	id    int
	state SlimeState
	// location Location
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
func (s *Slime) Location() {
	return
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
