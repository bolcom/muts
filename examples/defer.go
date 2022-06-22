package main

import (
	"log"

	. "github.com/bolcom/muts"
)

// go run defer.go fail
// go run defer.go panic
// go run defer.go test

func main() {
	Task("panic", func() {
		panic("panic")
	})

	Task("fail", func() {
		Chdir("/fdsafdsa")
	})

	Task("abort1", func() {
		log.Println("about to abort 1")
	})
	Defer("abort1")

	Task("abort2", func() {
		log.Println("about to abort 2")
	})
	Defer("abort2")

	RunTasksFromArgs()
}
