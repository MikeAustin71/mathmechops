package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

// BigIntMathDivide - This type is comprised of methods used to perform the
// division operation using *big.Int numeric types.
//
// Reference the 'big' math package: https://golang.org/pkg/math/big/
type BigIntMathDivide struct {
	Input BigIntPair
	// BigIntPair.Big1 = Dividend
	// BigIntPair.Big2 = Divisor

	Result BigIntNum // BigIntPair.Big1 = Quotient
	// BigIntPair.Big2 = Modulo

	ResultFracQuo BigIntNum // Quotient expressed with fractional digits
	// to the right of the decimal place.
}

// BigIntDividedByTwoToPower - Performs integer division by two using
// a 'right-shift' technique. Remainders from this division operation
// are discarded, only the integer quotient is returned.
//
// Example:
// ========
//
//	quotient =  dividend / 2^(exponent)
//
// In the example of 33,333 / 2^8:
//
//			(1) The fractional quotient of 33,333/256 (or 2^8) is 130.
//
//	   (2) This method will use a right shift technique 33,333 / 2^(8) to generate
//	   		a quotient of 130.
func (bIDivide BigIntMathDivide) BigIntDividedByTwoToPower(
	dividend *big.Int,
	exponent uint) (integerQuotient *big.Int) {

	integerQuotient = big.NewInt(0)

	integerQuotient.Rsh(dividend, exponent)

	return integerQuotient
}

// BigIntFracQuotient - Performs a division on integers of type *big.Int.
// The result is returned a as a type *big.Int with an accompanying precision
// specification. Taken together, the returned *big.Int quotient and precision
// specification describe a floating point numeric value with a a fixed number
// of digits after the decimal place. The input parameter, 'maxPrecision' is
// used to configure the maximum number of fractional digits to the right of
// the decimal place in the returned quotient.
//
// Example
// =======
// This division operation will produce a quotient which may include a fixed
// number of fractional digits to the right of the decimal place:
//
//	quotient = dividend / divisor
//
// For the example
//
//	quotient =	752.314 / 21.67894
//
// 'dividend' and 'divisor' would be configured as follows:
//
//										dividend 						= 752314
//	                 dividendPrecision		= 3
//	                 divisor 						= 2167894
//	                 divisorPrecision    = 5
//
// Assuming a 'maxPrecision' value of '30', the quotient would be
// calculated as follows:
//
//										quotient 						= 34702526968569496479071393712054
//	                 quotientPrecision   = 30
//
// Input Parameters
// ================
//
// dividend				*big.Int	- The 'dividend' value will be divided by the 'divisor'
//
//	to produce a 'quotient'.
//
// dividendPrecision	uint	- This unsigned integer value is a precision specification
//
//	associated with input parameter 'dividend'. The precision
//	specifies the number if digits to right of the decimal
//	place in the series of integer digits contained in
//	'dividend'.
//
// divisor				*big.Int	- The 'dividend' value will be divided by the 'divisor'
//
//	to produce a 'quotient'.
//
// divisorPrecision		uint	- This unsigned integer value is a precision specification
//
//	associated with input parameter 'divisor'. The precision
//	specifies the number if digits to right of the decimal
//	place in the series of integer digits contained in
//	'divisor'.
//
// 'maxPrecision' 		uint	-	Maximum precision specifies the maximum number of
//
//	decimal digits to which the result or 'quotient'
//	will calculated and returned to the caller. The
//	quotient may consist of actual fractional digits
//	which number less than 'maxPrecision'. However, if
//	the number of digits to the right of the decimal
//	place exceeds 'maxPrecision', the returned quotient
//	will be rounded to 'maxPrecision' fractional digits
//	to the right of the decimal place.
//
// Return Values
// =============
//
// quotient				*big.Int	- The result of the division operation expressed
//
//	as an integer.
//
// quotientPrecision	uint	- An unsigned integer which specifies the number of
//
//	fractional digits to the right of the decimal place
//	in the series of integer digits defined by 'quotient'
//
// err							 error	- If an error is encountered, this function will
//
//	                           return a quotient set equal to zero and an error
//	                           object will be returned containing an appropriate
//															error message. If the function completes the division
//															operation successfully, the returned 'quotient' will
//	                           be populated with the correct result and 'err' will
//															be set equal to 'nil'.
func (bIDivide BigIntMathDivide) BigIntFracQuotient(
	dividend,
	dividendPrecision,
	divisor,
	divisorPrecision,
	maxPrecision *big.Int) (quotient *big.Int, quotientPrecision *big.Int, err error) {

	ePrefix := "BigIntMathDivide.BigIntFracQuotient() "

	quotient = big.NewInt(0)
	quotientPrecision = big.NewInt(0)
	err = nil

	if dividend == nil {
		dividend = big.NewInt(0)
	}

	if divisor == nil {
		divisor = big.NewInt(0)
	}

	if dividendPrecision == nil {
		divisorPrecision = big.NewInt(0)
	}

	if divisorPrecision == nil {
		divisorPrecision = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if divisor.Cmp(bigZero) == 0 {
		err = fmt.Errorf(ePrefix + "Error - Divide by ZERO! Input parameter 'divisor' is ZERO!")
		return quotient, quotientPrecision, err
	}

	if dividend.Cmp(bigZero) == 0 {
		err = nil
		return quotient, quotientPrecision, err
	}

	// Prepare divisor
	// Setting to absolute value of divisor
	denomnatrDivsr := big.NewInt(0).Set(divisor)
	denomnatrDivsrSign := int64(1)

	if divisor.Cmp(bigZero) == -1 {
		denomnatrDivsrSign = -1
		denomnatrDivsr.Mul(denomnatrDivsr, big.NewInt(denomnatrDivsrSign))
	}

	// Prepare dividend
	// Setting to absolute value of dividend
	bigTen := big.NewInt(10)

	numratrDivdnd := big.NewInt(0).Set(dividend)
	numratrDivdndSign := int64(1)

	if dividend.Cmp(bigZero) == -1 {
		numratrDivdndSign = -1
		numratrDivdnd.Mul(numratrDivdnd, big.NewInt(numratrDivdndSign))
	}

	denomnatrDivsrShift := big.NewInt(0).Set(divisorPrecision)
	numratrDivdndShift := big.NewInt(0).Set(dividendPrecision)
	scale := big.NewInt(0)

	if numratrDivdndShift.Cmp(bigZero) == 1 {
		denomnatrDivsrShift.Sub(denomnatrDivsrShift, numratrDivdndShift)
		if denomnatrDivsrShift.Cmp(bigZero) == -1 {
			scale = big.NewInt(0).Exp(bigTen, big.NewInt(0).Mul(denomnatrDivsrShift, big.NewInt(-1)), nil)
			denomnatrDivsr.Mul(denomnatrDivsr, scale)
		}

		numratrDivdndShift = big.NewInt(0)
	}

	if denomnatrDivsrShift.Cmp(bigZero) == 1 {
		numratrDivdndShift.Sub(numratrDivdndShift, denomnatrDivsrShift)

		if numratrDivdndShift.Cmp(bigZero) == -1 {
			scale = big.NewInt(0).Exp(bigTen, big.NewInt(0).Mul(numratrDivdndShift, big.NewInt(-1)), nil)
			numratrDivdnd.Mul(numratrDivdnd, scale)
		}
	}

	// Do integer division
	scratch := big.NewInt(0)
	intQuotient, intRemndr := big.NewInt(0).QuoRem(numratrDivdnd, denomnatrDivsr, scratch)

	quotient = big.NewInt(0).Set(intQuotient)

	// Calculate fractional digits out to maxPrecision + 1
	bigOne := big.NewInt(1)
	iMaxPrecision := big.NewInt(0).Set(maxPrecision)
	iMaxPrecision.Add(iMaxPrecision, bigOne)
	iCnt := big.NewInt(0)
	lastNonZeroDigitIdx := big.NewInt(-1)

	for iCnt.Cmp(iMaxPrecision) == -1 {

		numratrDivdnd = big.NewInt(0).Mul(intRemndr, bigTen)
		intQuotient, intRemndr = big.NewInt(0).QuoRem(numratrDivdnd, denomnatrDivsr, scratch)

		if intQuotient.Cmp(bigZero) == 1 {
			lastNonZeroDigitIdx.Set(iCnt)
		}

		quotient.Mul(quotient, bigTen)
		quotient.Add(quotient, intQuotient)

		iCnt.Add(iCnt, bigOne)
	}

	if lastNonZeroDigitIdx.Cmp(big.NewInt(-1)) == 0 {

		scale = big.NewInt(0).Exp(bigTen, iMaxPrecision, nil)
		quotient.Quo(quotient, scale)
		quotientPrecision = big.NewInt(0)

	} else if lastNonZeroDigitIdx.Cmp(big.NewInt(0).Sub(iMaxPrecision, bigOne)) == -1 {
		// else if lastNonZeroDigitIdx < (i64MaxPrecision - 1)
		factor := big.NewInt(0).Sub(iMaxPrecision, lastNonZeroDigitIdx)
		factor.Sub(factor, bigOne)
		scale =
			big.NewInt(0).Exp(bigTen, factor, nil)

		quotient.Quo(quotient, scale)

		// uint(lastNonZeroDigitIdx + 1)
		quotientPrecision = big.NewInt(0).Add(lastNonZeroDigitIdx, bigOne)

	} else {

		//fmt.Println("before quotient: ", quotient.Text(10))
		quotient.Add(quotient, big.NewInt(5))

		quotient.Quo(quotient, bigTen)
		//fmt.Println("after quotient: ", quotient.Text(10))
		iMaxPrecision.Sub(iMaxPrecision, bigOne)

		// = uint(i64MaxPrecision)
		quotientPrecision.Set(iMaxPrecision)

	}

	if numratrDivdndSign != denomnatrDivsrSign {
		quotient.Mul(quotient, big.NewInt(-1))
	}

	err = nil

	return quotient, quotientPrecision, err
}

// BigIntNumQuotientMod - Performs a division operation on BigIntNum input
// parameters 'dividend' and 'divisor'.
//
// There are two BigIntNum return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumQuotientMod(
	dividend,
	divisor BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumQuotientMod() "

	if divisor.IsZero() {
		quotient = BigIntNum{}.NewBigInt(big.NewInt(0), 0)
		modulo = BigIntNum{}.NewBigInt(big.NewInt(0), 0)
		err = fmt.Errorf(ePrefix + "Error: Attempted to divide by zero!")
		return quotient, modulo, err
	}

	bPair := BigIntPair{}.NewBigIntNum(dividend, divisor)

	bPair.MaxPrecision = maxPrecision

	var err2 error

	quotient, modulo, err2 = BigIntMathDivide{}.PairQuotientMod(bPair)

	if err2 != nil {
		quotient = BigIntNum{}.NewBigInt(big.NewInt(0), 0)
		modulo = BigIntNum{}.NewBigInt(big.NewInt(0), 0)
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}.PairQuotientMod(bPair) "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, err2.Error())

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// BigIntNumIntQuotient - Performs integer division on two BigIntNum types passed as
// input parameters, 'dividend' and 'divisor'.
//
// The division operation performed by this method is T-Division or truncated division.
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// After completing the division operation, an integer quotient of type BigIntNum
// is returned.
//
// Examples:
// =========
//
//																					Return Value
//	 Divisor	divided by	Dividend		=		Integer Quotient
//			 5 				/ 				 2 				= 				 2
//	    5.25 		/ 				 2  			= 				 2
//	    2 				/ 				 4				= 				 0
//			-5 				/ 				 2 				= 				-2
//	   -5.25     /    			 2  			= 				-2
//	   -2 				/ 				 4				= 				 0
//			 5 				/ 				-2 				=					-2
//	    5.25 		/ 				-2 				= 				-2
//	    2 				/ 				-4				= 				 0
//			-5 				/ 				-2 				= 				 2
//	   -5.25     /    			-2 				= 				 2
//	   -2 				/ 				-4				= 				 0
//	    12.555		/ 			  -2.5			=			    -5
//	   -12.555		/ 			  -2.5			=			     5
//	    12.555		/ 			  -2				=			    -6
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumIntQuotient(
	dividend,
	divisor BigIntNum) (intQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumIntQuotient() "
	intQuotient = BigIntNum{}.NewBigInt(big.NewInt(0), 0)

	if divisor.IsZero() {
		err = fmt.Errorf(ePrefix + "Error: Attempted to divide by zero!")
		return intQuotient, err
	}

	numSeps := dividend.GetNumericSeparatorsDto()

	bPair := BigIntPair{}.NewBigIntNum(dividend, divisor)

	var errx error

	intQuotient, errx = BigIntMathDivide{}.PairIntQuotient(bPair)

	if errx != nil {
		intQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}.PairIntQuotient(bPair). "+
			"dividend='%v' divisor='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), errx.Error())

		return intQuotient, err
	}

	errx = intQuotient.SetNumericSeparatorsDto(numSeps)

	if errx != nil {
		intQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by intQuotient.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v'", errx.Error())

		return intQuotient, err
	}

	err = nil

	return intQuotient, err
}

