package mathops

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
)

// BigIntFixedDecimal
//
//	A light data transfer structure used to represent a
//	numeric value with a fixed number of decimal digits.
//	Used primarily for intensive and complex calculations.
type BigIntFixedDecimal struct {
	integerNum *big.Int
	//	All the numeric digits, both integer and
	//	fractional, necessary to define a fixed length
	//	floating point number.
	//
	//	The number of digits to the right of the decimal
	//	place is specified by the data field,
	//	BigIntFixedDecimal.precision.

	precision uint
	//	Specifies the number of digits to the right of the
	//	decimal place in the series of numeric digits
	//	represented by BigIntFixedDecimal.precision.
	//
	// 	Example: To represent the floating point number 52.459
	// 			a BigIntDecimal Structure would be configured
	//			as follows:
	//
	// 			BigIntFixedDecimal.integerNum	= 52459
	// 			BigIntFixedDecimal.precision	= 3

	lock *sync.Mutex
}

// Ceiling - Returns the ceiling integer value for the current
// BigIntFixedDecimal instance.
//
// Ceiling is defined as: The least, or lowest value integer,
// which is greater than or equal to the numeric value of the
// current BigIntFixedDecimal.
//
// Reference Wikipedia:
//
//	https://en.wikipedia.org/wiki/Floor_and_ceiling_functions
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
func (bigIFd *BigIntFixedDecimal) Ceiling() BigIntFixedDecimal {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	cmpZeroResult := bigIFd.integerNum.Cmp(big.NewInt(0))

	if cmpZeroResult == 0 {
		return new(BigIntFixedDecimal).NewZero(0)
	}

	ceiling := big.NewInt(0).Set(bigIFd.integerNum)

	if bigIFd.precision > 0 {

		scale := big.NewInt(0).Exp(
			big.NewInt(10),
			big.NewInt(int64(bigIFd.precision)),
			nil)

		ceiling.Quo(ceiling, scale)

		if cmpZeroResult == 1 {
			// signVal must be plus
			ceiling.Add(ceiling, big.NewInt(1))
		}

	}

	// else bigIFd.precision must be zero
	return new(BigIntFixedDecimal).New(ceiling, 0)
}

// ChangeSign - This method will change the sign of the
// current BigIntFixedDecimal numeric value. If the value
// is negative, this method will change the sign to
// positive. Likewise, if the sign is currently positive,
// calling this method will change the sign to negative.
func (bigIFd *BigIntFixedDecimal) ChangeSign() {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return
	}

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		return
	}

	bigIFd.integerNum.Neg(bigIFd.integerNum)
}

// Cmp - Compares the numeric values of two BigIntFixedDecimal
// instances.
//
// If the current BigIntFixedDecimal value is greater than input
// parameter 'fd2', the method returns '1'.
//
// If the current BigIntFixedDecimal value is equal to that of the
// input parameter 'fd2', the method returns '0'.
//
// If the current BigIntFixedDecimal value is less than the input
// parameter 'fd2', the method returns '-1'.
//
// Examples
// ========
//
//						BigIntFixedDecimal         'fd2'        Return
//									Value                Value        Value
//						------------------         ------       -------
//	                5                     2             1
//	                5.2                   5.1           1
//	                5.2                   5.2           0
//	           837123.4              837123.5          -1
//	                0                     0.1          -1
//	               35.123456             40.5          -1
//	               35.123456              2.5           1
func (bigIFd *BigIntFixedDecimal) Cmp(fd2 BigIntFixedDecimal) int {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	fd2.IsValid()

	if bigIFd.precision == fd2.precision {
		return bigIFd.integerNum.Cmp(fd2.integerNum)
	}

	bigTen := big.NewInt(10)

	if fd2.precision > bigIFd.precision {
		delta := int64(fd2.precision - bigIFd.precision)
		fdValue := big.NewInt(0).Set(bigIFd.integerNum)
		scale := big.NewInt(0).Exp(bigTen, big.NewInt(delta), nil)
		fdValue.Mul(fdValue, scale)

		return fdValue.Cmp(fd2.integerNum)

	}

	// MUST BE bigIFd.precision > fd2.precision
	delta := int64(bigIFd.precision - fd2.precision)
	fd2Value := big.NewInt(0).Set(fd2.integerNum)
	scale := big.NewInt(0).Exp(bigTen, big.NewInt(delta), nil)
	fd2Value.Mul(fd2Value, scale)

	return bigIFd.integerNum.Cmp(fd2Value)
}

// CmpZero Compares the current BigIntFixedDecimal to Zero and
// returns an integer flag as follows:
//
// +1 = BigIntFixedDecimal > 0
//
//	0 = BigIntFixedDecimal == 0
//
// -1 = BigINtFixedDecimal < 0
func (bigIFd *BigIntFixedDecimal) CmpZero() int {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	return bigIFd.integerNum.Cmp(big.NewInt(0))
}

// CopyIn - Receives a BigIntFixedDecimal type and copies the
// value to the current BigIntFixedDecimal instance.
func (bigIFd *BigIntFixedDecimal) CopyIn(fd BigIntFixedDecimal) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	fd.IsValid()

	intVal := fd.GetInteger()

	if intVal == nil {
		intVal = big.NewInt(0)
	}

	bigIFd.integerNum = big.NewInt(0).Set(intVal)
	bigIFd.precision = fd.GetPrecision()

}

