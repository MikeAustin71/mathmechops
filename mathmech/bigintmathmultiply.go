package mathops

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

// BigIntMathMultiply - This methods associated with this type are
// used to perform math operations using the *big.Int Type.
//
// If you are unfamiliar with the *big.Int type, reference:
// 						https://golang.org/pkg/math/big/
//
type BigIntMathMultiply struct {
	Input  BigIntPair
	Result BigIntNum
}

// MultiplyBigInts - Receives two *big.Int numbers and their associated precision
// specifications. This method then proceeds to perform a multiplication operation
// by multiplying the 'multiplier' by the 'multiplicand' to generate the 'product'.
// 'multiplier', 'multiplicand' and 'product' are configured as pairs of *big.Int
// integer numbers and precision specifications. Taken together, an integer number
// and precision specification are used to defined a fixed length floating point
// number.
//
// Examples
// ========
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// Consider the following multiplication example.
//
//							752.314 x 21.67894 = product
//
// 'multiplier' and 'multiplicand' would be configured as follows:
//									multiplier 						= 752314
//                  multiplierPrecision		= 3
//                  multiplicand 					= 2167894
//                  multiplicandPrecision = 5
//
// The 'product' value (16309.37006716) of 'multiplier' and 'multiplicand' would be calculated
// and configured as follows:
//
// 									product						= 1630937006716
//                  productPrecision	= 8
//
//
// Input Parameters
// ================
//
//	multiplier 						*big.Int	- The number to be multiplied by 'multiplicand'
//	multiplierPrecision 	*big.Int	- The 'multiplier' precision or numeric digits after
//																				the decimal point.
//	multiplicand 					*big.Int	- The number to be multiplied by the 'multiplier'.
//	multiplicandPrecision *big.Int	- The 'multiplicand' precision or numeric digits after
//																		the decimal point.
//
// Return Values
// =============
//
// product								*big.Int	- The product of the multiplier multiplied by
//                                		the multiplicand.
//
// productPrecision				*big.Int	- The precision specification for the returned
//                                		product. Here, the term precision is defined
// 																		as the number of fractional digits to the
//                                		right of the decimal place in the returned
//                                		'product'. 'productPrecision' is always equal
//																		to or greater than zero.
//
// Note: This method removes trailing fractional zeros from the result.
//			 Example: 3.1200 is returned as 3.12
//
func (bMultiply BigIntMathMultiply) BigIntMultiply(
	multiplier,
	multiplierPrecision,
	multiplicand,
	multiplicandPrecision *big.Int) (product *big.Int, productPrecision *big.Int, err error) {

	ePrefix := "BigIntMathMultiply.BigIntMultiply() "

	product = big.NewInt(0)
	productPrecision = big.NewInt(0)
	err = nil

	if multiplier == nil {
		multiplier = big.NewInt(0)
	}

	if multiplicand == nil {
		multiplicand = big.NewInt(0)
	}

	if multiplierPrecision == nil {
		multiplierPrecision = big.NewInt(0)
	}

	if multiplicandPrecision == nil {
		multiplicandPrecision = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if multiplierPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter multiplierPrecision is LESS THAN ZERO! "+
			"multiplierPrecision='%v' ", multiplierPrecision.Text(10))

		return product, productPrecision, err
	}

	if multiplicandPrecision.Cmp(bigZero) == -1 {

		err = fmt.Errorf(ePrefix+
			"Error: Input parameter multiplicandPrecision is LESS THAN ZERO! "+
			"multiplicandPrecision='%v' ", multiplicandPrecision.Text(10))

		return product, productPrecision, err
	}

	productPrecision.Add(multiplierPrecision, multiplicandPrecision)

	product = big.NewInt(0).Mul(multiplier, multiplicand)

	if product.Cmp(bigZero) == 0 {
		productPrecision = big.NewInt(0)
	}

	// Delete trailing fractional zeros
	// If productPrecision > 0

	if productPrecision.Cmp(bigZero) == 1 {
		bigOne := big.NewInt(1)
		scrap := big.NewInt(0)
		biBase10 := big.NewInt(10)
		biBaseZero := big.NewInt(0)
		newProduct, mod10 := big.NewInt(0).QuoRem(product, biBase10, scrap)

		for mod10.Cmp(biBaseZero) == 0 && productPrecision.Cmp(bigZero) == 1 {
			product.Set(newProduct)
			productPrecision.Sub(productPrecision, bigOne)
			newProduct, mod10 = big.NewInt(0).QuoRem(product, biBase10, scrap)
		}
	}

	err = nil

	return product, productPrecision, err
}

