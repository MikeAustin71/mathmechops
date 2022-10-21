package mathops

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

/*

	IntAry Overview And General Usage:
	==================================

	Source file 'intary.go' contains a structure, 'IntAry', which is designed
 	to perform a variety of math operations on integer strings.

	This Type is capable of performing highly accurate operations on very large
  numbers. For example, the directory 'MikeAustin71/mathopsgo/examples/eulersnumbercalc'
	contains an example which calculates Euler's Number out to 1,000 digits.

	The IntAry Type also has a backup and restore feature.

	Dependencies
	============
	The 'IntAry' object has the following dependency:

	nthroot.go - MikeAustin71/mathopsgo/mathops/nthroot.go


	INumMgr
	========

	The IntAry Type implements the INumMgr interface.

 Source Code Repository:
 =======================
 	https://github.com/MikeAustin71/mathopsgo.git

 Local File:
 ===========
 MikeAustin71\mathopsgo\mathops\intary.txt

*/

// FracIntAry - A fraction represented by a numerator and a denominator.
// Both numerator and denominator are of type intAry
type FracIntAry struct {
	Numerator   IntAry
	Denominator IntAry
}

// NewBigInts - Creates a new FracIntAry type from two *big.Int types passed
// as input parameters.
func (fIa FracIntAry) NewBigInts(numerator, denominator *big.Int) (FracIntAry, error) {
	ePrefix := "FracIntAry.NewBigInts() "

	iaNumerator, err := IntAry{}.NewBigInt(numerator, 0)

	if err != nil {
		return FracIntAry{},
			fmt.Errorf(ePrefix+"- Error returned by IntAry{}.NewBigInt(numerator, 0) "+
				"Error='%v' ", err)
	}

	iaDenominator, err := IntAry{}.NewBigInt(denominator, 0)

	if err != nil {
		return FracIntAry{},
			fmt.Errorf(ePrefix+"- Error returned by IntAry{}.NewBigInt(denominator, 0) "+
				"Error='%v' ", err)
	}

	newFracIntAry := FracIntAry{}.NewIntArys(&iaNumerator, &iaDenominator)

	return newFracIntAry, nil
}

// NewNumStrs - Creates a new FracIntAry type by passing input parameters numerator and
// denominator as number strings.
func (fIa FracIntAry) NewNumStrs(numerator, denominator string) (FracIntAry, error) {

	var err error

	ePrefix := "FracIntAry.NewNumStrs() "
	fIa2 := FracIntAry{}

	fIa2.Numerator, err = IntAry{}.NewNumStr(numerator)

	if err != nil {
		return FracIntAry{}, fmt.Errorf(ePrefix+
			"- Error returned from intAry{}.NewNumStr(numerator). Error= %v", err)
	}

	fIa2.Denominator, err = IntAry{}.NewNumStr(denominator)

	if err != nil {
		return FracIntAry{},
			fmt.Errorf(ePrefix+
				"- Error returned from intAry{}.NewNumStr(denominator). Error= %v", err)
	}

	return fIa2, nil
}

// NewIntArys - Creates a type FracIntAry by passing numerator and denominator
// input parameters of type *intAry
func (fIa FracIntAry) NewIntArys(numerator, denominator *IntAry) FracIntAry {

	fIa2 := FracIntAry{}

	fIa2.Numerator = numerator.CopyOut()
	fIa2.Denominator = denominator.CopyOut()

	return fIa2
}

// NewFracIntAry - Creates a FracIntAry instance from a single IntAry object.
// The IntAry input parameter is converted into an equivalent fraction.
func (fIa FracIntAry) NewFracIntAry(ia *IntAry) FracIntAry {

	fIa2 := FracIntAry{}

	if ia.GetPrecision() == 0 {

		fIa2.Numerator = ia.CopyOut()
		fIa2.Denominator = IntAry{}.NewOne(0)

		return fIa2
	}

	precision := ia.GetPrecision()
	fIa2.Numerator = ia.CopyOut()

	if precision > 0 {
		fIa2.Numerator.ShiftPrecisionRight(uint(precision))
	}

	fIa2.Denominator = IntAry{}.NewOne(0)
	IntAryMathMultiply{}.MultiplyByTenToPower(&fIa2.Denominator, uint(precision))

	return fIa2
}

// CopyOut - Creates and returns a copy of the current
// FracIntAry.
func (fIa *FracIntAry) CopyOut() FracIntAry {

	newFrac := FracIntAry{}

	newFrac.Numerator = fIa.Numerator.CopyOut()
	newFrac.Denominator = fIa.Denominator.CopyOut()

	return newFrac
}

// CopyIn - Receives a pointer to an incoming FracIntAry and copies
// the values into the current FracIntAry.
func (fIa *FracIntAry) CopyIn(fIa2 *FracIntAry) {

	fIa.Numerator = fIa2.Numerator.CopyOut()

	fIa.Denominator = fIa2.Denominator.CopyOut()

}

// GetRationalValue - Converts the fraction and returns the value as a
// big rational number (*big.Rat).
//
// Input parameter maxPrecision determines the maximum number of decimal
// places to the right of the decimal point contained in the result.
//
// If the value of maxPrecision is -1, maximum precision will default to
// 4096 decimal places. maxPrecision values less than -1 will trigger an
// error.
func (fIa *FracIntAry) GetRationalValue(maxPrecision int) (*big.Rat, error) {

	if maxPrecision < -1 {
		return big.NewRat(1, 1), fmt.Errorf("GetRationalValue() - maxPrecision is less than -1 and therefore INVALID. maxPrecision= %v", maxPrecision)
	}

	if maxPrecision == -1 {
		maxPrecision = 4096
	}

	if fIa.Numerator.GetPrecision() == 0 && fIa.Denominator.GetPrecision() == 0 {
		fRat, ok := big.NewRat(1, 1).SetString(fIa.Numerator.GetNumStr() + "/" + fIa.Denominator.GetNumStr())

		if !ok {
			return big.NewRat(1, 1), errors.New("GetRationalValue() - big.NewFloat(0).SetString(fracIa.GetNumStr()) Failed!")
		}

		return fRat, nil

	}

	newFloat, err := fIa.Numerator.DivideThisBy(&fIa.Denominator, 0, maxPrecision)

	if err != nil {
		return big.NewRat(1, 1), fmt.Errorf("GetRationalValue() - Error returned from fIa.Numerator.DivideThisBy(&fIa.Denominator, 42). Error= %v", err)
	}

	fRat, ok := big.NewRat(1, 1).SetString(newFloat.GetNumStr())

	if !ok {
		return big.NewRat(1, 1), errors.New("GetRationalValue() - big.NewFloat(0).SetString(fracIa.GetNumStr()) Failed!")
	}

	return fRat, nil

}

// GetLowestCommonDenom - Returns a FracIntAry which represents the lowest common
// denominator for the current FracIntAry.
//
// Note: if 'maxPrecision' is less than 0, it is automatically converted to '4,096'
// decimal places.
func (fIa *FracIntAry) GetLowestCommonDenom(maxPrecision int) (FracIntAry, error) {

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	ePrefix := "FracIntAry.GetLowestCommonDenom() "

	ratFrac, err := fIa.GetRationalValue(maxPrecision)

	if err != nil {
		return FracIntAry{},
			fmt.Errorf(ePrefix +
				"Error returned by fIa.GetRationalValue(4096) ")
	}

	newFAry, err := fIa.NewBigInts(ratFrac.Num(), ratFrac.Denom())

	if err != nil {
		return FracIntAry{},
			fmt.Errorf(ePrefix +
				"Error returned by fIa.NewBigInts(ratFrac.Num(), ratFrac.Denom()) ")
	}

	return newFAry, nil
}

// ReduceToLowestCommonDenom - Converts the value of the current FracIntAry
// to its lowest common denominator.
//
// Note: if 'maxPrecision' is less than 0, it is automatically converted to '4,096'
// decimal places.
func (fIa *FracIntAry) ReduceToLowestCommonDenom(maxPrecision int) error {

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	fIaLCD, err := fIa.GetLowestCommonDenom(maxPrecision)

	if err != nil {
		ePrefix := "FracIntAry.ReduceToLowestCommonDenom() "
		return fmt.Errorf(ePrefix+
			"Error returned by fIa.GetLowestCommonDenom(maxPrecision). "+
			"Error='%v' ", err.Error())
	}

	fIa.CopyIn(&fIaLCD)

	return nil
}

type BackUpIntAry struct {
	intAry                 []uint8
	intAryLen              int
	integerLen             int
	significantIntegerLen  int
	significantFractionLen int
	firstDigitIdx          int
	lastDigitIdx           int
	isZeroValue            bool
	isIntegerZeroValue     bool
	precision              int
	signVal                int
	decimalSeparator       rune // https://docs.oracle.com/cd/E19455-01/806-0169/overview-9/index.html
	thousandsSeparator     rune // https://docs.oracle.com/cd/E19455-01/806-0169/overview-9/index.html
	currencySymbol         rune // See Currency Symbol References below:
	// https://gist.github.com/bzerangue/5484121
	// http://symbologic.info/currency.htm
	// http://www.xe.com/symbols.php
}

func (iBa BackUpIntAry) New() BackUpIntAry {
	iAry := BackUpIntAry{}

	iAry.intAry = []uint8{}
	iAry.intAryLen = 0
	iAry.integerLen = 0
	iAry.significantIntegerLen = 0
	iAry.significantFractionLen = 0
	iAry.firstDigitIdx = -1
	iAry.lastDigitIdx = -1
	iAry.isZeroValue = true
	iAry.isIntegerZeroValue = true
	iAry.precision = 0
	iAry.signVal = 1
	iAry.decimalSeparator = '.'
	iAry.thousandsSeparator = ','
	iAry.currencySymbol = '$'

	return iAry
}

func (iBa *BackUpIntAry) Empty() {
	iBa.intAry = []uint8{}
	iBa.intAryLen = 0
	iBa.integerLen = 0
	iBa.significantIntegerLen = 0
	iBa.significantFractionLen = 0
	iBa.firstDigitIdx = -1
	iBa.lastDigitIdx = -1
	iBa.isZeroValue = true
	iBa.isIntegerZeroValue = true
	iBa.precision = 0
	iBa.signVal = 1
	iBa.decimalSeparator = '.'
	iBa.thousandsSeparator = ','
	iBa.currencySymbol = '$'

}

func (iBa *BackUpIntAry) CopyIn(iBa2 *BackUpIntAry) {
	iBa2.SetInternalFlags()
	iBa.Empty()

	iBa.intAry = make([]uint8, iBa2.intAryLen)
	for i := 0; i < iBa2.intAryLen; i++ {
		iBa.intAry[i] = iBa2.intAry[i]
	}

	iBa.intAryLen = iBa2.intAryLen
	iBa.integerLen = iBa2.integerLen
	iBa.significantIntegerLen = iBa2.significantIntegerLen
	iBa.significantFractionLen = iBa2.significantFractionLen
	iBa.firstDigitIdx = iBa2.firstDigitIdx
	iBa.lastDigitIdx = iBa2.lastDigitIdx
	iBa.isZeroValue = iBa2.isZeroValue
	iBa.isIntegerZeroValue = iBa2.isIntegerZeroValue
	iBa.precision = iBa2.precision
	iBa.signVal = iBa2.signVal
	iBa.decimalSeparator = iBa2.decimalSeparator
	iBa.thousandsSeparator = iBa2.thousandsSeparator
	iBa.currencySymbol = iBa2.currencySymbol

}

func (iBa *BackUpIntAry) CopyOut() BackUpIntAry {
	iBa.SetInternalFlags()
	iAry2 := BackUpIntAry{}.New()

	iAry2.intAry = make([]uint8, iBa.intAryLen)

	for i := 0; i < iBa.intAryLen; i++ {
		iAry2.intAry[i] = iBa.intAry[i]
	}

	iAry2.intAryLen = iBa.intAryLen
	iAry2.integerLen = iBa.integerLen
	iAry2.significantIntegerLen = iBa.significantIntegerLen
	iAry2.significantFractionLen = iBa.significantFractionLen
	iAry2.firstDigitIdx = iBa.firstDigitIdx
	iAry2.lastDigitIdx = iBa.lastDigitIdx
	iAry2.isZeroValue = iBa.isZeroValue
	iAry2.isIntegerZeroValue = iBa.isIntegerZeroValue
	iAry2.precision = iBa.precision
	iAry2.signVal = iBa.signVal
	iAry2.decimalSeparator = iBa.decimalSeparator
	iAry2.thousandsSeparator = iBa.thousandsSeparator
	iAry2.currencySymbol = iBa.currencySymbol

	return iAry2
}

func (iBa *BackUpIntAry) Equals(iBa2 *BackUpIntAry) bool {
	iBa.SetInternalFlags()
	iBa2.SetInternalFlags()

	if iBa.intAryLen != iBa2.intAryLen {
		return false
	}

	for i := 0; i < iBa2.intAryLen; i++ {
		if iBa.intAry[i] != iBa2.intAry[i] {
			return false
		}
	}

	if iBa.integerLen != iBa2.integerLen ||
		iBa.significantIntegerLen != iBa2.significantIntegerLen ||
		iBa.significantFractionLen != iBa2.significantFractionLen ||
		iBa.firstDigitIdx != iBa2.firstDigitIdx ||
		iBa.lastDigitIdx != iBa2.lastDigitIdx ||
		iBa.isZeroValue != iBa2.isZeroValue ||
		iBa.isIntegerZeroValue != iBa2.isIntegerZeroValue ||
		iBa.precision != iBa2.precision ||
		iBa.signVal != iBa2.signVal ||
		iBa.decimalSeparator != iBa2.decimalSeparator ||
		iBa.thousandsSeparator != iBa2.thousandsSeparator ||
		iBa.currencySymbol != iBa2.currencySymbol {

		return false
	}

	return true
}

func (iBa *BackUpIntAry) GetIntAryStats() IntAryStatsDto {

	iBa.SetInternalFlags()

	iStats := IntAryStatsDto{}

	iStats.IntAryLen = iBa.intAryLen
	iStats.IntegerLen = iBa.integerLen
	iStats.SignificantIntegerLen = iBa.significantIntegerLen
	iStats.SignificantFractionLen = iBa.significantFractionLen
	iStats.Precision = iBa.precision
	iStats.SignVal = iBa.signVal
	iStats.FirstDigitIdx = iBa.firstDigitIdx
	iStats.LastDigitIdx = iBa.lastDigitIdx
	iStats.IsZeroValue = iBa.isZeroValue
	iStats.IsIntegerZeroValue = iBa.isIntegerZeroValue
	iStats.DecimalSeparator = iBa.decimalSeparator
	iStats.ThousandsSeparator = iBa.thousandsSeparator
	iStats.CurrencySymbol = iBa.currencySymbol

	return iStats

}

func (iBa *BackUpIntAry) GetNumStr() string {

	iBa.SetInternalFlags()

	if iBa.decimalSeparator == 0 {
		iBa.decimalSeparator = '.'
	}

	var buffer bytes.Buffer

	if iBa.signVal < 0 {
		buffer.WriteRune('-')
	}

	intLen := iBa.intAryLen - iBa.precision

	for i := 0; i < intLen; i++ {
		buffer.WriteRune(rune(iBa.intAry[i] + 48))
	}

	if iBa.precision > 0 {
		buffer.WriteRune(iBa.decimalSeparator)

		for j := 0; j < iBa.precision; j++ {
			buffer.WriteRune(rune(iBa.intAry[intLen] + 48))
			intLen++
		}

	}

	return buffer.String()

}

func (iBa *BackUpIntAry) GetPrecision() int {
	return iBa.precision
}

func (iBa *BackUpIntAry) GetSignValue() int {
	return iBa.signVal
}

func (iBa *BackUpIntAry) SetInternalFlags() {

	iBa.intAryLen = len(iBa.intAry)

	if iBa.intAryLen == iBa.precision {
		iBa.intAry = append([]uint8{0}, iBa.intAry...)
		iBa.intAryLen++
	}

	if iBa.intAryLen < iBa.precision {

		deltaZeros := iBa.precision - iBa.intAryLen + 1
		zeroAry := make([]uint8, deltaZeros)
		iBa.intAry = append(zeroAry, iBa.intAry...)
		iBa.intAryLen += deltaZeros
	}

	iBa.firstDigitIdx = -1
	iBa.lastDigitIdx = -1

	lastIntIdx := iBa.intAryLen - iBa.precision - 1
	iBa.isZeroValue = true
	iBa.isIntegerZeroValue = true
	iBa.integerLen = iBa.intAryLen - iBa.precision

	for i := 0; i < iBa.intAryLen; i++ {
		if iBa.intAry[i] > 0 {
			iBa.isZeroValue = false

			if i < iBa.integerLen {
				iBa.isIntegerZeroValue = false
			}
		}
		// At minimum, there should be a single
		// leading zero before the decimal point.
		// Example 0.000.
		if i == lastIntIdx && iBa.intAry[i] == 0 {

			if iBa.firstDigitIdx == -1 {
				iBa.firstDigitIdx = i
			}

		}

		if iBa.intAry[i] > 0 {

			if iBa.firstDigitIdx == -1 {
				iBa.firstDigitIdx = i
			}

			iBa.lastDigitIdx = i
		}

	}

	iBa.significantIntegerLen = iBa.intAryLen - iBa.precision - iBa.firstDigitIdx

	if iBa.lastDigitIdx >= iBa.integerLen {
		iBa.significantFractionLen = iBa.precision - (iBa.lastDigitIdx - iBa.integerLen + 1)
	} else {
		iBa.significantFractionLen = 0
	}

}

func (iBa *BackUpIntAry) SetSignValue(signVal int) {

	if signVal < 0 {
		panic("BackUpIntAry.SetSignValue() - sign value less than zero!")
	}

	iBa.signVal = signVal

}

// IntAryStatsDto - used to transmit
// statistics and critical configuration
// data on an IntAry object.
type IntAryStatsDto struct {
	IntAryLen              int
	IntegerLen             int
	SignificantIntegerLen  int
	SignificantFractionLen int
	Precision              int
	SignVal                int
	FirstDigitIdx          int
	LastDigitIdx           int
	IsZeroValue            bool
	IsIntegerZeroValue     bool
	DecimalSeparator       rune
	ThousandsSeparator     rune
	CurrencySymbol         rune
}

// IntAry - Used to perform string
// based numeric math operations.
//
// Dependencies: NthRootOp - nthroot.go
type IntAry struct {
	intAry []uint8 // Storage is right to left.
	// Least significant digit is
	// stored in intAry[0]. Most
	// significant digit is stored
	// at intAry[intAryLen - 1]
	intAryLen              int
	integerLen             int
	significantIntegerLen  int
	significantFractionLen int
	firstDigitIdx          int
	lastDigitIdx           int
	isZeroValue            bool
	isIntegerZeroValue     bool
	precision              int
	signVal                int
	decimalSeparator       rune
	thousandsSeparator     rune
	currencySymbol         rune
	BackUp                 BackUpIntAry
}

// AddToThis - Adds the value of intAry parameter ia2 to the value
// of the current intAry object.
//
// Parameters:
//
// ia2 *intAry - Incoming intAry object whose value will be subtracted
//
//	from this current intAry value.
func (ia *IntAry) AddToThis(ia2 *IntAry) error {

	IntAryMathAdd{}.RunTotal(ia, ia2)

	return nil

}

// AddIntToThis - Adds an integer number to the value of the
// current IntAry object.
//
// Input Parameters:
//
//		num int 				-	The integer number to be added to the current IntAry object.
//
//	 precision uint 	- The precision which should be applied to the int64 input
//											parameter to designate the number of digits to the right
//											of the decimal point. Example:  num = 123456, precision = 3
//											Result = 123.456 will be added to the current value of
//	                   the current IntAry object. Note: If the value of input parameter
//											'precision' is negative, an error will be returned
//
// Example:
//
//	intDigits     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) AddIntToThis(num int, precision uint) {

	ia2 := IntAry{}.NewInt(num, precision)

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return
}

// AddInt64ToThis - Adds an integer (int64) to the value of the
// current IntAry object.
//
// Input Parameters:
//
//		int64Num int64 -	The integer number to be added to the current IntAry object.
//
//	 precision uint - 	The precision which should be applied to the int64 input
//											parameter to designate the number of digits to the right
//											of the decimal point. Example:  num = 123456, precision = 3
//											Result = 123.456 will be added to the current value of the
//	                   the current IntAry object. If precision is greater than
//	                   2147483647, it will be reduced to this maximum value.
//
// Example:
//
//	                                  Result Added
//	                                  to Current
//	int64Num     precision     	    	 IntAry
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) AddInt64ToThis(int64Num int64, precision uint) {

	precision = ia.validateUintToMaxPrecision(precision)

	ia2 := IntAry{}.NewInt64(int64Num, precision)

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return
}