// CopyIn - Receives a pointer to a BigIntFixedDecimal type and
// copies the value to the current BigIntFixedDecimal instance.
func (bigIFd *BigIntFixedDecimal) CopyInPtr(fd *BigIntFixedDecimal) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	fd.IsValid()

	intVal := fd.GetInteger()

	if intVal == nil {
		intVal = big.NewInt(0)
	}

	bigIFd.integerNum = big.NewInt(0).Set(intVal)
	bigIFd.precision = fd.GetPrecision()

}

// CopyOut - Returns a new BigIntFixedDecimal instance which is
// a deep copy of the current BigIntFixedDecimal instance.
func (bigIFd *BigIntFixedDecimal) CopyOut() BigIntFixedDecimal {

	if bigIFd.integerNum == nil {
		bigIFd.SetNumericValue(big.NewInt(0), bigIFd.precision)
	}

	return new(BigIntFixedDecimal).New(bigIFd.integerNum, bigIFd.precision)
}

// DivideByTenToPower - Divides the numeric value of the current
// BigIntFixedDecimal by 10 to the power of 'exponent'.
//
//	result = BigIntFixedDecimal / 10^exponent
//
// Input Parameter
// ===============
//
// exponent	uint	- The value of the current BigIntFixedDecimal
//
//	instance will be divided by ten raised to
//	the power of 'exponent'.
//
// This method will destroy and overwrite the previous value of
// the current BigIntFixedDecimal instance with the results of
// this calculation.
func (bigIFd *BigIntFixedDecimal) DivideByTenToPower(exponent uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		return
	}

	scale :=
		big.NewInt(0).Exp(
			big.NewInt(10),
			big.NewInt(int64(exponent)), nil)

	factor := new(BigIntFixedDecimal).New(scale, 0)

	result, _ :=
		new(BigIntMathDivide).FixedDecimalFracQuotient(
			bigIFd.CopyOut(), factor, bigIFd.precision+exponent)

	bigIFd.CopyIn(result)
}

// DivideByTwoToPower - Performs integer division by two using a 'right-shift'
// technique. Remainders from this division operation are discarded, only the
// integer quotient is returned. When the calculation is completed, the value
// of the integer quotient will replace the old value of the current BigIntFixedDecimal
// instance.
//
// Example:
// ========
//
//	quotient =  BigIntFixedDecimal / 2^(exponent)
//
// In this example BigIntFixedDecimal= 33,123.456, so 33,123.456/ 2^8:
//
//			(1) The fractional quotient of 33,123.456/256 (or 2^8) is 129.3885.
//
//	   (2) This method will use a right shift technique on the integer value
//	       33123456 / 2^(8) to generate an integer quotient of 129388.
//
// Consider the example BigIntFixedDecimal = 33123456 (no decimal fraction):
//
//	Dividing 33123456 / 2^8 = fractional quotient = 129388.5
//	Be careful when using this method.
//
// **************************************************************************
// (1) 	Be sure to consider the outcomes when sending a decimal fraction to
//
//	this method.
//
// (2)	Results returned by this method will always have precision = 0
//
// **************************************************************************
//
// Note:	This method will destroy and overwrite the previous value
//
//	of the current BigIntFixedDecimal instance with the results
//	of this calculation.
func (bigIFd *BigIntFixedDecimal) DivideByTwoToPower(exponent uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return
	}

	bigIFd.integerNum =
		BigIntMathDivide{}.BigIntDividedByTwoToPower(
			bigIFd.integerNum,
			exponent)

	bigIFd.precision = 0
}

// Empty - Reinitialize the current BigIntFixedDecimal instance
// to a zero value with zero precision.
func (bigIFd *BigIntFixedDecimal) Empty() {

	bigIFd.integerNum = big.NewInt(0)
	bigIFd.precision = 0

}

// Floor - Returns the floor integer value for the current
// BigIntFixedDecimal.
//
// In mathematics and computer science, the floor function
// is the function that takes as input a real number x dnd
// gives as output the greatest integer less than or equal
// to x.
//
// Source:
// https://en.wikipedia.org/wiki/Floor_and_ceiling_functions
//
// Examples
// ========
//
//					 BigIntFixedDecimal           Floor
//								Value                   Value
//					 ------------------         ----------
//	             0													0
//	             4                         4
//								3.2												3
//								2.9	                      2
//							 -2.7											 -3
//							 -2                        -2
func (bigIFd *BigIntFixedDecimal) Floor() BigIntFixedDecimal {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	cmpZeroResult := bigIFd.integerNum.Cmp(big.NewInt(0))

	if cmpZeroResult == 0 {
		return new(BigIntFixedDecimal).NewZero(0)
	}

	floor := big.NewInt(0).Set(bigIFd.integerNum)

	if bigIFd.precision > 0 {

		scale := big.NewInt(0).Exp(
			big.NewInt(10),
			big.NewInt(int64(bigIFd.precision)),
			nil)

		floor.Quo(floor, scale)

		if cmpZeroResult == -1 {
			// signVal must be -1
			floor.Add(floor, big.NewInt(-1))
		}

	}

	// else bigIFd.precision must be zero
	return new(BigIntFixedDecimal).New(floor, 0)
}

