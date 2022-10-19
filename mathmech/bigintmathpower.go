package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

type BigIntMathPower struct {
	Base     BigIntNum
	Exponent BigIntNum
	Result   BigIntNum
}

// BigIntPwr - Raises input parameter 'base' to the power of input parameter
// 'exponent'. Both 'base' and 'exponent' may be positive or negative integer
// or fractional values.
//
// This method uses the exponent method ('Exp') provided by the go "math/big" package.
//
//
func (bIPwr BigIntMathPower) BigIntPwr(
	base,
	basePrecision,
	exponent,
	exponentPrecision,
	maxPrecision *big.Int) (result *big.Int, resultPrecision *big.Int, err error) {

	result = big.NewInt(0)
	resultPrecision = big.NewInt(0)
	err = nil
	ePrefix := "BigIntMathPower.BigIntPwr() "
	var errx error

	bigZero := big.NewInt(0)

	if base.Cmp(bigZero) == 0 {
		// zero to any power is zero
		return result, resultPrecision, err
	}

	isZeroExponentPrecision := false

	if exponentPrecision.Cmp(bigZero) == 0 {
		isZeroExponentPrecision = true
	}

	exponentCmpZero := exponent.Cmp(bigZero)

	if exponentCmpZero == 0 {
		result = big.NewInt(1)
		return result, resultPrecision, err
	}

	exponentIsNegative := false

	if exponentCmpZero == -1 {
		exponentIsNegative = true
	}

	if exponentIsNegative == false &&
		isZeroExponentPrecision == true {

		result, resultPrecision, errx =
			BigIntMathPower{}.BigIntToPositiveIntegerPower(
				base,
				basePrecision,
				exponent,
				exponentPrecision,
				maxPrecision)

	} else if exponentIsNegative == true &&
		isZeroExponentPrecision == true {

		result, resultPrecision, errx =
			BigIntMathPower{}.BigIntToNegativeIntegerPower(
				base,
				basePrecision,
				exponent,
				exponentPrecision,
				maxPrecision)

	} else if exponentIsNegative == false &&
		isZeroExponentPrecision == false {

		result, resultPrecision, errx =
			BigIntMathPower{}.BigIntToPositiveFractionalPower(
				base,
				basePrecision,
				exponent,
				exponentPrecision,
				maxPrecision)

	} else if exponentIsNegative == true &&
		isZeroExponentPrecision == false {

		result, resultPrecision, errx =
			BigIntMathPower{}.BigIntToNegativeFractionalPower(
				base,
				basePrecision,
				exponent,
				exponentPrecision,
				maxPrecision)
	} else {
		errx = errors.New("Setup Configuration Error! INVALID Calculation!")
	}

	if errx != nil {
		result = big.NewInt(0)
		resultPrecision = big.NewInt(0)
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
	} else {
		err = nil
	}

	return result, resultPrecision, err
}

// BigIntToNegativeFractionalPower - Raises 'base' to the power of a
// negative fractional exponent, 'exponent'.
//
// As stated in the function name, this method expects to process only
// negative exponent values which have fractional digits to the right
// of the decimal place.
//
// Examples:
//
// base	basePrecision	exponent exponentPrecision result  resultPrecision
// ---- ------------- -------- ----------------- ------  ---------------
//
// 	  5			0					  -22						1						 29			 3 (to 3-decimal places)
//    (5^-2.2 = 0.02899118654710782125882456003526
//     The actual number of decimal places returned in the result is controlled by
//     input parameter, 'maxPrecision'.)
//
//   18     1            -34           1					 136		 3 (to 3-decimal places)
//  (1.8^-3.4 = 0.13554187298692911221722484380209
//     The actual number of decimal places returned in the result is controlled by
//     input parameter, 'maxPrecision'.)
//
// Input Parameters
// ================
//
// base 					*big.Int	- The base which will be raised to the
//                          	power of a positive integer exponent.
//
// basePrecision 			uint	- The precision specification for 'base'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// exponent				*big.Int	- The exponent to which 'base' will be raised
//                          	by this calculation. By method definition,
//                            'exponent' must be a negative value. If
//                            exponent is less than 0, an error will
//                            be triggered.
//
// exponentPrecision 	uint	- The precision specification for 'exponent'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place. For this method, 'exponentPrecision'
// 														must be greater than zero, thereby designating
// 														'exponent' as a fractional value. A value of
//                            zero for 'exponentPrecision' will trigger an error.
//
// maxPrecision				uint	- When this method calculates 'base' raised to the
//                            power of 'exponent', the maximum number of decimal
//                            digits to the right of the decimal place in the
//                            result will be limited by 'maxPrecision'.
//
// Return Values
// =============
//
// result					*big.Int	- If the calculation completes successfully, this
//                            return value will be populated with the value
//                            of 'base' raised to the power of 'exponent'.
//
// resultPrecision		uint	- The precision specification for 'result'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'result' and 'resultPrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// err							 error	- If the calculation encounters an error, an appropriate
//                            error message will be formatted and returned. If the
//                            calculation completes successfully, this return value
//                            will be set to 'nil'.
//
func (bIPwr BigIntMathPower) BigIntToNegativeFractionalPower(
	base,
	basePrecision,
	exponent,
	exponentPrecision,
	maxPrecision *big.Int) (result *big.Int, resultPrecision *big.Int, err error) {

	result = big.NewInt(0)
	resultPrecision = big.NewInt(0)
	err = nil

	ePrefix := "BigIntMathPower.BigIntToNegativeFractionalPower() "
	bigZero := big.NewInt(0)
	bigOne := big.NewInt(1)

	if base.Cmp(bigZero) == 0 {
		// base is zero result is zero
		return result, resultPrecision, err
	}

	if exponentPrecision.Cmp(bigZero) == 0 {
		err = fmt.Errorf(ePrefix+
			"Error: 'exponentPrecision' is zero. This is an integer exponent. "+
			"exponentPrecision='%v' ", exponentPrecision.Text(10))

		return result, resultPrecision, err
	}

	// exponentPrecision > 0
	cmpExponentZero := exponent.Cmp(bigZero)

	if cmpExponentZero == 1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'exponent' is a positive value! "+
			"exponent='%v' ", exponent.Text(10))

		return result, resultPrecision, err

	}

	if cmpExponentZero == 0 {
		// Any number raised to a zero power is one
		result = big.NewInt(1)
		return result, resultPrecision, err
	}

	if maxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'maxPrecision' is a negative value! "+
			"maxPrecision='%v' ", maxPrecision.Text(10))

		return result, resultPrecision, err
	}

	positiveExponent := big.NewInt(0).Set(exponent)
	positiveExponent.Neg(positiveExponent)

	internalMaxPrecision := big.NewInt(0).Add(maxPrecision, big.NewInt(2))

	positiveResult, positivePrecision, errx :=
		BigIntMathPower{}.BigIntToPositiveFractionalPower(
			base,
			basePrecision,
			positiveExponent,
			exponentPrecision,
			internalMaxPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return result, resultPrecision, err
	}

	result, resultPrecision, errx = BigIntMathDivide{}.BigIntFracQuotient(
		bigOne,
		big.NewInt(0),
		positiveResult,
		positivePrecision,
		maxPrecision)

	if errx != nil {
		result = big.NewInt(0)
		resultPrecision = big.NewInt(0)
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return result, resultPrecision, err
	}

	return result, resultPrecision, err
}

