package mathops

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

//	BigIntNum
//
//	Wraps a *big.Int integer and its associated precision
//	and sign Value. While the numeric value is stored as
//	an integer of type *big.Int, the BigIntNum type is
//	capable of storing decimal fractions.
//
//	All methods associated with this type all assume that
//	the *big.Int value stored by the BigIntNum Type is
//	configured in base 10.
//
//	INumMgr
//
//	The BigIntNum Type implements the INumMgr interface.
//
//	Source Code Repository:
//
//		https://github.com/MikeAustin71/mathopsgo.git
//
//

type BigIntNum struct {
	bigInt *big.Int
	//	The integer number which consists of
	//	both integer and fractional digits.
	//	This is a signed integer value.

	absBigInt *big.Int
	//	The absolute value of member variable
	//	'bigInt'.

	precision uint
	// Number of digits to the right of the decimal place.

	scaleFactor *big.Int
	// Scale Factor =  10^(precision)

	numberOfExpectedDigits *big.Int
	// Number of digits in the 'absBigInt' value

	sign int
	//	Valid values are -1 or +1.
	//	Indicates the sign of the
	//	'bigInt' integer.

	decimalSeparator rune
	// Character used to separate integer and fractional digits ('.')

	thousandsSeparator rune
	// Character used to separate thousands (1,000,000,000

	currencySymbol rune
	// Currency Symbol
}

// Ceiling - Returns the ceiling integer value of the current BigIntNum
// instance.
//
// Ceiling is defined as: The least, or lowest value integer, which is greater
// than or equal to the numeric value of the current BigIntNum.
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
func (bNum *BigIntNum) Ceiling() BigIntNum {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.IsZero() {
		return new(BigIntNum).NewBigInt(big.NewInt(0), 0)
	}

	if bNum.precision == 0 {
		return bNum.CopyOut()
	}

	scaleVal := big.NewInt(0).Exp(big.NewInt(10),
		big.NewInt(int64(bNum.precision)), nil)

	absQuotient := big.NewInt(0).Quo(bNum.absBigInt, scaleVal)

	// absQuotient IS NOT EQUAL TO bNum.absBigInt

	if bNum.sign > 0 {
		// bNum is positive
		absQuotient = big.NewInt(0).Add(absQuotient, big.NewInt(1))
		return new(BigIntNum).NewBigInt(absQuotient, 0)
	}

	// bNum is negative
	return new(BigIntNum).NewBigInt(big.NewInt(0).Neg(absQuotient), 0)
}

// ChangeSign - Changes the sign of the current BigIntNum value.
//
// If the value of BigIntNum is zero, the sign will remain unchanged
// and this method will return with no action taken.
//
// If the sign of the current BigIntNum value is positive (+), the sign
// will be changed to negative (-). Likewise, if the current sign is
// negative (-), the sign will be changed to positive (+).
func (bNum *BigIntNum) ChangeSign() {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.IsZero() {
		bNum.sign = 1
		return
	}

	bNum.bigInt = big.NewInt(0).Neg(bNum.bigInt)

	if bNum.bigInt.Cmp(big.NewInt(0)) == -1 {
		bNum.sign = -1
	} else {
		bNum.sign = 1
	}

	return
}

// Cmp
//
// Performs a comparison of two BigIntNum numeric
// values and returns an integer value indicating
// the relationship between the two numeric values
// (i.e. greater than, less than, or equal).
//
// Note: Unlike method CmpBigInt() below, this method
// does more than just compare the root *big.Int. In
// making this comparison, the method takes into
// account, numeric sign values and precision.
// Therefore, this method effectively compares numeric
// values. As such, this method provides a true and
// comprehensive picture of the relationship between
// two BigIntNum values.
//
// Return Values:
//
//	bNum == bigIntNum 	Returns  0
//	bNum > bigIntNum	Returns +1
//	bNum < bigIntNum	Returns -1
func (bNum *BigIntNum) Cmp(bigIntNum BigIntNum) int {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	_ = bigIntNum.IsValid("")

	b1Sign := bNum.GetSign()

	b2Sign := bigIntNum.GetSign()

	if b1Sign != b2Sign {
		return b1Sign
	}

	// The signs must be equal

	difference := BigIntMathSubtract{}.SubtractBigIntNums(
		bNum.CopyOut(),
		bigIntNum)

	if difference.IsZero() {
		return 0
	}

	return difference.GetSign()
}

// CmpBigInt - Compares the value of the *big.Int integer to that
// contained in an incoming BigIntNum.
//
// For a true comparison of BigIntNum values see 'Cmp'
//
// Return Values:
// bNum == bigIntNum 				Return  0
// bNum > bigIntNum					Return +1
// bNum < bigIntNum					Return -1
func (bNum *BigIntNum) CmpBigInt(bigIntNum BigIntNum) int {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	_ = bigIntNum.IsValid("")

	return bNum.bigInt.Cmp(bigIntNum.bigInt)

}

// CopyIn - Receives an incoming BigIntNum type and
// copies the value into the current BigIntNum instance.
func (bNum *BigIntNum) CopyIn(bigN BigIntNum) {

	_ = bigN.IsValid("")

	bNum.bigInt = big.NewInt(0).Set(bigN.bigInt)
	bNum.absBigInt = big.NewInt(0).Set(bigN.absBigInt)
	bNum.precision = bigN.precision
	bNum.scaleFactor = big.NewInt(0).Set(bigN.scaleFactor)
	bNum.numberOfExpectedDigits = big.NewInt(0).Set(bigN.numberOfExpectedDigits)
	bNum.sign = bigN.sign
	bNum.decimalSeparator = bigN.decimalSeparator
	bNum.thousandsSeparator = bigN.thousandsSeparator
	bNum.currencySymbol = bigN.currencySymbol
}

// CopyOut - Makes a deep copy of the current BigIntNum instance
// and returns it as a new BigIntNum instance.
func (bNum *BigIntNum) CopyOut() BigIntNum {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	b2 := new(BigIntNum).NewBigInt(big.NewInt(0).Set(bNum.bigInt), bNum.precision)

	b2.decimalSeparator = bNum.decimalSeparator

	b2.thousandsSeparator = bNum.thousandsSeparator

	b2.currencySymbol = bNum.currencySymbol

	b2.numberOfExpectedDigits = big.NewInt(0).Set(bNum.numberOfExpectedDigits)

	return b2
}

// Decrement - Subtracts a value of +1 (plus one) from the numeric
// value of the current BigIntNum instance.
//
// The numeric separators (decimal separator, thousands separator
// and currency symbol) from the original BigIntNum will remain
// unchanged.
func (bNum *BigIntNum) Decrement() {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	biNumOne := new(BigIntNum).NewOne(bNum.precision)

	bPair := BigIntPair{}.NewBigIntNum(bNum.CopyOut(), biNumOne)

	result := BigIntMathSubtract{}.SubtractPair(bPair)

	bNum.CopyIn(result)
}

// Divide - Performs a division operation. The current BigIntNum instance is the 'dividend'
// is divided by the input parameter, 'divisor'. The result of this division operation is
// the 'fracQuotient' which is returned as a BigIntNum type.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
//	bNum = dividend
//	dividend / divisor = quotient
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the resulting
// quotient. Precision is defined as the number of fractional digits to the right of
// the decimal place. Be advised that these calculations can support very large precision
// values.
//
// This returned BigIntNum 'fracQuotient' will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from the current BigIntNum instance (bNum).
func (bNum *BigIntNum) Divide(
	divisor BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	ePrefix := "BigIntNum.Divide() "

	if divisor.IsZero() {
		return BigIntNum{}, errors.New(ePrefix + "ERROR: Attempted Divide by zero!")
	}

	return BigIntMathDivide{}.BigIntNumFracQuotient(bNum.CopyOut(), divisor, maxPrecision)
}

