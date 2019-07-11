package codegen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

func NewGenFile(pkgName string, filename string) *GenFile {
	return &GenFile{
		PkgName:  pkgName,
		FileName: filename,
		Importer: &Importer{},
		Buffer:   &bytes.Buffer{},
	}
}

type GenFile struct {
	FileName string
	PkgName  string
	Data     interface{}
	*bytes.Buffer
	*Importer
}

func (f *GenFile) WithData(data interface{}) *GenFile {
	f.Data = data
	return f
}

func (f *GenFile) Block(tpl string) *GenFile {
	f.writeTo(f.Buffer, tpl)
	return f
}

func (f *GenFile) writeTo(writer io.Writer, tpl string) {
	t, parseErr := template.New(f.FileName).Parse(tpl)
	if parseErr != nil {
		panic(fmt.Sprintf("template Prase failded: %s", parseErr.Error()))
	}
	err := t.Execute(writer, f)
	if err != nil {
		panic(fmt.Sprintf("template Execute failded: %s", err.Error()))
	}
}

func (f *GenFile) String() string {
	if strings.HasSuffix(f.FileName, "go") {
		return fmt.Sprintf(`
package %s
%s
%s
`,
			f.PkgName,
			f.Importer.String(),
			f.Buffer.String(),
		)
	} else {
		return fmt.Sprintf(`%s`, f.Buffer.String())
	}

}

func (f *GenFile) OutputTo(outputs Outputs) {
	outputs.Add(f.FileName, f.String())
}
