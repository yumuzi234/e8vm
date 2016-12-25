package devs

// TableAction is a pending table action.
type TableAction struct {
	Bytes []byte
}

// TableRender is a rendering engine that receives table actions.
type TableRender interface {
	Act(a *TableAction)
}
