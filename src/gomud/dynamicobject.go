package gomud

// This is likely to change, later on, to something like a UUID
// that can be generated safely without consulting a global state.
type ObjectID uint64

const (
	NonObjectID = iota
	DoNotRouteID
	BroadcastID
)

type DynamicObject interface {
	ID() ObjectID
	Controller() ModelController
	SetController(ModelController)
	Edge() *Edge
	setEdge(*Edge)
	Place() *Place
	setPlace(*Place)
	Notify(Event)
	//// serialize creates a snapshot of the object
	//MarshalJSON() ([]byte, error)
	//// deserialize restores the object from a snapshot
	//UnmarshalJSON([]byte) error
}
