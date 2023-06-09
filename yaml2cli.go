//
// Copyright © 2023 Krishna Miriyala <krishnambm@gmail.com>
//

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func parseArgs() []string {
	var inputFiles []string
	pflag.StringSliceVarP(&inputFiles, "input", "i", []string{}, "Yaml configuration files in order of overrides")
	pflag.Parse()
	return inputFiles
}

func main() {
	fmt.Println(Yaml2CliParams(parseArgs()...)) //nolint
}

func Yaml2CliParams(inputFiles ...string) string {
	params := make(map[string]interface{})
	cliParams := make([]string, 0)
	for _, inputFile := range inputFiles {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(data, &params)
		if err != nil {
			log.Fatal(err)
		}
	}

	for key, value := range params {
		cliParam := ""
		if len(key) == 1 {
			cliParam = fmt.Sprintf("-%s", key)
		} else {
			cliParam = fmt.Sprintf("--%s", key)
		}

		switch v := value.(type) {
		case []interface{}:
			var values []string
			for _, val := range v {
				values = append(values, fmt.Sprintf("\"%v\"", val))
			}
			cliParam += " " + strings.Join(values, " ")
		default:
			cliParam += fmt.Sprintf(" %v", v)
		}

		cliParams = append(cliParams, cliParam)
	}

	return strings.Join(cliParams, " ")
}
