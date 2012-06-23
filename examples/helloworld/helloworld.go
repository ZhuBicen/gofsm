package main

import (
	fsm "github.com/ZhuBicen/gofsm"
	"fmt"
)

type Greeting struct {
	fsm.StateBase
}

func NewGreeting(name string, m fsm.StateMachine) *Greeting {
	this :=  &Greeting{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
	return this
}

func (this *Greeting) EntryAction() {
	fmt.Println("Hello World")
}

func (this *Greeting) ExitAction() {
	fmt.Println("Bye Bye World")
}

//State machine
type Machine struct {
	fsm.StateMachineBase
	greetingState fsm.State
}

func NewMachine() *Machine {
	return &Machine{}
}

func (this *Machine )InitFSM() {
	this.greetingState = NewGreeting("HelloState", this)
	this.SetInitialState(this.greetingState)
}

func main() {
	m := NewMachine()
	m.InitFSM()
	m.Terminate()
}

	