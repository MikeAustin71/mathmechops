package mathops

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

/*
	The source code repository for numstrdto.go is located at:
			https://github.com/MikeAustin71/mathopsgo.git

	The source file decimal.go is located in directory:
		MikeAustin71/mathopsgo/mathops/numstrdto.go
*/

// NumStrDto - This Type contains data fields and methods used
// to manage, store and transport number strings.
//
// The NumStrDto Type implements the INumMgr interface.
type NumStrDto struct {
	signVal int // An integer value indicating the numeric sign of this number string.
	// 		Valid values are +1 or -1
	absAllNumRunes []rune // An array of runes containing all the numeric digits in a number with
	//		no preceding plus or minus sign character. Example: 123.456 =
	//		[]rune{'1','2','3','4','5','6'}
	precision          uint // The number of digits to the right of the decimal point.
	thousandsSeparator rune // Separates thousands in the integer number: '1,000,000,000
	decimalSeparator   rune // Separates integer and fractional elements of a number. '123.456'
	currencySymbol     rune // Currency symbol used in currency string displays
}

// Add - Adds the value of input NumStrDto to the current NumStrDto
// instance
func (nDto *NumStrDto) Add(n2Dto NumStrDto) error {
	ePrefix := "NumStrDto.Add() "
	n1Dto := nDto.CopyOut()

	nResult, err := nDto.AddNumStrs(n1Dto, n2Dto)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by nDto.AddNumStrs(n1Dto, n2Dto). "+
			"Error='%v'", err.Error())
	}

	nDto.CopyIn(nResult)

	return nil
}

// AddNumStrs - Adds the values represented by two NumStrDto objects and
// returns the result as an NumStrDto.
func (nDto *NumStrDto) AddNumStrs(n1Dto NumStrDto, n2Dto NumStrDto) (NumStrDto, error) {

	n1DtoSetup, n2DtoSetup, _, _, err := nDto.FormatForMathOps(n1Dto, n2Dto)

	if err != nil {
		return NumStrDto{}, fmt.Errorf("AddNumStrs() - Error returned from nDto.FormatForMathOps(n1Dto, n2Dto). Error= %v", err)
	}

	newSignVal := n1DtoSetup.signVal

	if n1DtoSetup.signVal != n2DtoSetup.signVal {
		_ = n1DtoSetup.SetSignValue(1)
		_ = n2DtoSetup.SetSignValue(1)
		nDtoOut, err := nDto.SubtractNumStrs(n1DtoSetup, n2DtoSetup)

		if err != nil {
			return NumStrDto{}, fmt.Errorf("AddNumStrs() - Error returned from nDto.SubtractNumStrs(n1DtoSetup, n2DtoSetup). Error= %v", err)
		}

		if nDto.IsNumStrZeroValue(&nDtoOut) {
			newSignVal = 1
		}

		_ = nDtoOut.SetSignValue(newSignVal)

		return nDtoOut, nil
	}

	precision := n1DtoSetup.precision
	lenN1AllRunes := len(n1DtoSetup.absAllNumRunes)

	n3IntAry := make([]int, lenN1AllRunes+1)
	carry := 0
	n1 := 0
	n2 := 0
	n3 := 0

	for j := lenN1AllRunes - 1; j >= 0; j-- {

		n1 = int(n1DtoSetup.absAllNumRunes[j]) - 48
		n2 = int(n2DtoSetup.absAllNumRunes[j]) - 48

		n3 = n1 + n2 + carry

		carry = 0

		if n3 > 9 {
			n3 = n3 - 10
			carry = 1
		}

		n3IntAry[j+1] = n3

	}

	if carry > 0 {
		n3IntAry[0] = carry
	}

	return nDto.FindIntArraySignificantDigitLimits(n3IntAry, precision, newSignVal)

}

// CompareSignedValues - compares the signed numeric values
// of two NumStrDto objects.
//
// Return Values:
// -1 = n1Dto is less than n2Dto
//
//	0 = n1Dto is equal to n2Dto
//	1 = n1Dto is greater than n2Dto
//
// Examples:
//
//		n1        			n2           	Result
//		-9691.23				91.245				 	-1
//	 9691.23					91.245					 1
//	 -5							82							-1
//	  5							 5							 0
func (nDto *NumStrDto) CompareSignedValues(n1Dto, n2Dto *NumStrDto) int {

	cmpAbs := nDto.CompareAbsoluteValues(n1Dto, n2Dto)

	if cmpAbs == 0 {

		if n1Dto.signVal == n2Dto.signVal {
			return 0
		} else {
			// n1Dto.signVal != n2Dto.signVal
			if n1Dto.signVal == 1 {
				return 1
			}

			// n2Dto.signVal must == 1
			return -1

		}

	}

	if cmpAbs == 1 {

		if n1Dto.signVal == n2Dto.signVal {

			if n1Dto.signVal == 1 {
				return 1
			}

			// must be n1Dto.signVal == n2Dto.signVal && n1Dto.signVal == -1

			return -1

		}

		// must be n1Dto.signVal != n2Dto.signVal
		if n1Dto.signVal == 1 {
			return 1
		} else {
			// must be n2Dto.signVal == 1
			return -1
		}
	}

	// cmpAbs == -1

	if n2Dto.signVal == n1Dto.signVal {

		if n2Dto.signVal == 1 {
			// n1Dto.signVal && n2Dto.signVal must equal 1
			return -1
		} else {
			// n1Dto.signVal && n2Dto.signVal must equal -1
			return 1
		}

	}

	// must be n2Dto.signVal != n1Dto.signVal

	if n2Dto.signVal == -1 {
		return 1
	}

	// must be n2Dto.signVal == 1
	return -1
}

// CompareAbsoluteValues - compares the absolute numeric values
// of two NumStrDto objects. The signs (+ or -) of the two
// compared numeric values are ignored. Only the absolute
// numeric values are compared.
// Return Values:
// -1 = n1Dto is less than n2Dto
//
//	0 = n1Dto is equal to n2Dto
//	1 = n1Dto is greater than n2Dto
//
// Examples:
//
//		n1        			n2           	Result
//		-9691.23				91.245				 	 1
//	 9691.23					91.245					 1
//	 -5							82							-1
//	  5							 5							 0
func (nDto *NumStrDto) CompareAbsoluteValues(n1Dto, n2Dto *NumStrDto) int {

	n1DtoAbsFracRunes := n1Dto.GetAbsFracRunes()
	n2DtoAbsFracRunes := n2Dto.GetAbsFracRunes()

	n1DtoAbsIntRunes := n1Dto.GetAbsIntRunes()
	n2DtoAbsIntRunes := n2Dto.GetAbsIntRunes()

	lenN1IntRunes := len(n1DtoAbsIntRunes)
	lenN2IntRunes := len(n2DtoAbsIntRunes)

	isN1Zero := nDto.IsNumStrZeroValue(n1Dto)
	isN2Zero := nDto.IsNumStrZeroValue(n2Dto)

	if !isN1Zero && isN2Zero {
		return 1
	}

	if isN1Zero && !isN2Zero {
		return -1
	}

	if isN1Zero && isN2Zero {
		return 0
	}

	if lenN1IntRunes > lenN2IntRunes {
		return 1
	}

	if lenN1IntRunes < lenN2IntRunes {
		return -1
	}

	// lenN1IntRunes Must Be Equal to lenN2IntRunes

	for i := 0; i < lenN1IntRunes; i++ {
		n1 := n1DtoAbsIntRunes[i] - 48
		n2 := n2DtoAbsIntRunes[i] - 48

		if n1 > n2 {
			return 1
		}

		if n1 < n2 {
			return -1
		}
	}

	// All the integers are equal
	lenN1FracRunes := len(n1DtoAbsFracRunes)
	lenN2FracRunes := len(n2DtoAbsFracRunes)

	lenFracRunesToTest := lenN1FracRunes

	if lenN2FracRunes < lenN1FracRunes {
		lenFracRunesToTest = lenN2FracRunes
	}

	for j := 0; j < lenFracRunesToTest; j++ {
		n1 := n1DtoAbsFracRunes[j] - 48
		n2 := n2DtoAbsFracRunes[j] - 48
		if n1 > n2 {
			return 1
		}

		if n1 < n2 {
			return -1
		}

	}

	if lenN1FracRunes > lenN2FracRunes {
		return 1
	}

	if lenN1FracRunes < lenN2FracRunes {
		return -1
	}

	return 0
}

// CopyOut - Creates a copy of the current
// NumStrDto fields and returns a completely
// new instance of NumStrDto
func (nDto *NumStrDto) CopyOut() NumStrDto {
	nOut := NumStrDto{}

	nOut.signVal = nDto.signVal
	nOut.absAllNumRunes = nDto.absAllNumRunes
	nOut.precision = nDto.precision
	nOut.thousandsSeparator = nDto.thousandsSeparator
	nOut.decimalSeparator = nDto.decimalSeparator
	nOut.currencySymbol = nDto.currencySymbol

	return nOut
}

// CopyIn - Receives an incoming NumStrDto object
// and copies the information to the current NumStrDto
// data fields.
func (nDto *NumStrDto) CopyIn(nInDto NumStrDto) {

	nDto.Empty()

	nDto.signVal = nInDto.signVal
	nDto.absAllNumRunes = nInDto.absAllNumRunes
	nDto.precision = nInDto.precision
	nDto.thousandsSeparator = nInDto.thousandsSeparator
	nDto.decimalSeparator = nInDto.decimalSeparator
	nDto.currencySymbol = nInDto.currencySymbol

}

// Divide
//
// This method divides the current NumStrDto by input
// parameter 'n2Dto'. Maximum precision of the division
// result is controlled by the input parameter,
// 'maximumPrecision'.
//
// If 'maximumPrecision' is greater than or equal to zero
// ('0'), the number of digits to the right of the decimal
// place will not exceed 'maximumPrecision'.
//
// 'maximumPrecision' is set equal to minus one ('-1'),
// will be set to a maximum of 1,024 digits to the right
// of the decimal point.
//
// 'minimumPrecision' specifies the minimum precision of
// the final result.
//
// If 'minimumPrecision' is less than zero, it is
// automatically set to zero.
func (nDto *NumStrDto) Divide(n2Dto NumStrDto, minimumPrecision, maximumPrecision int) error {

	ePrefix := "NumStrDto.Divide() "

	ia1, err := IntAry{}.NewNumStrDto(nDto.CopyOut())

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by IntAry{}.NewNumStrDto(nDto.CopyOut()). "+
			"Error='%v'", err.Error())
	}

	ia2, err := IntAry{}.NewNumStrDto(n2Dto)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by IntAry{}.NewNumStrDto(n2Dto). "+
			"Error='%v'", err.Error())
	}

	iaResult, err := ia1.DivideThisBy(&ia2, minimumPrecision, maximumPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"ia1.DivideThisBy(&ia2, minimumPrecision, maximumPrecision). "+
			"minimumPrecision='%v' maximumPrecision='%v' Error='%v'",
			minimumPrecision, maximumPrecision, err.Error())
	}

	nResultDto, err := iaResult.GetNumStrDto()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by iaResult.GetNumStrDto(). "+
			"Error='%v'", err.Error())
	}

	if nResultDto.GetPrecision() < minimumPrecision {

		err = nResultDto.SetThisPrecision(uint(minimumPrecision), false)

		if err != nil {
			return fmt.Errorf(ePrefix+
				"Error returned by nResultDto.SetThisPrecision(uint(minimumPrecision), false). "+
				"Error='%v'", err.Error())
		}

	}

	nDto.CopyIn(nResultDto)

	return nil
}

// Equal - Returns true if the input parameter NumStrDto instance
// is equal in all respects to the current NumStrDto Instance
func (nDto *NumStrDto) Equal(n2Dto NumStrDto) bool {

	err := nDto.IsValid("")

	if err != nil {
		return false
	}

	err = n2Dto.IsValid("")

	if err != nil {
		return false
	}

	if nDto.GetNumStr() != n2Dto.GetNumStr() {
		return false
	}

	lenAbsRuneArray := len(nDto.absAllNumRunes)

	if nDto.signVal != n2Dto.signVal ||
		lenAbsRuneArray != len(n2Dto.absAllNumRunes) ||
		nDto.precision != n2Dto.precision ||
		nDto.thousandsSeparator != n2Dto.thousandsSeparator ||
		nDto.decimalSeparator != n2Dto.decimalSeparator ||
		nDto.currencySymbol != n2Dto.currencySymbol {
		return false
	}

	for i := 0; i < lenAbsRuneArray; i++ {
		if nDto.absAllNumRunes[i] != n2Dto.absAllNumRunes[i] {
			return false
		}
	}

	return true
}

// Empty - Sets all the fields in the NumStrDto
// to their initial or zero state.
func (nDto *NumStrDto) Empty() {
	nDto.signVal = 0
	nDto.absAllNumRunes = []rune{}
	nDto.precision = 0

	nDto.SetNumericSeparatorsToUSADefault()

}

// FindIntArraySignificantDigitLimits - Receives an array of integers and converts them
// to a number string consisting of significant digits. Leading and trailing zeros are
// eliminated. See Method: FindNumStrSignificantDigitLimits()
func (nDto *NumStrDto) FindIntArraySignificantDigitLimits(intArray []int, precision uint, signVal int) (NumStrDto, error) {

	lenIntArray := len(intArray)

	var absNumStr []rune

	for i := 0; i < lenIntArray; i++ {
		absNumStr = append(absNumStr, rune(intArray[i]+48))
	}

	return nDto.FindNumStrSignificantDigitLimits(absNumStr, precision, signVal)
}

