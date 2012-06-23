package main

import (
	fsm "github.com/ZhuBicen/gofsm"
	"fmt"
	"time"
)


const (
	EVENT_RESET_ID = iota
	EVENT_STARTSTOP_ID = iota
)

//active state 	
type ActiveState struct {
	fsm.StateBase
	elapsedTime time.Duration 
}

func NewActiveState(name string, m fsm.StateMachine) *ActiveState {
	this := &ActiveState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
}

func (this *ActiveState)HandleEvent(evt fsm.Event) bool {
	if evt.MessageId() == EVENT_RESET_ID {
		fmt.Println("Active state, Received EVENT_RESET_ID")
	}
}


//running state
type RunningState struct {
	fsm.StateBase
	startTime time.Time
}

func NewRunningState(name string, m fsm.StateMachine) *ActiveState {
	this := &RunningState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
}

func (this *RunningState) EntryAction() {
	this.startTime = time.Now()
}

func (this *RunningState) ExitAction() {
	stopWatch := this.StateMachine().(StopWatch)
	stopWatch.ActiveState.elapsedTime += time.Since(this.startTime)
}

func (this *RunningState) HandleEvent(evt fsm.Event) {
	if evt.MessageId() == EVENT_STARTSTOP_ID {
		stopWatch, _ := this.StateMachine().(StopWatch)
		stopWatch.StateTransition(stopWatch.StoppedState)
	}
}
	

//stopped state
type StoppedState struct {
	fsm.StateBase
}

func NewStoppedState(name string, m fsm.StateMachine) *ActiveState {
	this := &StoppedState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
}

func (this *StoppedState) HandleEvent(evt fsm.Event) {
	if evt.MessageId() == EVENT_STARTSTOP_ID {
		stopWatch, _ := this.StateMachine().(StopWatch)
		stopWatch.StateTransition(stopWatch.RunningState)
	}
}

//state machine
type StopWatch struct {
	fsm.StateMachineBase
	ActiveState fsm.State
	RunningState fsm.State
	StoppedState fsm.State
}

func NewStopWatch() *StopWatch {
	return &StopWatch{}
}

func (this *StopWatch) InitMachine() {
	this.ActiveState  = NewActiveState("ActiveState", this)
	this.RunningState = NewRunningState("RunningState", this)
	this.StoppedState = NewStoppedState("StoppedState", this)
}

func(this *StopWatch) ElapsedTime() time.Duration{
	activeState, _ := this.ActiveState.(ActiveState)
	return activeState.elapsedTime
}


func main() {
	sw := NewStopWatch()
	sw.InitMachine()
	sw.Terminate()

}