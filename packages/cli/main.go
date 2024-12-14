package cli

import (
	"errors"
	"os"
)

type translatedArgs struct {
	SpecificationPath string
}

func GetArgs() (map[string]string, error) {
	args := os.Args[1:]
	formatedArgs := make(map[string]string)

	for i := 0; i < len(args); i++ {
		key := args[i]

		if key[0] != '-' {
			message := "invalid argment: " + key
			return nil, errors.New(message)
		}

		formatedArgs[key] = args[i+1]
		i++
	}

	return formatedArgs, nil
}

func TranslateArgs(args map[string]string) translatedArgs {
	return translatedArgs{
		SpecificationPath: args["-p"],
	}
}
