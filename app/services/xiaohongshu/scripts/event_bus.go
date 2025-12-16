package scripts

import "github.com/asaskevich/EventBus"

var eventBus EventBus.Bus

func InitEventBus() {
	eventBus = EventBus.New()
}
func GetEventBus() EventBus.Bus {
	return eventBus
}
