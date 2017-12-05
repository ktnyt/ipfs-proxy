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
		return nil, fmt.Errorf("newProxy: %s", err)
	}

	p := &Proxy{
		prefix: prefix,
		level:  1,
		Msgs:   make(map[peer.ID]Message),
	}

	topic := p.Topic()

	var err error
	if p.sub, err = ipfs.PubSubSubscribe(topic); err != nil {
		return nil, fmt.Errorf("newProxy: %s", topic, err)
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

func NewProxyWithIpfs(prefix string, s *shell.Shell) (*Proxy, error) {
	ipfs = s

	return newProxy(prefix)
}

func (p *Proxy) Topic() string {
	return fmt.Sprintf("%s/%s", p.prefix, user[2:2+p.level])
}

func (p *Proxy) Ping(payload []byte) (err error) {
	m := NewMessage(payload)

	var data []byte
	if data, err = m.MarshalBinary(); err != nil {
		return fmt.Errorf("Proxy.Ping: %s", err)
	}

	if err = ipfs.PubSubPublish(p.Topic(), string(data)); err != nil {
		return fmt.Errorf("Proxy.Ping: %s", err)
	}

	return nil
}

func (p *Proxy) Next() (err error) {
	var rec shell.PubSubRecord
	if rec, err = p.sub.Next(); err != nil {
		return fmt.Errorf("Proxy.Next: %s", err)
	}

	data := rec.Data()

	if len(data) == 0 {
		return nil
	}

	var m Message
	if err = m.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("Proxy.Next: %s", err)
	}

	p.Msgs[rec.From()] = m

	return nil
}

func (p *Proxy) Spin(c chan error) {
	for {
		select {
		case <-c:
			break
		default:
			if err := p.Next(); err != nil {
				c <- err
			}
		}
	}
}

func (p *Proxy) Cancel() error {
	if p.sub == nil {
		return fmt.Errorf("Proxy.Cancel: not connected")
	}

	return p.sub.Cancel()
}
