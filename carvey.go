package main

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
)

const DST = "./pics"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: carvy <filepath>")
		os.Exit(1)
	}

	fmt.Println("Carving PNG files .....")

	_, err := os.Stat(DST)
	if os.IsNotExist(err) {
		err = os.MkdirAll(DST, os.ModeDir)
		handle(err)
	}

	fpath := os.Args[1]
	finfo, err := os.Stat(fpath)
	handle(err)

	src, err := os.Open(fpath)
	handle(err)

	err = carvePng(DST, src, finfo.Size())
	handle(err)
}

func handle(err error) {
	if err != nil {
		fmt.Printf("\n\n%v\n\n", err)
		os.Exit(1)
	}
}

func carvePng(dstPath string, read *os.File, size int64) error {
	buff := make([]byte, 1)
	var counter int8
	var carved []byte

	if _, err := read.Seek(0, 0); err != nil {
		return err
	}

	for i := int64(0); i < size; i++ {
		if _, err := read.Read(buff); err != nil {
			return err
		}

		switch counter {
		case 0:
			if buff[0] == 0x89 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 1:
			if buff[0] == 0x50 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 2:
			if buff[0] == 0x4e {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 3:
			if buff[0] == 0x47 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 4:
			if buff[0] == 0xae {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
			}
		case 5:
			if buff[0] == 0x42 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
				counter--
			}
		case 6:
			if buff[0] == 0x60 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
				counter -= 2
			}
		case 7:
			if buff[0] == 0x82 {
				carved = append(carved, buff[0])
				if err := writecarved(dstPath, "png", &carved); err != nil {
					return err
				}
				carved = nil
				counter = 0
			} else {
				carved = append(carved, buff[0])
				counter -= 3
			}
		}

	}

	return nil
}

func writecarved(dstPath, ext string, data *[]byte) error {
	name := filepath.Join(dstPath, getimgname(10)+"."+ext)
	if err := os.WriteFile(name, *data, 0644); err != nil {
		return err
	}
	_, fname := filepath.Split(name)
	fmt.Printf("\rImage file found: %s", fname)
	return nil
}

func getimgname(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