// DivideByFive - Divides the numerical value of the current BigIntNum by five ('5'). The
// result of this division operation is the 'fracQuotient' which is returned as a BigIntNum
// type.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
//	bNum / 5 = fracQuotient
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the resulting
// fractional quotient. Precision is defined as the number of fractional digits to the
// right of the decimal place. Be advised that these calculations can support very large
// precision values.
//
// This returned BigIntNum 'fracQuotient' will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from the current BigIntNum instance (bNum).
func (bNum *BigIntNum) DivideByFive(
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	ePrefix := "BigIntNum.DivideByFive() "

	var errs error

	fracQuotient, errs =
		BigIntMathDivide{}.BigIntNumDivideByFiveFracQuo(bNum.CopyOut(), maxPrecision)

	if errs != nil {

		fracQuotient = new(BigIntNum).New()

		err = fmt.Errorf(ePrefix+"Error returned by "+
			"BigIntMathDivide{}.BigIntNumDivideByFiveFracQuo(bNum.CopyOut(), maxPrecision) "+
			"bNum='%v' Error='%v' \n",
			bNum.GetNumStr(), errs.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// DivideByTen - Divides the numerical value of the current BigIntNum by ten ('10'). The
// result of this division operation is the 'fracQuotient' which is returned as a BigIntNum
// type.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
//	bNum / 10 = fracQuotient
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the resulting
// fractional quotient. Precision is defined as the number of fractional digits to the
// right of the decimal place. Be advised that these calculations can support very large
// precision values.
//
// This returned BigIntNum 'fracQuotient' will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from the current BigIntNum instance (bNum).
func (bNum *BigIntNum) DivideByTen(
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	ePrefix := "BigIntNum.DivideByTen() "

	var errs error

	fracQuotient, errs =
		BigIntMathDivide{}.BigIntNumDivideByTenFracQuo(bNum.CopyOut(), maxPrecision)

	if errs != nil {
		fracQuotient = new(BigIntNum).New()

		err = fmt.Errorf(ePrefix+"Error returned by "+
			"BigIntMathDivide{}.BigIntNumDivideByTenFracQuo(bNum.CopyOut(), maxPrecision) "+
			"bNum='%v' Error='%v' \n",
			bNum.GetNumStr(), errs.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// DivideByTenToPower - Divides the numerical value of the current BigIntNum
// instance by 10 raised to the power of the input parameter, 'exponent'.
//
//	bNum = bNum / (10^exponent)
//
// The original value of the current BigIntNum will be destroyed and overwritten
// by this method.
//
// The final BigIntNum will retain the original numeric separators (decimal separator,
// thousands separator and currency symbol).
func (bNum *BigIntNum) DivideByTenToPower(exponent uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	newPrecision := bNum.precision + exponent

	result := new(BigIntNum).NewBigInt(bNum.bigInt, newPrecision)

	bNum.CopyIn(result)
}

// DivideByThree - Divides the numerical value of the current BigIntNum by three ('3').
// The result of this division operation is the 'fracQuotient' which is returned as a
// BigIntNum type.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
//	bNum / 3 = fracQuotient
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the resulting
// fractional quotient. Precision is defined as the number of fractional digits to the
// right of the decimal place. Be advised that these calculations can support very large
// precision values.
//
// This returned BigIntNum 'fracQuotient' will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from the current BigIntNum instance (bNum).
func (bNum *BigIntNum) DivideByThree(
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntNum.DivideByThree() "

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	var errs error

	fracQuotient, errs =
		BigIntMathDivide{}.BigIntNumDivideByThreeFracQuo(bNum.CopyOut(), maxPrecision)

	if errs != nil {

		fracQuotient = new(BigIntNum).New()

		err = fmt.Errorf(ePrefix+"Error returned by "+
			"BigIntMathDivide{}.BigIntNumDivideByThreeFracQuo(bNum.CopyOut(), maxPrecision) "+
			"bNum='%v' Error='%v' \n",
			bNum.GetNumStr(), errs.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// DivideByTwo - Divides the numerical value of the current BigIntNum by two ('2'). The
// result of this division operation is the 'fracQuotient' which is returned as a BigIntNum
// type.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
// The return value may be an integer or a floating point value as determined by the results
// of the division operation. For floating point values, the number of digits to the right
// of the decimal place is limited by input parameter, 'maxPrecision'.
//
//	bNum / 2 = fracQuotient
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the resulting
// fractional quotient. Precision is defined as the number of fractional digits to the
// right of the decimal place. Be advised that these calculations can support very large
// precision values.
func (bNum *BigIntNum) DivideByTwo(
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntNum.DivideByTwo() "

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	var errs error

	fracQuotient, errs =
		BigIntMathDivide{}.BigIntNumDivideByTwoFracQuo(bNum.CopyOut(), maxPrecision)

	if errs != nil {

		fracQuotient = new(BigIntNum).New()

		err = fmt.Errorf(ePrefix+"Error returned by "+
			"BigIntMathDivide{}.BigIntNumDivideByTwoFracQuo(bNum.CopyOut(), maxPrecision) "+
			"bNum='%v' Error='%v' \n",
			bNum.GetNumStr(), errs.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// DivideByTwoQuoMod
//
// Divides the numerical value of the current BigIntNum
// by two ('2'). The result of this division operation
// is returned as an integer quotient, 'intQuotient',
// and a floating point modulo or remainder, 'modulo'.
//
// bNum / 2 = integer quotient and floating point modulo
//
// If 'modulo' equals zero ('0'), it signals the current
// BigIntNum numerical value is 'even'; that is, it is
// evenly divisible by two.
//
// Input parameter 'maxPrecision' is used to control the
// maximum precision of the resulting floating point
// 'modulo'. Precision is defined as the number of
// fractional digits to the right of the decimal place.
// Be advised that these calculations can support very
// large precision values.
//
// Both returned BigIntNum 'intQuotient' and 'modulo'
// BigIntNum types will contain numeric separators
// (decimal separator, thousands separator and
// currency symbol) copied from the current BigIntNum
// instance (bNum).
func (bNum *BigIntNum) DivideByTwoQuoMod(
	maxPrecision uint) (intQuotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntNum) DivideByTwoQuoMod() "
	var errs error

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	intQuotient, modulo, errs =
		BigIntMathDivide{}.BigIntNumDivideByTwoQuoMod(bNum.CopyOut(), maxPrecision)

	if errs != nil {

		intQuotient = new(BigIntNum).New()

		modulo = new(BigIntNum).New()

		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}."+
			"BigIntNumDivideByTwoQuoMod(bNum.CopyOut(), maxPrecision). "+
			"bNum='%v' Error='%v'\n",
			bNum.GetNumStr(), errs.Error())

		return intQuotient, modulo, err
	}

	err = nil

	return intQuotient, modulo, err
}

// Empty - Resets the BigIntNum data fields to their
// uninitialized or zero state.
func (bNum *BigIntNum) Empty() {

	bNum.bigInt = big.NewInt(0)
	bNum.absBigInt = big.NewInt(0)
	bNum.scaleFactor = big.NewInt(1)
	bNum.numberOfExpectedDigits = big.NewInt(0)
	bNum.sign = 1
	bNum.precision = 0

	bNum.SetNumericSeparatorsToDefaultIfEmpty()

}

// Equal - Compares two BigIntNum instances and returns 'true'
// if the two instances are equal in all respects.
//
// Be careful, two BigIntNum instances could have equal
// values with different precisions. In that case this
// method would return 'false'. To test for equivalent
// values, see method BigIntNum.EqualValue(), below.
//
// If they are not Equal, the method returns 'false'.
func (bNum *BigIntNum) Equal(b2 BigIntNum) bool {

	if bNum.bigInt.Cmp(b2.bigInt) != 0 {
		return false
	}

	if bNum.absBigInt.Cmp(b2.absBigInt) != 0 {
		return false
	}

	if bNum.scaleFactor.Cmp(b2.scaleFactor) != 0 {
		return false
	}

	if bNum.sign != b2.sign {
		return false
	}

	if bNum.precision != b2.precision {
		return false
	}

	return true
}

// EqualValue - Compares the values of the current BigIntNum instance
// and the input parameter BigIntNum, 'b2'. If the two numeric values
// are equal, this method returns 'true'.
func (bNum *BigIntNum) EqualValue(b2 BigIntNum) bool {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	_ = b2.IsValid("")

	difference := BigIntMathSubtract{}.SubtractBigIntNums(
		bNum.CopyOut(),
		b2)

	if difference.IsZero() {
		return true
	}

	return false
}

// ExtendPrecision - Extends the current precision.
//
// Precision is the number of fractional digits to the right
// of the decimal place. This method will extend the number of
// digits to the right of the decimal place by adding trailing
// zeros to the current numeric value of this 'BigIntNum' instance.
// The number of trailing zeros to be added is determined by the
// input parameter, 'deltaPrecision'.
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) ExtendPrecision(deltaPrecision uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if deltaPrecision == 0 {
		return
	}

	bNum.SetNumericSeparatorsToDefaultIfEmpty()
	numSeps := bNum.GetNumericSeparatorsDto()

	newPrecision := bNum.precision + deltaPrecision

	// bigInt == zero, set precision an return
	if bNum.bigInt.Cmp(big.NewInt(0)) == 0 {

		bNum.CopyIn(
			new(BigIntNum).NewBigInt(big.NewInt(0),
				newPrecision))

		_ = bNum.SetNumericSeparatorsDto(numSeps)

		return
	}

	base10 := big.NewInt(10)
	scaleVal := big.NewInt(0).Exp(base10, big.NewInt(int64(deltaPrecision)), nil)
	bigINum := big.NewInt(0).Set(bNum.bigInt)

	bigINum = big.NewInt(0).Mul(bigINum, scaleVal)

	bNum.CopyIn(new(BigIntNum).NewBigInt(bigINum, newPrecision))

	_ = bNum.SetNumericSeparatorsDto(numSeps)
}

// Floor
//
// Returns the greatest integer less than or equal to the
// numeric value of the current BigIntNum.
//
// Reference Wikipedia:
// https://en.wikipedia.org/wiki/Floor_and_ceiling_functions
//
//	Initial		Floor
//	 Value		Value
//	-------		-------
//	 5.95		  5
//	 5.05		  5
//	 5			  5
//	-5.05		 -6
//	 2.4		  2
//	 2.9		  2
//	-2.7		 -3
//	-2			 -2
func (bNum *BigIntNum) Floor() BigIntNum {

	if bNum.bigInt == nil {

		bNum.SetBigInt(
			big.NewInt(0),
			bNum.precision)
	}

	if bNum.IsZero() {
		return new(BigIntNum).NewBigInt(big.NewInt(0), 0)
	}

	if bNum.precision == 0 {
		return bNum.CopyOut()
	}

	scaleVal := big.NewInt(0).Exp(big.NewInt(10),
		big.NewInt(int64(bNum.precision)), nil)

	absQuotient := big.NewInt(0).Quo(bNum.absBigInt, scaleVal)

	// absQuotient IS NOT EQUAL TO bNum.absBigInt

	if bNum.sign > 0 {
		// bNum is positive
		return new(BigIntNum).NewBigInt(absQuotient, 0)
	}

	// bNum is negative
	absQuotient = big.NewInt(0).Add(absQuotient, big.NewInt(1))

	return new(BigIntNum).NewBigInt(
		big.NewInt(0).Neg(absQuotient),
		0)
}

//	FormatCurrencyStr
//
//	Formats the current BigIntNum numeric value as a currency
//	string.
//
//	If the Currency Symbol was not previously set for this
//	BigIntNum, the currency symbol is defaulted to the US
//	standard dollar sign, ('$').
//
//	To use other currency symbols, see method:
//		BigIntNum.SetCurrencySymbol().
//
//	If the Decimal Separator was not previously set for this
//	BigIntNum, the Decimal Separator is defaulted to the US
//	standard period ('.'). To use another character for
//	Decimal Separator, see method:
//		BigIntNum.SetDecimalSeparator().
//
//	If the Thousands Separator was not previously set for
//	this BigIntNum, the Thousands Separator is defaulted to
//	the US standard comma (','). To use another character for
//	Thousands Separator, see method:
//		BigIntNum.SetThousandsSeparator().
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	negValMode			NegativeValueFmtMode
//		Specifies the display mode for negative values:
//
//		LEADMINUSNEGVALFMTMODE
//			Negative values formatted with a leading minus
//			sign.
//				Example: -$123,456.78
//
//		PARENTHESESNEGVALFMTMODE
//			Negative values formatted with surrounding
//			parentheses.
//				Example: ($123,456.78)
//
//
//		ABSOLUTEPURENUMSTRFMTMODE
//			Formats a pure number string with absolute
//			(positive) integer value and no decimal place separator.
//				Example: ($12,345,678)
func (bNum *BigIntNum) FormatCurrencyStr(negValMode NegativeValueFmtMode) string {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.decimalSeparator == 0 {
		bNum.decimalSeparator = '.'
	}

	if bNum.thousandsSeparator == 0 {
		bNum.thousandsSeparator = ','
	}

	if bNum.currencySymbol == 0 {
		bNum.currencySymbol = '$'
	}

	outRunes := make([]rune, 0, 300)

	scratchNum := big.NewInt(0).Set(bNum.absBigInt)
	baseZero := big.NewInt(0)

	if scratchNum.Cmp(baseZero) == 0 {
		bNum.sign = 1

		outRunes = append(outRunes, bNum.currencySymbol)

		outRunes = append(outRunes, '0')

		if bNum.precision > 0 {

			if negValMode != ABSOLUTEPURENUMSTRFMTMODE {
				outRunes = append(outRunes, bNum.decimalSeparator)
			}

			cnt := int(bNum.precision)

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
	thouCnt := -1

	if bNum.precision == 0 {
		thouCnt = 0
	}

	for scratchNum.Cmp(baseZero) == 1 {

		if startIdx == 0 &&
			bNum.sign == -1 &&
			negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, ')')
		}

		modX := big.NewInt(0)
		scratchNum, modulo = big.NewInt(0).QuoRem(scratchNum, baseTen, modX)
		outRunes = append(outRunes, rune(modulo.Int64()+int64(48)))
		digitCnt++
		startIdx++

		if thouCnt > -1 {
			thouCnt++
		}

		if scratchNum.Cmp(baseZero) == 1 &&
			thouCnt == 3 {

			outRunes = append(outRunes, bNum.thousandsSeparator)
			startIdx++
			thouCnt = 0
		}

		if bNum.precision > 0 &&
			int(bNum.precision) == startIdx &&
			negValMode != ABSOLUTEPURENUMSTRFMTMODE {

			outRunes = append(outRunes, bNum.decimalSeparator)
			startIdx++
			thouCnt = 0
		}

	}

	if int(bNum.precision) >= digitCnt {

		delta := int(bNum.precision) - digitCnt + 1

		if negValMode == ABSOLUTEPURENUMSTRFMTMODE {
			delta--
		}

		for k := 0; k < delta; k++ {
			outRunes = append(outRunes, '0')
			startIdx++

			if bNum.precision > 0 &&
				int(bNum.precision) == startIdx &&
				negValMode != ABSOLUTEPURENUMSTRFMTMODE {

				outRunes = append(outRunes, bNum.decimalSeparator)
				startIdx++

			}
		}
	}

	startIdx--

	// append Currency Symbol
	outRunes = append(outRunes, bNum.currencySymbol)
	startIdx++

	// adjust for negative sign value
	if bNum.sign == -1 {

		if negValMode == LEADMINUSNEGVALFMTMODE {
			outRunes = append(outRunes, '-')
			startIdx++

		} else if negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, '(')
			startIdx += 2
		}

		// Must be negValMode == ABSOLUTEPURENUMSTRFMTMODE

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

// FormatNumStr - Formats the numeric value of the current BigIntNum
// instance as number string consisting of integer digits to the left
// of the decimal place and fractional digits to the right of the decimal
// point, if such fractional digits exist. The resulting number string
// will NOT contain a currency symbol or thousands separators.
//
// If the Decimal Separator was not previously set for this BigIntNum,
// the Decimal Separator is defaulted to the USA standard period ('.').
// To use another character for Decimal Separator, see method
// BigIntNum.SetDecimalSeparator().
//
// Output Examples: 123456.789 or -123456.789
//
// Input Parameters
// ================
//
// negValMode NegativeValueFmtMode -	Specifies the display mode for negative values:
//
//	LEADMINUSNEGVALFMTMODE 		-	Negative values formatted with
//													 		a leading minus sign.
//															Example: -123456.78
//
//	PARENTHESESNEGVALFMTMODE	-	Negative values formatted with
//															surrounding parentheses.
//															Example: (123456.78)
//
//	ABSOLUTEPURENUMSTRFMTMODE - Formats a pure number string with
//															absolute (positive) integer value
//															and no decimal place separator.
//															Example: (12345678)
func (bNum *BigIntNum) FormatNumStr(negValMode NegativeValueFmtMode) string {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.decimalSeparator == 0 {
		bNum.decimalSeparator = '.'
	}

	outRunes := make([]rune, 0, 300)

	scratchNum := big.NewInt(0).Set(bNum.absBigInt)
	baseZero := big.NewInt(0)

	if scratchNum.Cmp(baseZero) == 0 {
		bNum.sign = 1

		outRunes = append(outRunes, '0')

		if bNum.precision > 0 {

			if negValMode != ABSOLUTEPURENUMSTRFMTMODE {
				outRunes = append(outRunes, bNum.decimalSeparator)
			}

			cnt := int(bNum.precision)

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
			bNum.sign == -1 &&
			negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, ')')
		}

		modX := big.NewInt(0)
		scratchNum, modulo = big.NewInt(0).QuoRem(scratchNum, baseTen, modX)
		outRunes = append(outRunes, rune(modulo.Int64()+int64(48)))
		digitCnt++
		startIdx++

		if bNum.precision > 0 &&
			int(bNum.precision) == startIdx &&
			negValMode != ABSOLUTEPURENUMSTRFMTMODE {

			outRunes = append(outRunes, bNum.decimalSeparator)
			startIdx++
		}

	}

	if int(bNum.precision) >= digitCnt {

		delta := int(bNum.precision) - digitCnt + 1

		if negValMode == ABSOLUTEPURENUMSTRFMTMODE {
			delta--
		}

		for k := 0; k < delta; k++ {
			outRunes = append(outRunes, '0')
			startIdx++

			if bNum.precision > 0 &&
				int(bNum.precision) == startIdx &&
				negValMode != ABSOLUTEPURENUMSTRFMTMODE {

				outRunes = append(outRunes, bNum.decimalSeparator)
				startIdx++
			}
		}
	}

	startIdx--

	// adjust for negative sign value
	if bNum.sign == -1 {

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

// FormatThousandsStr - Returns the number string delimited with the
// BigIntNum ThousandsSeparator character plus the Decimal Separator
// character if applicable. See methods BigIntNum.SetThousandsSeparator()
// and BigIntNum.SetDecimalSeparator().
//
// If the Decimal Separator was not previously set for this BigIntNum,
// the Decimal Separator is defaulted to the USA standard period ('.').
// To use another character for Decimal Separator, see method
// BigIntNum.SetDecimalSeparator().
//
// If the Thousands Separator was not previously set for this BigIntNum,
// the Thousands Separator is defaulted to the USA standard comma (',').
// To use another character for Thousands Separator, see method
// BigIntNum.SetThousandsSeparator().
//
// Example:
// numStr = 1000000.234 converted to 1,000,000.234
//
// Input Parameters
// ================
//
// Input Parameters
// ================
//
// negValMode NegativeValueFmtMode -	Specifies the display mode for negative values:
//
//	LEADMINUSNEGVALFMTMODE 		-	Negative values formatted with
//													 		a leading minus sign.
//															Example: -123,456.78
//
//	PARENTHESESNEGVALFMTMODE	-	Negative values formatted with
//															surrounding parentheses.
//															Example: (123,456.78)
//
//
//	ABSOLUTEPURENUMSTRFMTMODE - Formats a pure number string with
//															absolute (positive) integer value
//															and no decimal place separator.
//															Example: (12,345,678)
func (bNum *BigIntNum) FormatThousandsStr(negValMode NegativeValueFmtMode) string {

	if bNum.decimalSeparator == 0 {
		bNum.decimalSeparator = '.'
	}

	if bNum.thousandsSeparator == 0 {
		bNum.thousandsSeparator = ','
	}

	outRunes := make([]rune, 0, 300)

	scratchNum := big.NewInt(0).Set(bNum.absBigInt)
	baseZero := big.NewInt(0)

	if scratchNum.Cmp(baseZero) == 0 {
		bNum.sign = 1

		outRunes = append(outRunes, '0')

		if bNum.precision > 0 {

			if negValMode != ABSOLUTEPURENUMSTRFMTMODE {
				outRunes = append(outRunes, bNum.decimalSeparator)
			}

			cnt := int(bNum.precision)

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
	thouCnt := -1

	if bNum.precision == 0 {
		thouCnt = 0
	}

	for scratchNum.Cmp(baseZero) == 1 {

		if startIdx == 0 &&
			bNum.sign == -1 &&
			negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, ')')
		}

		modX := big.NewInt(0)
		scratchNum, modulo = big.NewInt(0).QuoRem(scratchNum, baseTen, modX)
		outRunes = append(outRunes, rune(modulo.Int64()+int64(48)))
		digitCnt++
		startIdx++

		if thouCnt > -1 {
			thouCnt++
		}

		if scratchNum.Cmp(baseZero) == 1 &&
			thouCnt == 3 {

			outRunes = append(outRunes, bNum.thousandsSeparator)
			startIdx++
			thouCnt = 0
		}

		if bNum.precision > 0 &&
			int(bNum.precision) == startIdx &&
			negValMode != ABSOLUTEPURENUMSTRFMTMODE {

			outRunes = append(outRunes, bNum.decimalSeparator)
			startIdx++
			thouCnt = 0
		}

	}

	if int(bNum.precision) >= digitCnt {

		delta := int(bNum.precision) - digitCnt + 1

		if negValMode == ABSOLUTEPURENUMSTRFMTMODE {
			delta--
		}

		for k := 0; k < delta; k++ {
			outRunes = append(outRunes, '0')
			startIdx++

			if bNum.precision > 0 &&
				int(bNum.precision) == startIdx &&
				negValMode != ABSOLUTEPURENUMSTRFMTMODE {

				outRunes = append(outRunes, bNum.decimalSeparator)
				startIdx++
			}
		}
	}

	startIdx--

	// adjust for negative sign value
	if bNum.sign == -1 {

		if negValMode == LEADMINUSNEGVALFMTMODE {
			outRunes = append(outRunes, '-')
			startIdx++

		} else if negValMode == PARENTHESESNEGVALFMTMODE {

			outRunes = append(outRunes, '(')
			startIdx += 2
		}

		// Must Be negValMode == ABSOLUTEPURENUMSTRFMTMODE

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

// GetActualNumberOfDigits - Returns the number of numeric digits
// in the absolute value of this BigIntNum instance. In addition,
// a boolean value is returned indicating whether the absolute value
// is zero.
//
// Examples
// ========
//
//	      123.45														5
//	1,234,567                               7
//
// -1,234,567                               7
//
//						0																1
//	         0.00                            1
//	       012.34                            4
//	         0.1234													4
//	       - 0.1234													4
//	         0.123400												4
//	         0.0123400												4
//	 1,234,567.800													  8
//	         5                               1
func (bNum *BigIntNum) GetActualNumberOfDigits() (
	numberOfDigits *big.Int, isZeroValue bool, err error) {

	numberOfDigits = big.NewInt(0)
	isZeroValue = false
	err = nil

	numOfDigits, errs := BigIntMath{}.GetMagnitude(bNum.absBigInt)

	if errs != nil {
		ePrefix := "BigIntNum.GetActualNumberOfDigits() "
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMath{}.GetMagnitudeDigits(bNum.absBigInt) "+
			"bNum.absBigInt='%v' Error='%v' ", bNum.absBigInt.Text(10), err.Error())
		return numberOfDigits, isZeroValue, err
	}

	numberOfDigits = big.NewInt(0).Add(numOfDigits, big.NewInt(1))

	if bNum.absBigInt.Cmp(big.NewInt(0)) == 0 {
		isZeroValue = true
	}

	return numberOfDigits, isZeroValue, err
}

// GetAbsoluteNumStr - Returns the absolute integer value (positive value) of the
// *big.Int value encapsulated by this BigIntNum. No decimal place is included.
//
// If an error is encountered, an empty string is returned.
func (bNum *BigIntNum) GetAbsoluteNumStr() string {

	if bNum.GetSign() == 1 {
		return bNum.FormatNumStr(LEADMINUSNEGVALFMTMODE)
	}

	biNum := bNum.CopyOut()

	biNum.ChangeSign()

	return biNum.FormatNumStr(LEADMINUSNEGVALFMTMODE)
}

// GetAbsoluteBigIntNumValue
//
// Returns the absolute numeric value of this BigIntNum
// instance as a new BigIntNum Type.
//
// If the current BigIntNum value is '-123.456', this
// method will return '123.456'.
//
// If the current BigIntNum value is '123.456',
// this method will return '123.456'.
func (bNum *BigIntNum) GetAbsoluteBigIntNumValue() BigIntNum {

	return new(BigIntNum).NewBigInt(
		bNum.absBigInt,
		bNum.precision)
}

// GetAbsoluteBigIntValue
//
// Returns the absolute value of the *big.Int value
// encapsulated by the current BigIntNum.
//
// If the current BigIntNum value is '-123456', this
// method will return '123456'.
//
// If the current BigIntNum value is '123.456',
// this method will return '123456'.
func (bNum *BigIntNum) GetAbsoluteBigIntValue() *big.Int {

	return big.NewInt(0).Set(bNum.absBigInt)
}

// GetBigFloat - Returns the numeric value of the current
// BigIntNum as *big.Float type.
func (bNum *BigIntNum) GetBigFloat() *big.Float {

	numerator := big.NewInt(0).Set(bNum.bigInt)

	denominator := big.NewInt(0).Set(bNum.scaleFactor)

	bRat := big.NewRat(1, 1).SetFrac(numerator, denominator)

	return big.NewFloat(0).SetRat(bRat)
}

// GetBigInt - return the numeric value as an integer
// of type *big.int.
func (bNum *BigIntNum) GetBigInt() (*big.Int, error) {

	return big.NewInt(0).Set(bNum.bigInt), nil
}

// GetBigIntFixedDecimal
//
// Returns a BigIntFixedDecimal instance which contains a
// deep copy of the current BigIntNum values.
func (bNum *BigIntNum) GetBigIntFixedDecimal() BigIntFixedDecimal {

	return new(BigIntFixedDecimal).New(
		bNum.bigInt,
		bNum.precision)

}

// GetBigIntNum
//
// Returns a deep copy of the current BigIntNum instance.
//
// The returned BigIntNum will contain numeric separators
// (decimal separator, integer separator, integer grouping
// and currency symbol) copied from the current BigIntNum
// instance.
//
// Before returning the new BigIntNum copy, this method
// performs a validity test on the current BigIntNum
// instance. If the current BigIntNum instance is invalid,
// an error will be returned.
//
// This method is necessary in order to fulfill the
// requirements of the INumMgr interface.
func (bNum *BigIntNum) GetBigIntNum() (BigIntNum, error) {

	ePrefix := "BigIntNum.GetBigIntNum() "

	err := bNum.IsValid(ePrefix + "BigIntNum INVALID! ")

	if err != nil {
		return new(BigIntNum).New(), err
	}

	return bNum.CopyOut(), nil
}

// GetBigRat
// Returns the numeric value of the current BigIntNum as a
// type *big.Rat (Rational Number).
func (bNum *BigIntNum) GetBigRat() *big.Rat {

	numerator := big.NewInt(0).Set(bNum.bigInt)

	denominator := big.NewInt(0).Set(bNum.scaleFactor)

	return big.NewRat(1, 1).SetFrac(numerator, denominator)
}

// GetCurrencySymbol
//
// Returns the character currently designated as the
// currency symbol for this BigIntNum instance.
//
// In the United States, the currency symbol is the
// dollar sign ('$').
func (bNum *BigIntNum) GetCurrencySymbol() rune {

	if bNum.currencySymbol == 0 {
		bNum.currencySymbol = '$'
	}

	return bNum.currencySymbol

}

// GetDecimalSeparator - returns the character designated
// as the decimal separator for the current NumStrDto instance.
//
// In the USA, the decimal separator is the period character ('.').
//
// Example:		123.456
func (bNum *BigIntNum) GetDecimalSeparator() rune {

	if bNum.decimalSeparator == 0 {
		bNum.decimalSeparator = '.'
	}

	return bNum.decimalSeparator

}

// GetExpectedNumberOfDigits - Returns the number of expected numeric
// digits associated with this BigIntNum instance. The returned value
// is stored in data field, BigIntNum.numberOfExpectedDigits. The value
// is set by calling method BigIntNum.SetExpectedNumberOfDigits().
//
// This value is useful in tracking leading zeros.
func (bNum *BigIntNum) GetExpectedNumberOfDigits() *big.Int {
	return bNum.numberOfExpectedDigits
}

// GetFractionalPart
//
// Returns the fractional digits of the current
// BigIntNum instance as a new BigIntNum instance
// containing those correctly formatted fractional digits.
//
// Examples:
//
//	Current
//	BigIntNum		 Return
//	  Value			  Value
//	----------		---------
//	  123.456		  0.456
//	 -123.456		 -0.456
//	  123			  0
//	 -123			  0
func (bNum *BigIntNum) GetFractionalPart() BigIntNum {

	if bNum.IsZero() {
		return new(BigIntNum).NewBigInt(big.NewInt(0), 0)
	}

	scaleVal := big.NewInt(0).Exp(big.NewInt(10),
		big.NewInt(int64(bNum.precision)), nil)

	modulo := big.NewInt(0).Rem(bNum.bigInt, scaleVal)

	return new(BigIntNum).NewBigInt(modulo, bNum.precision)
}

// GetInt - Returns a type 'int' containing the 32-big integer
// value of the current BigIntNum instance.
//
// If the current BigIntNum value is greater than the maximum
// 'int' value, the maximum 32-bit integer value is returned
// in addition to an 'error'.
//
// If the current BigIntNum value is less than the minimum 'int'
// value, the minimum 32-bit integer value is returned along with
// an 'error'.
func (bNum *BigIntNum) GetInt() (int, error) {

	ePrefix := "BigIntNum) GetInt() "
	bIMaxInt := big.NewInt(int64(math.MaxInt32))
	bIMinInt := big.NewInt(int64(math.MinInt32))

	if bNum.bigInt.Cmp(bIMaxInt) == 1 {
		return math.MaxInt32, fmt.Errorf(ePrefix+"Error: BigIntNum Value is GREATER than Int32 Maximum! "+
			"Int32 Maximum Value='%v' BigIntNum Value='%v'", bIMaxInt.Text(10), bNum.GetNumStr())
	}

	if bNum.bigInt.Cmp(bIMinInt) == -1 {
		return math.MinInt32, fmt.Errorf(ePrefix+"Error: BigIntNum Value is LESS than Int32 Minmum! "+
			"Int32 Minimum Value='%v' BigIntNum Value='%v'", bIMinInt.Text(10), bNum.GetNumStr())
	}

	return int(bNum.bigInt.Int64()), nil

}

// GetIntegerPart
//
// Returns a BigIntNum equal to the integer value of
// the current BigIntNum. Decimal digits are truncated
// and not included in the return value.
//
// Examples:
//
//	 Current
//	BigIntNum		Return
//	  Value			Value
//	----------		---------
//
//	 123.456		 123
//	-123.456		-123
//	 123			 123
//	-123			-123
func (bNum *BigIntNum) GetIntegerPart() BigIntNum {

	if bNum.IsZero() {
		return new(BigIntNum).NewBigInt(big.NewInt(0), 0)
	}

	scaleVal := big.NewInt(0).Exp(big.NewInt(10),
		big.NewInt(int64(bNum.precision)), nil)

	quotient := big.NewInt(0).Quo(bNum.bigInt, scaleVal)

	return new(BigIntNum).NewBigInt(quotient, 0)
}

// GetIntegerValue
//
// Returns the internal *big.Int number for the current
// BigIntNum instance.
func (bNum *BigIntNum) GetIntegerValue() *big.Int {

	if bNum.bigInt == nil {
		bNum.bigInt = big.NewInt(0)
	}

	return big.NewInt(0).Set(bNum.bigInt)
}

// GetInverse
//
// Returns the value of one (1) divided by the current
// BigIntNum instance as a new BigIntNum Type.
func (bNum *BigIntNum) GetInverse(maxPrecision uint) (BigIntNum, error) {

	bINumOne := new(BigIntNum).NewOne(0)

	result, err := BigIntMathDivide{}.BigIntNumFracQuotient(bINumOne, bNum.CopyOut(), maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf("BigIntNum.GetInverse() - Error returned by BigIntMathDivide{}."+
				"BigIntNumFracQuotient(bINumOne, bNum.CopyOut(), maxPrecision)"+
				"bNum='%v' Error='%v' ", bNum.GetNumStr(), err.Error())
	}

	return result, nil
}

// GetNumberOfDigits - Returns the number of digits in the numeric value of the
// current BigIntNum instance. The count only includes numeric digits and as such
// EXCLUDES number signs ('-' or '+' ), thousands separators (',') and decimal
// separators.
//
// Examples:
// =========
//
//	Result=
//
// Numeric String                       Number of
//
//	Value								          Numeric Digits
//
// =============                      ==============
//
//	      123.45														5
//	1,234,567                               7
//
// -1,234,567.8															8
//
//						0																1
//	         0.00                            1
//	       012.34                            4
//	         0.1234													4
//	         0.123400												6
//	         0.0123400												6
//	 1,234,567.800													 10
//	         5                               1
//
// Note: The returned integer number will always be a positive number.
//
//	      Also, GetActualNumberOfDigits() will be faster for larger
//				 numbers.
func (bNum *BigIntNum) GetNumberOfDigits() int {

	scratchNum := big.NewInt(0).Set(bNum.absBigInt)
	baseZero := big.NewInt(0)
	digitCnt := 0

	if scratchNum.Cmp(baseZero) == 0 {

		digitCnt = 1

		return digitCnt
	}

	baseTen := big.NewInt(10)

	for scratchNum.Cmp(baseZero) == 1 {
		scratchNum = big.NewInt(0).Quo(scratchNum, baseTen)
		digitCnt++
	}

	return digitCnt
}

// GetNumericSeparatorsDto - Returns a structure containing the
// character or rune values for decimal place separator, thousands
// separator and currency symbol.
func (bNum *BigIntNum) GetNumericSeparatorsDto() NumericSeparatorDto {

	numSeps := NumericSeparatorDto{}
	numSeps.DecimalSeparator = bNum.GetDecimalSeparator()
	numSeps.ThousandsSeparator = bNum.GetThousandsSeparator()
	numSeps.CurrencySymbol = bNum.GetCurrencySymbol()

	return numSeps
}

// GetNumStr - Converts the current BigIntNum value to string of
// numbers which includes the decimal place and decimal digits
// if they exist.
func (bNum *BigIntNum) GetNumStr() string {

	return bNum.FormatNumStr(LEADMINUSNEGVALFMTMODE)

}

// GetNumStrDto - Converts the current BigIntNum value to a NumStrDto
// instance. The resulting number string includes the decimal place
// and decimal digits if they exist.
//
// The returned NumStrDto type contains numeric separators (decimal separator,
// thousands separator and currency symbol) copied from the current BigIntNum
// instance.
//
// This method performs a validity test on the current BigIntNum instance.
func (bNum *BigIntNum) GetNumStrDto() (NumStrDto, error) {

	ePrefix := "BigIntNum.GetNumStrDto() "

	err := bNum.IsValid(ePrefix + "BigIntNum INVALID! ")

	if err != nil {
		return NumStrDto{}, err
	}

	var nDto NumStrDto

	nDto, err =
		NumStrDto{}.NewBigInt(big.NewInt(0).Set(bNum.bigInt),
			bNum.precision)

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by NumStrDto{}.NewBigInt(bNum.bigInt, bNum.precision) "+
				"bNum.bigInt='%v' bNum.precision='%v' Error='%v'",
				bNum.bigInt.Text(10), bNum.precision, err.Error())
	}

	err = nDto.SetNumericSeparatorsDto(bNum.GetNumericSeparatorsDto())

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by nDto.SetNumericSeparatorsDto(bNum.GetNumericSeparatorsDto()). "+
				"Error='%v' \n", err.Error())
	}

	err = nDto.IsValid(ePrefix + "'nDto' INVALID! ")

	if err != nil {
		return NumStrDto{}.New(), err
	}

	return nDto, nil
}

// GetPrecision - Returns precision as an integer of
// type 'int'.
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal place in a string of numeric digits, go
// to the right most digit in the number string and count
// left 'precision' digits.
//
// Example:
//
//					1.234    	GetPrecision() = 3
//							5			GetPrecision() = 0
//				0.12345  		GetPrecision() = 5
//
//	Number String				precision				Fractional Number
//		123456								3								123.456
func (bNum *BigIntNum) GetPrecision() int {
	return int(bNum.precision)
}

// GetPrecisionBigInt - Returns the 'precision' of the current
// BigIntNum instance as a *big.Int Type.
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal place in a string of numeric digits, go
// to the right most digit in the number string and count
// left 'precision' digits.
//
// Example:
//
//		1.234    	GetPrecisionBigInt() = 3
//				5			GetPrecisionBigInt() = 0
//	0.12345  		GetPrecisionBigInt() = 5
func (bNum *BigIntNum) GetPrecisionBigInt() *big.Int {

	return big.NewInt(0).SetUint64(uint64(bNum.precision))

}

// GetPrecisionUint - Returns precision as an unsigned
// integer (uint).
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal place in a string of numeric digits, go
// to the right most digit in the number string and count
// left 'precision' digits.
//
// Example:
//
//					1.234    	GetPrecisionUint() = 3
//							5			GetPrecisionUint() = 0
//				0.12345  		GetPrecisionUint() = 5
//
//	Number String				precision				Fractional Number
//		123456								3								123.456
func (bNum *BigIntNum) GetPrecisionUint() uint {
	return bNum.precision
}

// GetScaleFactor - Returns the scale value of the current
// BigIntNum.  Scale value is a function of 'precision' or
// the number of digits to the right of the decimal place.
// Therefore, scale factor is defined by 10 raised to the
// power BigIntNum precision.
//
// Example:
// precision = 0 		Scale Factor = 10^0   	Scale Factor =    1
// precision = 1		Scale Factor = 10^1			Scale Factor =   10
// precision = 2		Scale Factor = 10^2			Scale Factor =  100
// precision = 3    Scale Factor = 10^3			Scale Factor = 1000
func (bNum *BigIntNum) GetScaleFactor() *big.Int {
	return big.NewInt(0).Set(bNum.scaleFactor)
}

// GetSign - Returns the numeric sign associated
// with the current numeric value encapsulated by
// this BigIntNum.
func (bNum *BigIntNum) GetSign() int {
	return bNum.sign
}

// GetSignedBigInt
//
// Returns the integer value of the current
// BigIntNum as a signed *big.Int Type.
func (bNum *BigIntNum) GetSignedBigInt() *big.Int {
	return bNum.bigInt
}

//	GetSciNotationNumber
//
//	Converts the numeric value of the current BigIntNum
//	instance into scientific notation and returns this
//	value as an instance of type SciNotationNum.
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	mantissaLen			uint
//
//		Specifies the length of the mantissa in the returned
//		scientific notation string.
//
//		If the value of 'mantissaLen' is less than two ('2'),
//		this method will automatically set the 'mantissaLen'
//		to a default value of two ('2').
//
// ----------------------------------------------------------------
//
//	Example Scientific Notation:
//
//		scientific notation string: '2.652e+8'
//		significand = '2.652'
//		significand integer digit = '2'
//		mantissa		= significand factional digits = '.652'
//		exponent    = '8'  (10^8)
func (bNum *BigIntNum) GetSciNotationNumber(mantissaLen uint) (SciNotationNum, error) {

	ePrefix := "BigIntNum.GetSciNotationNumber() "

	sciNotationNum := SciNotationNum{}.New()

	if mantissaLen < 2 {
		mantissaLen = 2
	}

	if bNum.IsZero() {

		_ = sciNotationNum.SetBigIntNumElements(
			bNum.CopyOut(),
			new(BigIntNum).NewZero(0))

		return sciNotationNum, nil
	}

	bigIntMaxUint32 :=
		big.NewInt(0).SetUint64(math.MaxUint32)

	bINumIntPart := bNum.GetIntegerPart()

	if !bINumIntPart.IsZero() {

		magnitudeBigInt, err := BigIntMath{}.GetMagnitude(bINumIntPart.bigInt)

		if err != nil {
			return SciNotationNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntMath{}.GetMagnitude(bINumIntPart.bigInt) "+
					"Error='%v'", err.Error())
		}

		if magnitudeBigInt.Cmp(bigIntMaxUint32) == 1 {
			return SciNotationNum{}.New(),
				errors.New(ePrefix + "Integer Magnitude greater than Max Uint32")
		}

		uintMagnitude := uint(magnitudeBigInt.Uint64())

		newBINum := new(BigIntNum).NewBigInt(bNum.bigInt, bNum.precision+uintMagnitude)

		_ = sciNotationNum.SetBigIntNumElements(
			newBINum,
			new(BigIntNum).NewBigInt(magnitudeBigInt, 0))

	} else {
		// Must be bINumFracPart > 0
		magnitudeBigInt, err := BigIntMath{}.GetMagnitude(bNum.bigInt)

		if err != nil {
			return SciNotationNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntMath{}.GetMagnitude(bNum.bigInt). "+
					"Error='%v'", err.Error())
		}

		if magnitudeBigInt.Cmp(bigIntMaxUint32) == 1 {
			return SciNotationNum{}.New(),
				errors.New(ePrefix + "Fractional Magnitude greater than Max Uint32")
		}

		uintMagnitude := uint(magnitudeBigInt.Uint64())

		bINumFracPart := new(BigIntNum).NewBigInt(bNum.bigInt, uintMagnitude)
		precisionFrac := int64(uintMagnitude) - int64(bNum.precision)

		bINumScale := new(BigIntNum).NewInt64Exponent(precisionFrac, 0)

		_ = sciNotationNum.SetBigIntNumElements(bINumFracPart, bINumScale)
	}

	sciNotationNum.SetMantissaLength(mantissaLen)

	return sciNotationNum, nil
}

//	GetSciNotationStr
//
//	Returns a string expressing the current BigIntNum
//	numerical value as scientific notation.
//
//	The returned Scientific Notation string is formatted
//	in E-Notation.
//
// ----------------------------------------------------------------
//
// # Reference:
//
//	https://en.wikipedia.org/wiki/Scientific_notation#E_notation
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	mantissaLen			uint
//
//		Specifies the length of the mantissa in the returned
//		scientific notation string. If the value of
//		'mantissaLen' is less than two ('2'), this method
//		will automatically set the 'mantissaLen' to a default
//		value of two ('2').
//
// ----------------------------------------------------------------
//
//	Example Scientific Notation:
//
//		scientific notation string: '2.652e+8'
//		significand = '2.652'
//		significand integer digit = '2'
//		mantissa		= significand factional digits = '.652'
//		exponent    = '8'  (10^8)
func (bNum *BigIntNum) GetSciNotationStr(mantissaLen uint) (string, error) {

	ePrefix := "BigIntNum.GetSciNotationStr() "

	sciNotation, err := bNum.GetSciNotationNumber(mantissaLen)

	if err != nil {
		return "",
			fmt.Errorf(ePrefix+
				"Error returned by bNum.GetSciNotationNumber(mantissaLen). "+
				"Error='%v'", err.Error())
	}

	result, err := sciNotation.GetSciNotationStr(mantissaLen)

	if err != nil {
		return "",
			fmt.Errorf(ePrefix+
				"Error returned by sciNotation.GetSciNotationStr(mantissaLen). "+
				"Error='%v'", err.Error())
	}

	return result, nil
}

// GetThisPointer
//
// Returns a pointer to the current instance of this
// BigIntNum.
func (bNum *BigIntNum) GetThisPointer() *BigIntNum {
	return bNum
}

// GetThousandsSeparator
//
// Returns a rune which represents the character
// currently used to separate thousands in the
// display of the current BigIntNum instance.
//
// In the United States (US), the thousands separator
// is a comma character.
//
// Example: 1,000,000,000
func (bNum *BigIntNum) GetThousandsSeparator() rune {

	if bNum.thousandsSeparator == 0 {
		bNum.thousandsSeparator = ','
	}

	return bNum.thousandsSeparator
}

// GetUInt - Returns a type 'uint' containing the 32-bit unsigned
// integer value of the current BigIntNum instance.
//
// If the current BigIntNum value is greater than the maximum
// 'uint' value, the maximum 32-bit unsigned integer value is returned
// in addition to an 'error'.
//
// If the current BigIntNum value is less than the minimum 'uint'
// value, the minimum 32-bit integer value of zero is returned along
// with an 'error'.
func (bNum *BigIntNum) GetUInt() (uint, error) {
	ePrefix := "BigIntNum.GetUInt() "

	bIMaxUint := big.NewInt(int64(math.MaxUint32))

	if bNum.bigInt.Cmp(big.NewInt(0)) == -1 {
		return uint(0),
			fmt.Errorf(ePrefix + "Error: BigIntNum is LESS THAN minimum 'uint' value of zero.")
	}

	if bNum.bigInt.Cmp(bIMaxUint) == 1 {
		return math.MaxUint32,
			fmt.Errorf("Error: BigIntNum is GREATER THAN maximum 'uint' value.")
	}

	return uint(bNum.bigInt.Uint64()), nil
}

// GetUInt64
//
// Returns the integer value of BigIntNum.bigInt as a
// 64-bit unsigned integer. If the value of
// BigIntNum.bigInt exceeds that of the maximum unsigned
// 64-bit integer value, an error will be returned.
func (bNum *BigIntNum) GetUInt64() (uint64, error) {

	bIntMaxUint64 := big.NewInt(0).SetUint64(uint64(math.MaxUint64))

	if bNum.bigInt.Cmp(bIntMaxUint64) == 1 {
		return uint64(0),
			fmt.Errorf("BigIntNum.GetUInt64() - Error: The value of this BigIntNum instance "+
				"exceeds the maximum value of the unsigned 64-bit integer. BigIntNum='%v' MaxUint64='%v' ",
				bNum.bigInt.Text(10), bIntMaxUint64.Text(10))
	}

	return bNum.bigInt.Uint64(), nil
}

// Inverse
//
// Returns the inverse of the current BigIntNum value.
//
// The inverse of the value is equal to one ('1') divided
// by the numeric value of the current BigIntNum.
//
// The BigIntNum return value for this operation will
// contain numeric separators (decimal separator, integer
// separator, integer grouping and currency symbol)
// copied from the original BigIntNum instance.
func (bNum *BigIntNum) Inverse(maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntNum.Inverse() "

	if bNum.IsZero() {
		return new(BigIntNum).NewZero(0), nil
	}

	bIOne := new(BigIntNum).NewOne(0)

	err := bIOne.SetNumericSeparatorsDto(bNum.GetNumericSeparatorsDto())

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by bIOne.SetNumericSeparatorsDto(bNum.GetNumericSeparatorsDto()) "+
				"Error='%v' \n", err.Error())
	}

	var inverse BigIntNum

	inverse, err = BigIntMathDivide{}.BigIntNumFracQuotient(bIOne, bNum.CopyOut(), maxPrecision)

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient(...) "+
				"Error='%v' \n", err.Error())
	}

	return inverse, nil
}