// FindNumStrSignificantDigitLimits
//
// Analyzes an array of characters which constitute a
// number string are returns the significant digits.
//
// Example:
// absAllRunes	precision	signVal		Result
// 001236700		4		 1			 123.67
// 000006700		4		 1			   0.67
// 001230000		4		 1			 123.0
func (nDto *NumStrDto) FindNumStrSignificantDigitLimits(absAllRunes []rune, precision uint, signVal int) (NumStrDto, error) {
	iPrecision := int(precision)
	firstIntIdx := -1
	lastIntIdx := -1
	lastFracIdx := -1

	isFractional := false

	if iPrecision > 0 {
		isFractional = true
	}

	lenAbsAllRunes := len(absAllRunes)
	lenAbsFracRunes := iPrecision
	lenAbsIntRunes := lenAbsAllRunes - lenAbsFracRunes

	if lenAbsAllRunes < 1 {
		return NumStrDto{}, errors.New("FindSignificantDigitLimits() - Error: absAllRunes has ZERO length!")
	}

	for i := 0; i < lenAbsAllRunes; i++ {

		if i < lenAbsIntRunes {

			if firstIntIdx == -1 && absAllRunes[i] > '0' && absAllRunes[i] <= '9' {
				firstIntIdx = i
			}

			lastIntIdx = i
		}

		if isFractional && i >= lenAbsIntRunes && absAllRunes[i] > '0' && absAllRunes[i] <= '9' {
			lastFracIdx = i
		}

	}

	if firstIntIdx == -1 {
		firstIntIdx = lastIntIdx
	}

	if isFractional && lastFracIdx == -1 {
		lastFracIdx = lenAbsIntRunes
	}

	numStrOut := ""

	if signVal < 0 {
		numStrOut = "-"
	}

	numStrOut += string(absAllRunes[firstIntIdx : lastIntIdx+1])
	if isFractional {
		numStrOut += string(nDto.decimalSeparator)
		numStrOut += string(absAllRunes[lastIntIdx+1 : lastFracIdx+1])
	}

	nOutDto, err := nDto.ParseNumStr(numStrOut)

	if err != nil {
		return NumStrDto{}, fmt.Errorf("FindSignificantDigitLimits() - Error retuned from nDto.ParseNumStr(numStrOut). numStrOut= '%v' Error= %v", numStrOut, err)
	}

	return nOutDto, nil
}

// FormatForMathOps - receives two NumStrDto objects and converts their strings
// such that both have the same number of integer and fractional digits. This will
// facilitate the performance of string based math operations such as addition and
// subtraction.
//
// The return values represent the formatted NumStrDto objects. The first NumStrDto
// returned always contains the larger absolute value. The second NumStrDto always
// contains the absolute numeric value which is less than or equal to the first
// NumStrDto object returned.
//
// The third parameter returned by this method is an int which will always be set to
// 1 or 0. 1 indicates that the absolute value of the first NumStrDto returned by
// this method is greater than the second NumStrDto returned by this method. If
// the int value returned is zero, it signals that the absolute values
// (not the signed values) of both returned NumStrDto objects are equal.
func (nDto *NumStrDto) FormatForMathOps(n1Dto, n2Dto NumStrDto) (n1DtoOut NumStrDto, n2DtoOut NumStrDto, compare int, isOrderReversed bool, err error) {

	ePrefix := "NumStrDto.FormatForMathOps() "

	lenN1AllRunes := 0
	lenN1IntRunes := 0
	lenN1FracRunes := 0
	lenN2AllRunes := 0
	lenN2IntRunes := 0
	lenN2FracRunes := 0

	err = n1Dto.IsValid(ePrefix + "n1Dto - ")

	if err != nil {
		return NumStrDto{}, NumStrDto{}, 0, false, err
	}

	err = n2Dto.IsValid(ePrefix + "n2Dto - ")

	if err != nil {
		return NumStrDto{}, NumStrDto{}, 0, false, err
	}

	compare = nDto.CompareAbsoluteValues(&n1Dto, &n2Dto)

	if compare == 1 {
		n1DtoOut = n1Dto.CopyOut()
		n2DtoOut = n2Dto.CopyOut()
	} else if compare == -1 {
		n1DtoOut = n2Dto.CopyOut()
		n2DtoOut = n1Dto.CopyOut()
		isOrderReversed = true
		compare = 1
	} else {
		// compare must be zero
		n1DtoOut = n1Dto.CopyOut()
		n2DtoOut = n2Dto.CopyOut()
	}

	n1DtoOutAbsIntRunes := n1DtoOut.GetAbsIntRunes()
	n1DtoOutAbsFracRunes := n1DtoOut.GetAbsFracRunes()

	n2DtoOutAbsIntRunes := n2DtoOut.GetAbsIntRunes()
	n2DtoOutAbsFracRunes := n2DtoOut.GetAbsFracRunes()

	if n1DtoOut.precision > n2DtoOut.precision {

		deltaPrecision := n1DtoOut.precision - n2DtoOut.precision

		for i := uint(0); i < deltaPrecision; i++ {
			n2DtoOut.absAllNumRunes = append(n2DtoOut.absAllNumRunes, '0')
			n2DtoOutAbsFracRunes = append(n2DtoOutAbsFracRunes, '0')
		}

		lenN2AllRunes = len(n2DtoOut.absAllNumRunes)
		lenN2IntRunes = len(n2DtoOutAbsIntRunes)
		lenN2FracRunes = len(n2DtoOutAbsFracRunes)

		n2DtoOut.precision = n1DtoOut.precision
		err = n2DtoOut.IsValid(ePrefix)

		if err != nil {
			return NumStrDto{}, NumStrDto{}, 0, false, err
		}

		lenN1AllRunes = len(n1DtoOut.absAllNumRunes)
		lenN1IntRunes = len(n1DtoOutAbsIntRunes)
		lenN1FracRunes = len(n1DtoOutAbsFracRunes)

	} else if n1DtoOut.precision < n2DtoOut.precision {

		deltaPrecision := n2DtoOut.precision - n1DtoOut.precision

		for i := uint(0); i < deltaPrecision; i++ {
			n1DtoOut.absAllNumRunes = append(n1DtoOut.absAllNumRunes, '0')
			n1DtoOutAbsFracRunes = append(n1DtoOutAbsFracRunes, '0')
		}

		lenN1AllRunes = len(n1DtoOut.absAllNumRunes)
		lenN1IntRunes = len(n1DtoOutAbsIntRunes)
		lenN1FracRunes = len(n1DtoOutAbsFracRunes)

		n1DtoOut.precision = n2DtoOut.precision
		err = n1DtoOut.IsValid(ePrefix)

		if err != nil {
			return NumStrDto{}, NumStrDto{}, 0, false, err
		}

		lenN2AllRunes = len(n2DtoOut.absAllNumRunes)
		lenN2IntRunes = len(n2DtoOutAbsIntRunes)
		lenN2FracRunes = len(n2DtoOutAbsFracRunes)

	} else {
		// n1DtoOut.precision == n2DtoOut.precision

		lenN1AllRunes = len(n1DtoOut.absAllNumRunes)
		lenN1IntRunes = len(n1DtoOutAbsIntRunes)
		lenN1FracRunes = len(n1DtoOutAbsFracRunes)

		lenN2AllRunes = len(n2DtoOut.absAllNumRunes)
		lenN2IntRunes = len(n2DtoOutAbsIntRunes)
		lenN2FracRunes = len(n2DtoOutAbsFracRunes)

	}

	if lenN2IntRunes > lenN1IntRunes {

		var absAllRunes []rune
		var absIntRunes []rune
		deltaRunes := lenN2IntRunes - lenN1IntRunes
		for i := 0; i < deltaRunes; i++ {
			absAllRunes = append(absAllRunes, '0')
			absIntRunes = append(absIntRunes, '0')
		}

		for j := 0; j < lenN1AllRunes; j++ {
			absAllRunes = append(absAllRunes, n1DtoOut.absAllNumRunes[j])

			if j < lenN1IntRunes {
				absIntRunes = append(absIntRunes, n1DtoOutAbsIntRunes[j])
			}

		}

		n1DtoOut.absAllNumRunes = absAllRunes
		n1DtoOutAbsIntRunes = absIntRunes
		lenN1AllRunes = len(n1DtoOut.absAllNumRunes)
		lenN1IntRunes = len(n1DtoOutAbsIntRunes)

		err = n1DtoOut.IsValid(ePrefix)

		if err != nil {
			return NumStrDto{}, NumStrDto{}, 0, false, err
		}

	} else if lenN1IntRunes > lenN2IntRunes {

		var absAllRunes []rune
		var absIntRunes []rune
		deltaRunes := lenN1IntRunes - lenN2IntRunes
		for i := 0; i < deltaRunes; i++ {
			absAllRunes = append(absAllRunes, '0')
			absIntRunes = append(absIntRunes, '0')
		}

		for j := 0; j < lenN2AllRunes; j++ {
			absAllRunes = append(absAllRunes, n2DtoOut.absAllNumRunes[j])

			if j < lenN2IntRunes {
				absIntRunes = append(absIntRunes, n2DtoOutAbsIntRunes[j])
			}

		}

		n2DtoOut.absAllNumRunes = absAllRunes
		n2DtoOutAbsIntRunes = absIntRunes
		lenN2AllRunes = len(n2DtoOut.absAllNumRunes)
		lenN2IntRunes = len(n2DtoOutAbsIntRunes)

		err := n2DtoOut.IsValid(ePrefix)

		if err != nil {
			return NumStrDto{}, NumStrDto{}, 0, false, err
		}

	}

	if lenN1AllRunes != lenN2AllRunes {
		return NumStrDto{}, NumStrDto{}, 0, false,
			fmt.Errorf("FormatForMathOps() - n1 and n2 AllNumRune arrays are NOT equal in length. "+
				"n1 length= '%v' n2 length= '%v'", lenN1AllRunes, lenN2AllRunes)
	}

	if lenN1IntRunes != lenN2IntRunes {
		return NumStrDto{}, NumStrDto{}, 0, false,
			fmt.Errorf("FormatForMathOps() - n1 and n2 IntRunes arrays are NOT equal in length. "+
				"n1 length= '%v' n2 length= '%v'", lenN1IntRunes, lenN2IntRunes)
	}

	if lenN1FracRunes != lenN2FracRunes {
		return NumStrDto{}, NumStrDto{}, 0, false,
			fmt.Errorf("FormatForMathOps() - n1 and n2 FracRunes arrays are NOT equal in length. "+
				"n1 length= '%v' n2 length= '%v'", lenN1FracRunes, lenN2FracRunes)
	}

	if n1DtoOut.precision != n2DtoOut.precision {
		return NumStrDto{}, NumStrDto{}, 0, false,
			fmt.Errorf("FormatForMathOps() - n1 and n2 FracRunes arrays are NOT equal in length. "+
				"n1 length= '%v' n2 length= '%v'", lenN1FracRunes, lenN2FracRunes)

	}

	err = n1DtoOut.IsValid(ePrefix + "n1DtoOut - ")

	if err != nil {
		return NumStrDto{}, NumStrDto{}, 0, false, err
	}

	err = n2DtoOut.IsValid(ePrefix + "n2DtoOut - ")

	if err != nil {
		return NumStrDto{}, NumStrDto{}, 0, false, err
	}

	return n1DtoOut, n2DtoOut, compare, isOrderReversed, nil
}

// FormatCurrencyStr - Formats the current NumStrDto numeric value as a currency string.
//
// If the Currency Symbol was not previously set for this NumStrDto, the currency symbol
// is defaulted to the USA standard dollar sign, ('$').
//
// If the Decimal Separator was not previously set for this NumStrDto, the Decimal Separator
// is defaulted to the USA standard period ('.').
//
// If the Thousands Separator was not previously set for this NumStrDto, the Thousands
// Separator is defaulted to the USA standard comma (',').
//
// Input Parameters
// ================
//
// negValMode NegativeValueFmtMode -	Specifies the display mode for negative values:
//
//	LEADMINUSNEGVALFMTMODE 		-	Negative values formatted with
//													 		a leading minus sign.
//															Example: -$123,456.78
//
//	PARENTHESESNEGVALFMTMODE	-	Negative values formatted with
//															surrounding parentheses.
//															Example: ($123,456.78)
func (nDto *NumStrDto) FormatCurrencyStr(negValMode NegativeValueFmtMode) (string, error) {

	ePrefix := "NumStrDto.FormatCurrencyStr() "

	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	err := nDto.IsValid("")

	if err != nil {
		return "",
			fmt.Errorf(ePrefix + "")
	}

	lenAllNumRunes := len(nDto.absAllNumRunes)

	lenOut := lenAllNumRunes

	lenIntRunes := lenAllNumRunes - int(nDto.precision)

	seps := lenIntRunes / 3

	mod := lenIntRunes - (seps * 3)

	if mod == 0 {
		seps--
	}

	// adjust for thousands delimiters
	lenOut += seps

	// adjust for negative sign value
	if nDto.signVal == -1 {
		if negValMode == LEADMINUSNEGVALFMTMODE {
			lenOut++
		} else {
			// MUST BE negValMode == PARENTHESESNEGVALFMTMODE
			lenOut += 2
		}
	}

	// adjust for decimal point
	if nDto.precision > 0 {
		lenOut++
	}

	// adjust for currency symbol
	lenOut++

	outRunes := make([]rune, lenOut)
	outIdx := lenOut - 1
	allNumsIdx := lenAllNumRunes - 1

	// If negative value and parenthesis formatting
	// specified, format trailing parenthesis.
	if nDto.signVal == -1 &&
		negValMode == PARENTHESESNEGVALFMTMODE {
		outRunes[outIdx] = ')'
		outIdx--
	}

	if nDto.precision > 0 {

		for i := 0; i < int(nDto.precision); i++ {
			outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
			outIdx--
			allNumsIdx--
		}

		outRunes[outIdx] = nDto.decimalSeparator
		outIdx--
	}

	sepCnt := 0

	for i := 0; i < lenIntRunes; i++ {

		sepCnt++

		if sepCnt == 4 && seps > 0 {
			sepCnt = 1
			seps--
			outRunes[outIdx] = nDto.thousandsSeparator
			outIdx--
		}

		outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
		outIdx--
		allNumsIdx--

	}

	outRunes[outIdx] = nDto.currencySymbol

	// If required, add leading negative
	// value sign
	if nDto.signVal == -1 &&
		negValMode == PARENTHESESNEGVALFMTMODE {
		outRunes[0] = '('
	} else if nDto.signVal == -1 {
		outRunes[0] = '-'
	}

	return string(outRunes), nil

}

