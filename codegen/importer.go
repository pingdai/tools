package codegen

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"strconv"
	"strings"
)

type ImportPkg struct {
	*build.Package
	Alias string
}

func (importPkg *ImportPkg) GetID() string {
	if importPkg.Alias != "" {
		return importPkg.Alias
	}
	return importPkg.Name
}

func (importPkg *ImportPkg) String() string {
	if importPkg.Alias != "" {
		return importPkg.Alias + " " + strconv.Quote(importPkg.ImportPath) + "\n"
	}
	return strconv.Quote(importPkg.ImportPath)
}

type Importer struct {
	Local string
	pkgs  map[string]*ImportPkg
}

func getPkgImportPathAndExpose(s string) (pkgImportPath string, expose string) {
	idxSlash := strings.LastIndex(s, "/")
	idxDot := strings.LastIndex(s, ".")
	if idxDot > idxSlash {
		return s[0:idxDot], s[idxDot+1:]
	}
	return s, ""
}

func (importer *Importer) String() string {
	buf := &bytes.Buffer{}
	if len(importer.pkgs) > 0 {
		buf.WriteString("import (\n")
		importer.WriteToImports(buf)
		buf.WriteString(")")
	}
	return buf.String()
}

func (importer *Importer) WriteToImports(w io.Writer) {
	if len(importer.pkgs) > 0 {
		for _, importPkg := range importer.pkgs {
			io.WriteString(w, importPkg.String()+"\n")
		}
	}
}

func (importer *Importer) Import(importPath string, alias bool) *ImportPkg {
	importPath = DeVendor(importPath)
	if importer.pkgs == nil {
		importer.pkgs = map[string]*ImportPkg{}
	}

	importPkg, exists := importer.pkgs[importPath]
	if !exists {
		pkg, err := build.Import(importPath, "", build.ImportComment)
		if err != nil {
			panic(err)
		}
		importPkg = &ImportPkg{
			Package: pkg,
		}
		if alias {
			importPkg.Alias = ToLowerSnakeCase(importPath)
		}
		importer.pkgs[importPath] = importPkg
	}

	return importPkg
}

func DeVendor(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func (importer *Importer) PureUse(importPath string, subPkgs ...string) string {
	pkgPath, expose := getPkgImportPathAndExpose(strings.Join(append([]string{importPath}, subPkgs...), "/"))

	importPkg := importer.Import(pkgPath, false)

	if expose != "" {
		if pkgPath == importer.Local {
			return expose
		}
		return fmt.Sprintf("%s.%s", importPkg.GetID(), expose)
	}

	return importPkg.GetID()
}
