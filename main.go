package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var (
	trimLines   = flag.Int("n", 100, "The number of lines to lbe left in the files")
	extention   = flag.String("ext", "", "The extention name of files you want to be trimmed")
	logFilePath = flag.String("path", "", "The path of the directoty or file you want to trim")
)

func init() {
	flag.Parse()
}
func main() {
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

		files, err := ioutil.ReadDir(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			f, err := os.OpenFile(file.Name(), os.O_RDWR, 0755)
			if err != nil {
				continue
			}
			trimFileBy(f, *trimLines)
		}
	} else {
		f, err := os.OpenFile(filePath, os.O_RDWR, 0755)
		if err != nil {
			return
		}
		trimFileBy(f, *trimLines)
	}
}

func trimFileBy(file *os.File, count int) error {

	if *extention != "" && !strings.HasSuffix(file.Name(), *extention) {
		return nil
	}

	seekPos := int64(*trimLines)
	finfo, err := file.Stat()
	if err != nil {
		return err
	}
	if finfo.Size() < seekPos {
		return nil
	}

	buf := make([]byte, seekPos)
	newLineCount := 0

	for newLineCount < *trimLines {
		file.Seek(finfo.Size()-seekPos, os.SEEK_SET)
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		for _, b := range buf {
			if b == 10 {
				newLineCount++
			}
		}
		seekPos += int64(n)
	}

	fileContent := make([]byte, seekPos)
	file.Seek(finfo.Size()-seekPos, os.SEEK_SET)
	file.Read(fileContent)
	os.Truncate(file.Name(), int64(0))
	file.Seek(0, os.SEEK_SET)
	file.WriteString(strings.Trim(string(fileContent), "\n"))
	file.Close()

	return nil
}
