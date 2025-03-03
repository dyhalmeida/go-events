package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type EventHandlerMock struct {
	ID string
}

func (eventHandlerMock *EventHandlerMock) Handle(event EventInterface, wg *sync.WaitGroup) {}

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
	suite.handler1 = EventHandlerMock{ID: "0001"}
	suite.handler2 = EventHandlerMock{ID: "0002"}
	suite.handler3 = EventHandlerMock{ID: "0003"}
	suite.event1 = EventMock{Name: "event mock 01", Payload: struct{ ID, name string }{ID: "01", name: "event mock 01"}}
	suite.event2 = EventMock{Name: "event mock 02", Payload: struct{ ID, name string }{ID: "02", name: "event mock 02"}}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherSuiteTest))
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Register() {

	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	assert.Equal(suite.T(), &suite.handler1, suite.eventDispatcher.handlers[suite.event1.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][1])
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Register_SameHandler() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.ErrorIs(ErrHandlerAlreadyExists, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler1))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler3))
}

type EventHandlerMocked struct {
	mock.Mock
}

func (EventHandlerMocked *EventHandlerMocked) Handle(event EventInterface, wg *sync.WaitGroup) {
	EventHandlerMocked.Called(event)
	wg.Done()
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Dispatch() {
	eventHandlerMocked1 := &EventHandlerMocked{}
	eventHandlerMocked2 := &EventHandlerMocked{}

	eventHandlerMocked1.On("Handle", &suite.event1)
	eventHandlerMocked2.On("Handle", &suite.event1)

	suite.eventDispatcher.Register(suite.event1.GetName(), eventHandlerMocked1)
	suite.eventDispatcher.Register(suite.event1.GetName(), eventHandlerMocked2)

	suite.eventDispatcher.Dispatch(&suite.event1)

	eventHandlerMocked1.AssertExpectations(suite.T())
	eventHandlerMocked2.AssertExpectations(suite.T())

	eventHandlerMocked1.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eventHandlerMocked2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Remove() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler1)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][0])

	suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler2)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
}
