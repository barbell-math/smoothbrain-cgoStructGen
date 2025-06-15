// A very simple library that is used to generate struct definitions shared
// between Go and C through CGO.
package sbcgostructgen

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"

	sberr "github.com/barbell-math/smoothbrain-errs"
)

//go:generate go-enum --marshal --names --values

type (
	// ENUM(
	//  void*,
	//	char*,
	//  int8_t,
	//  int16_t,
	//  int32_t,
	//  int64_t,
	//  uint8_t,
	//  uint16_t,
	//  uint32_t,
	//  uint64_t,
	//	float,
	//	double,
	//	bool,
	// )
	fieldType string

	// ENUM(
	//	None,
	//	Pntr,
	//	Array,
	// )
	typeMod int
	include string

	typeModifier struct {
		typeMod
		tModAmnt int
	}

	structField struct {
		typeModifier
		_type string
		name  string
	}

	CGoStructGen struct {
		includes map[include]struct{}
		structs  map[string][]structField
	}
)

var (
	InvalidTypeErr        = errors.New("Invalid Type")
	UnderspecifiedTypeErr = errors.New("Underspecified type")
	AnonymousNameErr      = errors.New("Anonymous name")

	reflectToEnumTypes = map[reflect.Kind]fieldType{
		reflect.Uintptr:       FieldTypeVoid,
		reflect.UnsafePointer: FieldTypeVoid,
		reflect.Int8:          FieldTypeInt8T,
		reflect.Int16:         FieldTypeInt16T,
		reflect.Int32:         FieldTypeInt32T,
		reflect.Int64:         FieldTypeInt64T,
		reflect.Uint8:         FieldTypeUint8T,
		reflect.Uint16:        FieldTypeUint16T,
		reflect.Uint32:        FieldTypeUint32T,
		reflect.Uint64:        FieldTypeUint64T,
		reflect.Float32:       FieldTypeFloat,
		reflect.Float64:       FieldTypeDouble,
		reflect.Bool:          FieldTypeBool,
		reflect.String:        FieldTypeChar,
	}

	reflectToIncludes = map[reflect.Kind]include{
		reflect.Int8:    "<stdint.h>",
		reflect.Int16:   "<stdint.h>",
		reflect.Int32:   "<stdint.h>",
		reflect.Int64:   "<stdint.h>",
		reflect.Uint8:   "<stdint.h>",
		reflect.Uint16:  "<stdint.h>",
		reflect.Uint32:  "<stdint.h>",
		reflect.Uint64:  "<stdint.h>",
		reflect.Float32: "<stdint.h>",
		reflect.Float64: "<stdint.h>",
		reflect.Bool:    "<stdbool.h>",
	}
)

func (s structField) String() string {
	switch s.typeModifier.typeMod {
	case TypeModPntr:
		return fmt.Sprintf("%s* %s", s._type, s.name)
	case TypeModArray:
		return fmt.Sprintf("%s %s[%d]", s._type, s.name, s.tModAmnt)
	case TypeModNone:
		fallthrough
	default:
		return fmt.Sprintf("%s %s", s._type, s.name)
	}
}
func (i include) String() string {
	return "#include " + string(i)
}

// Creates a new struct generator.
func New() *CGoStructGen {
	return &CGoStructGen{
		includes: map[include]struct{}{},
		structs:  map[string][]structField{},
	}
}

// Adds the supplied type and all of its sub-types to the struct generator. The
// type of T must be a struct or an error will be returned. The following field
// types are allowed in the struct:
//   - int8, int16, int32, int64
//   - uint8, uint16, uint32, uint64
//   - float32, float64
//   - string
//   - bool
//   - uintptr, unsafe.Pointer
//   - arrays and structs that are composed of the above types
//
// Types will be recursively added. Types that are duplicated between struct
// definitions will not be duplicated in the output C code.
//
// This funciton is intended to be called many times with the same value for the
// `ts` argument. The `ts` value will be updated with any newly-found structs.
func GenerateFor[T any](ts *CGoStructGen) error {
	var err error
	refType := reflect.TypeFor[T]()

	if refType.Kind() != reflect.Struct {
		err = sberr.Wrap(
			InvalidTypeErr, "Expected struct, got %s", refType.Kind(),
		)
		goto errExit
	}

	if err = checkType(refType, "", ts.structs); err != nil {
		goto errExit
	}
	generateCStructs(
		refType, "",
		"", typeModifier{typeMod: TypeModNone},
		ts.structs, ts.includes,
	)

errExit:
	return err
}

