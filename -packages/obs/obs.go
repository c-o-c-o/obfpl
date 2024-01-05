package obs

import "github.com/fsnotify/fsnotify"

type DirObserver struct {
	Msgch chan string
	path  string
	wtr   *fsnotify.Watcher
}

/*
path 監視するフォルダーのパス
*/
func NewDirObserver(path string) (*DirObserver, error) {
	dobs := &DirObserver{
		Msgch: make(chan string),
		path:  path,
	}

	wtr, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dobs.wtr = wtr
	return dobs, nil
}

func (o *DirObserver) Destroy() {
	o.wtr.Close()
	close(o.Msgch)
}

/*
処理の都合上レスポンスを上げるために、生成時のみ通知
*/
func (o *DirObserver) Observe(created func(filePath string)) {
	err := o.wtr.Add(o.path)
	if err != nil {
		o.Msgch <- err.Error()
		return
	}

	for {
		select {
		case event, ok := <-o.wtr.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				created(event.Name)
			}
		case err, ok := <-o.wtr.Errors:
			if !ok {
				return
			}

			o.Msgch <- err.Error()
		}
	}
}
