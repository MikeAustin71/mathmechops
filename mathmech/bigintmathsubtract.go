package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

// BigIntMathSubtract - Contains methods used to perform subtraction
// operations on *big.Int numeric types.
//
// 			minuend − subtrahend = difference
//
type BigIntMathSubtract struct {
	Input  BigIntPair // BigIntPair.Big1 = minuend  BigIntPair.Big2 = subtrahend
	Result BigIntNum  // The result of the subtraction otherwise known as the 'difference'
}

// BigIntSubtract - Performs the subtraction operation on two
// numeric values. The 'minuend' is the number from which the
// 'subtrahend' is subtracted in order to generate a result or
// difference between the two numbers.
//
// In the subtraction operation:
//
//				'minuend' - 'subtrahend' = difference or result
//
// This method provides for the subtraction of fixed length
// floating point values by means of integer and precision
// specification pairs.
//
// As an example, consider the following subtraction operation:
//
// 						752.314 - 21.67894 = 730.63506 = difference
//
// In this case the 'minuend', 'subtrahend' and 'difference' would be
// configured as follows:
//
//									minuend 						= 752314
//                  minuendPrecision		= 3
//                  subtrahend 					= 2167894
//                  subtrahendPrecision = 5
//
//                  difference					= 73063506
//                  differencePrecision = 5
//
// In this way, the method uses integer, precision pairs to define fixed
// length floating point numbers.
//
// Input Parameters
// ================
//
//	minuend 						*big.Int		- The number from which the subtrahend will be subtracted
//
//	minPrecision 				*big.Int		- The 'minuend' precision or numeric digits after
//																		the decimal point. 'minPrecision' must be greater
//                                    than or equal to zero.
//
//	subtrahend 					*big.Int		- The number to be subtracted from the 'minuend'.
//
//	subPrecision 				*big.Int 		- The 'subtrahend' precision or numeric digits after
//																		the decimal point. 'subPrecision' must be greater
//                                    than or equal to zero.
//
// Return Values
// =============
//
// difference 					*big.Int		- The difference or result of the subtraction
// 																		operation	returned as a *big.Int type.
//
// differencePrecision	*big.Int   	- The precision specification for the returned
//                          					subtraction 'result'. Precision specifies the
//                          					number of fractional digits to the right of the
//                          					decimal place. 'differencePrecision' will always
//                                    be greater than or equal to zero.
//
// Taken together, 'difference' and 'differencePrecision' can define a fixed
// length floating point number.
//
// Note: This function will delete all trailing fractional zeros from the result or
//       difference.
//
func (bSubtract BigIntMathSubtract) BigIntSubtract(
	minuend,
	minPrecision,
	subtrahend,
	subPrecision *big.Int) (difference *big.Int, differencePrecision *big.Int, err error) {

	ePrefix := "BigIntMathSubtract.BigIntSubtract() "

	difference = big.NewInt(0)
	differencePrecision = big.NewInt(0)
	err = nil

	if minuend == nil {
		minuend = big.NewInt(0)
	}

	if minPrecision == nil {
		minPrecision = big.NewInt(0)
	}

	if subtrahend == nil {
		subtrahend = big.NewInt(0)
	}

	if subPrecision == nil {
		subPrecision = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if minPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'minPrecision' is LESS THAN ZERO! "+
			"minPrecision='%v' ", minPrecision.Text(10))

		return difference, differencePrecision, err
	}

	if subPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'subPrecision' is LESS THAN ZERO! "+
			"subPrecision='%v' ", subPrecision.Text(10))

		return difference, differencePrecision, err
	}

	base10 := big.NewInt(10)
	deltaPrecision := big.NewInt(0)
	scale := big.NewInt(0)

	if minPrecision.Cmp(subPrecision) == 0 {
		// Precisions are equal.

		difference = big.NewInt(0).Sub(minuend, subtrahend)
		differencePrecision = big.NewInt(0).Set(minPrecision)

	} else if minPrecision.Cmp(subPrecision) == 1 {
		//  minPrecision > subPrecision
		deltaPrecision = big.NewInt(0).Sub(minPrecision, subPrecision)
		scale = big.NewInt(0).Exp(base10, deltaPrecision, nil)
		newSubInt := big.NewInt(0).Mul(subtrahend, scale)

		difference = big.NewInt(0).Sub(minuend, newSubInt)
		differencePrecision.Set(minPrecision)

	} else {
		// subPrecision must be GREATER THAN minPrecision
		deltaPrecision = big.NewInt(0).Sub(subPrecision, minPrecision)
		scale = big.NewInt(0).Exp(base10, deltaPrecision, nil)
		newMinuendInt := big.NewInt(0).Mul(minuend, scale)

		difference = big.NewInt(0).Sub(newMinuendInt, subtrahend)
		differencePrecision.Set(subPrecision)
	}

	if difference.Cmp(bigZero) == 0 {
		differencePrecision = big.NewInt(0)
	}

	// Delete trailing fractional zeros
	if differencePrecision.Cmp(bigZero) == 1 {
		// differencePrecision > 0
		scrap := big.NewInt(0)
		biBase10 := big.NewInt(10)
		biBaseZero := big.NewInt(0)
		newDifference, mod10 := big.NewInt(0).QuoRem(difference, biBase10, scrap)
		bigOne := big.NewInt(1)

		for mod10.Cmp(biBaseZero) == 0 && differencePrecision.Cmp(bigZero) == 1 {
			difference.Set(newDifference)
			differencePrecision.Sub(differencePrecision, bigOne)
			newDifference, mod10 = big.NewInt(0).QuoRem(difference, biBase10, scrap)
		}
	}

	err = nil

	return difference, differencePrecision, err
}

