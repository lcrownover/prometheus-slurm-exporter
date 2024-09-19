package types

type NodeState string

const (
	NodeStateAlloc         NodeState = "alloc"
	NodeStateComp          NodeState = "comp"
	NodeStateDown          NodeState = "down"
	NodeStateDrain         NodeState = "drain"
	NodeStateFail          NodeState = "fail"
	NodeStateErr           NodeState = "err"
	NodeStateIdle          NodeState = "idle"
	NodeStateMaint         NodeState = "maint"
	NodeStateMix           NodeState = "mix"
	NodeStateResv          NodeState = "resv"
	NodeStateNotResponding NodeState = "not_responding"
	NodeStateInvalidReg    NodeState = "invalid_reg"
)