// AddBigIntToThis - Adds the value of the *big.Int input
// parameter, 'num' to the value of the current IntAry object.
//
// Input Parameters:
//
// num *big.Int 	- The numeric value to be added to the
//
//	value of the current IntAry object. Note
//	that 'num' may be positive or negative.
//
// precision uint	- 'precision' indicates the number of digits
//
//	to be formatted to the right of the decimal
//	place.
//
// Example:
//
//	intDigits     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
//
// Usage:
// num := big.NewInt(123456)
// precision := uint(3)
// err := ia.AddBigIntToThis(num, precision)
//
// The result will equal the current value of the IntAry
// object plus, '123.456'
func (ia *IntAry) AddBigIntToThis(num *big.Int, precision int) error {

	ia2, err := IntAry{}.NewBigInt(num, precision)

	if err != nil {
		return err
	}

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return nil
}

// AddBigIntNumToThis - Adds the value of the current IntAry
// to the BigIntNum input parameter.
func (ia *IntAry) AddBigIntNumToThis(bINum BigIntNum) error {

	ia2, err := IntAry{}.NewBigIntNum(bINum)

	if err != nil {
		ePrefix := "IntAry.AddBigIntNumToThis() "
		return fmt.Errorf(ePrefix+"Error returned by IntAry{}.NewBigIntNum(bINum). "+
			"Error='%v'", err.Error())
	}

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return nil
}

// AddFloat32ToThis - Adds the value of input number (float32)
// to the current value of this IntAry object.
//
// Input Parameters:
//
// num float32 		- The number which will be added to the current
//
//	IntAry object.
//
// precision int	-	Input parameter 'precision' is used
//
//	to set the input precision of 'num'.
//
//	Input parameter 'precision' must be set
//	to a number greater than or equal to zero.
//	It may also be set to a value of -1
//	which causes the number to be formatted to
//	the smallest number of digits to right of
//	the decimal point.
//
// Usage:
// num := float32(123.456000)
// precision := 3 - signals that only the decimals 456 will included as input
// err := ia.AddFloat32ToThis(num, precision)
//
// This usage will result in a value of '123.456'
//
// If precision were set to -1, the resulting value
// would also be '123.456'
//
// Note: if precision is a positive value, and less than the current number
// of digits to the right of the decimal place, rounding may occur.
func (ia *IntAry) AddFloat32ToThis(num float32, precision int) error {

	if precision < -1 {
		return fmt.Errorf("AddFloat32ToThis() Error: Input parameter 'precision' is invalid. 'precision' must be greater than or equal to -1. 'precision'= '%v'", precision)
	}

	ia2, err := IntAry{}.NewFloat32(num, precision)

	if err != nil {
		return nil
	}

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return nil
}

// AddFloat64ToThis - Adds a floating point number (float64) to the
// current value of this IntAry object.
//
// Input Parameters:
//
//	num float64 - This float64 value will be added to the value of the
//								current IntAry. Positive and negative values are
//								accepted.
//
// precision int -	'precision' is applied to the input parameter
//
//	'num' (float64) to determine the number of
//	digits to the right of the decimal point which
//	will be applied to the resulting value.
//
//	Input parameter 'precision' must be set to a
//	number greater than or equal to zero.  It may
//	also be set to a value of -1 which causes the
//	number to be formatted to the smallest number
//	of digits to right of the decimal point.
//
//	Note: if precision is a positive number and less
//	than the current number of digits to the right
//	of the decimal place, rounding may occur.
func (ia *IntAry) AddFloat64ToThis(num float64, precision int) error {

	if precision < -1 {
		return fmt.Errorf("AddFloat64ToThis() Error: Input parameter 'precision' is invalid. 'precision' must be greater than or equal to -1. 'precision'= '%v'", precision)
	}

	ia2, err := IntAry{}.NewFloat64(num, precision)

	if err != nil {
		return err
	}

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return nil
}

// AddFloatBigToThis - Adds a *big.Float number to the current
// value of this IntAry object.
//
// Input Parameters:
//
//	num *big.Float	- This *big.Float value will be added to the value
//										of the current IntAry object. Positive and negative
//										values are accepted.
//
// precision int -	'precision' is applied to the input parameter
//
//	'num' (*big.Float) to determine the number of
//	digits to the right of the decimal point which
//	will be applied to the resulting value.
//
//	Input parameter 'precision' must be set to a
//	number greater than or equal to zero.  It may
//	also be set to a value of -1 which causes the
//	number to be formatted to the smallest number
//	of digits to right of the decimal point.
//
//	Note: if precision is a positive number and less
//	than the current number of digits to the right
//	of the decimal place, rounding may occur.
func (ia *IntAry) AddFloatBigToThis(num *big.Float, precision int) error {

	if precision < -1 {
		return fmt.Errorf("AddFloatBigToThis() Error: Input parameter 'precision' is invalid. 'precision' must be greater than or equal to -1. 'precision'= '%v'", precision)
	}

	ia2, err := IntAry{}.NewFloatBig(num, precision)

	if err != nil {
		return err
	}

	IntAryMathAdd{}.RunTotal(ia, &ia2)

	return nil
}

// AddMultipleToThis - Add the values of  multiple intAry objects to the current
// intAry value.
//
// convertToNumStr - boolean value determines whether the current intAry
//
//	object will convert the intAry value to a number string.
//	Set this parameter to 'false' if this method is called
//	multiple times in order to improve performance.
func (ia *IntAry) AddMultipleToThis(iaMany ...*IntAry) error {

	for _, iAry := range iaMany {

		IntAryMathAdd{}.RunTotal(ia, iAry)

	}

	return nil
}

// AddArrayLengthLeft - Adds leading zeros to the intAry
func (ia *IntAry) AddArrayLengthLeft(addLen int) {

	ia.SetIntAryLength()

	newLen := addLen + ia.intAryLen
	t := make([]uint8, newLen)

	for i := 0; i < newLen; i++ {

		if i < addLen {
			t[i] = 0
		} else {
			t[i] = ia.intAry[i-addLen]
		}

	}

	ia.intAry = t
	ia.SetIntAryLength()
}

// AddArrayLengthRight - Adds trailing zeros
// to the right of the current intAry.
func (ia *IntAry) AddArrayLengthRight(addLen int) {
	ia.SetIntAryLength()

	for i := 0; i < addLen; i++ {
		ia.intAry = append(ia.intAry, 0)
	}

	ia.SetIntAryLength()
}

// AppendToIntAry - appends an integer of
// type uint8 to the internal Integer Array
// of the current IntAry object.
func (ia *IntAry) AppendToIntAry(num uint8) {
	ia.intAry = append(ia.intAry, num)
	ia.intAryLen = len(ia.intAry)
	ia.integerLen = ia.intAryLen - ia.precision
	if num > 0 {
		ia.isZeroValue = false
	}
}

// Ceiling - Returns an IntAry which constitutes
// the mathematical ceiling of the current
// IntAry.
//
// Examples
// ========
//
//							Initial 		 Ceiling
//	 					 Value				Value
//							-------      -------
//	 						5.95					6
//	 						5.05					6
//	 						5							5
//							 -5.05			 	 -5
//	 						2.4				  	3
//	 						2.9					 	3
//							 -2.7				 	 -2
//							 -2					 	 -2
func (ia *IntAry) Ceiling() (IntAry, error) {

	err := ia.IsValid("Ceiling() - ")

	if err != nil {
		return IntAry{}, err
	}

	iAry2 := IntAry{}.New()

	intLen := ia.intAryLen - ia.precision

	intIdx := intLen - 1

	var hasFracDigits bool

	hasFracDigits, err = ia.HasFractionalDigits()

	if err != nil {
		return iAry2, err
	}

	if !hasFracDigits {
		iAry2 = ia.CopyOut()
		return iAry2, nil
	}

	if ia.signVal < 0 && hasFracDigits {

		t := make([]uint8, ia.intAryLen)

		for i := 0; i < intLen; i++ {
			t[i] = ia.intAry[i]
		}

		iAry2.intAry = t[0:]
		iAry2.intAryLen = ia.intAryLen
		iAry2.precision = ia.precision
		iAry2.signVal = ia.signVal
		return iAry2, nil
	}

	t := make([]uint8, ia.intAryLen+1)

	n1 := 0
	n2 := 0
	carry := 0
	adjFac := 1 * ia.signVal
	for i := intIdx; i >= 0; i-- {

		n1 = int(ia.intAry[i])

		if i == intIdx {
			if n1+adjFac < 0 {
				n2 = 10 + n1 + adjFac
				carry = -1
			} else if n1+adjFac > 9 {
				n2 = n1 + adjFac - 10
				carry = 1
			} else {
				n2 = n1 + adjFac
				carry = 0
			}

		} else {

			if n1+carry < 0 {
				n2 = 10 + n1
				carry = -1
			} else if n1+carry > 9 {
				n2 = n1 - 10
				carry = 1
			} else {
				n2 = n1 + carry
				carry = 0
			}

		}

		t[i+1] = uint8(n2)

	}

	if carry != 0 {
		t[0] = uint8(carry)
		iAry2.intAry = t[0 : ia.intAryLen+1]

	} else {

		iAry2.intAry = t[1 : ia.intAryLen+1]
	}

	iAry2.intAryLen = len(iAry2.intAry)
	iAry2.precision = ia.precision
	iAry2.signVal = ia.signVal
	return iAry2, nil
}

// ChangeSign - Changes the sign of the current IntAry instance.
// If the current IntAry numeric value is positive (+), this method
// will change the sign value to negative (-).
//
// Conversely if the current IntAry is a negative (-) numeric value,
// this method will change the sign to positive (+).
func (ia *IntAry) ChangeSign() {

	if ia.IsZero() {
		_ = ia.SetSign(1)
		return
	}

	if ia.signVal < 1 {
		_ = ia.SetSign(1)
	} else {
		_ = ia.SetSign(-1)
	}

}

// CompareSignedValues - Compares two IntAry signed numeric values.
//
// Returns:
// 0  = Current IntAry value is equal to the passed IntAry value.
// 1  = Current IntAry value is greater than the passed IntAry value.
// -1 = Current IntAry value is less than the passed IntAry value.
func (ia *IntAry) CompareSignedValues(iAry2 *IntAry) int {

	iCompare := ia.CompareAbsoluteValues(iAry2)

	if ia.isZeroValue && iAry2.isZeroValue {
		return 0
	}

	if ia.signVal != iAry2.signVal {

		if ia.signVal == 1 {
			return 1
		} else {
			return -1
		}
	}

	// Must be ia.signVal == iAry2.signVal

	if ia.signVal == 1 {
		return iCompare
	}

	// Must be ia.signVal && iAry2.signVal == -1

	return iCompare * -1

}

// CompareAbsoluteValues - Compares the absolute values of
// two IntAry instances.
// Returns:
// 0  = Current IntAry value is equal to the passed IntAry value.
// 1  = Current IntAry value is greater than the passed IntAry value.
// -1 = Current IntAry value is less than the passed IntAry value.
func (ia *IntAry) CompareAbsoluteValues(iAry2 *IntAry) int {

	ia.SetIntAryLength()
	iAry2.SetIntAryLength()

	iAry2.SetIsZeroValue()
	ia.SetIsZeroValue()

	if ia.isZeroValue && iAry2.isZeroValue {
		return 0
	}

	iaIntLen := ia.intAryLen - ia.precision
	iAry2IntLen := iAry2.intAryLen - iAry2.precision

	// Integer Lengths are Equal
	if iaIntLen == iAry2IntLen {
		for i := 0; i < iaIntLen; i++ {
			if ia.intAry[i] > iAry2.intAry[i] {
				return 1
			}

			if iAry2.intAry[i] > ia.intAry[i] {
				return -1
			}
		}
	}

	deltaStartIdx := 0

	// ia Integer Length is Greater than IAry2 Integer Length
	if iaIntLen > iAry2IntLen {
		deltaStartIdx = iaIntLen - iAry2IntLen

		for j := 0; j < iaIntLen; j++ {

			if j < deltaStartIdx {

				if ia.intAry[j] > 0 {
					return 1
				}

			} else {
				// i must be >= deltaStartIdx

				if ia.intAry[j] > iAry2.intAry[j-deltaStartIdx] {
					return 1
				}

				if iAry2.intAry[j-deltaStartIdx] > ia.intAry[j] {
					return -1
				}
			}
		}
	}

	// iAry2 Integer Length is Greater Than ia Integer Length
	if iAry2IntLen > iaIntLen {
		deltaStartIdx = iAry2IntLen - iaIntLen

		for k := 0; k < iAry2IntLen; k++ {

			if k < deltaStartIdx {
				if iAry2.intAry[k] > 0 {
					return -1
				}

			} else {
				// i must be >= deltaStartIdx

				if iAry2.intAry[k] > ia.intAry[k-deltaStartIdx] {
					return -1
				}

				if ia.intAry[k-deltaStartIdx] > iAry2.intAry[k] {
					return 1
				}
			}
		}
	}

	// If precision is zero, the intAry's are equivalent
	if ia.precision == 0 && iAry2.precision == 0 {
		return 0
	}

	// Integer Values are Equivalent. Now test
	// digits to the right of the decimal point.

	// Test fractional digits to right of decimal point
	iaFracIdx := iaIntLen
	iAry2FracIdx := iAry2IntLen
	// Test for case of Equal precision
	if ia.precision == iAry2.precision {
		for m := 0; m < ia.precision; m++ {

			if ia.intAry[iaFracIdx] > iAry2.intAry[iAry2FracIdx] {
				return 1
			}

			if iAry2.intAry[iAry2FracIdx] > ia.intAry[iaFracIdx] {
				return -1
			}

			iaFracIdx++
			iAry2FracIdx++
		}
	}

	iaFracIdx = iaIntLen
	iAry2FracIdx = iAry2IntLen
	// Test for case where ia precision Greater than iAry2 precision
	if ia.precision > iAry2.precision {

		for i := 0; i < ia.precision; i++ {

			if i < iAry2.precision {

				if ia.intAry[iaFracIdx] > iAry2.intAry[iAry2FracIdx] {
					return 1
				}

				if iAry2.intAry[iAry2FracIdx] > ia.intAry[iaFracIdx] {
					return -1
				}

				iaFracIdx++
				iAry2FracIdx++

			} else {
				if ia.intAry[iaFracIdx] > 0 {
					return 1
				}

				iaFracIdx++
			}
		}
	}

	iaFracIdx = iaIntLen
	iAry2FracIdx = iAry2IntLen
	// Test for case where iAry2 precision Greater than ia precision
	if iAry2.precision > ia.precision {

		for i := 0; i < iAry2.precision; i++ {

			if i < ia.precision {

				if ia.intAry[iaFracIdx] > iAry2.intAry[iAry2FracIdx] {
					return 1
				}

				if iAry2.intAry[iAry2FracIdx] > ia.intAry[iaFracIdx] {
					return -1
				}

				iaFracIdx++
				iAry2FracIdx++

			} else {
				if iAry2.intAry[iAry2FracIdx] > 0 {
					return -1
				}

				iAry2FracIdx++
			}
		}

	}

	// The two absolute numeric values must be equal
	return 0

}

func (ia *IntAry) CopyIn(iAry2 *IntAry, copyBackUp bool) {
	iAry2.SetInternalFlags()
	ia.Empty()

	ia.intAry = make([]uint8, iAry2.intAryLen)
	for i := 0; i < iAry2.intAryLen; i++ {
		ia.intAry[i] = iAry2.intAry[i]
	}

	ia.intAryLen = iAry2.intAryLen
	ia.integerLen = iAry2.integerLen
	ia.significantIntegerLen = iAry2.significantIntegerLen
	ia.significantFractionLen = iAry2.significantFractionLen
	ia.firstDigitIdx = iAry2.firstDigitIdx
	ia.lastDigitIdx = iAry2.lastDigitIdx
	ia.isZeroValue = iAry2.isZeroValue
	ia.isIntegerZeroValue = iAry2.isIntegerZeroValue
	ia.precision = iAry2.precision
	ia.signVal = iAry2.signVal
	ia.decimalSeparator = iAry2.decimalSeparator
	ia.thousandsSeparator = iAry2.thousandsSeparator
	ia.currencySymbol = iAry2.currencySymbol

	if copyBackUp {
		ia.BackUp = iAry2.BackUp.CopyOut()
	}
}

// CopyOut - Makes a deep copy of the current IntAry
// instance with backup and returns it as a new IntAry
// object.
func (ia *IntAry) CopyOut() IntAry {
	ia.SetInternalFlags()
	iAry2 := IntAry{}.New()

	iAry2.intAry = make([]uint8, ia.intAryLen)

	for i := 0; i < ia.intAryLen; i++ {
		iAry2.intAry[i] = ia.intAry[i]
	}

	iAry2.intAryLen = ia.intAryLen
	iAry2.integerLen = ia.integerLen
	iAry2.significantIntegerLen = ia.significantIntegerLen
	iAry2.significantFractionLen = ia.significantFractionLen
	iAry2.firstDigitIdx = ia.firstDigitIdx
	iAry2.lastDigitIdx = ia.lastDigitIdx
	iAry2.isZeroValue = ia.isZeroValue
	iAry2.isIntegerZeroValue = ia.isIntegerZeroValue
	iAry2.precision = ia.precision
	iAry2.signVal = ia.signVal
	iAry2.decimalSeparator = ia.decimalSeparator
	iAry2.thousandsSeparator = ia.thousandsSeparator
	iAry2.currencySymbol = ia.currencySymbol

	iAry2.BackUp.CopyIn(&ia.BackUp)

	return iAry2
}

// CopyOutNoBackup
//
// Makes a deep copy of the current IntAry instance
// with NO backup. The method then returns the deep
// copy as a new IntAry object.
func (ia *IntAry) CopyOutNoBackup() IntAry {
	ia.SetInternalFlags()
	iAry2 := IntAry{}.New()

	iAry2.intAry = make([]uint8, ia.intAryLen)

	for i := 0; i < ia.intAryLen; i++ {
		iAry2.intAry[i] = ia.intAry[i]
	}

	iAry2.intAryLen = ia.intAryLen
	iAry2.integerLen = ia.integerLen
	iAry2.significantIntegerLen = ia.significantIntegerLen
	iAry2.significantFractionLen = ia.significantFractionLen
	iAry2.firstDigitIdx = ia.firstDigitIdx
	iAry2.lastDigitIdx = ia.lastDigitIdx
	iAry2.isZeroValue = ia.isZeroValue
	iAry2.isIntegerZeroValue = ia.isIntegerZeroValue
	iAry2.precision = ia.precision
	iAry2.signVal = ia.signVal
	iAry2.decimalSeparator = ia.decimalSeparator
	iAry2.thousandsSeparator = ia.thousandsSeparator
	iAry2.currencySymbol = ia.currencySymbol

	iAry2.BackUp = BackUpIntAry{}.New()

	return iAry2
}

// CopyOutDigits
//
// Makes a deep copy of the current IntAry instance with
// NO backup. This method uses input parameter
// 'digitsToCopy' to copy a specified number of digits to
// the new copy. No backup of the original IntAry is
// returned in this copy.
func (ia *IntAry) CopyOutDigits(digitsToCopy int) IntAry {

	ia.SetInternalFlags()
	iAry2 := IntAry{}.New()

	iAry2.intAry = make([]uint8, digitsToCopy)

	for i := 0; i < digitsToCopy; i++ {
		iAry2.intAry[i] = ia.intAry[i]
	}

	iAry2.intAryLen = digitsToCopy
	if ia.integerLen < digitsToCopy {
		iAry2.precision = digitsToCopy - ia.integerLen
	} else {
		iAry2.precision = 0
	}

	iAry2.signVal = ia.signVal
	iAry2.decimalSeparator = ia.decimalSeparator
	iAry2.thousandsSeparator = ia.thousandsSeparator
	iAry2.currencySymbol = ia.currencySymbol
	iAry2.BackUp = BackUpIntAry{}.New()
	iAry2.SetInternalFlags()

	return iAry2
}

// CopyOutPtr - Similar to CopyOut. In this case, however,
// the method returns a pointer to a deep copy of the current
// IntAry instance.
func (ia *IntAry) CopyOutPtr() *IntAry {

	ia2 := ia.CopyOut()

	return &ia2

}

// CopyToBackUp - Makes a copy of the current
// intAry object and saves it to BackUp.
//
// See ResetFromBackUp() to retrieve the
// last backup copy.
func (ia *IntAry) CopyToBackUp() {

	iBa := BackUpIntAry{}.New()

	iBa.intAry = make([]uint8, ia.intAryLen)

	for i := 0; i < ia.intAryLen; i++ {
		iBa.intAry[i] = ia.intAry[i]
	}

	iBa.intAryLen = ia.intAryLen
	iBa.firstDigitIdx = ia.firstDigitIdx
	iBa.lastDigitIdx = ia.lastDigitIdx
	iBa.isZeroValue = ia.isZeroValue
	iBa.precision = ia.precision
	iBa.signVal = ia.signVal
	iBa.decimalSeparator = ia.decimalSeparator
	iBa.thousandsSeparator = ia.thousandsSeparator
	iBa.currencySymbol = ia.currencySymbol

	ia.BackUp = iBa
}