// FixedDecimalSubtract - Performs the subtraction operation on two
// BigIntFixedDecimal types. The subtraction result is also returned
// as a BigIntFixedDecimal type.
//
// Examples:
// =========
//
// In the subtraction operation:
//								'minuend' - 'subtrahend' = 'difference' or result
//                 752.314  -    21.67894  = 730.63506 = difference
//
// For this method 'minuend', 'subtrahend' and 'difference' are configured
// as BigIntFixedDecimal types.
//
// The BigIntFixedDecimal is used to defined fixed length floating point
// numbers and is defined as follows:
//
// type BigIntFixedDecimal struct {
//
//	integerNum *big.Int  -	All of the numeric digits, both integer and fractional,
// 													necessary to define a fixed length floating point number.
// 													The number of digits to the right of the decimal place
// 													is specified by the data field,
// 													BigIntFixedDecimal.precision.
//
//	precision  uint				- Specifies the number of digits to the right of the decimal
// 													place in the series of numeric digits represented by data
// 													field BigIntFixedDecimal.integerNum.
//
// }
//
//
// 	To represent the floating point number 52.459	a BigIntDecimal Structure
// 	would be configured as follows:
//
// 			BigIntFixedDecimal.integerNum	= 52459
// 			BigIntFixedDecimal.precision	= 3
//
// As an example consider the following subtraction operation:
// 						752.314 - 21.67894 = 730.63506 = difference
//
// In this case the 'minuend', 'subtrahend' and 'difference' consist of
// BigIntFixedDecimal types configured as follows:
//
//									minuend.integerNum		= 752314
//                  minuend.precision			= 3
//                  subtrahend.integerNum	= 2167894
//                  subtrahend.precision 	= 5
//
//                  difference.integerNum	= 73063506
//                  difference.precision 	= 5
//
//
// Input Parameters
// ================
//
//	minuend BigIntFixedDecimal		- The number from which the subtrahend will be
// 																	subtracted.
//
//	subtrahend BigIntFixedDecimal	- The number to be subtracted from the 'minuend'.
//
// Return Values
// =============
//
// difference BigIntFixedDecimal	- The difference or result of the subtraction
// 																	operation returned as a BigIntFixedDecimal type.
//                                      'minuend' - 'subtrahend' = 'difference'
//
func (bSubtract BigIntMathSubtract) FixedDecimalSubtract(
	minuend BigIntFixedDecimal,
	subtrahend BigIntFixedDecimal) (difference BigIntFixedDecimal) {

	difference = BigIntFixedDecimal{}.NewZero(0)

	minuend.IsValid()
	subtrahend.IsValid()

	// The error is ignored because minuend precision and
	// subtrahend precision will never be less than zero
	rBigInt, rBigIntPrecision, _ :=
		BigIntMathSubtract{}.BigIntSubtract(
			minuend.GetInteger(),
			minuend.GetPrecisionBigInt(),
			subtrahend.GetInteger(),
			subtrahend.GetPrecisionBigInt())

	difference.SetNumericValue(rBigInt, uint(rBigIntPrecision.Uint64()))

	return difference
}

// SubtractBigInts - Performs the subtraction operation.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// Input Parameters
// ================
//
//	minuend *big.Int					- The number from which the subtrahend will be subtracted
//	minPrecision uint,				- The 'minuend' precision or numeric digits after
//																	the decimal point.
//	subtrahend *big.Int,			- The number to be subtracted from the 'minuend'.
//	subPrecision uint  				- The 'subtrahend' precision or numeric digits after
//																	the decimal point.
//
// Return Values
// =============
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// default numeric separators (decimal separator, thousands separator and
// currency symbol).
//
// The returned BigIntNum 'result' will contain USA default numeric separators
// (decimal separator, thousands separator and currency symbol).
//
func (bSubtract BigIntMathSubtract) SubtractBigInts(
	minuend *big.Int,
	minPrecision uint,
	subtrahend *big.Int,
	subPrecision uint) BigIntNum {

	// Returned error is ignored because minuend and subtrahend precisions
	// will never be less than zero.
	result, resultPrecision, _ :=
		BigIntMathSubtract{}.BigIntSubtract(
			minuend,
			big.NewInt(0).SetUint64(uint64(minPrecision)),
			subtrahend,
			big.NewInt(0).SetUint64(uint64(subPrecision)))

	biNum, _ := BigIntNum{}.NewBigIntPrecision(result, resultPrecision)

	return biNum
}

// SubtractBigIntNums - Receives two 'BigIntNum' instances and proceeds to subtract
// b2 from b2.
//
// 								'minuend' - 'subtrahend' = difference or result
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum. The returned BigIntNum 'result' will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from the
// input parameter, 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractBigIntNums(minuend, subtrahend BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(minuend, subtrahend)

	bigIntNumResult := bSubtract.SubtractPair(bPair)

	return bigIntNumResult
}

