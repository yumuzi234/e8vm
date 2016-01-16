package build8

// Options contains building options
type Options struct {
	TestMaxCycle int64
	NoTest       bool
	JustStatic   bool
	JustImport   bool
}
