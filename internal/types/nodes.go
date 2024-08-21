package types

type NodeState int

const (
	NodeStateAlloc NodeState = iota
	NodeStateComp
	NodeStateDown
	NodeStateDrain
	NodeStateFail
	NodeStateErr
	NodeStateIdle
	NodeStateMaint
	NodeStateMix
	NodeStateResv
)