// BigIntNumModulo - Performs a modulo operation on BigIntNum input
// parameters 'dividend' and 'divisor'.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one BigIntNum value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division 'result' (modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumModulo(
	dividend,
	divisor BigIntNum,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumQuotientMod() "

	if divisor.IsZero() {
		modulo = BigIntNum{}.NewBigInt(big.NewInt(0), 0)
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair := BigIntPair{}.NewBigIntNum(dividend, divisor)

	bPair.MaxPrecision = maxPrecision

	var errx error

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// BigIntNumFracQuotient - Performs a division operation on BigIntNum input
// parameters 'dividend' and 'divisor'.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits
// as appropriate. Remember that the BigIntNum type specifies 'precision'. Precision
// is defined as the number of fractional digits to the right of the decimal place.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					   Quotient
//	 Dividend		divided by	Divisor		=		BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumFracQuotient(
	dividend,
	divisor BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumFracQuotient() "

	fracQuotient = new(BigIntNum).NewZero(0)
	err = nil

	binDividend, errx := dividend.GetBigInt()

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return fracQuotient, err
	}

	binDividendPrecision := dividend.GetPrecisionBigInt()

	binDivisor, errx := divisor.GetBigInt()

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return fracQuotient, err
	}

	binDivisorPrecision := divisor.GetPrecisionBigInt()

	biMaxPrecision := big.NewInt(0).SetUint64(uint64(maxPrecision))

	fracQuo, fracQuoPrecision, errx := BigIntMathDivide{}.BigIntFracQuotient(
		binDividend,
		binDividendPrecision,
		binDivisor,
		binDivisorPrecision,
		biMaxPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return fracQuotient, err
	}

	fracQuotient.SetBigInt(fracQuo, uint(fracQuoPrecision.Uint64()))
	fracQuotient.SetNumericSeparatorsDto(dividend.GetNumericSeparatorsDto())

	err = nil
	return fracQuotient, nil
}

// BigIntNumFracQuotientArray - Performs a division operation on BigIntNum input
// parameters 'dividends' and 'divisor'. 'dividends' is an array of BigIntNum types
// and the division is operation is performed on each element of the array using a
// single 'divisor'.
//
// The resulting quotients are returned as an array BigIntNum types representing the
// result of each division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotients. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// Each element of the returned BigIntNum array resulting from this division operation
// will contain contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from the first element of the input parameter 'dividends'
// array.
func (bIDivide BigIntMathDivide) BigIntNumFracQuotientArray(
	dividends []BigIntNum,
	divisor BigIntNum,
	maxPrecision uint) (fracQuoArray []BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumFracQuotientArray() "

	if divisor.IsZero() {
		fracQuoArray = []BigIntNum{}
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuoArray, err
	}

	lenAry := len(dividends)

	if lenAry == 0 {
		fracQuoArray = []BigIntNum{}
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	numSeps := dividends[0].GetNumericSeparatorsDto()

	fracQuoArray = make([]BigIntNum, lenAry, lenAry+20)
	var errx error

	for i := 0; i < lenAry; i++ {

		if i > 0 {
			errx = dividends[i].SetNumericSeparatorsDto(numSeps)

			if errx != nil {
				fracQuoArray = []BigIntNum{}

				err = fmt.Errorf(ePrefix+
					"Error returned by dividends[i].SetNumericSeparatorsDto(numSeps). "+
					"Index='%v' Error='%v'", i, errx.Error())

				return fracQuoArray, err
			}
		}

		bPair := BigIntPair{}.NewBigIntNum(dividends[i], divisor)

		bPair.MaxPrecision = maxPrecision

		fracQuoArray[i], errx =
			BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []BigIntNum{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
				"dividend='%v' divisor='%v' maxPrecision='%v' Index='%v' Error='%v'",
				bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, i, errx.Error())

			return fracQuoArray, err
		}

	}

	return fracQuoArray, nil
}

// BigIntNumDivideByTwoQuoMod - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value two ('2'). There are two BigIntNum return values: 'quotient' and
// 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by two
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	   4							/						 2			  =			 2							 0
//	   5  	 					/ 				 	 2  			= 		 2							 1
//	   2.5 					/ 				 	 2				= 	   1							 0.5
//		-12.555 				/ 				   2  			= 		-6							-0.555
//	   0 						/ 				 	 2 				= 		 0							 0
//	 -19	 						/ 				   2    		= 		-9							-1
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumDivideByTwoQuoMod(
	dividend BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	bINumTwo := BigIntNum{}.NewTwo(0)

	return bIDivide.BigIntNumQuotientMod(dividend, bINumTwo, maxPrecision)
}

// BigIntNumDivideByThreeQuoMod - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value three ('3'). There are two BigIntNum return values: 'quotient' and
// 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by three
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	 4							/						 3			  =			  1							 	1
//	 5  	 					/ 				 	 3  			= 		  1							 	2
//	 8 						/ 				 	 3				= 	    2							 	2
//	12							/            3				=				4								0
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumDivideByThreeQuoMod(
	dividend BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	bINumThree := BigIntNum{}.NewThree(0)

	return bIDivide.BigIntNumQuotientMod(dividend, bINumThree, maxPrecision)
}

// BigIntNumDivideByFiveQuoMod - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value five ('5'). There are two BigIntNum return values: 'quotient' and
// 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by five
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	 6							/						 5			  =			 1							 1
//	12  	 					/ 				 	 5  			= 		 2							 2
//	16 						/ 				 	 5				= 	   3							 1
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumDivideByFiveQuoMod(
	dividend BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	bINumFive := BigIntNum{}.NewFive(0)

	return bIDivide.BigIntNumQuotientMod(dividend, bINumFive, maxPrecision)
}

// BigIntNumDivideByTenQuoMod - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value 10 ('10'). There are two BigIntNum return values: 'quotient' and
// 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by 10
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	16							/						 10			  =			 1							 6
//	32  	 					/ 				 	 10  			= 		 3							 2
//	96 						/ 				 	 10				= 	   9							 6
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenQuoMod(
	dividend BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	bINumTen := BigIntNum{}.NewTen(0)

	return bIDivide.BigIntNumQuotientMod(dividend, bINumTen, maxPrecision)
}

// BigIntNumDivideByTenToPowerQuoMod - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the numeric value of ten
// (10) to the power of input parameter 'exponent'. There are two BigIntNum return
// values: 'quotient' and 'modulo'.
//
//	quotient, modulo = dividend / (10^exponent)
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by 10
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// The returned BigIntNum division 'result' (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenToPowerQuoMod(
	dividend, exponent BigIntNum,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumDivideByTenToPowerQuoMod() "
	quotient = BigIntNum{}.NewZero(0)
	modulo = BigIntNum{}.NewZero(0)

	scaleValue, err :=
		BigIntMathPower{}.Pwr(BigIntNum{}.NewTen(0), exponent, maxPrecision+10)

	if err != nil {
		return quotient, modulo,
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathPower{}.Pwr(10, exponent, maxPrecision + 10) "+
				"Error='%v'", err.Error())
	}

	return bIDivide.BigIntNumQuotientMod(dividend, scaleValue, maxPrecision)
}

// BigIntNumDivideByTwoFracQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value two ('2').
//
// The result of this division operation is returned as a BigIntNum type representing
// quotient as integer and fractional digits. Remember that the BigIntNum type specifies
// 'precision'. Precision is defined as the number of fractional digits to the right of
// the decimal place.
//
// Examples:
// =========
//
//	 Dividend		divided by	Divisor		  =		 	Precision	 		 	Result
//	 -------- 	  ----------	--------				  ---------	 	 		------
//		 10.5  				/ 				   2 				= 			  2  			 	 	 5.25
//		 10    				/ 				   2 				= 			 	0  			 	 	 5
//	 -12.555     	/    			   2  			= 		  	4						-6.2775
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByTwoFracQuo(
	dividend BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	bINumTwo := BigIntNum{}.NewTwo(0)

	return bIDivide.BigIntNumFracQuotient(dividend, bINumTwo, maxPrecision)
}

// BigIntNumDivideByThreeFracQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value three ('3').
//
// The result of this division operation is returned as a BigIntNum type representing
// quotient as integer and fractional digits. Remember that the BigIntNum type specifies
// 'precision'. Precision is defined as the number of fractional digits to the right of
// the decimal place.
//
// Examples:
// =========
//
//	 For this example, assume that 'maxPrecision' = 15
//
//	 Dividend		divided by	Divisor		  =		 	Precision	 		 	Result
//	 -------- 	  ----------	--------				  ---------	 	 		------
//		 9.5  				/ 				   3 				= 			  2  			 	 	 3.166666666666667
//		 10    				/ 				   3 				= 			 	0  			 	 	 3.333333333333333
//	  12						/						 3			  =					0						 4
//	 -12     			/    			   3  			= 		  	0						-4
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByThreeFracQuo(
	dividend BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	bINumThree := BigIntNum{}.NewThree(0)

	return bIDivide.BigIntNumFracQuotient(dividend, bINumThree, maxPrecision)
}

// BigIntNumDivideByFiveFracQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value five ('5').
//
// The result of this division operation is returned as a BigIntNum type representing
// quotient as integer and fractional digits. Remember that the BigIntNum type specifies
// 'precision'. Precision is defined as the number of fractional digits to the right of
// the decimal place.
//
// Examples:
// =========
//
//	 Dividend		divided by	Divisor		  =		 	Precision	 		 	Result
//	 -------- 	  ----------	--------				  ---------	 	 		------
//		 16.2  				/ 				   5 				= 			  2  			 	 	 3.24
//		 10    				/ 				   5 				= 			 	0  			 	 	 2
//	  12						/						 5			  =					1						 2.4
//	 -12     			/    			   5  			= 		  	1						-2.4
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByFiveFracQuo(
	dividend BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	bINumFive := BigIntNum{}.NewFive(0)

	return bIDivide.BigIntNumFracQuotient(dividend, bINumFive, maxPrecision)
}

// BigIntNumDivideByTenFracQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the integer numeric
// value ten ('10').
//
// The result of this division operation is returned as a BigIntNum type representing
// quotient as integer and fractional digits. Remember that the BigIntNum type specifies
// 'precision'. Precision is defined as the number of fractional digits to the right of
// the decimal place.
//
// Examples:
// =========
//
//	 For this example, assume that 'maxPrecision' = 15
//
//	 Dividend		divided by	Divisor		  =		 	Precision	 		 	Result
//	 -------- 	  ----------	--------				  ---------	 	 		------
//		 16.2  				/ 				  10				= 			  2  			 	 	 1.62
//		 10    				/ 				  10 				= 			 	0  			 	 	 1
//	  12						/						10			  =					1						 1.2
//	 -12     			/    			  10  			= 		  	1						-1.2
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenFracQuo(
	dividend BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	bINumTen := BigIntNum{}.NewTen(0)

	return bIDivide.BigIntNumFracQuotient(dividend, bINumTen, maxPrecision)
}

// BigIntNumDivideByTenToPowerFracQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the numeric value ten
// ('10') to the power of input parameter 'exponent'.
//
//	fractional quotient = dividend / (10^exponent)
//
// The result of this division operation is returned as a BigIntNum type representing
// quotient as integer and fractional digits. Remember that the BigIntNum type specifies
// 'precision'. Precision is defined as the number of fractional digits to the right of
// the decimal place.
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits.
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenToPowerFracQuo(
	dividend,
	exponent BigIntNum,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumDivideByTenToPowerFracQuo() "
	fracQuotient = BigIntNum{}.NewZero(0)

	scaleValue, err :=
		BigIntMathPower{}.Pwr(BigIntNum{}.NewTen(0), exponent, maxPrecision+10)

	if err != nil {
		return fracQuotient,
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathPower{}.Pwr(10, exponent, maxPrecision + 10) "+
				"Error='%v'", err.Error())
	}

	return bIDivide.BigIntNumFracQuotient(dividend, scaleValue, maxPrecision)
}

// BigIntNumDivideByTenToPowerIntQuo - Performs a division operation on BigIntNum input
// parameter 'dividend'. This method will divide 'dividend' by the numeric value ten
// ('10') to the power of input parameter 'exponent'.
//
//	integer quotient = dividend / (10^exponent)
//
// The result of this division operation is returned as a BigIntNum type representing
// 'quotient' as an integer value.
//
// The input parameter 'maxPrecision' is used to control the maximum precision of the
// internal calculations for those cases with 'exponent' is a fractional or floating
// point number. Be advised that this method is capable of calculating quotients with very
// long strings of fractional digits.
//
// The returned BigIntNum division result ('intQuotient') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenToPowerIntQuo(
	dividend,
	exponent BigIntNum,
	maxPrecision uint) (intQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumDivideByTenToPowerIntQuo() "
	intQuotient = BigIntNum{}.NewZero(0)

	scaleValue, err :=
		BigIntMathPower{}.Pwr(BigIntNum{}.NewTen(0), exponent, maxPrecision+10)

	if err != nil {
		return intQuotient,
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathPower{}.Pwr(10, exponent, maxPrecision + 10) "+
				"Error='%v'", err.Error())
	}

	return bIDivide.BigIntNumIntQuotient(dividend, scaleValue)
}

// BigIntNumDivideByTenToPowerMod - Performs a modulo operation on BigIntNum input
// parameters 'dividend' and ten (10) to the power of exponent.
//
//	modulo = dividend / (10^exponent)
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one BigIntNum value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
//	Dividend			  mod by	Power 		Divisor			=			Modulo/Remainder
//
// ---------				------	-----		  -------						----------------
//
//	1200.555					%				2		 			100			  =			 		0.555
//
// 10235.555					%				3		 		 1000			  =			 	235.555
//
// The returned BigIntNum division 'result' (modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) BigIntNumDivideByTenToPowerMod(
	dividend,
	exponent BigIntNum,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.BigIntNumDivideByTenToPowerIntQuo() "
	modulo = BigIntNum{}.NewZero(0)

	scaleValue, err :=
		BigIntMathPower{}.Pwr(BigIntNum{}.NewTen(0), exponent, maxPrecision+10)

	if err != nil {
		return modulo,
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathPower{}.Pwr(10, exponent, maxPrecision + 10) "+
				"Error='%v'", err.Error())
	}

	return bIDivide.BigIntNumModulo(dividend, scaleValue, maxPrecision)
}

// DecimalQuotientMod - Performs a division operation on Decimal type input
// parameters, 'dividend' and 'divisor'.
//
// There are two BigIntNum Type return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values. Therefore, the user is advised to set a relevant value
// for 'maxPrecision'.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results ('quotient' and 'modulo') will
// contain numeric separators (decimal separator, thousands separator and
// currency symbol) copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) DecimalQuotientMod(
	dividend,
	divisor Decimal,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.DecimalQuotientMod() "

	quotient = BigIntNum{}.New()
	modulo = BigIntNum{}.New()
	err = nil

	bPair, errx := BigIntPair{}.NewDecimal(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewDecimal(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return quotient, modulo, err
	}

	if bPair.Big2.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted divide by ZERO!")
		return quotient, modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	quotient, modulo, errx =
		BigIntMathDivide{}.PairQuotientMod(bPair)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}."+
			"PairQuotientMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		quotient = BigIntNum{}.New()
		modulo = BigIntNum{}.New()

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// DecimalFracQuotient - Performs a division operation on Decimal Type input
// parameters 'dividend' and 'divisor'.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
//
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					     Quotient
//	 Dividend		divided by	Divisor		=		  BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input parameter,
// 'dividend'.
func (bIDivide BigIntMathDivide) DecimalFracQuotient(
	dividend,
	divisor Decimal,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.DecimalFracQuotient() "

	if divisor.IsZero() {
		fracQuotient = BigIntNum{}.New()
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuotient, err
	}

	bPair, errx := BigIntPair{}.NewDecimal(dividend, divisor)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewDecimal(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return fracQuotient, err
	}

	bPair.MaxPrecision = maxPrecision

	fracQuotient, errx =
		BigIntMathDivide{}.PairFracQuotient(bPair)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	return fracQuotient, nil
}

// DecimalFracQuotientArray - Performs a division operation on Decimal input
// parameters 'dividends' and 'divisor'. 'dividends' is an array of Decimal
// Types. The division operation is performed on each element of the 'dividends'
// array using a single 'divisor'.
//
// The resulting quotients are returned as an array of Decimal Types. The values
// represent result of each division operation expressed as integer and fractional
// digits.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Precision is defined as the number of numeric
// digits to the right of the decimal point. Be advised that this method is
// capable of calculating quotients with very long strings of fractional digits.
// Therefore the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// Each element in the returned division results array ('fracQuoArray') will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// copied from the first element of the input parameter 'dividends' array.
func (bIDivide BigIntMathDivide) DecimalFracQuotientArray(
	dividends []Decimal,
	divisor Decimal,
	maxPrecision uint) (fracQuoArray []Decimal, err error) {

	ePrefix := "BigIntMathDivide.DecimalFracQuotientArray() "

	if divisor.IsZero() {
		fracQuoArray = []Decimal{}
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuoArray, err
	}

	lenAry := len(dividends)

	if lenAry == 0 {
		fracQuoArray = []Decimal{}
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	numSeps := dividends[0].GetNumericSeparatorsDto()

	fracQuoArray = make([]Decimal, lenAry, lenAry+20)

	var errx error

	for i := 0; i < lenAry; i++ {

		if i > 0 {

			errx = dividends[i].SetNumericSeparatorsDto(numSeps)

			if errx != nil {
				fracQuoArray = []Decimal{}
				err = fmt.Errorf(ePrefix+
					"Error returned by dividends[%v].SetNumericSeparatorsDto(numSeps). "+
					"Error='%v'", i, errx.Error())

				return fracQuoArray, err
			}
		}

		bPair, errx := BigIntPair{}.NewDecimal(dividends[i], divisor)

		if errx != nil {
			fracQuoArray = []Decimal{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewDecimal(dividends[i], divisor). "+
				"dividends[%v]='%v' divisor='%v' Error='%v'",
				i, dividends[i].GetNumStr(), divisor.GetNumStr(), errx.Error())

			return fracQuoArray, err
		}

		bPair.MaxPrecision = maxPrecision

		bINum, errx := BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []Decimal{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
				"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
				bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())
			return fracQuoArray, err
		}

		fracQuoArray[i], errx = bINum.GetDecimal()

		if errx != nil {
			fracQuoArray = []Decimal{}
			err = fmt.Errorf(ePrefix+
				"Error returned by bINum.GetDecimal(). Error='%v'",
				errx.Error())
			return fracQuoArray, err
		}

	}

	return fracQuoArray, nil
}

