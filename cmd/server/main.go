package main

import (
	"fmt"

	"github.com/rabbitmq/internal/server"
)

func main() {
	fmt.Println(`==> Starting server`)
	server.Run()
	fmt.Println(`==> Ending server`)
}
