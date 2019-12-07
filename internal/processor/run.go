package processor

import (
	"fmt"
)

//
// Run runs the processor server.
//
func Run() {

	for true {
		fmt.Println(`==> Handling rabbitmq queue`)
		err := handleProcess()
		if err != nil {
			fmt.Println(`==> Error: ` + err.Error())
		}
	}

}