// DecrementIntegerOne - Decrements the numeric value of the current
// intAry by subtracting '1'.
func (ia *IntAry) DecrementIntegerOne() error {

	ia.SetInternalFlags()

	if ia.isZeroValue || ia.isIntegerZeroValue {
		ia.signVal = -1
	}

	intLen := ia.intAryLen - ia.precision
	intIdx := intLen - 1
	lastIdx := ia.intAryLen - 1

	n1 := 0
	n2 := 0
	carry := 0

	ia.isZeroValue = true
	ia.isIntegerZeroValue = true

	for i := lastIdx; i >= 0; i-- {
		n1 = int(ia.intAry[i])

		if i > intIdx {
			//  i > intIdx
			// This must be a fractional digit
			// Retain fractional digits

			if n1 != 0 {
				ia.isZeroValue = false
			}

			continue

		} else if i == intIdx {

			n2 = n1 + (-1 * ia.signVal)

			if n2 < 0 {
				n2 = n1 + 10 - 1
				carry = 1

			} else if n2 > 9 {
				n2 = n1 + 1 - 10
				carry = 1

			} else {
				carry = 0
			}

		} else {
			// Must be i < intIdx

			n2 = n1 + ((ia.signVal * carry) * -1)

			if n2 < 0 {
				n2 = n1 + 10 - carry
				carry = 1
			} else if n2 > 9 {
				n2 = n1 - 10 + carry
				carry = 1
			} else {
				carry = 0
			}

		}

		if n2 != 0 {
			ia.isZeroValue = false
			ia.isIntegerZeroValue = false
		}

		ia.intAry[i] = uint8(n2)

	}

	if ia.isZeroValue && carry == 0 {
		ia.signVal = 1
	}

	if carry > 0 {

		ia.intAry = append([]uint8{1}, ia.intAry...)
		ia.intAryLen++

	} else if ia.intAry[0] == 0 && intLen > 1 {
		ia.intAry = ia.intAry[1:]
		ia.intAryLen--
	}

	return nil
}

// DivideByTwo
//
// Divides the current value of intAry by 2.
func (ia *IntAry) DivideByTwo() {

	IntAryMathDivide{}.DivideByTwo(ia)

}

// DivideByInt64 - Divide the current value of the intAry
// by an int64 'divisor' parameter passed to the method.
//
// If the quotient has a number of decimal places to the right of
// the decimal point which is greater than 'maxPrecision', the
// result is rounded to 'maxPrecision' decimal places.
//
// If 'maxPrecision' is set equal to -1, 'maxPrecision' is automatically
// set to 4,096.
//
// If 'maxPrecision' is less than -1, an error will be returned.
func (ia *IntAry) DivideByInt64(divisor int64, maxPrecision int) error {

	err := IntAryMathDivide{}.DivideByInt64(ia, divisor, maxPrecision)

	if err != nil {
		ePrefix := "IntAry.DivideByInt64() "
		return fmt.Errorf(ePrefix+
			"Error returned by IntAryMathDivide{}.DivideByInt64() "+
			"Error='%v' \n", err.Error())
	}

	return nil
}

// DivideByTenToPower - Divide the numerical value of the current IntAry
// instance  by 10 raised to the power of the input parameter, 'exponent'.
//
//	ia = ia / (10^exponent)
func (ia *IntAry) DivideByTenToPower(exponent uint) {

	IntAryMathDivide{}.DivideByTenToPower(ia, exponent)
}

// DivideThisBy - Divides the current value of intAry by the parameter iAry2.
// The result of this division is returned as an intAry.
//
// Maximum precision of the division result is controlled by the input
// parameter, 'maxPrecision'.
//
// If 'maxPrecision' is greater than or equal to zero ('0'),
// the number of digits to the right of the decimal place will
// not exceed 'maxPrecision'.
//
// If 'maxPrecision' is set equal to minus one ('-1'), 'maxPrecision'
// will be automatically set to a maximum of 4,096 digits to the right
// of the decimal point.
//
// 'minPrecision' specifies the minimum precision of the final result.
// If 'minPrecision' is less than zero, it is automatically set to zero.
func (ia *IntAry) DivideThisBy(iAry2 *IntAry, minPrecision, maxPrecision int) (IntAry, error) {

	quotient, err := IntAryMathDivide{}.Divide(ia, iAry2, minPrecision, maxPrecision)

	if err != nil {
		ePrefix := "IntAry.DivideThisBy() "
		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error='%v'\n", err.Error())
	}

	return quotient, nil
}

// Empty - Basically resets all the fields of the intAry
// structure to their 'zero' values.
func (ia *IntAry) Empty() {
	ia.intAry = []uint8{}
	ia.intAryLen = 0
	ia.integerLen = 0
	ia.significantIntegerLen = 0
	ia.significantFractionLen = 0
	ia.firstDigitIdx = -1
	ia.lastDigitIdx = -1
	ia.isZeroValue = true
	ia.isIntegerZeroValue = true
	ia.precision = 0
	ia.signVal = 1

	ia.SetNumericSeparatorsToUSADefault()

}

// EmptyBackUp - Deletes the values
// currently stored as backup to the
// current intAry object.
func (ia *IntAry) EmptyBackUp() {
	ia.BackUp = BackUpIntAry{}.New()
}

// Equal - Returns 'true' if all field values of the current
// intAry object are equal to all field values of the input
// parameter intAry object, 'iAry2'.
func (ia *IntAry) Equal(iAry2 IntAry) bool {
	return ia.Equals(&iAry2)
}

// Equals - Returns 'true' if all field values
// of the current intAry object are equal to all
// field values of the input parameter intAry object,
// 'iAry2'.
func (ia *IntAry) Equals(iAry2 *IntAry) bool {
	ia.SetInternalFlags()
	iAry2.SetInternalFlags()

	if ia.intAryLen != iAry2.intAryLen {

		return false
	}

	for i := 0; i < iAry2.intAryLen; i++ {
		if ia.intAry[i] != iAry2.intAry[i] {
			return false
		}
	}

	if ia.integerLen != iAry2.integerLen ||
		ia.significantIntegerLen != iAry2.significantIntegerLen ||
		ia.significantFractionLen != iAry2.significantFractionLen ||
		ia.firstDigitIdx != iAry2.firstDigitIdx ||
		ia.lastDigitIdx != iAry2.lastDigitIdx ||
		ia.isZeroValue != iAry2.isZeroValue ||
		ia.isIntegerZeroValue != iAry2.isIntegerZeroValue ||
		ia.precision != iAry2.precision ||
		ia.signVal != iAry2.signVal ||
		ia.decimalSeparator != iAry2.decimalSeparator ||
		ia.thousandsSeparator != iAry2.thousandsSeparator ||
		ia.currencySymbol != iAry2.currencySymbol {

		return false
	}

	return true
}

// Floor - Math 'Floor' function. Finds the
// integer number which is less than or
// equal to the value of the current intAry.
// Reference Wikipedia,
// https://en.wikipedia.org/wiki/Floor_and_ceiling_functions
//
// Examples
// ========
//
//							Initial 			Floor
//	 					 Value				Value
//							-------      -------
//	 						5.95					5
//	 						5.05					5
//	 						5							5
//							 -5.05			 	 -6
//	 						2.4				  	2
//	 						2.9					 	2
//							 -2.7				 	 -3
//							 -2					 	 -2
func (ia *IntAry) Floor() (IntAry, error) {

	err := ia.IsValid("Floor() - ")

	if err != nil {
		return IntAry{}, err
	}

	iAry2 := IntAry{}.New()

	if ia.isZeroValue {
		iAry2.SetIntAryToZero(uint(ia.precision))
	}

	hasFracDigits, err := ia.HasFractionalDigits()

	if err != nil {
		return iAry2, err
	}

	if !hasFracDigits {
		// There are NO non-zero digits to the
		// right of the decimal place
		iAry2 = ia.CopyOut()
		return iAry2, nil
	}

	intLen := ia.intAryLen - ia.precision
	intIdx := intLen - 1

	if hasFracDigits && ia.signVal > 0 {
		// There ARE non-zero digits to the right of the
		// decimal place
		t := make([]uint8, ia.intAryLen)

		for i := 0; i < intLen; i++ {
			t[i] = ia.intAry[i]
		}

		iAry2.intAry = t[0:]
		iAry2.intAryLen = ia.intAryLen
		iAry2.precision = ia.precision
		iAry2.signVal = ia.signVal
		return iAry2, nil
	}

	// The number has non-zero digits to
	// the right of the decimal place and
	// the number sign is minus (- or ia.signVal = -1)

	t := make([]uint8, ia.intAryLen+1)

	n1 := uint8(0)
	n2 := uint8(0)
	carry := uint8(0)

	for i := intIdx; i >= 0; i-- {

		n1 = ia.intAry[i]

		if i == intIdx {

			if n1+1 > 9 {
				n2 = n1 + 1 - 10
				carry = 1
			} else {
				n2 = n1 + 1
				carry = 0
			}

		} else {

			if n1+carry > 9 {
				n2 = n1 + carry - 10
				carry = 1
			} else {
				n2 = n1 + carry
				carry = 0
			}
		}

		t[i+1] = n2

	}

	if carry != 0 {
		t[0] = carry
		iAry2.intAry = t[0 : ia.intAryLen+1]

	} else {

		iAry2.intAry = t[1 : ia.intAryLen+1]
	}

	iAry2.precision = ia.precision
	iAry2.signVal = ia.signVal
	iAry2.intAryLen = len(iAry2.intAry)
	iAry2.SetIsZeroValue()

	return iAry2, nil
}

// GetAbsoluteValue
//
// Returns an intAry which represents/ the Absolute Value
// of the current intAry.
func (ia *IntAry) GetAbsoluteValue() IntAry {

	absIa := ia.CopyOut()

	absIa.SetAbsoluteValueThis()

	return absIa

}

// GetBigInt - Returns the current value of this intAry object expressed
// as a signed integer number of type *big.Int.
func (ia *IntAry) GetBigInt() (*big.Int, error) {

	ePrefix := "IntAry.GetBigInt() "

	err := ia.IsValid("")

	if err != nil {
		return big.NewInt(0),
			fmt.Errorf(ePrefix+"Error returned by ia.IsValid(). Error='%v'",
				err.Error())
	}

	result := big.NewInt(0).SetInt64(0)
	big10 := big.NewInt(0).SetInt64(10)

	for i := 0; i < ia.intAryLen; i++ {
		result = big.NewInt(0).Mul(result, big10)
		result = big.NewInt(0).Add(result, big.NewInt(0).SetInt64(int64(ia.intAry[i])))

	}

	if ia.signVal == -1 {

		result = big.NewInt(0).Neg(result)
	}

	return result, nil
}

// GetBigIntNum - Converts the numeric value of the current
// IntAry to 'BigIntNum' instance and returns it to the calling
// function.
//
// The returned BigIntNum will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from
// the current IntAry instance.
func (ia *IntAry) GetBigIntNum() (BigIntNum, error) {
	ePrefix := "IntAry.GetBigIntNum() "

	err := ia.IsValid(ePrefix)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by ia.IsValid(). "+
				"Error='%v' ", err.Error())
	}

	numSeps := ia.GetNumericSeparatorsDto()

	bInt, err := ia.GetBigInt()

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by ia.GetBigInt(). "+
				"Error='%v' ", err.Error())
	}

	bIntNum := new(BigIntNum).NewBigInt(bInt, uint(ia.precision))

	err = bIntNum.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by bIntNum.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return bIntNum, nil
}

// GetCurrencySymbol - returns a type 'rune'
// which represents the setting for currency
// symbol in the current IntAry object.
//
// To set the value of currency symbol, see
// the method SetCurrencySymbol().
func (ia *IntAry) GetCurrencySymbol() rune {

	if ia.currencySymbol == 0 {
		ia.currencySymbol = '$'
	}

	return ia.currencySymbol
}

// GetDecimalSeparator - returns a type 'rune'
// which represents the setting for decimal
// separator in the current IntAry object.
//
// A decimal separator separates integer digits
// from fractional digits which the value of
// the current IntAry object is expressed as a
// string.
//
// To set the value of decimal separator, see
// method SetDecimalSeparator().
func (ia *IntAry) GetDecimalSeparator() rune {

	if ia.decimalSeparator == 0 {
		ia.decimalSeparator = '.'
	}

	return ia.decimalSeparator
}

// GetFractionalDigits - Examines the current intAry and
// returns another intAry consisting of the fractional
// digits to the right of the decimal point from the
// current intAry object.
//
// Note: The sign Value of the returned int Ary is always
// positive or +1.
//
// The return intAry will display fractional digits with
// a leading integer digit of zero. Example '0.5678'
func (ia *IntAry) GetFractionalDigits() (IntAry, error) {

	err := ia.IsValid("GetFractionalDigits() - ")

	if err != nil {
		return IntAry{}, err
	}

	iAry2 := IntAry{}.New()

	iAry2.SetIntAryToZero(0)

	if ia.precision == 0 {
		return iAry2, nil
	}

	fracIdx := ia.intAryLen - ia.precision

	iAry2.intAry = make([]uint8, ia.precision+1)
	idx := 1
	for i := fracIdx; i < ia.intAryLen; i++ {
		iAry2.intAry[idx] = ia.intAry[i]
		idx++
	}

	iAry2.precision = ia.precision
	iAry2.signVal = 1
	iAry2.SetInternalFlags()

	return iAry2, nil
}

// GetIntegerDigits - Examines the current intAry object
// and returns a new intAry consisting of only the integer
// digits to the left of the decimal point in the current
// intAry object.
func (ia *IntAry) GetIntegerDigits() (IntAry, error) {

	err := ia.IsValid("GetFractionalDigits() - ")

	if err != nil {
		return IntAry{}, err
	}

	iAry2 := IntAry{}.New()

	if ia.isZeroValue {
		iAry2.SetIntAryToZero(0)
		return iAry2, nil
	}

	intLen := ia.intAryLen - ia.precision

	iAry2.intAry = make([]uint8, intLen)

	for i := 0; i < intLen; i++ {
		iAry2.intAry[i] = ia.intAry[i]
	}

	iAry2.signVal = ia.signVal
	iAry2.precision = 0

	iAry2.SetInternalFlags()

	if iAry2.isZeroValue {
		iAry2.signVal = 1
	}

	return iAry2, nil

}

// GetInt - returns the value of the current
// intAry as a 32-bit integer. If the intAry
// exceeds the maximum or minimum values for
// int, an error will be thrown.
//
// Minimum and Maximum Values for 32-bit Integer
// (int32 & int Types): -2147483648 to 2147483647 .
// Anything out this range will generate an error.
//
// Reference: https://golang.org/ref/spec#Numeric_types
func (ia *IntAry) GetInt() (int, error) {

	ePrefix := "IntAry.GetInt() "

	err := ia.IsValid("")

	if err != nil {
		return 0,
			fmt.Errorf(ePrefix+"Error returned by ia.IsValid(). "+
				"Error='%v' ", err.Error())
	}

	maxInt := big.NewInt(0).SetInt64(int64(math.MaxInt32))

	minInt := big.NewInt(0).SetInt64(int64(math.MinInt32))

	result, _ := ia.GetBigInt()

	compare := result.Cmp(maxInt)

	if compare == 1 {
		return 0,
			errors.New(ePrefix + "Error: the value of this intAry object " +
				"exceeds the maximum allowable value for the int type")
	}

	compare = result.Cmp(minInt)

	if compare == -1 {
		return 0,
			errors.New(ePrefix +
				"Error: the value of this intAry object is less " +
				"than the minimum allowable value for the int type")

	}

	return int(result.Int64()), nil
}

// GetInt64 - Returns the value of the current
// intAry as an int64 value. If the intAry exceeds
// the maximum or minimum values for int64, an error will be
// thrown.
//
// Minimum and Maximum Values for 64-bit Integer (int64 Type):
// -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
// Anything outside this range will generate an error.
//
// Reference: https://golang.org/ref/spec#Numeric_types
func (ia *IntAry) GetInt64() (int64, error) {

	ePrefix := "IntAry.GetInt64() "

	err := ia.IsValid("")

	if err != nil {
		return int64(0),
			fmt.Errorf(ePrefix+"Error returned by ia.IsValid() "+
				"Error='%v' ", err.Error())
	}

	maxI64 := big.NewInt(0).SetInt64(math.MaxInt64)

	minI64 := big.NewInt(0).SetInt64(math.MinInt64)

	result, _ := ia.GetBigInt()

	compare := result.Cmp(maxI64)

	if compare == 1 {
		return int64(0), errors.New("error: the value of this intAry object exceeds the maximum allowable value for the int64 type")
	}

	compare = result.Cmp(minI64)

	if compare == -1 {
		return int64(0), errors.New("error: the value of this intAry object is less than the minimum allowable value for the int64 type")
	}

	return result.Int64(), nil
}

// GetIntAryElement - Returns an element of the
// internal integer array maintained by the
// current IntAry object.
//
// Input Parameter:
// index int - 	The element returned is based on the
//
//	integer index passed to the method. If
//	the index is outside the range of
//	the internal array, an error is returned.
//
// Return Value:	If the 'index' passed to the method is valid,
//
//	the method will return an integer of type
//	uint8 representing the value of the internal
//	integer array at the specified index.
func (ia *IntAry) GetIntAryElement(index int) (uint8, error) {

	if index < 0 || index > (ia.intAryLen-1) {
		return 0,
			fmt.Errorf("Error: GetIntAryElement(index int) - Index is INVALID! index Out Of Array Bounds! "+
				"index= '%v'", index)
	}

	result := ia.intAry[index]

	return result, nil

}

// GetIntAryInt - Returns an element of the
// internal integer array maintained by the
// current IntAry object. The value returned is
// of Type 'int'.
//
// Input Parameter:
// index int - 	The element returned is based on the
//
//	integer index passed to the method. If
//	the index is outside the range of
//	the internal array, an error is returned.
//
// Return Value:	If the 'index' passed to the method is valid,
//
//	the method will return an integer of type
//	'int' representing the value of the internal
//	integer array at the specified index.
func (ia *IntAry) GetIntAryInt(index int) (int, error) {

	if index < 0 || index > (ia.intAryLen-1) {
		return 0,
			fmt.Errorf("Error: GetIntAryInt(index int) - Index is INVALID! index Out Of Array Bounds! "+
				"index= '%v'", index)
	}

	return int(ia.intAry[index]), nil
}

// GetIntAryRune - Returns an element of the
// internal integer array maintained by the
// current IntAry object. The value returned is
// of Type 'rune'.
//
// Input Parameter:
// index int - 	The element returned is based on the
//
//	integer index passed to the method. If
//	the index is outside the range of
//	the internal array, an error is returned.
//
// Return Value:	If the 'index' passed to the method is valid,
//
//	the method will return an integer of type
//	'rune' representing the value of the internal
//	integer array at the specified index.
func (ia *IntAry) GetIntAryRune(index int) (rune, error) {

	if index < 0 || index > (ia.intAryLen-1) {
		return 0,
			fmt.Errorf("Error: GetIntAryRune(index int) - Index is INVALID! index Out Of Array Bounds! "+
				"index= '%v'", index)
	}

	return rune(ia.intAry[index] + 48), nil

}

// GetIntAry - returns a deep copy of the current
// IntAry instance. This method is necessary in
// order to comply with the requirements of the
// INumMgr interface.
//
// The returned IntAry copy will also contain numeric
// separators (decimal separator, thousands separator
// and currency symbol) copied from the current IntAry
// instance.
//
// Before returning the new IntAry instance, this method
// performs a validity test on the current IntAry instance.
func (ia *IntAry) GetIntAry() (IntAry, error) {
	ePrefix := "IntAry.GetIntAry() "

	err := ia.IsValid(ePrefix + "IntAry INVALID! ")

	if err != nil {
		return IntAry{}.New(), err
	}

	return ia.CopyOut(), nil
}

// GetIntAryElements returns the internal integer
// array used by the IntAry object.
//
// ************************************************
// BE CAREFUL!! This returns a reference (pointer)
// to the internal integer array for this IntAry
// object.
// ************************************************
//
// Return parameters:
// []uint8 = an array of unsigned 8-bit integers
//
// int = The length of the []uint8 array returned
//
//	above.
//
// ************************************************
//
//	BE CAREFUL!!!!!!!
//
// ************************************************
// Be careful!! - This returns a slice and therefore
// a reference (i.e. pointer) to the internal array
// for this IntAry object. If you alter this array
// externally, it will alter the internal array
// values.
//
// IMPORTANT: If you alter the returned array, the
// current IntAry object may become invalid. After
// changing the returned array, Call method
// SetInternalFlags() to ensure that the IntAry
// object correctly reflects the changed value.
//
// For a an array which can be safely manipulated,
// see method GetIntAryDeepCopy().
//
// To Append Elements to the internal integer array,
// see method AppendToIntAry()
func (ia *IntAry) GetIntAryElements() ([]uint8, int) {

	ia.SetInternalFlags()

	return ia.intAry, ia.intAryLen
}

// GetIntAryDeepCopy - returns a deep copy of the
// internal integer array maintained by this IntAry
// object. Unlike the array returned by method
// GetIntAryElements(), the array returned by this method
// is not a reference or pointer.
//
// This means that the array returned by this method
// may altered externally without changing the value
// of the original internal array maintained by this
// IntAry object.
func (ia *IntAry) GetIntAryDeepCopy() ([]uint8, int) {
	ia.SetInternalFlags()
	ary := make([]uint8, ia.intAryLen)

	for i := 0; i < ia.intAryLen; i++ {
		ary[i] = ia.intAry[i]
	}

	return ary, ia.intAryLen
}

