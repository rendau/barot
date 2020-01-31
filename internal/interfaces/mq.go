package interfaces

import (
	"github.com/rendau/barot/internal/domain/entities"
)

// Mq is interface for mq
type Mq interface {
	PublishBannerEvent(event *entities.BannerEvent) error
}
