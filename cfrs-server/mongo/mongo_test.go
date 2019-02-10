// +build mongodb

package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMode(t *testing.T) {
	query := Create()
	assert.Equal(t, "mongodb", query.GetMode())
}
