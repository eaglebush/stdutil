package stdutil

import (
	"encoding/json"
	"strings"
)

var (
	ra func(s, old, new string) string = strings.ReplaceAll
	lw func(string) string             = strings.ToLower
)

type (
	// EventSubject is a struct to describe the subject event
	EventSubject struct {
		Application string // Application. This would form as the first segment
		Service     string // Service. This would form as the second segment
		Module      string // Module. This would form as the last segment
	}

	// Event contains the address and the event data
	Event struct {
		Index   int64       `json:"index,omitempty"`   // Optional. Index of this data to return to later
		Subject string      `json:"subject,omitempty"` // Subject of the event
		Data    interface{} `json:"data,omitempty"`    // Data of the event
	}
)

// NewEventSubjectBase properly creates a new event base.
func NewEventSubjectBase(application, service, module string) EventSubject {
	return EventSubject{
		Application: ra(lw(application), `.`, `-`),
		Service:     ra(lw(service), `.`, `-`),
		Module:      ra(lw(module), `.`, `-`),
	}
}

// GetEventSubjectMatch seeks the list of event by subject
func GetEventSubjectMatch(subject string, evtChans []EventSubject) *EventSubject {
	for _, e := range evtChans {
		if strings.EqualFold(subject, e.ToString(nil)) {
			return &e
		}
	}
	return nil
}

// GetEventModuleMatch seeks the list of events by module
func GetEventModuleMatch(module string, evtChans []EventSubject) *EventSubject {
	for _, e := range evtChans {
		if strings.EqualFold(module, e.Module) {
			return &e
		}
	}
	return nil
}

// ToString converts an EventSubject to a readable string
func (ec EventSubject) ToString(eventVerb *string) string {
	if eventVerb == nil || *eventVerb == "" {
		return lw(ec.Application) + `.` + lw(ec.Service) + `.` + lw(ec.Module)
	}
	return lw(ec.Application) + `.` + lw(ec.Service) + `.` + lw(ec.Module) + `.` + *eventVerb
}

// BuildEvent builds an event based on the inputs
func BuildEvent(subject EventSubject, eventVerb string, data interface{}, index int64) ([]byte, error) {
	return json.Marshal(
		Event{
			Index:   index,
			Subject: subject.ToString(&eventVerb),
			Data:    data,
		})
}