// BigIntToPositiveFractionalPower - Raises 'base' to the power
// of a positive fractional exponent, 'exponent'.
//
// As stated in the function name, this method expects to process only
// positive exponent values which have fractional digits to the right
// of the decimal place.
//
// Examples:
// base	basePrecision	exponent exponentPrecision result  resultPrecision
//
// 	  5			0					   22						1						 34493			 3 (to 3-decimal places)
//    (5^2.2 = 34.4932415365303708097515866054
//     The actual number of decimal places returned in the result is controlled by
//     input parameter, 'maxPrecision'.)
//
//   18     1            34            1					 7378				 3 (to 3-decimal places)
//  (1.8^3.4 = 7.3777938725727533349996174917827
//     The actual number of decimal places returned in the result is controlled by
//     input parameter, 'maxPrecision'.)
//
// Input Parameters
// ================
//
// base 					*big.Int	- The base which will be raised to the
//                          	power of a positive integer exponent.
//
// basePrecision 			uint	- The precision specification for 'base'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// exponent				*big.Int	- The exponent to which 'base' will be raised
//                          	by this calculation. By method definition,
//                            'exponent' must be a positive value. If
//                            exponent is less than 0, an error will
//                            be triggered.
//
// exponentPrecision 	uint	- The precision specification for 'exponent'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place. For this method, 'exponentPrecision'
// 														must be greater than zero, thereby designating
// 														'exponent' as a fractional value. A value of
// 														zero for 'exponentPrecision' will trigger an error.
//
// maxPrecision				uint	- When this method calculates 'base' raised to the
//                            power of 'exponent', the maximum number of decimal
//                            digits to the right of the decimal place in the
//                            result will be limited by 'maxPrecision'.
//
// Return Values
// =============
//
// result					*big.Int	- If the calculation completes successfully, this
//                            return value will be populated with the value
//                            of 'base' raised to the power of 'exponent'.
//
// resultPrecision		uint	- The precision specification for 'result'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'result' and 'resultPrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// err							 error	- If the calculation encounters an error, an appropriate
//                            error message will be formatted and returned. If the
//                            calculation completes successfully, this return value
//                            will be set to 'nil'.
//
func (bIPwr BigIntMathPower) BigIntToPositiveFractionalPower(
	base,
	basePrecision,
	exponent,
	exponentPrecision,
	maxPrecision *big.Int) (result *big.Int, resultPrecision *big.Int, err error) {

	result = big.NewInt(0)
	resultPrecision = big.NewInt(0)
	err = nil

	ePrefix := "BigIntMathPower.BigIntToPositiveFractionalPower() "

	bigZero := big.NewInt(0)

	bigOne := big.NewInt(1)

	internalMaxPrecision := big.NewInt(0).Mul(maxPrecision, big.NewInt(2))

	if base.Cmp(bigZero) == 0 {
		// base is zero result is zero
		return result, resultPrecision, err
	}

	basePrecisionCmpZero := basePrecision.Cmp(bigZero)

	if base.Cmp(bigOne) == 0 &&

		basePrecisionCmpZero == 0 {

		result.Set(bigOne)

		return result, resultPrecision, err
	}

	if basePrecisionCmpZero == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'basePrecision' is negative! "+
			"basePrecision='%v' ", basePrecision)

		return result, resultPrecision, err

	}

	exponentPrecisionCmpZero := exponentPrecision.Cmp(bigZero)

	if exponentPrecisionCmpZero == 0 {
		err = fmt.Errorf(ePrefix+
			"Error: 'exponentPrecision' is zero. This is an integer exponent. "+
			"exponentPrecision='%v' ", exponentPrecision)

		return result, resultPrecision, err
	}

	if exponentPrecisionCmpZero == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'exponentPrecision' is a negative value! "+
			"exponentPrecision='%v' ", exponentPrecision)

		return result, resultPrecision, err
	}

	// exponentPrecision > 0
	cmpExponentZero := exponent.Cmp(bigZero)

	if cmpExponentZero == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'exponent' is a negative value! "+
			"exponent='%v' ", exponent.Text(10))

		return result, resultPrecision, err
	}

	if cmpExponentZero == 0 {
		// Any number raised to a zero power is one
		result = big.NewInt(1)
		return result, resultPrecision, err
	}

	bigTen := big.NewInt(10)
	// Get exponent integer value
	scale := big.NewInt(0).Exp(
		bigTen,
		exponentPrecision,
		nil)

	var errx error

	scratch := big.NewInt(0)

	integerExponent, fractionalExponent := big.NewInt(0).QuoRem(exponent, scale, scratch)

	integerResult := big.NewInt(1)
	integerPrecision := big.NewInt(0)

	if integerExponent.Cmp(bigZero) != 0 {

		integerResult, integerPrecision, errx =
			BigIntMathPower{}.BigIntToPositiveIntegerPower(
				base,
				basePrecision,
				integerExponent,
				big.NewInt(0),
				internalMaxPrecision)

		if errx != nil {
			err = fmt.Errorf(ePrefix+"%v", errx.Error())
			return result, resultPrecision, err
		}

	}

	ratFrac := big.NewRat(1, 1).SetFrac(fractionalExponent, scale)

	ratFracExponentNumerator := ratFrac.Num()

	baseToFracExponentNumerator := big.NewInt(0).Exp(base, ratFracExponentNumerator, nil)
	baseToFracExponentNumeratorPrecision := big.NewInt(0).Mul(basePrecision, ratFracExponentNumerator)

	delta := big.NewInt(0)
	bigFive := big.NewInt(5)

	ratFracExponentDenominator := ratFrac.Denom()

	fdr := FixedDecimalNthRoot{}

	fracExponentRoot, fracExponentRootPrecision, errx :=
		fdr.CalculatePositiveIntegerNthRoot(
			baseToFracExponentNumerator,
			baseToFracExponentNumeratorPrecision,
			ratFracExponentDenominator,
			big.NewInt(0),
			internalMaxPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())
		return result, resultPrecision, err
	}

	result, resultPrecision, errx =
		BigIntMathMultiply{}.BigIntMultiply(
			integerResult,
			integerPrecision,
			fracExponentRoot,
			fracExponentRootPrecision)

	if errx != nil {
		result = big.NewInt(0)
		resultPrecision = big.NewInt(0)
		err = fmt.Errorf(ePrefix+
			"%v", errx.Error())

		return result, resultPrecision, err
	}

	bigFive = big.NewInt(5)

	if resultPrecision.Cmp(maxPrecision) == 1 {
		delta = big.NewInt(0).Sub(resultPrecision, maxPrecision)
		delta.Sub(delta, bigOne)
		scale.Exp(big.NewInt(10), delta, nil)
		result.Quo(result, scale)
		if result.Cmp(bigZero) == -1 {
			bigFive.Neg(bigFive)
		}
		result.Add(result, bigFive)
		result.Quo(result, bigTen)
		resultPrecision = big.NewInt(0).Set(maxPrecision)
	}

	return result, resultPrecision, err
}

