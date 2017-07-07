package template

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

type VarsFileArg struct {
	FS boshsys.FileSystem

	Vars StaticVariables
}

func (a *VarsFileArg) UnmarshalFlag(filePath string) error {
	if len(filePath) == 0 {
		return bosherr.Errorf("Expected file path to be non-empty")
	}

	bytes, err := a.FS.ReadFile(filePath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Reading variables file '%s'", filePath)
	}

	var vars StaticVariables

	err = yaml.Unmarshal(bytes, &vars)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deserializing variables file '%s'", filePath)
	}

	(*a).Vars = vars

	return nil
}

func (a *VarsFileArg) Complete(match string) []flags.Completion {
	files, _ := a.FS.Glob(match + "*")
	ret := make([]flags.Completion, len(files))

	for i, v := range files {
		ret[i].Item = v
	}

	return ret
}