// GetIntAryLength - returns the length of the
// internal integer array maintained by the current
// IntAry object.
func (ia *IntAry) GetIntAryLength() int {
	ia.SetInternalFlags()
	return ia.intAryLen
}

// GetIntAryStats - returns a series of descriptive
// statistics about the internal Integer Array
// maintained by the current IntAry object.
//
// The method returns an IntAryStatsDto object which
// includes the following information on the current
// IntAry object:
//
// IntAryLen  - Type int: Length of the internal integer array
//
// IntegerLen - Type int: Number of integer digits in the current
//
//	IntAry value.
//
// SignificantIntegerLen - 	Type int: The number of non-zero integer
//
//	digits in the current IntAry value
//
// SignificantFractionLen - Type int: The number of non-zero digits
//
//	to the right of the decimal point in
//	the current IntAry value
//
// precision							- Type int: The number of digits to the right
//
//	of the decimal point in the current
//	IntAry value
//
// SignVal                - Type int: Either +1 or -1 indicating the sign
//
//	of the current IntAry value
//
// FirstDigitIdx          - Type int: The IntAry index of the first
//
//	non-zero digit in the current
//	internal integer array
//
// LastDigitIdx          -  Type int: The IntAry index of the last non-zero
//
//	digit in the internal integer array
//
// IsZero           -	A boolean value indicating whether the
//
//	current IntAry value is zero
//
// IsIntegerZeroValue    - 	A boolean value indicating whether integer
//
//	value to the left of the decimal point is
//	zero
//
// DecimalSeparator      -  Type rune: The character used to separate
//
//	integer and fractional elements of the current
//	IntAry value when presented in string format
//
// ThousandsSeparator    -	Type rune: The character used to separate
//
//	thousands when the current IntAry value is
//	presented in string format
//
// CurrencySymbol        -  Type rune: The character used to designate
//
//	currency when the current IntAry value is
//	presented as a currency string.
func (ia *IntAry) GetIntAryStats() IntAryStatsDto {
	iStats := IntAryStatsDto{}
	ia.SetInternalFlags()

	iStats.IntAryLen = ia.intAryLen
	iStats.IntegerLen = ia.integerLen
	iStats.SignificantIntegerLen = ia.significantIntegerLen
	iStats.SignificantFractionLen = ia.significantFractionLen
	iStats.Precision = ia.precision
	iStats.SignVal = ia.signVal
	iStats.FirstDigitIdx = ia.firstDigitIdx
	iStats.LastDigitIdx = ia.lastDigitIdx
	iStats.IsZeroValue = ia.isZeroValue
	iStats.IsIntegerZeroValue = ia.isIntegerZeroValue
	iStats.DecimalSeparator = ia.decimalSeparator
	iStats.ThousandsSeparator = ia.thousandsSeparator
	iStats.CurrencySymbol = ia.currencySymbol

	return iStats
}

// GetMagnitude - Returns the magnitude of the
// integer portion of the current IntAry numeric
// value. The integer portion of the number is
// represented by the digits to the left of the
// decimal point.
//
// Magnitude is defined here as the power of 10
// which generates a value less than or equal to
// the integer portion of the current IntAry numeric
// value.
//
//	10^magnitude  <= IntAry value
//
// Note: If the current IntAry value is negative,
// an error will be generated.
func (ia *IntAry) GetMagnitude() (int, error) {

	ia.SetInternalFlags()

	if ia.signVal == -1 {
		return -1, fmt.Errorf("IntAry.GetMagnitude() Error: current IntAry value is negative! "+
			"value='%v' ", ia.GetNumStr())
	}

	return ia.intAryLen - ia.precision - ia.firstDigitIdx - 1, nil
}

// GetMagnitudeDigits - Returns the number of digits
// in the integer portion of the current IntAry numeric
// value.
func (ia *IntAry) GetMagnitudeDigits() int {
	ia.SetInternalFlags()
	return ia.intAryLen - ia.precision - ia.firstDigitIdx

}

// GetNumericSeparatorsDto - Returns a structure containing the
// character or rune values for decimal point separator, thousands
// separator and currency symbol.
func (ia *IntAry) GetNumericSeparatorsDto() NumericSeparatorDto {

	numSeps := NumericSeparatorDto{}
	numSeps.DecimalSeparator = ia.GetDecimalSeparator()
	numSeps.ThousandsSeparator = ia.GetThousandsSeparator()
	numSeps.CurrencySymbol = ia.GetCurrencySymbol()

	return numSeps
}

// GetNumStr - returns the current value
// of this intAry object as a number string.
func (ia *IntAry) GetNumStr() string {

	if ia.decimalSeparator == 0 {
		ia.decimalSeparator = '.'
	}

	ia.SetInternalFlags()

	var buffer bytes.Buffer

	if ia.signVal < 0 {
		buffer.WriteRune('-')
	}

	intLen := ia.intAryLen - ia.precision

	for i := 0; i < intLen; i++ {
		buffer.WriteRune(rune(ia.intAry[i] + 48))
	}

	if ia.precision > 0 {
		buffer.WriteRune(ia.decimalSeparator)

		for j := 0; j < ia.precision; j++ {
			buffer.WriteRune(rune(ia.intAry[intLen] + 48))
			intLen++
		}

	}

	return buffer.String()
	// ia.ConvertIntAryToNumStr()
	// return ia.numStrDto
}

// GetNumStrDto - Converts the current IntAry to a NumStrDto
// instance and returns it.
//
// The returned NumStrDto will contain numeric separators
// (decimal separator, thousands separator and currency symbol)
// copied from the current IntAry instance.
//
// Before returning the new NumStrDto instance, this method
// performs a validity check on the current IntAry instance.
func (ia *IntAry) GetNumStrDto() (NumStrDto, error) {

	ePrefix := "IntAry.GetNumStrDto() "

	err := ia.IsValid(ePrefix + " IntAry INVALID! ")

	if err != nil {
		return NumStrDto{}.New(), err
	}

	numSeps := ia.GetNumericSeparatorsDto()

	nDto, err := NumStrDto{}.NewNumStrWithNumSeps(ia.GetNumStr(), numSeps)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by NewNumStr(ia.GetNumStr()) "+
				"Error='%v'", err.Error())
	}

	return nDto, nil
}

// GetNthRootOfThis - Returns an intAry object equal to the 'nth Root'
// of the current intAry.
//
// Input Parameters:
// 'nthRoot' - uint value which specifies the root to calculate
// 'maxPrecision' - uint value which specifies the number of digits
//
//	to the right of the decimal place in the result.
func (ia *IntAry) GetNthRootOfThis(nthRoot, maxPrecision int) (IntAry, error) {

	iaNthRoot := IntAry{}.NewInt(nthRoot, 0)

	nthRt := NthRootOp{}

	return nthRt.GetNthRootIntAry(ia, &iaNthRoot, maxPrecision)
}

// GetPrecision - returns the precision value
// for the current intAry object.
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal point in a string of numeric digits, go
// to the right most digit in the number string and count
// left 'precision' digits.
//
// The value of 'precision' returned by this method will
// always be >= zero (greater than or equal to zero '0').
//
// Example:
//
//					1.234    	GetPrecision() = 3
//							5			GetPrecision() = 0
//				0.12345  		GetPrecision() = 5
//
//	Number String				precision				Fractional Number
//		123456								3								123.456
func (ia *IntAry) GetPrecision() int {
	return ia.precision
}

// GetPrecisionUint - returns the precision value for the
// current intAry object as an unsigned integer (uint).
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal point in a string of numeric digits, go
// to the right most digit in the number string and count
// left 'precision' digits.
//
// The value of 'precision' returned by this method will
// always be >= zero (greater than or equal to zero '0').
//
// Example:
//
//					1.234    	GetPrecision() = 3
//							5			GetPrecision() = 0
//				0.12345  		GetPrecision() = 5
//
//	Number String				precision				Fractional Number
//		123456								3								123.456
func (ia *IntAry) GetPrecisionUint() uint {
	return uint(ia.precision)
}

// GetRuneArray - Returns all the elements of the current
// IntAry as an array of runes.
func (ia *IntAry) GetRuneArray() []rune {

	ia.SetInternalFlags()
	aLen := ia.GetIntAryLength()

	if aLen == 0 {
		return []rune{}
	}

	outRunes := make([]rune, aLen)

	for i := 0; i < aLen; i++ {
		outRunes[i] = rune(ia.intAry[i] + 48)
	}

	return outRunes
}

// GetScaleFactorBigInt - Returns a pointer to a Big Integer
// (*big.Int) which specifies the scale factor associated
// with this IntAry value.
func (ia *IntAry) GetScaleFactor() (*big.Int, error) {
	ia.SetInternalFlags()

	err := ia.IsValid("GetScaleFactorBigInt() - ")

	if err != nil {
		return big.NewInt(0), errors.New("This intAry is invalid")
	}

	if ia.precision == 0 {
		return big.NewInt(int64(1)), nil
	}

	base10 := big.NewInt(0).SetInt64(int64(10))

	bigPrecision := big.NewInt(0).SetInt64(int64(ia.precision))

	scaleFactor := big.NewInt(0).Exp(base10, bigPrecision, nil)

	return scaleFactor, nil

}

// GetSciNotationNumber - Converts the numeric value of the current
// BigIntNum instance into scientific notation and returns this value
// as an instance of type SciNotationNum.
//
// Input Parameter
// ===============
//
// mantissaLen uint	- Specifies the length of the mantissa in the returned
//
//											scientific notation string. If the value of 'mantissaLen'
//											is less than two ('2'), this method will automatically set
//											the 'mantissaLen' to a default value of two ('2').
//
//											Example Scientific Notation:
//											----------------------------
//
//	 										scientific notation string: '2.652e+8'
//
//	 										significand = '2.652'
//	 										significand integer digit = '2'
//												mantissa		= significand factional digits = '.652'
//	 										exponent    = '8'  (10^8)
//
// Note the maximum number of digits which will be retained in the significand is
// 50,000.
func (ia *IntAry) GetSciNotationNumber(mantissaLen uint) (SciNotationNum, error) {

	ePrefix := "IntAry.GetSciNotationNumber() "

	if mantissaLen < 2 {
		mantissaLen = 2
	}

	sciNotationNum := SciNotationNum{}.New()

	err := ia.IsValid(ePrefix + "Current IntAry INVALID! ")

	if err != nil {
		return SciNotationNum{}.New(), err
	}

	if ia.isZeroValue {

		err = sciNotationNum.SetBigIntNumElements(
			new(BigIntNum).NewZero(mantissaLen), new(BigIntNum).NewZero(0))

		if err != nil {
			return SciNotationNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by sciNotationNum.SetBigIntNumElements(...). "+
					"Error='%v'", err.Error())
		}

		return sciNotationNum, nil
	}

	if !ia.isIntegerZeroValue {

		magnitudeInt, _ := ia.GetMagnitude()

		iaNew := ia.CopyOut()

		iaNew.DivideByTenToPower(uint(magnitudeInt))

		iaMagnitude := IntAry{}.NewInt(magnitudeInt, 0)

		if ia.precision > 50000 {
			_ = iaNew.SetPrecision(50000, true)
		}

		_ = sciNotationNum.SetIntAryElements(iaNew, iaMagnitude)

	} else {
		// Must be number with zero integers and fractional digits.
		// Example: 0.256

		iaFracPart, _ := ia.GetFractionalDigits()

		iaFracPart.MultiplyByTenToPower(uint(iaFracPart.precision))

		iaFracPart.OptimizeIntArrayLen(false)

		intMagnitudeFrac, err := iaFracPart.GetMagnitude()

		if err != nil {
			return SciNotationNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by IntAry{}.NewInt(intMagnitudeFrac, 0). "+
					"Error='%v'", err.Error())
		}

		iaFracPart.DivideByTenToPower(uint(intMagnitudeFrac))

		intMagnitudeFrac = intMagnitudeFrac - ia.precision

		iaMagnitude := IntAry{}.NewInt(intMagnitudeFrac, 0)

		err = sciNotationNum.SetIntAryElements(iaFracPart, iaMagnitude)

		if err != nil {
			return SciNotationNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by sciNotationNum.SetIntAryElements("+
					"iaFracPart, iaMagnitude ). Error='%v'",
					err.Error())
		}

	}

	sciNotationNum.SetMantissaLength(mantissaLen)

	return sciNotationNum, nil
}

// GetSciNotationStr - Returns a string expressing the current IntAry
// numerical value as scientific notation.
//
// Input parameter 'mantissaLen' is used to express the number of
// fractional digits displayed in the returned scientific notation
// string. If 'mantissaLen' is less than '2' (two), 'mantissaLen'
// will be automatically set to '2' (two).
//
// Input Parameter
// ===============
//
// mantissaLen uint	- Specifies the length of the mantissa in the returned
//
//											scientific notation string.
//
//											Example Scientific Notation:
//											----------------------------
//
//	 										scientific notation string: '2.652e+8'
//
//	 										significand = '2.652'
//	 										significand integer digit = '2'
//												mantissa		= significand factional digits = '.652'
//	 										exponent    = '8'  (10^8)
func (ia *IntAry) GetSciNotationStr(mantissaLen uint) (string, error) {

	ePrefix := "IntAry.GetSciNotationStr() "

	if mantissaLen < 2 {
		mantissaLen = 2
	}

	sciNotationNum, err := ia.GetSciNotationNumber(mantissaLen)

	if err != nil {
		return "",
			fmt.Errorf(ePrefix + "Error returned by ia.GetSciNotationNumber(mantissaLen) ")
	}

	result, err := sciNotationNum.GetSciNotationStr(mantissaLen)

	if err != nil {
		return "",
			fmt.Errorf(ePrefix + "Error returned by sciNotationNum.GetSciNotationStr(mantissaLen) ")
	}

	return result, nil
}

// GetSign - returns the sign of the current
// intAry object as an integer value.
//
// The sign value returned by this method should
// always be one of two values: +1 or -1 .
func (ia *IntAry) GetSign() int {
	return ia.signVal
}

// GetSquareRootOfThis - Returns an intAry object equal to the 'square root'
// of the current intAry.
//
// Input Parameters:
// 'maxPrecision' - uint value which specifies the number of digits
//
//	to the right of the decimal place in the result.
func (ia *IntAry) GetSquareRootOfThis(maxPrecision int) (IntAry, error) {
	nthRt := NthRootOp{}

	return nthRt.GetSquareRootIntAry(ia, maxPrecision)
}

// GetThisPointer - Returns a pointer to the current IntAry instance
func (ia *IntAry) GetThisPointer() *IntAry {

	return ia
}

// GetThousandsSeparator - returns a value of type 'rune'
// which represents the thousands separator associated
// with the current IntAry object.
//
// The thousands separator is used to separate thousands
// in the integer digits of this IntAry value when that
// value is expressed as a string.
//
// In the US, the thousands separator is typically the comma
// character (','). In this example, the thousands separator
// is a comma: Example - '1,000,000,000'.
//
// To set the value of thousands separator, see the method
// SetThousandsSeparator().
func (ia *IntAry) GetThousandsSeparator() rune {

	if ia.thousandsSeparator == 0 {
		ia.thousandsSeparator = ','
	}

	return ia.thousandsSeparator
}

// HasFractionalDigits - This method examines the
// current intAry object to determine if there
// are non-zero digits to the right of the decimal
// place. If all digits to the right of the decimal
// place are zero, this method returns 'false'
//
// If non-zero digits are present to the right of the
// decimal place, the method returns 'true'.
func (ia *IntAry) HasFractionalDigits() (bool, error) {

	err := ia.IsValid("HasFractionalDigits() - ")

	if err != nil {
		return false, err
	}

	if ia.precision == 0 {
		return false, nil
	}

	ia.SetIntAryLength()

	intLen := ia.intAryLen - ia.precision

	if intLen < 1 {
		return false, fmt.Errorf("HasFractionalDigits() Error - Int Array integer length is less than 1. intLen= '%v'", intLen)
	}

	for i := intLen; i < ia.intAryLen; i++ {
		if ia.intAry[i] > 0 {
			return true, nil
		}
	}

	return false, nil
}

// IncrementIntegerOne - Increment the value of the
// current intAry by adding '1'
func (ia *IntAry) IncrementIntegerOne() error {

	if ia.isZeroValue || ia.isIntegerZeroValue {
		ia.signVal = 1
	}

	intLen := ia.intAryLen - ia.precision
	intIdx := intLen - 1
	lastIdx := ia.intAryLen - 1

	n1 := 0
	n2 := 0
	carry := 0

	ia.isZeroValue = true
	ia.isIntegerZeroValue = true

	for i := lastIdx; i >= 0; i-- {
		n1 = int(ia.intAry[i])

		if i > intIdx {
			//  i > intIdx
			// This must be a fractional digit
			// Retain fractional digits

			if n1 != 0 {
				ia.isZeroValue = false
			}

			continue

		} else if i == intIdx {

			n2 = n1 + (1 * ia.signVal)

			if n2 < 0 {
				n2 = n1 + 10 - 1
				carry = 1

			} else if n2 > 9 {
				n2 = n1 + 1 - 10
				carry = 1

			} else {
				carry = 0
			}

		} else {
			// Must be i < intIdx

			n2 = n1 + ((ia.signVal * carry) * 1)

			if n2 < 0 {
				n2 = n1 + 10 - carry
				carry = 1
			} else if n2 > 9 {
				n2 = n1 - 10 + carry
				carry = 1
			} else {
				carry = 0
			}

		}

		if n2 != 0 {
			ia.isZeroValue = false
			ia.isIntegerZeroValue = false
		}

		ia.intAry[i] = uint8(n2)

	}

	if ia.isZeroValue && carry == 0 {
		ia.signVal = 1
	}

	if carry > 0 {

		ia.intAry = append([]uint8{1}, ia.intAry...)
		ia.intAryLen++

	} else if ia.intAry[0] == 0 && intLen > 1 {
		ia.intAry = ia.intAry[1:]
		ia.intAryLen--
	}

	return nil
}

// IsValid - Examines the current intAry and returns
// an error if the intAry object is found to be invalid.
func (ia *IntAry) IsValid(errName string) error {

	if errName == "" {
		errName = "IntAry.IsValid() "
	}

	ia.SetInternalFlags()

	if ia.signVal != -1 && ia.signVal != 1 {
		return fmt.Errorf("%v sign Value is INVALID! sign Value= '%v'\n", errName, ia.signVal)
	}

	if ia.precision < 0 {
		return fmt.Errorf("%v precision Value is INVALID! sign Value= '%v'\n", errName, ia.precision)
	}

	if ia.precision >= ia.intAryLen {
		return fmt.Errorf("%v error: precision is greater than or equal to IntArray length "+
			"- ia.precision= %v  ia.intAryLen= %v \n", errName, ia.precision, ia.intAryLen)

	}

	if ia.integerLen == 0 {
		return fmt.Errorf("%v error: integer length is zero - missing leading integer zero\n",
			errName)
	}

	return nil
}

// IsEvenNumber - Returns true if the current IntAry numeric value
// is evenly divisible by two (2) with no remainder.
//
// Even Number Definition:
//
//	https://www.mathsisfun.com/definitions/even-number.html
func (ia *IntAry) IsEvenNumber() bool {

	if ia.precision > 0 {
		return false
	}

	if ia.IsZero() {
		return true
	}

	tNum := ia.CopyOut()

	IntAryMathDivide{}.DivideByTwo(&tNum)

	if tNum.precision > 0 {
		return false
	}

	return true
}

// IsOne - Returns 'true' if the value of the current
// IntAry is minus one (-1).
//
// Examples:
// =========
//
// Value 			Result
// -----			------
// -1.0				true
// -1					true
// -1.0000			true
// -1.0001			false
// -2.0				false
func (ia *IntAry) IsMinusOne() bool {

	iaMinusOne := IntAry{}.NewOne(ia.precision)
	iaMinusOne.ChangeSign()

	if ia.Equals(&iaMinusOne) {
		return true
	}

	return false
}

// IsOne - Returns 'true' if the value of the current
// IntAry is '1'.
//
// Examples:
// =========
//
// Value 			Result
// -----			------
// 1.0				true
// 1					true
// 1.0000			true
// 1.0001			false
// 2.0				false
func (ia *IntAry) IsOne() bool {

	iaOne := IntAry{}.NewOne(ia.precision)

	if ia.Equals(&iaOne) {
		return true
	}

	return false
}

// IsZero - Analyzes the current IntAry to determine
// it is a zero value. If the IntAry is equal to a zero
// value, this method returns 'true'
//
// Even Number Definition:
//
//	https://www.mathsisfun.com/definitions/even-number.html
func (ia *IntAry) IsZero() bool {

	err := ia.IsValid("")

	if err != nil {
		return true
	}

	return ia.isZeroValue

}

