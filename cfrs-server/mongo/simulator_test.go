package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimulatorCreate(t *testing.T) {
	query := CreateSimulator()
	assert.NotNil(t, query)
}

func TestSimulatorGetMode(t *testing.T) {
	query := CreateSimulator()
	assert.Equal(t, "simulator", query.GetMode())
}

func TestSimulatorCounter(t *testing.T) {
	query := CreateSimulator()
	assert.Equal(t, uint64(0), query.GetRequestCounter())
	query.IncRequestCounter()
	assert.Equal(t, uint64(1), query.GetRequestCounter())
	query.ResetRequestCounter()
	assert.Equal(t, uint64(0), query.GetRequestCounter())
}
