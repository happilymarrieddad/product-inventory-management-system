





go: added nhooyr.io/websocket v1.8.10
go install github.com/onsi/ginkgo/v2/ginkgo
../../go/pkg/mod/github.com/onsi/ginkgo/v2@v2.15.0/ginkgo/internal/profiles_and_reports.go:12:2: missing go.sum entry for module providing package github.com/google/pprof/profile (imported by github.com/onsi/ginkgo/v2/ginkgo/internal); to add:
	go get github.com/onsi/ginkgo/v2/ginkgo/internal@v2.15.0
../../go/pkg/mod/github.com/onsi/ginkgo/v2@v2.15.0/ginkgo/labels/labels_command.go:15:2: missing go.sum entry for module providing package golang.org/x/tools/go/ast/inspector (imported by github.com/onsi/ginkgo/v2/ginkgo/labels); to add:
	go get github.com/onsi/ginkgo/v2/ginkgo/labels@v2.15.0
make: *** [Makefile:4: install.deps] Error 1
☁  old-world-api [master] 


if you see the above run `go mod tidy` and try again