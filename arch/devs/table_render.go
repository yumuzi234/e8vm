package devs

// TableAction is a pending table action.
type TableAction struct {
	Action string
	Pos    int
	Text   string
}

// TableRender is a rendering engine that receives table actions.
type TableRender interface {
	Act(a *TableAction)
}
