package main

import (
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

Plain text, directories, PDF, JPG, PNG, gif, MARC, zip, tgz, rar, zstd, mp3,
odt, doc, docx, xlsx, tar, tar.gz, gz, dmg, djvu, deb, rpm.

$ kat FILE
`

// Version of application.
const Version = "0.1.4"

// Viewer emits a view of a thing.
type Viewer interface {
	View() ([]byte, error)
}

// File has a name which often will be a path.
type File struct {
	Name string
}

// View is the generic via cat.
func (f *File) View() ([]byte, error) {
	if _, err := exec.LookPath("cat"); err != nil {
		return nil, fmt.Errorf("cat is required")
	}
	return exec.Command("cat", f.Name).Output()
}

// PDF file.
type PDF struct {
	File
}

// View PDF as text.
func (f *PDF) View() ([]byte, error) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return nil, fmt.Errorf("pdftotext is required")
	}
	return exec.Command("pdftotext", f.Name, "-").Output()
}

// Image file.
type Image struct {
	File
}

// View image with catimg.
func (f *Image) View() ([]byte, error) {
	if _, err := exec.LookPath("catimg"); err != nil {
		return nil, fmt.Errorf("catimg is required")
	}
	return exec.Command("catimg", "-w", "192", f.Name).Output()
}

// BinaryMARC21 bibliographic data.
type BinaryMARC21 struct {
	File
}

// View with yaz.
func (f *BinaryMARC21) View() ([]byte, error) {
	if _, err := exec.LookPath("yaz-marcdump"); err != nil {
		return nil, fmt.Errorf("yaz-marcdump is required")
	}
	return exec.Command("yaz-marcdump", f.Name).Output()
}

// Zipfile compressed file.
type Zipfile struct {
	File
}

// View lists zip file contents.
func (f *Zipfile) View() ([]byte, error) {
	if _, err := exec.LookPath("unzip"); err != nil {
		return nil, fmt.Errorf("unzip is required")
	}
	return exec.Command("unzip", "-l", f.Name).Output()
}

// Tarfile archive.
type Tarfile struct {
	File
}

// View lists tar contents.
func (f *Tarfile) View() ([]byte, error) {
	if _, err := exec.LookPath("tar"); err != nil {
		return nil, fmt.Errorf("tar is required")
	}
	return exec.Command("tar", "tf", f.Name).Output()
}

// MP3 audio.
type MP3 struct {
	File
}

// View display file information.
func (f *MP3) View() ([]byte, error) {
	if _, err := exec.LookPath("mp3info"); err != nil {
		return nil, fmt.Errorf("mp3info is required")
	}
	return exec.Command("mp3info", "-x", f.Name).Output()
}

// Rar archive.
type Rar struct {
	File
}

// View lists rar contents.
func (f *Rar) View() ([]byte, error) {
	if _, err := exec.LookPath("unrar"); err != nil {
		return nil, fmt.Errorf("unrar is required")
	}
	return exec.Command("unrar", "l", f.Name).Output()
}

// ODT document.
type ODT struct {
	File
}

// View extracts text.
func (f *ODT) View() ([]byte, error) {
	if _, err := exec.LookPath("docd"); err != nil {
		return nil, fmt.Errorf("docd is required")
	}
	return exec.Command("docd", "-input", f.Name).Output()
}

// Word document.
type Word struct {
	File
}

// View extracts text from Word.
func (f *Word) View() ([]byte, error) {
	if _, err := exec.LookPath("antiword"); err != nil {
		return nil, fmt.Errorf("antiword is required")
	}
	return exec.Command("antiword", f.Name).Output()
}

// XLSX spreadsheet.
type XLSX struct {
	File
}

// View lists the (first) spreadsheet.
func (f *XLSX) View() ([]byte, error) {
	// https://git.io/vXOHi
	if _, err := exec.LookPath("xlsx2tsv.py"); err != nil {
		return nil, fmt.Errorf("xlsx2tsv is required, https://git.io/vXOHi")
	}
	out, err := exec.Command("xlsx2tsv.py", f.Name).Output()
	if err == nil {
		return out, err
	}
	return exec.Command("xlsx2tsv.py", f.Name, "1").Output()
}

// DMG disk image.
type DMG struct {
	File
}

// View list info about the image.
func (f *DMG) View() ([]byte, error) {
	if _, err := exec.LookPath("hdiutil"); err != nil {
		return nil, fmt.Errorf("hdiutil is required")
	}
	return exec.Command("hdiutil", "imageinfo", f.Name).Output()
}

// Djvu document.
type Djvu struct {
	File
}

// View extracts text.
func (f *Djvu) View() ([]byte, error) {
	if _, err := exec.LookPath("djvutxt"); err != nil {
		return nil, fmt.Errorf("djvutxt is required")
	}
	return exec.Command("djvutxt", f.Name).Output()
}

// DebianPackage package.
type DebianPackage struct {
	File
}

// View file.
func (f *DebianPackage) View() ([]byte, error) {
	if _, err := exec.LookPath("dpkg"); err != nil {
		return nil, fmt.Errorf("dpkg is required")
	}
	return exec.Command("dpkg", "-c", f.Name).Output()
}

// RPM package.
type RPM struct {
	File
}

// View file.
func (f *RPM) View() ([]byte, error) {
	if _, err := exec.LookPath("rpm"); err != nil {
		return nil, fmt.Errorf("rpm is required")
	}
	return exec.Command("rpm", "-qplv", f.Name).Output()
}

// Gzip file.
type Gzip struct {
	File
}

// View file.
func (f *Gzip) View() ([]byte, error) {
	if _, err := exec.LookPath("gunzip"); err != nil {
		return nil, fmt.Errorf("gunzip is required")
	}
	return exec.Command("gunzip", "-c", f.Name).Output()
}

// Zstd file.
type Zstd struct {
	File
}

// View file.
func (f *Zstd) View() ([]byte, error) {
	if _, err := exec.LookPath("zstd"); err != nil {
		return nil, fmt.Errorf("zstd is required")
	}
	return exec.Command("zstd", "-c", "-d", "-T0", f.Name).Output()
}

// DispatchFile chooses a viewer for a given filename.
func DispatchFile(s string) (Viewer, error) {
	switch {
	case strings.HasSuffix(s, ".pdf"):
		return &PDF{File{Name: s}}, nil
	case strings.HasSuffix(s, ".jpg"):
		return &Image{File{Name: s}}, nil
	case strings.HasSuffix(s, ".png"):
		return &Image{File{Name: s}}, nil
	case strings.HasSuffix(s, ".gif"):
		return &Image{File{Name: s}}, nil
	case strings.HasSuffix(s, ".mrc"):
		return &BinaryMARC21{File{Name: s}}, nil
	case strings.HasSuffix(s, ".zip"):
		return &Zipfile{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tgz"):
		return &Tarfile{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tar"):
		return &Tarfile{File{Name: s}}, nil
	case strings.HasSuffix(s, ".tar.gz"):
		return &Tarfile{File{Name: s}}, nil
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
	case strings.HasSuffix(s, ".djvu"):
		return &Djvu{File{Name: s}}, nil
	case strings.HasSuffix(s, ".deb"):
		return &DebianPackage{File{Name: s}}, nil
	case strings.HasSuffix(s, ".rpm"):
		return &RPM{File{Name: s}}, nil
	case strings.HasSuffix(s, ".gz"):
		return &Gzip{File{Name: s}}, nil
	case strings.HasSuffix(s, ".zst"):
		return &Zstd{File{Name: s}}, nil
	default:
		return &File{Name: s}, nil
	}
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
