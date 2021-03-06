
set version=%1%
gox -output="bin\{{.OS}}_{{.Arch}}\{{.Dir}}" -osarch="linux/amd64 windows/amd64 darwin/amd64" -ldflags="-X github.com\mpopadic\go_n_find\cmd.Version=v%version%"

cd bin
cd windows_amd64 && zip ..\windows_amd64.zip *
cd ..
rd /s /q windows_amd64
cd linux_amd64 && zip ..\linux_amd64.zip *
cd ..
rd /s /q linux_amd64
cd darwin_amd64 && zip ..\darwin_amd64.zip *
cd ..
rd /s /q darwin_amd64