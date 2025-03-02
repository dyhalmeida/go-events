package events

import (
	"errors"
	"slices"
)

var ErrHandlerAlreadyExists = errors.New("handler already exists")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (dispatcher *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {

	if _, ok := dispatcher.handlers[eventName]; ok {
		if slices.Contains(dispatcher.handlers[eventName], handler) {
			return ErrHandlerAlreadyExists
		}
	}

	dispatcher.handlers[eventName] = append(dispatcher.handlers[eventName], handler)
	return nil
}

func (dispatcher *EventDispatcher) Clear() {
	dispatcher.handlers = make(map[string][]EventHandlerInterface)
}

func (dispatcher *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	_, ok := dispatcher.handlers[eventName]

	if ok && slices.Contains(dispatcher.handlers[eventName], handler) {
		return true
	}

	return false
}
