// +build mongodb

package mongo

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	query := Dial("mongodb")
	assert.NotNil(t, query)
}

func TestRecordRequest(t *testing.T) {
	query := Dial("mongodb")
	assert.NotNil(t, query)

	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	requestData := RequestData{
		Method:          r.Method,
		Remote:          r.RemoteAddr,
		Timestamp:       time.Now(),
		Url:             r.URL.String(),
		XB3ParentSpanId: r.Header.Get("x-b3-parentspanid"),
		XB3SpanId:       r.Header.Get("x-b3-spanid"),
		XB3TraceId:      r.Header.Get("x-b3-traceid"),
		XForwardedFor:   r.Header.Get("x-forwarded-for"),
		Tag:             "test",
	}

	result := query.RecordRequest(&requestData)
	assert.NotNil(t, result)
}

func TestResetRequest(t *testing.T) {
	query := Dial("mongodb")
	assert.NotNil(t, query)

	result := query.ResetAll()

	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Count)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Message)
}
