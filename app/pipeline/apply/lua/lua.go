package lua

import (
	lua "github.com/yuin/gopher-lua"
)

type Function struct {
	Name string
	Func func(flow *Flow) error
}

func NewFunction(name string, f func(flow *Flow) error) *Function {
	return &Function{
		Name: name,
		Func: f,
	}
}

type Flow struct {
	Args     []string
	state    *lua.LState
	rtnCount int
}

func NewFlow(s *lua.LState) *Flow {
	return &Flow{
		state:    s,
		rtnCount: 0,
	}
}

/*
offset 0
*/
func (f *Flow) GetString(argidx int) string {
	return f.state.ToString(argidx + 1)
}

func (f *Flow) ReturnBool(val bool) {
	f.state.Push(lua.LBool(val))
	f.rtnCount += 1
}

func (f *Flow) ReturnInt(val int) {
	f.state.Push(lua.LNumber(val))
	f.rtnCount += 1
}

func (f *Flow) ReturnString(val string) {
	f.state.Push(lua.LString(val))
	f.rtnCount += 1
}

func DoScript(script string, funcs []*Function) error {
	l := lua.NewState()
	defer l.Close()

	for i := 0; i < len(funcs); i++ {
		f := funcs[i]
		l.SetGlobal(f.Name, l.NewFunction(func(s *lua.LState) int {
			flow := NewFlow(s)
			if err := f.Func(flow); err != nil {
				l.Error(lua.LString(err.Error()), 0)
			}

			return flow.rtnCount
		}))

	}

	return l.DoString(script)
}
