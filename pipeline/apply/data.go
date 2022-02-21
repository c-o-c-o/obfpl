package apply

import (
	"obfpl/data"
	"syscall"
)

type ApplyContext struct {
	temp    string
	swaps   []string
	name    string
	group   map[string]string
	ctimes  map[string]syscall.Filetime
	repList map[string]string
	profile *data.Profile
	ext     map[string]string
	Loging  func(string)
}