// FormatNumStr - converts the numeric value of the current BigIntFixedDecimal
// instance to a number string. The returned number string will consist of a
// string of numeric digits. If the number contains fractional digits, the
// decimal separator period (.) will be used to separate integer and fractional
// digits within the string. There are no thousands separator present in the
// the returned string.
//
// The input parameter 'negValMode' is of type NegativeValueFmtMode.
// NegativeValueFmtMode encompasses a series of constants that are used
// to format negative values in a number string.
//
// Valid NegativeValueFormatMode's are defined as follows:
//
// LEADMINUSNEGVALFMTMODE 		-	Negative values formatted with
//
//	 		a leading minus sign.
//			Example: -123456.78
//
// PARENTHESESNEGVALFMTMODE	-	Negative values formatted with
//
//	surrounding parentheses.
//	Example: (123456.78)
//
// ABSOLUTEPURENUMSTRFMTMODE - Formats a pure integerNum string with
//
//	 absolute (positive) integer value
//	 and no decimal point separator.
//	Example: (12345678)
func (bigIFd *BigIntFixedDecimal) FormatNumStr(negValMode NegativeValueFmtMode) string {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	decimalSeparator := '.'
	baseZero := big.NewInt(0)
	bigIFdSign := 1

	if bigIFd.integerNum.Cmp(baseZero) == -1 {
		bigIFdSign = -1
	}

	absBigInt := big.NewInt(0).Set(bigIFd.integerNum)

	if absBigInt.Cmp(baseZero) == -1 {
		absBigInt.Neg(absBigInt)
	}

	outRunes := make([]rune, 0, 300)

	scratchNum := big.NewInt(0).Set(absBigInt)

	if scratchNum.Cmp(baseZero) == 0 {
		bigIFdSign = 1

		outRunes = append(outRunes, '0')

		if bigIFd.precision > 0 {

			if negValMode != ABSOLUTEPURENUMSTRFMTMODE {
				outRunes = append(outRunes, decimalSeparator)
			}

			cnt := int(bigIFd.precision)

			if negValMode == ABSOLUTEPURENUMSTRFMTMODE {
				cnt--
			}

			for h := 0; h < cnt; h++ {
				outRunes = append(outRunes, '0')
			}

		}

		return string(outRunes)
	}

	startIdx := 0
	modulo := big.NewInt(0)
	baseTen := big.NewInt(10)
	digitCnt := 0

	for scratchNum.Cmp(baseZero) == 1 {

		if startIdx == 0 &&
			bigIFdSign == -1 &&
			negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, ')')
		}

		modX := big.NewInt(0)
		scratchNum, modulo = big.NewInt(0).QuoRem(scratchNum, baseTen, modX)
		outRunes = append(outRunes, rune(modulo.Int64()+int64(48)))
		digitCnt++
		startIdx++

		if bigIFd.precision > 0 &&
			int(bigIFd.precision) == startIdx &&
			negValMode != ABSOLUTEPURENUMSTRFMTMODE {

			outRunes = append(outRunes, decimalSeparator)
			startIdx++
		}

	}

	if int(bigIFd.precision) >= digitCnt {

		delta := int(bigIFd.precision) - digitCnt + 1

		if negValMode == ABSOLUTEPURENUMSTRFMTMODE {
			delta--
		}

		for k := 0; k < delta; k++ {
			outRunes = append(outRunes, '0')
			startIdx++

			if bigIFd.precision > 0 &&
				int(bigIFd.precision) == startIdx &&
				negValMode != ABSOLUTEPURENUMSTRFMTMODE {

				outRunes = append(outRunes, decimalSeparator)
				startIdx++
			}
		}
	}

	startIdx--

	// adjust for negative sign value
	if bigIFdSign == -1 {

		if negValMode == LEADMINUSNEGVALFMTMODE {
			outRunes = append(outRunes, '-')
			startIdx++

		} else if negValMode == PARENTHESESNEGVALFMTMODE {
			outRunes = append(outRunes, '(')
			startIdx += 2
		}

		/*
				MUST BE negValMode == ABSOLUTEPURENUMSTRFMTMODE
			  Do NOT Display Sign Character

		*/
	}

	sortLimit := startIdx / 2
	tRune := rune(0)
	yCnt := 0

	for i := startIdx; i > sortLimit; i-- {
		tRune = outRunes[yCnt]
		outRunes[yCnt] = outRunes[i]
		outRunes[i] = tRune
		yCnt++
	}

	return string(outRunes)
}

// GetInteger - Returns the 'integerNum' for the current
// BigIntFixedDecimal instance. The returned *big.Int type
// contains all the numeric digits which comprise the fixed
// decimal numerical value represented by this BigIntFixedDecimal
// instance. If the value includes fractional digits, these
// too are included in the returned integer value.
//
// Example:
// ========
//
// BigIntFixedDecimal      Returned
//
//	Numeric Value        Integer Value
//
// ------------------    -------------
//
//	582.12345            58212345
func (bigIFd *BigIntFixedDecimal) GetInteger() *big.Int {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	return big.NewInt(0).Set(bigIFd.integerNum)
}