// FormatNumStr - Formats the numeric value of the current NumStrDto
// as number string consisting of integer digits to the left of the
// decimal point and fractional digits to the right of the decimal
// point, if such fractional digits exist. The resulting number string
// will NOT contain a currency symbol or thousands separators.
//
// Example: 123456.789
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
func (nDto *NumStrDto) FormatNumStr(negValMode NegativeValueFmtMode) (string, error) {

	ePrefix := "NumStrDto.FormatNumStr() "

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	err := nDto.IsValid("")

	if err != nil {
		return "", fmt.Errorf(ePrefix+"NumStrDto INVALID! Error='%v'",
			err.Error())
	}

	lenAllNumRunes := len(nDto.absAllNumRunes)

	lenOut := lenAllNumRunes

	lenIntRunes := lenAllNumRunes - int(nDto.precision)

	// adjust for negative sign value
	if nDto.signVal == -1 {
		if negValMode == LEADMINUSNEGVALFMTMODE {
			lenOut++
		} else {
			// MUST BE negValMode == PARENTHESESNEGVALFMTMODE
			lenOut += 2
		}

	}

	// adjust for decimal point
	if nDto.precision > 0 {
		lenOut++
	}

	outRunes := make([]rune, lenOut)
	outIdx := lenOut - 1

	if nDto.signVal == -1 &&
		negValMode == PARENTHESESNEGVALFMTMODE {
		outRunes[outIdx] = ')'
		outIdx--
	}

	allNumsIdx := lenAllNumRunes - 1

	if nDto.precision > 0 {

		for i := 0; i < int(nDto.precision); i++ {
			outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
			outIdx--
			allNumsIdx--
		}

		outRunes[outIdx] = nDto.decimalSeparator
		outIdx--
	}

	for i := 0; i < lenIntRunes; i++ {

		outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
		outIdx--
		allNumsIdx--

	}

	if nDto.signVal == -1 {
		if negValMode == LEADMINUSNEGVALFMTMODE {
			outRunes[0] = '-'
		} else {
			// MUST BE negValMode == PARENTHESESNEGVALFMTMODE
			outRunes[0] = '('
		}

	}

	return string(outRunes), nil
}

// FormatThousandsStr - Returns the number string delimited with the
// nDto.thousandsSeparator character plus the Decimal Separator if
// applicable.
//
// Example:
// numstr = 1000000.234 converted to 1,000,000.234
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
func (nDto *NumStrDto) FormatThousandsStr(negValMode NegativeValueFmtMode) (string, error) {

	ePrefix := "NumStrDto.FormatThousandsStr() "

	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	err := nDto.IsValid("")

	if err != nil {
		return "", fmt.Errorf(ePrefix+"NumStrDto INVALID! Error='%v'",
			err.Error())
	}

	lenAllNumRunes := len(nDto.absAllNumRunes)

	lenOut := lenAllNumRunes

	lenIntRunes := lenAllNumRunes - int(nDto.precision)

	seps := lenIntRunes / 3

	mod := lenIntRunes - (seps * 3)

	if mod == 0 {
		seps--
	}

	// adjust for thousands delimiters
	lenOut += seps

	// adjust for negative sign value
	if nDto.signVal == -1 {
		if negValMode == LEADMINUSNEGVALFMTMODE {
			lenOut++
		} else {
			// MUST BE negValMode == PARENTHESESNEGVALFMTMODE
			lenOut += 2
		}

	}

	// adjust for decimal point
	if nDto.precision > 0 {
		lenOut++
	}

	outRunes := make([]rune, lenOut)
	outIdx := lenOut - 1

	if nDto.signVal == -1 &&
		negValMode == PARENTHESESNEGVALFMTMODE {
		outRunes[outIdx] = ')'
		outIdx--
	}

	allNumsIdx := lenAllNumRunes - 1

	if nDto.precision > 0 {

		for i := 0; i < int(nDto.precision); i++ {
			outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
			outIdx--
			allNumsIdx--
		}

		outRunes[outIdx] = nDto.decimalSeparator
		outIdx--
	}

	sepCnt := 0

	for i := 0; i < lenIntRunes; i++ {

		sepCnt++

		if sepCnt == 4 && seps > 0 {
			sepCnt = 1
			seps--
			outRunes[outIdx] = nDto.thousandsSeparator
			outIdx--
		}

		outRunes[outIdx] = nDto.absAllNumRunes[allNumsIdx]
		outIdx--
		allNumsIdx--

	}

	if nDto.signVal == -1 {
		if negValMode == LEADMINUSNEGVALFMTMODE {
			outRunes[0] = '-'
		} else {
			outRunes[0] = '('
		}

	}

	return string(outRunes), nil
}

// GetAbsoluteBigInt
//
// Returns the absolute value of all numeric digits in
// the number string (nDto.absAllNumRunes). As such,
// Fractional digits to the right of the decimal are
// included in the consolidate integer number. All the
// numeric digits in the number string are therefore
// returned as a *big.Int.
//
// This method will fail if the NumStrDto has not been
// properly initialized with a valid number string.
func (nDto *NumStrDto) GetAbsoluteBigInt() (*big.Int, error) {

	ePrefix := "NumStrDto.GetAbsoluteBigInt() "

	err := nDto.IsValid("")

	if err != nil {
		return big.NewInt(0),
			fmt.Errorf(ePrefix+"This NumStrDto instance is INVALID! Error='%v'", err.Error())
	}

	lenAllNumRunes := len(nDto.absAllNumRunes)

	if lenAllNumRunes == 0 {
		s := ePrefix +
			"- The existing NumStrDto is a zero length number. " +
			"Re-initialize the NumStrDto object and try again."
		return big.NewInt(0), errors.New(s)

	}

	base10 := big.NewInt(int64(10))
	absBigInt := big.NewInt(0)

	for i := 0; i < lenAllNumRunes; i++ {

		absBigInt = big.NewInt(0).Mul(absBigInt, base10)
		absBigInt = big.NewInt(0).Add(absBigInt,
			big.NewInt(int64(nDto.absAllNumRunes[i]-48)))

	}

	return absBigInt, nil
}

// GetAbsAllNumRunes
//
// Returns an array of runes representing all the integer
// and fractional digits included in the current
// NumStrDto instance. The rune array returned will
// consist of numeric digits with no sign value prefixed.
//
// This effectively returns the absolute value of all
// integer and fractional digits combined in one rune
// array (there is no decimal point).
func (nDto *NumStrDto) GetAbsAllNumRunes() []rune {

	lenAbsAllNumRunes := len(nDto.absAllNumRunes)

	if lenAbsAllNumRunes == 0 {
		return []rune{}
	}

	outRunes := make([]rune, lenAbsAllNumRunes, lenAbsAllNumRunes+50)

	for i := 0; i < lenAbsAllNumRunes; i++ {
		outRunes[i] = nDto.absAllNumRunes[i]
	}

	return outRunes
}

// GetAbsFracRunes
//
// Returns all the fractional digits to the right of the
// decimal place in the current NumStrDto instance as an
// array of runes. The rune array is not signed; that is,
// the rune array does not contain a '+' or '-' character
// in the first array position. The rune array is
// therefore said to represent the absolute value of the
// fractional digits in the current NumStrDto numeric
// value.
func (nDto *NumStrDto) GetAbsFracRunes() []rune {

	precision := int(nDto.precision)

	lenAllNums := len(nDto.absAllNumRunes)

	if lenAllNums == 0 ||
		precision < 0 ||
		precision > lenAllNums {
		return []rune{}
	}

	absFracRunes := make([]rune, precision, precision+50)

	lenIntNums := lenAllNums - int(nDto.precision)

	for i := lenIntNums; i < lenAllNums; i++ {
		absFracRunes[i-lenIntNums] = nDto.absAllNumRunes[i]
	}

	return absFracRunes
}

// GetAbsFracRunesLength - Returns the length of the
// fractional digits in the number string.
func (nDto *NumStrDto) GetAbsFracRunesLength() int {
	return int(nDto.precision)
}

// GetAbsIntRunes
//
// Returns all the integer digits included in the current
// NumStrDto numeric value as an array of runes.
//
// The returned rune array does not contain a sign value
// in the first position and therefore represents the
// absolute or positive value of all the integer digits.
//
// The integer digits of a NumStrDto numeric includes all
// the digits to the left of the decimal point.
//
// If the current NumStrDto consists of zero integers and
// fractional digits (Example: '0.123456'), this method
// will return a rune array consisting one array element
// with a '0' value.
func (nDto *NumStrDto) GetAbsIntRunes() []rune {

	lenAllNum := len(nDto.absAllNumRunes)

	precision := int(nDto.precision)

	if lenAllNum == 0 ||
		precision < 0 ||
		precision >= lenAllNum {
		return []rune{}
	}

	lenIntNum := lenAllNum - precision

	absIntRunes := make([]rune, lenIntNum, lenIntNum+50)

	for i := 0; i < lenIntNum; i++ {
		absIntRunes[i] = nDto.absAllNumRunes[i]
	}

	return absIntRunes
}

// GetAbsNumStr - Returns all digits in the current NumStrDto numeric
// value as a pure, unsigned number string. If fractional digits exists
// they are included and NOT separated by a decimal separator.
//
// Examples:
// numstr				result
// ------       ------
// 123.45				12345
// -123.45			12345
func (nDto *NumStrDto) GetAbsNumStr() string {
	return string(nDto.absAllNumRunes)
}

// GetAbsIntRunesLength - Returns the length of the
// integer portion of the number string.
func (nDto *NumStrDto) GetAbsIntRunesLength() int {

	lenAllNums := len(nDto.absAllNumRunes)

	return lenAllNums - int(nDto.precision)
}

// GetBigInt
//
// Returns an integer of type *big.Int representing the
// signed integer value of NumStrDto.numStrDto. Decimal
// numbers like '-123.456' will be returned as signed
// integer values, '-123456'.
//
// This method will fail if the NumStrDto has not been
// properly initialized with a valid number string.
func (nDto *NumStrDto) GetBigInt() (*big.Int, error) {

	ePrefix := "NumStrDto.GetBigInt()"

	err := nDto.IsValid("")

	if err != nil {
		return big.NewInt(0),
			fmt.Errorf(ePrefix+"NumStrDto is INVALID! Error='%v' ", err.Error())
	}

	absBigInt, err := nDto.GetAbsoluteBigInt()

	if err != nil {
		s := fmt.Sprintf("GetBigInt() - Error returned from nDto.GetAbsoluteBigInt(). Error= %v", err)
		return big.NewInt(0), errors.New(s)
	}

	if nDto.signVal < 0 {
		return big.NewInt(0).Neg(absBigInt), nil
	}

	return big.NewInt(0).Set(absBigInt), nil
}

// GetBigIntNum - Converts the numeric value of the
// current NumStrDto to a 'BigIntNum' type and returns
// it to the calling function.
//
// The returned BigIntNum will contain numeric separators
// (decimal separator, thousands separator and currency
// symbol) copied from the current NumStrDto instance.
//
// Before returning the BigIntNum result, this method
// performs a validity test on the current NumStrDto instance.
func (nDto *NumStrDto) GetBigIntNum() (BigIntNum, error) {
	ePrefix := "NumStrDto.GetBigIntNum() "

	err := nDto.IsValid(ePrefix + " This NumStrDto INVALID! ")

	if err != nil {
		return new(BigIntNum).NewZero(0), err
	}

	numSeps := nDto.GetNumericSeparatorsDto()

	bInt, err := nDto.GetBigInt()

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by nDto.GetBigInt() "+
				"Error='%v' ", err.Error())
	}

	bIntNum := new(BigIntNum).NewBigInt(bInt, nDto.precision)

	err = bIntNum.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return new(BigIntNum).NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by bIntNum.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' \n", err.Error())
	}

	return bIntNum, nil
}

// GetCurrencySymbol - Returns the character currently designated
// as the currency symbol for this number string.
//
// In the USA, the currency symbol is the dollar sign ('$').
//
// For a list of Major Currency Unicode Symbols, see constants
// located in: MikeAustin71/mathopsgo/mathops/mathopsconstants.go
//
// Example: $123.45
func (nDto *NumStrDto) GetCurrencySymbol() rune {

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	return nDto.currencySymbol

}