// BigIntMultiplyByTenToPower - Multiplies a *big.Int number by ten
// to the power of 'exponent'. 'exponent' is an input parameter of
// type *big.Int.
//
// 				product = multiplier x 10^exponent
//
// The result is returned as a *big.Int type ('product') with a precision
// specification ('productPrecision').
//
// Input Parameters
// ================
//
// multiplier						*big.Int - 'multiplier' will be multiplied by 10 to power of
//                       						'exponent' to generate the result or 'product'.
//
// multiplierPrecision	*big.Int	- The precision specification for 'multiplier'.
// 																	Precision specifies the number of digits to the
//                                  right of the decimal place in 'multiplier'. This
//                                  value must be greater than or equal to zero.
//
// exponent							*big.Int	- Ten will be raised to the power of exponent and
//                                  multiplied by 'multiplier' to generate the result
//                                  or product. 'exponent' can be a negative value.
//
// Return Values
// =============
//
// product							*big.Int	- The result generated by multiplying 'multiplier' by
//                                  ten to the power of 'exponent'.
//
// productPrecision			*big.Int	- The precision specification for 'product'. Precision
//                                  specifies the number of digits to the right of the
//                                  decimal place in 'product'. This value will always
//                                  be greater than or equal to zero.
//
// err									error			- If 'multiplierPrecision' or 'exponent' are less than
//                                  zero, an error will be returned.
//
// Note: 	This method will remove trailing fractional zeros from the
// 				final result (product).
//
func (bMultiply BigIntMathMultiply) BigIntMultiplyByTenToPower(
	multiplier,
	multiplierPrecision,
	exponent *big.Int) (product *big.Int, productPrecision *big.Int, err error) {

	ePrefix := "BigIntMathMultiply.BigIntMultiplyByTenToPower() "

	product = big.NewInt(0)
	productPrecision = big.NewInt(0)
	err = nil

	if multiplier == nil {
		multiplier = big.NewInt(0)
	}

	if multiplierPrecision == nil {
		multiplierPrecision = big.NewInt(0)
	}

	if exponent == nil {
		exponent = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if multiplierPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'multiplierPrecision' is LESS THAN ZERO! "+
			"multiplierPrecision='%v' ", multiplierPrecision.Text(10))

		return product, productPrecision, err
	}

	bigTen := big.NewInt(10)
	bigOne := big.NewInt(1)

	if exponent.Cmp(bigZero) == -1 {
		// if exponent is less than zero
		// apply special processing...

		exponent.Neg(exponent)

		divisor := big.NewInt(0).Exp(bigTen, exponent, nil)

		quoFrac, quoFracPrecision, errx :=
			BigIntMathDivide{}.BigIntFracQuotient(
				bigOne,
				big.NewInt(0),
				divisor,
				big.NewInt(0),
				exponent)

		if errx != nil {
			err = fmt.Errorf(ePrefix+"%v", errx.Error())
			return product, productPrecision, err
		}

		product, productPrecision, errx =
			BigIntMathMultiply{}.BigIntMultiply(
				multiplier,
				multiplierPrecision,
				quoFrac,
				quoFracPrecision)

		if errx != nil {
			err = fmt.Errorf(ePrefix+"%v", errx.Error())
			product = big.NewInt(0)
			productPrecision = big.NewInt(0)
			return product, productPrecision, err
		}

		// Delete trailing fractional zeros
		// if productPrecision > 0
		if productPrecision.Cmp(bigZero) == 1 {
			scrap := big.NewInt(0)
			newProduct, mod10 := big.NewInt(0).QuoRem(product, bigTen, scrap)

			for mod10.Cmp(bigZero) == 0 && productPrecision.Cmp(bigZero) == 1 {
				product.Set(newProduct)
				productPrecision.Sub(productPrecision, bigOne)
				newProduct, mod10 = big.NewInt(0).QuoRem(product, bigTen, scrap)
			}
		}

		err = nil

		return product, productPrecision, err
	}

	delta := big.NewInt(0)
	scale := big.NewInt(0)

	if multiplier.Cmp(bigZero) == 0 {
		product = big.NewInt(0)
		productPrecision = big.NewInt(0)

	} else if exponent.Cmp(bigZero) == 0 {
		product.Set(multiplier)
		productPrecision.Set(multiplierPrecision)

	} else if exponent.Cmp(multiplierPrecision) == 0 {
		product.Set(multiplier)
		productPrecision = big.NewInt(0)

	} else if exponent.Cmp(multiplierPrecision) == 1 {
		// exponent > multiplierPrecision
		delta = big.NewInt(0).Sub(exponent, multiplierPrecision)
		scale = big.NewInt(0).Exp(bigTen, delta, nil)
		product = big.NewInt(0).Mul(multiplier, scale)
		productPrecision = big.NewInt(0)

	} else {
		// multiplierPrecision must be GREATER THAN exponent
		product.Set(multiplier)
		productPrecision.Sub(multiplierPrecision, exponent)
	}

	// Delete trailing fractional zeros
	// if productPrecision > 0
	if productPrecision.Cmp(bigZero) == 1 {
		scrap := big.NewInt(0)
		newProduct, mod10 := big.NewInt(0).QuoRem(product, bigTen, scrap)

		for mod10.Cmp(bigZero) == 0 && productPrecision.Cmp(bigZero) == 1 {
			product.Set(newProduct)
			productPrecision.Sub(productPrecision, bigOne)
			newProduct, mod10 = big.NewInt(0).QuoRem(product, bigTen, scrap)
		}
	}

	err = nil
	return product, productPrecision, err
}

// MultiplyBigIntByTwoToPower - Multiplies a *big.Int number by powers
// of two and returns the result as a *big.Int value and associated
// precision specification.
//
//
// 						product = multiplier X 2^exponent
//            productPrecision = multiplierPrecision
//
// Examples:
// =========
//	multiplier			multiplierPrecision 	exponent		product
//
//	12345										5								15				4045.2096
//								(0.12345 x 2^15 = 4045.2096)
// -------------------------------------------------------------
//
//    571										1									8			 14617.6
//                (57.1 x 2^8 = 14617.6)
// -------------------------------------------------------------
//
// Input Parameters
// ================
//
//
// multiplier						*big.Int - 'multiplier' will be multiplied by two to power of
//                       						'exponent' to generate the result or 'product'.
//
// multiplierPrecision	*big.Int	- The precision specification for 'multiplier'.
// 																	Precision specifies the number of digits to the
//                                  right of the decimal place in 'multiplier'. This
//                                  value must be greater than or equal to zero.
//
// exponent							uint			- Two will be raised to the power of exponent and
//                                  multiplied by 'multiplier' to generate the result
//                                  or product.
//
//
// Return Values
// =============
//
// product 							*big.Int	- The multiplication result. product will be set equal
//                                  to multiplier times two to the power of exponent.
//
// productPrecision			*big.Int  - The precision specification for 'product'. Precision
// 																	specifies the number of digits to the right of the
// 																	decimal place in 'product'. This value will always be
// 																	greater than or equal to zero.
//
//
// err									error			- If 'multiplierPrecision' is less than zero, an error
// 																	will be returned.
//
// Note: This method will delete trailing fractional zeros from
// 			 the returned result (product).
//
func (bMultiply BigIntMathMultiply) BigIntMultiplyByTwoToPower(
	multiplier,
	multiplierPrecision *big.Int,
	exponent uint) (product *big.Int, productPrecision *big.Int, err error) {

	product = big.NewInt(0)
	productPrecision = big.NewInt(0)
	err = nil

	if multiplier == nil {
		multiplier = big.NewInt(0)
	}

	if multiplierPrecision == nil {
		multiplierPrecision = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if multiplierPrecision.Cmp(bigZero) == -1 {
		ePrefix := "BigIntMathMultiply.BigIntMultiplyByTwoToPower() "
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'multiplierPrecision' is LESS THAN ZERO! "+
			"multiplierPrecision='%v' ", multiplierPrecision.Text(10))
		return product, productPrecision, err
	}

	if multiplier.Cmp(big.NewInt(0)) == 0 {
		product = big.NewInt(0)
		productPrecision = big.NewInt(0)
		err = nil
		return product, productPrecision, err
	}

	product = big.NewInt(0).Lsh(multiplier, exponent)
	productPrecision.Set(multiplierPrecision)

	// Delete trailing fractional zeros
	if productPrecision.Cmp(bigZero) == 1 {
		// productPrecision > 0
		scrap := big.NewInt(0)
		biBase10 := big.NewInt(10)
		bigOne := big.NewInt(1)
		biBaseZero := big.NewInt(0)
		newProduct, mod10 := big.NewInt(0).QuoRem(product, biBase10, scrap)

		for mod10.Cmp(biBaseZero) == 0 && productPrecision.Cmp(bigZero) == 1 {
			product.Set(newProduct)
			productPrecision.Sub(productPrecision, bigOne)
			newProduct, mod10 = big.NewInt(0).QuoRem(product, biBase10, scrap)
		}
	}

	return product, productPrecision, err
}

