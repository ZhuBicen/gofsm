package gofsm

import (
	"errors"
)

type StateMachine interface {
	Name() string
	SetName(string)
	SetInitialState(state State)
	ProcessEvent(event Event) (bool, error)
	StateTransition(state State)
	Terminate()
	//DeferEvent(Event)
	//SendInternalEvent(Event)
	InitFSM()
	CurrentState() State
}

type StateMachineBase struct {
	name string
	currentState State
	newState State
}

func (this *StateMachineBase) Name() string {
	return this.name
}

func (this *StateMachineBase) SetName(name string) {
	this.name = name
}

func (this *StateMachineBase) CurrentState() State {
	return this.currentState
}

func (this *StateMachineBase) SetInitialState(state State) {
	this.StateTransition(state)
	this.enterNewState()
}

func (this *StateMachineBase) Terminate() {
	if this.currentState != nil {
		this.currentState.ExitAction()
		this.currentState = nil
	}
}

func (this *StateMachineBase) ProcessEvent(event Event) (bool, error) {
	if this.currentState == nil {
		return false, errors.New("Strange error, no current state")
	}
	if !this.consumeEvent(event) {
		return false, nil
	}
	
	if this.newState != nil {
		this.enterNewState()
	}
	return true, nil

}

func (this *StateMachineBase) enterNewState() {
	for this.newState != nil {
		callExitActionsAndSetHistoryState(this.currentState, this.newState)
		callEntryActions(this.currentState, this.newState)
		this.currentState = this.newState
		this.newState = nil
		if cs, ok := this.currentState.(CompositeState);  ok {
			cs.InitTransition()
		}
	}
	
}

func (this *StateMachineBase) StateTransition(newState State) {
	this.newState = newState
}

func (this *StateMachineBase) consumeEvent(event Event) bool {
	tryingState := this.currentState
	for tryingState != nil {
		consumed := tryingState.HandleEvent(event)
		if consumed {
			return true
		}
		if !consumed {
			tryingState = tryingState.SuperState()
		}
	}
	return false
}

func (this *StateMachineBase) InitFSM() {

}	
	

func callExitActionsAndSetHistoryState(currentState State, newState State) {
	if currentState == nil {
		return
	}
	if currentState == newState {
		currentState.ExitAction()
		return
	}
	
	existingState := currentState

	for existingState != nil{

		targetState := newState
		for targetState != nil{
			//if current existing state is one of the super state of target state
			//there is no need to do the exit action
			if existingState == targetState.SuperState() {
				return
			}
			targetState = targetState.SuperState()
		}

		if existingState.SuperState() != nil {
			existingState.SuperState().SetShallowHistory(existingState)
			existingState.SuperState().SetDeepHistory(currentState)
		}
		
		existingState.ExitAction()
		existingState = existingState.SuperState()
	}
}

func callEntryActions(currentState State, newState State){
	entryState := newState
	entryStateSlice := make([]State, 0)
	for entryState != nil {

		entryStateSlice = append(entryStateSlice, entryState)
		
		entryState = entryState.SuperState()

		if entryState == nil {
			break
		}

		sourceState := currentState
		
		for sourceState != nil {
			//if curent entry state is one of the super state of the source state
			//there is no need to do the exit action
			if sourceState == entryState {
				return
			}
			sourceState = sourceState.SuperState()
		}
	}

	for i := len(entryStateSlice) -1 ; i >= 0; i-- {
		s := entryStateSlice[i]
		s.EntryAction()
	}
}