package model

type ObjectID uint64

const (
	DoNotRouteID = iota
	BroadcastID
)

type DynamicObject interface {
	ID() ObjectID
	Controller() ModelController
	SetController(ModelController)
	Edge() Edge
	setEdge(Edge)
	Place() *Place
	setPlace(*Place)
	Notify(Event)
	// serialize creates a snapshot of the object
	serialize() ([]byte, error)
	// deserialize restores the object from a snapshot
	deserialize([]byte) error
}
