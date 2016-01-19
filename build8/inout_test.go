package build8

// satisfying the interface
var _ Input = new(MemHome)
var _ Output = new(MemHome)

var _ Input = new(DirHome)
var _ Output = new(DirHome)