// IsEvenNumber
//
// Returns true if the current BigIntNum value is evenly
// divisible by 2.
//
// Even Number Definitions:
//
// https://www.mathsisfun.com/definitions/even-number.html
//
// In mathematics, parity is the property of an
// integer's inclusion in one of two categories:
// even or odd. An integer is even if it is evenly
// divisible by two and odd if it is not even.
//
// "Examples of even numbers include −4, 0, 82 and
// 178. In particular, zero is an even number."
//
// https://en.wikipedia.org/wiki/Parity_(mathematics)
func (bNum *BigIntNum) IsEvenNumber() (bool, error) {

	if bNum.precision > 0 {
		return false, nil
	}

	if bNum.IsZero() {
		return true, nil
	}

	_, mod, err := BigIntMathDivide{}.BigIntNumDivideByTwoQuoMod(bNum.CopyOut(), 50)

	if err != nil {
		ePrefix := "BigIntNum.IsEvenNumber() "
		return false,
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.BigIntNumQuotientMod(nthRoot, bigINumTwo, 0) "+
				"Error='%v' ", err.Error())
	}

	return mod.IsZero(), nil
}

// Increment
//
// Adds a value of +1 (plus one) to the numeric value of
// the current BigIntNum instance.
func (bNum *BigIntNum) Increment() {

	biNumOne := new(BigIntNum).NewOne(bNum.precision)

	bPair := BigIntPair{}.NewBigIntNum(
		bNum.CopyOut(),
		biNumOne)

	result := BigIntMathAdd{}.AddPair(bPair)

	bNum.CopyIn(result)

}

