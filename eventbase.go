package gofsm

type Event interface {
	MessageId() int
}

type EventBase struct{
	id int
}

func (this *EventBase) MessageId() int {
	return this.id
}

func NewEventBase(messageId int) *EventBase {
	return &EventBase{ id: messageId }
}