// GetCurrencyParen - Returns the number string delimited with the
// nDto.thousandsSeparator character and the currency symbol. If
// the value is negative, the number will be surrounded in parentheses.
//
// Example:
// numstr = 1000000.23
// GetThouStr() = $1,000,000.23
//
// numstr = -1000000.23
// GetThouStr() = ($1,000,000.23)
//
// Note: If the current NumStrDto is invalid, this method
// returns an empty string.
func (nDto *NumStrDto) GetCurrencyParen() string {

	outStr, err := nDto.FormatCurrencyStr(PARENTHESESNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr

}

// GetCurrencyStr - Returns the number string delimited with the
// nDto.thousandsSeparator character and the currency symbol.
// If the value is negative, a leading minus sign will be prefixed
// to the currency display.
//
// Example:
// numstr = 1000000.23
// GetCurrencyStr() = $1,000,000.23
//
// numstr = -1000000.23
// GetCurrencyStr() = -$1,000,000.23
//
// Note: If the current NumStrDto is invalid, this method
// returns an empty string.
func (nDto *NumStrDto) GetCurrencyStr() string {

	outStr, err := nDto.FormatCurrencyStr(LEADMINUSNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr
}

// GetDecimalSeparator - returns the character designated
// as the decimal separator for the current NumStrDto instance.
//
// In the USA, the decimal separator is the period character ('.').
//
// Example:		123.456
func (nDto *NumStrDto) GetDecimalSeparator() rune {

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	return nDto.decimalSeparator

}

// GetIntAry
//
// Converts the current NumStrDto instance to a Type IntAry
// and returns it to the calling function.
func (nDto *NumStrDto) GetIntAry() (IntAry, error) {
	ePrefix := "NumStrDto.GetIntAryElements() "

	err := nDto.IsValid(ePrefix)

	if err != nil {
		return IntAry{}, err
	}

	numSeps := nDto.GetNumericSeparatorsDto()

	ia, err := IntAry{}.NewNumStrWithNumSeps(nDto.GetNumStr(), numSeps)

	if err != nil {
		return IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by IntAry{}.NewNumStrWithNumSeps(nDto.GetNumStr(), numSeps). "+
				"nDto='%v' Error='%v'", nDto.GetNumStr(), err.Error())
	}

	return ia, nil
}

// GetNumericSeparatorsDto - Returns a structure containing the
// character or rune values for decimal point separator, thousands
// separator and currency symbol.
func (nDto *NumStrDto) GetNumericSeparatorsDto() NumericSeparatorDto {

	numSeps := NumericSeparatorDto{}
	numSeps.DecimalSeparator = nDto.GetDecimalSeparator()
	numSeps.ThousandsSeparator = nDto.GetThousandsSeparator()
	numSeps.CurrencySymbol = nDto.GetCurrencySymbol()

	return numSeps
}

// GetNumParen - Returns the numeric value of the current NumStrDto
// instance as a signed number string. The resulting number string
// will NOT contain a currency symbol or thousands separators. It
// will contain a decimal separator and fractional digits if such
// fractional digits exist.
//
// Note: If the current NumStrDto is invalid, this method will return
// an empty string.
//
// If the sign of the numeric value is negative, the resulting number
// string will be surrounded in parentheses.
//
// Examples:
// numeric value						result
//
//	 123456.78							123456.78
//	-123456.78             (123456.78)
func (nDto *NumStrDto) GetNumParen() string {

	outStr, err := nDto.FormatNumStr(PARENTHESESNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr

}

// GetNumStr - returns the numeric value of the current NumStrDto
// instance as a signed number string. The resulting number string
// will NOT contain a currency symbol or thousands separators. It
// will contain a decimal separator and fractional digits if such
// fractional digits exist.
//
// Note: If the current NumStrDto is invalid, this method will return
// an empty string.
//
// Examples:
//
//	 123456.78
//	-123456.78
func (nDto *NumStrDto) GetNumStr() string {

	outStr, err := nDto.FormatNumStr(LEADMINUSNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr
}

// GetNumStrDto - Returns a deep copy of the current NumStrDto
// instance.
//
// The returned NumStrDto instance will contain numeric
// separators (decimal separator, thousands separator
// and currency symbol) copied from the current NumStrDto
// instance.
//
// Before returning the NumStrDto result, this method
// performs a validity test on the current NumStrDto instance.
//
// This method is necessary in order to fulfill the requirements
// of the INumMgr interface.
func (nDto *NumStrDto) GetNumStrDto() (NumStrDto, error) {

	ePrefix := "NumStrDto.GetNumStrDto() "

	err := nDto.IsValid(ePrefix + "NumStrDto INVALID! ")

	if err != nil {
		return NumStrDto{}.New(), err
	}

	return nDto.CopyOut(), nil
}

// GetPrecision - Returns the precision of the current
// NumStrDto Instance.
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
func (nDto *NumStrDto) GetPrecision() int {
	return int(nDto.precision)
}

// GetPrecisionUint - Returns the precision of the
// current NumStrDto Instance as an unsigned integer
// (uint). precision represents the number of fractional
// digits to the right of the decimal point.
//
// precision is defined as the number of numeric digits to
// the right of the decimal place. To compute the location
// of the decimal point in a string of numeric digits, go
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
func (nDto *NumStrDto) GetPrecisionUint() uint {
	return nDto.precision
}

// GetRationalNumber - returns the sign value of the number string, plus the
// numeric value of the number string expressed as a Rational Number.
//
// This method will return an error if the NumStrDto fields are not properly
// initialized and populated.
//
// Returns
// =======
//
// sign value  						int 			- sign value of the number string
// big Rational Number		*big.Rat	- Number string expressed as a
//
//	rational number
//
// err										error			- In case of failure, an 'error' type
//
//	is returned. In case of success this
//	value is 'nil'
func (nDto *NumStrDto) GetRationalNumber() (int, *big.Rat, error) {

	ePrefix := "NumStrDto.GetRationalNumber() "

	err := nDto.IsValid("")

	if err != nil {
		return 0, big.NewRat(1, 1),
			fmt.Errorf(ePrefix+"This NumStrDto instance is INVALID! Error='%v'", err.Error())
	}

	ratZero := big.NewRat(0, 1)

	err = nDto.IsValid(ePrefix)

	if err != nil {
		return 0, ratZero, errors.New(ePrefix +
			"- Error: The existing NumStrDto is corrupted or improperly initialized. " +
			"Re-initialize the NumStrDto object and try again.")
	}

	signVal := nDto.signVal

	absInt, isOk := big.NewInt(0).SetString(string(nDto.absAllNumRunes), 10)

	if !isOk {
		return 0, ratZero, fmt.Errorf(ePrefix+
			"- Conversion of nDto.absAllNumRunes to big.Int Failed! "+
			"nDto.absIntRunes= '%v'", nDto.absAllNumRunes)
	}

	base10 := big.NewInt(10)

	bigPrecision := big.NewInt(int64(nDto.precision))

	scaleFactor := big.NewInt(0).Exp(base10, bigPrecision, nil)

	rationalNum := big.NewRat(1, 1).SetFrac(absInt, scaleFactor)

	return signVal, rationalNum, nil
}

// GetScaleFactor - returns the scale factor for this number
// string. Scale factor is defined by 10 raised to the power
// of nDto.precision.  nDto.precision is the number of
// digits to the right of the decimal point.
//
// This method will fail if the NumStrDto has not been properly
// initialized with a valid number string.
func (nDto *NumStrDto) GetScaleFactor() (*big.Int, error) {

	ePrefix := "NumStrDto.GetScaleFactor() "

	err := nDto.IsValid("")

	if err != nil {
		return big.NewInt(0),
			fmt.Errorf(ePrefix+"This NumStrDto instance is INVALID! Error='%v'", err.Error())
	}

	if len(nDto.absAllNumRunes) == 0 {
		s := ePrefix +
			"- The existing NumStrDto is a zero length number. " +
			"Re-initialize the NumStrDto object and try again."
		return big.NewInt(0), errors.New(s)

	}

	if nDto.precision == 0 {
		return big.NewInt(int64(1)), nil
	}

	base10 := big.NewInt(0).SetInt64(int64(10))

	bigPrecision := big.NewInt(0).SetInt64(int64(nDto.precision))

	scaleFactor := big.NewInt(0).Exp(base10, bigPrecision, nil)

	return scaleFactor, nil

}

// GetSciNotationNumber - Converts the numeric value of the current
// NumStrDto instance into scientific notation and returns this value
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
func (nDto *NumStrDto) GetSciNotationNumber(mantissaLen uint) (SciNotationNum, error) {

	ePrefix := "NumStrDto.GetSciNotationNumber() "

	bINum, err := nDto.GetBigIntNum()

	if err != nil {
		return SciNotationNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by nDto.GetBigIntNum(). Error='%v'",
				err.Error())
	}

	sciNotation, err := bINum.GetSciNotationNumber(mantissaLen)

	if err != nil {
		return SciNotationNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by bINum.GetSciNotationNumber(mantissaLen). Error='%v'",
				err.Error())
	}

	return sciNotation, nil
}

// GetSciNotationStr - Returns a string expressing the current NumStrDto
// numerical value as scientific notation.
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
func (nDto *NumStrDto) GetSciNotationStr(mantissaLen uint) (string, error) {

	ePrefix := "NumStrDto.GetSciNotationStr() "

	if mantissaLen < 2 {
		mantissaLen = 2
	}

	sciNotation, err := nDto.GetSciNotationNumber(mantissaLen)

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

// GetSign - Returns the sign value for this NumStrDto
// numeric value. Return values will be either +1 or -1.
func (nDto *NumStrDto) GetSign() int {
	return nDto.signVal
}

// GetThisPointer - Returns a pointer to the current NumStrDto instance.
func (nDto *NumStrDto) GetThisPointer() *NumStrDto {

	return nDto
}

// GetThouParen - Returns the number string delimited with the
// nDto.thousandsSeparator character. Negative values are
// surrounded in parentheses.
//
// Example:
// numstr = 1000000.234
// GetThouStr() = 1,000,000.234
//
// numstr = -1000000.234
// GetThouStr() = (1,000,000.234)
//
// Note: If the current NumStrDto is invalid, this method
// returns an empty string.
func (nDto *NumStrDto) GetThouParen() string {

	outStr, err := nDto.FormatThousandsStr(PARENTHESESNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr
}

// GetThouStr - Returns the number string delimited with the
// nDto.thousandsSeparator character plus the Decimal Separator
// if applicable.
//
// Example:
// numstr = 1000000.234
// GetThouStr() = 1,000,000.234
//
// numstr = -1000000.234
// GetThouStr() = -1,000,000.234
//
// Note: If the current NumStrDto is invalid, this method
// returns an empty string.
func (nDto *NumStrDto) GetThouStr() string {

	outStr, err := nDto.FormatThousandsStr(LEADMINUSNEGVALFMTMODE)

	if err != nil {
		return ""
	}

	return outStr
}

// GetThousandsSeparator - returns a rune which represents
// the character currently used to separate thousands in
// the display of the current NumStrDto number string.
//
// In the USA, the thousands separator is a comma character.
//
// Example: 1,000,000,000
func (nDto *NumStrDto) GetThousandsSeparator() rune {

	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	return nDto.thousandsSeparator
}

// GetZeroNumStrDto - returns a new NumStrDto initialized
// to zero value. If the parameter numFracDigits is set
// to a value greater than zero, then an equal number of
// zero characters will be added to the right of the
// decimal point.
// Examples:
// numFracDigits		Results NumStrOut
//
//	0									"0"
//	2									"0.00"
func (nDto *NumStrDto) GetZeroNumStrDto(numFracDigits uint) NumStrDto {

	// Set defaults for thousands separators,
	// decimal separators and currency Symbols
	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	n2Dto := NumStrDto{}.New()
	n2Dto.signVal = 1
	n2Dto.thousandsSeparator = nDto.thousandsSeparator
	n2Dto.decimalSeparator = nDto.decimalSeparator
	n2Dto.currencySymbol = nDto.currencySymbol
	n2Dto.signVal = 1
	n2Dto.precision = 0
	n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')

	if numFracDigits > 0 {

		for i := uint(0); i < numFracDigits; i++ {
			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')
		}

		n2Dto.precision = numFracDigits
	}

	return n2Dto
}

// HasNumericDigits - returns 'false' if the number
// string for the current NumStrDto instance is uninitialized
// and contains no numeric digits. In this case, the length
// of the number string is zero characters.
//
// If this method returns 'true' it signals that there is at
// least one numeric digit in the number string, even if that
// digit is zero.
func (nDto *NumStrDto) HasNumericDigits() bool {

	err := nDto.IsValid("NumStrDto.HasNumericDigits() ")

	if err != nil {
		return false
	}

	return true
}

// IsNumStrZeroValue - Returns 'true' if all the digits in the number
// string for the current NumStrDto instance are zero.
func (nDto *NumStrDto) IsNumStrZeroValue(numDto *NumStrDto) bool {

	lenAbsAllNumRunes := len(numDto.absAllNumRunes)

	isZeroVal := true

	for i := 0; i < lenAbsAllNumRunes; i++ {
		if numDto.absAllNumRunes[i] != '0' {
			isZeroVal = false
		}
	}

	return isZeroVal
}

// IsFractionalValue - Returns 'true' if the numeric value of the
// current NumStrDto object includes a fractional value; that is,
// the number has fractional digits to the right of the decimal
// point.
func (nDto *NumStrDto) IsFractionalValue() bool {

	if nDto.precision > 0 {
		return true
	}

	return false
}

// IsValid - Performs a diagnostic review of the current NumStrDto
// instance and returns 'nil' if the NumStrDto object is valid in all
// respects.
//
// If the NumStrDto instance is judged invalid, an error message is returned.
func (nDto *NumStrDto) IsValid(errName string) error {

	if errName == "" {
		errName = "NumStrDto.IsValid() "
	}

	// Set defaults for thousands separators,
	// decimal separators and currency Symbols
	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	lenAbsAllNumRunes := len(nDto.absAllNumRunes)

	if lenAbsAllNumRunes == 0 {
		return errors.New(errName +
			"- Error: Number string is a ZERO length string!")
	}

	if int(nDto.precision) >= lenAbsAllNumRunes {
		return errors.New(errName +
			"- Error: precision does match number string. Type is Corrupted!")
	}

	if nDto.signVal != 1 && nDto.signVal != -1 {
		return fmt.Errorf("%v - sign Value is INVALID. Should be +1 or -1. This sign Value is %v",
			errName, nDto.signVal)
	}

	for i := 0; i < lenAbsAllNumRunes; i++ {

		if nDto.absAllNumRunes[i] < '0' || nDto.absAllNumRunes[i] > '9' {

			return errors.New(errName + "- Error: Non-Numeric character found in number string!")

		}
	}

	return nil
}

// IsZero - Returns true if the value of the current NumStrDto
// instance is zero.
func (nDto *NumStrDto) IsZero() bool {

	lenRunes := len(nDto.absAllNumRunes)

	if lenRunes == 0 {
		return true
	}

	for i := 0; i < lenRunes; i++ {
		if nDto.absAllNumRunes[i] != '0' {
			return false
		}
	}

	return true
}

// Multiply - Multiplies the current NumStrDto by the input
// parameter NumStrDto and stores the result in the current
// NumStrDto.
func (nDto *NumStrDto) Multiply(n2Dto NumStrDto) error {
	ePrefix := "NumStrDto.Multiply() "

	n1Dto := nDto.CopyOut()

	nResult, err := nDto.MultiplyNumStrs(n1Dto, n2Dto)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by MultiplyNumStrs(n1Dto, n2Dto). "+
			"Error='%v'", err.Error())
	}

	nDto.CopyIn(nResult)

	return nil
}

// MultiplyNumStrs
//
// Multiplies two NumStrDto instances and returns the
// result as an separate NumStrDto instance.
func (nDto *NumStrDto) MultiplyNumStrs(n1Dto NumStrDto, n2Dto NumStrDto) (NumStrDto, error) {
	ePrefix := "NumStrDto.MultiplyNumStrs() "

	if err := n1Dto.IsValid(ePrefix + "- "); err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"- n1Dto, first NumStrDto is invalid! Error= %v", err)
	}

	if err := n2Dto.IsValid(ePrefix + "- "); err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"- n2Dto, second NumStrDto is invalid! Error= %v", err)
	}

	lenN1AbsAllRunes := len(n1Dto.absAllNumRunes)
	lenN2AbsAllRunes := len(n2Dto.absAllNumRunes)

	var n1Setup NumStrDto
	var n2Setup NumStrDto

	if lenN2AbsAllRunes > lenN1AbsAllRunes {
		n1Setup = n2Dto.CopyOut()
		n2Setup = n1Dto.CopyOut()
	} else if lenN1AbsAllRunes > lenN2AbsAllRunes {
		n1Setup = n1Dto.CopyOut()
		n2Setup = n2Dto.CopyOut()
	} else {
		// Must be lenN1AbsAllRunes == lenN2AbsAllRunes
		n1Setup = n1Dto.CopyOut()
		n2Setup = n2Dto.CopyOut()

	}

	newPrecision := n1Setup.precision + n2Setup.precision
	newSignVal := 1

	if n1Setup.signVal == n2Setup.signVal {
		newSignVal = 1
	} else {
		// Must be n1Setup.signVal != n2Setup.signVal
		newSignVal = -1
	}

	lenN1AbsAllRunes = len(n1Setup.absAllNumRunes)
	lenN2AbsAllRunes = len(n2Setup.absAllNumRunes)
	lenLevels := lenN2AbsAllRunes
	lenNumPlaces := (lenN1AbsAllRunes + lenN2AbsAllRunes) + 1

	intMAry := make([][]int, lenLevels)

	for i := 0; i < lenLevels; i++ {
		intMAry[i] = make([]int, lenNumPlaces)
	}

	intFinalAry := make([]int, lenNumPlaces+1)

	carry := 0
	levels := 0
	place := 0
	n1 := 0
	n2 := 0
	n3 := 0
	n4 := 0
	for i := lenN2AbsAllRunes - 1; i >= 0; i-- {

		place = (lenNumPlaces - 1) - levels

		for j := lenN1AbsAllRunes - 1; j >= 0; j-- {

			n1 = int(n1Setup.absAllNumRunes[j]) - 48
			n2 = int(n2Setup.absAllNumRunes[i]) - 48
			n3 = (n1 * n2) + carry
			n4 = int(math.Mod(float64(n3), 10.00))

			intMAry[levels][place] = n4

			carry = n3 / 10

			place--
		}

		intMAry[levels][place] = carry
		carry = 0
		levels++
	}

	carry = 0
	n1 = 0
	n2 = 0
	n3 = 0
	n4 = 0
	for i := 0; i < lenLevels; i++ {
		for j := lenNumPlaces - 1; j >= 0; j-- {

			n1 = intFinalAry[j+1]
			n2 = intMAry[i][j]
			n3 = n1 + n2 + carry
			n4 = 0

			if n3 > 9 {
				n4 = int(math.Mod(float64(n3), 10.0))
				carry = n3 / 10

			} else {
				n4 = n3
				carry = 0
			}

			intFinalAry[j+1] = n4
		}

		if carry > 0 {
			intFinalAry[0] = carry
		}

	}

	numStrOut, err := nDto.FindIntArraySignificantDigitLimits(intFinalAry, newPrecision, newSignVal)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"- Error returned from nDto.FindIntArraySignificantDigitLimits(intFinalAry,newPrecision, "+
				"newSignVal). Error= %v", err)
	}

	return numStrOut, nil
}