// BigIntToNegativeIntegerPower - Raises 'base' to the power of
// a negative integer exponent, 'exponent'.
//
// As stated in the function name, this method expects to
// process only negative, integer exponents.
//
// Examples:
// base	basePrecision	exponent exponentPrecision result  resultPrecision
//
// 	  5			0					  -2						0							4						2
//    (5^-2 = 0.04)
//
// 1131     2           -3            0						00069121  		8 (to 8-decimal places)
// (11.31^-3= 0.00069121345785745610965099525880031
//  The actual number of decimal places returned in
//  the result is controlled by	input parameter, 'maxPrecision'.)
//
// Input Parameters
// ================
//
// base 					*big.Int	- The base which will be raised to the
//                          	power of a negative integer exponent.
//
// basePrecision 			uint	- The precision specification for 'base'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// exponent				*big.Int	- The exponent to which 'base' will be raised
//                          	by this calculation. By method definition,
//                            exponent must be a negative value. If
//                            exponent is greater than -1, an error will
//                            be triggered.
//
// exponentPrecision 	uint	- The precision specification for 'exponent'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place. For this method, 'exponentPrecision' must
//                            be set to zero, thereby designating 'exponent'
// 														as an integer value. Values greater than zero
// 														will trigger an error.
//
// maxPrecision				uint	- When this method calculates 'base' raised to the
//                            power of 'exponent', the maximum number of decimal
//                            digits to the right of the decimal place in the
//                            resulting value will be limited by 'maxPrecision'.
//
// Return Values
// =============
//
// result					*big.Int	- If the calculation completes successfully, this
//                            return value will be populated with the value
//                            of 'base' raised to the power of 'exponent'.
//
// resultPrecision		uint	- The precision specification for 'result'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'result' and 'resultPrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// err							 error	- If the calculation encounters an error, an appropriate
//                            error message will be formatted and returned. If the
//                            calculation completes successfully, this return value
//                            will be set to 'nil'.
//
func (bIPwr BigIntMathPower) BigIntToNegativeIntegerPower(
	base,
	basePrecision,
	exponent,
	exponentPrecision,
	maxPrecision *big.Int) (result *big.Int, resultPrecision *big.Int, err error) {

	ePrefix := "BigIntMathPower) BigIntToNegativeIntegerPower() "
	result = big.NewInt(0)
	resultPrecision = big.NewInt(0)
	err = nil

	bigOne := big.NewInt(1)

	bigZero := big.NewInt(0)

	if base.Cmp(bigZero) == 0 {
		// base is zero result is zero
		return result, resultPrecision, err
	}

	cmpExponentZero := exponent.Cmp(bigZero)

	if cmpExponentZero == 1 {
		err = fmt.Errorf(ePrefix+"Error: Input parameter 'exponent' is positive! "+
			"exponent='%v' ", exponent.Text(10))

		return result, resultPrecision, err
	}

	if exponentPrecision.Cmp(bigZero) == 1 {
		err = fmt.Errorf(ePrefix+"Error: Input parameter 'exponent' is NOT an integer! "+
			"exponentPrecision='%v' ", exponentPrecision.Text(10))

		return result, resultPrecision, err
	}

	if maxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+"Error: Input parameter 'maxPrecision' is a negative value! "+
			"exponentPrecision='%v' ", maxPrecision.Text(10))

		return result, resultPrecision, err

	}

	if cmpExponentZero == 0 {
		// Any number raised to a zero power is one
		result = big.NewInt(1)
		return result, resultPrecision, err
	}

	var errx error

	if exponent.Cmp(big.NewInt(-1)) == 0 {
		// if exponent == -1, result is equal to inverse of base
		result, resultPrecision, errx =
			BigIntMathDivide{}.BigIntFracQuotient(bigOne,
				big.NewInt(0),
				base,
				basePrecision,
				maxPrecision)

		if errx != nil {
			err = fmt.Errorf(ePrefix+"%v", errx.Error())

			return result, resultPrecision, err
		}

		err = nil

		return result, resultPrecision, err
	}

	tempExponent := big.NewInt(0).Set(exponent)

	// Exponent now positive integer value
	tempExponent.Neg(tempExponent)

	tempResult := big.NewInt(0).Exp(base, tempExponent, nil)

	tempResultPrecision := big.NewInt(0).Mul(
		basePrecision, tempExponent)

	result, resultPrecision, errx =
		BigIntMathDivide{}.BigIntFracQuotient(
			bigOne,
			big.NewInt(0),
			tempResult,
			tempResultPrecision,
			maxPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+"%v", errx.Error())

		return result, resultPrecision, err
	}

	err = nil

	return result, resultPrecision, err
}

