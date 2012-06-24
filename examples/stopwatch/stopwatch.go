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
	fsm.CompositeStateBase
	elapsedTime time.Duration 
}

func NewActiveState(name string, m fsm.StateMachine) *ActiveState {
	this := &ActiveState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
	return this
}

func (this *ActiveState)HandleEvent(evt fsm.Event) bool {
	if evt.MessageId() == EVENT_RESET_ID {
		fmt.Println("Active state, Received EVENT_RESET_ID")
		this.elapsedTime = 0
		return true
	}
	return false
}


//running state
type RunningState struct {
	fsm.StateBase
	startTime time.Time
}

func NewRunningState(name string, m fsm.StateMachine) *RunningState {
	this := &RunningState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
	return this
}

func (this *RunningState) EntryAction() {
	fmt.Println("Entering RunningState")
	this.startTime = time.Now()
}

func (this *RunningState) ExitAction() {

	stopWatch := this.StateMachine().(*StopWatch)
	stopWatch.ActiveState.(*ActiveState).elapsedTime += time.Since(this.startTime)
	fmt.Println("Exiting RunningState", stopWatch.ActiveState.(*ActiveState).elapsedTime)
}

func (this *RunningState) HandleEvent(evt fsm.Event) bool{
	if evt.MessageId() == EVENT_STARTSTOP_ID {
		fmt.Println("RunningState, received EVENT_STARTSTOP_ID")
		stopWatch, _ := this.StateMachine().(*StopWatch)
		stopWatch.StateTransition(stopWatch.StoppedState)
		return true
	}
	return false
}
	

//stopped state
type StoppedState struct {
	fsm.StateBase
}

func NewStoppedState(name string, m fsm.StateMachine) *StoppedState {
	this := &StoppedState{}
	this.StateBase.SetStateMachine(m)
	this.StateBase.SetName(name)
	return this
}

func (this *StoppedState) HandleEvent(evt fsm.Event) bool {
	if evt.MessageId() == EVENT_STARTSTOP_ID {
		fmt.Println("StoppedState received EVENT_STARTSTOP_ID")
		stopWatch, _ := this.StateMachine().(*StopWatch)
		stopWatch.StateTransition(stopWatch.RunningState)
		return true
	}
	return false
}

func (this *StoppedState) EntryAction() {
	fmt.Println("Entering StoppedState")
}

func (this *StoppedState) ExitAction() {
	fmt.Println("Exiting StoppedState")
}

//state machine
type StopWatch struct {
	fsm.StateMachineBase
	ActiveState fsm.CompositeState
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
	this.RunningState.SetSuperState(this.ActiveState)
	this.StoppedState.SetSuperState(this.ActiveState)
	this.ActiveState.SetInitTransition(this.StoppedState)

	this.SetInitialState(this.ActiveState)
}

func(this *StopWatch) ElapsedTime() time.Duration{
	activeState, _ := this.ActiveState.(*ActiveState)
	return activeState.elapsedTime
}


func main() {
	myWatch := NewStopWatch()
	myWatch.InitMachine()
	myWatch.ProcessEvent(fsm.NewEventBase(EVENT_RESET_ID))
	myWatch.ProcessEvent(fsm.NewEventBase(EVENT_STARTSTOP_ID))
	time.Sleep(5*time.Second)
	myWatch.ProcessEvent(fsm.NewEventBase(EVENT_STARTSTOP_ID))
	fmt.Println("ElapsedTime:", myWatch.ElapsedTime())
	myWatch.Terminate()

}