package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

/*

  NthRootOp
	=========

	The source code repository for nthroot.go is located at:
			https://github.com/MikeAustin71/mathopsgo.git

	The source file decimal.go is located in directory:
		MikeAustin71/mathopsgo/mathops/nthroot.go


*/

// NthRootOp - Used to extract square roots and nth roots of positive and negative
// real numbers.
//
// Currently nth roots may be passed as positive or negative values. The maximum
// Nth Root integer value is +2,147,483,647. The minimum Nth Root integer value
// is two (-2,147,483,648).
//
// The technique employed to calculate nth roots is known as the
// "shifting nth-root algorithm".
//
// See: https://en.wikipedia.org/wiki/Shifting_nth_root_algorithm
//
// Dependencies: intAry - intary.go
type NthRootOp struct {
	NthRootIntAry      *IntAry
	NthRootBigInt      *big.Int
	NthRootInt         int
	Radicand           *IntAry
	BaseNumBundles     [][]int
	LenBaseNumBundles  int
	BaseNumBundlesIdx  int
	ResultAry          IntAry
	ResultIdx          int
	ResultPrecision    int
	RequestedPrecision int
	BigOne             *big.Int
	Big10              *big.Int
	Big10ToNthPower    *big.Int
	BigZero            *big.Int
	Big3               *big.Int
	Y                  *big.Int // Root Extracted thusfar
	YPrime             *big.Int // Next Value of Y
	Minuend            *big.Int
	Subtrahend         *big.Int
	R                  *big.Int // Let R be the remainder
	RPrime             *big.Int // Let RPrime be the new value of r for next iteration
	BaseNum            *big.Int // Base Number System - always 10
	Alpha              *big.Int // Next n-digits of the radicand
	Beta               *big.Int // Next Digit of the root
}

func (nthrt *NthRootOp) Empty() {
	nthrt.NthRootInt = 0
	nthrt.NthRootBigInt = big.NewInt(0)
	nRt := IntAry{}.NewZero(0)
	nthrt.NthRootIntAry = &nRt
	nRt = IntAry{}.NewZero(0)
	nthrt.Radicand = &nRt
	nthrt.BaseNumBundles = make([][]int, 0, 500)
	nthrt.LenBaseNumBundles = 0
	nthrt.BaseNumBundlesIdx = 0
	nthrt.ResultAry = IntAry{}.NewZero(0)
	nthrt.ResultIdx = 0
	nthrt.ResultPrecision = 0
	nthrt.RequestedPrecision = 0
	nthrt.Big10 = big.NewInt(0)
	nthrt.Big10ToNthPower = big.NewInt(0)
}

// GetNthRootFloat32 - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as a type float32. In addition, the caller must supply
// input parameters for 'precision', 'nthRoot' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the float32 parameter, 'radicand', which
// will be input and positioned to the right of the decimal place.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootFloat32(radicand float32, precision, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetNthRootFloat32() "
	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloat32(radicand, precision)

	if err != nil {
		return IntAry{}.New(), fmt.Errorf(ePrefix+
			"- Error ai.SetIntAryWithFloat32(radicand, precision) "+
			"radicand= %v, precision=%v Error= %v \n", radicand, precision, err)
	}

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	return nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)

}

// GetNthRootFloat64 - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as a type float64. In addition, the caller must supply
// input parameters for 'precision', 'nthRoot' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the float64 parameter, 'radicand', which
// will be input and positioned to the right of the decimal place.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootFloat64(radicand float64, precision, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetNthRootFloat64() "

	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloat64(radicand, int(precision))

	if err != nil {
		return IntAry{}.New(), fmt.Errorf(ePrefix+
			"- Error ai.SetIntAryWithFloat64(radicand, precision) "+
			"nthRoot= %v, precision=%v Error= %v ", nthRoot, precision, err)
	}

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	return nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)
}

// GetNthRootBigFloat - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as type Big Float (*big.Float). In addition, the caller
// must supply input parameters for the 'nthRoot' and 'maxPrecision'.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootBigFloat(radicand *big.Float, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp) GetNthRootBigFloat() "

	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloatBig(radicand, -1)

	if err != nil {
		return IntAry{}.New(), fmt.Errorf(ePrefix+
			"- Error ai.SetIntAryWithFloatBig(radicand) radicand= %v Error= %v ",
			radicand, err.Error())
	}

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	return nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)
}