// FixedDecimalMultiply - This method receives two BigIntFixedDecimal
// types and then proceeds to perform a multiplication operation by
// multiplying the 'multiplier' by the 'multiplicand' to generate the
// 'product'.
//
// In the multiplication operation, the number to be multiplied is called
// the "multiplicand", while the number of times the multiplicand is to
// be multiplied comes from the "multiplier". Usually the multiplier is
// placed first and the multiplicand is placed second.
//
// Examples
// ========
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// 'multiplier', 'multiplicand' and 'product' are BigIntFixed Decimal types
// which may be used to defined fixed length floating point numbers.
//
// The BigIntFixedDecimal structure is defined as
// type BigIntFixedDecimal struct {
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
// 	To represent the floating point number 52.459
// 	a BigIntDecimal Structure would be configured as follows:
// 			BigIntFixedDecimal.integerNum	= 52459
// 			BigIntFixedDecimal.precision	= 3
//
//  Consider the following multiplication example:
// 						product =	752.314 x 21.67894 = 16309.37006716
//
// 'multiplier' and 'multiplicand' would be configured as follows:
//									multiplier.integerNum		= 752314
//                  multiplier.precision		= 3
//                  multiplicand.integerNum	= 2167894
//                  multiplicand.precision	= 5
//
// The 'product' would be calculated as follows:
//									product.integerNum	= 1630937006716
//                  product.precision  	= 8
//
// Input Parameters
// ================
//
//	multiplier BigIntFixedDecimal		- The number to be multiplied by 'multiplicand'
//
//	multiplicand BigIntFixedDecimal - The number to be multiplied by the 'multiplier'.
//
// Return Values
// =============
//
// product			BigIntFixedDecimal	- The product of the 'multiplier' multiplied by
//                                		the 'multiplicand'.
//
func (bMultiply BigIntMathMultiply) FixedDecimalMultiply(
	multiplier BigIntFixedDecimal,
	multiplicand BigIntFixedDecimal) (product BigIntFixedDecimal) {

	product = BigIntFixedDecimal{}.NewZero(0)

	multiplier.IsValid()
	multiplicand.IsValid()

	result, resultPrecision, _ :=
		BigIntMathMultiply{}.BigIntMultiply(
			multiplier.GetInteger(),
			multiplier.GetPrecisionBigInt(),
			multiplicand.GetInteger(),
			multiplicand.GetPrecisionBigInt())

	biMaxUint := big.NewInt(int64(math.MaxUint32))

	if resultPrecision.Cmp(biMaxUint) > 1 {
		delta := big.NewInt(0).Sub(resultPrecision, biMaxUint)
		delta.Sub(delta, big.NewInt(1))
		bigTen := big.NewInt(10)
		scale := big.NewInt(0).Exp(bigTen, delta, nil)
		result.Quo(result, scale)
		bigFive := big.NewInt(5)

		if result.Cmp(big.NewInt(0)) == -1 {
			bigFive.Neg(bigFive)
		}

		result.Add(result, bigFive)
		result.Quo(result, bigTen)
		resultPrecision.Set(biMaxUint)
	}

	product.SetNumericValue(result, uint(resultPrecision.Uint64()))
	return product
}

// MultiplyBigIntByTwoToPower - Multiplies a *big.Int number by powers
// of two and returns the result as a BigIntNum type.
//
// Examples:
// ========
//
// 						product = multiplier X 2^exponent
//
// ------------------------------------------------------------
//	multiplier			multiplierPrecision 	exponent		product
// ------------------------------------------------------------
//
//	12345										5								15				4045.2096
//								(0.12345 x 2^15 = 4045.2096)
// -------------------------------------------------------------
//
//  571											1									8			 14617.6
//                (57.1 x 2^8 = 14617.6)
//
// Note: This method will delete trailing fractional zeros from
// 			 the returned product.
//
// The returned BigIntNum multiplication 'result' will contain default numeric
// separators (decimal separator, thousands separator and currency symbol)
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntByTwoToPower(
	multiplier *big.Int,
	multiplierPrecision,
	exponent uint) BigIntNum {

	// error is disregarded because 'multiplierPrecision' will
	// always be greater than or equal to zero.
	product, productPrecision, _ :=
		BigIntMathMultiply{}.BigIntMultiplyByTwoToPower(
			multiplier,
			big.NewInt(0).SetUint64(uint64(multiplierPrecision)),
			exponent)

	// error is disregarded because productPrecision will never be
	// greater than multiplierPrecision.
	biNum, _ := BigIntNum{}.NewBigIntPrecision(product, productPrecision)

	return biNum
}

// New - Creates a BigIntMathMultiply instance with data
// variables initialized to zero.
//
func (bMultiply BigIntMathMultiply) New() BigIntMathMultiply {

	b2Math := BigIntMathMultiply{}

	b2Math.Input = BigIntPair{}.New()

	baseZero := big.NewInt(0)

	b2Math.Result = BigIntNum{}.NewBigInt(baseZero, 0)

	return b2Math
}

// NewNewBigIntPair - Creates a new BigIntMathMultiply based on input parameter
// type, 'BigIntPair'
//
func (bMultiply BigIntMathMultiply) NewBigIntPairResult(bPair BigIntPair) BigIntMathMultiply {

	b2Math := BigIntMathMultiply{}

	b2Math.Input = bPair.CopyOut()

	return b2Math
}

// MultiplyBigInts - Receives two *big.Int numbers and their associated precision
// specifications. This method then proceeds to perform a multiplication operation
// by multiplying the 'multiplier' by the 'multiplicand'
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// Input Parameters
// ================
//
//	multiplier *big.Int					- The number to be multiplied by 'multiplicand'
//	multiplierPrecision uint,		- The 'multiplier' precision or numeric digits after
//																	the decimal point.
//	multiplicand *big.Int,			- The number to be multiplied by the 'multiplier'.
//	multiplicandPrecision uint  - The 'multiplicand' precision or numeric digits after
//																	the decimal point.
//
// Return Values
// =============
//
// This method performs the multiplication operation and returns the result or
// 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain default numeric
// separators (decimal separator, thousands separator and currency symbol)
//
func (bMultiply BigIntMathMultiply) MultiplyBigInts(
	multiplier *big.Int,
	multiplierPrecision uint,
	multiplicand *big.Int,
	multiplicandPrecision uint) BigIntNum {

	product := big.NewInt(0).Mul(multiplier, multiplicand)

	biNumProduct := BigIntNum{}.NewBigInt(
		product, multiplierPrecision+multiplicandPrecision)

	biNumProduct.TrimTrailingFracZeros()

	return biNumProduct
}

// MultiplyBigIntNums - Receives two BigIntNum types as input parameters and then
// performs a 'multiplication' operation on these BigIntNum's.
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter, 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNums(
	multiplier BigIntNum,
	multiplicand BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(multiplier, multiplicand)

	finalResult := bMultiply.MultiplyPair(bPair)

	return finalResult
}

