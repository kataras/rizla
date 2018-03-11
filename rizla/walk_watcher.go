package rizla

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type walkWatcher struct {
	// to stop the for loop
	stopChan           chan bool
	hasStoppedManually bool
	errListeners       []WatcherErrorListener
	changeListeners    []WatcherChangeListener
}

var _ Watcher = &walkWatcher{}

// newWalkWatcher returns a new golang's stdlib filepath.Walker's wrapper
// which watching with every x milleseconds the projects' directories.
func newWalkWatcher() Watcher {
	return &walkWatcher{
		stopChan: make(chan bool, 1),
	}
}

func (w *walkWatcher) OnError(evt WatcherErrorListener) {
	w.errListeners = append(w.errListeners, evt)
}

func (w *walkWatcher) OnChange(evt WatcherChangeListener) {
	w.changeListeners = append(w.changeListeners, evt)
}

func (w *walkWatcher) Stop() {
	w.stopChan <- true
}

var errDoneNot = errors.New("done")

// DefaultWalkLoopSleep it's the sleep time of the loop when --walk flag is used(`filepath#Walk`).
// Defaults to 1.3 second.
var DefaultWalkLoopSleep = 1350 * time.Second

func (w *walkWatcher) loop(p *Project, stopChan chan bool) {
	for {
		select {
		case stop := <-stopChan:
			if stop {
				return
			}
		default:
			filepath.Walk(p.dir, func(path string, info os.FileInfo, err error) error {

				if filepath.Ext(path) == goExt && info.ModTime().After(p.lastChange) {
					for i := range w.changeListeners {
						w.changeListeners[i](p, path)
					}
					return errDoneNot
				}

				return nil
			})
			// loop every 1.3 second.
			time.Sleep(DefaultWalkLoopSleep)
		}
	}
}

func (w *walkWatcher) Loop() {
	w.stopChan <- false

	for _, p := range projects {
		go w.loop(p, w.stopChan)
	}

	defer func() {
		for _, p := range projects {
			killProcess(p.proc, p.AppName)
		}
		if !w.hasStoppedManually {
			for i := range w.errListeners {
				w.errListeners[i](errUnexpected)
			}
		}
	}()

	for {
		select {
		case stop := <-w.stopChan:
			if stop {
				w.hasStoppedManually = true
				return
			}

		}
	}
}
