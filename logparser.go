package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

type Header struct {
	Op     int32
	Tty    int32
	Length int32
	Dir    int32
	Sec    uint32
	Usec   uint32
}

const (
	OP_OPEN = iota + 1
	OP_CLOSE
	OP_WRITE
	OP_EXEC
)

const (
	TYPE_INPUT = iota + 1
	TYPE_OUTPUT
	TYPE_INTERACT
)

var sssize = 24 // iLiiLL

func readNextBytes(file *os.File, number int32) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func GetFileSize(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func Playlog(path string, inputOnly bool, bothDir bool, colorify bool, maxdelay float64) {

	currtty := 0
	var prevtime float64 = 0.0
	prefdir := 0

	var filePos int32 = 0

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	fileSize, err := GetFileSize(path)

	color := false
	colorCode := ""

	for {
		header := Header{}
		data := readNextBytes(file, 24)

		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &header)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			break
		}

		// overjump header
		filePos += 24

		data = readNextBytes(file, header.Length)
		filePos += header.Length

		if currtty == 0 {
			currtty = int(header.Tty)
		}

		if int(header.Tty) == int(currtty) && int(header.Op) == int(OP_WRITE) {
			// the first stream seen is considered 'output'
			if prefdir == 0 {
				prefdir = int(header.Dir)

				// use the other direction
				if inputOnly {
					prefdir = TYPE_INPUT
					if int(header.Dir) == TYPE_INPUT {
						prefdir = TYPE_OUTPUT
					}
				}
			}
			if int(header.Dir) == TYPE_INTERACT {
				colorCode = "\033[36m"
				color = true
			} else if int(header.Dir) == TYPE_INPUT {
				colorCode = "\033[33m"
				color = true
			}
			if int(header.Dir) == prefdir || bothDir {
				curtime := float64(header.Sec+header.Usec) / 1000000

				if prevtime != 0.0 {
					var sleeptime float64 = curtime - float64(prevtime)
					if sleeptime > maxdelay {
						sleeptime = maxdelay
					}
					if maxdelay > 0 {
						sleepfinal := (sleeptime * 1000)
						time.Sleep(time.Duration(sleepfinal) * time.Millisecond)
					}
				}

				prevtime = curtime

				if colorify && color {
					fmt.Println(colorCode)
				}
				fmt.Printf(" %s", data)

				if colorify && color {
					fmt.Println("\033[0m")
					color = false
				}

			} else if header.Tty == int32(currtty) && header.Op == OP_CLOSE {
				break
			}
		}

		if int32(filePos) == int32(fileSize) {
			break
		}

	}
}

func main() {
	Playlog("./tty/LONG", true, true, true, 3.0)
}