// GetNthRootInt - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as a type int. In addition, the caller must supply
// input parameters for 'precision', 'nthRoot' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the int64 parameter, 'radicand', which
// will be positioned to the right of the decimal place. If 'precision' is a negative
// value, an error will be returned.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootInt(radicand, precision, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetNthRootIn() "

	ai := IntAry{}.New()

	if precision < 0 {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' is a negative value! precision='%v'",
				precision)
	}

	ai.SetIntAryWithInt(radicand, uint(precision))

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	iaResult, err := nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision) "+
				"Error='%v' ", err.Error())
	}

	return iaResult, nil
}

// GetNthRootInt64 - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as a type int64. In addition, the caller must supply
// input parameters for 'precision', 'nthRoot' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the int64 parameter, 'radicand', which
// will be positioned to the right of the decimal place. 'precision' MUST BE a positive
// value. Negative values will trigger an error.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootInt64(radicand int64, precision, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetNthRootInt64() "

	if precision < 0 {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' is a negative value. INVALID! "+
				"precision='%v'", precision)
	}

	ai := IntAry{}.New()

	ai.SetIntAryWithInt64(radicand, uint(precision))

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	iaResult, err := nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision). "+
				"Error='%v'", err.Error())
	}

	return iaResult, nil
}

// GetNthRootBigInt - Calculates the Nth Root of a positive real number ('radicand')
// passed to the method as type Big Int (*big.Int). In addition, the caller
// must supply input parameters for 'precision', 'nthRoot' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the *big.Int parameter, 'radicand', which
// will be positioned to the right of the decimal place.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootBigInt(radicand *big.Int, precision, nthRoot, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetNthRootBigInt()"

	ai := IntAry{}.New()

	err := ai.SetIntAryWithBigInt(radicand, precision)

	if err != nil {
		return IntAry{}.New(), fmt.Errorf(ePrefix+
			"- Error returned by ai.SetIntAryWithBigInt(radicand, precision) "+
			"radicand= %v, precision=%v Error= %v \n", radicand, precision, err.Error())
	}

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	return nthrt.GetNthRootIntAry(&ai, &iaNthRoot, maxPrecision)
}

// GetNthRootIntAry  - Calculates the Nth Root of a real number ('radicand')
// passed to the method as a pointer to type intAry.  In addition, the
// caller must supply input parameters for 'nthRoot' and 'maxPrecision'.
//
// 'nthRoot' specifies the root which will be calculated for parameter,
// 'radicand'. Example, square root, cube root, 4th root, 9th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the
// decimal place to which the Nth root will be calculated.  If 'maxPrecision'
// is less than zero, 'maxPrecision' will be automatically set to a value
// of '4,096'.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetNthRootIntAry(radicand, nthRoot *IntAry, maxPrecision int) (IntAry, error) {

	err := nthrt.calcNthRootGateway(radicand, nthRoot, maxPrecision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf("NthRootOp.GetNthRootIntAry() Error returned from calcNthRootGateway(..). "+
				"Error= %v", err.Error())
	}

	return nthrt.ResultAry.CopyOut(), nil
}

// GetSquareRootFloat32 - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as a type float32. In addition, the caller must supply
// input parameters for 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the float32 parameter, 'radicand', which
// will be input and positioned to the right of the decimal place.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootFloat32(
	radicand float32,
	precision, maxPrecision int) (IntAry, error) {

	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloat32(radicand, int(precision))

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf("NthRootOp.GetSquareRootFloat32() "+
				"- Error ai.SetIntAryWithFloat32(radicand, precision) "+
				"radicand= %v, precision=%v Error= %v ", radicand, precision, err.Error())
	}

	return nthrt.GetSquareRootIntAry(&ai, maxPrecision)

}

// GetSquareRootFloat64 - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as a type float64. In addition, the caller must supply
// input parameter 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the float64 parameter, 'radicand', which
// will be input and positioned to the right of the decimal place.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootFloat64(
	radicand float64,
	precision, maxPrecision int) (IntAry, error) {

	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloat64(radicand, precision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf("NthRootOp.GetSquareRootFloat64() "+
				"- Error ai.SetIntAryWithFloat64(radicand, precision) "+
				"radicand= %v, precision=%v Error= %v ", radicand, precision, err.Error())
	}

	return nthrt.GetSquareRootIntAry(&ai, maxPrecision)
}

// GetSquareRootBigFloat - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as type Big Float (*big.Float). In addition, the caller
// must supply input parameter for 'maxPrecision'.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.

// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootBigFloat(radicand *big.Float, maxPrecision int) (IntAry, error) {

	ai := IntAry{}.New()

	err := ai.SetIntAryWithFloatBig(radicand, -1)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf("NthRootOp.GetSquareRootBigFloat() "+
				"- Error ai.SetIntAryWithFloatBig(radicand) radicand= %v Error= %v ",
				radicand, err.Error())
	}

	return nthrt.GetSquareRootIntAry(&ai, maxPrecision)
}

