package observe

import "github.com/fsnotify/fsnotify"

func Folder(obsp string, detect func(path string, erch chan error), erch chan error) {
	wtr, err := fsnotify.NewWatcher()
	if err != nil {
		erch <- err
		return
	}
	err = wtr.Add(obsp)
	if err != nil {
		erch <- err
		return
	}
	defer wtr.Close()

	for {
		select {
		case event, ok := <-wtr.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				detect(event.Name, erch)
			}
		case err, ok := <-wtr.Errors:
			if !ok {
				return
			}

			erch <- err
		}
	}
}