// SubtractBigIntNumArray - Receives one BigIntNum which is classified as the 'minuend'.
// The second input parameter is an array of BigIntNum Types labeled, 'subtrahends'.
// The array of 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is an array of BigIntNum Types.
//
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum. The returned BigIntNum subtraction 'result' will contain numeric
// separators (decimal separator, thousands separator and currency symbol) copied
// from the input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractBigIntNumArray(
	minuend BigIntNum,
	subtrahends []BigIntNum) BigIntNum {

	numSeps := minuend.GetNumericSeparatorsDto()

	finalResult := minuend.CopyOut()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult
	}

	for i := 0; i < lenSubtrahends; i++ {

		bPair := BigIntPair{}.NewBigIntNum(finalResult, subtrahends[i])

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)

	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// SubtractBigIntNumOutputToArray - The first input parameter to this method
// is a BigIntNum Type labeled, 'minuend'.  The second input parameter is an
// array of BigIntNum types labeled 'subtrahends'. The 'minuend' is subtracted
// from each element of the 'subtrahends' array with the result output to another
// 'results' array of BigIntNum types which is then returned to the calling
// function.
//
// Example
// =======
// 										    subtrahends										 Output
// Minuend   				    	  Array											   Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
//
// Each of the BigIntNum instances included in the array of BigIntNum subtraction
// results returned by this method, will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractBigIntNumOutputToArray(
	minuend BigIntNum,
	subtrahends []BigIntNum) []BigIntNum {

	numSeps := minuend.GetNumericSeparatorsDto()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return []BigIntNum{}
	}

	resultsArray := make([]BigIntNum, lenSubtrahends)

	for i := 0; i < lenSubtrahends; i++ {

		bPair := BigIntPair{}.NewBigIntNum(minuend, subtrahends[i])

		resultsArray[i] = bSubtract.subtractPairNoNumSeps(bPair)
		_ = resultsArray[i].SetNumericSeparatorsDto(numSeps)
	}

	return resultsArray
}

// SubtractBigIntNumSeries - Receives one BigIntNum which is classified as
// the 'minuend'. The second input parameter, 'subtrahends' is a series of
// Type BigIntNum .
//
// The 'subtrahends' series is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is a series of BigIntNum Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The subtraction result returned by this method as a Type BigIntNum will
// contain numeric separators (decimal separator, thousands separator and currency
// symbol) copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractBigIntNumSeries(
	minuend BigIntNum,
	subtrahends ...BigIntNum) BigIntNum {

	numSeps := minuend.GetNumericSeparatorsDto()

	finalResult := minuend.CopyOut()

	for _, subtrahend := range subtrahends {

		bPair := BigIntPair{}.NewBigIntNum(finalResult, subtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// SubtractDecimal - Performs the subtraction operation on two Decimal Types.
//
// 				decMinuend - decSubtrahend = difference
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum. This resulting BigIntNum instance will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input
// parameter 'decMinuend'.
//
func (bSubtract BigIntMathSubtract) SubtractDecimal(
	decMinuend Decimal,
	decSubtrahend Decimal) (BigIntNum, error) {

	// This method will test the validity of decMinuend and decSubtrahend.
	bPair, err := BigIntPair{}.NewDecimal(decMinuend, decSubtrahend)

	if err != nil {
		ePrefix := "BigIntMathSubtract.SubtractDecimal() "
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewDecimal(decMinuend, decSubtrahend). "+
				"decMinuend='%v' decSubtrahend='%v' Error='%v'",
				decMinuend.GetNumStr(), decSubtrahend.GetNumStr(), err.Error())
	}

	finalResult := bSubtract.SubtractPair(bPair)

	return finalResult, nil
}

// SubtractDecimalArray - Receives one Decimal parameter which is classified as
// the 'minuend'. The second input parameter is an array of Decimal Types
// labeled, 'subtrahends'.
//
// The array of 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is an array of Decimal Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractDecimalArray(
	minuend Decimal,
	subtrahends []Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractDecimalArray() "

	numSeps := minuend.GetNumericSeparatorsDto()

	finalResult, err := minuend.GetBigIntNum()

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by minuend.GetBigIntNum(). Error='%v'",
				err.Error())
	}

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	for i := 0; i < lenSubtrahends; i++ {

		bigINumSubtrahend, err := subtrahends[i].GetBigIntNum()

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by subtrahends[i].GetBigInt(). "+
					"subtrahends[%v]='%v' Error='%v'", i, subtrahends[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(
			finalResult, bigINumSubtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v'", err.Error())
	}

	return finalResult, nil
}

// SubtractDecimalOutputToArray - The first input parameter to this method
// is a Decimal Type labeled, 'minuend'.  The second input parameter is an
// array of Decimal types labeled 'subtrahends'. The 'minuend' is subtracted
// from each element of the 'subtrahends' array with the result output to
// another 'results' array of Decimal types which is then returned to the
// calling function.
//
// Example
// =======
// 										    subtrahends										 Output
// Minuend   				    	   Array											 Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
//
// The array element of the []Decimal 'result' returned by this subtraction
// operation will contain numeric separators (decimal separator, thousands
// separator and currency symbol) copied from input parameter 'minuend'.
//
// The 'result' array ([]Decimal) returned by this subtraction operation will
// contain array elements with numeric separators (decimal separator, thousands
// separator and currency symbol) which have been copied from input parameter
// 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractDecimalOutputToArray(
	minuend Decimal,
	subtrahends []Decimal) ([]Decimal, error) {

	ePrefix := "BigIntMathSubtract.SubtractDecimalOutputToArray() "

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return []Decimal{},
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	bINumMinuend, err := minuend.GetBigIntNum()

	if err != nil {
		return []Decimal{},
			fmt.Errorf(ePrefix+
				"Error returned by minuend.GetBigIntNum(). "+
				"Error='%v'", err.Error())
	}

	resultsArray := make([]Decimal, lenSubtrahends)

	var bigINumSubtrahend, result BigIntNum

	var bPair BigIntPair

	for i := 0; i < lenSubtrahends; i++ {

		bigINumSubtrahend, err = subtrahends[i].GetBigIntNum()

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by subtrahends[i].GetBigIntNum(). "+
					"index='%v' Error='%v'", i, err.Error())
		}

		bPair = BigIntPair{}.NewBigIntNum(bINumMinuend, bigINumSubtrahend)

		result = bSubtract.subtractPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					"index='%v' Error='%v'", i, err.Error())
		}

		resultsArray[i], err = result.GetDecimal()

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetDecimal() "+
					"i='%v' result.Result='%v' Error='%v'. ",
					i, result.GetNumStr(), err.Error())
		}
	}

	return resultsArray, nil
}