// GetSquareRootInt - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as type Int. In addition, the caller must supply input
// parameters for 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the 'int' parameter, 'radicand', which
// will be positioned to the right of the decimal place. If 'precision' is a negative
// value an error will be returned.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootInt(radicand int, precision, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetSquareRootInt() "

	ai := IntAry{}.New()

	if precision < 0 {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' is a negative value! precision='%v'",
				precision)
	}

	ai.SetIntAryWithInt(radicand, uint(precision))

	return nthrt.GetSquareRootIntAry(&ai, maxPrecision)
}

// GetSquareRootInt32 - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as type Int32. In addition, the caller must supply input
// parameters for 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the 'int32' parameter, 'radicand', which
// will be positioned to the right of the decimal place. 'precision' MUST BE a positive
// value. Negative values will trigger an error.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootInt32(radicand int32, precision, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetSquareRootInt32() "

	if precision < 0 {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' is less than zero. INVALID! "+
				"precision='%v' ", precision)
	}

	ai := IntAry{}.New()

	ai.SetIntAryWithInt32(radicand, uint(precision))

	iaResult, err := nthrt.GetSquareRootIntAry(&ai, maxPrecision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.GetSquareRootIntAry(&ai, maxPrecision). "+
				"Error='%v' ", err.Error())
	}

	return iaResult, nil
}

// GetSquareRootInt64 - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as type Int64. In addition, the caller must supply input
// parameters for 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the 'int64' parameter, 'radicand', which
// will be positioned to the right of the decimal place. 'precision' MUST BE a positive
// value. Negative values will trigger an error.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootInt64(radicand int64, precision, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetSquareRootInt64() "

	if precision < 0 {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' is less than zero. INVALID! "+
				"precision='%v'", precision)
	}

	ai := IntAry{}.New()

	ai.SetIntAryWithInt64(radicand, uint(precision))

	iaResult, err := nthrt.GetSquareRootIntAry(&ai, maxPrecision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.GetSquareRootIntAry(&ai, maxPrecision). "+
				"Error='%v'", err.Error())
	}

	return iaResult, nil
}

// GetSquareRootBigInt - Calculates the Square Root of a positive real number ('radicand')
// passed to the method as a pointer to type Big Int (*big.Int). In addition, the caller
// must supply input parameters for 'precision' and 'maxPrecision'.
//
// 'precision' specifies the number of digits in the *big.Int parameter, 'radicand', which
// will be positioned to the right of the decimal place.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootBigInt(radicand *big.Int, precision, maxPrecision int) (IntAry, error) {
	ai := IntAry{}.New()

	err := ai.SetIntAryWithBigInt(radicand, precision)

	if err != nil {
		return IntAry{}.New(), fmt.Errorf("NthRootOp.GetSquareRootBigInt() - Error ai.SetIntAryWithBigInt(radicand, precision) radicand= %v, precision=%v Error= %v ", radicand, precision, err)
	}

	return nthrt.GetSquareRootIntAry(&ai, maxPrecision)
}

// GetSquareRootIntAry - Calculates the square Root of a positive real number ('radicand')
// passed to the method as a pointer to type intAry. In addition, the caller
// must supply input parameters for 'maxPrecision'.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the square root will be calculated.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) GetSquareRootIntAry(radicand *IntAry, maxPrecision int) (IntAry, error) {

	ePrefix := "NthRootOp.GetSquareRootIntAry() "

	iaNthRoot := IntAry{}.NewTwo(0)

	err := nthrt.calcNthRootGateway(radicand, &iaNthRoot, maxPrecision)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"- Error returned from nthrt.calcNthRootGateway(radicand, &iaNthRoot, maxPrecision). "+
				"Error= %v", err.Error())
	}

	return nthrt.ResultAry.CopyOut(), nil

}

// NewNthRoot - Returns the result of an Nth Root calculation. The result
// is returned as a type 'IntAry'.
//
// This method calculates the Nth Root of a real number ('radicand') passed
// to the method as a pointer to type intAry.  In addition, the caller must
// supply input parameters for 'nthRoot' and 'maxPrecision'.
//
// 'nthRoot' specifies the root which will be calculated for parameter,
// 'radicand'. Example, square root, cube root, 4th root, 9th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the
// decimal place to which the Nth root will be calculated. If 'maxPrecision'
// is less than zero, 'maxPrecision' will be automatically set to a value
// of '4,096'.
//
// The calculation result is returned as an intAry object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt NthRootOp) NewNthRoot(radicand, nthRoot *IntAry, maxPrecision int) (IntAry, error) {

	nthRtOp := NthRootOp{}

	return nthRtOp.GetNthRootIntAry(radicand, nthRoot, maxPrecision)
}