// IsValid - returns a boolean value signaling whether the
// current BigIntNum object is valid. For types of errors
// corrective action is performed to restore the BigIntNum
// instance.
func (bNum *BigIntNum) IsValid(errName string) error {

	if errName == "" {
		errName = "BigIntNum.IsValid() "
	}

	errName += "BigIntNum INVALID! "

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.sign != -1 && bNum.sign != 1 {
		bNum.Reset()
		return nil
	}

	if bNum.absBigInt == nil ||
		bNum.scaleFactor == nil {
		bNum.Reset()
		return nil
	}

	return nil
}

// IsZero - Returns a boolean signaling whether the current
// BigIntNum value is zero.
func (bNum *BigIntNum) IsZero() bool {

	if bNum.bigInt.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	return false
}

// Mod
//
// Performs a modulo operation where the current
// BigIntNum numeric value is the dividend and the
// divisor is the input parameter, 'divisor'.
//
// The modulo operation finds the remainder after
// division of one number by another (sometimes
// called modulus). Wikipedia
// https://en.wikipedia.org/wiki/Modulo_operation
//
//	dividend = bNum
//	dividend % divisor = modulo
//
// The result of this modulo operation is returned as
// a BigIntNum, 'modulo'. 'modulo' may consist of an
// integer or a floating point value consisting of
// integer and fractional digits.
//
// Input parameter 'maxPrecision' is used to control
// the maximum precision of the resulting floating point
// 'modulo'. Precision is defined as the number of
// fractional digits to the right of the decimal point.
// Be advised that these calculations can support very
// large precision values.
//
// The returned BigIntNum instance, 'modulo', will
// contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied
// from the current BigIntNum instance (bNum).
func (bNum *BigIntNum) Mod(
	divisor BigIntNum,
	maxPrecision uint) (modulo BigIntNum, err error) {

	return BigIntMathDivide{}.BigIntNumModulo(bNum.CopyOut(), divisor, maxPrecision)
}

// Multiply - Multiplies the numerical value of the current BigIntNum instance
// ('multiplier') times input parameter 'multiplicand'. The 'product' of this
// multiplication operation is returned as a BigIntNum.
//
//	multiplier = bNum
//	multiplier X multiplicand = product
//
// The BigIntNum instance returned by this method, 'product', will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from the current BigIntNum instance.
func (bNum *BigIntNum) Multiply(multiplicand BigIntNum) (product BigIntNum) {

	return BigIntMathMultiply{}.MultiplyBigIntNums(bNum.CopyOut(), multiplicand)
}

// MultiplyByFive - Multiplies the numerical value of the current BigIntNum
// instance times five (5). The product is returned as a BigIntNum.
//
//	product = bNum X 5
//
//
//	The BigIntNum instance returned by this method will
//	contain numeric separators (decimal separator, integer
//	separator, integer grouping and currency symbol)
//	copied from the original BigIntNum instance.
func (bNum *BigIntNum) MultiplyByFive() BigIntNum {

	return BigIntMathMultiply{}.MultiplyBigIntNumByFive(bNum.CopyOut())
}

// MultiplyByTen
//
// Multiplies the numerical value of the current BigIntNum
// instance times ten (10). The product is returned as a
// BigIntNum.
//
// product = bNum X 10
//
// The BigIntNum instance returned by this method will
// contain numeric separators (decimal separator, integer
// separator, integer grouping and currency symbol)
// copied from the original BigIntNum instance.
func (bNum *BigIntNum) MultiplyByTen() BigIntNum {

	return BigIntMathMultiply{}.MultiplyBigIntNumByTen(bNum.CopyOut())
}

// MultiplyByTenToPower
//
// Multiplies the numerical value of the current
// BigIntNum instance times ten to the power of
// 'exponent' (10^exponent). The product is returned
// as the new value for the current BigIntNum. The
// original value of the BigIntNum instance will be
// overwritten and destroyed.
//
//	bNum = bNum X 10^exponent
func (bNum *BigIntNum) MultiplyByTenToPower(exponent uint) {

	if bNum.precision >= exponent {

		bNum.CopyIn(
			new(BigIntNum).NewBigInt(bNum.bigInt, bNum.precision-exponent))

	} else {
		// exponent > bNum.precision
		scaleVal :=
			big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(exponent-bNum.precision)), nil)

		newVal := big.NewInt(0).Mul(bNum.bigInt, scaleVal)

		bNum.CopyIn(
			new(BigIntNum).NewBigInt(newVal, 0))
	}

	return
}