// GetIntegerFractionalParts - Returns two BigIntFixedDecimals comprising the integer
// and fractional parts of the current BigIntFixedDecimal numeric value.
//
// Examples:
// =========
//
// BigIntFixedDecimal              Returned       Integer    		 Returned					 Fraction
//
//	Numeric Value              Integer Value    Precision		Fractional Value	   Precision
//
// ------------------           -------------   ----------		----------------     ---------
//
//	859649.123456789								859649					0						0.123456789            9
//
// -859649.123456789							 -859649				 	0					 -0.123456789            9
func (bigIFd *BigIntFixedDecimal) GetIntegerFractionalParts() (integer BigIntFixedDecimal, fraction BigIntFixedDecimal) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	integer = new(BigIntFixedDecimal).NewZero(0)
	fraction = new(BigIntFixedDecimal).NewZero(0)

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		return integer, fraction
	}

	scale := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(uint64(bigIFd.precision)), nil)

	scratch := big.NewInt(0)

	intRadicand, fracRadicand := big.NewInt(0).QuoRem(bigIFd.integerNum, scale, scratch)

	integer = new(BigIntFixedDecimal).New(intRadicand, 0)

	fraction = new(BigIntFixedDecimal).New(fracRadicand, bigIFd.precision)

	return integer, fraction
}

// GetMagnitude - Returns the magnitude of the current BigIntFixedDecimal
// as a type *big.Int integer value.
//
// Magnitude as used here is defined as the power of 10 which generates a
// value less than or equal to the current BigIntFixedDecimal value
//
//	10^magnitude  <= BigIntFixedDecimal
//
// Examples
// ========
//
//	        BigIntFixedDecimal
//				   			Value									magnitude
//	             ------        				---------
//
//				  		 963,256										5
//										 2										0
//										32										1
//				 8,456,123,921					  				9
//	                  2.2									0
//	      8,456,123,912.123                9
//	           -643,212.123								5
//	                324.123456             2
func (bigIFd *BigIntFixedDecimal) GetMagnitude() (*big.Int, error) {

	bigZero := big.NewInt(0)

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return bigZero, nil
	}

	target := big.NewInt(0).Set(bigIFd.integerNum)
	bigTen := big.NewInt(10)

	if bigIFd.precision > 0 {
		target.Quo(
			target,
			big.NewInt(0).Exp(bigTen, big.NewInt(int64(bigIFd.precision)),
				nil))
	}

	magnitude, err := BigIntMath{}.GetMagnitude(target)

	if err != nil {
		ePrefix := "BigIntFixedDecimal.GetMagnitude() "
		return bigZero,
			fmt.Errorf(ePrefix+
				"Error returned "+
				"Error='%v' ", err.Error())
	}

	return magnitude, nil
}

// GetNumericValue - Returns the 'integerNum' and 'precision' values for the
// current BigIntFixedDecimal instance.
func (bigIFd *BigIntFixedDecimal) GetNumericValue() (*big.Int, uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	return big.NewInt(0).Set(bigIFd.integerNum), bigIFd.precision
}

// GetNumStr - Converts the current BigIntFixedDecimal value to
// string of numbers which includes the decimal place and decimal
// digits if they exist.
func (bigIFd *BigIntFixedDecimal) GetNumStr() string {

	return bigIFd.FormatNumStr(LEADMINUSNEGVALFMTMODE)

}

// GetPrecision - Returns the 'precision' value for the current
// BigIntFixedDecimal instance. 'precision' specifies the number
// of digits to the right of the decimal place in the
// BigIntFixedDecimal.integerNum.
func (bigIFd *BigIntFixedDecimal) GetPrecision() uint {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	return bigIFd.precision
}

// GetPrecisionBigInt - Returns the 'precision' value for the current
// BigIntFixedDecimal instance. 'precision' specifies the number
// of digits to the right of the decimal place in the
// BigIntFixedDecimal.integerNum. For this method, precision is
// returned as a *big.Int type.
func (bigIFd *BigIntFixedDecimal) GetPrecisionBigInt() *big.Int {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	return big.NewInt(0).SetUint64(uint64(bigIFd.precision))
}

// Inverse - Converts the current BigIntFixedDecimal
// to the inverse of its numeric value.
//
// Example:
// Current Value = '2'         Inverse= '1/2'
//
// Note: When called this method will destroy and
// overwrite the previous numeric value with the
// inverse value.
//
// Input Parameter
// ===============
//
// maxPrecision uint	- Defines the maximum precision for the
//
//	                     inverse value computed by this method.
//	                     As used here, 'maxPrecision' specifies
//	                     the maximum number of numeric digits to
//												the right of the decimal place.
func (bigIFd *BigIntFixedDecimal) Inverse(maxPrecision uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return
	}

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		return
	}

	bigOne := big.NewInt(1)

	inverseBigInt, inversePrecision, _ :=
		BigIntMathDivide{}.BigIntFracQuotient(
			bigOne,
			big.NewInt(0),
			bigIFd.integerNum,
			big.NewInt(0).SetUint64(uint64(bigIFd.precision)),
			big.NewInt(0).SetUint64(uint64(maxPrecision)))

	bigIFd.SetNumericValue(inverseBigInt, uint(inversePrecision))

	return
}

// IsInteger - Returns true if the numeric value of
// the current BigIntFixedDecimal is an integer value.
func (bigIFd *BigIntFixedDecimal) IsInteger() bool {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return true
	}

	if bigIFd.precision == 0 {
		return true
	}

	return false
}