// SubtractDecimalSeries - Receives one Decimal Type which is classified as
// the 'minuend'. The second input parameter, 'subtrahends' is a series of
// Type Decimal .
//
// The 'subtrahends' series is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is a series of Decimal Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractDecimalSeries(
	minuend Decimal,
	subtrahends ...Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractDecimalSeries() "

	numSeps := minuend.GetNumericSeparatorsDto()

	finalResult, err := minuend.GetBigIntNum()

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by minuend.GetBigIntNum(). Error='%v'", err.Error())
	}

	if len(subtrahends) == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends series is Empty!")
	}

	var bINumSubtrahend BigIntNum
	var bPair BigIntPair
	for _, subtrahend := range subtrahends {

		bINumSubtrahend, err = subtrahend.GetBigIntNum()

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by subtrahend.GetBigIntNum(). Error='%v'",
					err.Error())
		}

		bPair = BigIntPair{}.NewBigIntNum(finalResult, bINumSubtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// SubtractIntAry - Performs the subtraction operation on two
//
// IntAry Types.
// 				iaMinuend - iaSubtrahend = difference
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'iaMinuend'.
//
func (bSubtract BigIntMathSubtract) SubtractIntAry(
	iaMinuend IntAry,
	iaSubtrahend IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractIntAry() "

	// Method NewIntAry will test validity of iaMinuend and iaSubtrahend
	bPair, err := BigIntPair{}.NewIntAry(iaMinuend, iaSubtrahend)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+"Error returned from iaMinuend.GetBigIntNum(). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bSubtract.SubtractPair(bPair)

	return finalResult, nil
}

// SubtractIntAryArray - Receives one IntAry parameter which is classified as
// the 'minuend'. The second input parameter is an array of IntAry Types
// labeled, 'subtrahends'.
//
// The array of 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is an array of IntAry Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractIntAryArray(
	minuend IntAry,
	subtrahends []IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractInAryArray() "

	// NewIntAry method will test validity of 'minuend'
	finalResult, err := BigIntNum{}.NewIntAry(minuend)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewIntAry(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	for i := 0; i < lenSubtrahends; i++ {

		// Method NewIntAry will test validity of subtrahends[i]
		bigINum, err := BigIntNum{}.NewIntAry(subtrahends[i])

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(subtrahends[i]). "+
					"i='%v' subtrahend='%v' Error='%v'",
					i, subtrahends[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINum)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// SubtractIntAryOutputToArray - The first input parameter to this method
// is a IntAry Type labeled, 'minuend'.  The second input parameter is an
// array of IntAry types labeled 'subtrahends'. The 'minuend' is subtracted
// from each element of the 'subtrahends' array with the result output to
// another 'results' array of IntAry types which is then returned to the
// calling function.
//
// Example
// =======
// 										    subtrahends										 Output
//  Minuend   				    	Array													Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
// Each element in the []IntAry result returned by this subtraction operation will
// contain numeric separators (decimal separator, thousands separator and currency
// symbol) which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractIntAryOutputToArray(
	minuend IntAry,
	subtrahends []IntAry) ([]IntAry, error) {

	ePrefix := "BigIntMathSubtract.SubtractIntAryOutputToArray() "

	// NewIntAry will test validity of 'minuend'
	bINumMinuend, err := BigIntNum{}.NewIntAry(minuend)

	if err != nil {
		return []IntAry{},
			fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(minuend). "+
				"Error='%v' ", err.Error())
	}

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return []IntAry{},
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	resultsArray := make([]IntAry, lenSubtrahends)

	for i := 0; i < lenSubtrahends; i++ {

		// Method NewIntAry will test validity of subtrahends[i]
		bINumSubtrahend, err := BigIntNum{}.NewIntAry(subtrahends[i])

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(subtrahends[%v]). "+
					"Error='%v' \n", i, err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(bINumMinuend, bINumSubtrahend)

		result := bSubtract.subtractPairNoNumSeps(bPair)

		resultsArray[i], err = result.GetIntAry()

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetIntAryElements() "+
					"i='%v' result.Result='%v' Error='%v'. ",
					i, result.GetNumStr(), err.Error())
		}

		err = resultsArray[i].SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by resultsArray[i].SetNumericSeparatorsDto(numSeps). "+
					"Error='%v' \n", err.Error())
		}

	}

	return resultsArray, nil
}

