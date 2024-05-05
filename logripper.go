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

func max(A, B int64) int64 {
	if A > B {
		return A
	}
	return B
}

func TrimFileBy(opt Options) error {

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
	fileSeekPos := finfo.Size()
	fileSize := finfo.Size()

outer:
	for newLineCount < opt.LinesLeft {
		if fileSize <= 0 {
			break
		}

		_, err = opt.File.Seek(max(0, fileSeekPos-int64(opt.Bufsize)), io.SeekStart)
		if err != nil {
			break
		}
		n, err := opt.File.Read(buf)
		if err != nil {
			break
		}
		for index := n - 1; index >= 0; index-- {
			if buf[index] == NewLineConst {
				newLineCount++
			}
			if newLineCount >= opt.LinesLeft {
				fileSeekPos -= int64(n - index)
				fileSize -= int64(n - index)
				break outer
			}
		}

		fileSize -= int64(n)
		fileSeekPos -= int64(n)
	}

	fileSeekPos = max(0, fileSeekPos)
	_, err = opt.File.Seek(fileSeekPos, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(tempFile, opt.File)
	if err != nil {
		return err
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
