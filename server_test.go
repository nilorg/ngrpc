package ngrpc

import "testing"

func TestGetPort(t *testing.T) {
	addr := "0.0.0.0:5000"
	t.Log(getPort(addr))
}