// MultiplyBigIntNumArray - Receives one BigIntNum which is classified as the 'multiplier'.
// The second input parameter is an array of BigIntNum Types labeled, 'multiplicands'. The
// first element of the 'multiplicands' array is multiplied by the 'multiplier' to produce
// a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next element
// in the multiplicands array. This process is continued through the last element in the array
// and the combined, final 'product' is returned as a Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter, 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumArray(
	multiplier BigIntNum,
	multiplicands []BigIntNum) BigIntNum {

	finalResult := multiplier.CopyOut()
	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New()
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for i := 0; i < lenMultiplicands; i++ {

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicands[i])

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// MultiplyBigIntNumOutputToArray - Receives one input parameter of Type BigIntNum which
// is classified as the 'multiplier'. The second input parameter is an array of BigIntNum
// Types labeled, 'multiplicands'.
//
// Each element of the 'multiplicands' array is multiplied by the 'multiplier'. The result or
// 'product' is then stored in a results array which is returned to the calling function.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier', 3 is the 'multiplicand'
// and 15 is the 'product' or result.
//
// Example
// =======
// 										Multiplicands												Output
// Multiplier				    	Array														Array
//
//		3			x				multiplicands[0] = 2			=				  outputarray[0] =  6
//		3			x				multiplicands[1] = 3			=				  outputarray[1] =  9
//		3			x				multiplicands[2] = 4			=				  outputarray[2] = 12
//		3			x				multiplicands[3] = 5			=				  outputarray[3] = 15
//		3			x				multiplicands[4] = 6			=				  outputarray[4] = 18
//		3			x				multiplicands[5] = 7			=				  outputarray[5] = 21
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' in an Array of 'BigIntNums' ([] BigIntNums).
//
// Each element of the returned BigIntNum array 'result' will contain
// numeric separators (decimal separator, thousands separator and
// currency symbol) copied from input parameter, 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumOutputToArray(
	multiplier BigIntNum,
	multiplicands []BigIntNum) []BigIntNum {

	bINumInterimResult := multiplier.CopyOut()

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return []BigIntNum{}
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	resultArray := make([]BigIntNum, lenMultiplicands)

	for i := 0; i < lenMultiplicands; i++ {

		bPair := BigIntPair{}.NewBigIntNum(bINumInterimResult, multiplicands[i])

		resultArray[i] = bMultiply.multiplyPairNoNumSeps(bPair)
		_ = resultArray[i].SetNumericSeparatorsDto(numSeps)
	}

	return resultArray
}

// MultiplyBigIntNumSeries - Receives one input parameter of Type BigIntNum which is classified
//  as the 'multiplier'. The second input parameter is a series of BigIntNum Types labeled,
// 'multiplicands'. The first element of the 'multiplicands' series is multiplied by the 'multiplier'
// to produce a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next
// element in the multiplicands series. This process is continued through the last element in the
// series. Afterwards, the combined final 'product' is returned as a Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumSeries(
	multiplier BigIntNum,
	multiplicands ...BigIntNum) BigIntNum {

	finalResult := multiplier.CopyOut()

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New()
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for _, multiplicand := range multiplicands {

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicand)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// MultiplyBigIntNumByTwo - Receives a BigIntNum input parameter 'base' and then
// proceeds to multiply this value times two (2).
//
// 								product = base X 2
//
// The product of this multiplication operation is returned as a BigIntNum.
// This returned BigIntNum 'product' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input
// parameter,'base'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByTwo(base BigIntNum) BigIntNum {

	result := big.NewInt(0).Lsh(base.bigInt, 1)

	bINumResult := BigIntNum{}.NewBigInt(result, base.GetPrecisionUint())

	_ = bINumResult.SetNumericSeparatorsDto(base.GetNumericSeparatorsDto())

	return bINumResult
}

// MultiplyBigIntNumByTwo - Receives a BigIntNum input parameter 'base' and then
// proceeds to multiply this value times two to the power of 'exponent'.
//
// 								product = base X 2^exponent
//
// The product of this multiplication operation is returned as a BigIntNum.
// This returned BigIntNum 'product' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input
// parameter,'base'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByTwoToPower(base BigIntNum, exponent uint) BigIntNum {

	result := big.NewInt(0).Lsh(base.bigInt, exponent)

	bINumResult := BigIntNum{}.NewBigInt(result, base.GetPrecisionUint())

	_ = bINumResult.SetNumericSeparatorsDto(base.GetNumericSeparatorsDto())

	return bINumResult
}

// MultiplyBigIntNumByThree - Receives a BigIntNum input parameter 'base' and then
// proceeds to multiply this value times three (3).
//
// 								product = base X 3
//
// The product of this multiplication operation is returned as a BigIntNum.
// This returned BigIntNum 'product' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input
// parameter,'base'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByThree(base BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(base, BigIntNum{}.NewThree(0))

	return bMultiply.MultiplyPair(bPair)
}

// MultiplyBigIntNumByFive - Receives a BigIntNum input parameter 'base' and then
// proceeds to multiply this value times five (5).
//
// 								product = base X 5
//
// The product of this multiplication operation is returned as a BigIntNum.
// This returned BigIntNum 'product' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input
// parameter,'base'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByFive(base BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(base, BigIntNum{}.NewFive(0))

	return bMultiply.MultiplyPair(bPair)
}

// MultiplyBigIntNumByTen - Receives a BigIntNum input parameter 'base' and then
// proceeds to multiply this value times ten (10).
//
// 								product = base X 10
//
// The product of this multiplication operation is returned as a BigIntNum.
// This returned BigIntNum 'product' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input
// parameter,'base'.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByTen(base BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(base, BigIntNum{}.NewTen(0))

	return bMultiply.MultiplyPair(bPair)
}

// MultiplyBigIntNumByTenToPower - Receives two BigIntNum input parameters, 'base'
// and 'tenExponent'. This method then proceeds to multiply 'base' time 10 to the
// exponent, 'tenExponent'.
//
//		result = base X 10^tenExponent
//
// The exponent can be a negative value and/or a fractional value.
//
// Input parameter 'maxPrecision' is used to control the maximum precision for the
// result returned by this method. Precision is defined as the the number of fractional
// digits to the right of the decimal place. Maximum precision therefore controls
// the maximum number of decimal digits to the right of the decimal place. If the
// returned result from this operation contains a number of fractional digits which
// is greater than 'maxPrecision' the result will be rounded to 'maxPrecision' decimal
// places. Be advised that these calculations can support very large precision values.
//
// Return Value
// ============
// The return value is of type BigIntNum and represents the result of the
// multiplication operation described above. This returned BigIntNum multiplication
// 'result' will contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter,'base'.
//
// If the precision of the return value precision exceeds input parameter 'maxPrecision',
// the return value will be rounded to 'maxPrecision' decimal places.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByTenToPower(
	base, tenExponent BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyByTenToPower() "

	bigINumTen := BigIntNum{}.NewTen(0)

	scale, err := BigIntMathPower{}.Pwr(bigINumTen, tenExponent, maxPrecision+20)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathPower{}.Pwr(bigINumTen, tenExponent, 100) "+
				"Error='%v'", err.Error())
	}

	// result numSeps are copied from from 'base'
	result := bMultiply.MultiplyBigIntNums(base, scale)

	if result.precision > maxPrecision {
		result.RoundToDecPlace(maxPrecision)
	}

	return result, nil
}

