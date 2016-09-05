package arch

const (
	pageVoid      = 0
	pageInterrupt = 1
	pageBasicIO   = 2
	pageRPC       = 3

	pageScreenText  = 5
	pageScreenColor = 6
	pageSysInfo     = 7
	pageBootImage   = 8

	pageMin = 16
)

// Basic IO page layout.
const (
	consoleBase = 0x0   // 0-8
	bootArgBase = 0x8   // 8-c
	clicksBase  = 0x10  // 10-14
	serialBase  = 0x80  // 80-100
	romBase     = 0x100 // 100-180
)
