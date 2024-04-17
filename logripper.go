package logripper

import (
	"io"
	"os"
)

const NewLineConst = 10

type Options struct {
	Bufsize   int64
	LinesLeft int64
	File      *os.File
}

func TrimFileBy(opt Options) error {

	seekPos := opt.Bufsize
	finfo, err := opt.File.Stat()
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	tempFileName := tempFile.Name()
	fileName := opt.File.Name()

	buf := make([]byte, opt.Bufsize)
	newLineCount := int64(0)

outer:
	for newLineCount < opt.LinesLeft {
		opt.File.Seek(finfo.Size()-int64(seekPos), io.SeekEnd)
		n, err := opt.File.Read(buf)
		if err != nil {
			break
		}
		for index := n - 1; index >= 0; index-- {
			if buf[index] == NewLineConst {
				newLineCount++
			}
			if newLineCount >= opt.LinesLeft {
				seekPos += int64(n - index)
				_, err = tempFile.Seek(0, io.SeekStart)
				if err != nil {
					return err
				}

				_, err = tempFile.Write(buf[index+1 : n])
				if err != nil {
					return err
				}

				break outer
			}
		}
		seekPos += int64(n)
		tempFile.Seek(0, io.SeekStart)
		tempFile.Write(buf[:n])
	}

	err = tempFile.Close() // flush.
	if err != nil {
		return err
	}

	tempFile, err = os.Open(tempFileName)
	if err != nil {
		return err
	}

	err = os.Truncate(opt.File.Name(), 0)
	if err != nil {
		return err
	}
	opt.File.Close()

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, tempFile)
	if err != nil {
		return err
	}

	tempFile.Close()
	opt.File.Close()

	os.Remove(tempFileName)

	return nil
}