// MultiplyBigIntNumByTenToIntPower - Receives a BigIntNum input parameter, 'base'
// and a uint64 input parameter, 'tenExponent'. This method then proceeds to multiply
// 'base' time 10 to the exponent, 'tenExponent'.
//
//										result = base X 10^tenExponent
//
// Input parameter 'maxPrecision' is used to control the maximum precision for the
// result returned by this method. Precision is defined as the the number of fractional
// digits to the right of the decimal place. Maximum precision therefore controls
// the maximum number of decimal digits to the right of the decimal place. If the
// returned result from this operation contains a number of fractional digits which
// is greater than 'maxPrecision' the result will be rounded to 'maxPrecision' decimal
// places. Be advised that these calculations can support very large precision values.
//
// Return Value
// ============
// The return value is of type BigIntNum and represents the result of the
// multiplication operation described above. This returned BigIntNum multiplication
// 'result' will contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter,'base'.
//
// If the precision of the return value precision exceeds input parameter 'maxPrecision',
// the return value will be rounded to 'maxPrecision' decimal places.
//
func (bMultiply BigIntMathMultiply) MultiplyBigIntNumByTenToIntPower(
	base BigIntNum,
	tenExponent uint64,
	maxPrecision uint) BigIntNum {

	scale :=
		big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(tenExponent), nil)

	bINumScale := BigIntNum{}.NewBigInt(scale, 0)

	// result numSeps are copied from from 'base'
	result := bMultiply.MultiplyBigIntNums(base, bINumScale)

	if result.precision > maxPrecision {
		result.RoundToDecPlace(maxPrecision)
	}

	return result
}

// MultiplyDecimal - Receives two Decimal instances and multiplies their
// numeric values. The result or 'product' is returned as a 'BigIntNum'
// type.
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyDecimal(
	multiplier,
	multiplicand Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyDecimal() "

	// This method tests the validity of multiplier and multiplicand.
	bPair, err := BigIntPair{}.NewDecimal(multiplier, multiplicand)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewDecimal(multiplier, multiplicand). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bMultiply.MultiplyPair(bPair)

	return finalResult, nil
}

// MultiplyDecimalArray - Receives one Decimal which is classified as the 'multiplier'.
// The second input parameter is an array of Decimal Types labeled, 'multiplicands'. The
// first element of the 'multiplicands' array is multiplied by the 'multiplier' to produce
// a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next element
// in the multiplicands array. This process is continued through the last element in the array
// when the combined, final 'product' is returned as a Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyDecimalArray(
	multiplier Decimal,
	multiplicands []Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyDecimalArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method tests the validity of 'multiplier'.
	finalResult, err := BigIntNum{}.NewDecimal(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewDecimal(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for i := 0; i < lenMultiplicands; i++ {

		// This method tests the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewDecimal(multiplicands[i])

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewDecimal(multiplicands[i]) "+
					"multiplicands[%v]='%v' Error='%v'. ",
					i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return finalResult, nil
}

// MultiplyDecimalOutputToArray - Receives one input parameter of Type Decimal which
// is classified as the 'multiplier'. The second input parameter is an array of Decimal
// Types labeled, 'multiplicands'.
//
// Each element of the 'multiplicands' array is multiplied by the 'multiplier'. The result or
// 'product' is then stored in a results array which is returned to the calling function.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier', 3 is the 'multiplicand'
// and 15 is the 'product' or result.
//
// Example
// =======
// 										Multiplicands												Output
// Multiplier				    	Array														Array
//
//		3			x				multiplicands[0] = 2			=				  outputarray[0] =  6
//		3			x				multiplicands[1] = 3			=				  outputarray[1] =  9
//		3			x				multiplicands[2] = 4			=				  outputarray[2] = 12
//		3			x				multiplicands[3] = 5			=				  outputarray[3] = 15
//		3			x				multiplicands[4] = 6			=				  outputarray[4] = 18
//		3			x				multiplicands[5] = 7			=				  outputarray[5] = 21
//
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' in an Array of 'Decimals' ([] Decimals).
//
// The returned Decimal Array ([]Decimal) will contain numeric separators (decimal separator,
// thousands separator and currency symbol) copied from input parameter, 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyDecimalOutputToArray(
	multiplier Decimal,
	multiplicands []Decimal) ([]Decimal, error) {

	ePrefix := "BigIntMathMultiply.MultiplyDecimalOutputToArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return []Decimal{},
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method tests the validity of multiplier
	multiplierBINum, err := multiplier.GetBigIntNum()

	if err != nil {
		return []Decimal{},
			fmt.Errorf(ePrefix+
				"Error returned by multiplier.GetBigIntNum() "+
				" Error='%v'. ", err.Error())
	}

	resultArray := make([]Decimal, lenMultiplicands)

	for i := 0; i < lenMultiplicands; i++ {

		// This method tests the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewDecimal(multiplicands[i])

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewDecimal(multiplicands[i]) "+
					" multiplierBINum='%v' multiplicands[%v]='%v' Error='%v'. ",
					multiplierBINum.GetNumStr(), i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(multiplierBINum, multiplicandBINum)

		result := bMultiply.MultiplyPair(bPair)

		resultArray[i], err = result.GetDecimal()

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by result.GetDecimal() "+
					"index='%v' Error='%v'. ", i, err.Error())
		}
	}

	return resultArray, nil
}

