package sym8

// IsPublic checks if a symbol name is public.
func IsPublic(name string) bool {
	if name == "" {
		return false
	}
	r := name[0]
	return r >= 'A' && r <= 'Z'
}
