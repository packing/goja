package ftoa

import (
	"math"
	"strconv"
	"testing"
)

func _testFToStr(num float64, mode FToStrMode, precision int, expected string, t *testing.T) {
	buf := FToStr(num, mode, precision, nil)
	if s := string(buf); s != expected {
		t.Fatalf("expected: '%s', actual: '%s", expected, s)
	}
	if !math.IsNaN(num) && num != 0 && !math.Signbit(num) {
		_testFToStr(-num, mode, precision, "-"+expected, t)
	}
}

func testFToStr(num float64, mode FToStrMode, precision int, expected string, t *testing.T) {
	t.Run("", func(t *testing.T) {
		t.Parallel()
		_testFToStr(num, mode, precision, expected, t)
	})
}

func TestDtostr(t *testing.T) {
	testFToStr(0, ModeStandard, 0, "0", t)
	testFToStr(1, ModeStandard, 0, "1", t)
	testFToStr(9007199254740991, ModeStandard, 0, "9007199254740991", t)
	testFToStr(math.MaxInt64, ModeStandardExponential, 0, "9.223372036854776e+18", t)
	testFToStr(1e-5, ModeFixed, 1, "0.0", t)
	testFToStr(8.85, ModeExponential, 2, "8.8e+0", t)
	testFToStr(885, ModeExponential, 2, "8.9e+2", t)
	testFToStr(25, ModeExponential, 1, "3e+1", t)
	testFToStr(1e-6, ModeFixed, 7, "0.0000010", t)
	testFToStr(math.Inf(1), ModeStandard, 0, "Infinity", t)
	testFToStr(math.NaN(), ModeStandard, 0, "NaN", t)
	testFToStr(math.SmallestNonzeroFloat64, ModeExponential, 40, "4.940656458412465441765687928682213723651e-324", t)
}

func BenchmarkDtostrSmall(b *testing.B) {
	var buf [128]byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FToStr(math.Pi, ModeExponential, 0, buf[:0])
	}
}

func BenchmarkDtostrBig(b *testing.B) {
	var buf [128]byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FToStr(math.SmallestNonzeroFloat64, ModeExponential, 40, buf[:0])
	}
}

func BenchmarkAppendFloatBig(b *testing.B) {
	var buf [128]byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(buf[:0], math.SmallestNonzeroFloat64, 'e', 40, 64)
	}
}

func BenchmarkAppendFloatSmall(b *testing.B) {
	var buf [128]byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(buf[:0], math.Pi, 'e', -1, 64)
	}
}
