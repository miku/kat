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

const help = `kat - Preview.app for the command line
 _                        _            
(_)                      (_)           
(_)     _   _  _  _    _ (_) _  _      
(_)   _(_) (_)(_)(_) _(_)(_)(_)(_)     
(_) _(_)    _  _  _ (_)  (_)           
(_)(_)_   _(_)(_)(_)(_)  (_)     _     
(_)  (_)_(_)_  _  _ (_)_ (_)_  _(_)    
(_)    (_) (_)(_)(_)  (_)  (_)(_)    

Plain text, directories, PDF, JPG, PNG, MARC, zip, tgz, rar, mp3, odt, doc, docx, xlsx, tar, tar.gz, dmg.

$ kat FILE
`

const Version = "0.1.0"

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

type Zipfile struct {
	File
}

func (f *Zipfile) View() ([]byte, error) {
	return exec.Command("unzip", "-l", f.Name).Output()
}

type TGZ struct {
	File
}

func (f *TGZ) View() ([]byte, error) {
	return exec.Command("tar", "tf", f.Name).Output()
}

type MP3 struct {
	File
}

func (f *MP3) View() ([]byte, error) {
	return exec.Command("mp3info", "-x", f.Name).Output()
}

type Rar struct {
	File
}

func (f *Rar) View() ([]byte, error) {
	return exec.Command("unrar", "l", f.Name).Output()
}

type ODT struct {
	File
}

func (f *ODT) View() ([]byte, error) {
	return exec.Command("docd", "-input", f.Name).Output()
}

type Word struct {
	File
}

func (f *Word) View() ([]byte, error) {
	return exec.Command("antiword", f.Name).Output()
}

type XLSX struct {
	File
}

func (f *XLSX) View() ([]byte, error) {
	// https://git.io/vXOHi
	return exec.Command("xlsx2tsv.py", f.Name).Output()
}

type DMG struct {
	File
}

func (f *DMG) View() ([]byte, error) {
	return exec.Command("hdiutil", "imageinfo", f.Name).Output()
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
	case strings.HasSuffix(s, ".zip"):
		return &Zipfile{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tgz"):
		return &TGZ{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tar"):
		return &TGZ{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tar.gz"):
		return &TGZ{File{Name: s}}, nil
	case strings.HasSuffix(s, ".mp3"):
		return &MP3{File{Name: s}}, nil
	case strings.HasSuffix(s, ".rar"):
		return &Rar{File{Name: s}}, nil
	case strings.HasSuffix(s, ".odt"):
		return &ODT{File{Name: s}}, nil
	case strings.HasSuffix(s, ".docx"):
		return &ODT{File{Name: s}}, nil
	case strings.HasSuffix(s, ".doc"):
		return &Word{File{Name: s}}, nil
	case strings.HasSuffix(s, ".xlsx"):
		return &XLSX{File{Name: s}}, nil
	case strings.HasSuffix(s, ".dmg"):
		return &DMG{File{Name: s}}, nil
	default:
		return &File{Name: s}, nil
	}
	return nil, ErrUnsupportedFiletype
}

func main() {
	version := flag.Bool("version", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println(help)
		os.Exit(0)
	}

	for _, arg := range flag.Args() {
		if st, err := os.Stat(arg); err == nil {
			if st.IsDir() {
				out, err := exec.Command("tree", "-sh", arg).Output()
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
