package screen

// ScreenRender is an interface that renders the screen.
type ScreenRender interface {
	NeedUpdate() bool
	UpdateText(m map[uint32]byte)
	UpdateColor(m map[uint32]byte)
}
