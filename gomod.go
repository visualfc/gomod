package gomod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Module struct {
	Path      string       // module path
	Version   string       // module version
	Versions  []string     // available module versions (with -versions)
	Replace   *Module      // replaced by this module
	Time      *time.Time   // time version was created
	Update    *Module      // available update, if any (with -u)
	Main      bool         // is this the main module?
	Indirect  bool         // is this module only an indirect dependency of main module?
	Dir       string       // directory holding files for this module, if any
	GoMod     string       // path to go.mod file used when loading this module, if any
	GoVersion string       // go version used in module
	Retracted string       // retraction information, if any (with -retracted or -u)
	Error     *ModuleError // error loading module
}

func (m *Module) String() string {
	return fmt.Sprintf("{%v %v}", m.Path, m.GoMod)
}

type ModuleError struct {
	Err string // the error itself
}

type Package struct {
	List []*Module
}

func (p *Package) Root() *Module {
	return p.List[0]
}

func (p *Package) Lookup(pkg string) (path string, dir string, found bool) {
	for _, v := range p.List {
		if v.Path == pkg {
			return v.Path, v.Dir, true
		}
	}
	for _, v := range p.List {
		if strings.HasPrefix(pkg, v.Path+"/") {
			return pkg, filepath.Join(v.Dir, pkg[len(v.Path+"/"):]), true
		}
	}
	return "", "", false
}

func Load(dir string) (*Package, error) {
	var stdout, stderr bytes.Buffer
	stdout.WriteByte('[')
	cmd := exec.Command("go", "list", "-m", "-mod=readonly", "-json", "all")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	stdout.WriteByte(']')
	data := bytes.ReplaceAll(stdout.Bytes(), []byte("\n{"), []byte(",\n{"))
	var list []*Module
	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}
	// check std use vendor mod
	if list[0].Path == "std" {
		root := filepath.Join(list[0].Dir, "vendor")
		for i := 1; i < len(list); i++ {
			if list[i].Dir == "" {
				list[i].Dir = filepath.Join(root, list[i].Path)
			}
		}
	}
	return &Package{List: list}, nil
}
