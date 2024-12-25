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

func (parser *Parser) ParseEmbeddedFiles(embedded embed.FS, dirName string, mainAPIFile string) error {
	dir, err := embedded.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.IsDir() {
			err := parser.ParseEmbeddedFiles(embedded, filepath.Join(dirName, file.Name()), mainAPIFile)
			if err != nil {
				return err
			}
		} else {
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
func (parser *Parser) ParseAPIMultiSearchDirV2(embedded embed.FS, embedRoots []string, mainAPIFile string) error {
	for _, er := range embedRoots {
		err := parser.ParseEmbeddedFiles(embedded, er, mainAPIFile)
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

	if filePath := strings.Split(fileInfo.Path, "/"); filePath[len(filePath)-1] != "definition.go" {
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
	method := toSnakeCase(name)

	param := "// @Param %s body %s.%s true \"Request body\""
	result := "// @Success 200 {object} %s.%s"
	router := "// @Router /api/%s/%s [POST]"

	if len(decl.Params.List) > 1 {
		t := decl.Params.List[1].Type.(*ast.IndexExpr).Index.(*ast.Ident).Name
		param = fmt.Sprintf(param, t, packageName, t)
	}

	if len(decl.Results.List) != 0 {
		t := decl.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.Ident).Name
		result = fmt.Sprintf(result, packageName, t)
	}

	router = fmt.Sprintf(router, packageName, method)

	return &ast.CommentGroup{List: []*ast.Comment{{Text: param}, {Text: result}, {Text: router}}}
}
