package image

// CodeStart returns the virtual address of the first code section.
func CodeStart(secs []*Section) (uint32, bool) {
	for _, s := range secs {
		if s.Type == Code {
			return s.Addr, true
		}
	}
	return 0, false
}
