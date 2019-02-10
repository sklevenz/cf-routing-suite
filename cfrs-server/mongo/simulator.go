package mongo

import "sync"

type simulator struct {
	requestCounter uint64
	mutex          *sync.Mutex
}

func CreateSimulator() MongoDBQuery {
	return &simulator{0, &sync.Mutex{}}
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
