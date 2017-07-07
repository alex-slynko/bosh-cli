package cmd

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	flags "github.com/jessevdk/go-flags"
)

type FileArg struct {
	ExpandedPath string
	FS           boshsys.FileSystem
}

func (a *FileArg) UnmarshalFlag(data string) error {
	expandedPath, err := a.FS.ExpandPath(data)
	if err != nil {
		return bosherr.WrapErrorf(err, "Checking file path")
	}
	a.ExpandedPath = expandedPath

	if a.FS.FileExists(expandedPath) {
		stat, err := a.FS.Stat(expandedPath)
		if err != nil {
			return bosherr.WrapErrorf(err, "Checking file path")
		}

		if stat.IsDir() {
			return bosherr.Errorf("Path must not be directory")
		}
	}

	return nil
}

func (a *FileArg) Complete(match string) []flags.Completion {
	files, _ := a.FS.Glob(match + "*")
	fmt.Printf("%#v %s", a, match)
	ret := make([]flags.Completion, len(files))

	for i, v := range files {
		ret[i].Item = v
	}

	return ret
}
