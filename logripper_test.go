package logripper_test

import (
	"fmt"
	"io"
	"logripper"
	"os"
	"strings"
	"testing"
)

func writeToFile(t *testing.T, str, fileName string) {
	fileA, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = fileA.WriteString(str)
	if err != nil {
		t.Fatal(err.Error())
	}

	fileA.Close() // Flush file
}

func TestTrimFileBy(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	str := `	aaaa
	bbb
	ccc
	ddd
	eee`
	testCase := func(str string) {
		strArr := strings.Split(str, "\n")
		fileName := fmt.Sprintf("%s/%s", dir, "temp")
		defer os.Remove(fileName)

		for i := 1; i <= len(strArr); i++ {

			writeToFile(t, str, fileName)

			fileA, err := os.OpenFile("temp", os.O_RDWR, 0755) // open again
			if err != nil {
				t.Fatal(err.Error())
			}

			opt := logripper.Options{LinesLeft: int64(i), Bufsize: 10, File: fileA}

			err = logripper.TrimFileBy(opt)
			if err != nil {
				t.Fatal(err.Error())
			}

			fileA, err = os.Open("temp")
			if err != nil {
				t.Fatal(err.Error())
			}

			stat, err := fileA.Stat()
			if err != nil {
				t.Fatal(err.Error())
			}

			buf := make([]byte, stat.Size())
			fileA.Seek(0, io.SeekStart)
			_, err = fileA.Read(buf)
			if err != nil {
				t.Fatal(err.Error())
			}

			if strings.TrimSpace(string(buf)) != strings.TrimSpace(strings.Join(strArr[len(strArr)-i:], "\n")) {

				t.Fatalf(
					"File trimmed does not match expected result. Got (%s) expected (%s) index (%d)",
					strings.TrimSpace(string(buf)),
					strings.TrimSpace(strings.Join(strArr[len(strArr)-i:], "\n")),
					i,
				)
			}

			os.Remove(fileA.Name())
		}
	}

	testCase(str)
	str = ""
	for _, c := range strings.Split("ABCEDEFGJIHKLMNIOPQRSTUVWXYZabcdefghjkaifusjndsjb", "") {
		str += (strings.Repeat(c, 100) + "\n")
	}
	testCase(str)

}