// DecimalModulo - Performs a modulo operation on Decimal type
// input parameters, 'dividend' and 'divisor'. The result of this division
// operation is returned as a type BigIntNum.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one Decimal value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) DecimalModulo(
	dividend,
	divisor Decimal,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.DecimalModulo() "

	modulo = BigIntNum{}.New()

	bPair, errx := BigIntPair{}.NewDecimal(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewDecimal(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(),
			maxPrecision, errx.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// DecimalModuloToDecimal - Performs a modulo operation on Decimal input
// parameters 'dividend' and 'divisor'.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one Decimal value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
//
// The difference between this method and BigIntMathDivide.DecimalModulo()
// above, is that this method returns the resulting modulo value as a type
// 'Decimal'.
func (bIDivide BigIntMathDivide) DecimalModuloToDecimal(
	dividend,
	divisor Decimal,
	maxPrecision uint) (modulo Decimal, err error) {

	ePrefix := "BigIntMathDivide.DecimalModuloToDecimal() "

	bPair, errx := BigIntPair{}.NewDecimal(dividend, divisor)

	if errx != nil {
		modulo = Decimal{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewDecimal(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(),
			maxPrecision, errx.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		modulo = Decimal{}.New()
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	bINumModulo, errx := BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = Decimal{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	modulo, errx = bINumModulo.GetDecimal()

	if errx != nil {
		modulo = Decimal{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by bINumModulo.GetDecimal(). "+
			"Error='%v'", errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// FixedDecimalFracQuotient - Performs a division operation on objects
// of type BigIntFixedDecimal. The result is also returned a as a type
// BigIntFixedDecimal.
//
// Examples:
// =========
// This division operation will produce a quotient which may include a fixed
// length floating point number:
//
//	quotient = dividend / divisor
//
// The BigIntFixedDecimal structure is defined as
//
//	type BigIntFixedDecimal struct {
//		integerNum *big.Int  -	All of the numeric digits, both integer and fractional,
//														necessary to define a fixed length floating point number.
//														The number of digits to the right of the decimal place
//														is specified by the data field,
//														BigIntFixedDecimal.precision.
//
//		precision  uint				- Specifies the number of digits to the right of the decimal
//														place in the series of numeric digits represented by data
//														field BigIntFixedDecimal.integerNum.
//
// }
//
//		To represent the floating point number 52.459
//		a BigIntDecimal Structure would be configured as follows:
//				BigIntFixedDecimal.integerNum	= 52459
//				BigIntFixedDecimal.precision	= 3
//
//	 Consider the following division example
//										quotient =	752.314 / 21.67894
//
// 'dividend' and 'divisor' would be configured as follows:
//
//										dividend.integerNum	= 752314
//	                 dividend.precision	= 3
//	                 divisor.integerNum	= 2167894
//	                 divisor.precision		= 5
//
// Assuming a 'maxPrecision' value of '30', the quotient would be
// calculated as follows:
//
//										quotient.integerNum	= 34702526968569496479071393712054
//	                 quotient.precision  = 30
//
// Input Parameters
// ================
//
// dividend	BigIntFixedDecimal	- The 'dividend' value will be divided by the 'divisor'
//
//	to produce a 'quotient'.
//
// divisor	BigIntFixedDecimal	- The 'dividend' value will be divided by the 'divisor'
//
//	to produce a 'quotient'.
//
// 'maxPrecision' 				uint	-	Maximum precision will determine the maximum number
//
//	    of decimal digits to which the result or 'quotient'
//			will calculated and returned to the caller. The
//			quotient may consist of actual fractional digits
//			which number less than 'maxPrecision'. However, if
//			the number of digits to the right of the decimal
//			place exceeds 'maxPrecision', the returned quotient
//			will be rounded to 'maxPrecision' fractional digits
//			to the right of the decimal place.
//
// Return Values
// =============
//
// quotient	BigIntFixedDecimal	- The result of the division operation expressed
//
//	as a type BigIntFixedDecimal.
//
// err			error								- If an error is encountered, this function will
//
//	                               return a quotient set equal to zero and an
//	                               error object will be returned containing an
//	                               appropriate error message. If the function
//	                               completes the division operation successfully,
//	                               the returned 'quotient' will be populated with
//																	the correct result and 'err' will be set equal
//	                               to 'nil'.
func (bIDivide BigIntMathDivide) FixedDecimalFracQuotient(
	dividend BigIntFixedDecimal,
	divisor BigIntFixedDecimal,
	maxPrecision uint) (quotient BigIntFixedDecimal, err error) {

	quotient = BigIntFixedDecimal{}.NewZero(0)
	err = nil

	dividend.IsValid()
	divisor.IsValid()

	result, resultPrecision, errX :=
		BigIntMathDivide{}.BigIntFracQuotient(
			dividend.GetInteger(),
			dividend.GetPrecisionBigInt(),
			divisor.GetInteger(),
			divisor.GetPrecisionBigInt(),
			big.NewInt(0).SetUint64(uint64(maxPrecision)))

	if errX != nil {
		ePrefix := "BigIntMathDivide.FixedDecimalFracQuotient() "
		err = fmt.Errorf(ePrefix+"Error returned by "+
			"BigIntMathDivide{}.BigIntFracQuotient(). Error='%v' ",
			errX.Error())
		return quotient, err
	}

	quotient.SetNumericValue(result, uint(resultPrecision.Uint64()))
	err = nil

	return quotient, nil
}

// IntAryQuotientMod - Performs a division operation on IntAry type input
// parameters, 'dividend' and 'divisor'.
//
// There are two BigIntNum Type return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results ('quotient' and 'modulo') will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) IntAryQuotientMod(
	dividend,
	divisor IntAry,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.IntAryQuotientMod() "

	quotient = BigIntNum{}.New()
	modulo = BigIntNum{}.New()
	err = nil

	bPair, errx := BigIntPair{}.NewIntAry(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewIntAry(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), errx.Error())

		return quotient, modulo, err
	}

	if bPair.Big2.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted divide by ZERO!")
		return quotient, modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	quotient, modulo, errx =
		BigIntMathDivide{}.PairQuotientMod(bPair)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}."+
			"PairQuotientMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v'  Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		quotient = BigIntNum{}.New()
		modulo = BigIntNum{}.New()

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// IntAryFracQuotient - Performs a division operation on IntAry Type input
// parameters 'dividend' and 'divisor'.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
//
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					      Quotient
//	 Dividend		divided by	Divisor		=		  BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input parameter,
// 'dividend'.
func (bIDivide BigIntMathDivide) IntAryFracQuotient(
	dividend,
	divisor IntAry,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.IntAryFracQuotient() "

	if divisor.IsZero() {
		fracQuotient = BigIntNum{}.New()
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuotient, err
	}

	bPair, errx := BigIntPair{}.NewIntAry(dividend, divisor)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewIntAry(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return fracQuotient, err
	}

	bPair.MaxPrecision = maxPrecision

	fracQuotient, errx =
		BigIntMathDivide{}.PairFracQuotient(bPair)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	fracQuotient.TrimTrailingFracZeros()

	return fracQuotient, nil
}

// IntAryFracQuotientArray - Performs a division operation on IntAry input
// parameters 'dividends' and 'divisor'. 'dividends' is an array of IntAry
// Types. The division operation is performed on each element of the 'dividends'
// array using a single 'divisor'.
//
// The resulting quotients are returned as an array of IntAry Types. The values
// represent result of each division operation expressed as integer and fractional
// digits.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Precision is defined as the number of numeric
// digits to the right of the decimal point. Be advised that this method is
// capable of calculating quotients with very long strings of fractional digits.
// Therefore the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// The returned []IntAry division result ('fracQuoArray') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from the first
// element of the input parameter 'dividends' array.
func (bIDivide BigIntMathDivide) IntAryFracQuotientArray(
	dividends []IntAry,
	divisor IntAry,
	maxPrecision uint) (fracQuoArray []IntAry, err error) {

	ePrefix := "BigIntMathDivide.IntAryFracQuotientArray() "

	if divisor.IsZero() {
		fracQuoArray = []IntAry{}
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuoArray, err
	}

	lenAry := len(dividends)

	if lenAry == 0 {
		fracQuoArray = []IntAry{}
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	numSeps := dividends[0].GetNumericSeparatorsDto()

	fracQuoArray = make([]IntAry, lenAry, lenAry+20)

	for i := 0; i < lenAry; i++ {

		if i > 0 {

			errx := dividends[i].SetNumericSeparatorsDto(numSeps)

			if errx != nil {

				fracQuoArray = []IntAry{}

				err = fmt.Errorf(ePrefix+
					"Error returned by dividends[%v].SetNumericSeparatorsDto(numSeps). "+
					"Error='%v'", i, errx.Error())

				return fracQuoArray, err
			}

		}

		bPair, errx := BigIntPair{}.NewIntAry(dividends[i], divisor)

		if errx != nil {

			fracQuoArray = []IntAry{}

			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewIntAry(dividends[i], divisor). "+
				"dividends[%v]='%v' divisor='%v' Error='%v'",
				i, dividends[i].GetNumStr(), divisor.GetNumStr(), errx.Error())

			return fracQuoArray, err
		}

		bPair.MaxPrecision = maxPrecision

		bINum, errx := BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []IntAry{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
				"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
				bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())
			return fracQuoArray, err
		}

		fracQuoArray[i], errx = bINum.GetIntAry()

		if errx != nil {
			fracQuoArray = []IntAry{}
			err = fmt.Errorf(ePrefix+
				"Error returned by bINum.GetIntAryElements(). Error='%v'",
				errx.Error())
			return fracQuoArray, err
		}

	}

	return fracQuoArray, nil
}

// IntAryModulo - Performs a modulo operation on IntAry input
// parameters 'dividend' and 'divisor'.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one BigIntNum value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input
// parameter, 'dividend'.
func (bIDivide BigIntMathDivide) IntAryModulo(
	dividend,
	divisor IntAry,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.IntAryModulo() "

	bPair, errx := BigIntPair{}.NewIntAry(dividend, divisor)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewIntAry(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(),
			maxPrecision, errx.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// IntAryModuloToIntAry - Performs a modulo operation on IntAry input
// parameters 'dividend' and 'divisor'.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one IntAry value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned IntAry division result ('modulo') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input
// parameter, 'dividend'.
//
// The difference between this method and BigIntMathDivide.IntAryModulo() above,
// is that the division result, 'modulo', is returned as Type IntAry.
func (bIDivide BigIntMathDivide) IntAryModuloToIntAry(
	dividend,
	divisor IntAry,
	maxPrecision uint) (modulo IntAry, err error) {

	ePrefix := "BigIntMathDivide.IntAryModuloToIntAry() "

	bPair, errx := BigIntPair{}.NewIntAry(dividend, divisor)

	if errx != nil {
		modulo = IntAry{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewIntAry(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(),
			maxPrecision, errx.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		modulo = IntAry{}.New()
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	bINumModulo, errx := BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = IntAry{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	modulo, errx = bINumModulo.GetIntAry()

	if errx != nil {
		modulo = IntAry{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by bINumModulo.GetIntAryElements(). "+
			"Error='%v'", errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// INumMgrQuotientMod - Performs a division operation on types implementing the INumMgr
// interface. Input parameters, 'dividend' and 'divisor' must therefore implement the
// INumMgr interface.
//
// There are two BigIntNum Type return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values. Therefore, the user is advised to set a relevant value
// for 'maxPrecision'.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results ('quotient' and 'modulo') will contain numeric
// separators (decimal separator, thousands separator and currency symbol) copied from
// input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) INumMgrQuotientMod(
	dividend,
	divisor INumMgr,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.INumMgrQuotientMod() "

	quotient = BigIntNum{}.New()
	modulo = BigIntNum{}.New()
	err = nil

	bPair, errx := BigIntPair{}.NewINumMgr(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewINumMgr(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), maxPrecision, errx.Error())

		return quotient, modulo, err
	}

	if bPair.Big2.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted divide by ZERO!")
		return quotient, modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	quotient, modulo, errx =
		BigIntMathDivide{}.PairQuotientMod(bPair)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairQuotientMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		quotient = BigIntNum{}.New()
		modulo = BigIntNum{}.New()

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// INumMgrFracQuotient - Performs a division operation on input
// parameters 'dividend' and 'divisor' which implement the INumMgr
// interface.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
//
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					     Quotient
//	 Dividend		divided by	Divisor		=		  BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input parameter,
// 'dividend'.
func (bIDivide BigIntMathDivide) INumMgrFracQuotient(
	dividend,
	divisor INumMgr,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.INumMgrFracQuotient() "

	if divisor.IsZero() {
		fracQuotient = BigIntNum{}.New()
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuotient, err
	}

	// Validity tests are performed on 'dividend' and 'divisor'
	bPair, errx := BigIntPair{}.NewINumMgr(dividend, divisor)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewINumMgr(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return fracQuotient, err
	}

	bPair.MaxPrecision = maxPrecision

	fracQuotient, errx =
		BigIntMathDivide{}.PairFracQuotient(bPair)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	return fracQuotient, nil
}

// INumMgrFracQuotientArray - Performs a division operation on input parameters
// 'dividends' and 'divisor' which implement the INumMgr interface. 'dividends'
// is an array of types implementing the INumMgr interface. The division operation
// is performed on each element of the 'dividends' array using a single 'divisor'.
//
// The resulting quotients are returned as an array of types implementing the INumMgr
// Interface. The values represent the results of each division operation expressed
// as integer and fractional digits.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Precision is defined as the number of numeric
// digits to the right of the decimal point. Be advised that this method is
// capable of calculating quotients with very long strings of fractional digits.
// Therefore the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// The returned INumMgr array division result ('fracQuoArray') will contain numeric
// separators (decimal separator, thousands separator and currency symbol) copied
// from the first element of input parameter, 'dividends' array.
func (bIDivide BigIntMathDivide) INumMgrFracQuotientArray(
	dividends []INumMgr,
	divisor INumMgr,
	maxPrecision uint) (fracQuoArray []INumMgr, err error) {

	ePrefix := "BigIntMathDivide.INumMgrFracQuotientArray() "

	if divisor.IsZero() {
		fracQuoArray = []INumMgr{}
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuoArray, err
	}

	lenAry := len(dividends)

	if lenAry == 0 {
		fracQuoArray = []INumMgr{}
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	numSeps := dividends[0].GetNumericSeparatorsDto()

	fracQuoArray = make([]INumMgr, lenAry, lenAry+20)

	for i := 0; i < lenAry; i++ {

		if i > 0 {

			errx := dividends[i].SetNumericSeparatorsDto(numSeps)

			if errx != nil {
				fracQuoArray = []INumMgr{}
				err = fmt.Errorf(ePrefix+
					"Error returned by dividends[%v].SetNumericSeparatorsDto(numSeps). "+
					"Error='%v'", i, errx.Error())

				return fracQuoArray, err
			}

		}

		bPair, errx := BigIntPair{}.NewINumMgr(dividends[i], divisor)

		if errx != nil {
			fracQuoArray = []INumMgr{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewINumMgr(dividends[i], divisor). "+
				"dividends[%v]='%v' divisor='%v' Error='%v'",
				i, dividends[i].GetNumStr(), divisor.GetNumStr(), errx.Error())

			return fracQuoArray, err
		}

		bPair.MaxPrecision = maxPrecision

		bINum, errx := BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []INumMgr{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
				"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
				bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

			return fracQuoArray, err
		}

		fracQuoArray[i] = &bINum

	}

	err = nil

	return fracQuoArray, err
}

// INumMgrModulo - Performs a modulo operation on input parameters
// 'dividend' and 'divisor'. Both parameters must implement the INumMgr
// interface.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one BigIntNum value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) INumMgrModulo(
	dividend,
	divisor INumMgr,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.INumMgrModulo() "

	bPair, errx := BigIntPair{}.NewINumMgr(dividend, divisor)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewINumMgr("+
			"dividend, divisor) dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), err.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// NumStrQuotientMod - Performs a division operation on input parameters 'dividend' and 'divisor'
// which are comprised as number strings. Number strings are strings of numeric digits representing
// a numeric value. Number strings may include a leading minus sign (-) indicating a negative
// numeric value. Number strings may also include a decimal separator used to separate integer
// and fractional digits. The decimal separator character is specified by input parameter, 'numSeps'.
//
// There are two BigIntNum return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the dividend and divisor number strings. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// division operation.
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results ('quotient' and 'modulo') will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) designated by the input parameter, 'numSeps'.
func (bIDivide BigIntMathDivide) NumStrQuotientMod(
	dividend,
	divisor string,
	numSeps NumericSeparatorDto,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrQuotientMod() "

	quotient = BigIntNum{}.NewBigInt(big.NewInt(0), 0)

	modulo = BigIntNum{}.NewBigInt(big.NewInt(0), 0)

	numSeps.SetDefaultsIfEmpty()

	bigIDividend, errx := BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps) "+
			"dividend='%v' numSeps='%v' Error='%v' ",
			dividend, numSeps.String(), errx.Error())

		return quotient, modulo, err
	}

	bigIDivisor, errx := BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps) "+
			"divisor='%v' numSeps='%v' Error='%v' ",
			divisor, numSeps.String(), errx.Error())

		return quotient, modulo, err
	}

	if bigIDivisor.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted Divide By ZERO!")

		return quotient, modulo, err
	}

	bPair := BigIntPair{}.NewBigIntNum(bigIDividend, bigIDivisor)

	bPair.MaxPrecision = maxPrecision

	quotient, modulo, errx =
		BigIntMathDivide{}.PairQuotientMod(bPair)

	if errx != nil {
		quotient = BigIntNum{}.New()
		modulo = BigIntNum{}.New()

		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}."+
			"PairQuotientMod(bPair). "+
			" dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return quotient, modulo, err
	}

	return quotient, modulo, nil
}

// NumStrFracQuotient - Performs a division operation on input parameters
// 'dividend' and 'divisor' and returns the 'quotient' of this division
// operation.
//
// The input parameters, 'dividend' and 'divisor', are number strings consisting
// of strings of numeric digits representing a specific numeric value. Number
// strings may include a leading minus sign (-) indicating a negative numeric
// value. In addition, number strings may include decimal separators used to
// separate integer and fractional digits. The decimal separator character is
// specified by the input parameter, 'numSeps'.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the dividend and divisor number strings. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// division operation.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
//
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					     Quotient
//	 Dividend		divided by	Divisor		=		  BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result ('fracQuotient') will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) designated by the input parameter, 'numSeps'.
func (bIDivide BigIntMathDivide) NumStrFracQuotient(
	dividend,
	divisor string,
	numSeps NumericSeparatorDto,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrFracQuotient() "

	fracQuotient = BigIntNum{}.NewZero(0)

	numSeps.SetDefaultsIfEmpty()

	bigIDividend, errx := BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps) "+
			"dividend='%v' numSeps='%v' Error='%v' ",
			dividend, numSeps.String(), errx.Error())

		return fracQuotient, err
	}

	bigIDivisor, errx := BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps) "+
			"divisor='%v' numSeps='%v' Error='%v' ",
			divisor, numSeps.String(), errx.Error())

		return fracQuotient, err
	}

	if bigIDivisor.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted Divide By ZERO!")

		return fracQuotient, err
	}

	bPair := BigIntPair{}.NewBigIntNum(bigIDividend, bigIDivisor)

	bPair.MaxPrecision = maxPrecision

	fracQuotient, errx =
		BigIntMathDivide{}.PairFracQuotient(bPair)

	if errx != nil {
		fracQuotient = BigIntNum{}
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}.PairFracQuotient("+
			"bPair)"+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v' ",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// NumStrFracQuotientArray - Performs a division operation on input parameters
// 'dividends' and 'divisor' which are formatted as number strings.
//
// Number strings are strings of numeric digits representing a specific numeric
// value. Number strings may include a leading minus sign (-) indicating a negative
// numeric value. In addition, number strings may include a decimal separator used
// to separate integer and fractional digits. The decimal separator character is
// specified by the input parameter, 'numSeps'.
//
// 'dividends' is an array of number strings. The division operation is performed
// on each element of the 'dividends' array using a single 'divisor'.
//
// The resulting quotients are returned as an array of BigIntNum types. The values
// represent result of each division operation expressed as integer and fractional
// digits.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the dividends array and divisor number strings. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// division operation.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Precision is defined as the number of numeric
// digits to the right of the decimal point. Be advised that this method is
// capable of calculating quotients with very long strings of fractional digits.
// Therefore the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// Each element of the returned BigIntNum Array ('fracQuoArray') will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// specified by the 'numSeps' input parameter.
func (bIDivide BigIntMathDivide) NumStrFracQuotientArray(
	dividends []string,
	divisor string,
	numSeps NumericSeparatorDto,
	maxPrecision uint) (fracQuoArray []BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrFracQuotientArray() "

	fracQuoArray = []BigIntNum{}

	numSeps.SetDefaultsIfEmpty()

	lenAry := len(dividends)

	if lenAry == 0 {
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	fracQuoArray = make([]BigIntNum, lenAry, lenAry+20)

	bigINumDivisor, errx := BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps)

	if errx != nil {
		err = fmt.Errorf("Error returned by NewNumStrWithNumSeps(divisor, numSeps). "+
			"divisor='%v' numSeps='%v' Error='%v'. ",
			divisor, numSeps.String(), err.Error())
	}

	for i := 0; i < lenAry; i++ {

		bigINumDividend, errx := BigIntNum{}.NewNumStrWithNumSeps(dividends[i], numSeps)

		if errx != nil {
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.BigIntNum{}.NewNumStrWithNumSeps("+
				"dividends[i], numSeps) "+
				"dividends[%v]='%v' numSeps='%v' Error='%v' ",
				i, dividends[i], numSeps.String(), errx.Error())

			return fracQuoArray, err

		}

		bPair := BigIntPair{}.NewBigIntNum(bigINumDividend, bigINumDivisor)

		if errx != nil {
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewBigIntNum(bigINumDividend, bigINumDivisor) "+
				"bigINumDividend='%v', bigINumDivisor='%v' Index='%v' Error='%v' ",
				bigINumDividend.GetNumStr(), bigINumDivisor.GetNumStr(), i, errx.Error())

			return fracQuoArray, err
		}

		bPair.MaxPrecision = maxPrecision

		fracQuoArray[i], errx =
			BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []BigIntNum{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair) "+
				"dividend['%v']='%v' divisor='%v' maxPrecision='%v' Error='%v' ",
				i, bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), maxPrecision, errx.Error())

			return fracQuoArray, err
		}

	}

	err = nil

	return fracQuoArray, err
}

