package apply

import (
	"obfpl/data"
	"syscall"
)

type Context struct {
	profile *data.Profile
	swap    Swap
	ext     map[string]string
	group   map[string]string
	vars    map[string]string
	ctimes  map[string]syscall.Filetime
	name    string
	outpath string
	Loging  func(string)
}

type Swap struct {
	list []string
	idx  int
}
