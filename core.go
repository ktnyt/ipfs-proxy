package proxy

import (
	"fmt"

	"github.com/ipfs/go-ipfs-api"
	peer "github.com/libp2p/go-libp2p-peer"
)

var ipfs *shell.Shell
var user peer.ID

type Proxy struct {
	sub    *shell.PubSubSubscription
	level  uint8
	prefix string
	Msgs   map[peer.ID]Message
}

func setUser() error {
	idout, err := ipfs.ID()

	if err != nil {
		return err
	}

	id, err := peer.IDB58Decode(idout.ID)

	if err != nil {
		return err
	}

	user = id

	return nil
}

func newProxy(prefix string) (*Proxy, error) {
	if err := setUser(); err != nil {
		return nil, fmt.Errorf("failed to create proxy: %s", err)
	}

	p := &Proxy{
		prefix: prefix,
		level:  1,
		Msgs:   make(map[peer.ID]Message),
	}

	topic := p.Topic()

	var err error
	if p.sub, err = ipfs.PubSubSubscribe(topic); err != nil {
		return nil, fmt.Errorf("failed to join proxy '%s': %s", topic, err)
	}

	return p, nil
}

func NewLocalProxy(prefix string) (*Proxy, error) {
	ipfs = shell.NewLocalShell()

	return newProxy(prefix)
}

func NewProxy(prefix, url string) (*Proxy, error) {
	ipfs = shell.NewShell(url)

	return newProxy(prefix)
}

func (p *Proxy) Topic() string {
	return fmt.Sprintf("%s/%s", p.prefix, user[2:2+p.level])
}

func (p *Proxy) Ping(payload []byte) (err error) {
	m := NewMessage(payload)

	var data []byte
	if data, err = m.MarshalBinary(); err != nil {
		return fmt.Errorf("failed to ping: %s", err)
	}

	if err = ipfs.PubSubPublish(p.Topic(), string(data)); err != nil {
		return fmt.Errorf("failed to ping: %s", err)
	}

	return nil
}

func (p *Proxy) Spin(c chan error) {
	for {
		rec, err := p.sub.Next()

		if err != nil {
			c <- fmt.Errorf("failed to get next message: %s", err)
		}

		data := rec.Data()

		if len(data) == 0 {
			continue
		}

		var m Message
		if err := m.UnmarshalBinary(data); err != nil {
			c <- fmt.Errorf("failed to process message: %s", err)
		}

		p.Msgs[m.From()] = m
	}
}

func (p *Proxy) Cancel() error {
	if p.sub == nil {
		return fmt.Errorf("proxy not connected")
	}

	return p.sub.Cancel()
}