// IsEven - returns 'true' if the numeric value
// of the current BigIntFixedDecimal is even.
//
// ------------------------------------------------------
// "In mathematics, parity is the property of an
// integer's inclusion in one of two categories:
// even or odd. An integer is even if it is evenly
// divisible by two and odd if it is not even."
//
// "Examples of even numbers include −4, 0, 82 and 178."
// In particular, zero is an even number."
// ------------------------------------------------------
// https://en.wikipedia.org/wiki/Parity_(mathematics)
//
// Also, see 	https://www.mathsisfun.com/definitions/even-number.html
func (bigIFd *BigIntFixedDecimal) IsEven() bool {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return true
	}

	if bigIFd.precision > 0 {
		//BigIntFixedDecimal is NOT an integer
		return false
	}

	bigZero := big.NewInt(0)

	if bigIFd.integerNum.Cmp(bigZero) == 0 {
		return true
	}

	remainder := big.NewInt(0).Rem(bigIFd.integerNum, big.NewInt(2))

	if remainder.Cmp(bigZero) == 0 {
		return true
	}

	return false
}

// IsValid - IsValid test the validity of the internal
// data field BigIntFixedDecimal.integerNum. If this
// data field is 'nil', BigIntFixedDecimal.integerNum
// is set to zero and the function returns 'false'.
//
// Otherwise, the function returns 'true'.
func (bigIFd *BigIntFixedDecimal) IsValid() bool {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return false
	}

	return true
}

// IsZero - returns true only if the current BigIntFixedDecimal
// numeric value is equal to zero
func (bigIFd *BigIntFixedDecimal) IsZero() bool {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	return false
}

// MultiplyByTenToPower - Multiplies the numeric value of the current
// BigIntFixedDecimal by 10 to the power of 'exponent'.
//
//	result = BigIntFixedDecimal x 10^exponent
//
// Examples:
// =========
//
//		BigIntFixedDecimal
//				Value									Exponent							   Result
//	 ------------------        --------              -----------------
//	  105.6752										 0										      105.6752
//	  105.6752                    1									       1056.752
//	  105.6752                    2                        10567.52
//	  105.6752                    3                       105675.2
//	  105.6752                    8                  10567520000
//
// Input Parameter
// ===============
//
// exponent	uint	- The value of the current BigIntFixedDecimal
//
//	instance will be multiplied by ten raised to
//	the power of 'exponent'.
//
// Note:	(1)	This method will delete trailing fractional zeros from
//
//	 			the returned product.
//
//		(2)	This method will destroy and overwrite the previous value
//				of the current BigIntFixedDecimal instance with the results
//				of this calculation.
func (bigIFd *BigIntFixedDecimal) MultiplyByTenToPower(exponent uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		bigIFd.precision = 0
		return
	}

	scale :=
		big.NewInt(0).Exp(
			big.NewInt(10),
			big.NewInt(int64(exponent)), nil)

	factor := new(BigIntFixedDecimal).New(scale, 0)

	result := BigIntMathMultiply{}.FixedDecimalMultiply(bigIFd.CopyOut(), factor)

	result.TrimTrailingFracZeros()

	bigIFd.CopyIn(result)

}

// MultiplyByTwoToPower - Multiplies the numeric value of the current
// BigIntFixedDecimal by 2 to the power of 'exponent'.
//
//	product = BigIntFixedDecimal x 2^exponent
//
// When the calculation is completed, the value of 'product' will
// replace the old value of the current BigIntFixedDecimal instance.
//
// Examples:
// =========
//
//	 BigIntFixedDecimal 			 	exponent		  product
//	interNum		Precision
//
// -------------------------------------------------------------
//
//	12345						5								15				4045.2096
//								(0.12345 x 2^15 = 4045.2096)
//
// -------------------------------------------------------------
//
//	571						1								 8			 14617.6
//	            (57.1 x 2^8 = 14617.6)
//
// -------------------------------------------------------------
//
// Note:	(1)	This method will delete trailing fractional zeros from
//
//	 			the returned product.
//
//		(2)	This method will destroy and overwrite the previous value
//				of the current BigIntFixedDecimal instance with the results
//				of this calculation.
func (bigIFd *BigIntFixedDecimal) MultiplyByTwoToPower(exponent uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return
	}

	precision := big.NewInt(0)

	bigIFd.integerNum, precision, _ =
		BigIntMathMultiply{}.BigIntMultiplyByTwoToPower(
			bigIFd.integerNum,
			big.NewInt(0).SetUint64(uint64(bigIFd.precision)),
			exponent)

	bigIFd.precision = uint(precision.Uint64())

}

// New - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer	*big.Int	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	- Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) New(integer *big.Int, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(integer, precision)

	return num
}

