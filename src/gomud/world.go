package gomud

import (
	"io"
	"sync"
)

type World struct {
	DynamicObjects map[ObjectID]DynamicObject
	objLock        sync.RWMutex
	Places         map[PlaceID]*Place
	Edges          map[EdgeID]Edge
}

func (w *World) Load(placeStream, edgeStream, objectStream io.ReadCloser) error {
	// First load places
	// Then load edges, resolving places
	// Then load objects, resolving places/edges
	return nil
}

func NewWorld() *World {
	return &World{
		DynamicObjects: make(map[ObjectID]DynamicObject),
		Places:         make(map[PlaceID]*Place),
		Edges:          make(map[EdgeID]Edge),
	}
}
