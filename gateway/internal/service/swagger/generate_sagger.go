package swagger

import (
	"io"
	"log"
	"modular_chassis/gateway/internal/service/swagger/internal/gen"
	"strings"
)

func GenerateSwagger() (string, error) {
	return gen.New().BuildJson(&gen.Config{
		SearchDir:           "services",
		Excludes:            "",
		ParseExtension:      "",
		MainAPIFile:         "base.go",
		PropNamingStrategy:  "camelcase",
		OutputDir:           "./cmd/docs",
		OutputTypes:         strings.Split("go,json,yaml", ","),
		ParseVendor:         false,
		ParseDependency:     0,
		MarkdownFilesDir:    "",
		ParseInternal:       false,
		GeneratedTime:       false,
		RequiredByDefault:   false,
		CodeExampleFilesDir: "",
		ParseDepth:          100,
		InstanceName:        "",
		OverridesFile:       ".swaggo",
		ParseGoList:         true,
		Tags:                "",
		LeftTemplateDelim:   "{{",
		RightTemplateDelim:  "}}",
		PackageName:         "",
		Debugger:            log.New(io.Discard, "", log.LstdFlags),
		CollectionFormat:    "csv",
		PackagePrefix:       "",
		State:               "",
		ParseFuncBody:       false,
	})
}
