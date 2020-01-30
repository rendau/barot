package interfaces

import (
	"github.com/rendau/barot/internal/domain/entities"
)

type Mq interface {
	PublishBannerEvent(event *entities.BannerEvent) error
}
