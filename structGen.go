package sbcgostructgen

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	sberr "github.com/barbell-math/smoothbrain-errs"
)

//go:generate go-enum --marshal --names --values

type (
	Opts struct {
		ExitOnErr bool
		DryRun    bool
	}

	structField struct {
		Type fieldType
		Name string
	}

	// ENUM (
	//  void,
	// )
	fieldType int
)

var (
	InvalidTypeErr        = errors.New("Invalid Type")
	UnderspecifiedTypeErr = errors.New("Underspecified type")
)

func GenerateFor[T any](file string, opts *Opts) error {
	var err error
	refType := reflect.TypeFor[T]()
	cStructs := map[string][]structField{}
	includes := map[string]struct{}{}

	if refType.Kind() != reflect.Struct {
		err = sberr.Wrap(
			InvalidTypeErr, "Expected struct, got %s", refType.Kind(),
		)
		goto errExit
	}

	if _, err = os.Stat(file); err != nil && !opts.DryRun {
		goto errExit
	}

	if err = checkType(refType, ""); err != nil {
		goto errExit
	}

	if err = generateCStructs(
		refType, "", "", cStructs, includes,
	); err != nil {
		goto errExit
	}

errExit:
	if err != nil {
		fmt.Println("ERROR:", err)
		if opts.ExitOnErr {
			os.Exit(1)
		}
	}
	return err
}

func checkType(refType reflect.Type, fieldName string) error {
	switch refType.Kind() {
	case reflect.Map, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Complex64, reflect.Complex128:
		return sberr.Wrap(
			InvalidTypeErr,
			"Cannot translate a Go %s to C, field %s",
			refType.Kind(), fieldName,
		)
	case reflect.Int, reflect.Uint:
		return sberr.Wrap(
			UnderspecifiedTypeErr,
			"A %s can be varying sizes in C, specify bit size to fix (i.e. int32 instead of int), field %s",
			refType.Kind(), fieldName,
		)
	case reflect.UnsafePointer, reflect.Uintptr:
		fmt.Printf(
			"Cannot validate kind of %s, will be specified in C as a void*, field %s\n",
			refType.Kind(), fieldName,
		)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return checkType(refType.Elem(), fieldName)
	case reflect.Struct:
		for i := range refType.NumField() {
			iterField := refType.Field(i)

			var iterFieldName string
			if fieldName == "" {
				iterFieldName = iterField.Name
			} else {
				iterFieldName += fieldName + "." + iterField.Name
			}

			return checkType(iterField.Type, iterFieldName)
		}
	}
	return nil
}

func generateCStructs(
	refType reflect.Type,
	structName string,
	fieldName string,
	cStructs map[string][]structField,
	includes map[string]struct{},
) error {
	switch refType.Kind() {
	case reflect.UnsafePointer, reflect.Uintptr:
		cStructs[structName] = append(
			cStructs[structName],
			structField{Type: "void*", Name: fieldName},
		)
	case reflect.Int8:
		cStructs[structName] = append(
			cStructs[structName],
			fmt.Sprintf("int8_t %s", fieldName),
		)
		includes["<stdint.h>"] = struct{}{}
	case reflect.Int16:
		cStructs[structName] = append(
			cStructs[structName],
			fmt.Sprintf("int16_t %s", fieldName),
		)
		includes["<stdint.h>"] = struct{}{}
	case reflect.Int32:
		cStructs[structName] = append(
			cStructs[structName],
			fmt.Sprintf("int32_t %s", fieldName),
		)
		includes["<stdint.h>"] = struct{}{}
	case reflect.Int64:
		cStructs[structName] = append(
			cStructs[structName],
			fmt.Sprintf("int64_t %s", fieldName),
		)
		includes["<stdint.h>"] = struct{}{}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return generateCStructs(refType.Elem(), fieldName)
	case reflect.Struct:
		newStructName := refType.Name()
		// TODO - what to do about annon structs??
		cStructs[newStructName] = make([]string, 0)
		if fieldName != "" {
			cStructs[structName] = append(
				cStructs[structName],
				fmt.Sprintf("%s %s", newStructName, fieldName),
			)
		}

		for i := range refType.NumField() {
			iterField := refType.Field(i)
			return generateCStructs(
				iterField.Type, newStructName, iterField.Name, cStructs,
			)
		}
	}
	return nil
}
