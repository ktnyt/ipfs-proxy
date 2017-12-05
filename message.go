package proxy

import (
	"fmt"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
)

const (
	timeOffset = 15
	fromOffset = 46
)

type Message struct {
	time time.Time
	from peer.ID
	data []byte
}

func NewMessage(data []byte) Message {
	return Message{
		time: time.Now(),
		from: user,
		data: data,
	}
}

func (m *Message) Time() time.Time {
	return m.time
}

func (m *Message) From() peer.ID {
	return m.from
}

func (m *Message) Data() []byte {
	return m.data
}

func (m Message) MarshalBinary() ([]byte, error) {
	var err error

	var timeBytes []byte
	if timeBytes, err = m.time.MarshalBinary(); err != nil {
		return nil, fmt.Errorf("Message.MarshalBinary: %s", err)
	}

	fromBytes := []byte(peer.IDB58Encode(m.from))
	headBytes := append(timeBytes, fromBytes...)

	return append(headBytes, m.data...), nil
}

func (m *Message) UnmarshalBinary(data []byte) error {
	if len(data) < timeOffset+fromOffset {
		return fmt.Errorf("Message.UnmarshalBinary: data not sufficient")
	}

	var t time.Time
	if err := t.UnmarshalBinary(data[:timeOffset]); err != nil {
		return fmt.Errorf("Message.UnmarshalBinary: %s", err)
	}

	s := string(data[timeOffset : timeOffset+fromOffset])

	from, err := peer.IDB58Decode(s)
	if err != nil {
		return fmt.Errorf("Message.UnmarshalBinary: %s", err)
	}

	m.time = t
	m.from = from
	m.data = data[timeOffset+fromOffset:]

	return nil
}