// SetNthRootIntAry  - Calculates the Nth Root of a number ('radicand') passed to the
// method as a pointer to type intAry.  In addition, the caller must supply input
// parameters for 'nthRoot' and 'maxPrecision'.
//
// The difference between this method, 'SetNthRootIntAry' and 'OriginalNthRoot' is in
// the return value.  This method, 'SetNthRootIntAry' does not return the result. Instead,
// the calculation result is stored in the NthRootOp intAry Object, NthRootOp.ResultAry.
// This method is primarily for use by other low level routines seeking to improve performance
// by avoiding the return of a new intAry object.
//
// Nth root specifies the root which will be calculated for parameter, 'radicand'. Example,
// square root, cube root, 4th root, 9th root etc.
//
// 'maxPrecision' specifies the number of decimals to the right of the decimal place to
// which the Nth root will be calculated.
//
// The calculation result is stored in the NthRootOp field, 'NthRootOp.ResultAry'.
// 'NthRootOp.ResultAry' is an intAry Object.
//
// Note: A negative 'radicand' value with an even nthRoot will generate an error.
func (nthrt *NthRootOp) SetNthRootIntAry(radicand, nthRoot *IntAry, maxPrecision int) error {

	return nthrt.calcNthRootGateway(radicand, nthRoot, maxPrecision)
}

// calcNthRootGateway - This method is is the primary means by which the nth root
// calculation is accessed. All screening and validation of input parameters 'radicand'
// and 'nthRoot' are performed here. If 'radicand' and 'nthRoot' pass all tests for
// validity, this method proceeds to perform the nth root calculation and store the result
// in the NthRootOp data structure. The final result of the nth root calculation is therefore
// stored in data structure element, 'NthRootOp.ResultAry'.
//
// Note: If maxPrecision is less than zero it will automatically be set to a value of
// 4,096 decimal places.
func (nthrt *NthRootOp) calcNthRootGateway(radicand, nthRoot *IntAry, maxPrecision int) error {

	ePrefix := "NthRootOp.calcNthRootGateway() "

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	nthrt.ResultAry = IntAry{}.NewInt32(0, 0)

	// If radicand is zero, the result will always be zero.
	if radicand.IsZero() {
		nthrt.ResultAry.SetIntAryToZero(uint(maxPrecision))
		return nil
	}

	if nthRoot.IsZero() {
		nthrt.ResultAry.SetIntAryToOne(maxPrecision)
		return nil
	}

	nthRootPrecision := nthRoot.GetPrecision()
	nthRootSign := nthRoot.GetSign()

	if nthRootPrecision == 0 && nthRootSign == -1 {

		return nthrt.calcNegativeIntegerNthRoot(radicand, nthRoot, maxPrecision)

	} else if nthRoot.precision > 0 && nthRootSign == -1 {

		return nthrt.calcNegativeFractionalNthRoot(radicand, nthRoot, maxPrecision)

	} else if nthRoot.precision == 0 && nthRootSign == 1 {

		return nthrt.calcPositiveIntegerNthRoot(radicand, nthRoot, maxPrecision)

	} else if nthRoot.precision > 0 && nthRootSign == 1 {

		return nthrt.calcPositiveFractionalNthRoot(radicand, nthRoot, maxPrecision)

	}

	return fmt.Errorf(ePrefix+
		"Error - 'nthRoot' configuration failed to match acceptable calculation patterns! "+
		"nthRoot='%v' ", nthRoot.GetNumStr())

}

// calcPositiveIntegerNthRoot - Calculates the Nth Root of a radicand where
// nth root is both positive and an integer value.
func (nthrt *NthRootOp) calcPositiveIntegerNthRoot(
	radicand, nthRoot *IntAry,
	maxPrecision int) error {

	ePrefix := "NthRootOp.calcPositiveIntegerNthRoot() "

	if nthRoot.GetSign() != 1 {
		return fmt.Errorf(ePrefix+
			"Error expected postive 'nthRoot'. 'nthRoot' is negative! "+
			"nthRoot= %v", nthRoot.GetNumStr())
	}

	if radicand.GetSign() == -1 {

		if nthRoot.IsEvenNumber() {
			return errors.New(ePrefix +
				"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even.\n")
		}
	}

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	err := nthrt.initialize(radicand, nthRoot, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned from initialization. Error= %v",
			err.Error())
	}

	err = nthrt.doRootExtraction()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.doRootExtraction() "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// calcNegativeIntegerNthRoot - Calculates the Nth Root of a radicand where