// NumStrModulo - Performs a modulo operation on input parameters 'dividend'
// and 'divisor'. Both input parameters are formatted a number strings. Number
// strings are strings of numeric digits which represent a specific numeric value.
// Number strings may include a leading minus sign indicating a negative numeric
// value. In addition, number strings may also include a decimal separator used
// to separate integer and fractional digits. The decimal separator character is
// specified by input parameter, 'numSeps'.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one BigIntNum value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the dividend and divisor number strings. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// division operation.
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain default
// numeric separators (decimal separator, thousands separator and currency
// symbol).
func (bIDivide BigIntMathDivide) NumStrModulo(
	dividend,
	divisor string,
	numSeps NumericSeparatorDto,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrModulo() "

	modulo = BigIntNum{}.NewZero(0)

	numSeps.SetDefaultsIfEmpty()

	bigIDividend, errx := BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(dividend, numSeps) "+
			"dividend='%v' numSeps='%v' Error='%v' ",
			dividend, numSeps.String(), errx.Error())

		return modulo, err
	}

	bigIDivisor, errx := BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStrWithNumSeps(divisor, numSeps) "+
			"divisor='%v' numSeps='%v' Error='%v' ",
			divisor, numSeps.String(), errx.Error())

		return modulo, err
	}

	if bigIDivisor.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted Divide By ZERO!")

		return modulo, err
	}

	bPair := BigIntPair{}.NewBigIntNum(bigIDividend, bigIDivisor)

	bPair.MaxPrecision = maxPrecision

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// NumStrDtoQuotientMod - Performs a division operation on NumStrDto type input
// parameters, 'dividend' and 'divisor'.
//
// There are two BigIntNum Type return values: 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values. Therefore, the user is advised to set a relevant value
// for 'maxPrecision'.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results ('quotient' and 'modulo') will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter, 'dividend'.
func (bIDivide BigIntMathDivide) NumStrDtoQuotientMod(
	dividend,
	divisor NumStrDto,
	maxPrecision uint) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrDtoQuotientMod() "

	quotient = BigIntNum{}.New()
	modulo = BigIntNum{}.New()
	err = nil

	// This method will test the validity of dividend and divisor
	bPair, errx := BigIntPair{}.NewNumStrDto(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewNumStrDto(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return quotient, modulo, err
	}

	if bPair.Big2.IsZero() {
		err = errors.New(ePrefix + "Error: Attempted divide by ZERO!")
		return quotient, modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	quotient, modulo, errx =
		BigIntMathDivide{}.PairQuotientMod(bPair)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"Error returned by BigIntMathDivide{}."+
			"PairQuotientMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		quotient = BigIntNum{}.New()
		modulo = BigIntNum{}.New()

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// NumStrDtoFracQuotient - Performs a division operation on NumStrDto Type input
// parameters 'dividend' and 'divisor'.
//
// The resulting quotient is returned as a BigIntNum type representing the
// result of the division operation expressed as integer and fractional digits.
// Remember that the BigIntNum type specifies 'precision'. Precision is defined
// as the number of fractional digits to the right of the decimal place.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
//
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					      Quotient
//	 Dividend		divided by	Divisor		=		  BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result ('fracQuotient') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input parameter,
// 'dividend'.
func (bIDivide BigIntMathDivide) NumStrDtoFracQuotient(
	dividend,
	divisor NumStrDto,
	maxPrecision uint) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrDtoFracQuotient() "

	if divisor.IsZero() {
		fracQuotient = BigIntNum{}.New()
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuotient, err
	}

	bPair, errx := BigIntPair{}.NewNumStrDto(dividend, divisor)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewNumStrDto(dividend, divisor). "+
			"dividend='%v' divisor='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(), errx.Error())

		return fracQuotient, err
	}

	bPair.MaxPrecision = maxPrecision

	fracQuotient, errx =
		BigIntMathDivide{}.PairFracQuotient(bPair)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	fracQuotient.TrimTrailingFracZeros()

	return fracQuotient, nil
}

