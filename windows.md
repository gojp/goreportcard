## Go Report Card on MS Windows

Go Report Card does not officially support Windows, but users report that it works.

- install `gometalinter` directly (see below)
- use it is as usual (no need for any `make` step as mentioned [above](README.me)).
- have `wc` available, if possible

If You have any [severe issue](https://github.com/gojp/goreportcard/issues), members of our community such as [`@GoLangsam`](https://github.com/GoLangsam) will try to help.

Note: As we (the original authors) have no way to test it on Windows, we do not support this OS personally.

---
### the third party tools

Install [`gometalinter`](https://github.com/alecthomas/gometalinter)

	go get github.com/alecthomas/gometalinter
	go install github.com/alecthomas/gometalinter

and have it available in Your `%path%`.

You will also like the `wc` word count unix tool to be available in Your `%path%` (e.g. from GnuWin)

---
### start / launch

#### in place

	cd /D %GOPATH%\src\github.com\gojp\goreportcard
	go build github.com/gojp/goreportcard
	.\goreportcard.exe

#### installed / any other place

	cd /D %GOPATH%\src\github.com\gojp\goreportcard
	go install github.com/gojp/goreportcard
	cd /D %GOPATH%\bin\

move `goreportcard.exe` to the directory of Your choice, and note the following.

---
### the current / working directory

As of this writing, `goreportcard.exe` works from your *current* directory (`os.Getwd()`). Thus:

- Have the subdiretories `assets` and `templates` availabe there (use `xcopy` or `mklink /D`)

- Start `goreportcard.exe` from it's directory as *current* directory, so he can find his data
  (You may use a shortcut file - just make sure to set the `Start in` property.)

- Note: The database file `goreportcard.db` and the subdirectory `_repos` for the repository cache are created upon first start.

---
Enjoy - and be a happy `ʕ◔ϖ◔ʔ`

---