func checkType(
	refType reflect.Type,
	fieldName string,
	cStructs map[string][]structField,
) error {
	switch refType.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface,
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
			"WARN: Cannot validate kind of %s, will be specified in C as a void*, field %s\n",
			refType.Kind(), fieldName,
		)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Array, reflect.Pointer:
		return checkType(refType.Elem(), fieldName, cStructs)
	case reflect.Struct:
		newStructName := refType.Name()
		if newStructName == "" {
			return sberr.Wrap(
				AnonymousNameErr,
				"Anonymous structs are not supported, add a name, field %s\n",
				fieldName,
			)
		}
		if _, ok := cStructs[newStructName]; !ok {
			cStructs[newStructName] = make([]structField, 0)
		}

		for i := range refType.NumField() {
			iterField := refType.Field(i)

			var iterFieldName string
			if fieldName == "" {
				iterFieldName = iterField.Name
			} else {
				iterFieldName += fieldName + "." + iterField.Name
			}

			if err := checkType(
				iterField.Type, iterFieldName, cStructs,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateCStructs(
	refType reflect.Type, structName string,
	fieldName string, tMod typeModifier,
	cStructs map[string][]structField, includes map[include]struct{},
) {
	if e, ok := reflectToEnumTypes[refType.Kind()]; ok {
		cStructs[structName] = append(
			cStructs[structName],
			structField{
				_type:        e.String(),
				name:         fieldName,
				typeModifier: tMod,
			},
		)
		if i, ok := reflectToIncludes[refType.Kind()]; ok {
			includes[i] = struct{}{}
		}
		return
	}

	switch refType.Kind() {
	case reflect.Array:
		generateCStructs(
			refType.Elem(), structName,
			fieldName,
			typeModifier{typeMod: TypeModArray, tModAmnt: refType.Len()},
			cStructs, includes,
		)
	case reflect.Pointer:
		generateCStructs(
			refType.Elem(), structName,
			fieldName, typeModifier{typeMod: TypeModPntr},
			cStructs, includes,
		)
	case reflect.Struct:
		newStructName := refType.Name()
		if structName != "" {
			cStructs[structName] = append(
				cStructs[structName],
				structField{
					_type:        fmt.Sprintf("%s_t", newStructName),
					name:         fieldName,
					typeModifier: tMod,
				},
			)
		}

		if _, ok := cStructs[structName]; ok {
			// If the struct fields were already populated then don't add them
			// again
			return
		}
		for i := range refType.NumField() {
			iterField := refType.Field(i)
			generateCStructs(
				iterField.Type, newStructName,
				iterField.Name, typeModifier{typeMod: TypeModNone},
				cStructs, includes,
			)
		}
	default:
		// All errors should be caught by the [checkType] function
		panic(fmt.Sprintf(
			"An unexpected type was recieved: %s", refType.Name(),
		))
	}
}

// Writes all of the struct definitions that were previously added through calls
// to [GenerateFor] to the specified file.
func (t *CGoStructGen) WriteTo(file string, headerStr string) error {
	var err error
	var f *os.File

	f, err = os.Create(file)
	defer f.Close()
	if err != nil {
		goto errExit
	}

	t.templateHeader(f, headerStr)
	t.templateIncludes(f)
	t.templateExternCIf(f, func() {
		t.templateCStructs(f)
	})
	t.templateFooter(f)

errExit:
	return err
}

func (t *CGoStructGen) templateHeader(f *os.File, headerStr string) {
	f.WriteString("#ifndef ")
	f.WriteString(headerStr)
	f.WriteString("\n")
	f.WriteString("#define ")
	f.WriteString(headerStr)
	f.WriteString("\n\n")
	f.WriteString("// File generated by cgoStructGen - DO NOT EDIT\n")
	f.WriteString("// Struct definitions generated for C from Go struct definitions\n")
	f.WriteString("\n")
}

func (t *CGoStructGen) templateExternCIf(f *os.File, op func()) {
	f.WriteString("#ifdef __cplusplus\n")
	f.WriteString("extern \"C\" {\n")
	f.WriteString("#endif\n\n")

	op()

	f.WriteString("#ifdef __cplusplus\n")
	f.WriteString("}\n")
	f.WriteString("#endif\n\n")
}

func (t *CGoStructGen) templateIncludes(f *os.File) {
	includes := slices.Collect(maps.Keys(t.includes))
	slices.Sort(includes)
	for _, inc := range includes {
		f.WriteString(inc.String())
		f.WriteString("\n")
	}
	f.WriteString("\n")
}

func (t *CGoStructGen) templateCStructs(f *os.File) {
	for iterName, iterFields := range t.structs {
		f.WriteString("\ttypedef struct ")
		f.WriteString(iterName)
		f.WriteString("{\n")
		for _, iterField := range iterFields {
			f.WriteString("\t\t")
			f.WriteString(iterField.String())
			f.WriteString(";\n")
		}
		f.WriteString("\t} ")
		f.WriteString(iterName)
		f.WriteString("_t;\n\n")
	}
}

func (t *CGoStructGen) templateFooter(f *os.File) {
	f.WriteString("#endif\n")
}