// MultiplyByTenToPowerAdd
//
// Performs three operations on the current BigIntNum
// instance.
//
// (1)	First, the method multiplies the numerical value
//
//	of the current BigIntNum instance times ten to the
//	power of 'exponent' (10^exponent).
//
//		bNum1 = bNum X 10^exponent
//
// (2)	Second, the method adds input parameter 'addend' to
//
//	the product generated by operation (1), above.
//
//		bNum2 = bNum1 + 'addend'
//
// (3)	Third and finally, the original value of the current
//
//	BigIntNum instance will be overwritten and replaced
//	by the 'bNum2' value generated in operation (2),
//	above.
func (bNum *BigIntNum) MultiplyByTenToPowerAdd(
	exponent uint,
	addend BigIntNum) {

	var bx BigIntNum

	if bNum.precision >= exponent {

		bx = new(BigIntNum).NewBigInt(bNum.bigInt, bNum.precision-exponent)

	} else {
		// exponent > bNum.precision
		scaleVal :=
			big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(exponent-bNum.precision)), nil)

		newVal := big.NewInt(0).Mul(bNum.bigInt, scaleVal)

		bx = new(BigIntNum).NewBigInt(newVal, 0)
	}

	bNum.CopyIn(BigIntMathAdd{}.AddBigIntNums(bx, addend))

	return
}

// MultiplyByThree
//
// Multiplies the numerical value of the current
// BigIntNum instance times three (3). The product is
// returned as a new instance of BigIntNum.
//
//	product = bNum X 3
//
// The BigIntNum instance returned by this method will
// contain numeric separators (decimal separator, integer
// separator, integer grouping and currency symbol)
// copied from the original BigIntNum instance.
func (bNum *BigIntNum) MultiplyByThree() BigIntNum {

	return BigIntMathMultiply{}.MultiplyBigIntNumByThree(bNum.CopyOut())
}

// MultiplyByTwo
//
// Multiplies the numerical value of the current
// BigIntNum instance times two (2). The product
// is returned as a new instance of BigIntNum.
//
//	product = bNum X 2
//
// The BigIntNum instance returned by this method will
// contain numeric separators (decimal separator, integer
// separator, integer grouping and currency symbol)
// copied from the original BigIntNum instance.
func (bNum *BigIntNum) MultiplyByTwo() BigIntNum {

	return BigIntMathMultiply{}.MultiplyBigIntNumByTwo(bNum.CopyOut())
}

// New - returns a new BigIntNum instance initialized to zero.
//
// The BigIntNum instance returned by this method will contain USA
// default numeric separators (decimal separator, thousands separator
// and currency symbol).
func (bNum *BigIntNum) New() BigIntNum {
	b := BigIntNum{}
	b.Empty()
	return b
}

// NewWithNumSeps - returns a new BigIntNum instance initialized to zero.
//
// The BigIntNum instance returned by this method will contain USA
// default numeric separators (decimal separator, thousands separator
// and currency symbol).
func (bNum *BigIntNum) NewWithNumSeps(numSeps NumericSeparatorDto) BigIntNum {

	numSeps.SetDefaultsIfEmpty()

	b := BigIntNum{}
	b.Empty()
	_ = b.SetNumericSeparatorsDto(numSeps)

	return b
}

//	NewBigInt
//
//	Creates a new BigIntNum instance using a *big.Int
//	type and its associated precision.
//
//	The 'precision' parameter specifies the number of
//	digits to the right of the decimal point. The
//	Numeric value is equal to bigI x 10^(precision x -1).
//
//	This effectively locates the decimal place by counting
//	from the extreme right of the integer number,
//	'precision' places to the left. See the example below.
//
//	The new BigIntNum instance returned by this method will
//	contain USA default numeric separators (decimal separator,
//	thousands separator and currency symbol).
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	bigI				*big.Int
//
//		'bigI' is a type *big.Int and represents the integer
//		value of the number; that is, the numeric value
//		without decimal digits.
//
//	precision			int
//
//		This unsigned integer (always a positive value)
//		identifies the location of the decimal place in the
//		integer value 'bigI'.
//
//		The decimal place location is calculated by starting
//		with the right most digit in the integer number and
//		counting left, 'precision' places.
//
//		Example:
//
//			Integer 					Numeric
//			Value		precision		Value
//			-------		---------		-------
//			 123456			3			123.456
func (bNum *BigIntNum) NewBigInt(bigI *big.Int, precision uint) BigIntNum {

	if bigI == nil {
		bigI = big.NewInt(0)
	}

	b := BigIntNum{}
	b.Empty()
	b.SetBigInt(bigI, precision)
	return b
}

//	NewBigIntPrecision
//
//	Creates a new BigIntNum instance using a *big.Int
//	type and its associated precision.
//
//	The 'precision' parameter specifies the number of
//	fractional digits to the right of the decimal place.
//	The Numeric value is equal to:
//		bigI x 10^(precision x -1)
//
//	This effectively locates the decimal place by counting
//	from the extreme right of the integer number,
//	'precision' places to the left. See the example below.
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	bigI				*big.Int
//
//		'bigI' is a type *big.Int and represents the
//		integer value of the number; that is, the
//		numeric value without decimal digits.
//
//	precision			*big.Int
//
//		This integer value (always a positive value)
//		identifies the location of the decimal place
//		in the integer value 'bigI'.
//
//		The decimal place location is calculated by
//		starting with the right most digit in the
//		integer number and counting	left, 'precision'
//		places. If precision is greater than the maximum
//		value of an unsigned integer (+4,294,967,295 which
//		equals 2^32 − 1), an error will be triggered. Also,
//		if the 'precision' value is less than zero, an
//		error will be triggered.
//
//		Example:
//
//		Integer Value	precision	Numeric Value
//			123456			3			123.456
//
// ----------------------------------------------------------------
//
// # Return Values
//
//	BigIntNum
//
//		The new BigIntNum instance returned by this
//		method will contain USA default numeric
//		separators (decimal separator, thousands
//		separator and currency symbol).
//
//	error
//
//		If this method completes successfully, the
//		returned error Type is set equal to 'nil'.
//
//		If errors are encountered during processing, the
//		returned error Type will encapsulate an error
//		message.
func (bNum *BigIntNum) NewBigIntPrecision(bigInt, precision *big.Int) (BigIntNum, error) {

	ePrefix := "BigIntNum.NewBigIntPrecision() "

	if bigInt == nil {
		bigInt = big.NewInt(0)
	}

	if precision.Cmp(big.NewInt(0)) == -1 {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' IS LESS THAN ZERO! "+
				"precision='%v' ", precision.Text(10))
	}

	maxUint32 := big.NewInt(0).SetUint64(uint64(math.MaxUint32))

	if precision.Cmp(maxUint32) == 1 {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'precision' exceeds maximum limit of '4,294,967,295'! "+
				"precision='%v' ", precision.Text(10))
	}

	b := BigIntNum{}
	b.Empty()
	b.SetBigInt(bigInt, uint(precision.Uint64()))

	return b, nil
}

// NewBigIntExponent
//
// New bigInt Exponent returns a new BigIntNum instance
// in which the numeric value is set using an integer
// multiplied by 10 raised to the power of the
// 'exponent' parameter.
//
// Example:
//
//	numeric value = integer X 10^exponent
//
// If exponent is less than +1, precision is set equal
// to exponent and bigI is unchanged.
//
// If exponent is greater than 0, bigI is multiplied by
// 10 raised to the power of 'exponent' and precision
// is set equal to zero.
//
// Examples:
//
//	biNum :=
//		BigIntNum{}.NewBigIntExponent(big.NewInt(int64(123456)), -3) = "123.456"  precision = 3
//
//	biNum :=
//		BigIntNum{}.NewBigIntExponent(big.NewInt(int64(123456)), 3) = "123456.000" precision = 3
func (bNum *BigIntNum) NewBigIntExponent(bigI *big.Int, exponent int) BigIntNum {

	b := BigIntNum{}
	b.Empty()
	b.SetBigIntExponent(bigI, exponent)
	return b
}

//	NewBigFloat
//
//	Returns a new BigIntNum instance using a *big.Float
//	floating point number as an input parameter.  The
//	precision of the number is specified by the input
//	parameter, 'decimalPlaces'.
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	bigFloat			*big.Float
//
//		This *big.Float value will be converted into an
//		instance of BigIntNum and returned to the calling
//		function.
//
//	maxPrecision		uint
//
//		The maximum precision for the resulting BigIntNum
//		after conversion of input parameter 'f64'. Resulting
//		precision will never be greater than 'maxPrecision'
//		for type big.Float; however, actual precision may be
//		less than 'maxPrecision'.
func (bNum *BigIntNum) NewBigFloat(
	bigFloat *big.Float,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntNumNewFloat64() "

	b := new(BigIntNum).NewZero(0)

	err := b.SetBigFloat(bigFloat, maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by b.SetBigFloat(bigFloat, maxPrecision). "+
				"Error='%v' ", err.Error())
	}

	return b, nil
}

// NewBigIntFixedDecimal
//
// Creates and returns a new BigIntNum instance based
// on input parameter 'fd' of type BigIntFixedDecimal.
func (bNum *BigIntNum) NewBigIntFixedDecimal(fd BigIntFixedDecimal) BigIntNum {

	fd.IsValid()

	return new(BigIntNum).NewBigInt(fd.GetInteger(), fd.GetPrecision())

}

// NewFromIntFracStrings
//
// Creates a new BigIntNum instance based on a numeric
// value represented by separate integer and fractional
// components.
//
// Input parameters 'intStr' and 'fracStr' are strings
// representing the integer and fractional elements of
// the numeric value. These elements are combined by this
// method to create a numeric value which is then
// assigned to the new BigIntNum instance.
//
// Input parameter 'signVal' must be set to one of two
// values: +1 or -1. This value is used to signal the
// sign of the resulting numeric value. +1 generates
// a positive number and -1 generates a negative number.
func (bNum *BigIntNum) NewFromIntFracStrings(
	intStr, fracStr string, signVal int) (BigIntNum, error) {

	b2 := new(BigIntNum).NewZero(0)

	err := b2.SetIntFracStrings(intStr, fracStr, signVal)

	if err != nil {
		ePrefix := "BigIntNum.NewFromIntFracStrings() "
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by b2.SetIntFracStrings(intStr, fracStr, signVal). "+
				"Error='%v' \n", err.Error())
	}

	return b2, nil
}

//	NewFloat32
//
//	Returns a new BigIntNum instance using a float32
//	floating point input parameter.  The precision of
//	the number is specified by the input parameter,
//	'decimalPlaces'.
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	f32					float32
//
//		This float32 value will be converted into an
//		instance of BigIntNum.
//
//	maxPrecision		uint
//
//		The maximum precision for the result BigIntNum
//		after conversion of input parameter f64.
//
//		Precision will never be greater than 'maxPrecision';
//		however, actual precision may be less than
//		'maxPrecision'.
func (bNum *BigIntNum) NewFloat32(f32 float32, maxPrecision uint) (BigIntNum, error) {
	ePrefix := "BigIntNumNewFloat32() "

	b := new(BigIntNum).NewZero(0)

	err := b.SetFloat32(f32, maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by b.SetFloat32(f32, maxPrecision). "+
				"Error='%v' ", err.Error())
	}

	return b, nil
}

//	NewFloat64
//
//	Returns a new BigIntNum instance using a float64
//	floating point number as an input parameter.
//
//	The precision of the number is specified by the
//	input parameter, 'decimalPlaces'.
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	f64					float64
//
//		This float64 value will be converted into an
//		instance of BigIntNum.
//
//	maxPrecision		uint
//
//		The maximum precision for the result BigIntNum
//		after conversion of input parameter f64. Precision
//		will never be greater than 'maxPrecision'; however,
//		actual precision may be less than 'maxPrecision'.
func (bNum *BigIntNum) NewFloat64(f64 float64, maxPrecision uint) (BigIntNum, error) {
	ePrefix := "BigIntNumNewFloat64() "

	b := BigIntNum{}
	b.Empty()
	err := b.SetFloat64(f64, maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by b.SetFloat64(f64, decimalPlaces). "+
				"Error='%v' ", err.Error())
	}

	return b, nil
}

//	NewInt
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'intNum' which is passed as
//	type 'int'.
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint'
//
// ----------------------------------------------------------------
//
// # Usage:
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		intNum := int(123456)
//		precision := uint(3)
//		bINum := new(BigIntNum).NewInt(intNum, precision)
//		bINum is now equal to 123.456
//
//	Examples:
//
//									BigIntNum
//		intNum		precision		Result
//		------		---------		---------
//		123456			4			12.3456
//		123456			0			123456
//		123456			1			12345.6
func (bNum *BigIntNum) NewInt(intNum int, precision uint) BigIntNum {

	b2 := new(BigIntNum).NewBigInt(big.NewInt(int64(intNum)), precision)

	_ = b2.SetNumericSeparatorsDto(bNum.GetNumericSeparatorsDto())

	return b2
}

//	NewIntExponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//		numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	intNum				int
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		biNum := new(BigIntNum).NewIntExponent(123456, -3)
//	 	biNum is now equal to "123.456", precision = 3
//
//		biNum := new(BigIntNum).NewIntExponent(123456, 3)
//	 	biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//
//							BigIntNum
//		intNum	exponent	  Result
//		------	--------	---------
//		123456	  -3		123.456
//		123456	   3		123456.000
//		123456	   0		123456
func (bNum *BigIntNum) NewIntExponent(intNum int, exponent int) BigIntNum {

	bigI := big.NewInt(int64(intNum))

	b := new(BigIntNum).NewZero(0)

	b.SetBigIntExponent(bigI, exponent)

	return b
}

//	NewInt32
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'int32Num' which is passed
//	as type 'int32'.
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint'
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction
//	with the new(BigIntNum) syntax thereby allowing
//	BigIntNum type creation and initialization in one
//	step.
//
//		num :=		 int32(123456)
//		precision := uint(3)
//
//		bINum := new(BigIntNum).NewInt32(num, precision)
//		bINum is now equal to 123.456
//
//	Examples:
//									BigIntNum
//		int32Num	precision		  Result
//		--------	---------		---------
//		123456			4			12.3456
//		123456			0			123456
//		123456			1			12345.6
func (bNum *BigIntNum) NewInt32(int32Num int32, precision uint) BigIntNum {

	return new(BigIntNum).NewBigInt(big.NewInt(int64(int32Num)), precision)
}

//	NewInt32Exponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//		numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	int32Num			int32
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction
//	with the new(BigIntNum{}) syntax thereby allowing
//	BigIntNum type creation and initialization in one
//	step.
//
//		biNum := new(BigIntNum).NewInt32Exponent(123456, -3)
//		biNum is now equal to "123.456", precision = 3
//
//		biNum := new(BigIntNum).NewInt32Exponent(123456, 3)
//		biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//								BigIntNum
//		int32Num	exponent	Result
//		--------	--------	---------
//		123456			-3		123.456
//		123456			 3		123456.000
//		123456			 0		123456
func (bNum *BigIntNum) NewInt32Exponent(int32Num int32, exponent int) BigIntNum {

	bigI := big.NewInt(int64(int32Num))

	b := new(BigIntNum).NewZero(0)

	b.SetBigIntExponent(bigI, exponent)

	return b
}

//	NewInt64
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'int64Num' which is passed
//	as type 'int64'.
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint'
//
// ----------------------------------------------------------------
//
//	# Usage
//
//	This method is designed to be used in conjunction
//	with the new(BigIntNum) syntax thereby allowing
//	BigIntNum type creation and initialization in one
//	step.
//
//		int64Num := int64(123456)
//		precision := uint(3)
//		bINum := new(BigIntNum).NewInt64(int64Num, precision)
//		bINum is now equal to 123.456
//
//	Examples
//
//								BigIntNum
//		int64Num	precision	  Result
//		--------	---------	---------
//		123456			4		12.3456
//		123456			0		123456
//		123456			1		12345.6
func (bNum *BigIntNum) NewInt64(int64Num int64, precision uint) BigIntNum {

	return new(BigIntNum).NewBigInt(big.NewInt(int64Num), precision)
}

