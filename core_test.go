package proxy

import (
	"testing"
)

func TestCore(t *testing.T) {
	p, err := NewLocalProxy("/ip4/127.0.0.1/tcp/5001")
	defer p.Cancel()

	if err != nil {
		t.Error(err)
	}

	c := make(chan error)

	go p.Spin(c)

	data := []byte("Hello, world!")

	for len(p.Msgs) == 0 {
		select {
		case err := <-c:
			t.Error(err)
		default:
			if err := p.Ping(data); err != nil {
				t.Error(err)
			}
		}
	}
}
