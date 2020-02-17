package mock

import (
	"sync"

	"github.com/rendau/barot/internal/domain/entities"
)

// MessageQueueMock - is type for rabbit-mq client
type MessageQueueMock struct {
	q  []*entities.BannerEvent
	mu sync.Mutex
}

// NewMessageQueueMock - creates new MessageQueueMock instance
func NewMessageQueueMock() *MessageQueueMock {
	return &MessageQueueMock{
		q: make([]*entities.BannerEvent, 0),
	}
}

// PublishBannerEvent - publishes event to mq
func (q *MessageQueueMock) PublishBannerEvent(event *entities.BannerEvent) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.q = append(q.q, event)

	return nil
}

// PullAll is pulls all
func (q *MessageQueueMock) PullAll() []*entities.BannerEvent {
	q.mu.Lock()
	defer q.mu.Unlock()

	res := q.q

	q.q = make([]*entities.BannerEvent, 0)

	return res
}

// Clean is cleans
func (q *MessageQueueMock) Clean() {
	_ = q.PullAll()
}
