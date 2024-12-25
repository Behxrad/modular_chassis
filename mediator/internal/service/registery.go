package service

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"modular_chassis/echo/pkg"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	once        sync.Once
	registryIns *registry
)

type registry struct {
	services     map[string]interface{}
	serviceFuncs map[string]interface{}

	servicesInterfaceMethods  map[string][]string
	servicesMethodsInterfaces map[string][]string
}

func GetRegistry() *registry {
	once.Do(func() {
		if registryIns == nil {
			registryIns = &registry{
				services:     make(map[string]interface{}),
				serviceFuncs: make(map[string]interface{}),

				servicesInterfaceMethods:  make(map[string][]string),
				servicesMethodsInterfaces: make(map[string][]string),
			}
		}
	})
	return registryIns
}

func (r *registry) RegisterService(serviceImpl interface{}) {
	r.identifyServiceDefinitions("services")
	c := r.hasImplementedAnyDefinition(serviceImpl)
	fmt.Println(c)
	r.services["test"] = serviceImpl
}

func (r *registry) hasImplementedAnyDefinition(serviceImpl interface{}) string {
	defs := r.extractImplFuncDefs(serviceImpl)
	interfaces := r.servicesMethodsInterfaces[defs[1]]
	candidate := ""
	for _, interfaze := range interfaces {
		if len(r.servicesInterfaceMethods[interfaze]) == len(defs) {
			found := 0
			for _, method := range r.servicesInterfaceMethods[interfaze] {
				for _, def := range defs {
					if method == def {
						found++
						if found == len(defs) {
							return interfaze
						}
						break
					}
				}
			}
		}
	}
	return candidate
}

func (r *registry) extractImplFuncDefs(serviceImpl interface{}) []string {
	funcDefs := make([]string, 0)
	t := reflect.TypeOf(serviceImpl)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		paramCompile := regexp.MustCompile("(^[A-Za-z0-9]+$)|(^[A-Za-z0-9]+\\[.*]$)")
		genericArgCompile := regexp.MustCompile("\\w+\\.\\w+")
		replaceGenericArgCompile := regexp.MustCompile("\\[.*]")
		var params, results []string
		for i := 1; i < t.Method(0).Func.Type().NumIn(); i++ {
			find := paramCompile.Find([]byte(method.Func.Type().In(i).Name()))
			if find2 := genericArgCompile.Find([]byte(method.Func.Type().In(i).Name())); find2 != nil {
				find = replaceGenericArgCompile.ReplaceAll(find, []byte(fmt.Sprintf("[%s]", find2)))
			}
			params = append(params, string(find))
		}

		for i := 0; i < t.Method(0).Func.Type().NumOut(); i++ {
			find := paramCompile.Find([]byte(method.Func.Type().Out(i).Name()))
			if find2 := genericArgCompile.Find([]byte(method.Func.Type().Out(i).Name())); find2 != nil {
				find = replaceGenericArgCompile.ReplaceAll(find, []byte(fmt.Sprintf("[%s]", find2)))
			}
			results = append(results, string(find))
		}
		funcDefs = append(funcDefs, fmt.Sprintf("%s(%s)(%s)", method.Name, strings.Join(params, ","), strings.Join(results, ",")))
	}
	return funcDefs
}

func (r *registry) identifyServiceDefinitions(dirName string) error {
	dir, err := pkg.EmbeddedFiles.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.IsDir() {
			err := r.identifyServiceDefinitions(filepath.Join(dirName, file.Name()))
			if err != nil {
				return err
			}
		} else {
			data, err := pkg.EmbeddedFiles.ReadFile(filepath.Join(dirName, file.Name()))
			if err != nil {
				return err
			}
			fileSet := token.NewFileSet()
			parseFile, err := goparser.ParseFile(fileSet, filepath.Join(dirName, file.Name()), data, goparser.ParseComments)
			if err != nil {
				return err
			}
			r.identifyServiceInterfaces(parseFile)
		}
	}
	return nil
}

func (r *registry) identifyServiceInterfaces(file *ast.File) {
	for _, decl := range file.Decls {
		if s, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range s.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if ift, ok := ts.Type.(*ast.InterfaceType); ok {
						for _, method := range ift.Methods.List {
							decl := method.Type.(*ast.FuncType)

							packageName := file.Name.Name
							methodName := method.Names[0].Name

							var param, result string

							if len(decl.Params.List) > 1 {
								reqt := decl.Params.List[1].Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel.Name
								reqtp := decl.Params.List[1].Type.(*ast.IndexExpr).Index.(*ast.Ident).Name
								param = fmt.Sprintf("%s[%s.%s]", reqt, packageName, reqtp)
							}

							if len(decl.Results.List) != 0 {
								rest := decl.Results.List[0].Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel.Name
								restp := decl.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.Ident).Name
								result = fmt.Sprintf("%s[%s.%s]", rest, packageName, restp)
							}

							funcDef := fmt.Sprintf("%s(%s,%s)(%s,%s)", methodName, "Context", param, result, "error")
							packageInterface := fmt.Sprintf("%s.%s", packageName, ts.Name.Name)
							r.servicesInterfaceMethods[packageInterface] = append(r.servicesInterfaceMethods[packageInterface], funcDef)
							r.servicesMethodsInterfaces[funcDef] = append(r.servicesMethodsInterfaces[funcDef], packageInterface)
						}
					}
				}
			}
		}
	}
}