// MultiplyDecimalSeries - Receives one input parameter of Type Decimal which is classified
//  as the 'multiplier'. The second input parameter is a comma delimited series of Decimal
// Types labeled, 'multiplicands'. The first element of the 'multiplicands' series is multiplied
// by the 'multiplier' to produce a 'product'. That 'product' replaces the 'multiplier' and is
// multiplied by the next element in the multiplicands series. This process is continued through
// the last element in the series. Afterwards, the combined final 'product' is returned as a Type
// 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyDecimalSeries(
	multiplier Decimal,
	multiplicands ...Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyDecimalSeries() "

	if len(multiplicands) == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands series is Empty!")
	}

	// This method will test the validity of 'multiplier'
	finalResult, err := BigIntNum{}.NewDecimal(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewDecimal(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for _, multiplicand := range multiplicands {

		// This method will test the validity of 'multiplicand'
		multiplicandBINum, err := BigIntNum{}.NewDecimal(multiplicand)

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewDecimal(multiplicand) "+
					" multiplicand='%v' Error='%v'. ",
					multiplicand.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyIntAry - Receives two IntAry instances and multiplies their
// numeric values. The result or 'product' is returned as a 'BigIntNum'
// type.
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// Be careful, IntAry's can accommodate very, very large numbers.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyIntAry(
	multiplier,
	multiplicand IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyIntAry() "

	// This method will test the validity of 'multiplier' and
	// 'multiplicand'.
	bPair, err := BigIntPair{}.NewIntAry(multiplier, multiplicand)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewIntAry(multiplier, multiplicand). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bMultiply.MultiplyPair(bPair)

	return finalResult, nil
}

// MultiplyIntAryArray - Receives one IntAry which is classified as the 'multiplier'.
// The second input parameter is an array of IntAry Types labeled, 'multiplicands'. The
// first element of the 'multiplicands' array is multiplied by the 'multiplier' to produce
// a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next element
// in the multiplicands array. This process is continued through the last element in the array
// when the combined, final 'product' is returned as a Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyIntAryArray(
	multiplier IntAry,
	multiplicands []IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyIntAryArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method will test the validity of 'multiplier'
	finalResult, err := BigIntNum{}.NewIntAry(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewIntAry(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for i := 0; i < lenMultiplicands; i++ {

		// This method will test the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewIntAry(multiplicands[i])

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(multiplicands[i]) "+
					"multiplicands[%v]='%v' Error='%v'. ",
					i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)

	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyIntAryOutputToArray - Receives one input parameter of Type IntAry which
// is classified as the 'multiplier'. The second input parameter is an array of IntAry
// Types labeled, 'multiplicands'.
//
// Each element of the 'multiplicands' array is multiplied by the 'multiplier'. The result or
// 'product' is then stored in a results array which is returned to the calling function.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier', 3 is the 'multiplicand'
// and 15 is the 'product' or result.
//
// Example
// =======
// 										Multiplicands												Output
// Multiplier				    	Array														Array
//
//		3			x				multiplicands[0] = 2			=				  outputarray[0] =  6
//		3			x				multiplicands[1] = 3			=				  outputarray[1] =  9
//		3			x				multiplicands[2] = 4			=				  outputarray[2] = 12
//		3			x				multiplicands[3] = 5			=				  outputarray[3] = 15
//		3			x				multiplicands[4] = 6			=				  outputarray[4] = 18
//		3			x				multiplicands[5] = 7			=				  outputarray[5] = 21
//
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' in an Array of 'IntArys' ([]IntAry).
//
// Each element in the returned []IntAry will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyIntAryOutputToArray(
	multiplier IntAry,
	multiplicands []IntAry) ([]IntAry, error) {

	ePrefix := "BigIntMathMultiply.MultiplyIntAryOutputToArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return []IntAry{},
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method will test the validity of 'multiplier'
	multiplierBINum, err := multiplier.GetBigIntNum()

	if err != nil {
		return []IntAry{},
			fmt.Errorf(ePrefix+
				"Error returned by multiplier.GetBigIntNum() "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	resultArray := make([]IntAry, lenMultiplicands)

	for i := 0; i < lenMultiplicands; i++ {

		// This method will test the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewIntAry(multiplicands[i])

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(multiplicands[i]) "+
					" multiplierBINum='%v' multiplicands[%v]='%v' Error='%v'. ",
					multiplierBINum.GetNumStr(), i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(multiplierBINum, multiplicandBINum)

		result := bMultiply.MultiplyPair(bPair)

		resultArray[i], err = result.GetIntAry()

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by finalResult.Result.GetIntAryElements() "+
					"i='%v' Error='%v'. ", i, err.Error())
		}
	}

	return resultArray, nil
}

// MultiplyIntArySeries - Receives one input parameter of Type IntAry which is classified
//  as the 'multiplier'. The second input parameter is a comma separated series of IntAry Types
// labeled, 'multiplicands'. The first element of the 'multiplicands' series is multiplied by
// the 'multiplier' to produce a 'product'. That 'product' replaces the 'multiplier' and is
// multiplied by the next element in the multiplicands series. This process is continued through
// the last element in the series. Afterwards, the combined final 'product' is returned as a
// Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyIntArySeries(
	multiplier IntAry,
	multiplicands ...IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyIntArySeries() "

	if len(multiplicands) == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands series is Empty!")
	}

	// This method will test the validity of 'multiplier'.
	finalResult, err := BigIntNum{}.NewIntAry(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewIntAry(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for _, multiplicand := range multiplicands {

		// This method will test the validity of 'multiplicand'
		multiplicandBINum, err := BigIntNum{}.NewIntAry(multiplicand)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(multiplicand) "+
					" multiplicand='%v' Error='%v'. ",
					multiplicand.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyNumStr - Receives two number strings and multiplies their numeric
// values.
//
// 							n1NumStr x n2NumStr = product
//
// The result or 'product' is returned as a 'BigIntNum'
// type.
//
// The 'n1NumStr' and 'n2NumStr' strings passed to this method are number strings
// which consist of a string of numeric digits representing a numeric value. A leading
// minus sign (-) may be included to indicate a negative numeric value. This string
// of numeric digits may also include a delimiting decimal separator to identify
// fractional digits. The number strings are parsed based on the decimal separator
// character specified by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings 'n1NumStr' and 'n2NumStr'. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// multiplication operation.
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators
// (decimal separator, thousands separator and currency symbol) specified by input
// parameter, 'numSeps'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStr(
	n1NumStr,
	n2NumStr string,
	numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStr() "

	numSeps.SetDefaultsIfEmpty()

	bPair, err := BigIntPair{}.NewNumStrWithNumSeps(n1NumStr, n2NumStr, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewNumStrWithNumSeps("+
				"n1NumStr, n2NumStr, numSeps). Error='%v' ", err.Error())
	}

	finalResult := bMultiply.MultiplyPair(bPair)

	return finalResult, nil
}

// MultiplyNumStrArray - Receives one number string which is classified as the 'multiplier'.
// The second input parameter is an array of number strings labeled, 'multiplicands'. The
// first element of the 'multiplicands' array is multiplied by the 'multiplier' to produce
// a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next element
// in the multiplicands array. This process is continued through the last element in the array
// when the combined, final 'product' is returned as a Type 'BigIntNum'.
//
// The 'multiplier' string and 'multiplicands' string array passed to this method are number
// strings which consist of a string of numeric digits representing a numeric value. A leading
// minus sign (-) may be included to indicate a negative numeric value. This string of numeric
// digits may also include a delimiting decimal separator to identify fractional digits. Number
// strings are parsed based on the decimal separator character specified by input parameter,
// 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to parse the number
// strings contained in input parameters, 'multiplier' and 'multiplicands'. 'numSeps'
// represents the applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this multiplication
// operation.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// cumulative result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) specified by input parameter, 'numSeps'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrArray(
	multiplier string,
	multiplicands []string,
	numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrArray() "

	numSeps.SetDefaultsIfEmpty()

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands array is Empty!")

	}

	finalResult, err := BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps) "+
				" multiplier='%v' numSeps='%v' Error='%v'. ",
				multiplier, numSeps.String(), err.Error())
	}

	for i := 0; i < lenMultiplicands; i++ {

		multiplicandBINum, err := BigIntNum{}.NewNumStrWithNumSeps(multiplicands[i], numSeps)

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplicands[i], numSeps) "+
					"multiplicands[%v]='%v' numSeps='%v' Error='%v'. ",
					i, multiplicands[i], numSeps.String(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyNumStrOutputToArray - Receives one input parameter of Type string which
// is classified as the 'multiplier'. The second input parameter is an array of strings
// labeled, 'multiplicands'.
//
// Each element of the 'multiplicands' array is multiplied by the 'multiplier'. The result or
// 'product' is then stored in a results array which is returned to the calling function.
//
// The 'multiplier' string and 'multiplicands' string array passed to this method are number
// strings which consist of a string of numeric digits representing a numeric value. A leading
// minus sign (-) may be included to indicate a negative numeric value. This string of numeric
// digits may include a delimiting decimal separator to identify fractional digits. The number
// strings are parsed based on the decimal separator character specified by input parameter,
// 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to parse the number
// strings contained in input parameters, 'multiplier' and 'multiplicands'. 'numSeps'
// represents the applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the results string array returned by this
// multiplication operation.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier', 3 is the 'multiplicand'
// and 15 is the 'product' or result.
//
// Example
// =======
// 										Multiplicands												Output
// Multiplier				    	Array														Array
//
//		3			x				multiplicands[0] = 2			=				  outputarray[0] =  6
//		3			x				multiplicands[1] = 3			=				  outputarray[1] =  9
//		3			x				multiplicands[2] = 4			=				  outputarray[2] = 12
//		3			x				multiplicands[3] = 5			=				  outputarray[3] = 15
//		3			x				multiplicands[4] = 6			=				  outputarray[4] = 18
//		3			x				multiplicands[5] = 7			=				  outputarray[5] = 21
//
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' in an Array of 'NumStrs' ([] NumStrs).
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrOutputToArray(
	multiplier string,
	multiplicands []string,
	numSeps NumericSeparatorDto) ([]string, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrOutputToArray() "

	numSeps.SetDefaultsIfEmpty()

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return []string{},
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	multiplierBINum, err := BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps)

	if err != nil {
		return []string{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps) "+
				" multiplier='%v' numSeps='%v' Error='%v'. ",
				multiplier, numSeps.String(), err.Error())
	}

	resultArray := make([]string, lenMultiplicands)

	for i := 0; i < lenMultiplicands; i++ {

		multiplicandBINum, err := BigIntNum{}.NewNumStrWithNumSeps(multiplicands[i], numSeps)

		if err != nil {
			return []string{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplicands[i], numSeps) "+
					" multiplicands[%v]='%v' numSeps='%v' Error='%v'. ",
					i, multiplicands[i], numSeps.String(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(multiplierBINum, multiplicandBINum)

		result := bMultiply.MultiplyPair(bPair)

		resultArray[i] = result.GetNumStr()

	}

	return resultArray, nil
}

// MultiplyNumStrSeries - Receives one input parameter of Type string which is classified
//  as the 'multiplier'. The second input parameter is a comma delimited series of strings labeled,
// 'multiplicands'. The first element of the 'multiplicands' series is multiplied by the 'multiplier'
// to produce a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next
// element in the multiplicands series. This process is continued through the last element in the
// series. Afterwards, the combined final 'product' is returned as a Type 'BigIntNum'.
//
// The 'multiplier' string and 'multiplicands' string series passed to this method are number
// strings which consist of a string of numeric digits representing a numeric value. A leading
// minus sign(-) may be included to indicate a negative numeric value. This string of numeric
// digits may also include a delimiting decimal separator to identify fractional digits. These
// number strings are parsed based on the decimal separator character specified by input
// parameter 'numSeps'.
//
// Input parameter 'numSeps' is of type NumericSeparatorDto and is used to parse the number
// strings contained in input parameters, 'multiplier' and 'multiplicands'. 'numSeps'
// represents the applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the results string array returned by this
// multiplication operation.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) specified by input parameter, 'numSeps'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrSeries(
	numSeps NumericSeparatorDto,
	multiplier string,
	multiplicands ...string) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrSeries() "

	numSeps.SetDefaultsIfEmpty()

	if len(multiplicands) == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands series is Empty!")
	}

	finalResult, err := BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplier, numSeps) "+
				" multiplier='%v' numSeps='%v' Error='%v'. ",
				multiplier, numSeps.String(), err.Error())
	}

	for _, multiplicand := range multiplicands {

		multiplicandBINum, err := BigIntNum{}.NewNumStrWithNumSeps(multiplicand, numSeps)

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(multiplicand, numSeps) "+
					" multiplicand='%v' numSeps='%v' Error='%v'. ",
					multiplicand, numSeps.String(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"numSeps='%v' Error='%v'. ", numSeps.String(), err.Error())
	}

	return finalResult, nil
}

