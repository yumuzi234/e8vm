package builds

import (
	"fmt"
	"strings"
)

func isValidPathName(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '_' {
			continue
		}
		if r == '.' {
			if i == 0 {
				return false
			}
			continue
		}
		return false
	}
	return true
}

func isValidPath(p string) bool {
	parts := strings.Split(p, "/")
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if !isValidPathName(part) {
			return false
		}
	}
	return true
}

// CheckValidPath checks if the given path is a valid path.
func CheckValidPath(p string) error {
	if !isValidPath(p) {
		return fmt.Errorf("%q is not a valid path", p)
	}
	return nil
}

func checkValidDir(p string) error {
	if p == "" {
		return nil
	}
	return CheckValidPath(p)
}
