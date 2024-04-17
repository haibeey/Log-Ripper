package main

import (
	"flag"
	"log"
	"logripper"
	"os"
	"path"
	"strings"
)

const bufSize = 1000

var (
	linesLeft   = flag.Int("n", -1, "The number of lines to be left in the files")
	extension   = flag.String("ext", "", "The extension name of files you want to be trimmed")
	logFilePath = flag.String("path", "", "The path of the directoty or file you want to trim")
)

func init() {
	flag.Parse()
}

func main() {
	if len(*logFilePath) <= 0   {
		log.Println("Empty path . Aborting ")
		return
	}

	if *linesLeft <= 0{
		log.Println("Invalid Value for number of lines left")
		return
	}
	filePath := path.Clean(*logFilePath)
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	fileStat, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}


	if fileStat.IsDir() {

		files, err := os.ReadDir(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			f, err := os.OpenFile(file.Name(), os.O_RDWR, 0755)
			if err != nil {
				continue
			}

			if *extension != "" && !strings.HasSuffix(file.Name(), *extension) {
				log.Printf("File %s extension does not match\n",file.Name())
			}
			logripper.TrimFileBy(logripper.Options{File: f,Bufsize: bufSize,LinesLeft: int64(*linesLeft)})
		}
	} else {
		f, err := os.OpenFile(filePath, os.O_RDWR, 0755)
		if err != nil {
			return
		}
		if *extension != "" && !strings.HasSuffix(file.Name(), *extension) {
			log.Println("File extension does not match")
		}

		logripper.TrimFileBy(logripper.Options{File: f,Bufsize: bufSize,LinesLeft: int64(*linesLeft)})
	}
}