// NumStrDtoFracQuotientArray - Performs a division operation on NumStrDto input
// parameters 'dividends' and 'divisor'. 'dividends' is an array of NumStrDto
// Types. The division operation is performed on each element of the 'dividends'
// array using a single 'divisor'.
//
// The resulting quotients are returned as an array of NumStrDto Types. The values
// represent result of each division operation expressed as integer and fractional
// digits.
//
// The input parameter 'maxPrecision' is used to control the precision of the
// resulting fractional quotient. Precision is defined as the number of numeric
// digits to the right of the decimal point. Be advised that this method is
// capable of calculating quotients with very long strings of fractional digits.
// Therefore the user is advised to set a relevant value for 'maxPrecision'.
//
// Examples:
// =========
// Note: For all examples maximum precision is specified as '15'.
// ----------------------------------------------------------------------------
//
//	   	                                       Returned
//	 	Dividend		divided by	Divisor		=	       Array						=  	 Result
//	 	-------- 	  ----------	--------				-----------------	   		---------
//		 	 10.5  				 / 				2.5 			= 		 fracQuoArray[0]		=   	 4.2
//		 	 10    				 / 				2.5 			= 		 fracQuoArray[1]		=	 	   4
//		   11.5  				 /        2.5				=  		 fracQuoArray[2]		=		   4.6
//		    2.5					 /				2.5			  =			 fracQuoArray[3]    =      1
//			-12.555 			 / 				2.5 			= 		 fracQuoArray[4]    =  -   5.022
//		  - 2.5 				 / 			  2.5		    = 		 fracQuoArray[5]    =  -   1
//		   12.555 			 / 			  2.5 			= 		 fracQuoArray[6]    =      5.022
//		 -122.783 			 / 			  2.5 			= 		 fracQuoArray[7]    =  -  49.1132
//		-6847.231   	   /    	  2.5 			= 		 fracQuoArray[8]    =  -2738.8924
//		  - 2.5	 				 / 			  2.5		    = 		 fracQuoArray[9]    =  -   1
//		  -10						 /			  2.5				=			 fracQuoArray[10]   =  -   4
//		  -10.5					 /			  2.5				=			 fracQuoArray[11]   =  -   4.2
//
// The returned []NumStrDto division result ('fracQuoArray') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from the first
// element of the input parameter 'dividends' array.
func (bIDivide BigIntMathDivide) NumStrDtoFracQuotientArray(
	dividends []NumStrDto,
	divisor NumStrDto,
	maxPrecision uint) (fracQuoArray []NumStrDto, err error) {

	ePrefix := "BigIntMathDivide.NumStrDtoFracQuotientArray() "

	if divisor.IsZero() {
		fracQuoArray = []NumStrDto{}
		err = errors.New(ePrefix + "Error: Attempted divide by zero!")
		return fracQuoArray, err
	}

	lenAry := len(dividends)

	if lenAry == 0 {
		fracQuoArray = []NumStrDto{}
		err = errors.New(ePrefix + "Error: Input Parameter 'dividends' is an EMPTY Array!")
		return fracQuoArray, err
	}

	numSeps := dividends[0].GetNumericSeparatorsDto()

	fracQuoArray = make([]NumStrDto, lenAry, lenAry+20)

	for i := 0; i < lenAry; i++ {

		if i > 0 {

			errx := dividends[i].SetNumericSeparatorsDto(numSeps)

			if errx != nil {

				fracQuoArray = []NumStrDto{}

				err = fmt.Errorf(ePrefix+
					"Error returned by dividends[%v].SetNumericSeparatorsDto(numSeps). "+
					"Error='%v'", i, errx.Error())

				return fracQuoArray, err
			}

		}

		bPair, errx := BigIntPair{}.NewNumStrDto(dividends[i], divisor)

		if errx != nil {

			fracQuoArray = []NumStrDto{}

			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewNumStrDto(dividends[i], divisor). "+
				"dividends[%v]='%v' divisor='%v' Error='%v'",
				i, dividends[i].GetNumStr(), divisor.GetNumStr(), errx.Error())

			return fracQuoArray, err
		}

		bPair.MaxPrecision = maxPrecision

		bINum, errx := BigIntMathDivide{}.PairFracQuotient(bPair)

		if errx != nil {
			fracQuoArray = []NumStrDto{}
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.PairFracQuotient(bPair). "+
				"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
				bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(), bPair.MaxPrecision, errx.Error())
			return fracQuoArray, err
		}

		fracQuoArray[i], errx = bINum.GetNumStrDto()

		if errx != nil {
			fracQuoArray = []NumStrDto{}
			err = fmt.Errorf(ePrefix+
				"Error returned by bINum.GetNumStrDtoElements(). Error='%v'",
				errx.Error())
			return fracQuoArray, err
		}

	}

	return fracQuoArray, nil
}

