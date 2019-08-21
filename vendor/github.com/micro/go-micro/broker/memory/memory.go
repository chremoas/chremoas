// Package memory provides a memory broker
package memory

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/broker"
	maddr "github.com/micro/go-micro/util/addr"
	mnet "github.com/micro/go-micro/util/net"
)

type memoryBroker struct {
	opts broker.Options

	addr string
	sync.RWMutex
	connected   bool
	Subscribers map[string][]*memorySubscriber
}

type memoryEvent struct {
	topic   string
	message *broker.Message
}

type memorySubscriber struct {
	id      string
	topic   string
	exit    chan bool
	handler broker.Handler
	opts    broker.SubscribeOptions
}

func (m *memoryBroker) Options() broker.Options {
	return m.opts
}

func (m *memoryBroker) Address() string {
	return m.addr
}

func (m *memoryBroker) Connect() error {
	m.Lock()
	defer m.Unlock()

	if m.connected {
		return nil
	}

	addr, err := maddr.Extract("::")
	if err != nil {
		return err
	}
	i := rand.Intn(20000)
	// set addr with port
	addr = mnet.HostPort(addr, 10000+i)

	m.addr = addr
	m.connected = true

	return nil
}

func (m *memoryBroker) Disconnect() error {
	m.Lock()
	defer m.Unlock()

	if !m.connected {
		return nil
	}

	m.connected = false

	return nil
}

func (m *memoryBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *memoryBroker) Publish(topic string, message *broker.Message, opts ...broker.PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return errors.New("not connected")
	}

	subs, ok := m.Subscribers[topic]
	m.RUnlock()
	if !ok {
		return nil
	}

	p := &memoryEvent{
		topic:   topic,
		message: message,
	}

	for _, sub := range subs {
		if err := sub.handler(p); err != nil {
			return err
		}
	}

	return nil
}

func (m *memoryBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, errors.New("not connected")
	}
	m.RUnlock()

	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	sub := &memorySubscriber{
		exit:    make(chan bool, 1),
		id:      uuid.New().String(),
		topic:   topic,
		handler: handler,
		opts:    options,
	}

	m.Lock()
	m.Subscribers[topic] = append(m.Subscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		var newSubscribers []*memorySubscriber
		for _, sb := range m.Subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.Subscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil
}

func (m *memoryBroker) String() string {
	return "memory"
}

func (m *memoryEvent) Topic() string {
	return m.topic
}

func (m *memoryEvent) Message() *broker.Message {
	return m.message
}

func (m *memoryEvent) Ack() error {
	return nil
}

func (m *memorySubscriber) Options() broker.SubscribeOptions {
	return m.opts
}

func (m *memorySubscriber) Topic() string {
	return m.topic
}

func (m *memorySubscriber) Unsubscribe() error {
	m.exit <- true
	return nil
}

func NewBroker(opts ...broker.Option) broker.Broker {
	var options broker.Options
	rand.Seed(time.Now().UnixNano())
	for _, o := range opts {
		o(&options)
	}

	return &memoryBroker{
		opts:        options,
		Subscribers: make(map[string][]*memorySubscriber),
	}
}
