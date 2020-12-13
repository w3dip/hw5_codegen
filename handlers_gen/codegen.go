// go build handlers_gen/* && ./codegen api.go  api_handler.go
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

type tplFuncParams struct {
	Receiver          string
	HandlerMethodName string
	MethodName        string
	InputParamName    string
	Api
}

type serveHttpParams struct {
	Api
	HandlerMethodName string
}

type Api struct {
	Url    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

var (
	funcTpl = template.Must(template.New("funcTpl").Parse(`
func (srv *{{.Receiver}}) {{.HandlerMethodName}}(w http.ResponseWriter, r *http.Request) {
	{{if .Auth }}
	if r.Header.Get("X-Auth") != "100500" {
		makeOutput(w, ApiResponse{
			Error: "unauthorized",
		}, http.StatusForbidden)
		return
	}
	{{end}}

	{{if .Method }}
	if r.Method != "{{.Method}}" {
		makeOutput(w, ApiResponse{
			Error: "bad method",
		}, http.StatusNotAcceptable)
		return
	}
	{{end}}
	
	// заполнение структуры params
	params := {{.InputParamName}}{
		Login: r.FormValue("login"),
	}
	// валидирование параметров
	ctx := r.Context()
	var res interface{}
	res, err := srv.{{.MethodName}}(ctx, params)
	// прочие обработки
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err.(type) {
		case ApiError:
			err := err.(ApiError)
			makeOutput(w, ApiResponse{
				Error: err.Err.Error(),
			}, err.HTTPStatus)
		default:
			makeOutput(w, ApiResponse{
				Error: err.Error(),
			}, http.StatusInternalServerError)
		}
		return
	}
	makeOutput(w, ApiResponse{
		Response: &res,
	}, http.StatusOK)
}
`))

	apiResponseTpl = template.Must(template.New("serveHttpTpl").Parse(`
type ApiResponse struct {
	Error    string       ` + "`json:\"error\"`" + `
	Response *interface{} ` + "`json:\"response,omitempty\"`" + `
}

func makeOutput(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	result, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	_, err_write := io.WriteString(w, string(result))
	if err_write != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
`))

	serveHttpTpl = template.Must(template.New("serveHttpTpl").Parse(`
func (srv *{{.Receiver}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	{{range .Apis}}
	case "{{.Url}}":
		srv.{{.HandlerMethodName}}(w, r)
	{{end}}
	default:
		makeOutput(w, ApiResponse{
			Error: "unknown method",
		}, http.StatusNotFound)
	}
}
`))
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out, `import "io"`)
	fmt.Fprintln(out, `import "fmt"`)
	fmt.Fprintln(out) // empty line

	//apis := []serveHttpParams{}
	apisByReceiver := make(map[string][]serveHttpParams)
	for _, f := range node.Decls {
		func_decl, ok := f.(*ast.FuncDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.FuncDecl\n", f)
			continue
		}
		name := func_decl.Name.Name

		if func_decl.Doc == nil {
			fmt.Printf("SKIP func %s doesnt have comments\n", name)
			continue
		}

		needCodegen := false
		var api Api
		for _, comment := range func_decl.Doc.List {
			text := comment.Text
			prefix := "// apigen:api "
			hasPrefix := strings.HasPrefix(text, prefix)
			if hasPrefix {
				jsonStr := strings.ReplaceAll(text, prefix, "")
				//api = new(Api)
				err := json.Unmarshal([]byte(jsonStr), &api)
				if err != nil {
					log.Fatalln("Can't parse method comment")
					return
				}
				//if api.Method == "" {
				//	api.Method = http.MethodGet
				//}
				needCodegen = true
				break
			}
		}

		if !needCodegen {
			fmt.Printf("SKIP func %s doesnt have apigen:api mark\n", name)
			continue
		}

		fmt.Printf("Processing func name %s\n", name)
		fmt.Printf("\tgenerating handle method\n")

		var receiver_name string
		for _, field := range func_decl.Recv.List {
			receiver_type := field.Type.(*ast.StarExpr)
			receiver_name = receiver_type.X.(*ast.Ident).Name
			//switch field.(type) {
			//case ast.StarExpr:
			//
			//	break
			//}
		}

		inputParamName := func_decl.Type.Params.List[1].Type.(*ast.Ident).Name

		handlerMethodName := "handle" + name

		err := funcTpl.Execute(out, tplFuncParams{
			Receiver:          receiver_name,
			HandlerMethodName: handlerMethodName,
			MethodName:        name,
			InputParamName:    inputParamName,
			Api:               api,
		})

		if err != nil {
			log.Fatalln("Can't generate api handler for func ", name)
		}

		apis, ok := apisByReceiver[receiver_name]
		if !ok {
			apis = []serveHttpParams{}
		}
		apis = append(apis, serveHttpParams{
			Api:               api,
			HandlerMethodName: handlerMethodName,
		})
		apisByReceiver[receiver_name] = apis
		//SPECS_LOOP:
		//	for _, spec := range func_decl.Specs {
		//		currType, ok := spec.(*ast.TypeSpec)
		//		if !ok {
		//			fmt.Printf("SKIP %T is not ast.TypeSpec\n", spec)
		//			continue
		//		}
		//
		//		currStruct, ok := currType.Type.(*ast.StructType)
		//		if !ok {
		//			fmt.Printf("SKIP %T is not ast.StructType\n", currStruct)
		//			continue
		//		}
		//
		//		if g.Doc == nil {
		//			fmt.Printf("SKIP struct %#v doesnt have comments\n", currType.Name.Name)
		//			continue
		//		}
		//
		//		needCodegen := false
		//		for _, comment := range g.Doc.List {
		//			needCodegen = needCodegen || strings.HasPrefix(comment.Text, "// cgen: binpack")
		//		}
		//		if !needCodegen {
		//			fmt.Printf("SKIP struct %#v doesnt have cgen mark\n", currType.Name.Name)
		//			continue SPECS_LOOP
		//		}
		//
		//		fmt.Printf("process struct %s\n", currType.Name.Name)
		//		fmt.Printf("\tgenerating Unpack method\n")
		//
		//		fmt.Fprintln(out, "func (in *"+currType.Name.Name+") Unpack(data []byte) error {")
		//		fmt.Fprintln(out, "	r := bytes.NewReader(data)")
		//
		//	FIELDS_LOOP:
		//		for _, field := range currStruct.Fields.List {
		//
		//			if field.Tag != nil {
		//				tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		//				if tag.Get("cgen") == "-" {
		//					continue FIELDS_LOOP
		//				}
		//			}
		//
		//			fieldName := field.Names[0].Name
		//			fileType := field.Type.(*ast.Ident).Name
		//
		//			fmt.Printf("\tgenerating code for field %s.%s\n", currType.Name.Name, fieldName)
		//
		//			switch fileType {
		//			case "int":
		//				intTpl.Execute(out, tpl{fieldName})
		//			case "string":
		//				strTpl.Execute(out, tpl{fieldName})
		//			default:
		//				log.Fatalln("unsupported", fileType)
		//			}
		//		}
		//
		//		fmt.Fprintln(out, "	return nil")
		//		fmt.Fprintln(out, "}") // end of Unpack func
		//		fmt.Fprintln(out)      // empty line
		//
		//	}
	}

	for receiver := range apisByReceiver {

		fmt.Printf("Generating ServeHTTP func for receiver %s\n", receiver)

		err = serveHttpTpl.Execute(out,
			struct {
				Receiver string
				Apis     []serveHttpParams
			}{
				Receiver: receiver,
				Apis:     apisByReceiver[receiver],
			})

		if err != nil {
			log.Fatalln("Can't generate serve http")
		}
	}

	fmt.Printf("Generating ApiResponse struct\n")

	err = apiResponseTpl.Execute(out, nil)
	if err != nil {
		log.Fatalln("Can't generate ApiResponse")
	}
}
