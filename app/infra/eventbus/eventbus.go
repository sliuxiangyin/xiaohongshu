package eventbus

import (
	"github.com/asaskevich/EventBus"
)

var GlobalEventBus EventBus.Bus

func init() {
	GlobalEventBus = EventBus.New()
}
