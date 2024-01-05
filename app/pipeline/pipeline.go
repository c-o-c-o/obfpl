package pipeline

import (
	"obfpl/-packages/extgroup"
	"obfpl/-packages/sync"
	"obfpl/app/pipeline/profproc"
)

type Pipeline struct {
	ProfProc profproc.ProfProc
}

func NewPipeline(profilePath string, outPath string) (*Pipeline, error) {
	profProc, err := profproc.NewProfProc(profilePath, outPath)
	if err != nil {
		return nil, err
	}

	return &Pipeline{
		ProfProc: profProc,
	}, nil
}

func (p *Pipeline) CreateNotify(msgch chan<- string) func(string) {
	waiter := sync.NewClosedWaiter(msgch)
	extPool := extgroup.NewExtPool()

	return func(filePath string) {
		ext, err := p.ProfProc.SelectExt(filePath)
		if err != nil {
			msgch <- err.Error()
			return
		}

		files := extPool.Push(filePath, ext, p.ProfProc.CreateExtList())
		if files == nil {
			return
		}

		p.ProfProc.Call(files.Name, files.ExtGroup, waiter.Next())
	}
}