// nth root is both negative and an integer value.
func (nthrt *NthRootOp) calcNegativeIntegerNthRoot(
	radicand, nthRoot *IntAry,
	maxPrecision int) error {

	ePrefix := "NthRootOp.calcNegativeIntegerNthRoot() "

	if nthRoot.GetSign() != -1 {
		return fmt.Errorf(ePrefix+
			"Error expected negative 'nthRoot'. 'nthRoot' is positive! "+
			"nthRoot= %v", nthRoot.GetNumStr())
	}

	if radicand.GetSign() == -1 {

		if nthRoot.IsEvenNumber() {
			return errors.New(ePrefix +
				"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even.\n")
		}
	}

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	// Change sign from negative (-) to positive (+)
	nthRoot.ChangeSign()

	err := nthrt.initialize(radicand, nthRoot, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned from initialization. Error= %v",
			err.Error())
	}

	err = nthrt.doRootExtraction()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.doRootExtraction() "+
			"Error='%v' ", err.Error())
	}

	result, err := nthrt.ResultAry.Inverse(maxPrecision + 100)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.ResultAry.Inverse() "+
			"Error='%v' ", err.Error())

	}

	if result.GetPrecision() > maxPrecision {
		result.RoundToPrecision(maxPrecision)
	}

	// Change sign from positive (+), back to negative (-)
	nthRoot.ChangeSign()

	nthrt.Empty()
	nthrt.NthRootIntAry.CopyIn(nthRoot, false)
	nthrt.ResultAry.CopyIn(&result, false)
	nthrt.Radicand = radicand.CopyOutPtr()
	nthrt.ResultPrecision = nthrt.ResultAry.GetPrecision()
	nthrt.RequestedPrecision = maxPrecision
	nthrt.NthRootBigInt, err = nthrt.NthRootIntAry.GetBigInt()
	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.NthRootIntAry.GetBigInt() "+
			"Error='%v' ", err.Error())

	}

	return nil
}

// calcPositiveFractionalNthRoot - Calculates the Nth Root of a radicand where
// nth root is both negative and an integer value.
func (nthrt *NthRootOp) calcPositiveFractionalNthRoot(
	radicand, nthRoot *IntAry,
	maxPrecision int) error {

	ePrefix := "NthRootOp.calcPositiveFractionalNthRoot() "

	if nthRoot.GetPrecision() < 1 {
		return fmt.Errorf(ePrefix+
			"Error- Expected fractional 'nthRoot'. 'nthRoot' is an integer value. "+
			"nthRoot='%v' ", nthRoot.GetNumStr())
	}

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	internalMaxPrecision := maxPrecision + 100

	fracIntAry := FracIntAry{}.NewFracIntAry(nthRoot)

	err := fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision) "+
			"internalMaxPrecision='%v' Error='%v' ",
			internalMaxPrecision, err.Error())
	}

	newRadicand := radicand.CopyOut()

	err = IntAryMathPower{}.Pwr(
		&newRadicand,
		&fracIntAry.Denominator,
		0,
		internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by IntAryMathPower{}.Pwr(radicand, fracIntAry.Denominator) "+
			"Error='%v' ", err.Error())
	}

	if newRadicand.GetSign() == -1 {

		if fracIntAry.Numerator.IsEvenNumber() {
			return errors.New(ePrefix +
				"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even.\n")
		}
	}

	err = nthrt.initialize(&newRadicand, &fracIntAry.Numerator, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned from initialization. Error= %v",
			err.Error())
	}

	err = nthrt.doRootExtraction()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.doRootExtraction() "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// calcNegativeFractionalNthRoot - Calculates the Nth Root of a radicand where
