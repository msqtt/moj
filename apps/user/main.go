package main

func main() {
	// inject dependency
	app := InitializeApplication()

	app.Start()
}