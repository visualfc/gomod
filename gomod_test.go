package gomod

import (
	"fmt"
	"go/build"
	"path/filepath"
	"runtime"
	"testing"
)

func TestStd(t *testing.T) {
	goroot := runtime.GOROOT()
	pkg, err := Load(filepath.Join(goroot, "src/net"), &build.Default)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Root().Path != "std" {
		t.Fatal(pkg.Root().Path)
	}
	if pkg.Root().Dir != filepath.Join(goroot, "src") {
		t.Fatal(pkg.Root().Dir)
	}
	_, dir, found := pkg.Lookup("golang.org/x/net/dns/dnsmessage")
	if !found {
		t.Fail()
	}
	if dir != filepath.Join(goroot, "src/vendor/golang.org/x/net/dns/dnsmessage") {
		t.Fatal(dir)
	}
	fmt.Println("dep imports list")
	for _, dep := range pkg.DepImportList(true, true) {
		fmt.Println(dep)
	}
	fmt.Println("local imports list")
	for _, im := range pkg.LocalImportList(true) {
		fmt.Println(im)
	}
}
