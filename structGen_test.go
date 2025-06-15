package sbcgostructgen

import (
	"fmt"
	"os"
	"testing"
	"unsafe"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestGenerateForNonStruct(t *testing.T) {
	err := GenerateFor[int](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithMap(t *testing.T) {
	type s1 struct{ f1 map[int]int }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithSlice(t *testing.T) {
	type s1 struct{ f1 []int32 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithChan(t *testing.T) {
	type s1 struct{ f1 chan int }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)

	type s2 struct{ f1 <-chan int }
	err = GenerateFor[s2](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)

	type s3 struct{ f1 chan<- int }
	err = GenerateFor[s3](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithFunc(t *testing.T) {
	type s1 struct{ f1 func() }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithInterface(t *testing.T) {
	type s1 struct{ f1 fmt.Stringer }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithComplex64(t *testing.T) {
	type s1 struct{ f1 complex64 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithComplex128(t *testing.T) {
	type s1 struct{ f1 complex128 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForStructWithInt(t *testing.T) {
	type s1 struct{ f1 int }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestGenerateForStructWithUint(t *testing.T) {
	type s1 struct{ f1 uint }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestGenerateForStructWithUintptr(t *testing.T) {
	type s1 struct{ f1 uintptr }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithUnsafePointer(t *testing.T) {
	type s1 struct{ f1 unsafe.Pointer }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithInts(t *testing.T) {
	type s1 struct{ f1 int8 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)

	type s2 struct{ f1 int16 }
	err = GenerateFor[s2](New(Opts{}))
	sbtest.Nil(t, err)

	type s3 struct{ f1 int32 }
	err = GenerateFor[s3](New(Opts{}))
	sbtest.Nil(t, err)

	type s4 struct{ f1 int64 }
	err = GenerateFor[s4](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithUints(t *testing.T) {
	type s1 struct{ f1 uint8 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)

	type s2 struct{ f1 uint16 }
	err = GenerateFor[s2](New(Opts{}))
	sbtest.Nil(t, err)

	type s3 struct{ f1 uint32 }
	err = GenerateFor[s3](New(Opts{}))
	sbtest.Nil(t, err)

	type s4 struct{ f1 uint64 }
	err = GenerateFor[s4](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithBool(t *testing.T) {
	type s1 struct{ f1 bool }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithString(t *testing.T) {
	type s1 struct{ f1 string }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForStructWithArray(t *testing.T) {
	type s1 struct{ f1 [1]int32 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)

	type s2 struct{ f1 [1]int }
	err = GenerateFor[s2](New(Opts{}))
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestGenerateForStructOfStruct(t *testing.T) {
	type s2 struct{ f2 int32 }
	type s1 struct{ f1 s2 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)

	type s4 struct{ f2 int }
	type s3 struct{ f1 s4 }
	err = GenerateFor[s3](New(Opts{}))
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)

	type s6 struct{ f2 map[int32]int32 }
	type s5 struct{ f1 s6 }
	err = GenerateFor[s5](New(Opts{}))
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestGenerateForEmbededField(t *testing.T) {
	type s2 struct{ f2 int32 }
	type s1 struct{ s2 }
	err := GenerateFor[s1](New(Opts{}))
	sbtest.Nil(t, err)
}

func TestGenerateForSimpleStruct(t *testing.T) {
	type s1 struct{ f1 int8 }
	res := New(Opts{})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	sbtest.MapsMatch(
		t,
		map[include]struct{}{"<stdint.h>": {}},
		res.includes,
	)
	sbtest.Eq(t, 1, len(res.structs))
	sbtest.SlicesMatch(t, res.structs["s1"],
		[]structField{
			{
				typeModifier: typeModifier{typeMod: TypeModNone, tModAmnt: 0},
				_type:        "int8_t",
				name:         "f1",
			},
		},
	)
}

func TestWriteSingleStructOneField(t *testing.T) {
	type s1 struct{ f1 int8 }
	res := New(Opts{})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	err = res.WriteTo("./bs/testData/simpleStruct.h", "HEADER_GUARD")
	sbtest.Nil(t, err)

	data, err := os.ReadFile("./bs/testData/simpleStruct.h")
	sbtest.Nil(t, err)
	exp := `#ifndef HEADER_GUARD
#define HEADER_GUARD

// File generated by cgoStructGen - DO NOT EDIT
// Struct definitions generated for C from Go struct definitions

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	typedef struct s1{
		int8_t f1;
	} s1_t;

#ifdef __cplusplus
}
#endif

#endif
`
	sbtest.Eq(t, string(data), exp)
}

func TestWriteSingleStructMultipleFields(t *testing.T) {
	type s1 struct {
		f1 int8
		f2 uint8
		f3 float32
		f4 float64
		f5 bool
		f6 string
	}
	res := New(Opts{})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	err = res.WriteTo("./bs/testData/simpleStruct.h", "HEADER_GUARD")
	sbtest.Nil(t, err)

	data, err := os.ReadFile("./bs/testData/simpleStruct.h")
	sbtest.Nil(t, err)
	exp := `#ifndef HEADER_GUARD
#define HEADER_GUARD

// File generated by cgoStructGen - DO NOT EDIT
// Struct definitions generated for C from Go struct definitions

#include <math.h>
#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	typedef struct s1{
		int8_t f1;
		uint8_t f2;
		float_t f3;
		double_t f4;
		bool f5;
		char* f6;
	} s1_t;

#ifdef __cplusplus
}
#endif

#endif
`
	sbtest.Eq(t, string(data), exp)
}

func TestWriteMultipleStructsMultipleFields(t *testing.T) {
	type s1 struct {
		f1 int8
		f2 uint8
		f3 float32
		f4 float64
		f5 bool
		f6 string
	}
	type s2 struct {
		f7 [5]uint32
		f8 *int32
	}
	res := New(Opts{})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	err = GenerateFor[s2](res)
	sbtest.Nil(t, err)
	err = res.WriteTo("./bs/testData/simpleStruct.h", "HEADER_GUARD")
	sbtest.Nil(t, err)

	data, err := os.ReadFile("./bs/testData/simpleStruct.h")
	sbtest.Nil(t, err)
	exp := `#ifndef HEADER_GUARD
#define HEADER_GUARD

// File generated by cgoStructGen - DO NOT EDIT
// Struct definitions generated for C from Go struct definitions

#include <math.h>
#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	typedef struct s1{
		int8_t f1;
		uint8_t f2;
		float_t f3;
		double_t f4;
		bool f5;
		char* f6;
	} s1_t;

	typedef struct s2{
		uint32_t f7[5];
		int32_t* f8;
	} s2_t;

#ifdef __cplusplus
}
#endif

#endif
`
	sbtest.Eq(t, string(data), exp)
}

func TestWriteMultipleStructsReferencingEachOther(t *testing.T) {
	type s1 struct {
		f1 int8
		f2 uint8
		f3 float32
		f4 float64
		f5 bool
		f6 string
	}
	type s2 struct {
		f7  [5]uint32
		f8  *int32
		f9  s1
		f10 [10]s1
		f11 *s1
	}
	res := New(Opts{})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	err = GenerateFor[s2](res)
	sbtest.Nil(t, err)
	err = res.WriteTo("./bs/testData/simpleStruct.h", "HEADER_GUARD")
	sbtest.Nil(t, err)

	data, err := os.ReadFile("./bs/testData/simpleStruct.h")
	sbtest.Nil(t, err)
	exp := `#ifndef HEADER_GUARD
#define HEADER_GUARD

// File generated by cgoStructGen - DO NOT EDIT
// Struct definitions generated for C from Go struct definitions

#include <math.h>
#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	typedef struct s1{
		int8_t f1;
		uint8_t f2;
		float_t f3;
		double_t f4;
		bool f5;
		char* f6;
	} s1_t;

	typedef struct s2{
		uint32_t f7[5];
		int32_t* f8;
		s1_t f9;
		s1_t f10[10];
		s1_t* f11;
	} s2_t;

#ifdef __cplusplus
}
#endif

#endif
`
	sbtest.Eq(t, string(data), exp)
}

func TestWriteMultipleStructsRenaming(t *testing.T) {
	type s1 struct {
		f1 int8
		f2 uint8
		f3 float32
		f4 float64
		f5 bool
		f6 string
	}
	type s2 struct {
		f7  [5]uint32
		f8  *int32
		f9  s1
		f10 [10]s1
		f11 *s1
	}
	res := New(Opts{StructRename: map[string]string{"s1": "foo"}})
	err := GenerateFor[s1](res)
	sbtest.Nil(t, err)
	err = GenerateFor[s2](res)
	sbtest.Nil(t, err)
	err = res.WriteTo("./bs/testData/simpleStruct.h", "HEADER_GUARD")
	sbtest.Nil(t, err)

	data, err := os.ReadFile("./bs/testData/simpleStruct.h")
	sbtest.Nil(t, err)
	exp := `#ifndef HEADER_GUARD
#define HEADER_GUARD

// File generated by cgoStructGen - DO NOT EDIT
// Struct definitions generated for C from Go struct definitions

#include <math.h>
#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	typedef struct foo{
		int8_t f1;
		uint8_t f2;
		float_t f3;
		double_t f4;
		bool f5;
		char* f6;
	} foo_t;

	typedef struct s2{
		uint32_t f7[5];
		int32_t* f8;
		foo_t f9;
		foo_t f10[10];
		foo_t* f11;
	} s2_t;

#ifdef __cplusplus
}
#endif

#endif
`
	sbtest.Eq(t, string(data), exp)
}
