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