// NumStrDtoModulo - Performs a modulo operation on NumStrDto type
// input parameters, 'dividend' and 'divisor'. The result of this division
// operation is returned as a type BigIntNum.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one NumStrDto value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result ('modulo') will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from input
// parameter, 'dividend'.
func (bIDivide BigIntMathDivide) NumStrDtoModulo(
	dividend,
	divisor NumStrDto,
	maxPrecision uint) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.NumStrDtoModulo() "

	modulo = BigIntNum{}.New()

	// This method will test the validity of dividend and divisor.
	bPair, errx := BigIntPair{}.NewNumStrDto(dividend, divisor)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntPair{}.NewNumStrDto(dividend, divisor). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			dividend.GetNumStr(), divisor.GetNumStr(),
			maxPrecision, errx.Error())

		return modulo, err
	}

	if bPair.Big2.IsZero() {
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bPair.MaxPrecision = maxPrecision

	modulo, errx = BigIntMathDivide{}.PairMod(bPair)

	if errx != nil {
		modulo = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.PairMod(bPair). "+
			"dividend='%v' divisor='%v' maxPrecision='%v' Error='%v'",
			bPair.Big1.GetNumStr(), bPair.Big2.GetNumStr(),
			bPair.MaxPrecision, errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// NumStrDtoModuloToNumStrDto - Performs a modulo operation on NumStrDto type
// input parameters, 'dividend' and 'divisor'. The result of this division
// operation is returned as a type NumStrDto.
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// This method returns one NumStrDto value: 'modulo'.
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned NumStrDto division result ('modulo') will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from input parameter,'dividend'.
//
// This method is different from method BigIntMathDivide.NumStrDtoModulo()
// above, in that this method returns 'modulo' as Type NumStrDto.
func (bIDivide BigIntMathDivide) NumStrDtoModuloToNumStrDto(
	dividend,
	divisor NumStrDto,
	maxPrecision uint) (modulo NumStrDto, err error) {

	ePrefix := "BigIntMathDivide.NumStrDtoModuloToNumStrDto() "

	modulo = NumStrDto{}.New()

	if divisor.IsZero() {
		err = fmt.Errorf(ePrefix + "Error: Attempted to mod by zero!")
		return modulo, err
	}

	bINumModulo, errx :=
		BigIntMathDivide{}.NumStrDtoModulo(dividend, divisor, maxPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntMathDivide{}.NumStrDtoModulo(dividend, divisor,"+
			" maxPrecision). Error='%v'", errx.Error())

		return modulo, err
	}

	modulo, errx = bINumModulo.GetNumStrDto()

	if errx != nil {
		modulo = NumStrDto{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by bINumModulo.GetNumStrDto(). Error='%v'",
			errx.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// PairQuotientMod -  Receives a BigIntPair type as an input parameter. 'BigIntPair.Big1'
// is treated as the Dividend. 'BigIntPair.Big2' is considered the Divisor.
// 'BigIntPair.maxPrecision' is used to control the precision of the resulting
// fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits. Therefore, the user
// is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// The method performs a division operation on BigIntNum input
// parameters 'dividend' (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2).
// The result is a quotient and modulo (remainder). These values are returned as
// two BigIntNum types, 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results (quotient and modulo) will
// contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from the the input parameter Dividend,
// 'BigIntPair.Big1'.
func (bIDivide BigIntMathDivide) PairQuotientMod(
	bPair BigIntPair) (quotient, modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.PairQuotientMod() "

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	var err2 error

	quotient, modulo, err2 = bIDivide.pairQuotientModNoNumSeps(bPair)

	if err2 != nil {
		quotient = BigIntNum{}.NewZero(0)
		modulo = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+"Error returned by bIDivide.pairQuotientModNoNumSeps(bPair). "+
			"Error='%v'", err2.Error())

		return quotient, modulo, err
	}

	err2 = quotient.SetNumericSeparatorsDto(numSeps)

	if err2 != nil {
		quotient = BigIntNum{}.NewZero(0)
		modulo = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+"Error returned by quotient.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v'", err2.Error())

		return quotient, modulo, err
	}

	err2 = modulo.SetNumericSeparatorsDto(numSeps)

	if err2 != nil {
		quotient = BigIntNum{}.NewZero(0)
		modulo = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+"Error returned by modulo.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v'", err2.Error())

		return quotient, modulo, err
	}

	err = nil

	return quotient, modulo, err
}

// PairMod - Receives a BigIntPair type as an input parameter. 'BigIntPair.Big1'
// is treated as the Dividend. 'BigIntPair.Big2' is considered the Divisor.
// The method proceeds to performs a modulo operation on input parameters 'dividend'
// (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2). The modulo result is returned
// as a type BigIntNum.
//
// 'BigIntPair.maxPrecision' is used to control the precision of the resulting
// fractional modulo returned by this method. Be advised that this method is capable
// of calculating modulo values with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division result 'modulo' will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from the the input parameter Dividend, 'BigIntPair.Big1'.
func (bIDivide BigIntMathDivide) PairMod(
	bPair BigIntPair) (modulo BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.PairMod() "

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	var err2 error

	modulo, err2 = bIDivide.pairModNoNumSeps(bPair)

	if err2 != nil {
		modulo = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by bIDivide.pairModNoNumSeps(bPair). "+
			"Error='%v' ", err2.Error())

		return modulo, err
	}

	err2 = modulo.SetNumericSeparatorsDto(numSeps)

	if err2 != nil {
		modulo = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by modulo.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v' ", err2.Error())

		return modulo, err
	}

	err = nil

	return modulo, err
}

// PairIntQuotient - Receives a BigIntPair type as an input parameter. 'BigIntPair.Big1'
// is treated as the Dividend. 'BigIntPair.Big2' is considered the Divisor. This method
// performs integer division on input parameters, 'Dividend' (BigIntPair.Big1) and
// 'Divisor' (BigIntPair.Big2).
//
// The result of this division operation returns an integer quotient of Dividend
// (BigIntPair.Big1) divided by Divisor (BigIntPair.Big2).
//
// The division operation performed by this method is T-Division or truncated division.
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// After completing the division operation, an integer quotient of type BigIntNum
// is returned.
//
// Examples:
// =========
//
//																					Return Value
//	 Dividend divided by	Divisor			=		Integer Quotient
//			 5 				/ 				 2 				= 				 2
//	    5.25 		/ 				 2  			= 				 2
//	    2 				/ 				 4				= 				 0
//			-5 				/ 				 2 				= 				-2
//	   -5.25     /    			 2  			= 				-2
//	   -2 				/ 				 4				= 				 0
//			 5 				/ 				-2 				=					-2
//	    5.25 		/ 				-2 				= 				-2
//	    2 				/ 				-4				= 				 0
//			-5 				/ 				-2 				= 				 2
//	   -5.25     /    			-2 				= 				 2
//	   -2 				/ 				-4				= 				 0
//	    12.555		/ 			  -2.5			=			    -5
//	   -12.555		/ 			  -2.5			=			     5
//	    12.555		/ 			  -2				=			    -6
//
// The returned BigIntNum division result 'intQuotient' will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from the input parameter Dividend, 'BigIntPair.Big1'.
func (bIDivide BigIntMathDivide) PairIntQuotient(bPair BigIntPair) (intQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.PairIntQuotient() "

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	var err2 error

	intQuotient, err2 = bIDivide.pairIntQuotientNoNumSeps(bPair)

	if err2 != nil {
		intQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by bIDivide.pairIntQuotientNoNumSeps(bPair). "+
			"Error='%v' ", err2.Error())

		return intQuotient, err
	}

	err2 = intQuotient.SetNumericSeparatorsDto(numSeps)

	if err2 != nil {
		intQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by intQuotient.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v' ", err2.Error())

		return intQuotient, err
	}

	err = nil

	return intQuotient, err
}

// PairFracQuotient - Receives a BigIntPair type as an input parameter.
// 'BigIntPair.Big1' is treated as the Dividend. 'BigIntPair.Big2' is considered
// the divisor.
//
// 'BigIntPair.maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// This method performs a division operation on BigIntNum parameters 'dividend'
// (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2).
//
//	Dividend (BigIntPair.Big1) divided Divisor (BigIntPair.Big2) = quotient
//
// The resulting quotient is returned as a BigIntNum type representing the result
// of the division operation expressed as integer and fractional digits. The
// maximum number of fractional digits output to the result is controlled by
// BigIntPair.maxPrecision. Remember that the BigIntNum type specifies 'precision'.
// Precision is defined as the number of fractional digits to the right of the
// decimal place.
//
// Examples:
// =========
//
// Note: For all examples BigIntPair.maxPrecision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					   Quotient
//	 Dividend		divided by	Divisor		=		BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result 'fracQuotient' will contain numeric
// separators (decimal separator, thousands separator and currency symbol) copied
// from the input parameter Dividend, 'BigIntPair.Big1'.
func (bIDivide BigIntMathDivide) PairFracQuotient(
	bPair BigIntPair) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.pairFracQuotientNoNumSeps() "

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	var err2 error

	fracQuotient, err2 = bIDivide.pairFracQuotientNoNumSeps(bPair)

	if err2 != nil {
		fracQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by bIDivide.pairFracQuotientNoNumSeps(bPair). "+
			"Error='%v' ", err2.Error())

		return fracQuotient, err
	}

	err2 = fracQuotient.SetNumericSeparatorsDto(numSeps)

	if err2 != nil {
		fracQuotient = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned by fracQuotient.SetNumericSeparatorsDto(numSeps). "+
			"Error='%v' ", err2.Error())

		return fracQuotient, err
	}

	err = nil

	return fracQuotient, err
}

// pairFracQuotientNoNumSeps - Receives a BigIntPair type as an input parameter.
// 'BigIntPair.Big1' is treated as the Dividend. 'BigIntPair.Big2' is considered
// the divisor.
//
// 'BigIntPair.maxPrecision' is used to control the maximum precision of the
// resulting fractional quotient. Be advised that this method is capable of
// calculating quotients with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// This method performs a division operation on BigIntNum parameters 'dividend'
// (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2).
//
//	Dividend (BigIntPair.Big1) divided Divisor (BigIntPair.Big2) = quotient
//
// The resulting quotient is returned as a BigIntNum type representing the result
// of the division operation expressed as integer and fractional digits. The
// maximum number of fractional digits output to the result is controlled by
// BigIntPair.maxPrecision. Remember that the BigIntNum type specifies 'precision'.
// Precision is defined as the number of fractional digits to the right of the
// decimal place.
//
// Examples:
// =========
//
// Note: For all examples BigIntPair.maxPrecision is specified as '15'.
// ----------------------------------------------------------------------------
//
//																					   Quotient
//	 Dividend		divided by	Divisor		=		BigIntNum Integer 	Precision	 Result
//	 -------- 	  ----------	--------				-----------------	  ---------	 ------
//		 10.5  				/ 				2 				= 			525  							  2  			 5.25
//		 10    				/ 				2 				= 			5	  							  0  			 5
//	  11.5  				/         2.5				=  			46								  1				 4.6
//	   2.5					/				 12.555			=				199123855037834	   15				 0.199123855037834
//		-12.555 			/ 				2.5 			= 		 -5022							  3				-5.022
//	 -12.555     	/    			2  			  = 		 -62775							  4				-6.2775
//	 - 2.5 				/ 			 12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		 12.555				/ 			- 2.5			  =			 -5022								3				-5.022
//	  12.555 			/ 			- 2 				= 		 -62775								4				-6.2775
//	   2.5 				/ 			-12.555		  = 		 -199123855037834	   15				-0.199123855037834
//		-12.555 			/ 			- 2.5 			= 			5022								3				 5.022
//	 -12.555     	/    		- 2 				= 		  62775								4				 6.2775
//	 - 2.5	 				/ 			-12.555		  = 		  199123855037834	   15				 0.199123855037834
//	 -10						/				- 2					=				5														 5
//
// The returned BigIntNum division result 'fracQuotient' will contain default
// numeric separators (decimal separator, thousands separator and currency
// symbol).
func (bIDivide BigIntMathDivide) pairFracQuotientNoNumSeps(
	bPair BigIntPair) (fracQuotient BigIntNum, err error) {

	ePrefix := "BigIntMathDivide.pairFracQuotientNoNumSeps() "
	fracQuotient = BigIntNum{}.New()

	if bPair.Big2.bigInt.Cmp(big.NewInt(0)) == 0 {
		err = errors.New(ePrefix + "Attempted Divide by ZERO!")
		return fracQuotient, err
	}

	bPair.MakePrecisionsEqual()

	rDividend := big.NewRat(1, 1).SetInt(bPair.Big1.bigInt)
	rDivisor := big.NewRat(1, 1).SetInt(bPair.Big2.bigInt)
	rQuotient := big.NewRat(1, 1).Quo(rDividend, rDivisor)
	numStr := rQuotient.FloatString(int(bPair.MaxPrecision))

	fracQuotient, errx :=
		BigIntNum{}.NewNumStr(numStr)

	if errx != nil {
		fracQuotient = BigIntNum{}.New()
		err = fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStr(numStr). "+
			"numStr='%v' maxPrecision='%v' Error='%v'",
			numStr, bPair.MaxPrecision, errx.Error())

		return fracQuotient, err
	}

	fracQuotient.TrimTrailingFracZeros()

	fracQuotient.SetNumericSeparatorsToDefaultIfEmpty()

	err = nil

	return fracQuotient, err
}

