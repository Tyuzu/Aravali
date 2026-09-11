package mq

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrClosed = errors.New("mq is closed")
)

type memorySubscription struct {
	id      uint64
	subject string
	queue   string
	mq      *MemoryMQ
}

func (s *memorySubscription) Unsubscribe() error {
	return s.mq.unsubscribe(s.subject, s.queue, s.id)
}

type subscriber struct {
	id      uint64
	handler MessageHandler
}

type memoryMQ struct {
	mu          sync.RWMutex
	closed      bool
	nextSubID   uint64
	pubSub      map[string][]*subscriber
	queueGroups map[string]map[string][]*subscriber // subject -> queueName -> subscribers
	rrCounters  map[string]map[string]*uint64       // subject -> queueName -> counter
}

// MemoryMQ implements the MQ interface in-memory.
type MemoryMQ struct {
	*memoryMQ
}

// NewMemoryMQ creates a new in-memory message queue.
func NewMemoryMQ() *MemoryMQ {
	return &MemoryMQ{
		memoryMQ: &memoryMQ{
			pubSub:      make(map[string][]*subscriber),
			queueGroups: make(map[string]map[string][]*subscriber),
			rrCounters:  make(map[string]map[string]*uint64),
		},
	}
}

func (m *MemoryMQ) Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return ErrClosed
	}
	return nil
}

func (m *MemoryMQ) Publish(ctx context.Context, subject string, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return ErrClosed
	}

	msg := Message{
		Subject: subject,
		Data:    append([]byte(nil), data...), // Copy data to ensure immutability
	}

	// Gather broadcast subscribers
	var pubSubHandlers []MessageHandler
	for _, sub := range m.pubSub[subject] {
		pubSubHandlers = append(pubSubHandlers, sub.handler)
	}

	// Gather one handler per queue group (round-robin selection)
	var queueHandlers []MessageHandler
	if groups, ok := m.queueGroups[subject]; ok {
		for qName, subs := range groups {
			if len(subs) == 0 {
				continue
			}
			counter := m.rrCounters[subject][qName]
			idx := atomic.AddUint64(counter, 1) % uint64(len(subs))
			queueHandlers = append(queueHandlers, subs[idx].handler)
		}
	}
	m.mu.RUnlock()

	// Dispatch concurrently to avoid blocking the publisher
	allHandlers := append(pubSubHandlers, queueHandlers...)
	for _, h := range allHandlers {
		handler := h
		go func() {
			_ = handler(ctx, msg) // ACK/NACK handled by return code in application code
		}()
	}

	return nil
}

func (m *MemoryMQ) Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, ErrClosed
	}

	m.nextSubID++
	subID := m.nextSubID

	sub := &subscriber{
		id:      subID,
		handler: handler,
	}

	m.pubSub[subject] = append(m.pubSub[subject], sub)

	return &memorySubscription{
		id:      subID,
		subject: subject,
		mq:      m,
	}, nil
}

func (m *MemoryMQ) QueueSubscribe(ctx context.Context, subject string, queue string, handler MessageHandler) (Subscription, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, ErrClosed
	}

	m.nextSubID++
	subID := m.nextSubID

	sub := &subscriber{
		id:      subID,
		handler: handler,
	}

	if _, ok := m.queueGroups[subject]; !ok {
		m.queueGroups[subject] = make(map[string][]*subscriber)
		m.rrCounters[subject] = make(map[string]*uint64)
	}

	if _, ok := m.queueGroups[subject][queue]; !ok {
		var counter uint64
		m.rrCounters[subject][queue] = &counter
	}

	m.queueGroups[subject][queue] = append(m.queueGroups[subject][queue], sub)

	return &memorySubscription{
		id:      subID,
		subject: subject,
		queue:   queue,
		mq:      m,
	}, nil
}

func (m *MemoryMQ) unsubscribe(subject string, queue string, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if queue == "" {
		// Broadcast subscription
		subs := m.pubSub[subject]
		for i, sub := range subs {
			if sub.id == id {
				m.pubSub[subject] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(m.pubSub[subject]) == 0 {
			delete(m.pubSub, subject)
		}
	} else {
		// Queue subscription
		if groups, ok := m.queueGroups[subject]; ok {
			if subs, ok := groups[queue]; ok {
				for i, sub := range subs {
					if sub.id == id {
						m.queueGroups[subject][queue] = append(subs[:i], subs[i+1:]...)
						break
					}
				}
				if len(m.queueGroups[subject][queue]) == 0 {
					delete(m.queueGroups[subject], queue)
					delete(m.rrCounters[subject], queue)
				}
			}
			if len(m.queueGroups[subject]) == 0 {
				delete(m.queueGroups, subject)
				delete(m.rrCounters, subject)
			}
		}
	}

	return nil
}

// Close releases resources and prevents new operations.
func (m *MemoryMQ) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}