// nth root is both negative and an integer value.
func (nthrt *NthRootOp) calcNegativeFractionalNthRoot(
	radicand, nthRoot *IntAry,
	maxPrecision int) error {

	ePrefix := "NthRootOp.calcNegativeFractionalNthRoot() "

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	internalMaxPrecision := maxPrecision + 100

	if nthRoot.GetSign() != -1 {
		return fmt.Errorf(ePrefix+
			"Error: Expected negative 'nthRoot'. Instead, 'nthRoot' is positive! "+
			"nthRoot='%v' ", nthRoot.GetNumStr())
	}

	// Change sign from negative (-) to positive (+)
	nthRoot.ChangeSign()

	if nthRoot.GetPrecision() < 1 {
		return fmt.Errorf(ePrefix+
			"Error- Expected fractional 'nthRoot'. 'nthRoot' is an integer value. "+
			"nthRoot='%v' ", nthRoot.GetNumStr())
	}

	fracIntAry := FracIntAry{}.NewFracIntAry(nthRoot)

	err := fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by fracIntAry.ReduceToLowestCommonDenom(4096) "+
			"Error='%v' ", err.Error())
	}

	newRadicand := radicand.CopyOut()

	err = IntAryMathPower{}.Pwr(
		&newRadicand,
		&fracIntAry.Denominator,
		0,
		internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by IntAryMathPower{}.Pwr(radicand, fracIntAry.Denominator) "+
			"Error='%v' ", err.Error())
	}

	if newRadicand.GetSign() == -1 {

		if fracIntAry.Numerator.IsEvenNumber() {
			return errors.New(ePrefix +
				"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even. \n")
		}
	}

	err = nthrt.initialize(&newRadicand, &fracIntAry.Numerator, internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned from initialization. Error= %v",
			err.Error())
	}

	err = nthrt.doRootExtraction()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.doRootExtraction() "+
			"Error='%v' ", err.Error())
	}

	result, err := nthrt.ResultAry.Inverse(internalMaxPrecision + 10)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.ResultAry.Inverse() "+
			"Error='%v' ", err.Error())

	}

	if result.GetPrecision() > maxPrecision {
		result.RoundToPrecision(maxPrecision)
	}

	// Change sign from positive (+), back to negative (-)
	fracIntAry.Numerator.ChangeSign()

	nthrt.Empty()
	nthrt.NthRootIntAry.CopyIn(&fracIntAry.Numerator, false)
	nthrt.ResultAry.CopyIn(&result, false)
	nthrt.Radicand = newRadicand.CopyOutPtr()
	nthrt.ResultPrecision = nthrt.ResultAry.GetPrecision()
	nthrt.RequestedPrecision = maxPrecision
	nthrt.NthRootBigInt, err = nthrt.NthRootIntAry.GetBigInt()
	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.NthRootIntAry.GetBigInt() "+
			"Error='%v' ", err.Error())

	}

	return nil
}

// initialize - Initializes the data fields of the NthRootOp structure and validates the
// radicand and nthRoot numerical values passed to the Nth Root Calculation.
func (nthrt *NthRootOp) initialize(radicand, nthRoot *IntAry, maxPrecision int) error {

	ePrefix := "NthRootOp.initialize() "

	err := radicand.IsValid(ePrefix + "'radicand' Invalid - ")

	if err != nil {
		return err
	}

	err = nthRoot.IsValid(ePrefix + "'nthRoot' Invalid - ")

	if err != nil {
		return err
	}

	if nthRoot.GetSign() < 1 {
		return fmt.Errorf(ePrefix+
			"Error nthRoot is a negative number! nthRoot='%v'",
			nthRoot.GetNumStr())
	}

	if maxPrecision < 0 {
		return fmt.Errorf(ePrefix+
			"Error 'maxPrecision' is less than zero! maxPrecision='%v'",
			maxPrecision)
	}

	if nthRoot.IsOne() {
		return fmt.Errorf(ePrefix+
			"- Input Parameter 'nthRoot' INVALID! 'nthRoot' cannot equal 1. "+
			"nthRoot= %v\n", nthRoot.GetNumStr())
	}

	nthrt.NthRootInt, err = nthRoot.GetInt()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by nthRoot.GetInt() Error='%v' ",
			err.Error())
	}

	nthrt.NthRootIntAry = nthRoot.CopyOutPtr()

	bINum, err := nthRoot.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by nthRoot.GetBigInt() Error='%v' ",
			err.Error())
	}

	nthrt.NthRootBigInt = big.NewInt(0).Set(bINum)

	nthrt.Radicand = radicand
	nthrt.RequestedPrecision = maxPrecision

	err = nthrt.bundleInts()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"- Error returned from nthrt.bundleInts(). Error= %v", err.Error())
	}

	err = nthrt.bundleFracs()

	if err != nil {
		return fmt.Errorf("NthRootOp.initialize() - Error returned from nthrt.bundleFracs(). Error= %v", err)
	}

	err = nthrt.calcPrecision()

	if err != nil {
		return fmt.Errorf("NthRootOp.initialize() - Error returned from nthrt.CalcPrecision(). Error= %v", err)
	}

	// Set constants for calculations

	nthrt.BigOne = big.NewInt(1)
	nthrt.Big3 = big.NewInt(3)
	nthrt.Big10 = big.NewInt(10)
	nthrt.Big10ToNthPower = big.NewInt(0).Exp(nthrt.Big10, big.NewInt(int64(nthrt.NthRootInt)), nil)
	nthrt.BigZero = big.NewInt(0)
	nthrt.ResultAry = IntAry{}.New()
	nthrt.ResultAry.SetPrecision(nthrt.RequestedPrecision, false)
	nthrt.ResultAry.SetSign(radicand.GetSign())
	nthrt.Y = big.NewInt(0)
	nthrt.YPrime = big.NewInt(0)
	nthrt.R = big.NewInt(0)
	nthrt.RPrime = big.NewInt(0)
	nthrt.BaseNum = big.NewInt(10)
	nthrt.Alpha = big.NewInt(0)
	nthrt.Beta = big.NewInt(0)
	return nil
}

