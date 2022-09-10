package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"note/settings"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	err := readArgs()
	if err != nil {
		log.Fatalln(err)
	}
}

func readDate() (date string) {
	year, month, day := time.Now().Date()
	date = fmt.Sprintf("%04[1]d-%02[2]d-%02[3]d", year, month, day)
	return
}

func createDir() (dir string, err error) {
	date := readDate()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir = fmt.Sprintf(
		"%[1]s/%[2]s/%[3]s",
		homeDir,
		"Documents/notes",
		date,
	)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	return
}

func readArgs() (err error) {
	option := ""
	if len(os.Args) < 2 {
		option = "open"
	} else {
		option = os.Args[1]
	}
	switch option {
	case "open", "-o", "new", "-n":
		var dir string
		dir, err = createDir()
		if err != nil {
			return
		}
		var fileName string
		if len(os.Args) < 3 {
			fileName = fmt.Sprintf(
				"%[1]s",
				"notes.txt",
			)
		} else {
			fileName = fmt.Sprintf(
				"%[1]s%[2]s",
				strings.Join(os.Args[2:], "-"),
				".txt",
			)
		}
		filePath := fmt.Sprintf(
			"%[1]s/%[2]s",
			dir,
			fileName,
		)
		var file *os.File
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return
		}
		var bytes []byte
		bytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			return
		}
		timeStamp := time.Now().Format(time.RFC850)
		var text string
		if len(bytes) > 0 {
			text = fmt.Sprintf(
				"\n%[1]s\n|\t\t\t%[2]s\t\t\t\t|\n%[3]s\n\n",
				"+-------------------------------------------------------------------------------+",
				timeStamp,
				"+-------------------------------------------------------------------------------+",
			)
		} else {
			text = fmt.Sprintf(
				"%[1]s\n|\t\t\t%[2]s\t\t\t\t|\n%[3]s\n\n",
				"+-------------------------------------------------------------------------------+",
				timeStamp,
				"+-------------------------------------------------------------------------------+",
			)
		}
		textBytes := []byte(text)
		_, err = file.Write(textBytes)
		if err != nil {
			return
		}
		err = file.Close()
		if err != nil {
			return
		}
		switch settings.Option.Editor {
		case "nano":
			position := strings.Count(text+string(bytes), "\n")
			position++
			err = runNano(filePath, position)
			if err != nil {
				return
			}
		case "gedit":
			err = runGedit(filePath)
			if err != nil {
				return
			}
		}
	default:
		err = settings.ErrOption
	}
	return
}

func runNano(filePath string, position int) (err error) {
	row := fmt.Sprintf("+%[1]d", position)
	cmd := exec.Command("nano", row, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	return
}

func runGedit(filePath string) (err error) {
	cmd := exec.Command("gedit", filePath)
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}
