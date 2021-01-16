# kubensure

![Build](https://github.com/PhilRanzato/kubensure/workflows/Go/badge.svg)
![Lint](https://github.com/PhilRanzato/kubensure/workflows/GoLint/badge.svg)

Ensure consistency in Kubernetes cluster's resources

Additional commands cli-specific

```shell
dep ensure -add github.com/spf13/cobra/cobra
go get -u github.com/spf13/cobra/cobra
cd $GOPATH/src
cobra init github.com/PhilRanzato/kubensure/cli --pkg-name github.com/PhilRanzato/kubensure/cli
cd cli
go run main.go
# add a new command
cobra add test
go build -o $GOPATH/bin/kubensure.exe
kubensure -h
# add a new subcommand
cobra add subcmd -p cmd
```
