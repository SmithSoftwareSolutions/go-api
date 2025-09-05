package main

import (
	"fmt"
	"go/types"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	if len(os.Args) != 3 {
		failErr(fmt.Errorf("expected exactly two argumenst: <source type>"))
	}
	sourceType := os.Args[1]

	pkg := loadPackage("")

	obj := pkg.Types.Scope().Lookup(sourceType)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			sourceType, pkg))
	}

	if _, ok := obj.(*types.TypeName); !ok {
		failErr(fmt.Errorf("%v is not a named type", obj))
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		failErr(fmt.Errorf("type %v is not a struct", obj))
	}

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tagValue := structType.Tag(i)
		fmt.Println(field.Name(), tagValue, field.Type())
	}
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
