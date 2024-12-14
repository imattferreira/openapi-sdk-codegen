package main

import (
	"fmt"

	"github.com/imattferreira/openapi-sdk-codegen/src/cli"
	"github.com/imattferreira/openapi-sdk-codegen/src/codegen"
	"github.com/imattferreira/openapi-sdk-codegen/src/openapi"
)

func main() {
	args, err := cli.GetArgs()

	if err != nil {
		fmt.Println(err)
		return
	}

	translatedArgs := cli.TranslateArgs(args)

	if len(translatedArgs.SpecificationPath) == 0 {
		fmt.Println("Specification path not provided")
		return
	}

	specification, err := openapi.ReadFile(translatedArgs.SpecificationPath)

	if err != nil {
		fmt.Println("Error unmarshalling specification file")
		return
	}

	if specification.Version != "3.1.0" {
		fmt.Println("Codegen only supports OpenAPI v3.1.0")
		return
	}

	translator := openapi.Translator{}
	translated := translator.Translate(specification)

	codegen.Generate(translated)
}