// NewBigFloat - Creates a new NumStrDto instance from a Big Float value
// (*big.Float) and a precision specification.
func (nDto NumStrDto) NewBigFloat(
	bigFloat *big.Float, precision int) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewBigFloat() "

	numStr := bigFloat.Text('f', precision)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by n.ParseNumStr(numStr). "+
				"numStr='%v'  Error='%v'",
				numStr, err.Error())
	}

	return n2, nil

}

// NewBigInt - Creates a new NumStrDto instance from a signed big integer (*big.Int) and
// a precision specification.
func (nDto NumStrDto) NewBigInt(signedBigInt *big.Int, precision uint) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewBigInt() "

	n2, err := NumStrDto{}.ParseSignedBigInt(
		big.NewInt(0).Set(signedBigInt),
		precision)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by ParseSignedBigInt(signedBigInt, precision). "+
				"signedBigInt='%v' precision='%v'  Error='%v'",
				signedBigInt.Text(10), precision, err.Error())
	}

	err = n2.IsValid(ePrefix + "'n2' INVALID! ")

	if err != nil {
		return NumStrDto{}, err
	}

	return n2, nil
}

// NewBigIntNum - Receives a BigIntNum and converts it to a NumStrDto
// instance which is returned to the calling function.
func (nDto NumStrDto) NewBigIntNum(bINum BigIntNum) (NumStrDto, error) {
	ePrefix := "NumStrDto.NewBigIntNum() "
	n2, err := NumStrDto{}.ParseBigIntNum(bINum)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by ParseBigIntNum(bINum). "+
				"bINum='%v'  Error='%v'",
				bINum.GetNumStr(), err.Error())
	}

	return n2, nil

}

// NewFloat32 - Creates a new NumStrDto instance from a float32
// and precision specification.
func (nDto NumStrDto) NewFloat32(f32 float32, precision int) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewFloat32() "

	numStr := strconv.FormatFloat(float64(f32), 'f', precision, 32)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by n.ParseNumStr(numStr). "+
				"numStr='%v'  Error='%v'",
				numStr, err.Error())
	}

	return n2, nil

}

// NewFloat64 - Creates a new NumStrDto instance from a float64
// and precision specification.
func (nDto NumStrDto) NewFloat64(f64 float64, precision int) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewFloat64() "

	numStr := strconv.FormatFloat(f64, 'f', precision, 64)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by n.ParseNumStr(numStr). "+
				"numStr='%v'  Error='%v'",
				numStr, err.Error())
	}

	return n2, nil
}

// NewInt
//
// Creates a new NumStrDto from an int and a precision
// specification.
//
// Input parameter 'precision' indicates the number of
// digits to be formatted to the right of the decimal
// place.
//
// The 'NewInt' method is designed to used in conjunction
// with NumStrDto{} syntax thereby allowing NumStrDto type
// creation and initialization in one step.
//
// Example:
//
//	NumStrDto{}.NewInt(123456, 3)
//		Yields a new NumStrDto instance with a
//		numeric value of 123.456.
func (nDto NumStrDto) NewInt(intNum int, precision uint) NumStrDto {

	n2 := NumStrDto{}.NewInt64(int64(intNum), precision)

	return n2
}

// NewIntExponent - Returns a new NumStrDto instance. The numeric
// value is set using an integer multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = integer X 10^exponent
//
// Input parameter 'intNum' is of type int.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the Decimal{}
// syntax thereby allowing Decimal type creation and initialization in
// one step.
//
//		nDto := NumStrDto{}.NewIntExponent(123456, -3)
//	 -- nDto is now equal to "123.456", precision = 3
//
//		nDto := NumStrDto{}.NewIntExponent(123456, 3)
//	 -- decNum is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  intNum			 exponent			  	NumStrDto Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (nDto NumStrDto) NewIntExponent(intNum int, exponent int) NumStrDto {

	return NumStrDto{}.NewInt64Exponent(int64(intNum), exponent)
}

// NewInt32 - Creates a new NumStrDto instance from an int32 and a
// precision specification.
//
// Input parameter 'precision' indicates the number of digits to be
// formatted to the right of the decimal place.
//
// The 'NewInt32' method is designed to used in conjunction with
// NumStrDto{} syntax thereby allowing NumStrDto type creation and
// initialization in one step.
//
// Example: NumStrDto{}.NewInt32(123456, 3) yields a new NumStrDto
// instance with a numeric value of 123.456.
func (nDto NumStrDto) NewInt32(int32Num int32, precision uint) NumStrDto {

	n2 := NumStrDto{}.NewInt64(int64(int32Num), precision)

	return n2
}

// NewInt32Exponent - Returns a new NumStrDto instance. The numeric
// value is set using an int32 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int32 X 10^exponent
//
// For example, if exponent is -3, precision is set equal to 'int32Num'
// divided by 10^+3. Example:
//
//	  int32Num			exponent			NumStrDto Result
//		 123456		 		  -3							123.456
//
// If exponent is +3, int32Num is multiplied by 10 raised to the
// power of exponent and precision is set equal to exponent.
//
//	  int32Num			exponent			NumStrDto Result
//		 123456		 		   +3							123456.000
func (nDto NumStrDto) NewInt32Exponent(int32Num int32, exponent int) NumStrDto {

	return NumStrDto{}.NewInt64Exponent(int64(int32Num), exponent)
}

// NewInt64 - Creates a new NumStrDto instance from an int64 and a
// precision specification.
//
// Input parameter 'precision' indicates the number of digits to be
// formatted to the right of the decimal place.
//
// The 'NewInt64' method is designed to used in conjunction with
// NumStrDto{} syntax thereby allowing NumStrDto type creation and
// initialization in one step.
//
// Example: NumStrDto{}.NewInt64(123456, 3) yields a NumStrDto instance
// with a numeric value of 123.456.
func (nDto NumStrDto) NewInt64(i64 int64, precision uint) NumStrDto {
	ePrefix := "NumStrDto.NewInt64() "

	numStr := strconv.FormatInt(i64, 10)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	// This should never produce an error.
	if err != nil {
		sErr := fmt.Sprintf(ePrefix+
			"Fatal Error returned by NumStrDto{}.NewPtr().ParseNumStr(numStr). "+
			"numStr='%v' Error='%v'",
			numStr, err.Error())

		panic(sErr)
	}

	_ = n2.SetThisPrecision(precision, true)

	return n2
}

// NewInt64Exponent - Returns a new NumStrDto instance. The numeric
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
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//		nDto := NumStrDto{}.NewInt64Exponent(123456, -3)
//	 -- nDto is now equal to "123.456", precision = 3
//
//		nDto := NumStrDto{}.NewInt64Exponent(123456, 3)
//	 -- decNum is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  int64Num		 exponent			  	Decimal Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (nDto NumStrDto) NewInt64Exponent(int64Num int64, exponent int) NumStrDto {

	numStr := strconv.FormatInt(int64Num, 10)

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			numStr += "0"
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	var n2 NumStrDto

	if exponent == 0 {
		n2, _ = NumStrDto{}.NewNumStr(numStr)
	} else {
		n2, _ = nDto.ShiftPrecisionLeft(numStr, uint(exponent))
	}

	return n2
}

// NewUint - Creates a new NumStrDto instance from an uint and a
// precision specification.
//
// Input parameter 'precision' indicates the number of digits to be
// formatted to the right of the decimal place and is passed as type
// 'uint'
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//					uintNum := uint(123456)
//					precision := uint(3)
//					nDto := NumStrDto{}.NewUint(uintNum, precision)
//	       nDto is now equal to 123.456
//
// Examples:
// ---------
//
//	  uintNum			precision			NumStrDto Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (nDto NumStrDto) NewUint(uintNum uint, precision uint) NumStrDto {

	n2 := NumStrDto{}.NewUint64(uint64(uintNum), precision)

	return n2
}

// NewUintExponent - Returns a new NumStrDto instance. The numeric
// value is set using an uint value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int64 X 10^exponent
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//		nDto := NumStrDto{}.NewUintExponent(123456, -3)
//	 -- nDto is now equal to "123.456", precision = 3
//
//		nDto := NumStrDto{}.NewUintExponent(123456, 3)
//	 -- nDto is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  uintNum			exponent			NumStrDto Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (nDto NumStrDto) NewUintExponent(uintNum uint, exponent int) NumStrDto {

	return nDto.NewUint64Exponent(uint64(uintNum), exponent)
}

// NewUint32 - Creates a new NumStrDto instance from an uint32 and a
// precision specification.
//
// Input parameter 'precision' indicates the number of digits to be
// formatted to the right of the decimal place and is passed as type
// 'uint'.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//					uint32Num := uint32(123456)
//					precision := uint(3)
//					nDto := NumStrDto{}.NewUint32(uint32Num, precision)
//	       nDto is now equal to 123.456
//
// Examples:
// ---------
//
//	  uint32Num		precision			NumStrDto Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (nDto NumStrDto) NewUint32(uint32Num uint32, precision uint) NumStrDto {

	n2 := NumStrDto{}.NewUint64(uint64(uint32Num), precision)

	return n2
}

// NewUint32Exponent - Returns a new NumStrDto instance. The numeric
// value is set using an uint32 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int64 X 10^exponent
//
// Input parameter 'uint32Num' is of type uint32.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//		nDto := NumStrDto{}.NewUint32Exponent(123456, -3)
//	 -- nDto is now equal to "123.456", precision = 3
//
//		nDto := NumStrDto{}.NewUint32Exponent(123456, 3)
//	 -- nDto is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  uint32Num		exponent			NumStrDto Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (nDto NumStrDto) NewUint32Exponent(uint32Num uint32, exponent int) NumStrDto {

	return nDto.NewUint64Exponent(uint64(uint32Num), exponent)
}

// NewUint64 - Creates a new NumStrDto instance from an uint64 and a
// precision specification.
//
// Input parameter 'precision' indicates the number of digits to be
// formatted to the right of the decimal place.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//					uint64Num := uint64(123456)
//					precision := uint(3)
//					nDto := NumStrDto{}.NewUint64(uint64Num, precision)
//	       nDto is now equal to 123.456
//
// Examples:
// ---------
//
//	  uint64Num		precision			NumStrDto Result
//		 123456		 		   4							12.3456
//	  123456          0              123456
//	  123456          1              12345.6
func (nDto NumStrDto) NewUint64(uint64Num uint64, precision uint) NumStrDto {

	ePrefix := "NumStrDto.NewUint64() "

	numStr := strconv.FormatUint(uint64Num, 10)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)
	// This should NEVER produce an error
	if err != nil {
		sError := fmt.Sprintf(ePrefix+
			"Fatal Error returned by NumStrDto{}.NewPtr().ParseNumStr(numStr) "+
			"numStr='%v' Error='%v' ", numStr, err.Error())
		panic(sError)
	}

	_ = n2.SetThisPrecision(precision, true)

	_ = n2.SetNumericSeparatorsDto(nDto.GetNumericSeparatorsDto())

	return n2
}

