kat
===

Kat is like [Preview.app](https://en.wikipedia.org/wiki/Preview_(macOS)) for the command-line.

Installation
------------

```
$ go get github.com/miku/kat/cmd/kat
```

Idea
----

This project is nothing a shell script could not do. Displaying is outsourced suitable command-line tools.
My hope is to extend the list of supported files as the need arises and maybe add verbosity options, that do *the right thing*
depending on the filetype.

Usage
-----

```
kat - Preview.app for the command line
 _                        _            
(_)                      (_)           
(_)     _   _  _  _    _ (_) _  _      
(_)   _(_) (_)(_)(_) _(_)(_)(_)(_)     
(_) _(_)    _  _  _ (_)  (_)           
(_)(_)_   _(_)(_)(_)(_)  (_)     _     
(_)  (_)_(_)_  _  _ (_)_ (_)_  _(_)    
(_)    (_) (_)(_)(_)  (_)  (_)(_)    

Plain text, directories, PDF, JPG, PNG, gif, MARC, zip, tgz, rar, mp3, odt,
doc, docx, xlsx, tar, tar.gz, dmg, djvu, deb, rpm.

$ kat FILE
```

[![](docs/63qroppw4y5irohpkcyqpck0t.gif)](https://asciinema.org/a/63qroppw4y5irohpkcyqpck0t)
