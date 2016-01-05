package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPixiecoreHandler(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{validMACStr: testMachine},
	}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - machine config is rendered as Pixiecore JSON
	expectedJSON := `{"kernel":"/image/kernel","initrd":["/image/initrd_a","/image/initrd_b"],"cmdline":{"a":"b","c":""}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestPixiecoreHandler_InvalidMACAddress(t *testing.T) {
	store := &emptyStore{}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid MAC address /\n", w.Body.String())
}

func TestPixiecoreHandler_NoMatchingSpec(t *testing.T) {
	store := &emptyStore{}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