// BigIntToPositiveIntegerPower - Raises 'base' to the power of a
// positive integer exponent, 'exponent'.
//
// As stated in the function name, this method expects to
// process only positive, integer exponents.
//
// Examples:
// base	basePrecision	exponent exponentPrecision result  resultPrecision
//
// 	  5			0					   2						0						 25						0
//    (5^2 = 25)
//
//   18     1            3            0					 5832						3
// 		(1.8^3= 5.832)
//
// Input Parameters
// ================
//
// base 					*big.Int	- The base which will be raised to the
//                          	power of a positive integer exponent.
//
// basePrecision 			uint	- The precision specification for 'base'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// exponent				*big.Int	- The exponent to which 'base' will be raised
//                          	by this calculation. By method definition,
//                            'exponent' must be a positive value. If
//                            exponent is less than 0, an error will
//                            be triggered.
//
// exponentPrecision 	uint	- The precision specification for 'exponent'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'base' and 'basePrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place. For this method, 'exponentPrecision'
// 														must be set to zero, designating exponent as an integer
//                            value. Values greater than zero will trigger an
//                            error.
//
// maxPrecision				uint	- When this method calculates 'base' raised to the
//                            power of 'exponent', the maximum number of decimal
//                            digits to the right of the decimal place in the
//                            resulting value will be limited by 'maxPrecision'.
//
// Return Values
// =============
//
// result					*big.Int	- If the calculation completes successfully, this
//                            return value will be populated with the value
//                            of 'base' raised to the power of 'exponent'.
//
// resultPrecision		uint	- The precision specification for 'result'.
//                          	Precision defines the number of numeric
//                          	digits to the right of the decimal place.
//                          	Taken together, 'result' and 'resultPrecision'
//                          	define a numeric value with a fixed number
//                            of decimal digits to the right of the decimal
//                            place.
//
// err							 error	- If the calculation encounters an error, an appropriate
//                            error message will be formatted and returned. If the
//                            calculation completes successfully, this return value
//                            will be set to 'nil'.
//
func (bIPwr BigIntMathPower) BigIntToPositiveIntegerPower(
	base,
	basePrecision,
	exponent,
	exponentPrecision,
	maxPrecision *big.Int) (result *big.Int, resultPrecision *big.Int, err error) {

	result = big.NewInt(0)
	resultPrecision = big.NewInt(0)
	err = nil

	ePrefix := "BigIntMathPower.BigIntToPositiveIntegerPower() "
	bigZero := big.NewInt(0)

	if base.Cmp(bigZero) == 0 {
		// base is zero result is zero
		return result, resultPrecision, err
	}

	cmpExponentZero := exponent.Cmp(bigZero)

	if cmpExponentZero == 0 {
		// Any number raised to a zero power is one
		result = big.NewInt(1)
		return result, resultPrecision, err
	}

	if cmpExponentZero == -1 {
		err = fmt.Errorf(ePrefix+"Error: Input parameter 'exponent' is negative! "+
			"exponent='%v' ", exponent.Text(10))

		return result, resultPrecision, err
	}

	// exponentPrecision > 0
	if exponentPrecision.Cmp(bigZero) == 1 {
		err = fmt.Errorf(ePrefix+"Error: Input parameter 'exponent' is NOT an integer! "+
			"exponentPrecision='%v' ", exponentPrecision)
		return result, resultPrecision, err
	}

	bigOne := big.NewInt(1)

	if exponent.Cmp(bigOne) == 0 {
		result.Set(base)
		resultPrecision = basePrecision
		err = nil
		return result, resultPrecision, err
	}

	result.Exp(base, exponent, nil)

	resultPrecision = big.NewInt(0).Mul(
		basePrecision,
		exponent)

	if resultPrecision.Cmp(maxPrecision) == 1 {

		bigTen := big.NewInt(10)
		delta := big.NewInt(0).Sub(resultPrecision, maxPrecision)
		delta.Sub(delta, bigOne)
		scale := big.NewInt(0).Exp(bigTen, delta, nil)
		result.Quo(result, scale)

		roundFive := big.NewInt(5)

		if result.Cmp(bigZero) == -1 {
			roundFive.Neg(roundFive)
		}
		result.Add(result, roundFive)
		result.Quo(result, bigTen)
		resultPrecision.Set(maxPrecision)

	}

	err = nil

	return result, resultPrecision, err
}

// BigIntPwrIteration - Raises input parameter 'base' to the power of input parameter
// 'exponent'.
//
// This method of raising a base to an exponent uses iterative multiplication and manages
// the internal precision of each iterative multiplication. If, during the process of
// multiplying the base time itself, the internal precision exceeds the 'internalMaxPrecision'
// limit, that intermediate number is rounded down to 'internalMaxPrecision'.
//
// If the precision of the final result exceeds the limit imposed by input parameter,
// 'outputMaxPrecision', that final result will be rounded to 'outputMaxPrecision'
// digits to the right of the decimal place.
//
// Input Parameter
// ===============
//
// base							*big.Int	-	The base which will be raised to the power of 'exponent'.
//
//                        							baseToPwr = base^exponent
//
// basePrecision				uint	- The number of digits to the right of the decimal place
//                        			in the numeric sequence represented by 'base'.
//
// exponent							uint	- This function will raise 'base' to the power of 'exponent'.
//
//                        							baseToPwr = base^exponent
//
//
// internalMaxPrecision uint	- This value is imposed as a limit on the precision of
//                              internal calculations necessary to compute the result
//                              of this power operation. If during the calculation an
//                              interim or intermediate result is generated which exceeds
//                              this limit, that intermediate result will be rounded to
//             									'internalMaxPrecision'. The term precision defines the
//                              number of digits to the right of the decimal place.
//
//                              If 'internalMaxPrecision' is less than 'outputMaxPrecision',
//                              'internalMaxPrecision' will be automatically set to a value
//                              of 'outputMaxPrecision' + 100.
//
// outputMaxPrecision		uint	- This value is imposed as a limit on the precision of
//                              the final calculated result of the power operation.
// 															If the number of digits to the right of the decimal
//                              point in the final calculated result exceeds this limit,
// 															that final result will be rounded to 'outputMaxPrecision'
//                              digits to the right of the decimal place. The term precision
// 															defines the number of digits to the right of the decimal
// 															place.
//
// Return Values
// =============
//
// baseToPwr						*big.Int	-	This function returns the result of 'base' raised
//                               		to the power of 'exponent'. This result, 'baseToPwr'
//                                  is returned as a type *big.Int.
//
//                                  				baseToPwr = base^exponent
//
// baseToPwrPrecision		uint			- Specifies the number of digits to the the right of the
//                                  decimal place in the numeric sequence represented
// 																	by the calculation result, 'baseToPwr'.
//
func (bIPwr BigIntMathPower) BigIntPwrIteration(
	base *big.Int,
	basePrecision,
	exponent,
	internalMaxPrecision,
	outputMaxPrecision uint) (baseToPwr *big.Int, baseToPwrPrecision uint) {

	baseToPwrPrecision = 0

	exponent++

	bigIZero := big.NewInt(0)

	if base.Cmp(bigIZero) == 0 {
		baseToPwr = big.NewInt(0)
		baseToPwrPrecision = 0
		return baseToPwr, baseToPwrPrecision
	}

	if base.Cmp(big.NewInt(1)) == 0 {
		baseToPwr = big.NewInt(0).Set(base)
		baseToPwrPrecision = basePrecision
		return baseToPwr, baseToPwrPrecision
	}

	if internalMaxPrecision < outputMaxPrecision {
		internalMaxPrecision = outputMaxPrecision + 100
	}

	bigIPlusFive := big.NewInt(5)
	bigIMinusFive := big.NewInt(-5)
	bigITen := big.NewInt(10)
	roundFactor := big.NewInt(0)
	cmpResult := 0

	for i := uint(0); i < exponent; i++ {

		if i == 0 {
			baseToPwr = big.NewInt(1)
			baseToPwrPrecision = 0
		} else if i == 1 {
			baseToPwr = big.NewInt(0).Set(base)
			baseToPwrPrecision = basePrecision
		} else {
			baseToPwr = big.NewInt(0).Mul(baseToPwr, base)
			baseToPwrPrecision = baseToPwrPrecision + basePrecision
		}

		if baseToPwrPrecision > internalMaxPrecision {
			// beforeRounding := BigIntNum{}.NewBigInt(baseToPwr, baseToPwrPrecision)
			delta := baseToPwrPrecision - internalMaxPrecision - 1
			scale := big.NewInt(0).Exp(bigITen, big.NewInt(int64(delta)), nil)

			cmpResult = baseToPwr.Cmp(bigIZero)
			if cmpResult == 1 {
				// baseToPwr is GREATER Than zero
				roundFactor = big.NewInt(0).Mul(bigIPlusFive, scale)
			} else if cmpResult == -1 {
				// baseToPwr is LESS Than zero
				roundFactor = big.NewInt(0).Mul(bigIMinusFive, scale)
			} else {
				baseToPwrPrecision = 0
				continue
			}

			baseToPwr = big.NewInt(0).Add(baseToPwr, roundFactor)
			scale = big.NewInt(0).Mul(scale, bigITen)
			baseToPwr = big.NewInt(0).Quo(baseToPwr, scale)
			baseToPwrPrecision = internalMaxPrecision

		}

	}

	if baseToPwrPrecision > outputMaxPrecision {
		delta := baseToPwrPrecision - outputMaxPrecision - 1
		scale := big.NewInt(0).Exp(bigITen, big.NewInt(int64(delta)), nil)
		round := big.NewInt(0).Mul(bigIPlusFive, scale)
		if baseToPwr.Cmp(bigIZero) == -1 {
			round = big.NewInt(0).Mul(round, big.NewInt(-1))
		}
		baseToPwr = big.NewInt(0).Add(baseToPwr, round)
		scale = big.NewInt(0).Mul(scale, bigITen)
		baseToPwr = big.NewInt(0).Quo(baseToPwr, scale)
		baseToPwrPrecision = outputMaxPrecision

	}

	return baseToPwr, baseToPwrPrecision
}

