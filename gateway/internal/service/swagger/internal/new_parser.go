package swag

import (
	"embed"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func (parser *Parser) ParseEmbeddedFiles(embedded embed.FS, dirName string, mainAPIFile string, targetedFilesName string) error {
	dir, err := embedded.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.IsDir() {
			err := parser.ParseEmbeddedFiles(embedded, filepath.Join(dirName, file.Name()), mainAPIFile, targetedFilesName)
			if err != nil {
				return err
			}
		} else {
			if file.Name() != targetedFilesName && file.Name() != mainAPIFile {
				continue
			}
			data, err := embedded.ReadFile(filepath.Join(dirName, file.Name()))
			if err != nil {
				return err
			}
			if file.Name() != mainAPIFile {
				err = parser.packages.ParseFile(dirName, filepath.Join(dirName, file.Name()), data, ParseAll)
				if err != nil {
					return err
				}
			} else {
				err := parser.ParseGeneralAPIInfoV2(data, mainAPIFile)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ParseAPIMultiSearchDirV2 is like ParseAPI but for multiple search dirs.
func (parser *Parser) ParseAPIMultiSearchDirV2(embedded embed.FS, embedRoots []string, mainAPIFile string, targetedFileName string) error {
	for _, er := range embedRoots {
		err := parser.ParseEmbeddedFiles(embedded, er, mainAPIFile, targetedFileName)
		if err != nil {
			return err
		}
	}

	var err error
	parser.parsedSchemas, err = parser.packages.ParseTypes()
	if err != nil {
		return err
	}

	err = parser.packages.RangeFiles(parser.ParseRouterAPIInfoV2)
	if err != nil {
		return err
	}

	return parser.checkOperationIDUniqueness()
}

func (parser *Parser) ParseGeneralAPIInfoV2(data []byte, mainAPIFile string) error {
	fileTree, err := goparser.ParseFile(token.NewFileSet(), "", data, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse source files %s: %s", mainAPIFile, err)
	}

	parser.swagger.Swagger = "2.0"

	for _, comment := range fileTree.Comments {
		comments := strings.Split(comment.Text(), "\n")
		if !isGeneralAPIComment(comments) {
			continue
		}

		err = parseGeneralAPIInfo(parser, comments)
		if err != nil {
			return err
		}
	}

	return nil
}

func (parser *Parser) ParseRouterAPIInfoV2(fileInfo *AstFileInfo) error {
	if (fileInfo.ParseFlag & ParseOperations) == ParseNone {
		return nil
	}

	// parse File.Comments instead of File.Decls.Doc if ParseFuncBody flag set to "true"
	if parser.ParseFuncBody {
		for _, astComments := range fileInfo.File.Comments {
			if astComments.List != nil {
				if err := parser.parseRouterAPIInfoComment(astComments.List, fileInfo); err != nil {
					return err
				}
			}
		}

		return nil
	}

	for _, decl := range fileInfo.File.Decls {
		if s, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range s.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if ift, ok := ts.Type.(*ast.InterfaceType); ok {
						for _, method := range ift.Methods.List {
							funcDoc := parser.CreateCommentsBasedOnFuncDecl(method.Names[0].Name, method.Type.(*ast.FuncType), fileInfo)
							if funcDoc != nil && funcDoc.List != nil {
								if err := parser.parseRouterAPIInfoComment(funcDoc.List, fileInfo); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (parser *Parser) CreateCommentsBasedOnFuncDecl(name string, decl *ast.FuncType, file *AstFileInfo) *ast.CommentGroup {
	packageName := toSnakeCase(file.File.Name.Name)
	methodName := toSnakeCase(name)

	param := "// @Param %s body %s.%s true \"Request body\""
	result := "// @Success 200 {object} %s.%s"
	router := "// @Router /api/%s/%s [POST]"

	defer func() {
		if err := recover(); err != nil {
			PrintError(packageName, methodName)
		}
	}()

	if len(decl.Params.List) == 2 {
		t := decl.Params.List[1].Type.(*ast.Ident).Name
		for i, field := range decl.Params.List[1].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
			if field.Names == nil {
				decl.Params.List[1].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).
					Fields.List[i].Tag = &ast.BasicLit{Value: `swaggerignore:"true"`}
			}
		}
		param = fmt.Sprintf(param, t, packageName, t)
	} else {
		PrintError(packageName, methodName)
		return nil
	}

	if len(decl.Results.List) == 2 {
		t := decl.Results.List[0].Type.(*ast.Ident).Name
		for i, field := range decl.Results.List[0].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
			if field.Names == nil {
				decl.Results.List[0].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).
					Fields.List[i].Tag = &ast.BasicLit{Value: `swaggerignore:"true"`}
			}
		}
		result = fmt.Sprintf(result, packageName, t)
	} else {
		PrintError(packageName, methodName)
		return nil
	}

	router = fmt.Sprintf(router, packageName, methodName)

	return &ast.CommentGroup{List: []*ast.Comment{{Text: param}, {Text: result}, {Text: router}}}
}

func PrintError(packageName string, methodName string) {
	fmt.Printf("\033[33mFunc %s.%s does not follow below pattern to be exposed:\n"+
		"Func(Context,{Request})({Response},error)\033[0m\n", packageName, methodName)
}
