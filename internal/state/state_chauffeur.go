package state

import "fmt"

var StateMachine *stateMachine

func InitStateMachine() {
	StateMachine = &stateMachine{
		CurrentState: make(map[int64]State),
	}
}

func (stateMachine *stateMachine) Set(stateType StateType, chatID int64, data *StateContext) State {
	stateFactory, exists := statesFactory[stateType]
	if !exists {
		panic(fmt.Sprintf("no state factory was found for a state with type = %s", stateType))
	}

	stateMachine.CurrentState[chatID] = stateFactory(data)
	return stateMachine.CurrentState[chatID]
}

func (stateMachine *stateMachine) Get(chatID int64) State {
	handler, exists := stateMachine.CurrentState[chatID]
	if !exists {
		handler = stateMachine.Set(Idle, chatID, &StateContext{})
	}
	return handler
}