// BigIntegerPwrIteration - Raises input parameter 'base' to the power of input parameter
// 'exponent'. This version of the function uses *big.Int type for all input parameters
//
// This method of raising a base to an exponent uses iterative multiplication and manages
// the internal precision of each iterative multiplication. If, during the process of
// multiplying the base time itself, the internal precision exceeds the 'internalMaxPrecision'
// limit, that intermediate number is rounded down to 'internalMaxPrecision'.
//
// If the precision of the final result exceeds the limit imposed by input parameter,
// 'outputMaxPrecision', that final result will be rounded to 'outputMaxPrecision'
// digits to the right of the decimal place.
//
// Input Parameter
// ===============
//
// base							*big.Int	-	The base which will be raised to the power of 'exponent'.
//
//                        							baseToPwr = base^exponent
//
// basePrecision		*big.Int	- The number of digits to the right of the decimal place
//                        			in the numeric sequence represented by 'base'.
//
// exponent					*big.Int	- This function will raise 'base' to the power of 'exponent'.
//
//                        							baseToPwr = base^exponent
//
//                              'exponent' is an integer value with zero precision.
//
//
// internalMaxPrecision *big.Int	- This value is imposed as a limit on the precision of
//                              		internal calculations necessary to compute the result
//                              		of this power operation. If during the calculation an
//                              		interim or intermediate result is generated which exceeds
//                              		this limit, that intermediate result will be rounded to
//             											'internalMaxPrecision'. The term precision defines the
//                              		number of digits to the right of the decimal place.
//
//                              		If 'internalMaxPrecision' is less than 'outputMaxPrecision',
//                              		'internalMaxPrecision' will be automatically set to a value
//                              		of 'outputMaxPrecision' + 100.
//
// outputMaxPrecision		*big.Int	- This value is imposed as a limit on the precision of
//                              		the final calculated result of the power operation.
// 																	If the number of digits to the right of the decimal
//                              		point in the final calculated result exceeds this limit,
// 																	that final result will be rounded to 'outputMaxPrecision'
//                              		digits to the right of the decimal place. The term precision
// 																	defines the number of digits to the right of the decimal
// 																	place. If 'outputMaxPrecision' is less than zero, an
//                                  error will be returned.
//
// Return Values
// =============
//
// baseToPwr						*big.Int	-	This function returns the result of 'base' raised
//                               		to the power of 'exponent'. This result, 'baseToPwr'
//                                  is returned as a type *big.Int.
//
//                                  				baseToPwr = base^exponent
//
// baseToPwrPrecision		*big.Int	- Specifies the number of digits to the the right of the
//                                  decimal place in the numeric sequence represented
// 																	by the calculation result, 'baseToPwr'.
//
// err									error			- If the function fails to complete successfully, this value
//                                  is configured with an appropriate error message and returned
//                                  to the caller. If the function completes successfully, this
//                                  value is set to 'nil'.
//
func (bIPwr BigIntMathPower) BigIntegerPwrIteration(
	base,
	basePrecision,
	exponent,
	internalMaxPrecision,
	outputMaxPrecision *big.Int) (baseToPwr *big.Int, baseToPwrPrecision *big.Int, err error) {

	ePrefix := "BigIntMathPower.BigIntegerPwrIteration() "
	baseToPwr = big.NewInt(0)
	baseToPwrPrecision = big.NewInt(0)
	err = nil

	bigZero := big.NewInt(0)
	bigOne := big.NewInt(1)

	if basePrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+"Error: Input Parameter 'basePrecision' is a negative value! "+
			"basePrecision='%v' ", basePrecision.Text(10))
		return baseToPwr, baseToPwrPrecision, err
	}

	if outputMaxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+"Error: Input Parameter 'outputMaxPrecision' is a negative value! "+
			"outputMaxPrecision='%v' ", outputMaxPrecision.Text(10))
		return baseToPwr, baseToPwrPrecision, err
	}

	if base.Cmp(bigZero) == 0 {
		return baseToPwr, baseToPwrPrecision, err
	}

	if base.Cmp(bigOne) == 0 &&
		baseToPwrPrecision.Cmp(bigZero) == 0 {
		baseToPwr = big.NewInt(0).Set(bigOne)
		baseToPwrPrecision.Set(basePrecision)
		return baseToPwr, baseToPwrPrecision, err
	}

	if base.Cmp(big.NewInt(-1)) == 0 &&
		baseToPwrPrecision.Cmp(bigZero) == 0 {
		baseToPwr = big.NewInt(-1)
		baseToPwrPrecision.Set(basePrecision)
		return baseToPwr, baseToPwrPrecision, err
	}

	if internalMaxPrecision.Cmp(outputMaxPrecision) != 1 {
		internalMaxPrecision = big.NewInt(0).Add(outputMaxPrecision, big.NewInt(100))
	}

	cycles := big.NewInt(0).Add(exponent, bigOne)
	bigIPlusFive := big.NewInt(5)
	bigIMinusFive := big.NewInt(-5)
	bigITen := big.NewInt(10)
	roundFactor := big.NewInt(0)
	cmpResult := 0

	for i := big.NewInt(0); i.Cmp(cycles) == -1; i.Add(i, bigOne) {

		if i.Cmp(bigZero) == 0 {
			baseToPwr = big.NewInt(1)
			baseToPwrPrecision = big.NewInt(0)
		} else if i.Cmp(bigOne) == 0 {
			baseToPwr = big.NewInt(0).Set(base)
			baseToPwrPrecision = big.NewInt(0).Set(basePrecision)
		} else {
			baseToPwr = big.NewInt(0).Mul(baseToPwr, base)
			baseToPwrPrecision.Add(baseToPwrPrecision, basePrecision)
		}

		if baseToPwrPrecision.Cmp(internalMaxPrecision) == 1 {
			// beforeRounding := BigIntNum{}.NewBigInt(baseToPwr, baseToPwrPrecision)
			delta := big.NewInt(0).Sub(baseToPwrPrecision, internalMaxPrecision)
			delta.Sub(delta, bigOne)
			scale := big.NewInt(0).Exp(bigITen, delta, nil)

			cmpResult = baseToPwr.Cmp(bigZero)
			if cmpResult == 1 {
				// baseToPwr is GREATER Than zero
				roundFactor = big.NewInt(0).Mul(bigIPlusFive, scale)
			} else if cmpResult == -1 {
				// baseToPwr is LESS Than zero
				roundFactor = big.NewInt(0).Mul(bigIMinusFive, scale)
			} else {
				baseToPwrPrecision = big.NewInt(0)
				continue
			}

			baseToPwr = big.NewInt(0).Add(baseToPwr, roundFactor)
			scale = big.NewInt(0).Mul(scale, bigITen)
			baseToPwr = big.NewInt(0).Quo(baseToPwr, scale)
			baseToPwrPrecision = big.NewInt(0).Set(internalMaxPrecision)

		}

	}

	if baseToPwrPrecision.Cmp(outputMaxPrecision) == 1 {
		delta := big.NewInt(0).Sub(baseToPwrPrecision, outputMaxPrecision)
		delta.Sub(delta, bigOne)

		scale := big.NewInt(0).Exp(bigITen, delta, nil)
		roundFactor = big.NewInt(0).Mul(bigIPlusFive, scale)
		if baseToPwr.Cmp(bigZero) == -1 {
			roundFactor = big.NewInt(0).Mul(roundFactor, big.NewInt(-1))
		}
		baseToPwr = big.NewInt(0).Add(baseToPwr, roundFactor)
		scale = big.NewInt(0).Mul(scale, bigITen)
		baseToPwr = big.NewInt(0).Quo(baseToPwr, scale)
		baseToPwrPrecision = big.NewInt(0).Set(outputMaxPrecision)

	}

	return baseToPwr, baseToPwrPrecision, err
}

