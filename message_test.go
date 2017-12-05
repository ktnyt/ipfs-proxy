package proxy

import (
	"testing"
	"time"

	"github.com/ipfs/go-ipfs-api"
)

func init() {
	ipfs = shell.NewLocalShell()
	if err := setUser(); err != nil {
		panic(err)
	}
}

func TestNewMessage(t *testing.T) {
	data := "Hello, world!"
	before := time.Now()
	m := NewMessage([]byte(data))
	after := time.Now()

	if before.After(m.Time()) {
		t.Errorf("timestamp %s is before %s: should be after", m.Time(), before)
	}

	if after.Before(m.Time()) {
		t.Errorf("timestamp %s is after %s: should be before", m.Time(), after)
	}

	if m.From() != user {
		t.Errorf("owner '%s' does not match user '%s'", m.From(), user)
	}

	if string(m.Data()) != data {
		t.Errorf("data '%s' does not match '%s'", m.Data(), data)
	}
}

func TestMessageMarshalUnmarshalBinary(t *testing.T) {
	data := "Hello, world!"
	m0 := NewMessage([]byte(data))

	marshaled, err := m0.MarshalBinary()

	if err != nil {
		t.Fatalf("failed to marshal message: %s", err)
	}

	var m1 Message
	if err := m1.UnmarshalBinary(marshaled); err != nil {
		t.Fatalf("failed to unmarshal message: %s", err)
	}

	if !m0.Time().Equal(m1.Time()) {
		t.Errorf("timestamps %s and %s do not match", m0.Time(), m1.Time())
	}

	if m0.From() != m1.From() {
		t.Errorf("owners '%s' and '%s' do not match", m0.From(), m1.Data())
	}

	if string(m0.Data()) != string(m1.Data()) {
		t.Errorf("data '%s' and '%s' do not match", m0.Data(), m1.Data())
	}
}