// SubtractIntArySeries - Receives one IntAry Type which is classified as
// the 'minuend'. The second input parameter, 'subtrahends' is a series of
// Type IntAry .
//
// The 'subtrahends' series is subtracted from the 'minuend'.
//
// In the subtraction operation:
// 								b1 - b2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractIntArySeries(
	minuend IntAry,
	subtrahends ...IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractIntArySeries() "

	// Method New IntAry will test the validity of minuend
	finalResult, err := BigIntNum{}.NewIntAry(minuend)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewIntAry(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	if len(subtrahends) == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends series is Empty!")
	}

	for i, subtrahend := range subtrahends {

		// Method NewIntAry will test the validity of IntAry subtrahend
		bigINumSubtrahend, err := BigIntNum{}.NewIntAry(subtrahend)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(subtrahend). "+
					"index='%v' subtrahend='%v' Error='%v'",
					i, subtrahend.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINumSubtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// SubtractINumMgr - Receives two objects which implement the INumMgr Interface
// and subtracts their numeric values.
//
// The 'subtrahend' numeric value is subtracted from the 'minuend' numeric value.
//
// In the subtraction operation:
//
//								'minuend' - 'subtrahend' = difference or result
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractINumMgr(
	minuend,
	subtrahend INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractINumMgr() "

	bPair, err := BigIntPair{}.NewINumMgr(minuend, subtrahend)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by NBigIntPair{}.NewINumMgr(minuend, subtrahend). "+
				"minuend.GetNumStr()='%v', subtrahend.GetNumStr()='%v' Error='%v' ",
				minuend.GetNumStr(), subtrahend.GetNumStr(), err.Error())
	}

	finalResult := bSubtract.SubtractPair(bPair)

	return finalResult, nil
}

// SubtractINumMgrArray - Receives two input parameters. The first parameter
// is an INumMgr instance which is classified as the 'minuend'. The second
// parameter is an array of INumMgr instances which resents the 'subtrahends'.
// A subtraction operation is performed on the 'minuend' and the 'subtrahends'
// array. The numeric value of the 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
//
//								'minuend' - 'subtrahend' = difference or result
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry. This allows the user to mix different types in
// a single array and add their numeric values.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractINumMgrArray(
	minuend INumMgr,
	subtrahends []INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractINumMgrArray() "

	finalResult, err := BigIntNum{}.NewINumMgr(minuend)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewINumMgr(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	for i := 0; i < lenSubtrahends; i++ {

		bPair, err := BigIntPair{}.NewINumMgr(&finalResult, subtrahends[i])

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(&finalResult.Result, &subtrahends[i]). "+
					" i='%v' subtrahends[i].GetNumStr()='%v' Error='%v' ", i, subtrahends[i].GetNumStr(), err.Error())
		}

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())

	}

	return finalResult, nil
}

// SubtractINumMgrOutputToArray - The first input parameter to this method
// is an object which implements the INumMgr interface ('minuend').  The
// second input parameter is an array of INumMgr interface types labeled
// 'subtrahends'. The 'minuend' is subtracted from each element of the
// 'subtrahends' array with the result output to another 'results' array
// of INumMgr interface types which is then returned to the calling function.
//
// Example
// =======
// 										    subtrahends										 Output
//  Minuend   				    	Array													Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
// Note: The underlying type for the returned results array is 'BigIntNum', a
// type which implements the INumMgr Interface.
//
//
// Each element of the []INumMgr 'result' array returned by this subtraction
// operation will contain numeric separators (decimal separator, thousands
// separator and currency symbol) which were copied from input parameter
// 'minuend'.
//

func (bSubtract BigIntMathSubtract) SubtractINumMgrOutputToArray(
	minuend INumMgr,
	subtrahends []INumMgr) ([]INumMgr, error) {

	ePrefix := "BigIntMathSubtract.SubtractINumMgrOutputToArray() "

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return []INumMgr{},
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	resultsArray := make([]INumMgr, lenSubtrahends)

	for i := 0; i < lenSubtrahends; i++ {

		bPair, err := BigIntPair{}.NewINumMgr(minuend, subtrahends[i])

		if err != nil {
			return []INumMgr{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(minuend, subtrahends[i]) "+
					"minuend='%v' subtrahends[%v]='%v' Error='%v'. ",
					minuend.GetNumStr(), i, subtrahends[i].GetNumStr(), err.Error())
		}

		result := bSubtract.subtractPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []INumMgr{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					"Index='%v' Error='%v' \n", i, err.Error())
		}

		resultsArray[i] = &result

	}

	return resultsArray, nil
}

// SubtractINumMgrSeries - Receives two input parameters. The first parameter
// is an INumMgr instance which is classified as the 'minuend'. The second
// parameter is a series of INumMgr instances which resents the 'subtrahends'.
// A subtraction operation is performed on the 'minuend' and the 'subtrahends'.
// The numeric value of the 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
//
//					'minuend' - 'subtrahend' = difference or result
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry. This allows the user to mix different types in
// a single array and add their numeric values.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractINumMgrSeries(
	minuend INumMgr,
	subtrahends ...INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractINumMgrSeries() "

	finalResult, err := BigIntNum{}.NewINumMgr(minuend)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewINumMgr(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	if len(subtrahends) == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends series is Empty!")
	}

	for i, subtrahend := range subtrahends {

		bPair, err := BigIntPair{}.NewINumMgr(&finalResult, subtrahend)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(&finalResult.Result, &subtrahend). "+
					" i='%v' subtrahend.GetNumStr()='%v' Error='%v' ",
					i, subtrahend.GetNumStr(), err.Error())
		}

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v'\n", err.Error())
	}

	return finalResult, nil
}