// FixedDecimalPwrIteration - Raises input parameter 'base' to the power of input parameter 'exponent'.
// This method of raising a base to an exponent uses iterative multiplication and manages
// the internal precision of each iterative multiplication. If, during the process of
// multiplying the base times itself, the internal precision exceeds the 'internalMaxPrecision'
// limit, that intermediate number is rounded down to 'internalMaxPrecision' digits to the
// right of the decimal place.
//
// If the precision of the final result exceeds the limit imposed by input parameter,
// 'outputMaxPrecision', that final result will be rounded to 'outputMaxPrecision'
// digits to the right of the decimal place.
//
// Input Parameter
// ===============
//
// base	BigIntFixedDecimal	-	The base which will be raised to the power of 'exponent'.
//                            The BigIntFixedDecimal type describes a numeric value with
//                            a fixed number of digits to the right of the decimal place.
//                            The type includes a *big.Int integer value and a precision
//                            specification.
//
//                        							baseToPwr = base^exponent
//
//
// exponent							uint	- This function will raise 'base' to the power of 'exponent'.
//
//                        							baseToPwr = base^exponent
//
//
// internalMaxPrecision uint	- This value is imposed as a limit on the precision of
//                              internal calculations necessary to compute the result
//                              of this power operation. If during the calculation an
//                              interim or intermediate result is generated which exceeds
//                              this limit, that intermediate result will be rounded to
//             									'internalMaxPrecision' digits to the right of the decimal
//                              place. The term precision defines the number of digits to
// 															the right of the decimal place.
//
//                              If 'internalMaxPrecision' is less than 'outputMaxPrecision',
//                              'internalMaxPrecision' will be automatically set to a value
//                              of 'outputMaxPrecision' + 100.
//
// outputMaxPrecision		uint	- This value is imposed as a limit on the precision of
//                              the final calculated result of the power operation.
// 															If the number of digits to the right of the decimal
//                              point in the final calculated result exceeds this limit,
// 															that final result will be rounded to 'outputMaxPrecision'
//                              digits to the right of the decimal place. The term precision
// 															defines the number of digits to the right of the decimal
// 															place.
//
// Return Values
// =============
//
// baseToPwr	BigIntFixedDecimal	-	This function returns the result of 'base' raised
//                               		to the power of 'exponent'. This result, 'baseToPwr'
//                                  is returned as a type BigIntFixedDecimal. The BigIntFixedDecimal
// 																	type describes a numeric value with a fixed number of digits
// 																	to the right of the decimal place.
//
//                                  				baseToPwr = base^exponent
//
func (bIPwr BigIntMathPower) FixedDecimalPwrIteration(
	base BigIntFixedDecimal,
	exponent,
	internalMaxPrecision,
	outputMaxPrecision uint) (baseToPwr BigIntFixedDecimal) {

	baseToPwr = BigIntFixedDecimal{}.NewZero(0)

	bigIPwr, bigIPwrPrecision, _ := BigIntMathPower{}.BigIntegerPwrIteration(
		base.GetInteger(),
		base.GetPrecisionBigInt(),
		big.NewInt(0).SetUint64(uint64(exponent)),
		big.NewInt(0).SetUint64(uint64(internalMaxPrecision)),
		big.NewInt(0).SetUint64(uint64(outputMaxPrecision)))

	baseToPwr, _ = BigIntFixedDecimal{}.NewBigIntPrecision(bigIPwr, bigIPwrPrecision)

	return baseToPwr
}

// MinimumRequiredPrecision - designed to be used with the power function
// below. This method will compute the minimum number of decimal places
// required to support the result of raising a 'base' value to a specified
// exponent. Both the 'base' and the 'exponent' are passed to this function
// as type BigIntNum.
//
// For example, raising the value 3.12 to the power of 4 means that the
// result will require at least 8-decimal places to the right of the
// decimal in order to display a correct result. In the following example
// with base ='3.12' and exponent = '4', this method will return '8'.
//
// 		Example: 3.12^4 = 94.75854336 (8-digits to the right of the decimal)
//
// The calculated minimum required precision is returned as type 'uint'.
//
// If the minimum required precision exceeds the maximum value for type
// 'uint' (+4,294,967,295, which equals 2^32 − 1), an error message is returned
// in addition to the maximum uint value (+4,294,967,295).
//
func (bIPwr BigIntMathPower) MinimumRequiredPrecision(
	base, exponent BigIntNum) (uint, error) {

	basePrecision := BigIntNum{}.NewUint(base.GetPrecisionUint(), 0)

	tExponent := exponent.CopyOut()

	if tExponent.GetSign() == -1 {
		tExponent.ChangeSign()
	}

	minRequiredPrecision :=
		BigIntMathMultiply{}.MultiplyBigIntNums(basePrecision, tExponent)

	minRequiredPrecision.RoundToDecPlace(0)

	return minRequiredPrecision.GetUInt()
}

