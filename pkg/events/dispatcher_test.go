package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EventMock struct {
	Name    string
	Payload any
}

func (eventMock *EventMock) GetName() string {
	return eventMock.Name
}

func (eventMock *EventMock) GetPayload() any {
	return eventMock.Payload
}

func (eventMock *EventMock) GetDateTime() time.Time {
	return time.Now()
}

type EventHandlerMock struct{}

func (eventHadlerMock *EventHandlerMock) Handler(event EventInterface) {}

type EventDispatcherSuiteTest struct {
	suite.Suite
	event1          EventMock
	event2          EventMock
	handler1        EventHandlerMock
	handler2        EventHandlerMock
	handler3        EventHandlerMock
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherSuiteTest) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler1 = EventHandlerMock{}
	suite.handler2 = EventHandlerMock{}
	suite.handler3 = EventHandlerMock{}
	suite.event1 = EventMock{Name: "event mock 01", Payload: struct{ ID, name string }{ID: "01", name: "event mock 01"}}
	suite.event2 = EventMock{Name: "event mock 02", Payload: struct{ ID, name string }{ID: "02", name: "event mock 02"}}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherSuiteTest))
}
