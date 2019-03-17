package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimulatorDial(t *testing.T) {
	query := Dial("simulator")
	assert.NotNil(t, query)
}

func TestSimulatorDialWrongMode(t *testing.T) {
	defer func () {
		assert.NotNil(t, recover())
	}()
	Dial("xxx")
}


func TestSimulatorReset(t *testing.T) {
	query := Dial("simulator")
	assert.NotNil(t, query)

	result := query.ResetAll()

	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Count)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Message)
}

func TestSimulatorRecordRequest(t *testing.T) {
	query := Dial("simulator")
	assert.NotNil(t, query)

	result := query.RecordRequest(&RequestData{})

	assert.NotNil(t, result)
	assert.True(t, result.Count > int64(0))
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Message)
}
