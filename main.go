package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/worker/persist"
)

func main() {

	e := engine.Engine{
		Scheduler:   &scheduler.Scheduler{},
		WorkerCount: 2000,
		ItemChan:    persist.ItemSaver(),
	}

	e.Run()

}