//	NewInt64Exponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//			numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	int64Num			int64
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		biNum := new(BigIntNum).NewInt64Exponent(123456, -3)
//		biNum is now equal to "123.456", precision = 3
//
//		biNum := new(BigIntNum).NewInt64Exponent(123456, 3)
//		biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//
//								BigIntNum
//		int64Num	exponent	  Result
//		--------	--------	---------
//		123456		-3			123.456
//		123456		 3			123456.000
//		123456		 0			123456
func (bNum *BigIntNum) NewInt64Exponent(int64Num int64, exponent int) BigIntNum {

	bigI := big.NewInt(int64Num)

	b := new(BigIntNum).NewZero(0)

	b.SetBigIntExponent(bigI, exponent)

	return b
}

// NewIntFracStr
//
// Creates a new BigIntNum instance based on a numeric
// value represented by separate strings containing
// integer and fractional digits.
//
// Input parameters 'intStr' and 'fracStr' are strings
// representing the integer and fractional components of
// a numeric value. They are combined by this method to
// create a numeric value which is assigned to the
// BigIntNum instance returned by this method.
//
// Input parameter 'signVal' must be set to one of two
// values: +1 or -1. This value is used to signal the
// sign of the resulting numeric value. +1 generates a
// positive number and -1 generates a negative number.
// If input parameters 'inStr' or 'fracStr' contain a
// leading minus or plus sign character, it will be
// ignored.
//
// The sign of the resulting BigIntNum numeric value is
// controlled exclusively by input parameter, 'signVal'.
func (bNum *BigIntNum) NewIntFracStr(intStr, fracStr string, signVal int) (BigIntNum, error) {

	bIntNum := new(BigIntNum).NewZero(0)

	err := bIntNum.SetIntFracStrings(intStr, fracStr, signVal)

	if err != nil {

		ePrefix := "BigIntNum.NewIntFracStr() "

		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by bIntNum.SetIntFracStrings(intStr, fracStr, signVal) "+
				"Error='%v' \n", err.Error())

	}

	return bIntNum, nil
}

// NewNumStr
//
// Receives a number string as input and returns a new
// BigIntNum instance.
//
// This method assumes that the input parameter 'numStr'
// is a string of numeric digits which may be delimited
// by default United States (US) numeric separators.
//
// Default US numeric separators are defined as:
//
//	decimal separator	= '.'
//	integer separator	= ','
//	integer grouping	= thousands
//	currency symbol		= '$'
//
// If the subject 'numStr' employs other national or
// cultural numeric separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewNumStr(numStr string) (BigIntNum, error) {

	ePrefix := "BigIntNum.NewNumStr() "

	b := new(BigIntNum).NewZero(0)
	err := b.SetNumStr(numStr)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by b.SetNumStr(numStr). "+
				"numStr='%v' Error='%v'",
				numStr, err.Error())
	}

	return b, nil
}

// NewNumStrWithNumSeps
//
// Receives a number string as input and returns a new
// BigIntNum instance. The input parameter 'numSeps'
// contains numeric separators (decimal separator,
// integer separator, integer grouping and currency
// symbol) which will be used to parse the number
// string.
//
// In addition, the numeric separators contained in
// input parameter 'numSeps' will be copied to the
// returned BigIntNum instance.
func (bNum *BigIntNum) NewNumStrWithNumSeps(
	numStr string,
	numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntNum.NewNumStrWithNumSeps() "

	numSeps.SetDefaultsIfEmpty()

	b2 := new(BigIntNum).New()

	err := b2.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by b2.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' \n", err.Error())

	}

	err = b2.SetNumStr(numStr)

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by b2.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' \n", err.Error())
	}

	return b2, nil
}

// NewNumStrMaxPrecision
//
// Receives a number string as input and returns a new
// BigIntNum instance. If the resulting precision
// exceeds input parameter 'maxPrecision', the
// returned BigIntNum result will be rounded to
// 'maxPrecision' decimal places.
func (bNum *BigIntNum) NewNumStrMaxPrecision(
	numStr string,
	maxPrecision uint) (BigIntNum, error) {

	b := new(BigIntNum).NewZero(0)

	err := b.SetNumStr(numStr)

	if err != nil {

		ePrefix := "BigIntNum.NewNumStr() "

		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by b.SetNumStr(numStr). "+
				"numStr='%v' Error='%v' ",
				numStr, err.Error())
	}

	if b.precision > maxPrecision {
		b.RoundToDecPlace(maxPrecision)
	}

	return b, nil
}

// NewNumStrDto
//
// Receives a NumStrDto instance as input and returns
// a new BigIntNum instance.
func (bNum *BigIntNum) NewNumStrDto(nDto NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntNum.NewNumStrDto() "

	err := nDto.IsValid(ePrefix + "'nDto' INVALID! ")

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned from nDto.IsValid(\"\"). "+
				"NumStr='%v' Error='%v'", nDto.GetNumStr(), err.Error())
	}

	var bigI *big.Int

	bigI, err = nDto.GetBigInt()

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by nDto.GetBigInt(). "+
				"Error='%v'", err.Error())
	}

	b := new(BigIntNum).NewZero(0)

	b.SetBigInt(bigI, uint(nDto.GetPrecision()))

	return b, nil
}

// NewOne
//
// Returns a BigIntNum Type with a value equal to '1'
// (one).
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '1', set
// 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		 1
//		1		 1.0
//		2		 1.00
//		3		 1.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewOne(precision uint) BigIntNum {

	b := new(BigIntNum).NewZero(0)

	if precision == 0 {
		b.SetBigInt(big.NewInt(1), 0)
		return b
	}

	scaleVal :=
		big.NewInt(0).Exp(big.NewInt(10),
			big.NewInt(int64(precision)),
			nil)

	newVal := big.NewInt(0).Mul(big.NewInt(1), scaleVal)

	b.SetBigInt(newVal, precision)

	return b
}

// NewTwo
//
// Returns a BigIntNum Type with a value equal to  '2'
// (two).
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '2', set
// 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		 2
//		1		 2.0
//		2		 2.00
//		3		 2.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewTwo(precision uint) BigIntNum {

	b := new(BigIntNum).NewZero(0)

	if precision == 0 {
		b.SetBigInt(big.NewInt(2), 0)
		return b
	}

	scaleVal :=
		big.NewInt(0).Exp(big.NewInt(10),
			big.NewInt(int64(precision)),
			nil)

	newVal := big.NewInt(0).Mul(big.NewInt(2), scaleVal)

	b.SetBigInt(newVal, precision)

	return b
}

// NewThree
//
// Returns a BigIntNum Type with a value equal to  '3'
// (three).
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '3', set
// 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		  3
//		2		  3.00
//		3		  3.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewThree(precision uint) BigIntNum {

	b := new(BigIntNum).NewZero(0)

	if precision == 0 {
		b.SetBigInt(big.NewInt(3), 0)
		return b
	}

	scaleVal := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
	newVal := big.NewInt(0).Mul(big.NewInt(3), scaleVal)
	b.SetBigInt(newVal, precision)

	return b
}

// NewFive
//
// Returns a BigIntNum with integer value of  '5' (five).
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '5', set
// 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		  5
//		1		  5.0
//		2		 5.00
//		3		 5.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewFive(precision uint) BigIntNum {

	b := new(BigIntNum).NewZero(0)

	if precision == 0 {
		b.SetBigInt(big.NewInt(5), 0)
		return b
	}

	scaleVal :=
		big.NewInt(0).Exp(big.NewInt(10),
			big.NewInt(int64(precision)),
			nil)

	newVal :=
		big.NewInt(0).Mul(
			big.NewInt(5),
			scaleVal)

	b.SetBigInt(newVal, precision)

	return b
}

// NewTen
//
// Returns a BigIntNum with integer value of  '10' (ten).
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '10', set
// 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		  10
//		1		  10.0
//		2		  10.00
//		3		  10.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewTen(precision uint) BigIntNum {

	b := new(BigIntNum).NewZero(0)

	if precision == 0 {
		b.SetBigInt(big.NewInt(10), 0)
		return b
	}

	scaleVal :=
		big.NewInt(0).Exp(big.NewInt(10),
			big.NewInt(int64(precision)),
			nil)

	newVal :=
		big.NewInt(0).Mul(big.NewInt(10),
			scaleVal)

	b.SetBigInt(newVal, precision)

	return b
}

// NewZero
//
// Returns a BigIntNum instance with a value equal to
// zero.
//
// The number of zeros created after the decimal
// placeholder (fractional digits) is determined by the
// input parameter 'precision'.
//
// To create an integer with a value equal to '0',
// set 'precision' equal to zero (0).
//
// # Examples
//
// 'precision' values and resulting outcomes.
//
//	precision
//	  value		Result
//	---------	------
//		0		  0
//		2		  0.00
//		3		  0.000
//
// The new BigIntNum instance returned by this method
// will contain United States (US) default numeric
// separators (decimal separator, integer separator,
// integer grouping and currency symbol).
//
// To implement other national or cultural numeric
// separators, see method:
//
//	BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewZero(precision uint) BigIntNum {

	b := BigIntNum{}
	b.Empty()
	b.SetBigInt(big.NewInt(0), precision)

	return b

}

//	NewUint
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'uintNum' which is passed
//	as type 'uint' (Unsigned Integer).
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint'
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		uintNum := uint(123456)
//		precision := uint(3)
//		bINum := new(BigIntNum).NewUint(uintNum, precision)
//		bINum is now equal to 123.456
//
//	Examples:
//
//								BigIntNum
//		uintNum		exponent	  Result
//		--------	--------	---------
//		123456		  -4		12.3456
//		123456		   4		123456.0000
//		123456		   0		123456
//		123456		  -1		12345.6
//
//	The new BigIntNum instance returned by this method
//	will contain United States (US) default numeric
//	separators (decimal separator, integer separator,
//	integer grouping and currency symbol).
//
//	To implement other national or cultural numeric
//	separators, see method:
//		BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewUint(uintNum uint, precision uint) BigIntNum {

	return new(BigIntNum).NewBigInt(big.NewInt(0).SetUint64(uint64(uintNum)), precision)
}

//	NewUintExponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//		numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
//	# Input Parameters
//
//	uintNum				uint
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		biNum := BigIntNum{}.NewUintExponent(123456, -3)
//		biNum is now equal to "123.456", precision = 3
//
//		biNum := BigIntNum{}.NewUintExponent(123456, 3)
//		biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//
//								BigIntNum
//		int64Num	exponent	  Result
//		--------	--------	---------
//		123456		  -3		123.456
//		123456		   3		123456.000
//		123456		   0		123456
func (bNum *BigIntNum) NewUintExponent(uintNum uint, exponent int) BigIntNum {

	baseBInt := big.NewInt(int64(uintNum))

	b2 := new(BigIntNum).New()

	b2.SetBigIntExponent(baseBInt, exponent)

	return b2
}

//	NewUint32
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'uint32Num' which is passed
//	as type 'uint32' (Unsigned 32-bit Integer).
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint32'.
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		uint32Num := uint32(123456)
//		precision := uint(3)
//		bINum := BigIntNum{}.NewUint32(uint32Num, precision)
//		bINum is now equal to 123.456
//
//	Examples:
//
//								BigIntNum
//		uintNum		exponent	  Result
//		--------	--------	---------
//		123456		  4			 12.3456
//		123456		  0			 123456
//		123456		  1			 12345.6
//
//	The new BigIntNum instance returned by this method
//	will contain United States (US) default numeric
//	separators (decimal separator, integer separator,
//	integer grouping and currency symbol).
//
//	To implement other national or cultural numeric
//	separators, see method:
//		BigIntNum.NewNumStrWithNumSeps()
func (bNum *BigIntNum) NewUint32(uint32Num uint32, precision uint) BigIntNum {

	return new(BigIntNum).NewBigInt(
		big.NewInt(0).SetUint64(uint64(uint32Num)),
		precision)
}

//	NewUint32Exponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//		numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	uint32Num			uint32
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		biNum := new(BigIntNum).NewUint32Exponent(123456, -3)
//		biNum is now equal to "123.456", precision = 3
//
//		biNum := new(BigIntNum).NewUint32Exponent(123456, 3)
//		biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//
//								BigIntNum
//		uintNum		exponent	  Result
//		--------	--------	---------
//		123456		  -3		123.456
//		123456		   3		123456.000
//		123456		   0		123456
func (bNum *BigIntNum) NewUint32Exponent(uint32Num uint32, exponent int) BigIntNum {

	baseBInt := big.NewInt(int64(uint32Num))

	b2 := new(BigIntNum).New()

	b2.SetBigIntExponent(baseBInt, exponent)

	return b2
}

//	NewUint64
//
//	Creates a new BigIntNum instance initialized to the
//	value of input parameter 'uint64Num' which is passed
//	as type 'uint64' (64-bit Unsigned Integer).
//
//	Input parameter 'precision' indicates the number of
//	digits to be formatted to the right of the decimal
//	place and is passed as type 'uint64'.
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//		uint64Num := uint64(123456)
//		precision := uint(3)
//		bINum := BigIntNum{}.NewUint64(uint64Num, precision)
//		bINum is now equal to 123.456
//
//	Examples:
//
//								BigIntNum
//		uintNum		exponent	  Result
//		--------	--------	---------
//		123456		   4		12.3456
//		123456		   0		123456
//		123456		   1		12345.6
func (bNum *BigIntNum) NewUint64(uint64Num uint64, precision uint) BigIntNum {

	return new(BigIntNum).NewBigInt(
		big.NewInt(0).SetUint64(uint64Num),
		precision)
}

//	NewUint64Exponent
//
//	This method returns a new BigIntNum instance in which
//	the numeric value is set using an integer multiplied
//	by 10 raised to the power of the 'exponent' parameter.
//
//		numeric value = integer X 10^exponent
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	uint64Num			uint64
//
//	exponent			int
//
// ----------------------------------------------------------------
//
// # Usage
//
//	This method is designed to be used in conjunction with
//	the new(BigIntNum) syntax thereby allowing BigIntNum
//	type creation and initialization in one step.
//
//
//		biNum := new(BigIntNum).NewUint64Exponent(123456, -3)
//		biNum is now equal to "123.456", precision = 3
//
//		biNum := new(BigIntNum).NewUint64Exponent(123456, 3)
//		biNum is now equal to "123456.000", precision = 3
//
//	Examples:
//
//								BigIntNum
//		uintNum		exponent	  Result
//		--------	--------	---------
//		123456		   -3		123.456
//		123456			3		123456.000
//		123456			0		123456
func (bNum *BigIntNum) NewUint64Exponent(uint64Num uint64, exponent int) BigIntNum {

	baseBInt := big.NewInt(0).SetUint64(uint64Num)

	b2 := new(BigIntNum).New()

	b2.SetBigIntExponent(baseBInt, exponent)

	return b2
}

// Reset - Resets the current BigIntNum to a new
// valid BigIntNum using the BigIntNum components
// BigIntNum.bigInt and BigIntNum.precision. This
// method is usually called after method bNum.IsValid()
// returns false.
func (bNum *BigIntNum) Reset() {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), uint(0))
		return
	}

	if bNum.sign != 1 && bNum.sign != -1 {
		newNum := big.NewInt(0).Set(bNum.bigInt)
		bNum.SetBigInt(newNum, bNum.precision)
		return
	}

	if bNum.absBigInt == nil || bNum.scaleFactor == nil {
		newNum := big.NewInt(0).Set(bNum.bigInt)
		bNum.SetBigInt(newNum, bNum.precision)
		return
	}

	return
}