// Inverse - Returns the inverse of the the current intAry's
// value.
//
// Input Parameter:
//
//		maxPrecision int -	determines the number of digits to the
//	 										right of the decimal point in the result.
//
//												Note: if 'maxPrecision' is set equal to negative
//															one (-1), the maximum number of decimals is
//															set to 4096 digits to the right of the decimal
//															place
func (ia *IntAry) Inverse(maxPrecision int) (IntAry, error) {

	if maxPrecision < 0 {
		maxPrecision = 4096
	}

	if maxPrecision < 0 {
		return IntAry{}.New(), errors.New("InverseOfThis() - Input parameter 'maxPrecision' is INVALID. 'maxPrecision' cannot be less than zero.")
	}

	internalPrecision := maxPrecision + 50

	iaOne := IntAry{}.NewInt(1, 0)

	iaInverse, err := iaOne.DivideThisBy(ia, 0, internalPrecision)

	if err != nil {
		return IntAry{}.New(), fmt.Errorf("InverseOfThis() - Error returned from iaOne.DivideThisBy(ia, maxPrecision). Error= %v", err)
	}

	if iaInverse.GetPrecision() > maxPrecision {
		_ = iaInverse.RoundToPrecision(maxPrecision)
	}

	return iaInverse, nil
}

// MultiplyByTwoToPower Multiply the existing value
// of the IntAry by 2 to the power of the passed in
// parameter.
func (ia *IntAry) MultiplyByTwoToPower(power uint) {

	IntAryMathMultiply{}.MultiplyByTwoToPower(ia, power)

}

// MultiplyByTenToPower - The value of intAry is multiplied
// by 10 to the power of the passed in parameter.
func (ia *IntAry) MultiplyByTenToPower(power uint) {

	IntAryMathMultiply{}.MultiplyByTenToPower(ia, power)

}

// MultiplyThisBy -Multiplies the current IntAry by intAry ia2 and stores the multiplication result
// in the IntAry instance.
//
// Parameters
// ==========
//
// 'ia2' - 				Pointer to an intAry object. In this multiplication operation, 'ia2'
//
//	is the multiplier.
//
// 'minimumResultPrecision' int -
//
//	'minimumResultPrecision' will determine the minimum number of digits computed
//	to the right of the decimal place in the final result.
//	If 'minimumResultPrecision' is set to a value of -1, all significant digits (digits
//	greater than zero) will be returned to the right of the decimal place. Remember
//	that the maximum number of decimal digits returned will be controlled by parameter
//	'maxResultPrecision'
//
// 'maxResultPrecision' 		int -
//
//									'maxResultPrecision' will determine the maximum
//									number of digits to the right of the decimal
//									place in the result.
//
//									Valid values are -1 and values >= zero ('0')
//	       				Values less than -1 will trigger an error.
//
//									A value of -1 signals that no limit will be placed on
//									the number of decimals places to right of the decimal
//									point in the result. Be advised that a very, very large number
//									of decimal digits may be accommodated by the IntAry Type.
func (ia *IntAry) MultiplyThisBy(ia2 *IntAry, minimumPrecision, maxPrecision int) error {

	return IntAryMathMultiply{}.Multiply(ia, ia2, ia, minimumPrecision, maxPrecision)

}

// Multiply - Multiplies intAry ia1 by intAry ia2 and stores the multiplication result
// in intAry iaResult.
//
// Parameters
// ==========
//
// 'ia1' - 				Pointer to an intAry object. In this multiplication operation, 'ia1'
//
//	is the multiplicand.
//
// 'ia2' - 				Pointer to an intAry object. In this multiplication operation, 'ia2'
//
//									is the multiplier.
//
//	 'iaResult' -	Pointer to an intAry object which will be populated with the result
//									of the multiplication operation. The multiplication operation is achieved
//									by multiplying 'ia1' by 'ia2'.
//
// 'minimumResultPrecision' int -
//
//	'minimumResultPrecision' will determine the minimum number of digits computed
//	to the right of the decimal place in the final result.
//	If 'minimumResultPrecision' is set to a value of -1, all significant digits (digits
//	greater than zero) will be returned to the right of the decimal place. Remember
//	that the maximum number of decimal digits returned will be controlled by parameter
//	'maxResultPrecision'
//
// 'maxResultPrecision' 		int -
//
//									'maxResultPrecision' will determine the maximum
//									number of digits to the right of the decimal
//									place in the result.
//
//									Valid values are -1 and values >= zero ('0')
//	       				Values less than -1 will trigger an error.
//
//									A value of -1 signals that no limit will be placed on
//									the number of decimals places to right of the decimal
//									point in the result. Be advised that a very, very large number
//									of decimal digits may be accommodated by the IntAry Type.
func (ia *IntAry) Multiply(ia1, ia2, iaResult *IntAry, minimumResultPrecision, maxResultPrecision int) error {

	return IntAryMathMultiply{}.Multiply(ia1, ia2, iaResult, minimumResultPrecision, maxResultPrecision)
}

// New() - Creates and returns a new blank intAry object.
//
// The returned IntAry instance will contain USA default numeric separators
// (decimal separator, thousands separator and currency symbol).
//
// Usage: ia := intAry{}.New()
func (ia IntAry) New() IntAry {
	iAry := IntAry{}
	iAry.intAry = []uint8{}
	iAry.intAryLen = 0
	iAry.integerLen = 0
	iAry.significantIntegerLen = 0
	iAry.significantFractionLen = 0
	iAry.firstDigitIdx = -1
	iAry.lastDigitIdx = -1
	iAry.isZeroValue = true
	iAry.isIntegerZeroValue = true
	iAry.precision = 0
	iAry.signVal = 1
	iAry.decimalSeparator = '.'
	iAry.thousandsSeparator = ','
	iAry.currencySymbol = '$'
	iAry.BackUp = BackUpIntAry{}.New()

	return iAry
}

// NewWithNumSeps() - Creates and returns a new blank intAry object.
//
// The returned IntAry instance will contain numeric separators (decimal
// separator, thousands separator and currency symbol) as specified
// by input parameter, 'numSeps'.
//
// Usage: ia := intAry{}.New()
func (ia IntAry) NewWithNumSeps(numSeps NumericSeparatorDto) IntAry {

	numSeps.SetDefaultsIfEmpty()

	ia2 := IntAry{}.New()

	_ = ia2.SetNumericSeparatorsDto(numSeps)

	return ia2
}

// NewBigInt - Creates a new intAry object initialized
// to the value of input parameter 'num' which is passed
// as type '*big.Int'. Note that 'num' may be a positive
// or negative number.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. 'precision'
// must be a positive value. Negative 'precision' values will
// trigger an error.
//
// Example:
//
//	intDigits     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
//
// Usage:
// num := big.NewInt(123456)
// precision := uint(3)
// ia, err := intAry{}.NewBigInt(num, precision)
func (ia IntAry) NewBigInt(num *big.Int, precision int) (IntAry, error) {
	ePrefix := "IntAry.NewBigInt() "

	if precision < 0 {
		return IntAry{},
			fmt.Errorf(ePrefix+"Error: Input parameter 'precision' is a negative value! "+
				"precision='%v' ", precision)
	}

	iAry := IntAry{}.New()
	err := iAry.SetIntAryWithBigInt(num, precision)

	if err != nil {
		return IntAry{}, err
	}

	return iAry, nil

}

// NewBigIntNum - Creates a new intAry object initialized
// to the value of input parameter 'bINum' which is passed
// as type 'BigIntNum'.
//
// If the 'BigIntNum' precision value exceeds the positive
// maximum value of 'int' (32-bit integer = 2,147,483,647 ),
// an error will be returned.
//
// Usage:
// bINum, _ := new(BigIntNum).NewNumStr("1234.5678")
// ia, err := intAry{}.NewBigIntNim(bINum)
func (ia IntAry) NewBigIntNum(bINum BigIntNum) (IntAry, error) {

	ePrefix := "IntAry.NewBigIntNum() "

	iAry := IntAry{}.New()

	if bINum.precision > uint(math.MaxInt32) {
		return iAry,
			fmt.Errorf(ePrefix+"Error: Input parameter bINum has a 'precision' value "+
				"which exceeds the MaxInt32 Value. MaxInt32='%v' bINum.precision='%v' ",
				math.MinInt32, bINum.precision)
	}

	err := iAry.SetIntAryWithBigIntNum(bINum)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+"Error returned by iAry.SetIntAryWithBigIntNum(bINum). "+
				"bINum='%v' Error='%v'",
				bINum.GetNumStr(), err.Error())
	}

	return iAry, nil
}

// NewFive - Creates a new IntAry instance with a
// value of '3'.
//
// Note: 'precision' values less than zero will be
// converted to zero.
func (ia IntAry) NewFive(precision int) IntAry {

	ia1 := IntAry{}

	_ = ia1.SetIntAryToFive(precision)

	return ia1
}

// NewFloat32 - Creates a new intAry object initialized
// to the value of input parameter 'num' which is passed
// as type 'float32'. Input parameter 'precision' is used
// to set the input precision of 'num'.
//
// Input parameter 'precision' must be set to a number greater than or equal
// to zero.  It may also be set to a value of -1 which causes the number
// to be formatted to the smallest number of digits to right of the decimal
// point.
//
// Usage:
// num := float32(123.456000)
// precision := 3 - signals that only the decimals 456 will included as input
// ia, err := intAry{}.NewFloat32(num, precision)
//
// The value resulting from the usage example above would be '123.456'.
// If the value of 'precision' was set to -1, the resulting value would
// also be '123.456'
//
// Note: if precision is a positive value, and less than the current number
// of digits to the right of the decimal place, rounding may occur.
func (ia IntAry) NewFloat32(num float32, precision int) (IntAry, error) {

	iAry := IntAry{}.New()

	if precision < -1 {
		return iAry, fmt.Errorf("NewFloat32() Error: 'precision' INVALID. 'precision' must be greater than or equal to -1. precision='%v'", precision)
	}

	err := iAry.SetIntAryWithFloat32(num, precision)

	if err != nil {
		return IntAry{}, err
	}

	return iAry, nil
}

// NewFloat64 - Creates a new intAry object initialized
// to the value of input parameter 'num' which is passed
// as type 'float64'. Input parameter 'precision' is used
// to set the input precision of 'num'.
//
// Input parameter 'precision' must be set to a number greater than or equal
// to zero.  It may also be set to a value of -1 which causes the number
// to be formatted to the smallest number of digits to right of the decimal
// point.
//
// Input parameter 'precision' must be set to a
//
//	number greater than or equal to zero.  It may
//	also be set to a value of -1 which causes the
//	number to be formatted to the smallest number
//	of digits to right of the decimal point.
//
// Usage:
// num := float64(123.456000)
// precision := 3 - signals that only the decimals 456 will included as input
// ia, err := intAry{}.NewFloat64(num, precision)
func (ia IntAry) NewFloat64(num float64, precision int) (IntAry, error) {

	iAry := IntAry{}.New()

	if precision < -1 {
		return iAry, fmt.Errorf("NewFloat64() Error: 'precision' INVALID. 'precision' must be greater than or equal to -1. precision='%v'", precision)
	}

	err := iAry.SetIntAryWithFloat64(num, precision)

	if err != nil {
		return IntAry{}, err
	}

	return iAry, nil
}

// NewFloatBig - Creates a new intAry object initialized
// to the value of input parameter 'num' which is passed
// as type '*big.Float'.
//
// Input Parameters:
//
// num *big.Float -	This floating value will be used to
//
//	to set the value of a new IntAry object
//
// precision int -	'precision' is applied to the input parameter
//
//	'num' (*big.Float) to determine the number of
//	digits to the right of the decimal point which
//	will be applied to the resulting value.
//
//	Input parameter 'precision' must be set to a
//	number greater than or equal to zero.  It may
//	also be set to a value of -1 which causes the
//	number to be formatted to the smallest number
//	of digits to right of the decimal point.
//
// Usage:
// num, err := big.NewFloatBig(123.4560, 3)
//
// Result:
// num = '123.456 - Note: precision parameter may result in rounding.
func (ia IntAry) NewFloatBig(num *big.Float, precision int) (IntAry, error) {

	iAry := IntAry{}.New()

	if precision < -1 {
		return iAry, fmt.Errorf("NewFloatBig() Error: 'precision' INVALID. 'precision' must be greater than or equal to -1. precision='%v'", precision)
	}

	err := iAry.SetIntAryWithFloatBig(num, precision)

	if err != nil {
		return IntAry{}, err
	}

	return iAry, nil
}

// NewInt - Creates a new intAry object initialized to the value
// of input parameter 'intNum' which is passed as type 'int'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the Decimal{}
// syntax thereby allowing Decimal type creation and initialization in
// one step.
//
//					intNum := int(123456)
//					precision := uint(3)
//					dec := Decimal{}.NewInt(intNum, precision)
//	       dec is now equal to 123.456
//
// Examples:
// ---------
//
//	  intNum				precision			Decimal Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewInt(intNum int, precision uint) IntAry {

	iAry := IntAry{}.New()
	iAry.SetIntAryWithInt(intNum, precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewIntExponent - Returns a new IntAry instance. The numeric
// value is set using an int value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int X 10^exponent
//
// Input parameter 'intNum' is of type int.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewIntExponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewIntExponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  intNum		 exponent			  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewIntExponent(intNum int, exponent int) IntAry {

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			intNum *= 10
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithInt(intNum, uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewInt32 - Creates a new intAry object initialized
// to the value of input parameter 'int32Num' which is passed
// as type 'int32'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//					int32Num := int64(123456)
//					precision := uint(3)
//					iAry := IntAry{}.NewInt32(int32Num, precision)
//	       iAry is now equal to 123.456
//
// Examples:
// ---------
//
//	  int32Num			precision			 IntAry Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewInt32(int32Num int32, precision uint) IntAry {

	iAry := IntAry{}.New()
	iAry.SetIntAryWithInt32(int32Num, precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry

}

// NewInt32Exponent - Returns a new IntAry instance. The numeric
// value is set using an int32 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int32 X 10^exponent
//
// Input parameter 'int32Num' is of type int32.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewInt32Exponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewInt32Exponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  int32Num		 exponent			  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewInt32Exponent(int32Num int32, exponent int) IntAry {

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			int32Num *= 10
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithInt32(int32Num, uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewInt64 - Creates a new intAry object initialized
// to the value of input parameter 'int64Num' which is passed
// as type 'int64'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//					int64Num := int64(123456)
//					precision := uint(3)
//					iAry := IntAry{}.NewInt64(int64Num, precision)
//	       iAry is now equal to 123.456
//
// Examples:
// ---------
//
//	  int64Num			precision			 IntAry Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewInt64(int64Num int64, precision uint) IntAry {

	iAry := IntAry{}.New()
	precision = ia.validateUintToMaxPrecision(precision)
	iAry.SetIntAryWithInt64(int64Num, precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewInt64Exponent - Returns a new IntAry instance. The numeric
// value is set using an int64 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int64 X 10^exponent
//
// Input parameter 'int64Num' is of type int64.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewInt64Exponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewInt64Exponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  int64Num		 exponent			  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewInt64Exponent(int64Num int64, exponent int) IntAry {

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			int64Num *= 10
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithInt64(int64Num, uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewIntFracStr - Creates a new IntAry instance based on a numeric value represented
// by separate integer and fractional components.
//
// Input parameters 'intStr' and 'fracStr' are strings representing the integer and
// fractional components. They are combined by this method to create a numeric value
// which is assigned to the new IntAry instance.
//
// Input parameter 'signVal' must be set to one of two values: +1 or -1. This value is
// used to signal the sign of the resulting numeric value. +1 generates a positive number
// and -1 generates a negative number. If input parameters 'inStr' or 'fracStr' contain
// a leading minus or plus sign character, it will be ignored. The sign of the resulting
// numeric value is controlled strictly by input parameter, 'signVal'.
func (ia IntAry) NewIntFracStr(intStr, fracStr string, signVal int) (IntAry, error) {

	ia2 := IntAry{}.NewZero(0)

	ia2.SetNumericSeparatorsToDefaultIfEmpty()

	err := ia2.SetIntAryWithIntFracStr(intStr, fracStr, signVal)

	if err != nil {
		ePrefix := "IntAry.NewIntFracStr() "

		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by ia2.SetIntAryWithIntFracStr(intStr, fracStr, signVal) "+
				"Error='%v' ", err.Error())
	}

	return ia2, nil
}

// NewNumStr - Creates a new intAry object initialized
// to the value of input parameter 'numStr', a string of
// numbers and delimiters which is passed as type 'string'.
//
// This method assumes that the input parameter 'numStr' is a string
// of numeric digits which may be delimited by default USA numeric
// separators. Default USA numeric separators are defined as:
//
//	 	decimal separator = '.'
//	   thousands separator = ','
//			currency symbol = '$'
//
// If the subject 'numStr' employs other national or cultural numeric
// separators, see method IntAry.NewNumStrWithNumSeps(), below.
//
// Usage: ia := intAry{}.NewNumStr("123.456")
func (ia IntAry) NewNumStr(numStr string) (IntAry, error) {

	iAry := IntAry{}.New()
	err := iAry.SetIntAryWithNumStr(numStr)

	if err != nil {
		return IntAry{},
			fmt.Errorf("IntAry.NewNumStr() Error returned by  "+
				"iAry.SetIntAryWithNumStr(numStr). numStr='%v', Error='%v' ",
				numStr, err.Error())
	}

	return iAry, nil

}

// NewNumStrWithNumSeps - Receives a number string as input and returns a
// new IntAry instance. The input parameter 'numSeps' contains numeric
// separators (decimal separator, thousands separator and currency symbol)
// which will be used to parse the number string.
//
// In addition, the numeric separators contained in input parameter 'numSeps'
// will be copied to the returned IntAry instance.
func (ia IntAry) NewNumStrWithNumSeps(
	numStr string,
	numSeps NumericSeparatorDto) (IntAry, error) {

	ePrefix := "IntAry.NewNumStrWithNumSeps() "

	iAry := IntAry{}.New()

	numSeps.SetDefaultsIfEmpty()

	err := iAry.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by  Ary.SetIntAryWithNumStr(numStr). "+
				"Error='%v' ", err.Error())
	}

	err = iAry.SetIntAryWithNumStr(numStr)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by iAry.SetIntAryWithNumStr(numStr). "+
				"numStr='%v', Error='%v' ", numStr, err.Error())
	}

	return iAry, nil
}

// NewNumStrMaxPrecision - Creates a new intAry object initialized
// to the value of input parameter 'num' which is passed
// as type 'string'.
//
// If the number of decimal points to the right of the decimal place
// exceeds input parameter 'maxPrecision', the resulting numeric value
// will be rounded to 'maxPrecision' decimal places.
//
// Usage: ia := intAry{}.NewNumStr("123.456", 3)
func (ia IntAry) NewNumStrMaxPrecision(num string, maxPrecision int) (IntAry, error) {

	iAry := IntAry{}.New()
	err := iAry.SetIntAryWithNumStrMaxPrecision(num, maxPrecision)

	if err != nil {
		return IntAry{},
			fmt.Errorf("IntAry.NewNumStr() Error returned by  "+
				"iAry.SetIntAryWithNumStrMaxPrecision(num). num='%v', Error='%v' ",
				num, err.Error())
	}

	return iAry, nil
}

// NewNumStrDto - Creates, initializes and returns an IntAry
// Type using an input parameter of Type NumStrDto.
func (ia IntAry) NewNumStrDto(numDto NumStrDto) (IntAry, error) {

	ePrefix := "IntAry.NewNumStrDto() "

	err := numDto.IsValid(ePrefix)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by IsValid(ePrefix). "+
				"Error='%v' ", err.Error())
	}

	iAry := IntAry{}.New()

	err = iAry.SetIntAryWithNumStr(numDto.GetNumStr())

	if err != nil {
		return IntAry{},
			fmt.Errorf("IntAry.NewNumStr() Error returned by  "+
				"iAry.SetIntAryWithNumStr(numDto.NumStrOut). "+
				"numDto.NumStrOut='%v', Error='%v' ",
				numDto.GetNumStr(), err.Error())
	}

	return iAry, nil
}

// NewOne - Creates a new IntAry with a value of '1'.
// Note: 'precision' values less than zero will be
// converted to zero.
func (ia IntAry) NewOne(precision int) IntAry {

	ia1 := IntAry{}

	_ = ia1.SetIntAryToOne(precision)

	return ia1
}

// NewPtr - Returns a pointer to a new IntAry instance.
func (ia IntAry) NewPtr() *IntAry {

	ia2 := IntAry{}.New()

	return &ia2
}

// NewTen - Creates a new IntAry instance with a
// value of '10'.
//
// Note: 'precision' values less than zero will be
// converted to zero.
func (ia IntAry) NewTen(precision int) IntAry {

	ia1 := IntAry{}

	_ = ia1.SetIntAryToTen(precision)

	return ia1
}

// NewThree - Creates a new IntAry instance with a
// value of '3'.
//
// Note: 'precision' values less than zero will be
// converted to zero.
func (ia IntAry) NewThree(precision int) IntAry {

	ia1 := IntAry{}

	_ = ia1.SetIntAryToThree(precision)

	return ia1
}