// MultiplyNumStrDto - Receives two NumStrDto instances and multiplies their
// numeric values. The result or 'product' is returned as a 'BigIntNum'
// type.
//
// In the multiplication operation, the number to be multiplied is called the
// "multiplicand", while the number of times the multiplicand is to be multiplied
// comes from the "multiplier". Usually the multiplier is placed first and the
// multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrDto(
	multiplier,
	multiplicand NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrDto() "

	// This method tests the validity of 'multiplier' and 'multiplicand'
	bPair, err := BigIntPair{}.NewNumStrDto(multiplier, multiplicand)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewNumStrDto(multiplier, multiplicand). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bMultiply.MultiplyPair(bPair)

	return finalResult, nil
}

// MultiplyNumStrDtoArray - Receives one NumStrDto which is classified as the 'multiplier'.
// The second input parameter is an array of NumStrDto Types labeled, 'multiplicands'. The
// first element of the 'multiplicands' array is multiplied by the 'multiplier' to produce
// a 'product'. That 'product' replaces the 'multiplier' and is multiplied by the next element
// in the multiplicands array. This process is continued through the last element in the array
// when the combined, final 'product' is returned as a Type 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrDtoArray(
	multiplier NumStrDto,
	multiplicands []NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrDtoArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method tests the validity of 'multiplier'
	finalResult, err := BigIntNum{}.NewNumStrDto(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrDto(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for i := 0; i < lenMultiplicands; i++ {

		// This method tests the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewNumStrDto(multiplicands[i])

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(multiplicands[i]) "+
					"multiplicands[%v]='%v' Error='%v'. ",
					i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyNumStrDtoOutputToArray - Receives one input parameter of Type NumStrDto which
// is classified as the 'multiplier'. The second input parameter is an array of NumStrDto
// Types labeled, 'multiplicands'.
//
// Each element of the 'multiplicands' array is multiplied by the 'multiplier'. The result or
// 'product' is then stored in a results array which is returned to the calling function.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier', 3 is the 'multiplicand'
// and 15 is the 'product' or result.
//
// Example
// =======
// 										Multiplicands												Output
// Multiplier				    	Array														Array
//
//		3			x				multiplicands[0] = 2			=				  outputarray[0] =  6
//		3			x				multiplicands[1] = 3			=				  outputarray[1] =  9
//		3			x				multiplicands[2] = 4			=				  outputarray[2] = 12
//		3			x				multiplicands[3] = 5			=				  outputarray[3] = 15
//		3			x				multiplicands[4] = 6			=				  outputarray[4] = 18
//		3			x				multiplicands[5] = 7			=				  outputarray[5] = 21
//
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' in an Array of 'NumStrDtos' ([] NumStrDtos).
//
//
// Each element in the returned []NumStrDto array will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrDtoOutputToArray(
	multiplier NumStrDto,
	multiplicands []NumStrDto) ([]NumStrDto, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrDtoOutputToArray() "

	lenMultiplicands := len(multiplicands)

	if lenMultiplicands == 0 {
		return []NumStrDto{},
			errors.New(ePrefix + "Error: multiplicands array is Empty!")
	}

	// This method will test the validity of 'multiplier'
	multiplierBINum, err := BigIntNum{}.NewNumStrDto(multiplier)

	if err != nil {
		return []NumStrDto{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrDto(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	resultArray := make([]NumStrDto, lenMultiplicands)

	for i := 0; i < lenMultiplicands; i++ {

		// This method will test the validity of multiplicands[i]
		multiplicandBINum, err := BigIntNum{}.NewNumStrDto(multiplicands[i])

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(multiplicands[i]) "+
					" multiplierBINum='%v' multiplicands[%v]='%v' Error='%v'. ",
					multiplierBINum.GetNumStr(), i, multiplicands[i].GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(multiplierBINum, multiplicandBINum)

		result := bMultiply.multiplyPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps) "+
					"Error='%v'. ", err.Error())
		}

		resultArray[i], err = result.GetNumStrDto()

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by finalResult.Result.GetNumStrDto() "+
					"i='%v' Error='%v'. ", i, err.Error())
		}
	}

	return resultArray, nil
}

