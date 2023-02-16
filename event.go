package stdutil

import (
	"strings"
)

// EventChannel is a struct to describe the channel where the event is contained
type EventChannel struct {
	Application string // Application. This would form as the first segment
	Service     string // Service. This would form as the second segment
	Module      string // Module. This would form as the last segment
	Stream      string // Stream. The stream name when using a JetStream subscription
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
func NewEventChannel(application, service, module string) EventChannel {

	application = strings.ReplaceAll(strings.ToLower(application), `.`, ``)
	service = strings.ReplaceAll(strings.ToLower(service), `.`, ``)
	module = strings.ReplaceAll(strings.ToLower(module), `.`, ``)

	return EventChannel{
		Application: application,
		Service:     service,
		Module:      module,
	}
}

// NewStreamEventChannel properly creates a new event channel with stream name
func NewStreamEventChannel(application, service, module, stream string) EventChannel {

	application = strings.ReplaceAll(strings.ToLower(application), `.`, ``)
	service = strings.ReplaceAll(strings.ToLower(service), `.`, ``)
	module = strings.ReplaceAll(strings.ToLower(module), `.`, ``)

	return EventChannel{
		Application: application,
		Service:     service,
		Module:      module,
		Stream:      stream,
	}
}

// GetEventSubjectMatch seeks the list of event channels by module
func GetEventSubjectMatch(subject string, evtchans []EventChannel) *EventChannel {

	for _, e := range evtchans {
		if strings.EqualFold(subject, e.ToString()) {
			return &e
		}
	}

	return nil
}

// GetEventModuleMatch seeks the list of event channels by module
func GetEventModuleMatch(module string, evtchans []EventChannel) *EventChannel {

	for _, e := range evtchans {
		if strings.EqualFold(module, e.Module) {
			return &e
		}
	}

	return nil
}

// ToString composes the event channel to a proper channel name
func (ec *EventChannel) ToString() string {
	return strings.ToLower(ec.Application) + `.` + strings.ToLower(ec.Service) + `.` + strings.ToLower(ec.Module)
}