// NewTwo - Creates a new IntAry instance with a
// value of '2'.
//
// Note: 'precision' values less than zero will be
// converted to zero.
func (ia IntAry) NewTwo(precision int) IntAry {

	ia1 := IntAry{}

	_ = ia1.SetIntAryToTwo(precision)

	return ia1
}

// NewUint - Creates a new intAry object initialized to the
// value of input parameter 'uintNum' which is passed as type
// 'uint'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//					uintNum := uint(123456)
//					precision := uint(3)
//					iAry := IntAry{}.NewUint(uintNum, precision)
//	       iAry is now equal to 123.456
//
// Examples:
// ---------
//
//	 uintNum			precision			 IntAry Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewUint(uintNum uint, precision uint) IntAry {

	iAry := IntAry{}.New()
	precision = ia.validateUintToMaxPrecision(precision)
	iAry.SetIntAryWithUint64(uint64(uintNum), precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewUintExponent - Returns a new IntAry instance. The numeric
// value is set using an uint value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = uint X 10^exponent
//
// Input parameter 'uintNum' is of type uint.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewUintExponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewUintExponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	 uintNum		    exponent		  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewUintExponent(uintNum uint, exponent int) IntAry {

	uintTen := uint(10)

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			uintNum *= uintTen
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithUint64(uint64(uintNum), uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewUint32 - Creates a new intAry object initialized to the
// value of input parameter 'uint32Num' which is passed as type
// 'uint32'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//					uint32Num := uint32(123456)
//					precision := uint(3)
//					iAry := IntAry{}.NewUint32(uint32Num, precision)
//	       iAry is now equal to 123.456
//
// Examples:
// ---------
//
//	 uint32Num			precision			 IntAry Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewUint32(uint32Num uint32, precision uint) IntAry {

	iAry := IntAry{}.New()
	precision = ia.validateUintToMaxPrecision(precision)
	iAry.SetIntAryWithUint64(uint64(uint32Num), precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewUint32Exponent - Returns a new IntAry instance. The numeric
// value is set using an uint32 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = uint32 X 10^exponent
//
// Input parameter 'uint32Num' is of type uint32.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewUint32Exponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewUint32Exponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	 uint32Num		  exponent		  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewUint32Exponent(uint32Num uint32, exponent int) IntAry {

	uint32Ten := uint32(10)

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			uint32Num *= uint32Ten
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithUint64(uint64(uint32Num), uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewUint64 - Creates a new intAry object initialized to the
// value of input parameter 'uint64Num' which is passed as type
// 'uint64'.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//					uint64Num := uint64(123456)
//					precision := uint(3)
//					iAry := IntAry{}.NewUint64(uint64Num, precision)
//	       iAry is now equal to 123.456
//
// Examples:
// ---------
//
//	 uint64Num			precision			 IntAry Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (ia IntAry) NewUint64(uint64Num uint64, precision uint) IntAry {

	iAry := IntAry{}.New()
	precision = ia.validateUintToMaxPrecision(precision)
	iAry.SetIntAryWithUint64(uint64Num, precision)
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewUint64Exponent - Returns a new IntAry instance. The numeric
// value is set using an uint64 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = uint64 X 10^exponent
//
// Input parameter 'uint64Num' is of type uint64.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the IntAry{}
// syntax thereby allowing IntAry type creation and initialization in
// one step.
//
//		iAry := IntAry{}.NewUint64Exponent(123456, -3)
//	 -- iAry is now equal to "123.456", precision = 3
//
//		iAry := IntAry{}.NewUint64Exponent(123456, 3)
//	 -- iAry is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	 uint64Num		  exponent		  	IntAry Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (ia IntAry) NewUint64Exponent(uint64Num uint64, exponent int) IntAry {

	uint64Ten := uint64(10)

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			uint64Num *= uint64Ten
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	iAry := IntAry{}.New()
	iAry.SetIntAryWithUint64(uint64Num, uint(exponent))
	_ = iAry.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return iAry
}

// NewZero - Creates a new IntAry instance and sets
// the value to zero.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the
// IntAry{} syntax thereby allowing IntAry type creation and
// initialization in one step.
//
//		iAry := IntAry{}.NewZero(0)
//	 -- iAry is now equal to "0", precision = 0
//
//		iAry := IntAry{}.NewZero(3)
//	 -- iAry is now equal to "0.000", precision = 3
func (ia IntAry) NewZero(precision uint) IntAry {

	ia2 := IntAry{}

	ia2.SetIntAryToZero(precision)

	_ = ia2.SetNumericSeparatorsDto(ia.GetNumericSeparatorsDto())

	return ia2
}

// OptimizeIntArrayLen - Eliminates Leading
// zeros from the front or integer portion
// of the integer string.
//
// If parameter 'optimizeFracDigits' is set
// equal to 'true', trailing zeros to the
// right of the decimal place will also be
// eliminated.
func (ia *IntAry) OptimizeIntArrayLen(optimizeFracDigits bool) {

	ia.SetInternalFlags()

	if ia.isZeroValue {
		return
	}

	integerLen := ia.intAryLen - ia.precision - ia.firstDigitIdx

	if optimizeFracDigits {

		ia.intAry = ia.intAry[ia.firstDigitIdx : ia.lastDigitIdx+1]
		ia.intAryLen = ia.lastDigitIdx - ia.firstDigitIdx + 1
	} else {
		ia.intAry = ia.intAry[ia.firstDigitIdx:]
		ia.intAryLen = ia.intAryLen - ia.firstDigitIdx
	}

	ia.precision = ia.intAryLen - integerLen

}

//	Pow
//
//
//	Raises the value of the current intAry to the power
//	designated by the input parameter 'power'.
//
//	This method calls ia.pwrByTwos() which raises a
//	number to a specified power using the exponentiation
//	by squaring algorithm.
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	power						int
//
//		The input parameter 'power' may be either a
//		positive or negative integer.
//
//	maxResultPrecision			int
//
//		'maxResultPrecision' will determine the maximum
//		number of digits to the right of the decimal Place
//		in the result.
//
//		Valid values are -1 and values >= zero ('0')
//		Values less than -1 will trigger an error.
//
//		A value of -1 signals that no limit will be placed
//		on the number of decimals places to right of the
//		decimal	point in the result.
//
//	internalPrecision			int
//
//		'internalPrecision' will control the number of
//		digits of accuracy to the right of the decimal
//		point maintained by internal multiplication
//		operations used in raising the intAry value to
//		the designated power.
//
//		Valid values are -1 and values >= zero ('0')
//		Values less than -1 will trigger an error.
//
//		A value of -1 signals that no limit will be
//		placed on the number of decimals places to
//		right of the decimal point during internal
//		multiplication operations.
func (ia *IntAry) Pow(power, maxResultPrecision, internalPrecision int) error {

	ePrefix := "IntAry.PowInt() "

	/*
		pwr := big.NewInt(int64(power))
		return ia.pwrByTwos(pwr, maxResultPrecision, internalPrecision)
	*/
	iaPower := IntAry{}.NewInt(power, 0)

	err := IntAryMathPower{}.Pwr(ia, &iaPower, 0, maxResultPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by IntAryMathPower{}.Pwr(). "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// PowThisSquared - Raises the value of the current intAry object
// to a power of '2'.  Essentially the new value of this intAry object
// is equal to the original value squared.
func (ia *IntAry) PowThisSquared() error {

	return ia.MultiplyThisBy(ia, ia.GetPrecision(), -1)
}

// PowByTwos - Raises the value of the current intAry to the power indicated by the parameter,
// 'power'.
//
// Input Parameters:
// =================
// 'power' int -
//
//	The input parameter 'power' may be either
//	a positive or negative integer.
//
// 'maxResultPrecision' int -
//
//					'maxResultPrecision' will determine the maximum
//					number of digits to the right of the decimal
//					place in the result.
//
//					Valid values are -1 and values >= zero ('0')
//	       Values less than -1 will trigger an error.
//
//					A value of -1 signals that no limit will be placed on
//					the number of decimals places to right of the decimal
//					point in the result.
//
//		'internalPrecision' int -
//					'internalPrecision' will control the number of digits of
//					accuracy to the right of the decimal point maintained by
//					internal multiplication operations used in raising the intAry
//					value to the designated power.
//
//					Valid values are -1 and values >= zero ('0')
//	       Values less than -1 will trigger an error.
//
//					A value of -1 signals that no limit will be placed on
//					the number of decimals places to right of the decimal
//					point during internal multiplication operations.
func (ia *IntAry) PowByTwos(power *big.Int, maxResultPrecision, internalPrecision int) error {

	return ia.pwrByTwos(power, maxResultPrecision, internalPrecision)
}

// pwrByTwos - Raises a *big.Int 'base', to the specified 'power'
// using the Exponentiation by squaring algorithm.
// See:
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring#Computation_by_powers_of_2
// This method is based on revised code taken in part from Ye Lin Aung.
// https://stackoverflow.com/questions/30182129/calculating-large-exponentiation-in-golang
// Algorithm modified by Mike Rapp to achieve improved performance.
//
// Input Parameters:
// =================
// 'power' int -
//
//	The input parameter 'power' may be either
//	a positive or negative integer.
//
// 'maxResultPrecision' int -
//
//					'maxResultPrecision' will determine the maximum
//					number of digits to the right of the decimal
//					place in the result.
//
//					Valid values are -1 and values >= zero ('0')
//	       Values less than -1 will trigger an error.
//
//					A value of -1 signals that no limit will be placed on
//					the number of decimals places to right of the decimal
//					point in the result.
//
//		'internalPrecision' int -
//					'internalPrecision' will control the number of digits of
//					accuracy to the right of the decimal point maintained by
//					internal multiplication operations used in raising the intAry
//					value to the designated power.
//
//					Valid values are -1 and values >= zero ('0')
//	       Values less than -1 will trigger an error.
//
//					A value of -1 signals that no limit will be placed on
//					the number of decimals places to right of the decimal
//					point during internal multiplication operations.
func (ia *IntAry) pwrByTwos(power *big.Int, maxResultPrecision, internalPrecision int) error {

	if maxResultPrecision < -1 {
		return fmt.Errorf("Error: Parameter maxResultPrecision is less than -1. maxResultPrecision= %v", maxResultPrecision)
	}

	if internalPrecision < -1 {
		return fmt.Errorf("Error: Parameter internalPrecision is less than -1. internalPrecision= %v", internalPrecision)
	}

	ia.SetInternalFlags()
	tPower := big.NewInt(0).Set(power)
	one := big.NewInt(1)
	zero := big.NewInt(0)
	two := big.NewInt(2)

	if ia.isZeroValue {
		ia.SetIntAryToZero(0)
		return nil
	}

	if tPower.Cmp(two) == -1 {

		if tPower.Cmp(zero) == -1 {

			ia2, err := ia.Inverse(internalPrecision)

			if err != nil {
				return fmt.Errorf("intAry.pwrByTwos = Error returned from ia.Inverse(-1): Error= %v", err)
			}

			ia.CopyIn(&ia2, false)
			tPower = big.NewInt(0).Mul(tPower, big.NewInt(-1))

		} else if tPower.Cmp(one) == 0 {
			// no change in value. x^1 == x
			return nil
		} else if tPower.Cmp(zero) == 0 {
			_ = ia.SetIntAryToOne(0)
			return nil
		}

	}

	tBase := ia.CopyOut()

	_ = ia.SetIntAryToOne(0)

	for tPower.Cmp(zero) == 1 {
		//temp, _:= intAry{}.NewNumStr("0")

		if big.NewInt(0).Mod(tPower, two).Cmp(one) == 0 {
			//temp = big.NewInt(0).Mul(result, tBase)
			//result = big.NewInt(0).Set(temp)
			err := ia.MultiplyThisBy(&tBase, -1, internalPrecision)
			//fmt.Println("ia precision = ", ia.GetPrecision())

			if err != nil {
				return fmt.Errorf("intAry.pwrByTwos() - Error From result.MultiplyThisBy(&tBase, true). Error= %v", err)
			}

			if tPower.Cmp(one) == 0 {
				if maxResultPrecision > -1 && maxResultPrecision < ia.GetPrecision() {
					_ = ia.SetPrecision(maxResultPrecision, true)
				}

				return nil
			}
		}

		err := tBase.MultiplyThisBy(&tBase, -1, internalPrecision)

		if err != nil {
			return fmt.Errorf("intAry.pwrByTwos() - Error From tBase.MultiplyThisBy(&temp, true). Error= %v", err)
		}

		tPower = big.NewInt(0).Div(tPower, two)
	}

	if maxResultPrecision == -1 {
		return nil
	}

	if maxResultPrecision > -1 && maxResultPrecision < ia.GetPrecision() {
		_ = ia.SetPrecision(maxResultPrecision, true)
	}

	return nil
}

// PrefixToIntAry - Adds an integer of type
// uint8 to the front of the existing internal
// integer array maintained by the current
// IntAry object.
func (ia *IntAry) PrefixToIntAry(num uint8) {
	ia.intAry = append([]uint8{num}, ia.intAry...)
	ia.intAryLen = len(ia.intAry)
	ia.integerLen = ia.intAryLen - ia.precision
	if num > 0 {
		ia.isZeroValue = false
	}
}

// ResetFromBackUp - Retrieves data from the
// last saved backup and populates the current
// intAry object.
func (ia *IntAry) ResetFromBackUp() {
	ia.BackUp.SetInternalFlags()

	ia.intAry = make([]uint8, ia.BackUp.intAryLen)
	for i := 0; i < ia.BackUp.intAryLen; i++ {
		ia.intAry[i] = ia.BackUp.intAry[i]
	}

	ia.intAryLen = ia.BackUp.intAryLen
	ia.firstDigitIdx = ia.BackUp.firstDigitIdx
	ia.lastDigitIdx = ia.BackUp.lastDigitIdx
	ia.isZeroValue = ia.BackUp.isZeroValue
	ia.precision = ia.BackUp.precision
	ia.signVal = ia.BackUp.signVal
	ia.decimalSeparator = ia.BackUp.decimalSeparator
	if ia.decimalSeparator == 0 {
		ia.decimalSeparator = '.'
	}

}

// RoundToPrecision - Rounds the value of the current IntAry instance
// to a precision specified by the 'roundToPrecision' parameter.
func (ia *IntAry) RoundToPrecision(roundToPrecision int) error {

	if roundToPrecision < 0 {
		return fmt.Errorf("RoundToPrecision() - Error: roundToPrecision is less than ZERO! roundToPrecision= '%v'", roundToPrecision)
	}

	if ia.precision == 0 {
		return nil
	}

	err := ia.IsValid("RoundToPrecision() - ")

	if err != nil {
		return err
	}

	if ia.isZeroValue {
		ia.SetIntAryToZero(uint(roundToPrecision))
		return nil
	}

	if roundToPrecision == ia.precision {
		return nil
	}

	if roundToPrecision > ia.precision {
		deltaPrecision := roundToPrecision - ia.precision

		for i := 0; i < deltaPrecision; i++ {
			ia.intAry = append(ia.intAry, 0)
		}

		ia.intAryLen = len(ia.intAry)
		ia.precision = roundToPrecision
		ia.SetInternalFlags()
		return nil
	}

	// roundToPrecision must be < ia.precision

	intLen := ia.intAryLen - ia.precision
	newIntAryLen := intLen + roundToPrecision
	fracIdx := intLen

	fracRoundIdx := fracIdx + roundToPrecision

	t := make([]uint8, ia.intAryLen+1)
	n1 := uint8(0)
	n2 := uint8(0)
	carry := uint8(0)

	for i := fracRoundIdx; i >= 0; i-- {

		n1 = ia.intAry[i]

		if i == fracRoundIdx {
			n2 = n1 + 5
		} else {
			n2 = n1 + carry
		}

		if n2 > 9 {
			carry = 1
			n2 = n2 - 10
		} else {
			carry = 0
		}

		t[i+1] = n2
	}

	ia.intAry = []uint8{}

	if carry > 0 {
		t[0] = carry
		ia.intAry = t[0 : newIntAryLen+1]
	} else {
		ia.intAry = t[1 : newIntAryLen+1]
	}
	ia.precision = roundToPrecision
	ia.intAryLen = len(ia.intAry)
	ia.SetInternalFlags()
	return nil
}

// SetAbsoluteValueThis - Converts the current
// value of this intAry object to its
// absolute value.
func (ia *IntAry) SetAbsoluteValueThis() {
	_ = ia.SetSign(1)
}

// SetCurrencySymbol is used to set the value of the currency
// separator character. The currency separator character is used
// to designate currency when formatting the value of this IntAry
// is expressed as a string. In the US, the currency is the dollar
// character ('$'). In the following example, the '$' character
// designates a currency value. Example - '$1,000,000,000.00' .
//
// The SetCurrencySymbol method can be used to change the currency
// separator character in accordance with the customs of countries other
// than the US. For example, in England use of the pound sterling currency
// symbol ('£') would be appropriate.
//
// To examine the current setting for currency symbol, see
// the method GetCurrencySymbol().
//
// Note the default thousands separator character is the comma (',').
func (ia *IntAry) SetCurrencySymbol(currencySymbol rune) {
	ia.currencySymbol = currencySymbol
}

// SetDecimalSeparator - sets the decimal separator character
// (usually a decimal point) which is used to separate
// integer digits from fractional digits when the value
// of this IntAry is expressed as a string.
//
// For US usage, the '.' in the following example is
// the decimal separator. Example: 123.456
//
// SetDecimalSeparator can be used to change the default
// decimal separator to a value other than '.'
//
// To examine the current setting for decimal separator, see
// method GetDecimalSeparator().
//
// Note: The default decimal separator character is '.'
func (ia *IntAry) SetDecimalSeparator(decimalSeparator rune) {
	ia.decimalSeparator = decimalSeparator
}

// SetElement - Sets the value of an IntAry element
// at a given index.
func (ia *IntAry) SetElement(index, val int) error {

	ePrefix := "IntAry.SetElement() "

	if val > math.MaxUint8 {

		return fmt.Errorf(ePrefix+
			"Error: Input parameter 'val' Exceeds Maximum for Unsigned Integer! MaxUint8='%v'",
			math.MaxUint8)

	}

	if val < 0 {
		return errors.New(ePrefix + "Error: Input parameter 'val' is less than ZERO!")
	}

	if index < 0 {
		return errors.New(ePrefix + "Error: Input parameter 'index' is less than ZERO!")
	}

	if index > ia.GetIntAryLength()-1 {
		return errors.New(ePrefix + "Error: Input parameter 'index' EXCEEDS Int Array Length!")
	}

	ia.intAry[index] = uint8(val)

	return nil

}

// SetEqualArrayLengths - Compares an intAry object
// to the current intAry and ensures that the lengths
// of both IntArrays are equal.
func (ia *IntAry) SetEqualArrayLengths(iAry2 *IntAry) {
	iAry2.SetInternalFlags()
	ia.SetInternalFlags()

	iaIntLen := ia.intAryLen - ia.precision
	iAry2IntLen := iAry2.intAryLen - iAry2.precision

	if iaIntLen > iAry2IntLen {
		iAry2.AddArrayLengthLeft(iaIntLen - iAry2IntLen)
	}

	if iAry2IntLen > iaIntLen {
		ia.AddArrayLengthLeft(iAry2IntLen - iaIntLen)
	}

	if ia.precision > iAry2.precision {
		iAry2.AddArrayLengthRight(ia.precision - iAry2.precision)
		iAry2.precision = ia.precision
	}

	if iAry2.precision > ia.precision {
		ia.AddArrayLengthRight(iAry2.precision - ia.precision)
		ia.precision = iAry2.precision
	}

	return
}

// SetIntAryToFive - Sets the value of the intAry object to one ('1').
func (ia *IntAry) SetIntAryToFive(precision int) error {

	if precision < 0 {
		return fmt.Errorf("SetIntAryToFive() - Error: precision is less than ZERO! precision= '%v'", precision)
	}

	ia.intAryLen = 1 + precision
	ia.precision = precision
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.intAry[0] = 5
	ia.signVal = 1
	ia.isZeroValue = false
	ia.isIntegerZeroValue = false
	ia.firstDigitIdx = 0
	ia.lastDigitIdx = 0

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	return nil
}

// SetIntAryToOne - Sets the value of the intAry object to one ('1').
func (ia *IntAry) SetIntAryToOne(precision int) error {

	if precision < 0 {
		return fmt.Errorf("SetIntAryToOne() - Error: precision is less than ZERO! precision= '%v'", precision)
	}

	ia.intAryLen = 1 + precision
	ia.precision = precision
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.intAry[0] = 1
	ia.signVal = 1
	ia.isZeroValue = false
	ia.isIntegerZeroValue = false
	ia.firstDigitIdx = 0
	ia.lastDigitIdx = 0

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	return nil
}

// SetIntAryToTwo - Sets the value of the intAry object to one ('1').
func (ia *IntAry) SetIntAryToTwo(precision int) error {

	if precision < 0 {
		return fmt.Errorf("SetIntAryToTwo() - Error: precision is less than ZERO! precision= '%v'", precision)
	}

	ia.intAryLen = 1 + precision
	ia.precision = precision
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.intAry[0] = 2
	ia.signVal = 1
	ia.isZeroValue = false
	ia.isIntegerZeroValue = false
	ia.firstDigitIdx = 0
	ia.lastDigitIdx = 0

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	return nil
}

// SetIntAryToThree - Sets the value of the intAry object to one ('1').
func (ia *IntAry) SetIntAryToThree(precision int) error {

	if precision < 0 {
		return fmt.Errorf("SetIntAryToThree() - Error: precision is less than ZERO! precision= '%v'", precision)
	}

	ia.intAryLen = 1 + precision
	ia.precision = precision
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.intAry[0] = 3
	ia.signVal = 1
	ia.isZeroValue = false
	ia.isIntegerZeroValue = false
	ia.firstDigitIdx = 0
	ia.lastDigitIdx = 0

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	return nil
}

// SetIntAryToTen - Sets the value of the intAry object to ten ('10')
func (ia *IntAry) SetIntAryToTen(precision int) error {

	if precision < 0 {
		return fmt.Errorf("SetIntAryToTen() - Error: precision is less than ZERO! precision= '%v'", precision)
	}

	ia.intAryLen = 2 + precision
	ia.precision = precision
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.intAry[0] = 1
	ia.signVal = 1
	ia.isZeroValue = false
	ia.isIntegerZeroValue = false
	ia.firstDigitIdx = 0
	ia.lastDigitIdx = 0

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	return nil
}

// SetIntAryToZero - Sets the value of the intAry object to zero ('0').
func (ia *IntAry) SetIntAryToZero(precision uint) {

	precision = ia.validateUintToMaxPrecision(precision)

	ia.intAryLen = 1 + int(precision)
	ia.precision = int(precision)
	ia.intAry = make([]uint8, ia.intAryLen)
	ia.signVal = 1

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	ia.SetInternalFlags()

	return
}

// SetIntAryWithInt - Sets the value of the current intAry object
// to that of the input parameter 'intDigits', an integer of type
// 'int'.
//
// Input parameter 'precision' to indicate the number of digits
// to the right of the decimal place. Input parameter 'precision'
// is of type uint.
//
// The numeric sign (plus or minus) of the resulting intAry value
// is determined by the sign of input parameter,'intDigits'.
// Example:
//
//	intDigits     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) SetIntAryWithInt(intDigits int, precision uint) {
	quotient := 0
	mod := 0

	precision = ia.validateUintToMaxPrecision(precision)

	ia.intAry = []uint8{}
	ia.intAryLen = 0
	ia.precision = int(precision)
	ia.signVal = 1

	if intDigits < 0 {
		intDigits = intDigits * -1
		ia.signVal = -1
	}

	if intDigits == 0 {
		ia.SetIntAryToZero(precision)
		return
	}

	for true {

		if intDigits == 0 {
			break
		}

		quotient = intDigits / 10

		mod = intDigits - (quotient * 10)

		ia.intAry = append(ia.intAry, uint8(mod))
		ia.intAryLen++

		intDigits = quotient

	}

	n1 := uint8(0)
	lastIdx := ia.intAryLen - 1
	totalLen := ia.intAryLen / 2
	for i := 0; i < totalLen; i++ {
		n1 = ia.intAry[i]
		ia.intAry[i] = ia.intAry[lastIdx]
		ia.intAry[lastIdx] = n1
		lastIdx--
	}

	ia.SetInternalFlags()

	return
}

// SetIntAryWithInt32 - Sets the value of the current intAry object
// to that of the input parameter 'intDigits', a 32-bit integer.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// The numeric sign (plus or minus) of the resulting intAry value
// is determined by the sign of input parameter 'int32Num'.

// Example:
//
//	int32Num     precision     	       result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) SetIntAryWithInt32(int32Num int32, precision uint) {

	precision = ia.validateUintToMaxPrecision(precision)

	tenI32 := int32(10)
	quotient := int32(0)
	mod := int32(0)

	ia.intAry = []uint8{}
	ia.intAryLen = 0
	ia.precision = int(precision)
	ia.signVal = 1

	if int32Num < 0 {
		int32Num = int32Num * -1
		ia.signVal = -1
	}

	if int32Num == 0 {
		ia.SetIntAryToZero(precision)
		return
	}

	for true {

		if int32Num == 0 {
			break
		}

		quotient = int32Num / tenI32

		mod = int32Num - (quotient * tenI32)

		ia.intAry = append(ia.intAry, uint8(mod))
		ia.intAryLen++

		int32Num = quotient

	}

	n1 := uint8(0)
	lastIdx := ia.intAryLen - 1
	totalLen := ia.intAryLen / 2
	for i := 0; i < totalLen; i++ {
		n1 = ia.intAry[i]
		ia.intAry[i] = ia.intAry[lastIdx]
		ia.intAry[lastIdx] = n1
		lastIdx--
	}

	ia.SetInternalFlags()

	return
}

// SetIntAryWithInt64 - Sets the value of the current intAry
// object to that of the input parameter 'int64Num', a 64-bit
// integer.
//
// Input parameter 'precision' indicates the number of digits
// to be formatted to the right of the decimal place. Input
// parameter 'precision' is of type uint. The maximum value
// allowed for 'precision' is 2147483645 (the max int32 value
// minus 2). If 'precision' exceeds this maximum value it will
// be reset to that maximum value.
//
// The numeric sign (plus or minus) of the resulting intAry value
// is determined by the sign of input parameter 'int64Num'.
//
// Example:
//
//	int64Num     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) SetIntAryWithInt64(int64Num int64, precision uint) {

	precision = ia.validateUintToMaxPrecision(precision)

	quotient := int64(0)
	mod := int64(0)
	i64Ten := int64(10)

	ia.intAry = []uint8{}
	ia.intAryLen = 0
	ia.precision = int(precision)
	ia.signVal = 1

	if int64Num < 0 {
		int64Num = int64Num * int64(-1)
		ia.signVal = -1
	}

	if int64Num == 0 {
		ia.SetIntAryToZero(precision)
		return
	}

	for true {

		if int64Num == 0 {
			break
		}

		quotient = int64Num / i64Ten

		mod = int64Num - (quotient * i64Ten)

		ia.intAry = append(ia.intAry, uint8(mod))
		ia.intAryLen++

		int64Num = quotient

	}

	n1 := uint8(0)
	lastIdx := ia.intAryLen - 1
	totalLen := ia.intAryLen / 2
	for i := 0; i < totalLen; i++ {
		n1 = ia.intAry[i]
		ia.intAry[i] = ia.intAry[lastIdx]
		ia.intAry[lastIdx] = n1
		lastIdx--
	}

	ia.SetInternalFlags()

	return
}

// SetIntAryWithIntFracStr - Sets the value of the current IntAry instance based on
// a numeric value represented by separate integer and fractional components.
//
// Input parameters 'intStr' and 'fracStr' are strings representing the integer and
// fractional components. They are combined by this method to create a numeric value
// which is assigned to the current IntAry instance.
//
// Input parameter 'signVal' must be set to one of two values: +1 or -1. This value is
// used to signal the sign of the resulting numeric value. +1 generates a positive number
// and -1 generates a negative number. If input parameters 'inStr' or 'fracStr' contain
// a leading minus or plus sign character, it will be ignored. The sign of the resulting
// numeric value is controlled strictly by input parameter, 'signVal'.
func (ia *IntAry) SetIntAryWithIntFracStr(intStr, fracStr string, signVal int) error {

	ePrefix := "IntAry.SetIntFracStrings() "

	cleanIntRuneAry := make([]rune, 0, 100)

	zeroChar := uint8('0')
	nineChar := uint8('9')

	lStr := len(intStr)

	if lStr == 0 {
		return errors.New(ePrefix + "Error: Input Parameter 'intStr' is zero Length!")
	}

	isFirstRune := true

	// Create pure number string from 'intStr'
	for i := 0; i < lStr; i++ {

		if intStr[i] >= zeroChar &&
			intStr[i] <= nineChar {

			if isFirstRune && signVal == -1 {
				cleanIntRuneAry = append(cleanIntRuneAry, '-')
			}

			isFirstRune = false

			cleanIntRuneAry = append(cleanIntRuneAry, rune(intStr[i]))
		}
	}

	if len(cleanIntRuneAry) == 0 {
		cleanIntRuneAry = append(cleanIntRuneAry, '0')
	}

	lStr = len(fracStr)

	if lStr > 0 {

		isFirstRune = true

		for j := 0; j < lStr; j++ {

			if fracStr[j] >= zeroChar &&
				fracStr[j] <= nineChar {

				if isFirstRune {
					cleanIntRuneAry = append(cleanIntRuneAry, ia.GetDecimalSeparator())
					isFirstRune = false
				}

				cleanIntRuneAry = append(cleanIntRuneAry, rune(fracStr[j]))
			}

		}
	}

	err := ia.SetIntAryWithNumStr(string(cleanIntRuneAry))

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bNum.SetNumStr(string(cleanIntRuneAry)). "+
			"cleanIntRuneAry='%v' Error='%v' ", string(cleanIntRuneAry), err.Error())
	}

	return nil
}

