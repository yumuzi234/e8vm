package builds

import (
	"path"
)

func listDirs(in Input, p string, lst []string) ([]string, error) {
	dirs, err := in.ListDirs(p)
	if err != nil {
		return nil, err
	}

	for i, dir := range dirs {
		dirs[i] = path.Join(p, dir)
	}

	for _, dir := range dirs {
		lst = append(lst, dir)
		lst, err = listDirs(in, dir, lst)
		if err != nil {
			return nil, err
		}
	}

	return lst, nil
}

// ListDirs lists all directories and sub directories under a path, including
// this path.
func ListDirs(in Input, p string) ([]string, error) {
	return listDirs(in, p, nil)
}