// Pwr - Raises 'base' to the power of 'exponent'.  Both 'base' and 'exponent' are Type BigIntNum.
// Upon computing the result of 'base' raised to the power of 'exponent' (base^exponent), the result
// is returned as Type BigIntNum.
//
// Examples:
// =========
//
//	base				exponent				maxPrecision	  		result
// -------			---------  			-------------	 			-------
//
//	2							 4								17								16
//  2							-4								17								 0.0625
//  3.7						 2.8							30								38.991040735983142451443031376258
//  3.7						-2.8							32 								 0.02564691737189623971146450249457
// -2							-3.8							32								 0.07179364718731468792491418417362
// -2              3.8              30								13.928809012737986226180320279676
//
// The return value, a type BigIntNum, represents the result of the base^exponent operation described above.
// This returned BigIntNum 'result' will contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter,'base'.
//
func (bIPwr BigIntMathPower) Pwr(base, exponent BigIntNum, maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathPower.Pwr() "

	result := BigIntNum{}.NewWithNumSeps(base.GetNumericSeparatorsDto())
	var err error

	if base.IsZero() {
		// If 'base' is zero, return zero value.
		return result, nil
	}

	if exponent.IsZero() {
		return BigIntNum{}.NewOne(0), nil
	}

	numSeps := base.GetNumericSeparatorsDto()

	bigOne := BigIntNum{}.NewOne(exponent.GetPrecisionUint())

	if exponent.Equal(bigOne) {
		result = base.CopyOut()
		return base, nil
	}

	if exponent.GetPrecisionUint() == uint(0) {

		if exponent.GetSign() > 0 {
			// exponent is a positive number
			result, err = bIPwr.bigIntNumRaiseToPositiveIntegerPower(base, exponent)

			if err != nil {
				return BigIntNum{}.NewZero(0),
					fmt.Errorf(ePrefix+
						"Error returned by bIPwr.bigIntNumRaiseToPositiveIntegerPower("+
						"base, exponent, maxPrecision) "+
						"base='%v' exponent='%v' Error='%v' ",
						base.GetNumStr(), exponent.GetNumStr(), err.Error())
			}

		} else {
			// exponent must be a negative number
			result, err = bIPwr.bigIntNumRaiseToNegativeIntegerPower(base, exponent, maxPrecision)

			if err != nil {
				return BigIntNum{}.NewZero(0),
					fmt.Errorf(ePrefix+"Error returned by bIPwr.bigIntNumRaiseToNegativeIntegerPower("+
						"base, exponent, maxPrecision) "+
						"base='%v' exponent='%v' maxPrecision='%v' Error='%v' ",
						base.GetNumStr(), exponent.GetNumStr(), maxPrecision, err.Error())
			}
		}

	} else {
		// precision must be greater than zero. exponent is a fractional number

		if exponent.GetSign() > 0 {
			// fractional exponent is a positive number

			result, err = bIPwr.bigIntNumRaiseToPositiveFractionalPower(base, exponent, maxPrecision)

			if err != nil {
				return BigIntNum{},
					fmt.Errorf(ePrefix+
						" Error='%v' \n", err.Error())
			}

		} else {
			// fractional exponent must be a negative number

			result, err = bIPwr.bigIntNumRaiseToNegativeFractionalPower(base, exponent, maxPrecision)

			if err != nil {
				return BigIntNum{},
					fmt.Errorf(ePrefix+
						"Error='%v' \n", err.Error())
			}

		}
	}

	if result.precision > maxPrecision {
		result.SetPrecision(maxPrecision)
	}

	err = result.SetNumericSeparatorsDto(numSeps)

	return result, nil
}

// bigIntNumRaiseToNegativeFractionalPower - Assumes that input parameter 'exponent' is negative and a
// fractional number (precision > 0 - has fractional digits). If 'exponent' is positive or if
// 'exponent' is an integer number, an error is returned.
//
// If 'exponent' is both a negative number and a fractional number, this method proceeds to
// raise input parameter 'base' to the power of 'exponent' and return the result as a BigIntNum
// type.
//
func (bIPwr BigIntMathPower) bigIntNumRaiseToNegativeFractionalPower(
	base,
	exponent BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathPower.bigIntNumRaiseToNegativeFractionalPower() "

	if exponent.GetPrecisionUint() == 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error - The Exponent is an Integer Number. "+
				"Only Fractional Exponents can be processed by this method! exponent='%v' ",
				exponent.GetNumStr())
	}

	if exponent.GetSign() > 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error = The Exponent is a Positive Number. "+
				"Only Negative Exponents can be processed by this method! exponent='%v'",
				exponent.GetNumStr())
	}

	exponentAbsVal := exponent.GetAbsoluteBigIntNumValue()

	bINumResult, err :=
		bIPwr.bigIntNumRaiseToPositiveFractionalPower(
			base,
			exponentAbsVal,
			maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error='%v' ", err.Error())
	}

	inverseResult, err := bINumResult.GetInverse(maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by bINumResult.GetInverse(maxPrecision). "+
				"bINumResult='%v' maxPrecision='%v' Error='%v' ",
				bINumResult.GetNumStr(), maxPrecision, err.Error())
	}

	return inverseResult, nil
}

// bigIntNumRaiseToPositiveFractionalPower - Assumes that input parameter 'exponent' is positive and a
// fractional number (precision > 0 - has fractional digits). If 'exponent' is negative or if
// 'exponent' is an integer number, an error is returned.
//
// If 'exponent' is both a positive number and a fractional number, this method proceeds to
// raise input parameter 'base' to the power of 'exponent' and return the result as a BigIntNum
// type.
//
func (bIPwr BigIntMathPower) bigIntNumRaiseToPositiveFractionalPower(
	base,
	exponent BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathPower.bigIntNumRaiseToPositiveFractionalPower() "

	if exponent.GetPrecisionUint() == 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error - The Exponent is an Integer Number. "+
				"Only Fractional Exponents can be processed by this method! exponent='%v' ",
				exponent.GetNumStr())
	}

	if exponent.GetSign() < 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error = The Exponent is a Negative Number. "+
				"Only Positive Exponents can be processed by this method! exponent='%v'",
				exponent.GetNumStr())
	}

	modXZero := big.NewInt(0)

	scaleFactor := big.NewInt(0).Exp(
		big.NewInt(10),
		big.NewInt(int64(exponent.precision)),
		modXZero)

	ratExponent := big.NewRat(1, 1).SetFrac(exponent.bigInt, scaleFactor)

	bINumNumerator := BigIntNum{}.NewBigInt(ratExponent.Num(), 0)

	pwr1Result, err := bIPwr.bigIntNumRaiseToPositiveIntegerPower(base, bINumNumerator)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by bIPwr.raiseToPositive"+
				"IntegerPower(base, bINumNumerator)"+
				"base='%v' bINumNumerator='%v' Error='%v' \n",
				base.GetNumStr(), bINumNumerator.GetNumStr(), err.Error())
	}

	nthRoot := BigIntNum{}.NewBigInt(ratExponent.Denom(), 0)

	// If nthRoot is zero, the result will always be '1'
	if nthRoot.IsZero() {
		return BigIntNum{}.NewOne(maxPrecision), nil
	}

	bigINumOne := BigIntNum{}.NewOne(0)

	// Error if nthRoot == 1
	if nthRoot.Cmp(bigINumOne) == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix +
				"Error - Input Parameter 'nthRoot' INVALID! 'nthRoot' cannot equal 1.\n")
	}

	if pwr1Result.GetSign() == -1 {

		isEvenNum, err := nthRoot.IsEvenNumber()

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by nthRoot.IsEvenNumber() "+
					"nthRoot='%v' Error='%v'\n", nthRoot.GetNumStr(), err.Error())
		}

		if isEvenNum {
			return BigIntNum{}.NewZero(0),
				errors.New(ePrefix +
					"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even. \n")
		}

	}

	biNumResult, err := BigIntMathNthRoot{}.GetNthRoot(pwr1Result, nthRoot, maxPrecision)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by dBigIntMathNthRoot{}.OriginalNthRoot(...) "+
				"Error='%v'", err.Error())
	}

	return biNumResult, nil
}