// SetIntAryWithUint64 - Sets the value of the current intAry
// object to that of the input parameter 'intDigits', a 64-bit
// unsigned integer.
//
// Note: Input parameter 'precision' to indicate the number of
// digits to the right of the decimal place.
//
// Input parameter, 'signVal' must be set to one of two values: -1 or +1.
// 'signVal' determines the numeric sign of the resulting intAry value,
// either plus or minus.
//
// Example:
//
//	intDigits  precision  signVal		result
//	946254  			3					1				946.254
//	946254				0					1				946254
//	946254  			3					-1			-946.254
//	946254				0					-1			-946254
func (ia *IntAry) SetIntAryWithUint64(intDigits uint64, precision uint) {

	precision = ia.validateUintToMaxPrecision(precision)

	ia.signVal = 1

	if intDigits == 0 {
		ia.SetIntAryToZero(precision)
		return
	}

	ia.precision = int(precision)

	quotient := uint64(0)
	mod := uint64(0)
	ten := uint64(10)

	ia.intAry = []uint8{}
	ia.intAryLen = 0
	for true {

		if intDigits == 0 {
			break
		}

		quotient = intDigits / ten
		mod = intDigits - (quotient * ten)

		ia.intAry = append(ia.intAry, uint8(mod))
		ia.intAryLen++

		intDigits = quotient

	}

	n1 := uint8(0)
	lastIdx := ia.intAryLen - 1
	totalLen := ia.intAryLen / 2
	for i := 0; i < totalLen; i++ {
		n1 = ia.intAry[i]
		ia.intAry[i] = ia.intAry[lastIdx]
		ia.intAry[lastIdx] = n1
		lastIdx--
	}

	ia.SetInternalFlags()

	return
}

// SetIntAryWithBigInt - Sets the current value of the intAry to the value
// of input parameter 'intDigits'. The sign value (plus or minus) is taken
// from the input parameter, 'intDigits'. The precision or number of digits
// to the right of the decimal point, is determined by the input parameter,
// 'precision'. Input parameter 'precision' must be passed as a positive value.
// Negative 'precision' values will trigger an error.
//
// Example:
//
//	intDigits     precision     	    result
//	946254  			   3							   946.254
//	946254				   0							   946254
//	-946254  			   3					      -946.254
//	-946254				   0						    -946254
func (ia *IntAry) SetIntAryWithBigInt(intDigits *big.Int, precision int) error {

	if precision < 0 {
		ePrefix := "IntAry.SetIntAryWithBigInt() "

		return fmt.Errorf(ePrefix+"Error: Input parameter 'precision' is a negative value! "+
			"precision='%v' ", precision)

	}

	bigZero := big.NewInt(0)
	quotient := big.NewInt(0)
	mod := big.NewInt(0)
	big10 := big.NewInt(10)
	modX := big.NewInt(0)

	xIntDigits := big.NewInt(0).Set(intDigits)

	compare := bigZero.Cmp(xIntDigits)

	ia.intAry = []uint8{}
	ia.intAryLen = 0
	ia.precision = precision

	ia.signVal = 1

	if compare == 1 {
		bigMinus1 := big.NewInt(0).SetInt64(int64(-1))
		xIntDigits = big.NewInt(0).Mul(xIntDigits, bigMinus1)
		ia.signVal = -1
	}

	if compare == 0 {
		ia.SetIntAryToZero(uint(precision))
		return nil
	}

	for true {

		compare := bigZero.Cmp(xIntDigits)

		if compare == 0 {
			break
		}

		quotient, mod = big.NewInt(0).QuoRem(xIntDigits, big10, modX)

		ia.intAry = append(ia.intAry, uint8(mod.Int64()))
		ia.intAryLen++

		xIntDigits.Set(quotient)

	}

	n1 := uint8(0)
	lastIdx := ia.intAryLen - 1
	totalLen := ia.intAryLen / 2
	for i := 0; i < totalLen; i++ {
		n1 = ia.intAry[i]
		ia.intAry[i] = ia.intAry[lastIdx]
		ia.intAry[lastIdx] = n1
		lastIdx--
	}

	ia.SetInternalFlags()

	return nil
}

// SetIntAryWithBigIntNum - Sets the current value of the intAry to the value
// of input parameter 'bigINum', a BigIntNum type.
//
// If the input parameter 'BigIntNum' precision value exceeds the positive
// maximum value of 'int' (32-bit integer = 2,147,483,647 ), an error will
// be returned.
//
// The value of Numeric Separators contained in BigINum will NOT be copied
// into the current IntAry instance. IntAry will retain its
// current numeric separator values.
func (ia *IntAry) SetIntAryWithBigIntNum(bigINum BigIntNum) error {

	ePrefix := "IntAry.SetIntAryWithBigIntNum() "

	if bigINum.precision > uint(math.MaxInt32) {
		return fmt.Errorf(ePrefix+"Error: Input parameter bigINum has a 'precision' value "+
			"which exceeds the MaxInt32 Value. MaxInt32='%v' bigINum.precision='%v' ",
			math.MinInt32, bigINum.precision)
	}

	bInt, err := bigINum.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by bigINum.GetBigInt(). Error='%v'",
			err.Error())
	}

	err = ia.SetIntAryWithBigInt(bInt, int(bigINum.GetPrecisionUint()))

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by ia.SetIntAryWithBigInt(bInt, "+
			"bigINum.GetPrecisionUint())).  Error='%v'",
			err.Error())
	}

	return nil
}

// SetIntAryWithFloat32 - Sets the current value of the intAry based on the
// input parameter, 'floatNum'. 'floatNum is of type float32, a 32-bit
// floating point number.
//
// Input parameter 'precision' must be set to a number greater than or equal
// to zero.  It may also be set to a value of -1 which causes the number
// to be formatted to the smallest number of digits to right of the decimal
// point.
//
// Usage:
// num := float32(123.456000)
// precision := 3 - signals that only the decimals 456 will included as input
// err := ia.SetIntAryWithFloat32(num, precision)
//
// The value resulting from the usage example above would be '123.456'.
// If the value of 'precision' was set to -1, the resulting value would
// also be '123.456'
//
// Note: if precision is a positive value, and less than the current number
// of digits to the right of the decimal place, rounding may occur.
func (ia *IntAry) SetIntAryWithFloat32(floatNum float32, precision int) error {

	if precision < -1 {
		return fmt.Errorf("SetIntAryWithFloat32() - Input Parameter 'precision' invalid. 'precision' must be greater than or equal to -1. precison= %v", precision)
	}

	numStr := strconv.FormatFloat(float64(floatNum), 'f', -1, 32)

	err := ia.SetIntAryWithNumStr(numStr)

	if err != nil {
		return fmt.Errorf("SetIntAryWithFloat32() - Error returned from ia.SetIntAryWithNumStr(numStrDto). numStrDto= '%v'  Error= '%v'", numStr, err)
	}

	if precision > -1 {

		_ = ia.SetPrecision(precision, true)

		if err != nil {
			return fmt.Errorf("SetIntAryWithFloat32() - Error returned from ia.SetPrecision(precision, true). precision='%v' Error= '%v'", precision, err)
		}

	}

	return nil
}

// SetIntAryWithFloat32 - Sets the current value of the intAry based on the
// input parameter, 'floatNum'. 'floatNum' is of type float64, a 64-bit
// floating point number.
//
// Input parameter 'precision' must be set to a number greater than or equal
// to zero.  It may also be set to a value of -1 which causes the number
// to be formatted to the smallest number of digits to right of the decimal
// point.
//
// Note: if precision is a positive value, and less than the current number
// of digits to the right of the decimal place, rounding	may occur.
func (ia *IntAry) SetIntAryWithFloat64(floatNum float64, precision int) error {

	if precision < -1 {
		return fmt.Errorf("SetIntAryWithFloat64() Error: Invalid input parameter 'precision'. 'precision' must be greater than or equal to -1. precision= '%v'", precision)
	}

	numStr := strconv.FormatFloat(floatNum, 'f', -1, 64)

	err := ia.SetIntAryWithNumStr(numStr)

	if err != nil {
		return fmt.Errorf("SetIntAryWithFloat64() - ia.SetIntAryWithNumStr(numStrDto) returned error. numStrDto='%v' Error='%v'", numStr, err)
	}

	if precision > -1 {
		err := ia.SetPrecision(precision, true)

		if err != nil {
			return fmt.Errorf("SetIntAryWithFloat64() - Error returned from ia.SetPrecision(precision, true). precision='%v' Error= '%v'", precision, err)
		}

	}

	return nil
}

// SetIntAryWithFloatBig - Sets the current value of the intAry based on the
// input parameter, 'floatNum'.  'floatNum' is of type *big.Float.
//
// Input Parameters:
//
// floatNum *big.Float - 	The floating point value to which the current IntAry object
//
//	will be set.
//
// precision int 			 - 	The number of digits to the right of the decimal point
//
//	which will be used to set the value of the current IntAry
//	object
//
//	Input parameter 'precision' must be set to a number greater
//	than or equal to zero.  It may also be set to a value of -1
//	which causes the number to be formatted to the smallest number
//	of digits to right of the decimal point.
//
//	Note: if precision is a positive value, and less than the current
//	number of digits to the right of the decimal place, rounding may occur.
func (ia *IntAry) SetIntAryWithFloatBig(floatNum *big.Float, precision int) error {

	if precision < -1 {
		return fmt.Errorf("SetIntAryWithFloatBig() Error: Invalid input parameter 'precision'. 'precision' must be greater than -1. precision= '%v'", precision)
	}

	var numStr string

	if precision == -1 {
		numStr = floatNum.Text('f', precision)
	} else {
		numStr = floatNum.Text('f', precision+1)
	}

	err := ia.SetIntAryWithNumStr(numStr)

	if err != nil {
		return fmt.Errorf("SetIntAryWithFloatBig() - Error returned from ia.SetIntAryWithNumStr(numStrDto). numStrDto='%v'  Error= '%v'", numStr, err)
	}

	if precision > -1 {
		err := ia.SetPrecision(precision, true)

		if err != nil {
			return fmt.Errorf("SetIntAryWithFloatBig() - Error returned from ia.SetPrecision(precision, true). precision='%v' Error= '%v'", precision, err)
		}

	}

	return nil
}

// SetIntAryWithIntAry - Sets the value of the current intAry based on an array of
// integers i.e. []int.
//
// Input parameter 'precision' will determine the number of digits to the right
// of the decimal place.
//
// Input parameter 'signVal' must be either +1 or -1 indicating the sign of the
// number represented by the integer array. Note: If signVal is not equal to +1 or -1,
// an error is generated.
func (ia *IntAry) SetIntAryWithIntAry(iAry2 []int, precision uint, signVal int) error {

	if signVal != 1 && signVal != -1 {
		return fmt.Errorf("SetIntAryWithIntAry() - Error: signVal parameter is INVALID! signVal must be -1 or +1. signVal='%v'", signVal)
	}

	lIAry2 := len(iAry2)
	ia.intAry = make([]uint8, lIAry2)
	for i := 0; i < lIAry2; i++ {
		ia.intAry[i] = uint8(iAry2[i])
	}

	ia.intAryLen = lIAry2
	ia.precision = int(precision)
	ia.signVal = signVal
	ia.SetInternalFlags()
	return nil
}

// SetIntAryWithUint8Ary - is designed to set the value of the current IntAry
// object by passing in a pointer to an unsigned integer array ([]uint8). []uint8
// is the type of array native to the IntAry object.
//
// Input parameter 'precision' will determine the number of digits to the right
// of the decimal place.
//
// Input parameter 'signVal' must be either +1 or -1 indicating the sign of the
// number represented by the integer array. Note: If signVal is not equal to +1 or -1,
// an error is generated.
func (ia *IntAry) SetIntAryWithUint8Ary(iAry2 []uint8, precision uint, signVal int) error {

	if signVal != 1 && signVal != -1 {
		return fmt.Errorf("SetIntAryWithUint8Ary() - Error: signVal parameter is INVALID! signVal must be -1 or +1. signVal='%v'", signVal)
	}

	lIAry2 := len(iAry2)
	ia.intAry = make([]uint8, lIAry2)

	for i := 0; i < lIAry2; i++ {
		ia.intAry[i] = uint8(iAry2[i])
	}

	ia.intAryLen = lIAry2
	ia.precision = int(precision)
	ia.signVal = signVal
	ia.SetInternalFlags()

	if ia.isIntegerZeroValue && ia.integerLen > 1 {
		ia.OptimizeIntArrayLen(false)
	}

	return nil
}

