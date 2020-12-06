package stdutil

import (
	"strings"
)

// EventChannel is a struct to describe the channel where the event is contained
type EventChannel struct {
	Application string
	Service     string
	Module      string
}

// EventData contains the data of the event.
type EventData struct {
	ID    string      `json:"id,omitempty"`
	Index int64       `json:"index,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// Event contains the channel and the event data
type Event struct {
	Channel string    `json:"channel,omitempty"`
	Data    EventData `json:"data,omitempty"`
}

// NewEventChannel properly creates a new event channel
func NewEventChannel(application string, service string, module string) EventChannel {

	application = strings.ReplaceAll(strings.ToLower(application), `.`, ``)
	service = strings.ReplaceAll(strings.ToLower(service), `.`, ``)
	module = strings.ReplaceAll(strings.ToLower(module), `.`, ``)

	return EventChannel{
		Application: application,
		Service:     service,
		Module:      module,
	}
}

// GetMatch event channel
func GetMatch(subject string, evtchans []EventChannel) *EventChannel {

	for _, e := range evtchans {
		if subject == e.ToString() {
			return &e
		}
	}

	return nil
}

// ToString composes the event channel to a proper channel name
func (ec *EventChannel) ToString() string {
	return strings.ToLower(ec.Application) + `.` + strings.ToLower(ec.Service) + `.` + strings.ToLower(ec.Module)
}
