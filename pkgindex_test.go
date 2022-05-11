package gomod

import (
	"go/build"
	"runtime"
	"testing"
)

func TestGoroot(t *testing.T) {
	var pkgs PathPkgsIndex
	pkgs.LoadIndex(build.Default, runtime.GOROOT())
	pkgs.Sort()
	for _, v := range pkgs.Indexs {
		for _, pkg := range v.Pkgs {
			if pkg.IsCommand() {
				continue
			}
			//	t.Log(pkg.ImportPath, pkg.Goroot)
		}
	}
}
