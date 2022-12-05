package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	CURRENT_DIR_PREFIX = "./"
	ROOT_DIR_PREFIX    = "/"
	BOOTNODES_KEY      = "bootnodes"
)

func main() {
	args := os.Args[1:]
	if args[0] == "-h" {
		println("Instrument takes bootnodes from source genesis file to target.\nJust pass 2 params to with app.\n1 - path to target genesis file, 2 - path to source genesis file")
		return
	}
	if len(args) != 2 {
		panic(fmt.Errorf("you must path 2 params. For help use -h"))
	}
	println("Genesis replacer runned")
	// Read files
	targetGenesisFilePath := args[0]
	targetGenesisFile := readGenesisFile(targetGenesisFilePath)
	sourceGenesisFilePath := args[1]
	sourceGenesisFile := readGenesisFile(sourceGenesisFilePath)

	// Replace content
	bootnodeReplacedGenesisFile := replaceBootNodes(targetGenesisFile, sourceGenesisFile)
	resultGenesisFile := replaceUserData(bootnodeReplacedGenesisFile, sourceGenesisFile)

	// Write result
	writeGenesisFile(resultGenesisFile, targetGenesisFilePath)
}

func replaceBootNodes(targetGenesisFile map[string]interface{}, sourceGenesisFile map[string]interface{}) map[string]interface{} {
	resultGenesisFile := targetGenesisFile

	sourceBootnodes := sourceGenesisFile[BOOTNODES_KEY]
	resultGenesisFile[BOOTNODES_KEY] = sourceBootnodes

	return resultGenesisFile
}

func replaceUserData(targetGenesisFile map[string]interface{}, sourceGenesisFile map[string]interface{}) map[string]interface{} {
	resultGenesisFile := targetGenesisFile

	sourceGenesis := sourceGenesisFile["genesis"].(map[string]interface{})
	resultGenesis := resultGenesisFile["genesis"].(map[string]interface{})

	resultGenesis["extraData"] = sourceGenesis["extraData"]
	resultGenesisFile["genesis"] = resultGenesis

	return resultGenesisFile
}

func readGenesisFile(path string) map[string]interface{} {
	resultPath := path
	if !strings.HasPrefix(resultPath, ROOT_DIR_PREFIX) {
		resultPath = strings.TrimPrefix(resultPath, CURRENT_DIR_PREFIX)
		workDirectory, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		resultPath = fmt.Sprintf("%s/%s", workDirectory, resultPath)
	}
	bytes, readingError := os.ReadFile(resultPath)
	if readingError != nil {
		panic(readingError)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(fmt.Errorf("error in file %s", path))
	}
	return data
}

func writeGenesisFile(genesisFile map[string]interface{}, path string) {
	bytes, marshallError := json.Marshal(genesisFile)
	if marshallError != nil {
		panic(marshallError)
	}

	diskWriteError := os.WriteFile(path, bytes, 0777)
	if diskWriteError != nil {
		panic(diskWriteError)
	}
}
