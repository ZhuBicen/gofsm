package gofsm

type State interface {
	Name() string
	
	SetSuperState(*State)
	SuperState() *State
	
	SetShallowHistory(*State)
	ShallowHistory() *State

	SetDeepHistory(*State)
	DeepHistory() *State
	
	HandleEvent(*Event) bool

	EntryAction()
	ExitAction()
}


type CompositeState interface{
	State
	InitTransition()
}

type StateBase struct {
	name string
	superState *State
	deepHistory *State
	shallowHistory *State
	initState *State
}

func (this *StateBase) GetName() {
	return name
}

func (this *StateBase) SetSuperState(superSate *State) {
	this.superState = superState
}

func (this *StateBase) SuperState() *State {
	return this.superState
}

func (this *StateBase) SetDeepHistory(deepHistory *State) {
	this.deepHistory = deepHistory
}

func (this *StateBase) DeepHistory() *State {
	return this.deepHistory
}

func (this *StateBase) SetShallowHistory(shallowHistory *State) {
	this.shallowHistory = shallowHistory
}

func (this *StateBase) ShallowHistory() *State {
	return this.shallowHistory
}

func (this *StateBase) HandleEvent(*Event) bool{
	return false
}

func (this *StateBase) EntryAction() {
}

func (this *StateBase) ExitAction() {

}

type CompositeStateBase struct {
	StateBase
	initState *StateBase
}

func (this *CompositeStateBase) SetInitTransition(initState *State) {
	this.initState = initState
}

func (this *CompositeStateBase) InitTransition() {
	if this.initState {
		this.fsm.StateTransition(this.InitState)
	}
}