func (ia *IntAry) SetIntAryWithIntAryObj(iAry2 *IntAry, copyBackup bool) error {

	err := iAry2.IsValid("SetIntAryWithIntAryObj()")

	if err != nil {
		return fmt.Errorf("SetIntAryWithIntAryObj - Input parameter iAry2 is INVALID! Error= %v", err)
	}

	ia.CopyIn(iAry2, copyBackup)

	return nil
}

// SetIntAryWithNumStrMaxPrecision - receives a raw number string and sets the
// fields of the internal intAry structure to the appropriate
// values.
//
// A second input parameter specifies the maximum allowable precision for the IntAry.
// If IntAry precision exceeds 'maxPrecision', IntAry precision is rounded to
// 'maxPrecision'.
//
// The term 'precision' defines the number of numeric digits to the right of the
// decimal point or decimal separator.
func (ia *IntAry) SetIntAryWithNumStrMaxPrecision(str string, maxPrecision int) error {
	ePrefix := "IntAry.SetIntAryWithNumStrMaxPrecision() "

	err := ia.SetIntAryWithNumStr(str)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by ia.SetIntAryWithNumStr(str)"+
			"str='%v' Error='%v' ",
			str, err.Error())
	}

	if ia.precision > maxPrecision {
		_ = ia.RoundToPrecision(maxPrecision)
	}

	return nil
}

// SetIntAryWithNumStr - receives a raw number string and sets the
// fields of the internal intAry structure to the appropriate
// values.
func (ia *IntAry) SetIntAryWithNumStr(str string) error {

	ePrefix := "IntAry.SetIntAryWithNumStr() "

	if len(str) == 0 {
		return errors.New(ePrefix + "Error: received zero length number string")
	}

	numSeps := ia.GetNumericSeparatorsDto()

	ia.Empty()

	err := ia.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by ia.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v' \n", err.Error())
	}

	ia.signVal = 1
	baseRunes := []rune(str)
	lBaseRunes := len(baseRunes)
	isStartRunes := false
	isEndRunes := false
	isFractionalValue := false

	for i := 0; i < lBaseRunes && isEndRunes == false; i++ {

		if baseRunes[i] == '+' ||
			baseRunes[i] == ' ' ||
			baseRunes[i] == ia.thousandsSeparator ||
			baseRunes[i] == ia.currencySymbol {

			continue

		}

		if baseRunes[i] == ',' && ia.decimalSeparator != ',' {
			continue
		}

		if isStartRunes == true &&
			isEndRunes == false &&
			isFractionalValue &&
			baseRunes[i] == ia.decimalSeparator {

			continue
		}

		if baseRunes[i] == '-' &&
			isStartRunes == false && isEndRunes == false &&
			i+1 < lBaseRunes &&
			((baseRunes[i+1] >= '0' && baseRunes[i+1] <= '9') ||
				baseRunes[i+1] == ia.decimalSeparator) {

			ia.signVal = -1
			isStartRunes = true
			continue

		} else if isEndRunes == false &&
			baseRunes[i] >= '0' && baseRunes[i] <= '9' {

			ia.intAry = append(ia.intAry, uint8(baseRunes[i]-48))
			isStartRunes = true

			if isFractionalValue {
				ia.precision++
			}

		} else if isEndRunes == false &&
			i+1 < lBaseRunes &&
			baseRunes[i+1] >= '0' && baseRunes[i+1] <= '9' &&
			baseRunes[i] == ia.decimalSeparator {

			isFractionalValue = true
			continue

		} else if isStartRunes && !isEndRunes {

			isEndRunes = true

		}

	}

	ia.SetSignificantDigitIdxs()

	if ia.intAryLen == 0 || ia.isZeroValue {
		ia.SetIntAryToZero(uint(ia.precision))
		return nil
	}

	// Validate intAry object
	err = ia.IsValid(ePrefix + "- ")

	if err != nil {
		return err
	}

	return nil
}

// SetIntAryWithNumStrDto - Sets the numeric value of the current
// IntAry instance to that of input parameter 'nDto' which is of
// Type NumStrDto.
//
// Note that the Numeric Separator values are NOT copied into the
// current IntAry instance. The current IntAry numeric separators
// remain unchanged.
func (ia *IntAry) SetIntAryWithNumStrDto(nDto NumStrDto) error {

	ePrefix := "IntAry.SetIntAryWithNumStrDto() "

	err := nDto.IsValid(ePrefix + "'nDto' Invalid! ")

	if err != nil {
		return err
	}

	err = ia.SetIntAryWithNumStr(nDto.GetNumStr())

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by ia.SetIntAryWithNumStr(nDto.GetNumStr()). "+
			"Error='%v' \n", err.Error())
	}

	return nil
}

// SetIntAryLength - Calculates the current IntAry string length
// and sets internal variable 'ia.intAryLen'.
func (ia *IntAry) SetIntAryLength() {
	ia.intAryLen = len(ia.intAry)
}

// SetIsZeroValue - Analyzes the value
// of the intAry and sets a flag
// if the value of intAry evaluates
// to zero.
func (ia *IntAry) SetIsZeroValue() {
	ia.intAryLen = len(ia.intAry)

	ia.isZeroValue = true
	ia.isIntegerZeroValue = true

	intLen := ia.intAryLen - ia.precision

	for i := 0; i < ia.intAryLen; i++ {

		if i < intLen && ia.intAry[i] > 0 {
			ia.isIntegerZeroValue = false
		}

		if ia.intAry[i] > 0 {
			ia.isZeroValue = false
			return
		}
	}

	// ia.isZeroValue == true
	// signVal must be 1
	ia.signVal = 1
}

// SetInternalFlags - Sets Array Lengths and
// test for zero values
func (ia *IntAry) SetInternalFlags() {
	ia.SetSignificantDigitIdxs()
}

// SetNumericSeparators - Used to assign values for the Decimal and Thousands separators
// as well as the Currency Symbol to be used in displaying the current number string.
//
// Different nations and cultures use different symbols to delimit numerical values. In the
// USA and many other countries, a period character ('.') is used to delimit integer and
// fractional digits within a numeric value (123.45). Likewise, thousands may be delimited
// by a comma (','). Currency signs very by nationality. For instance, the USA, Canada and
// several other countries use the dollar sign ($) as a currency symbol.
//
// For a list of major world currency symbols see:
//
//		MikeAustin71\mathopsgo\mathops\mathopsconstants.go
//	 http://www.xe.com/symbols.php
//
// Note: If zero values are submitted as input for separator values, those values will default
// to USA standards.
//
// USA Examples:
//
// Decimal Separator period ('.') 		= 123.456
// Thousands Separator comma (',') 		= 1,000,000,000
// Currency Symbol dollar sign ('$')	= $123
func (ia *IntAry) SetNumericSeparators(
	decimalSeparator,
	thousandsSeparator,
	currencySymbol rune) {

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	ia.SetDecimalSeparator(decimalSeparator)
	ia.SetThousandsSeparator(thousandsSeparator)
	ia.SetCurrencySymbol(currencySymbol)
}

// SetNumericSeparatorsToDefaultIfEmpty - If numeric separators are
// set to zero or nil, this method will set those numeric
// separators to the USA defaults. This means that the
// Decimal separator is set to ('.'), the Thousands separator
// is set to (',') and the currency symbol is set to '$'.
//
// If the numeric separators were previously set to a value
// other than zero or nil, that value is not altered by this
// method.
//
// Effectively, this method ensures that numeric separators
// are set to valid values.
func (ia *IntAry) SetNumericSeparatorsToDefaultIfEmpty() {

	if ia.decimalSeparator == 0 {
		ia.decimalSeparator = '.'
	}

	if ia.thousandsSeparator == 0 {
		ia.thousandsSeparator = ','
	}

	if ia.currencySymbol == 0 {
		ia.currencySymbol = '$'
	}

}

// SetNumericSeparatorsToUSADefault - Sets Numeric separators:
//
//	Decimal Point Separator
//	Thousands Separator
//	Currency Symbol
//
// to United States of America (USA) defaults.
//
// Call specific methods to set numeric separators for other countries or
// cultures:
//
//	ia.SetDecimalSeparator()
//	ia.SetThousandsSeparator()
//	ia.SetCurrencySymbol()
func (ia *IntAry) SetNumericSeparatorsToUSADefault() {
	ia.SetDecimalSeparator('.')
	ia.SetThousandsSeparator(',')
	ia.SetCurrencySymbol('$')
}

// SetNumericSeparatorsDto - Sets the values of numeric separators:
//
//	decimal point separator
//	thousands separator
//	currency symbol
//
// based on values transmitted through input parameter 'customSeparators'.
//
// If any of the values contained in input parameter 'customSeparators' are set
// to zero, an error will be returned.
func (ia *IntAry) SetNumericSeparatorsDto(customSeparators NumericSeparatorDto) error {

	ePrefix := "IntAry.SetNumericSeparatorsDto() "

	if customSeparators.DecimalSeparator == 0 {
		return errors.New(ePrefix +
			"Error: Input Parameter customSeparators.DecimalSeparator is set to '0' - Invalid rune!")
	}

	if customSeparators.ThousandsSeparator == 0 {
		return errors.New(ePrefix +
			"Error: Input Parameter customSeparators.ThousandsSeparator is set to '0' - Invalid rune!")
	}

	if customSeparators.CurrencySymbol == 0 {
		return errors.New(ePrefix +
			"Error: Input Parameter customSeparators.CurrencySymbol is set to '0' - Invalid rune!")
	}

	ia.decimalSeparator = customSeparators.DecimalSeparator
	ia.thousandsSeparator = customSeparators.ThousandsSeparator
	ia.currencySymbol = customSeparators.CurrencySymbol

	return nil
}

// SetPrecision - Sets the precision of the current
// IntAry instance to the value of input parameter
// 'precision'.
//
// If input parameter 'roundResult' is set to 'true',
// the resulting numeric value will be rounded to the
// number of decimal places specified by input parameter
// 'precision'.
//
// If input parameter 'roundResult' is set to 'false',
// the resulting numeric value will be truncated to the
// number of decimal places specified by input parameter
// 'precision'.
//
// If 'precision' is set to a value less than zero,
// an error will be returned.
//
// If 'precision' is greater than the existing precision,
// trailing zeros will be added
// Examples:
//
//		 Original       			'newPrecision'				Resulting
//	   Value								input parameter			  Value
//	 --------------				---------------     -------------
//		654.123456									9							 654.123456000
//		654.123456									4							 654.1235
//
// -654.123456									9							-654.123456000
// -654.123456									4							-654.1235
//
//			0													3								 0.000
//	   0.000000									0								 0
func (ia *IntAry) SetPrecision(precision int, roundResult bool) error {

	if precision < 0 {
		return fmt.Errorf("SetPrecision() - Error: 'precision' value is less than ZERO! precision= '%v'", precision)
	}

	err := ia.IsValid("SetPrecision() - ")

	if err != nil {
		return err
	}

	if ia.isZeroValue {
		ia.SetIntAryToZero(uint(precision))
		return nil
	}

	if precision == ia.precision {
		return nil
	}

	if precision > ia.precision {
		deltaPrecision := precision - ia.precision

		for i := 0; i < deltaPrecision; i++ {
			ia.intAry = append(ia.intAry, 0)
		}

		ia.precision = precision
		ia.intAryLen = len(ia.intAry)
		ia.SetInternalFlags()
		return nil
	}

	// Must ia.precision > precision

	if roundResult {
		_ = ia.RoundToPrecision(precision)
		return nil
	}

	intLen := ia.intAryLen - ia.precision
	newAryLen := intLen + precision
	ia.intAry = ia.intAry[0:newAryLen]

	ia.intAryLen = newAryLen
	ia.precision = precision

	return nil
}

// SetSign - Can be used to change the sign value
// of the current intAry value. The new sign value
// will be set according the input parameter, 'signVal'.
//
// 'signVal' has only two valid values, -1 or +1. If
// any value other than -1 or +1 is detected, an error
// will be thrown.
func (ia *IntAry) SetSign(signVal int) error {

	ia.SetInternalFlags()

	if ia.isZeroValue {
		return nil
	}

	if signVal != -1 && signVal != 1 {
		return fmt.Errorf("SetSign() - Input parameter 'signVal' is INVALID. ")
	}

	if signVal != ia.signVal {
		ia.signVal = signVal
	}

	return nil
}

// SetSignificantDigitIdxs - Finds the first
// significant digit (the first numeric digit
// greater than zero) and sets index value in
// the local field variable, 'firstDigitIdx'.
//
// In addition, this method also identifies the
// Last Significant Digit (the last non-zero value
// in the intAry) and records that index in the
// local field variable, 'lastDigitIdx'.
func (ia *IntAry) SetSignificantDigitIdxs() {

	ia.SetNumericSeparatorsToDefaultIfEmpty()

	ia.intAryLen = len(ia.intAry)

	if ia.intAryLen == ia.precision {
		ia.intAry = append([]uint8{0}, ia.intAry...)
		ia.intAryLen++
	}

	if ia.intAryLen < ia.precision {

		deltaZeros := ia.precision - ia.intAryLen + 1
		zeroAry := make([]uint8, deltaZeros)
		ia.intAry = append(zeroAry, ia.intAry...)
		ia.intAryLen += deltaZeros
	}

	ia.firstDigitIdx = -1
	ia.lastDigitIdx = -1

	ia.integerLen = 0
	ia.significantIntegerLen = 0
	ia.significantFractionLen = 0

	lastIntIdx := ia.intAryLen - ia.precision - 1
	ia.isZeroValue = true
	ia.isIntegerZeroValue = true
	ia.integerLen = ia.intAryLen - ia.precision

	for i := 0; i < ia.intAryLen; i++ {
		if ia.intAry[i] > 0 {
			ia.isZeroValue = false

			if i < ia.integerLen {
				ia.isIntegerZeroValue = false
			}
		}

		// At minimum, there should be a single
		// leading zero before the decimal point.
		// Example 0.000.
		if i == lastIntIdx && ia.intAry[i] == 0 {

			if ia.firstDigitIdx == -1 {
				ia.firstDigitIdx = i
			}

		}

		if ia.intAry[i] > 0 {

			if ia.firstDigitIdx == -1 {
				ia.firstDigitIdx = i
			}

			ia.lastDigitIdx = i
		}

	}

	ia.significantIntegerLen = ia.intAryLen - ia.precision - ia.firstDigitIdx

	if ia.lastDigitIdx >= ia.integerLen {
		ia.significantFractionLen = ia.lastDigitIdx - ia.integerLen + 1
	} else {
		ia.significantFractionLen = 0
	}
}

// SetNumericSeparators - Used to assign values for the Decimal and Thousands separators as well
// as the Currency Symbol to be used in displaying the current intAry number string.
//
// Different nations and cultures use different symbols to delimit numerical values. In the
// USA and many other countries, a period character ('.') is used to delimit integer and
// fractional digits within a numeric value (123.45). Likewise, thousands may be delimited
// by a comma (','). Currency signs very by nationality. For instance, the USA, Canada and
// several other countries use the dollar sign ($) as a currency symbol.
//
// For a list of major world currency symbols see:
//
//		MikeAustin71\mathopsgo\mathops\mathopsconstants.go
//	 http://www.xe.com/symbols.php
//
// Note: If zero values are submitted as input for separator values, those values will default
// to USA standards.
//
// USA Examples:
//
// Decimal Separator period ('.') 		= 123.456
// Thousands Separator comma (',') 		= 1,000,000,000
// Currency Symbol dollar sign ('$')	= $123
func (ia *IntAry) SetSeparators(decimalSeparator, thousandsSeparator, currencySymbol rune) {

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	ia.decimalSeparator = decimalSeparator
	ia.thousandsSeparator = thousandsSeparator
	ia.currencySymbol = currencySymbol
}

// ShiftPrecisionLeft - Shifts the relative position of a decimal point within a number
// string. The position of the decimal point is shifted 'shiftPrecision' positions to
// the left of the current decimal point position.
//
// This is equivalent to: result = signedNumStr / 10^precision or signedNumStr divided
// by 10 raised to the power of precision.
//
// See the Examples section below.
//
// Input Parameters
// ================
//
//	shiftPrecision		uint		- The number of digits by which the current decimal point
//															point position in the current IntAry numeric value, will
//															be shifted to the left.
//
// Examples
// ========
//
//											Requested
//	                     Shift
//	 signedNumStr			precision				Result
//		 "123456.789"				  3						"123.456789"
//		 "123456.789"				  2						"1234.56789"
//		 "123456.789"	   		  6					  "0.123456789"
//		 "123456789"					6						"123.456789"
//		 "123"								5						"0.00123"
//	  "0"									3						"0.000"
//		 "0.000"							2						"0.00000"
//	 "123456.789"					0						"123456.789"		- zero 'shiftPrecision' has no effect on
//																												original number string
//
// "-123456.789"          0          "-123.456789"
// "-123456.789"          3          "-123.456789"
// "-123456789"						6					 "-123.456789"
func (ia *IntAry) ShiftPrecisionLeft(shiftPrecision uint) {

	IntAryMathDivide{}.DivideByTenToPower(ia, shiftPrecision)

}

// ShiftPrecisionRight - Shifts the existing precision of the current IntAry numeric value.
// The position of the decimal point is shifted 'shiftPrecision' positions to the right.
//
// This is equivalent to: result = IntAry X 10^shiftPrecision or IntAry Multiplied by 10
// raised to the power of input parameter 'shiftPrecision'.
//
// Input Parameters:
// =================
//
//	shiftPrecision		uint		- The number of digits by which the current decimal point
//															point position in the current IntAry numeric value will
//															be shifted to the right.
//
// Examples:
// =========
//
//	  Input
//	Parameter
//
// IntAry Value		shiftPrecision		Result
// ------------   --------------    ------
//
//	"123456.789"				3						"123456789"
//	"123456.789"				2						"12345678.9"
//	"123456.789"        6					  "123456789000"
//	"123456789"	 			  6						"123456789000000"
//	"123"               5	          "12300000"
//	"0"								  3						"0"
//	"123456.789"				0						"123456.789"		- zero has no effect on original number string
//
// "-123456.789"        0          "-123456.789"
// "-123456.789"        3          "-123456789"
// "-123456789"			    6					 "-123456789000000"
func (ia *IntAry) ShiftPrecisionRight(shiftPrecision uint) {

	IntAryMathMultiply{}.MultiplyByTenToPower(ia, shiftPrecision)

}

// SetThousandsSeparator is used to set the value of the thousands
// separator character. The thousands separator character is used
// to separate thousands when the value of this IntAry object is
// expressed as a string. In the US, the thousands separator character
// is a comma. In the following example, the ',' character separates
// thousands in the integer number: Example - '1,000,000,000'.
//
// The SetThousandsSeparator method can be used to change the thousands
// separator character in accordance with the customs of countries other
// than the US.
//
// To examine the current setting for thousands separator, see
// method GetIntAryStats().
//
// Note the default thousands separator character is the comma (',').
func (ia *IntAry) SetThousandsSeparator(thousandsSeparator rune) {
	ia.thousandsSeparator = thousandsSeparator
}

// SubtractFromThis - Subtracts the value of parameter
// 'ia2' from the current intAry object.
//
// Input Parameters:
// =================
//
// ia2 *intAry - Incoming intAry object whose value will be subtracted
//
//	from this current intAry value.
func (ia *IntAry) SubtractFromThis(ia2 *IntAry) error {

	IntAryMathSubtract{}.SubtractTotal(ia, ia2)

	return nil

}

// String - Returns the numeric value of the current IntAry as
// a string.
func (ia *IntAry) String() string {
	return ia.GetNumStr()
}

// SubtractMultipleFromThis - This method will subtract multiple intAry values from the
// current intAry value. There are two input parameters:
//
// convertToNumStr bool - If true the result will be converted to a number string after
//
//	the final subtraction operation.
//
// iaMany ...*intAry - An unlimited series of pointers to intAry objects which will be
//
//	subtracted from the current intAry Value.
func (ia *IntAry) SubtractMultipleFromThis(iaMany ...*IntAry) error {

	for _, iAry := range iaMany {

		IntAryMathSubtract{}.SubtractTotal(ia, iAry)

	}

	return nil
}

// getMaximumPrecision - Returns the maximum allowable precision for
// the IntAry type as an unsigned integer (uint). Currently the
// maximum allowable precision is computed as:
//
//	21474836472147483647 - 2 = 21474836472147483645
//
// Effectively this is the maximum value of an int32 minus 2.
func (ia *IntAry) getMaximumPrecision() uint {

	return uint(math.MaxInt32) - uint(2)
}

// validateUintToMaxPrecision - The purpose of this method is to ensure that
// a precision value does not exceed the maximum precision capacity of
// the IntAry Type. The maximum allowable precision value is computed
// as an unsigned int.
func (ia *IntAry) validateUintToMaxPrecision(origPrecision uint) uint {

	maxTypePrecision := ia.getMaximumPrecision()

	if origPrecision > maxTypePrecision {
		return maxTypePrecision
	}

	return origPrecision
}