// SubtractNumStr - Receives two number strings and proceeds to subtract
// n2 from n1.
//
// The 'n1' and 'n2' strings passed to this method are number strings in that they
// consist of a string of numeric digits representing a numeric value. Number strings
// may contain leading minus signs (-) to indicate a negative numeric value. The
// numeric digits string may also include a delimiting decimal separator to identify
// fractional digits. The number strings are parsed based on the decimal separator
// character specified by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings 'n1NumStr' and 'n2NumStr'. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// subtraction operation.
//
// In the subtraction operation:
// 								n1 - n2 = difference or result
//								'minuend' - 'subtrahend' = difference or result
//								n1 = 'minuend'
//								n2 = 'subtrahend'
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// specified by input parameter, 'numSeps'.
//
func (bSubtract BigIntMathSubtract) SubtractNumStr(
	n1 string,
	n2 string,
	numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStr() "

	numSeps.SetDefaultsIfEmpty()

	bPair, err := BigIntPair{}.NewNumStrWithNumSeps(n1, n2, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewNumStr(n1, n2). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bSubtract.SubtractPair(bPair)

	return finalResult, nil
}

// SubtractNumStrArray - Receives one string parameter which is classified as
// the 'minuend'. The second input parameter is an array of strings labeled,
// 'subtrahends'.
//
// The array of 'subtrahends' is subtracted from the 'minuend' and the summary
// result of this transaction operation is returned as a BigIntNum type.
//
// The 'minuend' string and 'subtrahends' string array are passed to this method
// as number strings. As such, each number string consists of a string of numeric
// digits representing a numeric value. Number strings may include a leading minus
// sign indicating a negative numeric value. These strings of numeric digits may
// also include a delimiting decimal separator to identify fractional digits. The
// number strings are parsed based on the decimal separator character specified by
// input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings 'minuend' and 'subtrahends'. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// subtraction operation.
//
// In the subtraction operation:
//								'minuend' - 'subtrahend' = difference or result
//
// In this method, the 'subtrahend' is an array of string Types.
//
// After all subtraction operations have been completed, the summary 'difference'
// or 'result' is returned as a Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// specified by input parameter, 'numSeps'.
//
func (bSubtract BigIntMathSubtract) SubtractNumStrArray(
	minuend string,
	subtrahends []string,
	numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrArray() "

	numSeps.SetDefaultsIfEmpty()

	finalResult, err := BigIntNum{}.NewNumStrWithNumSeps(minuend, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(minuend, numSeps). "+
				"minuend='%v' numSeps='%v' Error='%v'",
				minuend, numSeps.String(), err.Error())
	}

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	for i := 0; i < lenSubtrahends; i++ {

		bigINum, err := BigIntNum{}.NewNumStrWithNumSeps(subtrahends[i], numSeps)

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(subtrahends[i], numSeps). "+
					"i='%v' subtrahend='%v' Error='%v'",
					i, subtrahends[i], err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINum)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"numSeps='%v' Error='%v'",
				numSeps.String(), err.Error())
	}

	return finalResult, nil
}

// SubtractNumStrOutputToArray - The first input parameter to this method
// is a string Type labeled, 'minuend'.  The second input parameter is an
// array of strings labeled 'subtrahends'. The 'minuend' is subtracted
// from each element of the 'subtrahends' array with the result output to
// another 'results' array of strings which is then returned to the calling
// function.
//
// The 'minuend' string and 'subtrahends' string array are passed to this method
// as number strings. Each number string consists of a string of numeric digits
// representing a numeric value. Number strings may include a leading minus sign
// (-) indicating a negative numeric value. Numeric strings may also include a
// delimiting decimal separator to identify fractional digits. These number
// strings are parsed based on the decimal separator character specified by input
// parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings 'minuend' and 'subtrahends'. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. 'numSeps'
// is also used in configuring the return value for this subtraction operation.
//
// Example
// =======
// 										    subtrahends										 Output
//  Minuend   				    	Array													Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
//
func (bSubtract BigIntMathSubtract) SubtractNumStrOutputToArray(
	minuend string,
	subtrahends []string,
	numSeps NumericSeparatorDto) ([]string, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrOutputToArray() "

	numSeps.SetDefaultsIfEmpty()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		return []string{},
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	resultsArray := make([]string, lenSubtrahends)

	for i := 0; i < lenSubtrahends; i++ {

		bPair, err := BigIntPair{}.NewNumStrWithNumSeps(minuend, subtrahends[i], numSeps)

		if err != nil {
			return []string{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewNumStrWithNumSeps(minuend, subtrahends[i], numSeps) "+
					"minuend='%v' subtrahends[%v]='%v' numSeps='%v' Error='%v'. ",
					minuend, i, subtrahends[i], numSeps.String(), err.Error())
		}

		result := bSubtract.SubtractPair(bPair)

		resultsArray[i] = result.GetNumStr()

		if err != nil {
			return []string{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetNumStr() "+
					"i='%v' result.Result='%v' Error='%v'. ",
					i, result.GetNumStr(), err.Error())
		}
	}

	return resultsArray, nil
}

