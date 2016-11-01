package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var ErrUnsupportedFiletype = errors.New("unsupported file type")

type Viewer interface {
	View() ([]byte, error)
}

type File struct {
	Name string
}

func (f *File) View() ([]byte, error) {
	return exec.Command("cat", f.Name).Output()
}

type PDF struct {
	File
}

func (f *PDF) View() ([]byte, error) {
	return exec.Command("pdftotext", f.Name, "-").Output()
}

type Image struct {
	File
}

func (f *Image) View() ([]byte, error) {
	return exec.Command("catimg", "-w", "192", f.Name).Output()
}

type BinaryMARC21 struct {
	File
}

func (f *BinaryMARC21) View() ([]byte, error) {
	return exec.Command("marcdump", f.Name).Output()
}

func DispatchFile(s string) (Viewer, error) {
	switch {
	case strings.HasSuffix(s, ".pdf"):
		return &PDF{File{Name: s}}, nil
	case strings.HasSuffix(s, ".jpg"):
		return &Image{File{Name: s}}, nil
	case strings.HasSuffix(s, ".png"):
		return &Image{File{Name: s}}, nil
	case strings.HasSuffix(s, ".mrc"):
		return &BinaryMARC21{File{Name: s}}, nil
	default:
		return &File{Name: s}, nil
	}
	return nil, ErrUnsupportedFiletype
}

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		if st, err := os.Stat(arg); err == nil {
			if st.IsDir() {
				out, err := exec.Command("ls", arg).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(out))
			} else {
				v, err := DispatchFile(arg)
				if err != nil {
					log.Fatal(err)
				}
				b, err := v.View()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(b))
			}
		}
	}
}
