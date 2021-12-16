package apply

import (
	"obfpl/data"
	"syscall"
)

type ApplyContext struct {
	name    string
	swaps   []string
	group   map[string]string
	ctimes  map[string]syscall.Filetime
	repList map[string]string
	profile *data.Profile
	Loging  func(string)
}
