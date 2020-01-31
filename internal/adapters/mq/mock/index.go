package mock

import (
	"sync"

	"github.com/rendau/barot/internal/domain/entities"
)

// St - is type for rabbit-mq client
type St struct {
	q  []*entities.BannerEvent
	mu sync.Mutex
}

// NewSt - creates new St instance
func NewSt() *St {
	return &St{
		q: make([]*entities.BannerEvent, 0),
	}
}

// PublishBannerEvent - publishes event to mq
func (q *St) PublishBannerEvent(event *entities.BannerEvent) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.q = append(q.q, event)

	return nil
}

// PullAll is pulls all
func (q *St) PullAll() []*entities.BannerEvent {
	q.mu.Lock()
	defer q.mu.Unlock()

	res := q.q

	q.q = make([]*entities.BannerEvent, 0)

	return res
}

// Clean is cleans
func (q *St) Clean() {
	_ = q.PullAll()
}