// MultiplyNumStrDtoSeries - Receives one input parameter of Type NumStrDto which is classified
//  as the 'multiplier'. The second input parameter is a comma delimited series of NumStrDto
// Types labeled, 'multiplicands'. The first element of the 'multiplicands' series is multiplied
// by the 'multiplier' to produce a 'product'. That 'product' replaces the 'multiplier' and is
// multiplied by the next element in the multiplicands series. This process is continued through
// the last element in the series. Afterwards, the combined final 'product' is returned as a Type
// 'BigIntNum'.
//
// In the multiplication operation, the number to be multiplied is called the "multiplicand",
// while the number of times the multiplicand is to be multiplied comes from the "multiplier".
// Usually the multiplier is placed first and the multiplicand is placed second.
//
// For example, in the problem 5 x 3 equals 15, the 5 is the 'multiplier',
// 3 is the 'multiplicand' and 15 is the 'product' or result.
//
//							multiplier x multiplicand = product or result
//
// This method performs the multiplication operation described above and afterwards returns the
// result or 'product' as a BigIntNum type.
//
// The returned BigIntNum multiplication 'result' will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'multiplier'.
//
func (bMultiply BigIntMathMultiply) MultiplyNumStrDtoSeries(
	multiplier NumStrDto,
	multiplicands ...NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathMultiply.MultiplyNumStrDtoSeries() "

	if len(multiplicands) == 0 {
		return BigIntNum{}.New(),
			errors.New(ePrefix + "Error: multiplicands series is Empty!")
	}

	// This method will test the validity of 'multiplier'
	finalResult, err := BigIntNum{}.NewNumStrDto(multiplier)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrDto(multiplier) "+
				" multiplier='%v' Error='%v'. ",
				multiplier.GetNumStr(), err.Error())
	}

	numSeps := multiplier.GetNumericSeparatorsDto()

	for _, multiplicand := range multiplicands {

		// This method will test the validity of multiplicand
		multiplicandBINum, err := BigIntNum{}.NewNumStrDto(multiplicand)

		if err != nil {
			return BigIntNum{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(multiplicand) "+
					" multiplicand='%v' Error='%v'. ",
					multiplicand.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, multiplicandBINum)

		finalResult = bMultiply.multiplyPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v'. ", err.Error())
	}

	return finalResult, nil
}

// MultiplyPair - Receives a BigIntPair instance and proceeds to multiply
// bPair.Big1 by bPair.Big2. Both 'Big1' and 'Big2' are of type 'BigIntNum'.
//
// bPair.Big1 x bPair.Big2 = Result
//
// The result of this multiplication operation is returned as a BigIntNum
// type.
//
// The returned BigIntNum multiplication 'Result' will contain numeric
// separators (decimal separator, thousands separator and currency symbol)
// copied from bPair.Big1.
//
func (bMultiply BigIntMathMultiply) MultiplyPair(bPair BigIntPair) BigIntNum {

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	finalResult := bMultiply.multiplyPairNoNumSeps(bPair)

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// multiplyPairNoNumSeps - Receives a BigIntPair instance and proceeds to multiply
// bPair.Big1 by bPair.Big2. Both 'Big1' and 'Big2' are of type 'BigIntNum'.
//
// bPair.Big1 x bPair.Big2 = Result
//
// The result of this multiplication operation is returned as a BigIntNum
// type.
//
// The returned BigIntNum multiplication 'Result' will contain default numeric
// separators (decimal separator, thousands separator and currency symbol).
//
func (bMultiply BigIntMathMultiply) multiplyPairNoNumSeps(bPair BigIntPair) BigIntNum {

	b3 := big.NewInt(0).Mul(bPair.GetBig1BigInt(), bPair.GetBig2BigInt())

	bResult := BigIntNum{}.NewBigInt(
		b3,
		bPair.Big1.GetPrecisionUint()+bPair.Big2.GetPrecisionUint())

	bResult.TrimTrailingFracZeros()

	return bResult
}
