package main

import (
	"fmt"

	"github.com/rabbitmq/internal/processor"
)

func main() {
	fmt.Println(`==> Starting server`)
	processor.Run()
	fmt.Println(`==> Ending server`)
}
