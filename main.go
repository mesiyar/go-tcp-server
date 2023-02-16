package main

import "talking/core"

func main() {
	server := core.NewServer("127.0.0.1", 8099)
	server.Start()
}
