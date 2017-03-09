package linux

import (
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

type dalUtil struct {
	envDep env.FactoryEnv
}

func (d dalUtil) getCommandData(parser procParser.Parser, cfg procParser.Config, command string, arg ...string) (*procParser.Data, error) {
	reader, err := d.envDep.GetEnv().GetCommandReader(command, arg...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := parser.Parse(cfg, reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d dalUtil) getFileData(parser procParser.Parser, cfg procParser.Config, filePath string) (*procParser.Data, error) {
	reader, err := d.envDep.GetEnv().GetFileReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := parser.Parse(cfg, reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d dalUtil) getProcData(data *procParser.Data, splitFromKey string) []map[string][]string {
	var keyValArr []map[string][]string
	var m map[string][]string
	for i := 0; i < len(data.Lines); i++ {
		key := data.Lines[i].Values[0]
		if key == splitFromKey {
			if m != nil {
				keyValArr = append(keyValArr, m)
			}
			m = make(map[string][]string)
		}
		m[key] = data.Lines[i].Values
	}
	//Add the last/first set to the map array
	if m != nil {
		keyValArr = append(keyValArr, m)
	}
	return keyValArr
}