// bundleInts - computes the bundle size and creates the
// bundle array elements necessary to process the integer
// digits in the original number.
func (nthrt *NthRootOp) bundleInts() error {

	intNums, err := nthrt.Radicand.GetIntegerDigits()

	intNumsStats := intNums.GetIntAryStats()

	if err != nil {
		return fmt.Errorf("NthRootOp.bundleInts() - Error= %v", err)
	}

	if intNumsStats.IntAryLen < 1 {
		return errors.New("NthRootOp.bundleInts() intNums array length less than 1")
	}

	bundleSize := 0

	if intNumsStats.IntAryLen <= nthrt.NthRootInt {

		bundleSize = 1

	} else {

		bundleSize = intNumsStats.IntAryLen / nthrt.NthRootInt

		if intNumsStats.IntAryLen > ((intNumsStats.IntAryLen / nthrt.NthRootInt) * nthrt.NthRootInt) {

			bundleSize++

		}

	}

	nthrt.BaseNumBundles = make([][]int, bundleSize)

	bundleIdx := bundleSize - 1
	intAryIdx := 0
	intAry, _ := intNums.GetIntAryElements()
	for j := intNumsStats.IntAryLen - 1; j >= 0; j -= nthrt.NthRootInt {
		bundle := make([]int, nthrt.NthRootInt)
		intAryIdx = j
		for k := nthrt.NthRootInt - 1; k >= 0; k-- {
			if intAryIdx < 0 {
				break
			}

			bundle[k] = int(intAry[intAryIdx])
			intAryIdx--
		}

		nthrt.BaseNumBundles[bundleIdx] = append(nthrt.BaseNumBundles[bundleIdx], bundle...)

		bundleIdx--
	}

	return nil
}

// bundleFracs - computes the numeric bundle size and creates
// the bundle array elements required to process the fractional
// digits in the original number.
func (nthrt *NthRootOp) bundleFracs() error {

	if nthrt.Radicand.GetPrecision() < 1 {
		return nil
	}

	fracNums, err := nthrt.Radicand.GetFractionalDigits()

	iFracNumStats := fracNums.GetIntAryStats()
	iFAry, _ := fracNums.GetIntAryElements()

	if err != nil {
		return fmt.Errorf("NthRootOp.bundleFracs() Error Returned from nthrt.OriginalRadicand.GetFractionalDigits() - Error= %v", err)
	}

	for i := 1; i < iFracNumStats.IntAryLen; i += nthrt.NthRootInt {

		bundle := make([]int, nthrt.NthRootInt)

		for j := 0; j < nthrt.NthRootInt; j++ {

			if i+j < iFracNumStats.IntAryLen {
				bundle[j] = int(iFAry[i+j])
			}
		}

		nthrt.BaseNumBundles = append(nthrt.BaseNumBundles, bundle)

	}

	return nil
}

// CalcPrecision - calculates the bundle size and creates
// the bundle array elements necessary to process the nth
// root to the requested number of decimal places to the
// right of the decimal point.
func (nthrt *NthRootOp) calcPrecision() error {

	if nthrt.Radicand.GetPrecision() < 0 {
		return fmt.Errorf("NthRootOp.CalcPrecision() - Existing precision is less than zero! OriginalRadicand.precision= %v", nthrt.Radicand.GetPrecision())
	}

	existingPrecision := nthrt.Radicand.GetPrecision() / nthrt.NthRootInt

	if nthrt.Radicand.GetPrecision() != existingPrecision*nthrt.NthRootInt {
		existingPrecision++
	}

	bundle := make([]int, nthrt.NthRootInt)

	if nthrt.RequestedPrecision <= existingPrecision {

		// nthrt.BaseNumBundles = append(nthrt.BaseNumBundles, bundle)

		nthrt.RequestedPrecision = existingPrecision

		// return nil
	}

	nthrt.RequestedPrecision++

	deltaPrecision := nthrt.RequestedPrecision - existingPrecision

	for i := 0; i < deltaPrecision; i++ {
		nthrt.BaseNumBundles = append(nthrt.BaseNumBundles, bundle)
	}

	nthrt.LenBaseNumBundles = len(nthrt.BaseNumBundles)

	nthrt.ResultPrecision = nthrt.RequestedPrecision

	return nil
}