// NewUint64Exponent - Returns a new NumStrDto instance. The numeric
// value is set using an uint64 value multiplied by 10 raised to the
// power of the 'exponent' parameter.
//
//	numeric value = int64 X 10^exponent
//
// Input parameter 'uint64Num' is of type uint64.
//
// Input parameter 'exponent' is of type int.
//
// Usage:
// ------
// This method is designed to be used in conjunction with the NumStrDto{}
// syntax thereby allowing NumStrDto type creation and initialization in
// one step.
//
//		nDto := NumStrDto{}.NewUint64Exponent(123456, -3)
//	 -- nDto is now equal to "123.456", precision = 3
//
//		nDto := NumStrDto{}.NewUint64Exponent(123456, 3)
//	 -- nDto is now equal to "123456.000", precision = 3
//
// Examples:
// ---------
//
//	  uint64Num		exponent			NumStrDto Result
//		 123456		 		  -3							123.456
//		 123456		 		   3							123456.000
//	  123456          0              123456
func (nDto NumStrDto) NewUint64Exponent(uint64Num uint64, exponent int) NumStrDto {

	ePrefix := "NumStrDto.NewUint64Exponent() "
	numStr := strconv.FormatUint(uint64Num, 10)

	if exponent > 0 {
		for i := 0; i < exponent; i++ {
			numStr += "0"
		}
	}

	if exponent < 0 {
		exponent = exponent * -1
	}

	var n2 NumStrDto
	var err error

	if exponent == 0 {
		n2, err = NumStrDto{}.NewNumStr(numStr)
		// This should never produce an error.
		if err != nil {
			sErr := fmt.Sprintf(ePrefix+
				"Fatal Error returned by NumStrDto{}.NewNumStr(numStr). "+
				"numStr='%v' uint64Num='%v' Error='%v'",
				numStr, uint64Num, err.Error())
			panic(sErr)
		}

	} else {
		n2, err = nDto.ShiftPrecisionLeft(numStr, uint(exponent))
		// This should never produce an error.
		if err != nil {
			sErr := fmt.Sprintf(ePrefix+
				"Fatal Error returned by nDto.ShiftPrecisionLeft("+
				"numStr, uint(exponent)). "+
				"numStr='%v' uint64Num='%v' exponent='%v' Error='%v'",
				numStr, uint64Num, exponent, err.Error())
			panic(sErr)
		}

	}

	_ = n2.SetNumericSeparatorsDto(nDto.GetNumericSeparatorsDto())

	return n2
}

// NewRational - Creates a new NumStrDto instance from a rational number and a precision
// specification.
//
// For information on Big Rational Numbers (*big.Rat), see https://golang.org/pkg/math/big/
func (nDto NumStrDto) NewRational(bigRat *big.Rat, precision int) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewRational() "
	numStr := bigRat.FloatString(precision)

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by n.ParseNumStr(numStr). "+
				"numStr='%v'  Error='%v'",
				numStr, err.Error())
	}

	return n2, nil
}

// NewNumStr - Used to create a populated NumStrDto instance.
// using a valid number string as an input parameter.
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
// separators, see method NumStrDto.NewNumStrWithNumSeps(), below.
//
// Usage Example:
//
//	n, err := NumStrDto{}.NewNumStr("123.456")
func (nDto NumStrDto) NewNumStr(numStr string) (NumStrDto, error) {

	ePrefix := "NumStrDto.NewNumStr() "

	n := NumStrDto{}.New()

	n2, err := n.ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"Error returned by n.ParseNumStr(numStrDto). "+
				"numStrDto='%v'  Error='%v'",
				numStr, err.Error())
	}

	return n2, nil
}

// NewNumStrWithNumSeps - Receives a number string as input and returns a
// new NumStrDto instance. The input parameter 'numSeps' contains numeric
// separators (decimal separator, thousands separator and currency symbol)
// which will be used to parse the number string.
//
// In addition, the numeric separators contained in input parameter 'numSeps'
// will be copied to the returned NumStrDto instance.
func (nDto NumStrDto) NewNumStrWithNumSeps(
	numStr string,
	numSeps NumericSeparatorDto) (NumStrDto, error) {

	ePrefix := "IntAry.NewNumStrWithNumSeps() "

	n := NumStrDto{}.New()

	numSeps.SetDefaultsIfEmpty()

	err := n.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"Error returned by  Ary.SetIntAryWithNumStr(numStr). "+
				"Error='%v' ", err.Error())
	}

	n2, err := n.ParseNumStr(numStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"Error returned by n.ParseNumStr(numStr). "+
				"numStr='%v', Error='%v' ", numStr, err.Error())
	}

	return n2, nil
}

// New - Used to create empty NumStrDto types.
// This message initializes the NumStrDto
// fields. This method will return the newly
// create type (not a pointer to the type).
// Example:
// n := NumStrDto{}.New()
// n2, err := n.ParseNumStr("123.456")
//
// Compare this method of object creation
// with that shown in the NewPtr() method.
func (nDto NumStrDto) New() NumStrDto {
	n := NumStrDto{}
	n.Empty()

	return n
}

// NewPtr - Used to create and initialize
// NumStrDto fields.  This method returns
// a pointer to the newly created NumStrDto
// type. As such, this method may be used
// to streamline the initialization process.
// Example:
// n, err := NumStrDto{}.NewPtr().ParseNumStr("123.456")
func (nDto NumStrDto) NewPtr() *NumStrDto {
	n := NumStrDto{}
	n.Empty()
	return &n
}

// NewZero
//
// Creates a new NumStrDto with a value of zero and a
// precision specified by input parameter 'precision'.
func (nDto NumStrDto) NewZero(precision uint) NumStrDto {

	numStr := "0"

	if precision > 0 {
		numStr += "."

		for i := uint(0); i < precision; i++ {
			numStr += "0"
		}
	}

	n2, err := NumStrDto{}.NewPtr().ParseNumStr(numStr)

	// This should NEVER produce an error.
	if err != nil {
		ePrefix := "NumStrDto.NewZero() "
		sErr := fmt.Sprintf(ePrefix+
			"Error returned by NumStrDto{}.NewPtr().ParseNumStr(numStr). "+
			"numStr='%v' Error='%v' ", numStr, err.Error())
		panic(sErr)
	}

	return n2
}

// ParseBigIntNum - Receives a BigIntNum instance and coverts it to a NumStrDto
// instance which is returned to the calling function.
func (nDto NumStrDto) ParseBigIntNum(biNum BigIntNum) (NumStrDto, error) {

	ePrefix := "NumStrDto.ParseBigIntNum() "

	nDto.SetNumericSeparatorsToDefaultIfEmpty()
	numSeps := nDto.GetNumericSeparatorsDto()

	n2Dto := NumStrDto{}.New()

	n2Dto.SetCurrencySymbol(biNum.GetCurrencySymbol())
	n2Dto.SetDecimalSeparator(biNum.GetDecimalSeparator())
	n2Dto.SetThousandsSeparator(biNum.GetThousandsSeparator())
	_ = n2Dto.SetSignValue(biNum.GetSign())
	n2Dto.precision = biNum.GetPrecisionUint()

	scratchNum := big.NewInt(0).Set(biNum.bigInt)

	if n2Dto.signVal < 0 {
		scratchNum.Neg(scratchNum)
	}

	bigZero := big.NewInt(0)
	bigTen := big.NewInt(int64(10))
	modulo := big.NewInt(0)
	modX := big.NewInt(0)
	n2Dto.absAllNumRunes = make([]rune, 0, 100)

	if scratchNum.Cmp(bigZero) == 0 {

		n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')

	} else {

		for scratchNum.Cmp(bigZero) == 1 {

			scratchNum, modulo = big.NewInt(0).QuoRem(scratchNum, bigTen, modX)

			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, rune(modulo.Int64()+int64(48)))
		}
	}

	lenAllNumRunes := len(n2Dto.absAllNumRunes)

	if int(n2Dto.precision) >= lenAllNumRunes {

		deltaNumRunes := int(n2Dto.precision) - lenAllNumRunes + 1

		for k := 0; k < deltaNumRunes; k++ {
			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')
			lenAllNumRunes++
		}

	}

	tRune := rune(0)

	if lenAllNumRunes > 1 {
		xLen := lenAllNumRunes - 1
		sortLimit := xLen / 2
		yCnt := 0
		for i := xLen; i > sortLimit; i-- {
			tRune = n2Dto.absAllNumRunes[yCnt]
			n2Dto.absAllNumRunes[yCnt] = n2Dto.absAllNumRunes[i]
			n2Dto.absAllNumRunes[i] = tRune
			yCnt++
		}
	}

	err := n2Dto.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+"Error returned by n2Dto.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' \n", err.Error())
	}

	err = n2Dto.IsValid("")

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+
				"NumStrDto INVALID! Error='%v'",
				err.Error())
	}

	return n2Dto, nil
}

// ParseSignedBigInt - receives a signed *Big Int number and precision parameter. It then
// generates and returns a new NumStrDto type.
func (nDto NumStrDto) ParseSignedBigInt(signedBigInt *big.Int, precision uint) (NumStrDto, error) {

	ePrefix := "NumStrDto.ParseSignedBigInt() "

	nDto.SetNumericSeparatorsToDefaultIfEmpty()
	numSeps := nDto.GetNumericSeparatorsDto()

	n2Dto := NumStrDto{}.New()

	n2Dto.SetCurrencySymbol(nDto.GetCurrencySymbol())
	n2Dto.SetDecimalSeparator(nDto.GetDecimalSeparator())
	n2Dto.SetThousandsSeparator(nDto.GetThousandsSeparator())
	n2Dto.precision = precision
	scratchNum := big.NewInt(0).Set(signedBigInt)
	bigZero := big.NewInt(0)
	n2Dto.signVal = 1

	if scratchNum.Cmp(bigZero) == -1 {
		scratchNum.Neg(scratchNum)
		n2Dto.signVal = -1
	}

	bigTen := big.NewInt(int64(10))
	modulo := big.NewInt(0)
	n2Dto.absAllNumRunes = make([]rune, 0, 100)

	if scratchNum.Cmp(bigZero) == 0 {

		n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')

	} else {

		for scratchNum.Cmp(bigZero) == 1 {
			modulo = big.NewInt(0).Rem(scratchNum, bigTen)
			scratchNum = big.NewInt(0).Quo(scratchNum, bigTen)
			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, rune(modulo.Int64()+int64(48)))
		}
	}

	lenAllNumRunes := len(n2Dto.absAllNumRunes)

	if int(n2Dto.precision) >= lenAllNumRunes {

		deltaNumRunes := int(n2Dto.precision) - lenAllNumRunes + 1

		for k := 0; k < deltaNumRunes; k++ {
			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, '0')
			lenAllNumRunes++
		}

	}

	tRune := rune(0)

	if lenAllNumRunes > 1 {
		xLen := lenAllNumRunes - 1
		sortLimit := xLen / 2
		yCnt := 0
		for i := xLen; i > sortLimit; i-- {
			tRune = n2Dto.absAllNumRunes[yCnt]
			n2Dto.absAllNumRunes[yCnt] = n2Dto.absAllNumRunes[i]
			n2Dto.absAllNumRunes[i] = tRune
			yCnt++
		}
	}

	err := n2Dto.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+"Error returned by n2Dto.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' \n", err.Error())
	}

	err = n2Dto.IsValid("")

	if err != nil {
		return NumStrDto{}.New(),
			fmt.Errorf(ePrefix+
				"NumStrDto INVALID! Error='%v'",
				err.Error())
	}

	return n2Dto, nil
}

// ParseNumStr - receives a raw string and converts to a properly
// formatted number string. The string is returned via a NumStrDto type.
// Returned number strings may consist of a leading negative sign ('-')
// numeric digits and may include a decimal separator ('.'). The NumStrDto
// breaks the string down into sign, Integer and Fractional components.
//
// The numeric separators (decimal separator, thousands separator and
// currency symbol) taken from the current NumStrDto instance will be
// copied to the NumStrDto instance returned by this method.
func (nDto *NumStrDto) ParseNumStr(str string) (NumStrDto, error) {

	ePrefix := "NumStrDto.ParseNumStr() "

	if len(str) == 0 {
		return NumStrDto{}, errors.New(ePrefix + "Received zero length number string as input!")
	}

	nDto.SetNumericSeparatorsToDefaultIfEmpty()
	numSeps := nDto.GetNumericSeparatorsDto()
	n2Dto := NumStrDto{}.New()

	n2Dto.signVal = 1
	_ = n2Dto.SetNumericSeparatorsDto(numSeps)
	baseRunes := []rune(str)
	lBaseRunes := len(baseRunes)
	isStartRunes := false
	isEndRunes := false
	isMinusSignFound := false
	//lCurRunes := len(NumStrCurrencySymbols)
	//isSkip := false
	isFractionalValue := false

	var absFracRunes []rune
	var absIntRunes []rune

	for i := 0; i < lBaseRunes && isEndRunes == false; i++ {

		if baseRunes[i] != '-' &&
			baseRunes[i] != n2Dto.decimalSeparator &&
			(baseRunes[i] < '0' || baseRunes[i] > '9') {

			continue

		} else if baseRunes[i] == '-' &&
			isMinusSignFound == false &&
			isStartRunes == false && isEndRunes == false &&
			i+1 < lBaseRunes &&
			((baseRunes[i+1] >= '0' && baseRunes[i+1] <= '9') ||
				baseRunes[i+1] == n2Dto.decimalSeparator) {

			isMinusSignFound = true
			n2Dto.signVal = -1
			isStartRunes = true
			continue

		} else if isEndRunes == false &&
			baseRunes[i] >= '0' && baseRunes[i] <= '9' {

			n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, baseRunes[i])
			isStartRunes = true

			if isFractionalValue {
				absFracRunes = append(absFracRunes, baseRunes[i])
			} else {
				absIntRunes = append(absIntRunes, baseRunes[i])
			}

		} else if isEndRunes == false &&
			i+1 < lBaseRunes &&
			baseRunes[i+1] >= '0' && baseRunes[i+1] <= '9' &&
			baseRunes[i] == n2Dto.decimalSeparator {

			isFractionalValue = true
			continue

		}

		if i == lBaseRunes-1 {

			isEndRunes = true

		}

	}

	lenAbsAllNumRunes := len(n2Dto.absAllNumRunes)

	if lenAbsAllNumRunes == 0 {
		nZeroNumStr := nDto.GetZeroNumStrDto(0)
		return nZeroNumStr, nil
	}

	lenAbsIntNumRunes := len(absIntRunes)
	if lenAbsIntNumRunes == 0 {
		absIntRunes = append(absIntRunes, '0')
	}

	lenAbsAllNumRunes = len(n2Dto.absAllNumRunes)
	lenAbsIntNumRunes = len(absIntRunes)
	lenAbsFracNumRunes := len(absFracRunes)

	isZeroVal := true

	for i := 0; i < lenAbsAllNumRunes; i++ {
		if n2Dto.absAllNumRunes[i] != '0' {
			isZeroVal = false
		}
	}

	if isZeroVal {
		nZeroDto := nDto.GetZeroNumStrDto(uint(lenAbsFracNumRunes))
		return nZeroDto, nil
	}

	if isFractionalValue {
		n2Dto.precision = uint(len(absFracRunes))
	}

	if lenAbsAllNumRunes != lenAbsIntNumRunes+lenAbsFracNumRunes {
		n2Dto.absAllNumRunes = []rune{}
		newLenAbsAllNumRunes := lenAbsIntNumRunes + lenAbsFracNumRunes

		for i := 0; i < newLenAbsAllNumRunes; i++ {
			if i < lenAbsIntNumRunes {
				n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, absIntRunes[i])
			} else {
				n2Dto.absAllNumRunes = append(n2Dto.absAllNumRunes, absFracRunes[i-lenAbsIntNumRunes])
			}
		}

		lenAbsAllNumRunes = len(n2Dto.absAllNumRunes)
	}

	// Validate n2Dto object
	err := n2Dto.IsValid(ePrefix)

	if err != nil {
		return NumStrDto{}, err
	}

	return n2Dto, nil

}

