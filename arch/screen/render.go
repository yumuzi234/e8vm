package screen

// Render is an interface that renders the screen.
type Render interface {
	NeedUpdate() bool
	UpdateText(m map[uint32]byte)
	UpdateColor(m map[uint32]byte)
}