// NewBigIntPrecision - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// bigInt			*big.Int	- Specifies the sequence of numerical digits in the numeric value.
//
// precision	*big.Int	- Specifies the number of digits to the right of the decimal point
//
//													in input parameter, 'integer'. If 'precision' is greater than the
//	                       maximum value of an unsigned integer (+4,294,967,295, which equals
//													2^32 − 1), an error will be triggered. Also, if 'precision' is less
//	                       than zero, an error will be triggered.
func (bigIFd *BigIntFixedDecimal) NewBigIntPrecision(
	bigInt, precision *big.Int) (BigIntFixedDecimal, error) {

	ePrefix := "BigIntFixedDecimal.NewBigIntPrecision() "

	if precision.Cmp(big.NewInt(0)) == -1 {
		return new(BigIntFixedDecimal).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' LESS THAN ZERO! "+
				"precision='%v' ", precision.Text(10))
	}

	maxUint32 := big.NewInt(0).SetUint64(uint64(math.MaxUint32))

	if precision.Cmp(maxUint32) == 1 {
		return new(BigIntFixedDecimal).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' exceeds maximum limit of '4,294,967,295'! "+
				"precision='%v' ", precision.Text(10))
	}

	num := BigIntFixedDecimal{}

	num.SetNumericValue(bigInt, uint(precision.Uint64()))

	return num, nil
}

// NewInt - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer	int				- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	- Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewInt(integer int, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(int64(integer)), precision)

	return num
}

// NewInt32 - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer		 int32	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	- Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewInt32(integer int32, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(int64(integer)), precision)

	return num
}

// NewInt64 - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer		 int32	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	- Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewInt64(integer int64, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(integer), precision)

	return num
}

// NewNumStr - Receives a number string as input and returns
// a new BigIntFixedDecimal instance.
//
// A number string is a string of numeric digits which may,
// or may not, be prefixed with a minus sign ('-') indicating
// a negative number. If the numeric string of digits is prefixed
// by a left parenthesis and suffixed by a corresponding right
// parenthesis, this also indicates a negative value.
//
// The numeric string of digits may also contain a period
// ('.') which is treated as a decimal separator and used to
// separate integer and fractional digits within the number
// string.
//
// The only decimal separator recognized by this method is the
// period ('.').
func (bigIFd *BigIntFixedDecimal) NewNumStr(numStr string) (BigIntFixedDecimal, error) {

	ePrefix := "BigIntFixedDecimal.NewNumStr() "

	fixedDecimal := BigIntFixedDecimal{}

	err := fixedDecimal.SetNumStr(numStr)

	if err != nil {
		return BigIntFixedDecimal{},
			fmt.Errorf(ePrefix+"Error returned by fixedDecimal.SetNumStr(numStr). "+
				"Error='%v' ", err.Error())
	}

	return fixedDecimal, nil
}

// NewUInt - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer		  uint	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	- Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewUInt(integer uint, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(0).SetUint64(uint64(integer)), precision)

	return num
}

// NewUInt32 - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer		  uint32	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	  - Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewUInt32(integer uint32, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(0).SetUint64(uint64(integer)), precision)

	return num
}

// NewUInt64 - Creates and returns a new BigIntFixedDecimal type based on input parameters,
// 'integer' and 'precision'.
//
// Input Parameters
// ================
//
// integer		  uint32	- Specifies the sequence of numerical digits in the numeric value.
//
// precision		uint	  - Specifies the number of digits to the right of the decimal point
//
//	in input parameter, 'integer'.
func (bigIFd *BigIntFixedDecimal) NewUInt64(integer uint64, precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(0).SetUint64(integer), precision)

	return num
}

// NewZero - Creates and returns a new BigIntFixedDecimal type with a zero value. The
// number of digits to the right of the decimal place is specified by input parameter,
// precision.
//
// Input Parameters
// ================
//
// precision		uint	- Specifies the number of digits to the right of the decimal point.
func (bigIFd *BigIntFixedDecimal) NewZero(precision uint) BigIntFixedDecimal {

	num := BigIntFixedDecimal{}

	num.SetNumericValue(big.NewInt(0), precision)

	return num
}

