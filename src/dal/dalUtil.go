package dal

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

func (d dalUtil) getProcData(data *procParser.Data, splitFromKey string, parentMapKey string) map[string]map[string][]string {
	keyValArr := make(map[string]map[string][]string)
	var (
		m         map[string][]string
		parentKey string
	)
	for i := 0; i < len(data.Lines); i++ {
		if len(data.Lines[i].Values) == 0 {
			continue
		}
		key := data.Lines[i].Values[0]
		if parentKey != "" {
			keyValArr[parentKey] = m
			parentKey = ""
		}
		if key == splitFromKey {
			m = make(map[string][]string)
		}
		if m != nil && key == parentMapKey {
			parentKey = data.Lines[i].Values[1]
		}
		m[key] = data.Lines[i].Values
	}
	//Add the last/first set to the map array
	if m != nil && parentKey != "" {
		keyValArr[parentKey] = m
	}
	return keyValArr
}

func (d dalUtil) getDataFromMap(key string, data *procParser.Data) int64 {
	if _, exists := data.Map[key]; !exists {
		return 0
	}
	val, err := procParser.GetInt64(data.Map[key].Values[1])
	if err != nil {
		return 0
	}
	// unit conversion to Bytes
	return procParser.GetBytes(val, data.Map[key].Values[2])
}
