package logripper_test

import (
	"fmt"
	"logripper"
	"os"
	"strings"
	"testing"
)

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

	strArr := strings.Split(str, "\n")
	fileName := fmt.Sprintf("%s/%s", dir, "temp")
	defer os.Remove(fileName)

	for i := 1; i <= len(strArr); i++ {

		fileA, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err.Error())
		}

		_, err = fileA.WriteString(str)
		if err != nil {
			t.Fatal(err.Error())
		}

		fileA.Close() // Flush file

		fileA, err = os.OpenFile("temp", os.O_RDWR, 0755) // open again
		if err != nil {
			t.Fatal(err.Error())
		}

		opt := logripper.Options{LinesLeft: int64(i), Bufsize: 100, File: fileA}

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
		_, err = fileA.Read(buf)
		if err != nil {
			t.Fatal(err.Error())
		}

		if string(buf) != strings.Join(strArr[len(strArr)-i:], "\n") {
			t.Fatalf(
				"File trimmed does not match expected result. Got (%s) expected (%s) index (%d)",
				string(buf),
				strings.Join(strArr[len(strArr)-i:], ""),
				i,
			)
		}

		os.Remove(fileA.Name())
	}

}