// pairIntQuotientNoNumSeps - Performs integer division on two BigIntNum types passed as
// BigIntPair input parameters. 'dividend' is BigIntPair.Big1 and 'divisor' is
// BigIntPair.Big2.
//
// The result of this division operation returns an integer quotient of Dividend
// (BigIntPair.Big1) divided by Divisor (BigIntPair.Big2).
//
// The division operation performed by this method is T-Division or truncated division.
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// After completing the division operation, an integer quotient of type BigIntNum
// is returned.
//
// Examples:
// =========
//
//																					Return Value
//	 Divisor	divided by	Dividend		=		Integer Quotient
//			 5 				/ 				 2 				= 				 2
//	    5.25 		/ 				 2  			= 				 2
//	    2 				/ 				 4				= 				 0
//			-5 				/ 				 2 				= 				-2
//	   -5.25     /    			 2  			= 				-2
//	   -2 				/ 				 4				= 				 0
//			 5 				/ 				-2 				=					-2
//	    5.25 		/ 				-2 				= 				-2
//	    2 				/ 				-4				= 				 0
//			-5 				/ 				-2 				= 				 2
//	   -5.25     /    			-2 				= 				 2
//	   -2 				/ 				-4				= 				 0
//	    12.555		/ 			  -2.5			=			    -5
//	   -12.555		/ 			  -2.5			=			     5
//	    12.555		/ 			  -2				=			    -6
//
// The returned BigIntNum division result 'intQuotient' will contain default
// numeric separators (decimal separator, thousands separator and currency
// symbol).
func (bIDivide BigIntMathDivide) pairIntQuotientNoNumSeps(
	bPair BigIntPair) (intQuotient BigIntNum, err error) {

	intQuotient = BigIntNum{}.New()

	if bPair.Big2.IsZero() {
		ePrefix := "BigIntMathDivide.pairIntQuotientNoNumSeps() "
		err = fmt.Errorf(ePrefix + "Error: Attempted to divide by zero!")
		return intQuotient, err
	}

	bPair.MakePrecisionsEqual()

	bigIQuotient := big.NewInt(0).Quo(bPair.Big1.bigInt, bPair.Big2.bigInt)

	intQuotient = BigIntNum{}.NewBigInt(bigIQuotient, 0)

	intQuotient.SetNumericSeparatorsToDefaultIfEmpty()

	err = nil

	return intQuotient, err
}

// pairModNoNumSeps - Receives a BigIntPair type as an input parameter. 'BigIntPair.Big1'
// is treated as the Dividend. 'BigIntPair.Big2' is considered the Divisor.
// The method proceeds to performs a modulo operation on input parameters 'dividend'
// (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2). The modulo result is returned
// as a type BigIntNum.
//
// 'BigIntPair.maxPrecision' is used to control the precision of the resulting
// fractional modulo returned by this method. Be advised that this method is capable
// of calculating modulo values with very long strings of fractional digits. Therefore,
// the user is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// The modulo operation finds the remainder after division of one number
// by another (sometimes called modulus).
//
//	Wikipedia https://en.wikipedia.org/wiki/Modulo_operation
//
// The calculation of 'modulo' is based on T-Division (Truncate Division). See
// "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of
// Utrecht Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
//
// So, for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	  						q = D div d = f(D/d)
//								r = D mod d = D − d ·q
//
// The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Examples:
// =========
//
// Dividend			  mod by			Divisor			=			Modulo/Remainder
// --------				------			-------						----------------
//
//	  12.555					%						 2.5			=			 0.055
//	  12.555  	 			% 				 	 2  			= 		 0.555
//	   2.5 					% 				 	12.555		= 	   2.5
//		-12.555 				% 				   2.5 			= 		-0.055
//	 -12.555     		%    			 	 2  			= 		-0.555
//	 - 2.5 					% 				 	12.555		= 		-2.5
//		 12.555					% 				 - 2.5			=			 0.055
//	  12.555 				% 				 - 2 				= 		 0.555
//	   2.5 				  % 				 -12.555		= 		 2.5
//		-12.555 				% 				 - 2.5 			= 		-0.055
//	 -12.555     		%    			 - 2 				= 		-0.555
//	 - 2.5	 					% 				 -12.555		= 		-2.5
//
// The returned BigIntNum division results (quotient and modulo) will
// contain default numeric separators (decimal separator, thousands
// separator and currency symbol).
func (bIDivide BigIntMathDivide) pairModNoNumSeps(bPair BigIntPair) (modulo BigIntNum, err error) {

	modulo = BigIntNum{}.New()

	if bPair.Big2.bigInt.Cmp(big.NewInt(0)) == 0 {
		ePrefix := "BigIntMathDivide.pairModNoNumSeps() "
		err = errors.New(ePrefix + "Error: Attempted divide by ZERO!")
		return modulo, err
	}

	bPair.MakePrecisionsEqual()

	moduloBigI := big.NewInt(0).Rem(bPair.Big1.bigInt, bPair.Big2.bigInt)

	modulo = BigIntNum{}.NewBigInt(moduloBigI, bPair.Big2.GetPrecisionUint())

	if modulo.precision > bPair.MaxPrecision {
		modulo.RoundToDecPlace(bPair.MaxPrecision)
	}

	modulo.TrimTrailingFracZeros()

	modulo.SetNumericSeparatorsToDefaultIfEmpty()

	err = nil

	return modulo, err
}

// pairQuotientModNoNumSeps -  Receives a BigIntPair type as an input parameter. 'BigIntPair.Big1'
// is treated as the Dividend. 'BigIntPair.Big2' is considered the Divisor.
// 'BigIntPair.maxPrecision' is used to control the precision of the resulting
// fractional quotient. Be advised that this method is capable of calculating
// quotients with very long strings of fractional digits. Therefore, the user
// is advised to set a relevant 'BigIntPair.maxPrecision' value.
//
//	type BigIntPair struct {
//				Big1							BigIntNum  // The Dividend
//				Big2							BigIntNum	 // The Divisor
//				maxPrecision			uint			 // Controls Precision
//	}
//
// The method performs a division operation on BigIntNum input
// parameters 'dividend' (BigIntPair.Big1) and 'divisor' (BigIntPair.Big2).
// The result is a quotient and modulo (remainder). These values are returned as
// two BigIntNum types, 'quotient' and 'modulo'.
//
// The calculation of 'quotient' and 'modulo' is based on T-Division (Truncate Division).
// See "Division and Modulus for Computer Scientists", DAAN LEIJEN, University of Utrecht
// Dept. of Computer Science, PO.Box 80.089, 3508 TB Utrecht The Netherlands:
// https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote-letter.pdf
// Also available at ../notes/divmodnote-letter.pdf.
// So for q=quotient; D=Dividend d=Divisor r=Remainder or 'modulo' :
//
//	q = D div d = f(D/d) r = D mod d = D − d ·q
//
// 'quotient' is the integer result of dividing the 'dividend' by the 'divisor'
//
// 'modulo' - The modulo operation finds the remainder after division of one
// number by another. (r = D mod d = D − d ·q)
//
// Input parameter 'maxPrecision' is used to control the maximum precision of the
// resulting 'modulo'. Precision is defined as the the number of fractional digits
// to the right of the decimal point. Be advised that these calculations can support
// very large precision values.
//
// Examples:
// =========
//
// Dividend			divided by		Divisor			=		Quotient			Modulo/Remainder
//
//	  12.555					/						 2.5			=			 5							 0.055
//	  12.555  	 			/ 				 	 2  			= 		 6							 0.555
//	   2.5 					/ 				 	12.555		= 	   0							 2.5
//		-12.555 				/ 				   2.5 			= 		-5							-0.055
//	 -12.555     		/    			 	 2  			= 		-6							-0.555
//	 - 2.5 					/ 				 	12.555		= 		 0							-2.5
//		 12.555					/ 				 - 2.5			=			-5							 0.055
//	  12.555 				/ 				 - 2 				= 		-6							 0.555
//	   2.5 				  / 				 -12.555		= 		 0							 2.5
//		-12.555 				/ 				 - 2.5 			= 		 5							-0.055
//	 -12.555     		/    			 - 2 				= 		 6							-0.555
//	 - 2.5	 					/ 				 -12.555		= 		 0							-2.5
//
// The returned BigIntNum division results (quotient and modulo) will
// contain default numeric separators (decimal separator, thousands
// separator and currency symbol).
func (bIDivide BigIntMathDivide) pairQuotientModNoNumSeps(
	bPair BigIntPair) (quotient, modulo BigIntNum, err error) {

	quotient = BigIntNum{}.New()
	modulo = BigIntNum{}.New()

	if bPair.Big2.bigInt.Cmp(big.NewInt(0)) == 0 {

		ePrefix := "BigIntMathDivide.pairQuotientModNoNumSeps() "

		err = errors.New(ePrefix +
			"Error: Attempted Divide By ZERO!")

		return quotient, modulo, err
	}

	bPair.MakePrecisionsEqual()

	scratch := big.NewInt(0)

	quotientBigI, moduloBigI := big.NewInt(0).QuoRem(
		bPair.Big1.bigInt,
		bPair.Big2.bigInt,
		scratch)

	quotient = BigIntNum{}.NewBigInt(quotientBigI, 0)

	modulo = BigIntNum{}.NewBigInt(moduloBigI, bPair.Big2.GetPrecisionUint())

	if modulo.precision > bPair.MaxPrecision {
		modulo.RoundToDecPlace(bPair.MaxPrecision)
	}

	modulo.TrimTrailingFracZeros()

	quotient.SetNumericSeparatorsToDefaultIfEmpty()

	modulo.SetNumericSeparatorsToDefaultIfEmpty()

	err = nil

	return quotient, modulo, err
}
