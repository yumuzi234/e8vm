package build8

// OnlyFile returns the only file in the map. The map of files must have
// exactly one file, otherwise returns nil.
func OnlyFile(src map[string]*File) *File {
	if len(src) != 1 {
		return nil
	}
	for _, f := range src {
		return f
	}
	panic("unreachable")
}