// SubtractNumStrSeries - Receives one 'string' Type which is classified as
// the 'minuend'. The second input parameter, 'subtrahends' is a series of
// 'strings' .
//
// The 'minuend' string and 'subtrahends' string series are passed to this method
// as number strings. Each number string consists of a string of numeric digits
// representing a numeric value. Number strings may include a leading minus sign
// (-) indicating a negative numeric value. Number strings may also include a
// delimiting decimal separator to identify fractional digits. These number
// strings are parsed based on the decimal separator character specified by input
// parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings contained in the 'minuend' and 'subtrahends' input
// parameters. 'numSeps' represents the applicable decimal separator, thousands
// separator and currency symbol. 'numSeps' is also used in configuring the return
// value for this subtraction operation.
//
// The 'subtrahends' series is subtracted from the 'minuend' to produce a
// cumulative result.
//
// In the subtraction operation:
//
//								'minuend' - 'subtrahend' = difference or result
//
// In this method, the 'subtrahend' is a series of strings.
//
// After the subtraction operation, the cumulative 'difference' or 'result' is returned
// as a Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// specified by the input parameter, 'numSeps'.
//
func (bSubtract BigIntMathSubtract) SubtractNumStrSeries(
	numSeps NumericSeparatorDto,
	minuend string,
	subtrahends ...string) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrSeries() "

	numSeps.SetDefaultsIfEmpty()

	finalResult, err := BigIntNum{}.NewNumStrWithNumSeps(minuend, numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(minuend, numSeps). "+
				"minuend='%v' numSeps='%v' Error='%v'", minuend, numSeps.String(), err.Error())
	}

	if len(subtrahends) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends series is Empty!")
	}

	for i, subtrahend := range subtrahends {

		bigINumSubtrahend, err := BigIntNum{}.NewNumStrWithNumSeps(subtrahend, numSeps)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(subtrahend, numSeps). "+
					"index='%v' subtrahend='%v' numSeps='%v' Error='%v'",
					i, subtrahend, numSeps.String(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINumSubtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"numSeps='%v' Error='%v'", numSeps.String(), err.Error())
	}

	return finalResult, nil
}

// SubtractNumStrDto - Receives two 'NumStrDto' instances and proceeds to subtract
// 'nDtoSubtrahend' from 'nDtoMinuend'.
//
// 			nDtoMinuend - nDtoSubtrahend = difference or result
//
// In the subtraction operation:
//							  'minuend' - 'subtrahend' = difference or result
//								nDtoMinuend 		= 'minuend'
//								nDtoSubtrahend 	= 'subtrahend'
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'nDtoMinuend'.
//

func (bSubtract BigIntMathSubtract) SubtractNumStrDto(
	nDtoMinuend NumStrDto,
	nDtoSubtrahend NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrDto() "

	err := nDtoMinuend.IsValid(ePrefix + "'nDtoMinuend' INVALID! ")

	if err != nil {
		return BigIntNum{}.New(), err
	}

	err = nDtoSubtrahend.IsValid(ePrefix + "'nDtoSubtrahend' INVALID! ")

	if err != nil {
		return BigIntNum{}.New(), err
	}

	bPair, err := BigIntPair{}.NewNumStrDto(nDtoMinuend, nDtoSubtrahend)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewNumStrDto(nDtoMinuend, nDtoSubtrahend). "+
				"Error='%v' ", err.Error())
	}

	numSeps := nDtoMinuend.GetNumericSeparatorsDto()

	finalResult := bSubtract.subtractPairNoNumSeps(bPair)

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return finalResult, nil
}