//	RoundToDecPlace
//
//	Rounds the numeric value of the current
//	BigIntFixedDecimal instance to a specified number of
//	decimal places.
//
//	'precision' equals the number of digits to the right
//	of the decimal place.
//
//	Example:
//
//		integer= 123456; precision = 3; Numeric Value= 123.456
//
//	If the value of BigIntFixedDecimal.bigInt is zero
//	('0'), that zero value will remain unaltered. However,
//	the BigIntNum.precision value will be set equal to
//	input parameter, 'precision'.
//
//	If the number of decimal places specified for rounding
//	("precision") is equal to the current
//	BigIntFixedDecimal.precision, no action is taken.
//
//	If the number of decimal places specified for rounding
//	('precision') is greater than the current
//	BigIntFixedDecimal.precision value, trailing zeros are
//	added to the current BigIntFixedDecimal.bigInt value
//	and BigIntNum.precision is set equal to input parameter,
//	'precision'.
//
//	Finally, if the number of decimal places specified for
//	rounding ('precision') is less than the current
//	BigIntNum.precision value, the fractional digits will
//	be rounded in accordance with the input parameter,
//	'precision'.
//
//	Examples:
//
//		Original	'precision'			Resulting
//		 Value		input parameter		  Value
//	 --------------	---------------     -------------
//
//		 654.123456		9				 654.123456000
//		 654.123456		4				 654.1235
//		-654.123456		9				-654.123456000
//		-654.123456		4				-654.1235
//		   0			3				   0.000
//		   0.000000		0				   0
//
// Note: This method does NOT trim or delete trailing fractional zero
// digits.
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	precision			uint
//
//		As used here, the term 'precision' defines the number of
//		fractional digits to the right of the radix point or
//		decimal separator (a.k.a. decimal point).
//
//		Example:
//
//			integer= 123456;
//			precision = 3;
//			Numeric Value= 123.456
//
// ----------------------------------------------------------------
//
// # Return Values
//
//	NONE
func (bigIFd *BigIntFixedDecimal) RoundToDecPlace(precision uint) {

	if bigIFd.lock == nil {
		bigIFd.lock = new(sync.Mutex)
	}

	bigIFd.lock.Lock()

	defer bigIFd.lock.Unlock()

	if bigIFd.integerNum == nil {
		bigIFd.SetNumericValue(big.NewInt(0), precision)
		return
	}

	if bigIFd.precision == precision {
		// Nothing to do. Specified 'precision' is already implemented.
		return
	}

	cmpToZeroResult := bigIFd.integerNum.Cmp(big.NewInt(0))

	// bigInt == zero, set precision an return
	if cmpToZeroResult == 0 {
		bigIFd.CopyIn(
			new(BigIntFixedDecimal).NewZero(precision))
		return
	}

	scale := big.NewInt(0)
	base10 := big.NewInt(10)

	// If existing precision is less than new specified precision,
	// add trailing zeros, set new precision parameter and return.
	if bigIFd.precision < precision {

		deltaPrecision := precision - bigIFd.precision

		scale =
			big.NewInt(0).Exp(
				base10,
				big.NewInt(int64(deltaPrecision)),
				nil)

		bigIFd.integerNum.Mul(bigIFd.integerNum, scale)

		bigIFd.precision += deltaPrecision

		return
	}

	// Must be: bigIFd.precision >  precision

	bigNumRound5 :=
		new(BigIntFixedDecimal).NewInt(5, precision+1)

	if cmpToZeroResult == -1 {
		bigNumRound5.integerNum.Mul(
			bigNumRound5.integerNum,
			big.NewInt(-1))
	}

	result :=
		new(BigIntMathAdd).FixedDecimalAdd(bigIFd.CopyOut(), bigNumRound5)

	// 10^deltaPrecision
	scale.Exp(big.NewInt(10),
		big.NewInt(int64(bigIFd.precision-precision)), nil)

	result.integerNum.Quo(result.integerNum, scale)

	bigIFd.SetNumericValue(result.integerNum, precision)
}

// SetIntegerValue
//
// Sets the BigIntFixedDecimal.integerNum or integer
// value for the current BigIntFixedDecimal instance.
func (bigIFd *BigIntFixedDecimal) SetIntegerValue(integer *big.Int) {

	if bigIFd.lock == nil {
		bigIFd.lock = new(sync.Mutex)
	}

	bigIFd.lock.Lock()

	defer bigIFd.lock.Unlock()

	if integer == nil {

		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0

	} else {

		bigIFd.integerNum = big.NewInt(0).Set(integer)
		bigIFd.precision = 0

	}

}

//	SetNumStr
//	Initializes the current BigIntFixedDecimal instance of
//	the numeric value of the number string input
//	parameter.
//
//	A number string is a string of numeric digits which may
//	or may not be prefixed with a minus sign ('-'). The
//	numeric string of digits may also contain a decimal
//	separator period ('.') which is used to separate integer
//	and fractional digits within the number string.
//
// The only decimal separator recognized by this method is the
// period ('.').
func (bigIFd *BigIntFixedDecimal) SetNumStr(numStr string) error {

	if bigIFd.lock == nil {
		bigIFd.lock = new(sync.Mutex)
	}

	bigIFd.lock.Lock()

	defer bigIFd.lock.Unlock()

	ePrefix := "BigIntFixedDecimal.SetNumStr() "

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	if len(numStr) == 0 {
		return errors.New(ePrefix +
			"Error: Input parameter 'numStr' is an EMPTY string!")
	}

	decimalSeparator := '.'

	baseRunes := []rune(numStr)
	lBaseRunes := len(baseRunes)

	newPrecision := uint(0)
	newAbsBigInt := big.NewInt(0)
	baseTen := big.NewInt(10)
	hasLeftParen := false
	hasRightParen := false
	hasMinusSign := false
	startFractionalDigits := false
	isStartNumericDigits := false
	isEndNumericDigits := false
	numOfNumericDigits := uint(0)

	for i := 0; i < lBaseRunes; i++ {

		if isEndNumericDigits == true {
			continue
		}

		if baseRunes[i] == '-' && isStartNumericDigits == false {
			hasMinusSign = true
			continue
		}

		if baseRunes[i] == '(' && isStartNumericDigits == false {
			hasLeftParen = true
			continue
		}

		if baseRunes[i] == ')' &&
			isStartNumericDigits == true &&
			hasLeftParen == true {

			hasRightParen = true
			isEndNumericDigits = true

		}

		if baseRunes[i] == decimalSeparator {
			startFractionalDigits = true
			continue
		}

		if baseRunes[i] >= '0' &&
			baseRunes[i] <= '9' {

			isStartNumericDigits = true

			newAbsBigInt.Mul(newAbsBigInt, baseTen)

			newAbsBigInt.Add(newAbsBigInt, big.NewInt(int64(baseRunes[i]-48)))

			numOfNumericDigits++

			if startFractionalDigits == true {
				newPrecision++
			}

			continue
		}

	}

	if numOfNumericDigits == 0 {
		return fmt.Errorf(ePrefix+
			"Error: No numeric digits were found in input parameter 'numStr'. "+
			"numStr='%v'", numStr)
	}

	if hasMinusSign == true || (hasLeftParen == true && hasRightParen == true) {

		newAbsBigInt.Neg(newAbsBigInt)

	}

	bigIFd.integerNum.Set(newAbsBigInt)
	bigIFd.precision = newPrecision

	return nil
}

