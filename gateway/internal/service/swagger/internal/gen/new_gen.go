package gen

import (
	"fmt"
	"log"
	"modular_chassis/echo/pkg"
	swag "modular_chassis/gateway/internal/service/swagger/internal"
	"os"
	"strings"
)

func (g *Gen) BuildJson(config *Config) (string, error) {
	log.Println("Swagger generator definitions check ...")

	if config.Debugger != nil {
		g.debug = config.Debugger
	}
	if config.InstanceName == "" {
		config.InstanceName = swag.Name
	}

	searchDirs := strings.Split(config.SearchDir, ",")

	if config.LeftTemplateDelim == "" {
		config.LeftTemplateDelim = "{{"
	}

	if config.RightTemplateDelim == "" {
		config.RightTemplateDelim = "}}"
	}

	var overrides map[string]string

	if config.OverridesFile != "" {
		overridesFile, err := open(config.OverridesFile)
		if err != nil {
			// Don't bother reporting if the default file is missing; assume there are no overrides
			if !(config.OverridesFile == DefaultOverridesFile && os.IsNotExist(err)) {
				return "", fmt.Errorf("could not open overrides file: %w", err)
			}
		} else {
			g.debug.Printf("Using overrides from %s", config.OverridesFile)

			overrides, err = parseOverrides(overridesFile)
			if err != nil {
				return "", err
			}
		}
	}

	g.debug.Printf("Generate swagger docs....")

	p := swag.New(
		swag.SetParseDependency(config.ParseDependency),
		swag.SetMarkdownFileDirectory(config.MarkdownFilesDir),
		swag.SetDebugger(config.Debugger),
		swag.SetExcludedDirsAndFiles(config.Excludes),
		swag.SetParseExtension(config.ParseExtension),
		swag.SetCodeExamplesDirectory(config.CodeExampleFilesDir),
		swag.SetStrict(config.Strict),
		swag.SetOverrides(overrides),
		swag.ParseUsingGoList(config.ParseGoList),
		swag.SetTags(config.Tags),
		swag.SetCollectionFormat(config.CollectionFormat),
		swag.SetPackagePrefix(config.PackagePrefix),
	)

	p.PropNamingStrategy = config.PropNamingStrategy
	p.ParseVendor = config.ParseVendor
	p.ParseInternal = config.ParseInternal
	p.RequiredByDefault = config.RequiredByDefault
	p.HostState = config.State
	p.ParseFuncBody = config.ParseFuncBody

	if err := p.ParseAPIMultiSearchDirV2(pkg.EmbeddedFiles, searchDirs, config.MainAPIFile, "definition.go"); err != nil {
		return "", err
	}

	b, err := g.jsonIndent(p.GetSwagger())
	if err != nil {
		return "", err
	}

	log.Println("Done")
	return string(b), nil
}
