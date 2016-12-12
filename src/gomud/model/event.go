package model

import "time"

type Event interface {
	ObjectID() ObjectID
}

type TimeTick struct {
	OID  ObjectID
	Type string
	When time.Time
}

func (tt TimeTick) ObjectID() ObjectID {
	return tt.ObjectID()
}
func NewTimeTick(when time.Time) TimeTick {
	return TimeTick{
		OID:  BroadcastID,
		Type: "timetick",
		When: when,
	}
}

type PoisonPill struct{}

func (pp PoisonPill) ObjectID() ObjectID {
	return DoNotRouteID
}

type SetPlace struct {
	OID     ObjectID
	object  DynamicObject
	Type    string
	PlaceId PlaceID
	place   *Place
}

func (sp SetPlace) ObjectID() ObjectID {
	return sp.OID
}
func NewSetPlace(object DynamicObject, place *Place) SetPlace {
	return SetPlace{
		OID:     object.ID(),
		object:  object,
		Type:    "setplace",
		PlaceId: place.ID,
		place:   place,
	}
}

type SetEdge struct {
	OID         ObjectID
	object      DynamicObject
	Type        string
	EID         EdgeID
	edge        Edge
	FromPlaceID PlaceID
	fromPlace   *Place
}

func (se SetEdge) ObjectID() ObjectID {
	return se.OID
}
func NewSetEdge(object DynamicObject, edge Edge, fromPlace *Place) SetEdge {
	return SetEdge{
		OID:         object.ID(),
		object:      object,
		Type:        "setedge",
		EID:         edge.ID(),
		edge:        edge,
		FromPlaceID: fromPlace.ID,
		fromPlace:   fromPlace,
	}
}

type InsertObject struct {
	NewObject DynamicObject
}

func (io InsertObject) ObjectID() ObjectID {
	return io.NewObject.ID()
}

type InsertPlace struct {
	NewPlace *Place
}

func (ip InsertPlace) ObjectID() ObjectID {
	return NonObjectID
}

type InsertEdge struct {
	NewEdge *edge
}

func (ie InsertEdge) ObjectID() ObjectID {
	return NonObjectID
}