// RoundToDecPlace
//
// Rounds the current BigIntNum instance to a specified
// number of decimal places.
//
// The rounding procedure applied is 'HalfAwayFromZero'.
//
// 'precision' equals the number of digits to the right
// of the decimal place.
//
// Example:
//
//	integer= 123456
//	precision = 3
//	Numeric Value= 123.456
//
// If the value of BigIntNum.bigInt is zero ('0'), that
// zero value will remain unaltered. However, the
// BigIntNum.precision value will be set equal to input
// parameter, 'precision'.
//
// If the number of decimal places specified for rounding
// ('precision') is equal to the current
// BigIntNum.precision, no action is taken.
//
// If the number of decimal places specified for rounding
// ('precision') is greater than the current
// BigIntNum.precision value, trailing zeros are added to
// the current BigIntNum.bigInt value and
// BigIntNum.precision is set equal to input parameter,
// 'precision'.
//
// Finally, if the number of decimal places specified for
// rounding ('precision') is less than the current
// BigIntNum.precision value, the fractional digits will
// be rounded in accordance with the input parameter,
// 'precision'.
//
// Examples:
//
//	  Original		 'precision'		Resulting
//	   Value		input parameter		   Value
//	--------------	---------------		----------
//	 654.123456			   9			 654.123456000
//	 654.123456			   4			 654.1235
//	-654.123456			   9			-654.123456000
//	-654.123456			   4			-654.1235
//	   0				   3			   0.000
//	   0.000000			   0			   0
//
// Existing numeric separators (decimal separator,
// thousands separator and currency symbol) remain
// unchanged and are not altered by this method.
func (bNum *BigIntNum) RoundToDecPlace(precision uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.precision == precision {
		// Nothing to do. Specified 'precision' is already implemented.
		return
	}

	bNum.SetNumericSeparatorsToDefaultIfEmpty()
	numSeps := bNum.GetNumericSeparatorsDto()

	// bigInt == zero, set precision an return
	if bNum.bigInt.Cmp(big.NewInt(0)) == 0 {

		bNum.CopyIn(new(BigIntNum).NewBigInt(big.NewInt(0), precision))

		_ = bNum.SetNumericSeparatorsDto(numSeps)

		return
	}

	// If existing precision is less than new specified precision,
	// add trailing zeros, set new precision parameter and return.
	if bNum.precision < precision {

		deltaPrecision := precision - bNum.precision

		bNum.ExtendPrecision(deltaPrecision)

		_ = bNum.SetNumericSeparatorsDto(numSeps)

		return
	}

	// Must be: bNum.precision >  precision

	bigNumRound5 :=
		new(BigIntNum).NewBigInt(big.NewInt(5), precision+1)

	bigNumBase :=
		new(BigIntNum).NewBigInt(bNum.absBigInt, bNum.precision)

	result :=
		new(BigIntMathAdd).AddBigIntNums(bigNumBase, bigNumRound5)

	// 10^deltaPrecision
	scaleVal := big.NewInt(0).Exp(big.NewInt(10),
		big.NewInt(int64(bNum.precision-precision)), nil)

	result.bigInt = big.NewInt(0).Quo(result.bigInt, scaleVal)

	if bNum.sign < 0 {
		result.bigInt = big.NewInt(0).Neg(result.bigInt)
	}

	_ = bNum.SetNumericSeparatorsDto(numSeps)

	bNum.SetBigInt(result.bigInt, precision)
}

// SetBigInt - Sets the value of the current BigIntNum instance using
// the input parameters *big.Int integer and precision.
//
// The 'precision' parameter specifies the number of digits to the right
// of the decimal place. The Numeric value is equal to bigI x 10^(precision x -1).
// This effectively locates the decimal place by counting from the extreme right
// of the integer number, 'precision' places to the left. See the example below.
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	bigI				*big.Int
//
//		'bigI' is a type *big.Int and represents the
//		integer value of the number; that is, the
//		numeric value without decimal digits.
//
//
//	precision			uint
//
//		This unsigned integer (always a positive value)
//		identifies the location of the decimal place in
//		the integer value 'bigI'.
//
//		The decimal place location is calculated by
//		starting with the right most digit in the
//		integer number and counting	left, 'precision'
//		places.
//
//		Example:
//
//			Integer					Numeric
//			Value		precision	Value
//			123456			3		123.456
func (bNum *BigIntNum) SetBigInt(bigI *big.Int, precision uint) {

	numSeps := bNum.GetNumericSeparatorsDto()

	bNum.Empty()

	bNum.bigInt = big.NewInt(0).Set(bigI)
	bNum.precision = precision
	base10 := big.NewInt(0).SetInt64(int64(10))
	bigPrecision := big.NewInt(0).SetInt64(int64(bNum.precision))
	bNum.scaleFactor = big.NewInt(0).Exp(base10, bigPrecision, nil)
	bNum.numberOfExpectedDigits = big.NewInt(0)
	result := bNum.bigInt.Cmp(big.NewInt(0))

	if result == -1 {
		bNum.sign = -1
		minusOne := big.NewInt(0).SetInt64(-1)
		bNum.absBigInt = big.NewInt(0).Mul(bNum.bigInt, minusOne)
	} else {
		bNum.sign = 1
		bNum.absBigInt = big.NewInt(0).Set(bNum.bigInt)
	}

	_ = bNum.SetNumericSeparatorsDto(numSeps)

}

// SetBigIntExponent - Sets the numeric value using an integer
// multiplied by 10 raised to the power of the 'exponent'
// parameter.
//
//	numeric value = integer X 10^exponent
//
// Input parameter 'bigI' is of type *big.Int.
//
// Input parameter 'exponent' is of type int.
//
// If exponent is less than +1, precision is set equal to exponent and
// bigI is unchanged. Example:
//
//	   bigI				exponent			BigIntNum Result
//		 123456		 		  -3							123.456
//
// If exponent is greater than 0, bigI is multiplied by 10 raised to the
// power of exponent and precision is set equal to exponent.
//
//	   bigI				exponent			BigIntNum Result
//		 123456		 		   3							123456.000
func (bNum *BigIntNum) SetBigIntExponent(bigI *big.Int, exponent int) {

	if exponent < 1 {
		precision := uint(exponent * -1)
		bNum.SetBigInt(bigI, precision)
		return
	}

	// exponent must be greater than zero.
	// scale left exponent places and set precision to zero

	big10 := big.NewInt(10)
	scale := big.NewInt(int64(exponent))
	scaleValue := big.NewInt(0).Exp(big10, scale, nil)
	newBigI := big.NewInt(0).Mul(bigI, scaleValue)

	bNum.SetBigInt(newBigI, uint(exponent))

	return
}

//	SetBigFloat
//
//	Sets the value of a BigIntNum using a *big.Float
//	floating point input parameter.  The precision of
//	the number is specified by the input parameter,
//	'maxPrecision'.
//
//	Input parameter 'bigFloat' is first converted to a
//	rational number of type *Rat. Next, the floating
//	point value is calculated as in integer value and
//	copied to the current instance of BigIntNum.
//
//	As part of the conversion of a BigFloat to a rational
//	number the Accuracy flag is returned describing the
//	rounding error associated with the conversion. The
//	Accuracy Flag is set as:
//
//		Below Accuracy == -1
//		Exact Accuracy == 0
//		Above Accuracy == +1
//
//	If Accuracy == 0, no error is issued by this method.
//	However, if Accuracy == -1 or Accuracy == +1 an error
//	will be returned. All conversions must be exact.
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
//
// ----------------------------------------------------------------
//
// # Input Parameters
//
//	bigFloat			*big.Float
//
//		The numeric value of this big.Float parameter will be
//		copied to the current instance of BigIntNum.
//
// maxPrecision uint  - The maximum precision for the resulting BigIntNum after
//
//	conversion of input parameter 'ratNum'. Precision will
//	never be greater than 'maxPrecision'; however, actual
//	precision may be less than 'maxPrecision'.
func (bNum *BigIntNum) SetBigFloat(
	bigFloat *big.Float,
	maxPrecision uint) error {

	ePrefix := "NumStrDto.NewBigFloat() "

	rat, accuracyFlag := bigFloat.Rat(nil)

	if accuracyFlag == -1 {
		return errors.New(ePrefix +
			"Error: Conversion of input parameter 'bigFloat' resulted in Accuracy Flag == -1 " +
			"or 'Below Accuracy' Conversion is NOT Exact!")
	}

	if accuracyFlag == 1 {
		return errors.New(ePrefix +
			"Error: Conversion of input parameter 'bigFloat' resulted in Accuracy Flag == +1 " +
			"or 'Above Accuracy' Conversion is NOT Exact!")
	}

	err := bNum.SetBigRat(rat, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bNum.SetBigRat(rat, maxPrecision). "+
			"Error='%v' \n", err.Error())
	}

	return nil
}

// SetBigRat - Sets the current BigIntNum value to that of input parameter
// 'ratNum', a rational number or *big.Rat type.
//
// maxPrecision uint  - The maximum precision for the resulting BigIntNum after
//
//	conversion of input parameter 'ratNum'. Precision will
//	never be greater than 'maxPrecision'; however, actual
//	precision may be less than 'maxPrecision'.
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) SetBigRat(ratNum *big.Rat, maxPrecision uint) error {

	ePrefix := "BigIntNum.SetBigRat() "

	numSeps := bNum.GetNumericSeparatorsDto()

	numerator := big.NewInt(0).Set(ratNum.Num())

	denominator := big.NewInt(0).Set(ratNum.Denom())

	biPair := BigIntPair{}.NewBase(numerator, 0, denominator, 0)
	biPair.MaxPrecision = maxPrecision

	biNum, err := BigIntMathDivide{}.pairFracQuotientNoNumSeps(biPair)
	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.pairFracQuotientNoNumSeps(biPair). "+
			"Error='%v'\n", err.Error())
	}

	if biNum.GetPrecisionUint() > maxPrecision {
		biNum.SetPrecision(maxPrecision)
	}

	_ = biNum.SetNumericSeparatorsDto(numSeps)

	bNum.CopyIn(biNum)

	return nil
}

// SetCurrencySymbol - assigns the input parameter rune as the
// currency symbol to be used by the BigIntNum when generating
// number strings for display.
//
// In the USA, the currency symbol is the dollar sign ('$').
//
// Note: If a zero value is submitted as input, Currency Symbol
// will default to the USA dollar sign ('$').
//
// Example: $123.45
func (bNum *BigIntNum) SetCurrencySymbol(currencySymbol rune) {

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	bNum.currencySymbol = currencySymbol
}

// SetDecimalSeparator - Assigns a rune or character to the internal
// data field, 'decimalSeparator'. The Decimal Separator is used to
// separate the integer and fractional elements of a number string.
//
// The BigIntNum Type uses this character when generating number strings
// for display.
//
// In the USA, the Decimal Separator is a period character ('.').
//
// Note: If a zero value is submitted as input, the Decimal Separator
// will default to the USA standard period character ('.').
//
// Example: 123.456
func (bNum *BigIntNum) SetDecimalSeparator(decimalSeparator rune) {

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	bNum.decimalSeparator = decimalSeparator
}

// SetIntFracStrings - Sets the value of the current BigIntNum instance based on
// a numeric value represented by separate integer and fractional components.
//
// Input parameters 'intStr' and 'fracStr' are strings representing the integer and
// fractional components. They are combined by this method to create a numeric value
// which is assigned to the current BigIntNum instance.
//
// Input parameter 'signVal' must be set to one of two values: +1 or -1. This value is
// used to signal the sign of the resulting numeric value. +1 generates a positive number
// and -1 generates a negative number. If input parameters 'inStr' or 'fracStr' contain
// a leading minus or plus sign character, it will be ignored. The sign of the resulting
// numeric value is controlled strictly by input parameter, 'signVal'.
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) SetIntFracStrings(intStr, fracStr string, signVal int) error {

	ePrefix := "BigIntNum.SetIntFracStrings() "

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
					cleanIntRuneAry = append(cleanIntRuneAry, bNum.GetDecimalSeparator())
					isFirstRune = false
				}

				cleanIntRuneAry = append(cleanIntRuneAry, rune(fracStr[j]))
			}

		}
	}

	err := bNum.SetNumStr(string(cleanIntRuneAry))

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bNum.SetNumStr(string(cleanIntRuneAry)). "+
			"cleanIntRuneAry='%v' Error='%v' ", string(cleanIntRuneAry), err.Error())
	}

	return nil
}

// SetFloat32 - Sets the value of a BigIntNum using a float32 floating point
// input parameter.  The precision of the number is specified by the input
// parameter, 'decimalPlaces'.
//
// Input Parameters
// ================
//
// f32 float32				- This float32 value will be converted into an instance of
//
//	BigIntNum.
//
// maxPrecision uint  - The maximum precision for the resulting BigIntNum after
//
//	conversion of input parameter 'ratNum'. Precision will
//	never be greater than 'maxPrecision'; however, actual
//	precision may be less than 'maxPrecision'.
func (bNum *BigIntNum) SetFloat32(f32 float32, maxPrecision uint) error {

	ePrefix := "BigIntNum.SetFloat32() "

	rat := big.NewRat(1, 1).SetFloat64(float64(f32))

	err := bNum.SetBigRat(rat, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bNum.SetBigRat(rat, maxPrecision). "+
			"Error='%v' \n", err.Error())
	}

	return nil
}

// SetFloat64 - Sets the value of a BigIntNum using a float64 floating point
// input parameter.  The precision of the number is specified by the input
// parameter 'decimalPlaces'.
//
// Input Parameters
// ================
//
// f64 float64				- This float64 value will be converted into an instance of
//
//	BigIntNum.
//
// maxPrecision uint  - The maximum precision for the resulting BigIntNum after
//
//	conversion of input parameter 'f64'. Resulting precision
//	will never be greater than 'maxPrecision'; however, actual
//	precision may be less than 'maxPrecision'.
func (bNum *BigIntNum) SetFloat64(f64 float64, maxPrecision uint) error {

	ePrefix := "BigIntNum.SetFloat64() "

	rat := big.NewRat(1, 1).SetFloat64(f64)

	err := bNum.SetBigRat(rat, maxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bNum.SetBigRat(rat, maxPrecision). "+
			"Error='%v' \n", err.Error())
	}

	return nil
}