func (nthrt *NthRootOp) doRootExtraction() error {

	nthrt.Y = big.NewInt(0)
	nthrt.Minuend = big.NewInt(0)
	nthrt.Subtrahend = big.NewInt(0)
	nthrt.R = big.NewInt(0)

	for i := 0; i < nthrt.LenBaseNumBundles; i++ {

		nthrt.findNextRoot(i)

	}

	nthrt.ResultAry.SetSign(nthrt.Radicand.GetSign())
	nthrt.ResultAry.OptimizeIntArrayLen(false)
	err := nthrt.ResultAry.RoundToPrecision(nthrt.RequestedPrecision - 1)

	if err != nil {
		return fmt.Errorf("doRootExtraction() - Error returned from sqrt.ResultAry.RoundToPrecision(sqrt.ResultPrecision). nthrt.ResultPrecision= %v  Error= %v", nthrt.ResultPrecision, err)
	}

	return nil
}

func (nthrt *NthRootOp) findNextRoot(bundleIdx int) {

	bundle := nthrt.getBundleBigInt(bundleIdx)

	// alpha = next n-digits of radicand
	nthrt.Alpha = big.NewInt(0).Set(bundle)

	nthrt.RPrime = big.NewInt(-1)

	// nthrt.R already set
	// n = nthrt.NthRootInt already set
	// nthrt.Y already set
	// nthrt.BaseNum = 10
	// nthrt.Big10ToNthPower = BaseNum^n

	itatr := big.NewInt(9)
	term_1a := big.NewInt(0)
	term_1b := big.NewInt(0)

	term_2a1 := big.NewInt(0)
	term_2a2 := big.NewInt(0)
	term_2a := big.NewInt(0)

	term_2b := big.NewInt(0)
	term_2b1 := big.NewInt(0)
	term_2b2 := big.NewInt(0)

	term_1a = big.NewInt(0).Mul(nthrt.Big10ToNthPower, nthrt.R)
	term_1b = big.NewInt(0).Set(nthrt.Alpha)
	nthrt.Minuend = big.NewInt(0).Add(term_1a, term_1b)

	for itatr.Cmp(nthrt.BigZero) > -1 &&
		nthrt.RPrime.Cmp(nthrt.BigZero) == -1 {

		nthrt.Beta = big.NewInt(0).Set(itatr)
		nthrt.YPrime = big.NewInt(0).Mul(nthrt.Y, nthrt.Big10)
		nthrt.YPrime = big.NewInt(0).Add(nthrt.YPrime, nthrt.Beta)

		term_2a1 = big.NewInt(0).Mul(nthrt.BaseNum, nthrt.Y)
		term_2a2 = big.NewInt(0).Add(term_2a1, nthrt.Beta)
		term_2a = big.NewInt(0).Exp(term_2a2, big.NewInt(int64(nthrt.NthRootInt)), nil)

		term_2b1 = big.NewInt(0).Set(nthrt.Big10ToNthPower)
		term_2b2 = big.NewInt(0).Exp(nthrt.Y, big.NewInt(int64(nthrt.NthRootInt)), nil)

		term_2b = big.NewInt(0).Mul(term_2b1, term_2b2)

		nthrt.Subtrahend = big.NewInt(0).Sub(term_2a, term_2b)

		nthrt.RPrime = big.NewInt(0).Sub(nthrt.Minuend, nthrt.Subtrahend)

		itatr = big.NewInt(0).Sub(itatr, nthrt.BigOne)
	}

	nthrt.R = big.NewInt(0).Set(nthrt.RPrime)
	nthrt.Y = big.NewInt(0).Set(nthrt.YPrime)
	nthrt.ResultAry.AppendToIntAry(uint8(nthrt.Beta.Int64()))

}

func (nthrt *NthRootOp) getBundleBigInt(idx int) *big.Int {

	bigBundleVal := big.NewInt(0)

	for i := 0; i < nthrt.NthRootInt; i++ {

		bigBundleVal = big.NewInt(0).Mul(bigBundleVal, nthrt.Big10)
		bigBundleVal = big.NewInt(0).Add(bigBundleVal, big.NewInt(0).SetInt64(int64(nthrt.BaseNumBundles[idx][i])))

	}

	return bigBundleVal
}
