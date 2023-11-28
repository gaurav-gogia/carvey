package main

import (
	"bytes"
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

	err = carvePng(DST, os.Args[1])
	handle(err)
}

func handle(err error) {
	if err != nil {
		fmt.Printf("\n\n%v\n\n", err)
		os.Exit(1)
	}
}

var (
	HEADER = []byte{0x89, 0x50, 0x4e, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	FOOTER = []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
)

func carvePng(dstPath string, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	var count int
	start := bytes.Index(data, HEADER)
	end := bytes.Index(data, FOOTER)

	for start != -1 || end != -1 {
		end = end + len(FOOTER)
		if start < end {
			count++
			fdata := data[start : end-1]
			err = writeCarved(count, dstPath, fdata)
			if err != nil {
				return err
			}
		}

		data = data[end:]
		start = bytes.Index(data, HEADER)
		end = bytes.Index(data, FOOTER)
	}

	return nil
}

func writeCarved(count int, dstPath string, data []byte) error {
	fname := fmt.Sprintf("IMAGE_%d.png", count)
	fpath := filepath.Join(dstPath, fname)
	if err := os.WriteFile(fpath, data, 0644); err != nil {
		return err
	}
	fmt.Printf("\rCarved: %s", fname)
	return nil
}
