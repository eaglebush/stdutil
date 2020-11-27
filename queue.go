package stdutil

import (
	"strings"
)

// EventChannel is a struct to describe the channel where the event is contained
type EventChannel struct {
	Application string
	Schema      string
	Module      string
	Events      []string
}

// NewEventChannel properly creates a new event channel
func NewEventChannel(application string, schema string, module string) EventChannel {

	application = strings.ReplaceAll(strings.ToLower(application), `.`, ``)
	schema = strings.ReplaceAll(strings.ToLower(schema), `.`, ``)
	module = strings.ReplaceAll(strings.ToLower(module), `.`, ``)

	return EventChannel{
		Application: application,
		Schema:      schema,
		Module:      module,
	}
}

// HandleEvent add event to channel and its handler
// func (ec *EventChannel) HandleEvent(eventid string) {

// 	var exists bool

// 	exists = false

// 	evt := Event{
// 		ID:      strings.TrimSpace(strings.ToLower(eventid)),
// 		Handler: handler,
// 	}

// 	for i, e := range ec.Events {
// 		if e.ID == evt.ID {
// 			ec.Events[i].Handler = handler
// 			exists = true
// 			break
// 		}
// 	}

// 	if !exists {
// 		ec.Events = append(ec.Events, evt)
// 	}

// }

// ToString composes the event channel to a proper channel name
func (ec *EventChannel) ToString() string {
	return strings.ToLower(ec.Application) + `.` + strings.ToLower(ec.Schema) + `.` + strings.ToLower(ec.Module)
}
