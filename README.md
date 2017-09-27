# go_n_find

CLI tool for finding files and folders and renaming them

Build Command:

```cmd

go build -ldflags="-X github.com/mpopadic/go_n_find/cmd.Version=v1.0.0 -X github.com/mpopadic/go_n_find/cmd.BuildTime=${BUILD_TIME}"

```

```cmd

gox -output="bin/{{.OS}}_{{.Arch}}/{{.Dir}}" -osarch="linux/amd64 windows/amd64" -ldflags="-X github.com/mpopadic/go_n_find/cmd.Version=v1.0.0 -X github.com/mpopadic/go_n_find/cmd.BuildTime=${BUILD_TIME}"

```