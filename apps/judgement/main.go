package main

func main() {
	app := InitializeApplication()

	app.Start()
	defer app.Stop()
}