// SetNumericValue
// Sets the 'integerNum' and 'precision' values for the
// current BigIntFixedDecimal instance. Taken together,
// 'integerNum' and 'precision' describe a numeric value
// with a fixed number of fractional digits to the right
// of the decimal place.
func (bigIFd *BigIntFixedDecimal) SetNumericValue(integer *big.Int, precision uint) {

	if integer == nil {
		bigIFd.integerNum = big.NewInt(0)
	} else {
		bigIFd.integerNum = big.NewInt(0).Set(integer)
	}

	bigIFd.precision = precision

}

// SetPrecision - Sets the 'precision' value for the current BigIntFixedDecimal
// instance. 'precision' specifies the number of fractional digits to the right
// of the decimal place.
func (bigIFd *BigIntFixedDecimal) SetPrecisionValue(precision uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
	}

	bigIFd.precision = precision

}

// TrimTrailingFracZeros - This method will delete non-significant
// trailing zeros from the fractional digits of the current
// BigIntFixedDecimal numerical value.
//
// Examples:
//
//	Initial Value			Trimmed Value
//		456.123000 			 456.123
//			0.000					 0
//			7.0						 7
//	 -456.123000			-456.123
func (bigIFd *BigIntFixedDecimal) TrimTrailingFracZeros() {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
		return
	}

	if bigIFd.precision == 0 {
		return
	}

	// bigIFd.precision must be GREATER THAN ZERO
	// Delete trailing fractional zeros

	scrap := big.NewInt(0)
	biBase10 := big.NewInt(10)
	biBaseZero := big.NewInt(0)
	newintegerNum, mod10 := big.NewInt(0).QuoRem(bigIFd.integerNum, biBase10, scrap)

	for mod10.Cmp(biBaseZero) == 0 && bigIFd.precision > 0 {
		bigIFd.integerNum.Set(newintegerNum)
		bigIFd.precision--
		newintegerNum, mod10 = big.NewInt(0).QuoRem(bigIFd.integerNum, biBase10, scrap)
	}

}

// TruncToDecPlace - Truncates the current BigIntFixedDecimal to the
// number of decimal places specified by input parameter 'precision'.
// No rounding occurs, the trailing digits are simply truncated or
// deleted in order to achieve the specified number of decimal places.
//
// 'precision' equals the number of digits to the right of the decimal
// place.
//
// If the value of BigIntFixedDecimal.integerNum is zero ('0'), that
// zero value will remain unaltered. However, BigIntFixedDecimal.precision
// will be set equal to input parameter, 'precision'.
//
// If the number of decimal places specified for truncation ('precision") is
// equal to the current BigIntFixedDecimal.precision, no action is taken and
// the original BigIntFixedDecimal numeric value remains unchanged.
//
// If the number of decimal places specified for truncation ('precision') is
// greater than the current BigIntFixedDecimal.precision, trailing zeros
// are added to the current BigIntFixedDecimal.integerNum value and
// BigIntFixedDecimal.precision is set equal to input parameter, 'precision'.
//
// If 'precision' is less than the current BigIntFixedDecimal.precision
// value, the current BigIntFixedDecimal numeric value is truncated to
// the specified 'precision' value and NO rounding occurs.
//
// Examples:
//
//		 Original       			'newPrecision'				Resulting
//	   Value								input parameter			  Value
//	 --------------				---------------     -------------
//		654.123456									9							 654.123456000
//		654.123456									4							 654.1234 (no rounding)
//
// -654.123456									9							-654.123456000
// -654.123456									4							-654.1234 (no rounding)
//
//			0													3								 0.000
//	   0.000000									0								 0
func (bigIFd *BigIntFixedDecimal) TruncToDecPlace(precision uint) {

	if bigIFd.integerNum == nil {
		bigIFd.integerNum = big.NewInt(0)
		bigIFd.precision = 0
	}

	if bigIFd.precision == precision {
		// Nothing to do. Specified 'precision' is already implemented.
		return
	}

	// bigInt == zero, set precision an return
	if bigIFd.integerNum.Cmp(big.NewInt(0)) == 0 {
		bigIFd.precision = precision
		return
	}

	scale := big.NewInt(0)
	big10 := big.NewInt(10)
	delta := uint(0)

	// If existing precision is less than new specified precision,
	// add trailing zeros, set new precision parameter and return.
	if bigIFd.precision < precision {
		delta = precision - bigIFd.precision
		scale.Exp(
			big10,
			big.NewInt(int64(delta)),
			nil)
		bigIFd.integerNum.Mul(bigIFd.integerNum, scale)
		bigIFd.precision += delta
		return
	}

	// Must be bigIFd.precision > precision
	delta = bigIFd.precision - precision
	scale.Exp(big10, big.NewInt(int64(delta)), nil)
	bigIFd.integerNum.Quo(bigIFd.integerNum, scale)
	bigIFd.precision = precision

}
