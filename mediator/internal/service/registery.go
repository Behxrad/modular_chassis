package service

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"modular_chassis/echo/pkg"
	"modular_chassis/echo/pkg/errs"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/echo/pkg/utils/utils"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	servicesDefinitionRoot = "services"
)

var (
	once        sync.Once
	registryIns *registry
)

type methodInfo struct {
	Function reflect.Value
	Request  reflect.Type
	Response reflect.Type
}

type registry struct {
	serviceMethods map[string]methodInfo

	servicesInterfaceMethods  map[string][]string
	servicesMethodsInterfaces map[string][]string
}

func init() {
	fmt.Println("Registry definitions check ...")
	err := GetRegistry().identifyServiceDefinitions(servicesDefinitionRoot)
	if err != nil {
		return
	}
}

func GetRegistry() *registry {
	once.Do(func() {
		if registryIns == nil {
			registryIns = &registry{
				serviceMethods: make(map[string]methodInfo),

				servicesInterfaceMethods:  make(map[string][]string),
				servicesMethodsInterfaces: make(map[string][]string),
			}
		}
	})
	return registryIns
}

func (r *registry) GetService(domain, method string) methodInfo {
	return r.serviceMethods[fmt.Sprintf("%s.%s", domain, method)]
}

func (r *registry) GetServiceRequestModel(domain, method string) (reflect.Type, error) {
	mInfo := r.GetService(domain, method)
	if !mInfo.Function.IsValid() {
		return nil, errs.NewServiceErrorCode(services.ServiceNotFound)
	}
	return mInfo.Request, nil
}

func (r *registry) RegisterService(serviceImpl interface{}) {
	implMethods := r.extractImplFuncDefs(serviceImpl)
	interfaceName := r.identifyImplementedServiceInterface(implMethods)
	if interfaceName == "" {
		return
	}

	afterMethodNameCompile := regexp.MustCompile("\\(.*\\)")
	for signature, value := range implMethods {
		d := strings.Split(interfaceName, ".")[0]
		m := utils.ToSnakeCase(string(afterMethodNameCompile.ReplaceAll([]byte(signature), []byte{})))
		r.serviceMethods[fmt.Sprintf("%s.%s", d, m)] = value
	}
}

func (r *registry) identifyImplementedServiceInterface(implMethods map[string]methodInfo) string {
	var candidate, firstKey string
	for k := range implMethods {
		firstKey = k
		break
	}
	interfaces := r.servicesMethodsInterfaces[firstKey]
	for _, interfaze := range interfaces {
		if len(r.servicesInterfaceMethods[interfaze]) == len(implMethods) {
			found := 0
			for _, interfaceMethod := range r.servicesInterfaceMethods[interfaze] {
				for implMethod := range implMethods {
					if interfaceMethod == implMethod {
						found++
						if found == len(implMethods) {
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

func (r *registry) extractImplFuncDefs(serviceImpl interface{}) map[string]methodInfo {
	funcDefs := make(map[string]methodInfo)
	t := reflect.TypeOf(serviceImpl)
	v := reflect.ValueOf(serviceImpl)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		paramCompile := regexp.MustCompile("(^[A-Za-z0-9]+$)|(^[A-Za-z0-9]+\\[.*]$)")
		genericArgCompile := regexp.MustCompile("\\w+\\.\\w+")
		replaceGenericArgCompile := regexp.MustCompile("\\[.*]")
		var params, results []string
		for i := 1; i < method.Func.Type().NumIn(); i++ {
			find := paramCompile.Find([]byte(method.Func.Type().In(i).Name()))
			if find2 := genericArgCompile.Find([]byte(method.Func.Type().In(i).Name())); find2 != nil {
				find = replaceGenericArgCompile.ReplaceAll(find, []byte(fmt.Sprintf("[%s]", find2)))
			}
			params = append(params, string(find))
		}

		for i := 0; i < method.Func.Type().NumOut(); i++ {
			find := paramCompile.Find([]byte(method.Func.Type().Out(i).Name()))
			if find2 := genericArgCompile.Find([]byte(method.Func.Type().Out(i).Name())); find2 != nil {
				find = replaceGenericArgCompile.ReplaceAll(find, []byte(fmt.Sprintf("[%s]", find2)))
			}
			results = append(results, string(find))
		}
		funcDefs[fmt.Sprintf("%s(%s)(%s)", method.Name, strings.Join(params, ","),
			strings.Join(results, ","))] = methodInfo{
			Function: v.Method(i),
			Request:  t.Method(i).Func.Type().In(2),
			Response: t.Method(i).Func.Type().Out(0),
		}
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
			if file.Name() != "definition.go" {
				continue
			}
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

							defer func() {
								if err := recover(); err != nil {
									fmt.Printf("\033[33mFunc %s.%s does not follow below pattern to be exposed:\n"+
										"Func(Context,{Request})({Response},error)\033[0m\n", packageName, methodName)
								}
							}()

							if len(decl.Params.List) == 2 {
								param = decl.Params.List[1].Type.(*ast.Ident).Name
							}

							if len(decl.Results.List) == 2 {
								result = decl.Results.List[0].Type.(*ast.Ident).Name
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