// SetNumStr - Initializes the current BigIntNum instance
// of the numeric value of the number string input parameter.
// A number string is a string of numeric digits which may
// or may not be prefixed with a minus sign ('-'). The numeric
// string of digits may also contain a decimal separator such
// as a period ('.'). The decimal separator may be set by the
// user. See Method BigIntNum.SetDecimalSeparator(). The decimal
// separator is used to separate integer and fractional numeric
// digits within the number string.
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) SetNumStr(numStr string) error {

	ePrefix := "BigIntNum.SetNumStr() "

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), 0)
	}

	if len(numStr) == 0 {
		return errors.New(ePrefix + "Error: Input parameter 'numStr' is an EMPTY string!")
	}

	baseRunes := []rune(numStr)
	lBaseRunes := len(baseRunes)

	numSeps := bNum.GetNumericSeparatorsDto()

	newSign := 1
	newPrecision := uint(0)
	newAbsBigInt := big.NewInt(0)
	baseTen := big.NewInt(10)
	isStartNumericDigits := false
	isEndNumericDigits := false
	isFractionalValue := false
	hasMinusSign := false
	hasLeftParen := false
	hasRightParen := false
	numOfNumericDigits := 0

	for i := 0; i < lBaseRunes; i++ {

		if isEndNumericDigits {
			continue
		}

		if baseRunes[i] == ',' && bNum.decimalSeparator != ',' {
			continue
		}

		if baseRunes[i] == '-' {
			hasMinusSign = true
			continue
		}

		if baseRunes[i] == '(' {

			if isStartNumericDigits == false {
				hasLeftParen = true
			}
			continue
		}

		if baseRunes[i] == ')' {

			if isStartNumericDigits == true &&
				hasLeftParen == true {
				hasRightParen = true
				isEndNumericDigits = true
			}

			continue
		}

		if baseRunes[i] == bNum.decimalSeparator {
			isFractionalValue = true
			continue
		}

		if baseRunes[i] >= '0' && baseRunes[i] <= '9' &&
			isEndNumericDigits == false {

			newAbsBigInt = big.NewInt(0).Mul(newAbsBigInt, baseTen)

			newAbsBigInt = big.NewInt(0).Add(newAbsBigInt,
				big.NewInt(int64(baseRunes[i]-48)))

			isStartNumericDigits = true
			numOfNumericDigits++

			if isFractionalValue {
				newPrecision++
			}

		}

	}

	if numOfNumericDigits == 0 {
		return fmt.Errorf(ePrefix+
			"Error: No numeric digits were found in input parameter 'numStr'. "+
			"numStr='%v'", numStr)
	}

	if hasMinusSign == true ||
		(hasLeftParen == true && hasRightParen == true) {
		newSign = -1
	}

	bNum.Empty()
	bNum.sign = newSign
	bNum.precision = newPrecision
	bNum.absBigInt = big.NewInt(0).Set(newAbsBigInt)

	if bNum.sign == 1 {
		bNum.bigInt = big.NewInt(0).Set(newAbsBigInt)
	} else {
		bNum.bigInt = big.NewInt(0).Neg(newAbsBigInt)
	}

	bNum.scaleFactor = big.NewInt(0).Exp(baseTen,
		big.NewInt(int64(newPrecision)),
		nil)

	err := bNum.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by err := bNum.SetNumericSeparatorsDto(numSeps) "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// SetExpectedNumberOfDigits - Sets the number of expected digits associated with the
// Absolute Value of this 'BigIntNum.absBigInt'. The value is stored in the data
// field, 'BigIntNum.numberOfExpectedDigits'.
//
// Useful in tracking leading zeros.
func (bNum *BigIntNum) SetExpectedNumberOfDigits(numOfDigits *big.Int) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	bNum.numberOfExpectedDigits = big.NewInt(0).Set(numOfDigits)
}

// SetExpectedToActualNumberOfDigits - Sets the 'Expected' number of numeric
// digits associated with this BigIntNum, to the actual number of numeric digits
// in the BigIntNum value at the time when this method is called.
func (bNum *BigIntNum) SetExpectedToActualNumberOfDigits() {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	actNumOfDigits, _, _ := bNum.GetActualNumberOfDigits()

	bNum.numberOfExpectedDigits = big.NewInt(0).Set(actNumOfDigits)

}

// SetPrecision - Sets a new 'precision' value for the current
// BigIntNum instance. The new 'precision' is specified by the
// uint type input parameter, 'newPrecision'.
//
// Precision is defined as the number of numeric digits to right
// of the decimal place.
//
// If 'newPrecision' is equal to the current BigIntNum.precision value,
// no action is taken and the original BigIntNum numeric value remains
// unchanged.
//
// If 'newPrecision' is greater than the current BigIntNum.precision
// value, trailing zeros are added to the fractional digits to the
// right of the decimal place.
//
// If 'newPrecision' is less than the current BigIntNum precision
// value, the current BigIntNum numeric value is rounded to the
// specified 'newPrecision' value.
//
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
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) SetPrecision(newPrecision uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if newPrecision == bNum.precision {
		return
	}

	if bNum.precision > newPrecision {
		bNum.RoundToDecPlace(newPrecision)
		return
	}

	deltaPrecision := newPrecision - bNum.precision
	// bNum.precision must be less than newPrecision
	bNum.ExtendPrecision(deltaPrecision)
	return

}

// SetNumericSeparators
//
// Used to assign values for the Decimal and Thousands
// separators as well as the Currency Symbol to be used
// in displaying the current number string.
//
// Different nations and cultures use different symbols
// to delimit numerical values. In the USA and many other
// countries, a period character ('.') is used to delimit
// integer and fractional digits within a numeric value
// (123.45). Likewise, thousands may be delimited by a
// comma (','). Currency signs very by nationality. For
// instance, the US, Canada and several other countries
// use the dollar sign ($) as a currency symbol.
//
// Note: If zero values are submitted as input for
// separator values, those values will default to US
// standards.
//
// US Examples:
//
// Decimal Separator period ('.')		= 123.456
// Thousands Separator comma (',')		= 1,000,000,000
// Currency Symbol dollar sign ('$')	= $123
func (bNum *BigIntNum) SetNumericSeparators(
	decimalSeparator,
	thousandsSeparator,
	currencySymbol rune) {

	bNum.SetNumericSeparatorsToDefaultIfEmpty()

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	bNum.decimalSeparator = decimalSeparator
	bNum.thousandsSeparator = thousandsSeparator
	bNum.currencySymbol = currencySymbol
}

// SetNumericSeparatorsDto - Sets the values of numeric separators:
//
//	decimal place separator
//	thousands separator
//	currency symbol
//
// based on values transmitted through input parameter 'customSeparators'.
//
// If any of the values contained in input parameter 'customSeparators' is set
// to zero, an error will be returned.
func (bNum *BigIntNum) SetNumericSeparatorsDto(customSeparators NumericSeparatorDto) error {

	ePrefix := "BigIntNum.SetNumericSeparatorsDto() "

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

	bNum.decimalSeparator = customSeparators.DecimalSeparator
	bNum.thousandsSeparator = customSeparators.ThousandsSeparator
	bNum.currencySymbol = customSeparators.CurrencySymbol

	return nil
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
func (bNum *BigIntNum) SetNumericSeparatorsToDefaultIfEmpty() {

	if bNum.decimalSeparator == 0 {
		bNum.decimalSeparator = '.'
	}

	if bNum.thousandsSeparator == 0 {
		bNum.thousandsSeparator = ','
	}

	if bNum.currencySymbol == 0 {
		bNum.currencySymbol = '$'
	}

}

// SetNumericSeparatorsToUSADefault
//
// Sets Numeric separators to United States
// defaults:
//
//	Decimal Point Separator = "."
//	Thousands Separator = ","
//	Currency Symbol = "$"
//
// Call specific methods to set numeric separators
// for other countries or cultures:
//
//	bNum.SetDecimalSeparator()
//	bNum.SetThousandsSeparator()
//	bNum.SetCurrencySymbol()
func (bNum *BigIntNum) SetNumericSeparatorsToUSADefault() {

	bNum.SetDecimalSeparator('.')

	bNum.SetThousandsSeparator(',')

	bNum.SetCurrencySymbol('$')
}

// SetSignValue - Sets the sign value of the current BigIntNum
// to positive (+1) or negative (-1).
//
// If a value other than +1 or -1 is transmitted as an input
// parameter, an error will be returned.
func (bNum *BigIntNum) SetSignValue(signVal int) error {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if signVal == 1 || signVal == -1 {

		if bNum.GetSign() == signVal {
			return nil
		}

		bNum.ChangeSign()

		return nil

	}

	return fmt.Errorf("BigIntNum.SetSignValue() Error: Input parameter 'signVal' "+
		"must be +1 or -1. signVal='%v' ", signVal)
}

// ShiftPrecisionLeft - Shifts precision of the current BigIntNum
// numeric value to the left by 'shiftLeftPlaces' decimal places. This
// is a 'relative' shift-left operation. The shift left operation is
// therefore performed with the current decimal place position as the
// starting point.
//
// This operation is equivalent to:	result = Decimal value / 10^shiftLeftPlaces
// or signed number divided by 10 raised to the power of shiftLeftPlaces.
//
// This method performs a relative shift left of the decimal place position.
// Be careful, this is NOT Shift Number Left operation. This is Shift Precision
// Left which means that the decimal place will be shifted left.
//
// See Examples below.
//
// Input Parameters
// ================
//
//	shiftLeftPlaces int	- The number of positions the decimal place will be
//												shifted left from its current position.
//
// Examples:
// =========
//
//	shift-left
//
// signed Number		  places				Result
//
//	"123456.789"				3						"123.456789"
//	"123456.789"				2						"1234.56789"
//	"123456.789"        6					  "0.123456789"
//	"123456789"	 			  6						"123.456789"
//	"123"               5	          "0.00123"
//	"0"								  3						"0"
//	"123456.789"				0						"123456.789"		- zero has no effect on original number string
//
// "-123456.789"        0          "-123456.789"
// "-123456.789"        3          "-123.456789"
// "-123456789"			    6					 "-123.456789"
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) ShiftPrecisionLeft(shiftLeftPlaces uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if shiftLeftPlaces == 0 ||
		bNum.IsZero() {
		return
	}

	newPrecision := bNum.precision + shiftLeftPlaces
	bNum.SetBigInt(bNum.bigInt, newPrecision)

	return
}

// ShiftPrecisionRight - Shifts precision of the current BigIntNum
// numeric value to the right by 'shiftRightPlaces' decimal places. This
// is a 'relative' shift-right operation. The shift right operation is
// therefore performed with the current decimal place position as the
// starting point.
//
// This is equivalent to: result = Decimal value X 10^shiftRightPrecision or
// Decimal numeric value multiplied by 10 raised to the power of
// shiftRightPrecision.
//
// This method performs a relative shift right of the decimal place position.
// Be careful, this is NOT a Shift Number Right operation. This is Shift Precision
// Right which means that the decimal place will be shifted right.
//
// See Examples below.
//
// Input Parameters
// ================
//
//	shiftRightPlaces int	- The number of positions the decimal place will be
//													shifted right from its current position.
//
// Examples:
// =========
//
//	shift-right
//
// signed Number		  places				Result
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
//
// Existing numeric separators (decimal separator, thousands separator
// and currency symbol) remain unchanged and are not altered by this method.
func (bNum *BigIntNum) ShiftPrecisionRight(shiftRightPlaces uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if shiftRightPlaces == 0 ||
		bNum.IsZero() {
		return
	}

	if shiftRightPlaces <= bNum.precision {

		newPrecision := bNum.precision - shiftRightPlaces

		bNum.SetBigInt(bNum.bigInt, newPrecision)

		return
	}

	// shiftRightPlaces > bNum.precision

	newPrecision := shiftRightPlaces - bNum.precision

	bigITen := big.NewInt(10)

	exponent := big.NewInt(int64(newPrecision))

	scaleFactor := big.NewInt(0).Exp(bigITen, exponent, nil)

	newValue := big.NewInt(0).Mul(bNum.bigInt, scaleFactor)

	bNum.SetBigInt(newValue, 0)

	return
}

// SetThousandsSeparator
//
// Sets the value of the character which will be used to
// separate thousands in the display of the NumStrDto number
// string.
//
// In the United States (US) the typical thousands separator
// is the comma character (',').
//
// If a zero input value is submitted, the Thousands
// Separator will default to the comma character.
//
//	Example:
//	1,000,000
func (bNum *BigIntNum) SetThousandsSeparator(thousandsSeparator rune) {

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	bNum.thousandsSeparator = thousandsSeparator

}

// TrimTrailingFracZeros - This method will delete non-significant
// trailing zeros from the fractional digits of the current BigIntNum
// numerical value.
//
// Examples:
//
//	Initial Value			Trimmed Value
//		456.123000 			 456.123
//			0.000					 0
//			7.0						 7
//	 -456.123000			-456.123
func (bNum *BigIntNum) TrimTrailingFracZeros() {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), bNum.precision)
	}

	if bNum.precision == 0 {
		return
	}

	biBaseZero := big.NewInt(0)

	if bNum.bigInt.Cmp(biBaseZero) == 0 {
		bNum.precision = 0
		bNum.scaleFactor = big.NewInt(1)
		return
	}

	// bNum.precision must be greater than zero
	biBase10 := big.NewInt(10)
	scrap := big.NewInt(0)
	newBigIntNum, mod10 := big.NewInt(0).QuoRem(bNum.bigInt, biBase10, scrap)
	doReset := false

	for mod10.Cmp(biBaseZero) == 0 && bNum.precision > 0 {
		bNum.bigInt.Set(newBigIntNum)
		bNum.precision--
		newBigIntNum, mod10 = big.NewInt(0).QuoRem(bNum.bigInt, biBase10, scrap)
		doReset = true
	}

	if doReset {
		if bNum.sign < 0 {
			bNum.absBigInt = big.NewInt(0).Neg(bNum.bigInt)
		} else {
			bNum.absBigInt = big.NewInt(0).Set(bNum.bigInt)
		}

		bigPrecision := big.NewInt(0).SetInt64(int64(bNum.precision))
		bNum.scaleFactor = big.NewInt(0).Exp(biBase10, bigPrecision, nil)

	}

	return
}

// TruncToDecPlace
//
// Truncates the current BigIntNum to the number of
// decimal places specified by input parameter
// 'precision'. No rounding occurs, the trailing digits
// are simply truncated or deleted in order to achieve
// the specified number of decimal places.
//
// 'precision' equals the number of digits to the right
// of the decimal place.
//
// If the value of BigIntNum.bigInt is zero ('0'), that
// zero value will remain unaltered. However,
// BigIntNum.precision will be set equal to input
// parameter, 'precision'.
//
// If the number of decimal places specified for
// truncation ('precision') is equal to the current
// BigIntNum.precision, no action is taken and the
// original BigIntNum numeric value remains unchanged.
//
// If the number of decimal places specified for
// truncation ('precision') is greater than the current
// BigIntNum.precision, trailing zeros will be added to
// the current BigIntNum.bigInt value and
// BigIntNum.precision is set equal to input parameter,
// 'precision'.
//
// If 'precision' is less than the current
// BigIntNum.precision value, the current BigIntNum
// numeric value is truncated to the specified
// 'precision' value and NO rounding occurs.
//
// Examples:
//
//	Original		'newPrecision'		  Resulting
//	Value			input parameter			Value
//	--------------	---------------     -------------
//	 654.123456			 9				 654.123456000
//	 654.123456			 4				 654.1234 (no rounding)
//	-654.123456			 9				-654.123456000
//	-654.123456			 4				-654.1234 (no rounding)
//	   0				 3				   0.000
//	   0.000000			 0				   0
func (bNum *BigIntNum) TruncToDecPlace(precision uint) {

	if bNum.bigInt == nil {
		bNum.SetBigInt(big.NewInt(0), precision)
	}

	if bNum.precision == precision {
		// Nothing to do. Specified 'precision' is already implemented.
		return
	}

	// bigInt == zero, set precision an return
	if bNum.bigInt.Cmp(big.NewInt(0)) == 0 {
		bNum.precision = precision
		return
	}

	// If existing precision is less than new specified precision,
	// add trailing zeros, set new precision parameter and return.
	if bNum.precision < precision {

		bNum.ExtendPrecision(precision - bNum.precision)
		return
	}

	// Must be bNum.precision > precision
	base10 := big.NewInt(10)
	deltaPrecision := big.NewInt(int64(bNum.precision - precision))
	newBigInt := big.NewInt(0).Set(bNum.absBigInt)
	newScaleVal := big.NewInt(0).Exp(base10, deltaPrecision, nil)
	newBigInt = big.NewInt(0).Quo(newBigInt, newScaleVal)

	if bNum.sign < 1 {
		newBigInt = big.NewInt(0).Neg(newBigInt)
	}

	bNum.SetBigInt(newBigInt, precision)
}
