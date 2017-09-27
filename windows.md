## Go Report Card on MS Windows

You just use it is as usual (no need for the unix scripts).

### the third party tools

- You ***must*** have [`gometalinter`]() installed and available in Your `%path%`

```
go get github.com/alecthomas/gometalinter
go install github.com/alecthomas/gometalinter
```

- You *may* have the `wc` word count unix tool installed and available in Your `%path%` (e.g. from GnuWin)

### the current / working directory

As of this writing, `goreportcard.exe` works from your *current directory* (`os.Getwd()`). Thus:

No matter what You use:

	go run github.com/gojp/goreportcard 

	go build github.com/gojp/goreportcard

	go install github.com/gojp/goreportcard

- You ***must*** start `goreportcard.exe` from it's directory as current directory, so he can find his data.

- You ***must*** have the subdiretories `assets` and `templates` directly below Your `goreportcard.exe` (use `xcopy` or `mklink /D`)

- Note: The database file `goreportcard.db` and the subdirectory `_repos` for the repository cache are created upon first start.

If You have any [severe issue](https://github.com/gojp/goreportcard/issues), members of our community such as [`@GoLangsam`](https://github.com/GoLangsam) will try to help.

Note: As we (the original authors) have no way to test it on Windows, we do not support this OS personally.

