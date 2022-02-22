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
	vars    map[string]string
	profile *data.Profile
	ext     map[string]string
	lastdst string
	Loging  func(string)
}
