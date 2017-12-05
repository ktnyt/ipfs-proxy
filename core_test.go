package proxy

import (
	"testing"
)

func TestCore(t *testing.T) {
	p, err := NewLocalProxy("/ip4/127.0.0.1/tcp/5001")

	if err != nil {
		t.Error(err)
	}

	defer p.Cancel()

	data := []byte("Hello, world!")

	go p.Spin()

	for len(p.Msgs) == 0 {
		select {
		case err := <-p.Comm:
			t.Error(err)
		default:
			if err := p.Ping(data); err != nil {
				t.Error(err)
			}
		}
	}
}
