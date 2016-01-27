/*
muts - Go package with utilities to create Make-like files in Go

Example of a make.go

	func main() {
		flag.Parse()
		if len(*DeployableVersion) == 0 || len(*BuildNumber) == 0 || len(*BuildDate) == 0 {
			log.Fatal("one of the required flags is missing")
		}
		Tasks["clean"] = taskClean
		Tasks["build"] = taskBuild
		Tasks["compile"] = taskCompile
		Tasks["compile_tester"] = taskCompileTester
		Tasks["compile_dbtester"] = taskCompileDBTester
		Tasks["compile_composer"] = taskCompileComposer
		Tasks["apidocs"] = func() { Call("tar -zcf ./target/swagger.tar.gz swagger") }
		Tasks["dashboard"] = func() { Call("tar -zcf ./target/gui.tar.gz dashboard") }
		Tasks["readme"] = func() { Call("cp -v readme.md ./target/") }
		Tasks["startsh"] = func() { Call("cp -v start.sh ./target/") }
		Tasks["db"] = taskArchiveDB
		Tasks["app"] = taskArchiveApp
		Tasks["standalone"] = taskStandalone
		Tasks["unit"] = func() { Call("godep go test ./internal/...") }
		Tasks["local"] = func() { Call("godep go build -o target/boqs main.go") }
		Tasks["publish"] = taskPublish
		RunTasksFromArgs()
	}

Use 	it like this

	go run make.go -version=1.1.1.1 -buildnumber=42 -builddate=`date +%Y:%m:%d.%H:%M:%S` build


Some background

This package contains a collection of small helper functions to create scripts the easy way.
Most of the time, shell scripting is fine but soon it can become complex once you need functions,loops and decision trees.
So why not use the Go language and its rich SDK to write real programs which can be organized much easier.

It all started with the Call function that mimics what you would write in a shell script.
It can be used both with a single line command and one that is composed of a list of strings.
The CallWait version lets you wait for the program to finish or return the process ID for stopping it later.

	Call("zip", "-q", "-r", fmt.Sprintf("%s/sql/boqs-db-%s.zip", versionDir, *DeployableVersion), ".")

Next, we added the concept of a simple task (without the dependencies).
A task is just a no-argument function.
By putting these tasks in the global Tasks map, you can execute them just by passing their names to your program:

	go run make.go clean build unit

The last feature to mention is the Workspace variable that refers to the directory in which the program was started.
Task execution may change this directory (Chdir) so to keep things simpler, the current directory is reset after each task.

Most functions will produce a log entry.
If an error occurs then the program exits (calling the Fatalln function).

*/
package muts