// ScaleNumStr - Shifts the position of the decimal point left or right depending
// on the value of input parameter 'scaleMode'.
//
// Input Parameters
// ================
//
// signedNumStr					 string -		A valid Signed Number String
//
// shiftPrecision 	 			 uint -		The number of positions which the decimal point
//
//	will be shifted. If 'shiftPrecision is Equal to
//	zero, no action will be taken, no error will be
//	issued and the original signedNumStr will be
//	returned.
//
// scaleMode	PrecisionScaleMode -	A constant with one of two Scale Mode values.
//
//	SCALEPRECISIONLEFT - 	Shifts the decimal point
//												from its current position
//												to the left.
//
//	SCALEPRECISIONRIGHT - Shifts the decimal point
//												from its current position
//												to the right.
//
// Note: 	See Methods NumStrDto.ShiftPrecisionRight() and NumStrDto.ShiftPrecisionLeft()
//
//	for additional information.
func (nDto *NumStrDto) ScaleNumStr(signedNumStr string,
	shiftPrecision uint,
	scaleMode PrecisionScaleMode) (NumStrDto, error) {

	ePrefix := "NumStrDto.ScaleNumStr() "

	n2Dto := NumStrDto{}

	var err error

	if scaleMode == SCALEPRECISIONLEFT {

		n2Dto, err = nDto.ShiftPrecisionLeft(signedNumStr, shiftPrecision)

		if err != nil {
			return NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned from nDto.ShiftPrecisionLeft(signedNumStr, shiftPrecision) "+
					"signedNumStr='%v' shiftPrecision='%v' scaleMode='%v' Error='%v' ",
					signedNumStr, shiftPrecision, scaleMode.String(), err.Error())

		}

	} else if scaleMode == SCALEPRECISIONRIGHT {

		n2Dto, err = nDto.ShiftPrecisionRight(signedNumStr, shiftPrecision)

		if err != nil {
			return NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned from nDto.ShiftPrecisionRight(signedNumStr, shiftPrecision) "+
					"signedNumStr='%v' shiftPrecision='%v' scaleMode='%v' Error='%v' ",
					signedNumStr, shiftPrecision, scaleMode.String(), err.Error())
		}

	} else {

		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"Error! Scale Mode is INVALID! "+
				"Scale Mode is NOT Equal to SCALEPRECISIONLEFT or SCALEPRECISIONRIGHT. scaleMode='%v' ",
				scaleMode.String())

	}

	return n2Dto, nil
}

// SetCurrencySymbol - assigns the input parameter rune as the
// currency symbol to be used by the NumStrDto.
//
// In the USA, the currency symbol is the dollar sign ('$').
//
// Note: If a zero value is submitted as input, Currency Symbol
// will default to the USA dollar sign ('$').
//
// For a list of Major Currency Unicode Symbols, see constants
// located in: MikeAustin71/mathopsgo/mathops/mathopsconstants.go
//
// Example: $123.45
func (nDto *NumStrDto) SetCurrencySymbol(currencySymbol rune) {

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	nDto.currencySymbol = currencySymbol
}

// SetDecimalSeparator - Assigns a rune or character to the internal
// data field, 'decimalSeparator'. The Decimal Separator is used to
// separate the integer and fractional elements of a number string.
//
// In the USA, the Decimal Separator is a period character ('.').
//
// Note: If a zero value is submitted as input, the Decimal Separator
// will default to the USA standard period character ('.').
//
// Example: 123.456
func (nDto *NumStrDto) SetDecimalSeparator(decimalSeparator rune) {

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	nDto.decimalSeparator = decimalSeparator
}

// SetThousandsSeparator - Sets the value of the character which will be
// used to separate thousands in the display of the NumStrDto number
// string. In the USA the typical thousands separator is the comma.
//
// If if a zero value is submitted, the Thousands Separator will default
// to the comma character.
//
// Example:
// 1,000,000
func (nDto *NumStrDto) SetThousandsSeparator(thousandsSeparator rune) {

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	nDto.thousandsSeparator = thousandsSeparator

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
//	signedNumStr		string		- A valid number string. The leading digit may optionally
//															be a '+' or '-' indicating numeric sign value. If '+'
//															or '-'	characters are not present in the first character
//															position, the number is assumed to represent a positive
//															numeric value ('+').
//
//	shiftPrecision		uint		- The number of digits by which the current decimal point
//															point position in the number string, 'signedNumStr' will
//															be shifted to the left.
//
// Returns
// =======
//
//	NumStrDto				- If successful, the method returns the result of the Shift Left precision
//										operation in the form of a 'NumStrDto' instance
//
//	error						- If successful, the 'error' type is set to 'nil'. In case of an error,
//										the 'error' instance returned will hold the error message.
//
// Examples
// ========
//
//	                   Shift-Left
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
func (nDto *NumStrDto) ShiftPrecisionLeft(
	signedNumStr string,
	shiftLeftPrecision uint) (NumStrDto, error) {

	ePrefix := "NumStrDto.ShiftPrecisionLeft() "

	if len(signedNumStr) == 0 {
		return NumStrDto{}, errors.New(ePrefix +
			"Received zero length number string!")
	}

	// Set defaults for thousands separators,
	// decimal separators and currency Symbols
	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	n1, err := NumStrDto{}.NewPtr().ParseNumStr(signedNumStr)

	if err != nil {
		return NumStrDto{}, fmt.Errorf(ePrefix+
			"Received Error from NumStrDto.ParseNumStr(signedNumStr). "+
			"str= '%v' Error= %v",
			signedNumStr, err)
	}

	n2 := NumStrDto{}.New()

	n2.thousandsSeparator = nDto.thousandsSeparator
	n2.decimalSeparator = nDto.decimalSeparator
	n2.currencySymbol = nDto.currencySymbol
	n2.signVal = n1.signVal
	n2.precision = shiftLeftPrecision + n1.precision
	iTotalSpecPrecision := int(n2.precision)
	lenAbsAllNumRunes := len(n1.absAllNumRunes)
	lenAbsIntRunes := n1.GetAbsIntRunesLength()
	lenAbsFracRunes := n1.GetAbsFracRunesLength()

	if nDto.IsNumStrZeroValue(&n1) {

		return nDto.GetZeroNumStrDto(n2.precision), nil
	}

	if iTotalSpecPrecision == lenAbsAllNumRunes {

		n2.absAllNumRunes = append(n2.absAllNumRunes, '0')

	} else if iTotalSpecPrecision > lenAbsAllNumRunes {

		deltaPrecision := iTotalSpecPrecision - lenAbsAllNumRunes + 1

		for i := 0; i < deltaPrecision; i++ {
			n2.absAllNumRunes = append(n2.absAllNumRunes, '0')
		}

	}

	for j := 0; j < lenAbsAllNumRunes; j++ {
		n2.absAllNumRunes = append(n2.absAllNumRunes, n1.absAllNumRunes[j])
	}

	lenAbsAllNumRunes = len(n2.absAllNumRunes)
	lenAbsFracRunes = iTotalSpecPrecision
	lenAbsIntRunes = lenAbsAllNumRunes - lenAbsFracRunes

	if lenAbsIntRunes <= 0 {
		return NumStrDto{}, fmt.Errorf(ePrefix+
			"Calculated number of integer digits is less than or equal to ZERO. "+
			"lenAbsIntRunes= '%v' ",
			lenAbsIntRunes)
	}

	lenAbsFracRunes = n2.GetAbsFracRunesLength()

	err = n2.IsValid(ePrefix)

	if err != nil {
		return NumStrDto{}, err
	}

	return n2, nil
}

// ShiftPrecisionRight - Shifts the existing precision of a number string. The position of
// the decimal point is shifted 'precision' positions to the right.
//
// This is equivalent to: result = signedNumStr X 10^precision or signedNumStr Multiplied
// by 10 raised to the power of precision.
//
// Examples:
// signedNumStr			precision			Result
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
func (nDto *NumStrDto) ShiftPrecisionRight(signedNumStr string, precision uint) (NumStrDto, error) {

	ePrefix := "NumStrDto.ShiftPrecisionRight() "

	if len(signedNumStr) == 0 {
		return NumStrDto{}, errors.New(ePrefix + "Received zero length number string as input!")
	}

	// Set defaults for thousands separators,
	// decimal separators and currency Symbols
	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	n1, err := NumStrDto{}.NewPtr().ParseNumStr(signedNumStr)

	if err != nil {
		return NumStrDto{}, fmt.Errorf(ePrefix+"- Received Error from NumStrDto.ParseNumStr(signedNumStr). str= '%v' Error= %v", signedNumStr, err)
	}

	n2 := NumStrDto{}.New()

	iTotalSpecPrecision := 0
	iPrecision := int(precision)
	iN1Precision := int(n1.precision)

	if iN1Precision > 0 && iPrecision < iN1Precision {
		iTotalSpecPrecision = iN1Precision - iPrecision
	} else {
		iTotalSpecPrecision = 0
	}

	n2.thousandsSeparator = nDto.thousandsSeparator
	n2.decimalSeparator = nDto.decimalSeparator
	n2.currencySymbol = nDto.currencySymbol
	n2.signVal = n1.signVal
	n2.precision = uint(iTotalSpecPrecision)

	lenAbsAllNumRunes := len(n1.absAllNumRunes)

	if nDto.IsNumStrZeroValue(&n1) {

		return nDto.GetZeroNumStrDto(0), nil
	}

	if int(precision) > int(n1.precision) {

		for i := 0; i < lenAbsAllNumRunes; i++ {
			n2.absAllNumRunes = append(n2.absAllNumRunes, n1.absAllNumRunes[i])
		}

		deltaPrecision := int(precision) - int(n1.precision)

		for i := 0; i < deltaPrecision; i++ {
			n2.absAllNumRunes = append(n2.absAllNumRunes, '0')
		}

	} else {

		for i := 0; i < lenAbsAllNumRunes; i++ {
			n2.absAllNumRunes = append(n2.absAllNumRunes, n1.absAllNumRunes[i])
		}

	}

	lenAbsAllNumRunes = len(n2.absAllNumRunes)
	lenAbsFracRunes := iTotalSpecPrecision
	lenAbsIntRunes := lenAbsAllNumRunes - lenAbsFracRunes

	if lenAbsIntRunes <= 0 {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"- Calculated number of integer digits is less than or equal to ZERO. "+
				"lenAbsIntRunes= '%v' ", lenAbsIntRunes)
	}

	lenAbsFracRunes = n2.GetAbsFracRunesLength()

	err = n2.IsValid(ePrefix)

	if err != nil {
		return NumStrDto{}, err
	}

	return n2, nil
}

// SetNumericSeparators - Used to assign values for the Decimal and Thousands separators as well
// as the Currency Symbol to be used in displaying the current number string.
//
// Note: If zero values are submitted as input, the values will default to USA standards.
//
// USA Examples:
//
// Decimal Separator period ('.') 		= 123.456
// Thousands Separator comma (',') 		= 1,000,000,000
// Currency Symbol dollar sign ('$')	= $123
func (nDto *NumStrDto) SetNumericSeparators(
	decimalSeparator,
	thousandsSeparator,
	currencySymbol rune) {

	if decimalSeparator == 0 {
		decimalSeparator = '.'
	}

	if thousandsSeparator == 0 {
		thousandsSeparator = ','
	}

	if currencySymbol == 0 {
		currencySymbol = '$'
	}

	nDto.decimalSeparator = decimalSeparator
	nDto.thousandsSeparator = thousandsSeparator
	nDto.currencySymbol = currencySymbol
}

// SetNumericSeparatorsDto - Sets the values of numeric separators:
//
//	decimal point separator
//	thousands separator
//	currency symbol
//
// based on values transmitted through input parameter 'customSeparators'.
//
// If any of the values contained in input parameter 'customSeparators' is set
// to zero or nil, an error will be returned.
func (nDto *NumStrDto) SetNumericSeparatorsDto(customSeparators NumericSeparatorDto) error {

	ePrefix := "NumStrDto.SetNumericSeparatorsDto() "

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

	nDto.decimalSeparator = customSeparators.DecimalSeparator
	nDto.thousandsSeparator = customSeparators.ThousandsSeparator
	nDto.currencySymbol = customSeparators.CurrencySymbol

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
func (nDto *NumStrDto) SetNumericSeparatorsToDefaultIfEmpty() {

	if nDto.GetDecimalSeparator() == 0 {
		nDto.SetDecimalSeparator('.')
	}

	if nDto.GetThousandsSeparator() == 0 {
		nDto.SetThousandsSeparator(',')
	}

	if nDto.GetCurrencySymbol() == 0 {
		nDto.SetCurrencySymbol('$')
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
//	nDto.SetDecimalSeparator()
//	nDto.SetThousandsSeparator()
//	nDto.SetCurrencySymbol()
func (nDto *NumStrDto) SetNumericSeparatorsToUSADefault() {
	nDto.SetDecimalSeparator('.')
	nDto.SetThousandsSeparator(',')
	nDto.SetCurrencySymbol('$')
}

// SetNumStr - Sets the value of the current NumStrDto instance
// to the number string received as input.
func (nDto *NumStrDto) SetNumStr(numStr string) error {

	ePrefix := "NumStrDto.SetNumStr() "

	numSeps := nDto.GetNumericSeparatorsDto()

	n2, err := NumStrDto{}.NewNumStr(numStr)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by NumStrDto{}.NewNumStr(numStr). "+
			"numStr='%v' Error='%v' ", numStr, err.Error())
	}

	err = n2.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by n2.SetNumericSeparatorsDto(numSeps) "+
			"Error='%v' \n", err.Error())
	}

	nDto.CopyIn(n2)

	return nil

}

