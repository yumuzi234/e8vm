package table

// Action is a pending table action.
type Action struct {
	Action string
	Pos    int
	Text   string
}

// Render is a rendering engine that receives table actions.
type Render interface {
	Act(a *Action)
}