// SubtractNumStrDtoArray - Receives one NumStrDto parameter which is classified as
// the 'minuend'. The second input parameter is an array of NumStrDto Types
// labeled, 'subtrahends'.
//
// The array of 'subtrahends' is subtracted from the 'minuend'.
//
// In the subtraction operation:
//								'minuend' - 'subtrahend' = difference or result
// 								b1 - b2 = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is an array of NumStrDto Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractNumStrDtoArray(
	minuend NumStrDto,
	subtrahends []NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrDtoArray() "

	err := minuend.IsValid(ePrefix + "'minuend' INVALID! ")

	if err != nil {
		return BigIntNum{}.New(), err
	}

	finalResult, err := BigIntNum{}.NewNumStrDto(minuend)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrDto(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	lenSubtrahends := len(subtrahends)

	if lenSubtrahends == 0 {
		_ = finalResult.SetNumericSeparatorsDto(numSeps)
		return finalResult,
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	for i := 0; i < lenSubtrahends; i++ {

		err := subtrahends[i].IsValid(ePrefix +
			fmt.Sprintf("subtrahends[%v] INVALID!", i))

		bigINum, err := BigIntNum{}.NewNumStrDto(subtrahends[i])

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(subtrahends[i]). "+
					"i='%v' subtrahend='%v' Error='%v'",
					i, subtrahends[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINum)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// SubtractNumStrDtoOutputToArray - The first input parameter to this method
// is a NumStrDto Type labeled, 'minuend'.  The second input parameter is an
// array of NumStrDto types labeled 'subtrahends'. The 'minuend' is subtracted
// from each element of the 'subtrahends' array with the result output to
// another 'results' array of NumStrDto types which is then returned to the
// calling function.
//
// Example
// =======
// 										    subtrahends										 Output
// Minuend   				    	   Array											 Array
//
//		10			-					subtrahends[0] = 2			=				  outputarray[0] =  8
//		10			-					subtrahends[1] = 3			=				  outputarray[1] =  7
//		10			-					subtrahends[2] = 4			=				  outputarray[2] =  6
//		10			-					subtrahends[3] = 5			=				  outputarray[3] =  5
//		10			-					subtrahends[4] = 6			=				  outputarray[4] =  4
//		10			-					subtrahends[5] = 9			=				  outputarray[5] =  1
//
//
// Each array element in the []NumStrDto 'result' returned by this subtraction
// operation will contain numeric separators (decimal separator, thousands
// separator and currency symbol) which were copied from input parameter
// 'minuend'.
//

func (bSubtract BigIntMathSubtract) SubtractNumStrDtoOutputToArray(
	minuend NumStrDto,
	subtrahends []NumStrDto) ([]NumStrDto, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrDtoOutputToArray() "

	lenSubtrahends := len(subtrahends)

	err := minuend.IsValid(ePrefix + "'minuend' INVALID! ")

	if err != nil {
		return []NumStrDto{}, err
	}

	if lenSubtrahends == 0 {
		return []NumStrDto{},
			errors.New(ePrefix + "Error: subtrahends array is Empty!")
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	resultsArray := make([]NumStrDto, lenSubtrahends)

	for i := 0; i < lenSubtrahends; i++ {

		err := subtrahends[i].IsValid(ePrefix +
			fmt.Sprintf("subtrahends[%v] INVALID! ", i))

		bPair, err := BigIntPair{}.NewNumStrDto(minuend, subtrahends[i])

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewNumStrDto(minuend, subtrahends[i]) "+
					"minuend='%v' subtrahends[%v]='%v' Error='%v'. ",
					minuend.GetNumStr(), i, subtrahends[i].GetNumStr(), err.Error())
		}

		result := bSubtract.subtractPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					"Error='%v' \n", err.Error())
		}

		resultsArray[i], err = result.GetNumStrDto()

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetNumStrDto() "+
					"i='%v' result.Result='%v' Error='%v'. ",
					i, result.GetNumStr(), err.Error())
		}
	}

	return resultsArray, nil
}

// SubtractNumStrDtoSeries - Receives one NumStrDto Type which is classified as
// the 'minuend'. The second input parameter, 'subtrahends' is a series of
// Type NumStrDto .
//
// The 'subtrahends' series is subtracted from the 'minuend'.
//
// In the subtraction operation:
//								'minuend' - 'subtrahend' = difference or result
// 								b1 - b2 = difference or result
//								b1 = 'minuend'
//								b2 = 'subtrahend'
//
// In this method, the 'subtrahend' is a series of NumStrDto Types.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'minuend'.
//
func (bSubtract BigIntMathSubtract) SubtractNumStrDtoSeries(
	minuend NumStrDto,
	subtrahends ...NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathSubtract.SubtractNumStrDtoSeries() "

	err := minuend.IsValid(ePrefix + "'minuend' INVALID! ")

	if err != nil {
		return BigIntNum{}.New(), err
	}

	finalResult, err := BigIntNum{}.NewNumStrDto(minuend)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrDto(minuend). "+
				"minuend='%v' Error='%v'", minuend.GetNumStr(), err.Error())
	}

	numSeps := minuend.GetNumericSeparatorsDto()

	if len(subtrahends) == 0 {

		_ = finalResult.SetNumericSeparatorsDto(numSeps)

		return finalResult,
			errors.New(ePrefix + "Error: subtrahends series is Empty!")
	}

	for i, subtrahend := range subtrahends {

		err := subtrahend.IsValid(ePrefix +
			fmt.Sprintf("subtrahend index='%v' INVALID! ", i))

		bigINumSubtrahend, err := BigIntNum{}.NewNumStrDto(subtrahend)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(subtrahend). "+
					"index='%v' subtrahend='%v' Error='%v'",
					i, subtrahend.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINumSubtrahend)

		finalResult = bSubtract.subtractPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult, nil
}

// SubtractPair - Performs the subtraction operation. This method receives a type
// 'BigIntPair' and proceeds to subtract bPair.Big2 from bPair.Big1.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter bPair.Big1, the minuend.
//
func (bSubtract BigIntMathSubtract) SubtractPair(bPair BigIntPair) BigIntNum {

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	finalResult := bSubtract.subtractPairNoNumSeps(bPair)

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// subtractPairNoNumSeps - Performs the subtraction operation. This method receives a type
// 'BigIntPair' and proceeds to subtract bPair.Big1 from bPair.Big2.
//
// After the subtraction operation, the 'difference' or 'result' is returned as a
// Type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// default numeric separators (decimal separator, thousands separator and currency
// symbol).
//
func (bSubtract BigIntMathSubtract) subtractPairNoNumSeps(bPair BigIntPair) BigIntNum {

	bPair.MakePrecisionsEqual()

	b3 := big.NewInt(0).Sub(bPair.GetBig1BigInt(), bPair.GetBig2BigInt())

	bResult := BigIntNum{}.NewBigInt(b3, bPair.Big2.GetPrecisionUint())

	bResult.SetNumericSeparatorsToDefaultIfEmpty()

	return bResult
}