// SetPrecision - parses the incoming number string and applies the designated 'precision'. 'precision'
// determines the number of digits to the right of the decimal place. The boolean parameter 'roundResult'
// is used to apply rounding in those cases where 'precision' dictates a reduction in the number of
// digits to the right of the decimal place. See 'Examples' below.
//
// Input Parameters
// ================
//
//		signedNumStr 	string	- A valid number string
//
//		precision 		uint		- The 'precision' values designates the number of places to the right of the
//														decimal point which will be realized upon completion of this operation.
//
//		roundResult 	bool		- If the 'precision' value is less than the current number of places to the
//														right of the decimal point, this method will truncate the existing fractional
//														digits. If 'roundResult' is set to true, this truncation operation will
//														include rounding the last digit.
//
//
//	 Examples
//	 ========
//
//			------------ Input Parameters ------------        	---- Result -------
//
// Test
//
//	 No		signedNumStr			precision			roundResult
//		1	 		"123456789"				  7							false						 "123456789.0000000"
//		2 		"123456789"				  7							true						 "123456789.0000000"
//	 3  	 "-123456789"				  7							false						"-123456789.0000000"
//	 4  	 "-123456789"				  7							true						"-123456789.0000000"
//		5  	 "123456.789"					2							true						 "123456.79"
//	 6  	 "123456.789"					2							false						 "123456.78"
//		7  	 "123456.789"      		5             false						 "123456.78900"
//	 8  	 "123.456789"         1             false            "123.4"
//	 9  	 "123.456789"         1             true             "123.5"
//
// 10   "-123.456789"         1             false           "-123.4"
// 11 	"-123.456789"         1             true            "-123.5"
// 12  	 "123456.789"					0							true	           "123457"
// 13 	"-123456.789"         0							true						"-123457"
// 14 	 "123456.789"					0							false	           "123456"
// 15 	"-123456.789"					0							false						"-123456"
// 16 	 "123457"             1							false						 "123457.0"
// 17 	 "123457"             1							true						 "123457.0"
// 18 	"-123457"							1							false						"-123457.0"
// 19 	"-123457"							1							true						"-123457.0"
func (nDto *NumStrDto) SetPrecision(
	signedNumStr string,
	precision uint,
	roundResult bool) (NumStrDto, error) {

	ePrefix := "NumStrDto.SetPrecision() "

	if len(signedNumStr) == 0 {
		return NumStrDto{},
			errors.New(ePrefix + "Received zero length number string!")
	}

	// Set defaults for thousands separators,
	// decimal separators and currency Symbols
	if nDto.thousandsSeparator == 0 {
		nDto.thousandsSeparator = ','
	}

	if nDto.decimalSeparator == 0 {
		nDto.decimalSeparator = '.'
	}

	if nDto.currencySymbol == 0 {
		nDto.currencySymbol = '$'
	}

	n0 := NumStrDto{}.New()
	n0.thousandsSeparator = nDto.thousandsSeparator
	n0.decimalSeparator = nDto.decimalSeparator
	n0.currencySymbol = nDto.currencySymbol

	n1, err := n0.ParseNumStr(signedNumStr)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+
				"Error returned from ns.ParseNumString(signedNumStr). "+
				"signedNumStr='%v' Error= %v", signedNumStr, err)
	}

	n2 := NumStrDto{}.New()

	n2.signVal = n1.signVal
	n2.precision = precision
	n2.thousandsSeparator = nDto.thousandsSeparator
	n2.decimalSeparator = nDto.decimalSeparator
	n2.currencySymbol = nDto.currencySymbol
	n2AbsIntRunes := n2.GetAbsIntRunes()
	n2AbsFracRunes := n2.GetAbsFracRunes()

	iSpecPrecision := int(precision)
	lenN1AbsAllNumRunes := len(n1.absAllNumRunes)
	n1AbsIntRunes := n1.GetAbsIntRunes()
	n1AbsFracRunes := n1.GetAbsFracRunes()
	lenN1AbsIntRunes := len(n1AbsIntRunes)
	lenN1AbsFracRunes := len(n1AbsFracRunes)
	totalRunes := 0

	if roundResult && lenN1AbsFracRunes > 0 &&
		iSpecPrecision < lenN1AbsFracRunes {

		absAllNumsToRound, isOk := big.NewInt(0).SetString(string(n1.absAllNumRunes), 10)

		if !isOk {
			return NumStrDto{},
				fmt.Errorf(ePrefix+"Error: Failed to convert string to big.Int(). "+
					"big.Int.SetString(n1.absAllNumRunes). n1.absAllNumRunes='%v' ",
					string(n1.absAllNumRunes))
		}

		bigDeltaPrecision := big.NewInt(int64(lenN1AbsFracRunes - iSpecPrecision - 1))
		base10 := big.NewInt(int64(10))
		roundUp5 := big.NewInt(int64(5))
		roundScaleFactor := big.NewInt(0).Exp(base10, bigDeltaPrecision, nil)
		roundUpNum := big.NewInt(0).Mul(roundUp5, roundScaleFactor)
		roundedAbsAllNums := big.NewInt(0).Add(absAllNumsToRound, roundUpNum)
		actualDeltaPrecision := big.NewInt(int64(lenN1AbsFracRunes - iSpecPrecision))
		actualDeltaScaleFactor := big.NewInt(0).Exp(base10, actualDeltaPrecision, nil)
		actualAbsAllNums := big.NewInt(0).Div(roundedAbsAllNums, actualDeltaScaleFactor)
		n1.absAllNumRunes = []rune{}
		n1AbsIntRunes = []rune{}
		n1AbsFracRunes = []rune{}
		n1.absAllNumRunes = []rune(actualAbsAllNums.String())
		lenN1AbsAllNumRunes = len(n1.absAllNumRunes)

		for i := 0; i < lenN1AbsAllNumRunes; i++ {

			if i < lenN1AbsIntRunes {
				n1AbsIntRunes = append(n1AbsIntRunes, n1.absAllNumRunes[i])
			} else {
				n1AbsFracRunes = append(n1AbsFracRunes, n1.absAllNumRunes[i])
			}
		}

		lenN1AbsIntRunes = len(n1AbsIntRunes)
		lenN1AbsFracRunes = len(n1AbsFracRunes)

		if lenN1AbsAllNumRunes != (lenN1AbsIntRunes + lenN1AbsFracRunes) {

			return NumStrDto{},
				fmt.Errorf(ePrefix+"Error on Rounding. lenN1AbsAllNumRunes != "+
					"(lenN1AbsIntRunes + lenN1AbsFracRunes). lenN1AbsAllNumRunes= '%v' "+
					"lenN1AbsIntRunes= '%v' lenN1AbsFracRunes= '%v'",
					lenN1AbsAllNumRunes, lenN1AbsIntRunes, lenN1AbsFracRunes)
		}

	}

	if lenN1AbsIntRunes == 0 {
		n2.absAllNumRunes = append(n2.absAllNumRunes, '0')
		n2AbsIntRunes = append(n2AbsIntRunes, '0')
	}

	totalRunes = lenN1AbsIntRunes + iSpecPrecision

	for i := 0; i < totalRunes; i++ {

		if i < lenN1AbsAllNumRunes {
			n2.absAllNumRunes = append(n2.absAllNumRunes, n1.absAllNumRunes[i])
		} else {
			n2.absAllNumRunes = append(n2.absAllNumRunes, '0')
		}

		if i < lenN1AbsIntRunes {

			n2AbsIntRunes = append(n2AbsIntRunes, n1.absAllNumRunes[i])

		} else {

			if i < lenN1AbsAllNumRunes {
				n2AbsFracRunes = append(n2AbsFracRunes, n1.absAllNumRunes[i])
			} else {
				n2AbsFracRunes = append(n2AbsFracRunes, '0')
			}
		}
	}

	err = n2.IsValid(ePrefix)

	if err != nil {

		return NumStrDto{}, err
	}

	return n2, nil
}

// SetSignValue - Sets the sign of the numeric value
// for the current NumStrDto. Only two values are
// allowed: +1 and -1. If any other value is passed
// an error is thrown
func (nDto *NumStrDto) SetSignValue(newSignVal int) error {

	ePrefix := "NumStrDto.SetSignValue() "

	if newSignVal != -1 && newSignVal != 1 {
		return fmt.Errorf(ePrefix+
			"Invalid sign value passed. sign must be +1 or -1. "+
			"This sign value= %v", newSignVal)
	}

	nDto.signVal = newSignVal

	return nil
}

// SetThisPrecision - Sets precision for the current NumStrDto instance.
// 'precision' identifies the number of decimal places to the right of the
// decimal point.
//
// Input Parameters
// ================
//
//	precision 		uint		- The 'precision' values designates the number of places to the right of the
//													decimal point which will be realized upon completion of this operation. The
//													precision operation will be performed on the number string contained in the
//													NumStrDto instance.
//
//	roundResult 	bool		- If the 'precision' value is less than the current number of places to the
//													right of the decimal point, this method will truncate the existing fractional
//													digits. If 'roundResult' is set to true, this truncation operation will
//													include rounding the last digit.
func (nDto *NumStrDto) SetThisPrecision(
	precision uint,
	roundResult bool) error {

	ePrefix := "NumStrDto.SetThisPrecision() "

	n2, err := nDto.SetPrecision(nDto.GetNumStr(), precision, roundResult)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by nDto.SetPrecision(signedNumStr, precision, "+
			"roundResult). nDto.numStrDto='%v' precision='%v', roundResult='%v'",
			nDto.GetNumStr(), precision, roundResult)
	}

	nDto.CopyIn(n2)

	return nil
}

// Subtract - Subtracts the value of an input NumStrDto from the
// current NumStrDto instance.
func (nDto *NumStrDto) Subtract(n2Dto NumStrDto) error {

	ePrefix := "NumStrDto.Subtract() "

	n1Dto := nDto.CopyOut()

	nResult, err := nDto.SubtractNumStrs(n1Dto, n2Dto)

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by nDto.SubtractNumStrs(n1Dto, n2Dto). "+
			"Error='%v'", err.Error())
	}

	nDto.CopyIn(nResult)

	return nil
}

// SubtractNumStrs - Subtracts the numeric values represented by two NumStrDto
// objects.
func (nDto *NumStrDto) SubtractNumStrs(n1Dto, n2Dto NumStrDto) (NumStrDto, error) {

	ePrefix := "NumStrDto.SubtractNumStrs() "

	n1NumDto, n2NumDto, compare, isReversed, err := nDto.FormatForMathOps(n1Dto, n2Dto)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"- Error from nDto.FormatForMathOps(n1Dto, n2Dto). "+
				"Error= %v", err)
	}

	if compare == 0 {
		return nDto.GetZeroNumStrDto(n1NumDto.precision), nil
	}

	newSignVal := n1NumDto.signVal
	precision := n1NumDto.precision

	if n1NumDto.signVal != n2NumDto.signVal {

		err = n1NumDto.SetSignValue(1)

		if err != nil {
			return NumStrDto{},
				fmt.Errorf(ePrefix+"- Error from n1NumDto.SetSignValue(1). Error= %v", err)
		}

		err = n2NumDto.SetSignValue(1)

		if err != nil {
			return NumStrDto{},
				fmt.Errorf(ePrefix+"- Error from n2NumDto.SetSignValue(1). Error= %v", err)
		}

		nOutDto, err := nDto.AddNumStrs(n1NumDto, n2NumDto)

		if err != nil {
			return NumStrDto{},
				fmt.Errorf(ePrefix+"- Error from nDto.AddNumStrs(n1NumDto, n2NumDto). "+
					"Error= %v", err)
		}

		_ = nOutDto.SetSignValue(newSignVal)

		return nOutDto, nil
	}

	// Change sign for subtraction
	newSignVal = n1NumDto.signVal

	if isReversed {
		newSignVal = newSignVal * -1
	}

	lenN1AllRunes := len(n1NumDto.absAllNumRunes)

	n1IntAry := make([]int, lenN1AllRunes)
	n2IntAry := make([]int, lenN1AllRunes)
	n3IntAry := make([]int, lenN1AllRunes)

	for i := 0; i < lenN1AllRunes; i++ {

		n1IntAry[i] = int(n1NumDto.absAllNumRunes[i]) - 48
		n2IntAry[i] = int(n2NumDto.absAllNumRunes[i]) - 48

	}

	carry := 0
	n1 := 0
	n2 := 0
	n3 := 0
	// Main Subtraction Routine
	for j := lenN1AllRunes - 1; j >= 0; j-- {

		n1 = n1IntAry[j]
		n2 = n2IntAry[j]
		n3 = 0

		if n1-carry-n2 < 0 {
			n1 += 10
			n3 = n1 - n2 - carry
			carry = 1
		} else {
			n3 = n1 - n2 - carry
			carry = 0
		}

		n3IntAry[j] = n3

	}

	nOutDto, err := nDto.FindIntArraySignificantDigitLimits(n3IntAry, precision, newSignVal)

	if err != nil {
		return NumStrDto{},
			fmt.Errorf(ePrefix+"- Error from final nDto.FindIntArraySignificantDigitLimits"+
				"(n3IntAry, precision, newSignVal). precision='%v' newSignVal='%v' Error= %v",
				precision, newSignVal, err)
	}

	return nOutDto, nil
}
