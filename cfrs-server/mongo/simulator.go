package mongo

import (
	"log"
	"sync"
)

type simulator struct {
	requestCounter uint64
	mutex          *sync.Mutex
}

var query MongoDBQuery

func CreateSimulator() MongoDBQuery {
	if (query == nil) {
		query = &simulator{0, &sync.Mutex{}}
	}

	log.Printf("MongoDB connected: %v", query.GetMode())
	return query
}

func (sim *simulator) GetMode() string {
	return "simulator"
}

func (sim *simulator) GetRequestCounter() uint64 {
	return sim.requestCounter
}

func (sim *simulator) IncRequestCounter() {
	sim.mutex.Lock()
	sim.requestCounter++
	sim.mutex.Unlock()
}

func (sim *simulator) ResetRequestCounter() {
	sim.mutex.Lock()
	sim.requestCounter = 0
	sim.mutex.Unlock()
}
