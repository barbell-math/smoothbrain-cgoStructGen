package sbcgostructgen

import (
	"fmt"
	"testing"
	"unsafe"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestCheckTypeNonStruct(t *testing.T) {
	err := GenerateFor[int]("", &Opts{ExitOnErr: false, DryRun: true})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithMap(t *testing.T) {
	err := GenerateFor[struct{ f1 map[int]int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithChan(t *testing.T) {
	err := GenerateFor[struct{ f1 chan int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
	err = GenerateFor[struct{ f1 <-chan int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
	err = GenerateFor[struct{ f1 chan<- int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithFunc(t *testing.T) {
	err := GenerateFor[struct{ f1 func() }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithInterface(t *testing.T) {
	err := GenerateFor[struct{ f1 fmt.Stringer }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithComplex64(t *testing.T) {
	err := GenerateFor[struct{ f1 complex64 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithComplex128(t *testing.T) {
	err := GenerateFor[struct{ f1 complex128 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}

func TestCheckTypeStructWithInt(t *testing.T) {
	err := GenerateFor[struct{ f1 int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestCheckTypeStructWithUint(t *testing.T) {
	err := GenerateFor[struct{ f1 uint }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestCheckTypeStructWithUintptr(t *testing.T) {
	err := GenerateFor[struct{ f1 uintptr }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithUnsafePointer(t *testing.T) {
	err := GenerateFor[struct{ f1 unsafe.Pointer }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithInts(t *testing.T) {
	err := GenerateFor[struct{ f1 int8 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 int16 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 int32 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 int64 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithUints(t *testing.T) {
	err := GenerateFor[struct{ f1 uint8 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 uint16 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 uint32 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 uint64 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithBool(t *testing.T) {
	err := GenerateFor[struct{ f1 bool }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithString(t *testing.T) {
	err := GenerateFor[struct{ f1 string }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
}

func TestCheckTypeStructWithArray(t *testing.T) {
	err := GenerateFor[struct{ f1 [1]int32 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 [1]int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestCheckTypeStructWithSlice(t *testing.T) {
	err := GenerateFor[struct{ f1 []int32 }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 []int }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
}

func TestCheckTypeStructOfStruct(t *testing.T) {
	err := GenerateFor[struct{ f1 struct{ f2 int32 } }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.Nil(t, err)
	err = GenerateFor[struct{ f1 struct{ f2 int } }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, UnderspecifiedTypeErr, err)
	err = GenerateFor[struct{ f1 struct{ f2 map[int32]int32 } }]("", &Opts{
		ExitOnErr: false, DryRun: true,
	})
	sbtest.ContainsError(t, InvalidTypeErr, err)
}
