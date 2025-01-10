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
	NodeStatePlanned       NodeState = "planned"
	NodeStateNotResponding NodeState = "not_responding"
	NodeStateInvalid       NodeState = "invalid"
	NodeStateInvalidReg    NodeState = "invalid_reg"
	NodeStateDynamicNorm   NodeState = "dynamic_norm"
)
