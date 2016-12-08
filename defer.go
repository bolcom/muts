package muts

import "sync"

var deferList = new(deferFuncs)

var abortHook = new(sync.Once)

type deferFuncs struct {
	sync.Mutex
	tasks []string
}

func (d *deferFuncs) add(task string) {
	d.Lock()
	defer d.Unlock()
	d.tasks = append([]string{task}, d.tasks...)
}

func (d *deferFuncs) run() {
	d.Lock()
	defer d.Unlock()
	RunTasks(d.tasks...)
	d.tasks = []string{}
}

// Defer adds a task to the list of tasks that are run when Abort or panic is called.
func Defer(task string) {
	abortHook.Do(func() {
		oldAbort := Abort
		Abort = func(v ...interface{}) {
			deferList.run()
			oldAbort(v...)
		}
	})
	deferList.add(task)
}
