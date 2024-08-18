package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("incus-compose: No valid command provided")
		return
	}
	args := os.Args[1:]

	switch args[0] {
	case "--version", "-v":
		fmt.Println("incus-compose 0.0.1")
	case "convert":
		if len(args) < 2 {
			fmt.Println("incus-compose convert: No docker compose file provided")
			return
		}
		convertDockerComposeToIncusCompose(args[1])
		fmt.Println("Summary: Convert docker-compose.yaml to incus-compose.yaml")
		fmt.Println("Usage: incus-compose convert <path/to/docker-compose.yaml> -o <path/to/incus-compose.yaml>(default: ./incus-compose.yaml)")
	case "up":
		up()
	default:
		fmt.Println("incus-compose: No valid command provided")
	}
}