// bigIntNumRaiseToNegativeIntegerPower - Assumes that input parameter 'exponent' is negative and an
// integer number (precision==0 - no fractional digits). If 'exponent' is positive or if
// it is NOT an integer number, an error is returned.
//
// If 'exponent' is both a negative number and an integer number, this method proceeds to
// raise input parameter 'base' to the power of 'exponent' and return the result as a BigIntNum
// type.
//
func (bIPwr BigIntMathPower) bigIntNumRaiseToNegativeIntegerPower(
	base,
	exponent BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathPower.bigIntNumRaiseToNegativeIntegerPower()"

	if exponent.GetPrecisionUint() > 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error - The Exponent is a Fractional Number. "+
				"Only Integer Exponents can be processed by this method! exponent='%v' ",
				exponent.GetNumStr())
	}

	if exponent.GetSign() > 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error = The Exponent is a Positive Number. "+
				"Only Negative Exponents can be processed by this method! exponent='%v'",
				exponent.GetNumStr())
	}

	bigIBasePrecision := big.NewInt(int64(base.GetPrecisionUint()))

	bigINewPrecision := big.NewInt(0).Mul(bigIBasePrecision, exponent.absBigInt)

	newPrecision := uint(bigINewPrecision.Int64())

	result := big.NewInt(0).Exp(base.bigInt, exponent.absBigInt, nil)

	bINumResult1 := BigIntNum{}.NewBigInt(result, newPrecision)

	inverseResult, err := bINumResult1.GetInverse(maxPrecision)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by bINumResult1.GetInverse(maxPrecision)"+
				"bINumResult1='%v' maxPrecision='%v' Error='%v'",
				bINumResult1.GetNumStr(), maxPrecision, err.Error())
	}

	return inverseResult, nil
}

// bigIntNumRaiseToPositiveIntegerPower - Assumes that input parameter 'exponent' is a positive
// integer number. If 'exponent' is negative or if 'exponent' is NOT an integer (precision > 0),
// an error will be triggered.
//
// If the exponent is both positive and an integer number, this method proceeds to raise the 'base'
// parameter to the power of 'exponent' and returns the result as a 'BigIntNum' type.
//
//	Examples:
//  =========
//
//	base^exponent = result
//
//	Base			Exponent		Result
//	 2					2						 4
//   3					4						81
//   4.2				3						74.088
// 	10					3					1000
//  -4.2				3					 -74.088
//	-2.9				4					  70.7281
//  -2          3.8
func (bIPwr BigIntMathPower) bigIntNumRaiseToPositiveIntegerPower(
	base,
	exponent BigIntNum) (BigIntNum, error) {

	ePrefix := "BigIntMathPower.bigIntNumRaiseToPositiveIntegerPower() "

	if exponent.GetPrecisionUint() > 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error - The Exponent is a Fractional Number. "+
				"Only Integer Exponents can be processed by this method! exponent='%v' ",
				exponent.GetNumStr())
	}

	if exponent.GetSign() < 0 {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error = The Exponent is a Negative Number. "+
				"Only Positive Exponents can be processed by this method! exponent='%v'",
				exponent.GetNumStr())
	}

	bigIBasePrecision := big.NewInt(int64(base.GetPrecisionUint()))

	bigINewPrecision := big.NewInt(0).Mul(bigIBasePrecision, exponent.bigInt)

	newPrecision := uint(bigINewPrecision.Int64())

	result := big.NewInt(0).Exp(base.bigInt, exponent.bigInt, nil)

	return BigIntNum{}.NewBigInt(result, newPrecision), nil
}

// computeInternalPrecision - Returns computed internal precision for variables used
// in intermediate calculations. Returned 'internalPrecision' is based on requested
// maximum precision for a specific BigIntMathPower calculation.
/*
func (bIPwr BigIntMathPower) computeMaxInternalPrecision(maxPrecision *big.Int) *big.Int {

	internalPrecision := big.NewInt(10)

	maxPrecisionCmp25 := maxPrecision.Cmp(big.NewInt(50))

	maxPrecisionCmp50 := maxPrecision.Cmp(big.NewInt(50))

	maxPrecisionCmp200 := maxPrecision.Cmp(big.NewInt(200))

	maxPrecisionCmpTwoThou := maxPrecision.Cmp(big.NewInt(2000))

	maxPrecisionCmpFiveThou := maxPrecision.Cmp(big.NewInt(5000))

	maxPrecisionCmpTenThou := maxPrecision.Cmp(big.NewInt(10000))

	if maxPrecisionCmp25 < 1 {

		internalPrecision.Add(maxPrecision, big.NewInt(50))

	} else if  maxPrecisionCmp25 >= 0 &&   maxPrecisionCmp50 == -1 {

		internalPrecision.Add(maxPrecision, big.NewInt(65))

	} else if maxPrecisionCmp50 >= 0  && maxPrecisionCmp200 == -1 {
		internalPrecision.Add(maxPrecision, big.NewInt(75))

	} else if maxPrecisionCmp200 >= 0 && maxPrecisionCmpTwoThou ==-1 {
		internalPrecision.Add(maxPrecision, big.NewInt(0).Quo(maxPrecision, big.NewInt(4)))

	}else if maxPrecisionCmpTwoThou >= 0 && maxPrecisionCmpFiveThou == -1 {
		internalPrecision.Add(maxPrecision, big.NewInt(0).Quo(maxPrecision, big.NewInt(8)))

	}else if maxPrecisionCmpFiveThou >= 0 && maxPrecisionCmpTenThou == -1 {
		internalPrecision.Add(maxPrecision, big.NewInt(0).Quo(maxPrecision, big.NewInt(100)))

	} else if maxPrecisionCmpTenThou >= 0  {
		internalPrecision.Add(maxPrecision, big.NewInt(100))
	}

	return internalPrecision

}
*/
