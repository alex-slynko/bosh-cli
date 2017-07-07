package cmd

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/cppforlife/go-patch/patch"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

type OpsFileArg struct {
	FS boshsys.FileSystem

	Ops patch.Ops
}

func (a *OpsFileArg) UnmarshalFlag(filePath string) error {
	if len(filePath) == 0 {
		return bosherr.Errorf("Expected file path to be non-empty")
	}

	bytes, err := a.FS.ReadFile(filePath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Reading ops file '%s'", filePath)
	}

	var opDefs []patch.OpDefinition

	err = yaml.Unmarshal(bytes, &opDefs)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deserializing ops file '%s'", filePath)
	}

	ops, err := patch.NewOpsFromDefinitions(opDefs)
	if err != nil {
		return bosherr.WrapErrorf(err, "Building ops")
	}

	(*a).Ops = ops

	return nil
}

func (a *OpsFileArg) Complete(match string) []flags.Completion {
	files, _ := a.FS.Glob(match + "*")
	fmt.Printf("%#v %s", a, match)
	ret := make([]flags.Completion, len(files))

	for i, v := range files {
		ret[i].Item = v
	}

	return ret
}
