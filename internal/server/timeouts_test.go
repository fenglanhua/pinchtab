package server

import (
	"testing"
	"time"
)

func TestServerTimeoutOrdering(t *testing.T) {
	readHeader := 10 * time.Second
	read := 30 * time.Second
	write := 60 * time.Second
	idle := 120 * time.Second

	if readHeader >= read {
		t.Errorf("ReadHeaderTimeout (%v) should be less than ReadTimeout (%v)", readHeader, read)
	}
	if read >= write {
		t.Errorf("ReadTimeout (%v) should be less than WriteTimeout (%v)", read, write)
	}
	if write >= idle {
		t.Errorf("WriteTimeout (%v) should be less than IdleTimeout (%v)", write, idle)
	}
}
