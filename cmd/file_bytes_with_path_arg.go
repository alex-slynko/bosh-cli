package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	flags "github.com/jessevdk/go-flags"
)

type FileBytesWithPathArg struct {
	FS boshsys.FileSystem

	Bytes []byte
	Path  string
}

func (a *FileBytesWithPathArg) UnmarshalFlag(data string) error {
	if len(data) == 0 {
		return bosherr.Errorf("Expected file path to be non-empty")
	}

	absPath, err := a.FS.ExpandPath(data)
	if err != nil {
		return bosherr.WrapErrorf(err, "Getting absolute path '%s'", data)
	}

	bytes, err := a.FS.ReadFile(absPath)
	if err != nil {
		return err
	}

	(*a).Bytes = bytes
	(*a).Path = absPath

	return nil
}

func (a *FileBytesWithPathArg) Complete(match string) []flags.Completion {
	files, _ := a.FS.Glob(match + "*")
	ret := make([]flags.Completion, len(files))

	for i, v := range files {
		ret[i].Item = v
	}

	return ret
}
