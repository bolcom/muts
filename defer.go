package muts

import "sync"

var deferList = new(deferFuncs)

type deferFuncs struct {
	sync.Mutex
	tasks []string
}

func (d *deferFuncs) add(task string) {
	d.Lock()
	defer d.Unlock()
	d.tasks = append(d.tasks, task)
}

func (d *deferFuncs) run() {
	d.Lock()
	defer d.Unlock()
	RunTasks(d.tasks...)
	d.tasks = []string{}
}

// Defer adds a task to the list of tasks that are run when CallDeferTasks is called
func Defer(task string) {
	deferList.add(task)
}

// CallDeferTasks() runs each task from the list. After this the list is empty.
func CallDeferTasks() {
	deferList.run()
}
