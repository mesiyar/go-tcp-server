package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	cli, err := net.Dial("tcp", "127.0.0.1:8099")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			var body string
			fmt.Scanln(&body)
			cli.Write([]byte(body))
		}
	}()

	go func() {
		for {
			buf := [512]byte{}
			n, err := cli.Read(buf[:])
			if err != nil {
				fmt.Println("recv failed, err:", err)
				return
			}
			fmt.Println(string(buf[:n]))
		}
	}()

	select {

	}
}
