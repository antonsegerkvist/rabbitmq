package processor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func handleMessage(process *Process) error {

	originalLocation :=
		"storage/original/" +
			strconv.FormatInt(int64(process.Filesize), 10) +
			"_" +
			strconv.FormatUint(uint64(process.Checksum), 10) +
			process.Extension

	fmt.Println(originalLocation)

	previewLocation :=
		"storage/preview/" +
			strconv.FormatInt(int64(process.Filesize), 10) +
			"_" +
			strconv.FormatUint(uint64(process.Checksum), 10) +
			".png"

	fmt.Println(previewLocation)

	stringCommand :=
		"gm convert " +
			originalLocation +
			` -thumbnail 700x700 -gravity center -extent 700x700 ` +
			previewLocation

	args := strings.Split(stringCommand, " ")
	command := exec.Command(args[0], args[1:]...)

	err := command.Run()
	if err != nil {
		return err
	}

	return nil

}
