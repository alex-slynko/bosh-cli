package template

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	flags "github.com/jessevdk/go-flags"
)

type VarFileArg struct {
	FS boshsys.FileSystem

	Vars StaticVariables
}

func (a *VarFileArg) UnmarshalFlag(data string) error {
	pieces := strings.SplitN(data, "=", 2)
	if len(pieces) != 2 {
		return bosherr.Errorf("Expected var '%s' to be in format 'name=path'", data)
	}

	if len(pieces[0]) == 0 {
		return bosherr.Errorf("Expected var '%s' to specify non-empty name", data)
	}

	if len(pieces[1]) == 0 {
		return bosherr.Errorf("Expected var '%s' to specify non-empty path", data)
	}

	absPath, err := a.FS.ExpandPath(pieces[1])
	if err != nil {
		return bosherr.WrapErrorf(err, "Getting absolute path '%s'", pieces[1])
	}

	bytes, err := a.FS.ReadFile(absPath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Reading variable from file '%s'", absPath)
	}

	(*a).Vars = StaticVariables{pieces[0]: string(bytes)}

	return nil
}

func (a *VarFileArg) Complete(match string) []flags.Completion {
	files, _ := a.FS.Glob(match + "*")
	ret := make([]flags.Completion, len(files))

	for i, v := range files {
		ret[i].Item = v
	}

	return ret
}
