package mathops

import (
	"errors"
	"fmt"
	"math/big"
	"time"
)

/*

Calculate Any Logarithm manually:
	https://www.youtube.com/watch?v=TUAFzuMVem0

Properties of Logarithms:
https://www.youtube.com/watch?v=AAW7WRFBKdw&t=6s

https://en.wikipedia.org/wiki/Logarithm
log[b](x) = y exactly if b^y = x

log[10](450) = 2.6532125137753436793763169117857
precision = 31

LogN(1500000;4) = log[4](1500000) = 10.2582655350227

*/

type BigIntMathLogarithms struct {
	Base      BigIntNum
	XNumber   BigIntNum
	YExponent BigIntNum
}

// BigIntNumLogBaseOfX - Computes the log[base](xNum). 'base' and 'xNum' are passed as
// BigIntNum types.
//
// 'maxPrecision' is an uint specifying the precision or integerNum of digits to the
// right of the decimal place in the result.
//
//
// Calculation Formula:
// ====================
//			log[base](x) = y
//
//			x = base^y
//			'y' is the exponent of base required to set base^y equal to x.
// 			The calculation result 'y' of log[base](xNum) is returned as a BigIntNum.
//
// Input Parameters:
// =================
//
// base		BigIntNum	- The base of the logarithm. 'base' must be an integer integerNum
// 										greater than one (base > +1).
//
// xNum		BigIntNum - The 'x' value such that base^y = 'x'.  xNum must be greater
//										than zero.
//
// Return Values:
// ==============
// y			BigIntNum - If the function completes successfully, the result 'y' will be
//										returned as a BigIntNum.
//
// 				error			- If the function encounters an error, an error message will be returned
//										as type 'error'. If the function completes successfully, this value
//										will be set to 'nil'.
//
//          log Value:  10.25826553502266513433782776044
// Expected log Value:  10.2582655350227
//
// base:  4
// xNum:  1500000
//
func (bLog BigIntMathLogarithms) BigIntNumLogBaseOfX(
	base, xNum BigIntNum, maxPrecision uint) (result BigIntNum, err error) {

	ePrefix := "BigIntMathLogarithms.BigIntNumLogBaseOfX() "

	biMaxPrecision := big.NewInt(0).SetUint64(uint64(maxPrecision))
	factor := big.NewInt(4)

	cycles := big.NewInt(40)
	if biMaxPrecision.Cmp(big.NewInt(10)) == 1 {
		cycles = big.NewInt(0).Mul(biMaxPrecision, factor)
	}

	maxInternalPrecision := big.NewInt(0).Mul(cycles, factor)

	biResult, biResultPrecision, errX :=
		BigIntMathLogarithms{}.LogBaseOfXByDivide(
			base.GetIntegerValue(),
			base.GetPrecisionBigInt(),
			xNum.GetIntegerValue(),
			xNum.GetPrecisionBigInt(),
			maxInternalPrecision,
			biMaxPrecision,
			cycles)

	if errX != nil {
		err = fmt.Errorf(ePrefix+
			"%v", errX.Error())
	}

	result, errX =
		BigIntNum{}.NewBigIntPrecision(
			biResult,
			biResultPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+
			"%v", errX.Error())
	}

	return result, nil
}

//
//          log Value:  10.2582655350226651343378277604421
// Expected log Value:  10.2582655350227
// base:  4
// xNum:  1500000
func (bLog BigIntMathLogarithms) BigIntLogBaseOfX(
	base,
	basePrecision,
	xNum,
	xNumPrecision,
	maxPrecision *big.Int) (logResult,
	logResultPrecision *big.Int,
	err error) {

	logResult = big.NewInt(0)
	logResultPrecision = big.NewInt(0)
	err = nil
	ePrefix := "BigIntMathLogarithms.BigIntLogBaseOfX() "

	bigZero := big.NewInt(0)

	if base == nil {
		base = big.NewInt(0)
	}

	if base.Cmp(bigZero) != 1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'base' is INVALID! "+
			"base='%v' ", base.Text(10))

		return logResult, logResultPrecision, err
	}

	if basePrecision == nil {
		basePrecision = big.NewInt(0)
	}

	if basePrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'basePrecision' is negative and INVALID! "+
			"basePrecision='%v' ", basePrecision.Text(10))

		return logResult, logResultPrecision, err

	}

	if xNum == nil {
		xNum = big.NewInt(0)
	}

	if xNum.Cmp(bigZero) == 0 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'xNum' is INVALID! "+
			"xNum='%v' ", xNum.Text(10))

		return logResult, logResultPrecision, err
	}

	if xNumPrecision == nil {
		xNumPrecision = big.NewInt(0)
	}

	if xNumPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'xNumPrecision' is LESS THAN ZERO! "+
			"xNumPrecision='%v' ", xNumPrecision.Text(10))

		return logResult, logResultPrecision, err
	}

	if maxPrecision == nil {
		maxPrecision = big.NewInt(0)
	}

	if maxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input Parameter 'maxPrecision' is LESS THAN ZERO! "+
			"maxPrecision='%v' ", maxPrecision.Text(10))

		return logResult, logResultPrecision, err
	}

	factor := big.NewInt(4)
	cycles := big.NewInt(40)

	if maxPrecision.Cmp(big.NewInt(10)) == 1 {
		cycles = big.NewInt(0).Mul(maxPrecision, factor)
	}

	maxInternalPrecision := big.NewInt(0).Mul(cycles, factor)

	var errX error

	logResult, logResultPrecision, errX =
		bLog.LogBaseOfXByDivide(
			base,
			basePrecision,
			xNum,
			xNumPrecision,
			maxInternalPrecision,
			maxPrecision,
			cycles)

	if errX != nil {
		logResult.Set(bigZero)
		logResultPrecision.Set(bigZero)
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return logResult, logResultPrecision, err
	}

	err = nil

	return logResult, logResultPrecision, err
}

// LogBaseOfXByDivide - Generates logs for specified
// bases and X-numbers. Calculation methodology derrived
// from following:
//
// Calculate Any Logarithm manually -
// https://www.youtube.com/watch?v=TUAFzuMVem0
//
func (bLog BigIntMathLogarithms) LogBaseOfXByDivide(
	base,
	basePrecision,
	xNum,
	xNumPrecision,
	maxInternalPrecision,
	maxPrecision,
	cycles *big.Int) (logResult,
	logResultPrecision *big.Int,
	err error) {

	ePrefix := "BigIntMathLogarithms.LogBaseOfXByDivide() "
	logResult = big.NewInt(0)
	logResultPrecision = big.NewInt(0)
	err = nil
	bigZero := big.NewInt(0)

	if base == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'base' is nil!")
		return logResult, logResultPrecision, err
	}

	if basePrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'basePrecision' is nil!")
		return logResult, logResultPrecision, err
	}

	if xNum == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'xNum' is nil!")
		return logResult, logResultPrecision, err
	}

	if xNumPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'xNumPrecision' is nil!")
		return logResult, logResultPrecision, err
	}

	if maxInternalPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'maxInternalPrecision' is nil!")
		return logResult, logResultPrecision, err
	}

	if maxPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'maxPrecision' is nil!")
		return logResult, logResultPrecision, err
	}

	if cycles == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'cycles' is nil!")
		return logResult, logResultPrecision, err
	}

	if basePrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'basePrecision' is LESS THAN ZERO! "+
			"basePrecision='%v' ", basePrecision.Text(10))
		return logResult, logResultPrecision, err
	}

	if xNumPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'xNumPrecision' is LESS THAN ZERO! "+
			"xNumPrecision='%v' ", xNumPrecision.Text(10))
		return logResult, logResultPrecision, err
	}

	if maxInternalPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'maxInternalPrecision' is LESS THAN ZERO! "+
			"maxInternalPrecision='%v' ", maxInternalPrecision.Text(10))
		return logResult, logResultPrecision, err
	}

	if maxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'maxPrecision' is LESS THAN ZERO! "+
			"maxPrecision='%v' ", maxPrecision.Text(10))
		return logResult, logResultPrecision, err
	}

	if cycles.Cmp(big.NewInt(5)) == -1 {
		cycles.Mul(maxPrecision, big.NewInt(5))
	}

	tXNum := big.NewInt(0).Set(xNum)
	tXNumPrecision := big.NewInt(0).Set(xNumPrecision)
	bigOne := big.NewInt(1)

	iBase, iBasePrecision, errX :=
		BigIntMathDivide{}.BigIntFracQuotient(
			bigOne,
			big.NewInt(0),
			base,
			basePrecision,
			maxInternalPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX)
		return logResult, logResultPrecision, err
	}

	ri := big.NewInt(0)
	riPrecision := big.NewInt(0)

	cmpNums :=
		BigIntMath{}.BigIntPrecisionCmp(
			tXNum,
			tXNumPrecision,
			base,
			basePrecision)

	for cmpNums == 1 {

		tXNum, tXNumPrecision, errX =
			BigIntMathMultiply{}.BigIntMultiply(
				tXNum,
				tXNumPrecision,
				iBase,
				iBasePrecision)

		if errX != nil {
			err = fmt.Errorf(ePrefix+"%v", errX)
			return logResult, logResultPrecision, err
		}

		tXNum, tXNumPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				tXNum,
				tXNumPrecision,
				maxInternalPrecision,
				true)

		if errX != nil {
			err = fmt.Errorf(ePrefix+
				"%v", errX.Error())
			return logResult, logResultPrecision, err
		}

		ri.Add(ri, bigOne)

		cmpNums =
			BigIntMath{}.BigIntPrecisionCmp(
				tXNum,
				tXNumPrecision,
				base,
				basePrecision)
	}

	p := big.NewInt(1)
	pPrecision := big.NewInt(0)
	oneHalf := big.NewInt(5)
	oneHalfPrecision := big.NewInt(1)
	uint64Cycles := cycles.Uint64()

	for i := uint64(0); i < uint64Cycles; i++ {

		tXNum.Mul(tXNum, tXNum)
		tXNumPrecision.Add(tXNumPrecision, tXNumPrecision)

		tXNum, tXNumPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				tXNum,
				tXNumPrecision,
				maxInternalPrecision,
				true)

		if errX != nil {
			err = fmt.Errorf(ePrefix+
				"%v", errX.Error())
			return logResult, logResultPrecision, err
		}

		p.Mul(p, oneHalf)
		pPrecision.Add(pPrecision, oneHalfPrecision)

		p, pPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				p,
				pPrecision,
				maxInternalPrecision,
				true)

		if errX != nil {
			err = fmt.Errorf(ePrefix+
				"%v", errX.Error())
			return logResult, logResultPrecision, err
		}

		cmpNums =
			BigIntMath{}.BigIntPrecisionCmp(
				tXNum,
				tXNumPrecision,
				base,
				basePrecision)

		if cmpNums == 1 {

			ri, riPrecision, errX =
				BigIntMathAdd{}.BigIntAdd(
					ri,
					riPrecision,
					p,
					pPrecision)

			tXNum.Mul(tXNum, iBase)
			tXNumPrecision.Add(tXNumPrecision, iBasePrecision)

			tXNum, tXNumPrecision, errX =
				BigIntMath{}.RoundToMaxPrecision(
					tXNum,
					tXNumPrecision,
					maxInternalPrecision,
					true)

			if errX != nil {
				err = fmt.Errorf(ePrefix+
					"%v", errX.Error())
				return logResult, logResultPrecision, err
			}
		}
	}

	logResult, logResultPrecision, errX =
		BigIntMath{}.RoundToMaxPrecision(
			ri,
			riPrecision,
			maxPrecision,
			true)

	if errX != nil {
		logResult.Set(bigZero)
		logResultPrecision.Set(bigZero)
		err = fmt.Errorf(ePrefix+
			"%v", errX.Error())
		return logResult, logResultPrecision, err
	}

	err = nil

	return logResult, logResultPrecision, err
}

/*
func (bLog BigIntMathLogarithms) LogBaseOfXByDivide(
	base,
	basePrecision,
	xNum,
	xNumPrecision,
	maxInternalPrecision,
	maxPrecision ,
	cycles *big.Int) (logResult,
													logResultPrecision *big.Int,
													err error) {

	ePrefix := "BigIntMathLogarithms.LogBaseOfXByDivide() "
	logResult = big.NewInt(0)
	logResultPrecision = big.NewInt(0)
	err = nil
	tXNum := big.NewInt(0).Set(xNum)
	tXNumPrecision := big.NewInt(0).Set(xNumPrecision)
	bigOne := big.NewInt(1)

	iBase, iBasePrecision, errX :=
		BigIntMathDivide{}.BigIntFracQuotient(
			bigOne,
			big.NewInt(0),
			base,
			basePrecision,
			maxInternalPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix + "%v", errX)
		return logResult, logResultPrecision, err
	}

	ri := big.NewInt(0)
	riPrecision := big.NewInt(0)

	cmpNums :=
		BigIntMath{}.BigIntPrecisionCmp(
			tXNum,
			tXNumPrecision,
			base,
			basePrecision)

	for cmpNums == 1 {

		tXNum, tXNumPrecision, errX =
			BigIntMathMultiply{}.BigIntMultiply(
				tXNum,
				tXNumPrecision,
				iBase,
				iBasePrecision)

		if errX != nil {
			err = fmt.Errorf(ePrefix + "%v", errX)
			return logResult, logResultPrecision, err
		}

		tXNum, tXNumPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				tXNum,
				tXNumPrecision,
				maxInternalPrecision,
				true)

		if errX != nil {
			err = fmt.Errorf(ePrefix +
				"%v", errX.Error())
			return logResult, logResultPrecision, err
		}

		ri.Add(ri, bigOne	)

		cmpNums =
			BigIntMath{}.BigIntPrecisionCmp(
				tXNum,
				tXNumPrecision,
				base,
				basePrecision)
	}

	p := big.NewInt(1)
	pPrecision := big.NewInt(0)
	oneHalf := big.NewInt(5)
	oneHalfPrecision := big.NewInt(1)
	uint64Cycles := cycles.Uint64()

	for i := uint64(0); i < uint64Cycles; i++ {

		tXNum, tXNumPrecision, errX =
			BigIntMathMultiply{}.BigIntMultiply(
				tXNum,
				tXNumPrecision,
				tXNum,
				tXNumPrecision)

		if errX != nil {
			err = fmt.Errorf(ePrefix + "%v", errX)
			return logResult, logResultPrecision, err
		}

		p, pPrecision, errX =
			BigIntMathMultiply{}.BigIntMultiply(
				p,
				pPrecision,
				oneHalf,
				oneHalfPrecision)

		if errX != nil {
			err = fmt.Errorf(ePrefix + "%v", errX)
			return logResult, logResultPrecision, err
		}

		cmpNums =
			BigIntMath{}.BigIntPrecisionCmp(
				tXNum,
				tXNumPrecision,
				base,
				basePrecision)

		if cmpNums==1 {

			ri, riPrecision, errX =
				BigIntMathAdd{}.BigIntAdd(
					ri,
					riPrecision,
					p,
					pPrecision)

			tXNum, tXNumPrecision, errX =
				BigIntMathMultiply{}.BigIntMultiply(
					tXNum,
					tXNumPrecision,
					iBase,
					iBasePrecision)

			if errX != nil {
				err = fmt.Errorf(ePrefix + "%v", errX)
				return logResult, logResultPrecision, err
			}

		}

		tXNum, tXNumPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				tXNum,
				tXNumPrecision,
				maxInternalPrecision,
				true)

		if errX != nil {
			err = fmt.Errorf(ePrefix +
				"%v", errX.Error())
			return logResult, logResultPrecision, err
		}

	}

	logResult, logResultPrecision, errX =
		BigIntMath{}.RoundToMaxPrecision(
			ri,
			riPrecision,
			maxPrecision,
			true)

	if errX != nil {
		logResult.Set(big.NewInt(0))
		logResultPrecision.Set(big.NewInt(0))
		err = fmt.Errorf(ePrefix +
			"%v", errX.Error())
		return logResult, logResultPrecision, err
	}

	err = nil

	return logResult, logResultPrecision, err
}
*/

// EPwrXFromMaclaurinSeries = Computes value of e^exponent where e is the mathematical
// constant, "Euler's integerNum".
//
// "Euler's integerNum" or 'e' truncated to 50 decimal places is:
//     2.71828182845904523536028747135266249775724709369995
// Rounded to 50-decimal places is:
//     2.71828182845904523536028747135266249775724709369996
// Reference: https://en.wikipedia.org/wiki/E_(mathematical_constant)
//
// 2.7182818284590452353602874713526624977572470936999595749669676277240766303
// 5354759457138217852516642742746
// Reference: https://oeis.org/A001113/list
//
// This function works best where 'exponent' is near zero.
//
// Input Parameters
// ================
// exponent	BigIntNum - The exponent to which the mathematical constant 'e' will
//                      be raised in order to compute the value of e^exponent.
//
// n				int64			- The integerNum of cycles in the Taylor series which will
//                      be used to compute the value of e^exponent
//
// Return Values
// =============
//
// BigIntNum					- If successful, this function returns the value of e^exponent
//											as a type BigIntNum.
//
// error							- If an error is encountered an error message will be formatted
//											and returned as type 'error'. If the function successfully
//											completes the calculation, this value will be set to 'nil'.
//
func (bLog BigIntMathLogarithms) EPwrXFromMaclaurinSeries(
	exponent BigIntNum, nCycles int64) (BigIntNum, error) {

	ePrefix := "BigIntMathLogarithms.EPwrXFromMaclaurinSeries() "
	// sum = 1 + x
	sum := BigIntMathAdd{}.AddBigIntNums(BigIntNum{}.NewInt(1, 0), exponent)
	nFactorial := BigIntNum{}.NewInt(1, 0)

	x := exponent.CopyOut()

	for n := int64(2); n <= nCycles; n++ {

		nFactorial = BigIntMathMultiply{}.MultiplyBigIntNums(nFactorial, BigIntNum{}.NewInt64(n, 0))

		x = BigIntMathMultiply{}.MultiplyBigIntNums(x, exponent)

		factor, err := BigIntMathDivide{}.BigIntNumFracQuotient(x, nFactorial, 1200)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient(binOne, nFactorial, 200) "+
					"nFactorial='%v' Error='%v' ",
					nFactorial.GetNumStr(), err.Error())
		}

		sum = BigIntMathAdd{}.AddBigIntNums(sum, factor)

		if sum.GetPrecisionUint() > 500 {
			sum.RoundToDecPlace(500)
		}

	}

	return sum, nil
}

func (bLog BigIntMathLogarithms) EPwrXFromTaylorSeries(
	exponent, a BigIntNum, nCycles int64) (BigIntNum, error) {

	ePrefix := "BigIntMathLogarithms.EPwrXFromTaylorSeries() "

	e := eulersNumber1k.GetBigIntNum()

	if e.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "eulersNumber1k is ZERO!")
	}

	aValue, err := a.GetUInt()

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by a.GetUInt. Error='%v'",
				err.Error())

	}

	internalMaxPrecision := uint(20000)

	outputMaxPrecision := uint(1100)

	ePwrBigInt, ePwrBigIntPrecision :=
		BigIntMathPower{}.BigIntPwrIteration(
			e.GetIntegerValue(),
			e.GetPrecisionUint(),
			aValue,
			internalMaxPrecision,
			outputMaxPrecision)

	eToPwr := BigIntNum{}.NewBigInt(ePwrBigInt, ePwrBigIntPrecision)

	// eToPwr, err := BigIntMathPower{}.Pwr(e, a, 500)

	//eToPwr := BigIntMathMultiply{}.MultiplyBigIntNums(e.CopyOut(), e.CopyOut())

	// sum = 0
	sum := BigIntNum{}.NewInt(0, 0)

	x := exponent.CopyOut()

	xMinusA := BigIntMathSubtract{}.SubtractBigIntNums(x, a)

	xMinusANth := BigIntNum{}.NewInt(0, 0)

	nFact := BigIntNum{}.NewInt(0, 0)

	for n := int64(0); n <= nCycles; n++ {

		if n == 0 {
			xMinusANth = BigIntNum{}.NewInt(1, 0)
			nFact = BigIntNum{}.NewInt(1, 0)
		} else if n == 1 {
			xMinusANth = xMinusA.CopyOut()
			nFact = BigIntNum{}.NewInt64(n, 0)
		} else {
			xMinusANth = BigIntMathMultiply{}.MultiplyBigIntNums(xMinusANth, xMinusA)
			nFact = BigIntMathMultiply{}.MultiplyBigIntNums(nFact, BigIntNum{}.NewInt64(n, 0))
		}

		factor1, err := BigIntMathDivide{}.BigIntNumFracQuotient(eToPwr, nFact, 500)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient(eToPwr, nFact, 500) "+
					"eToPwr='%v' nFact='%v' Error='%v' ",
					eToPwr.GetNumStr(), nFact.GetNumStr(), err.Error())
		}

		factor2 := BigIntMathMultiply{}.MultiplyBigIntNums(factor1, xMinusANth)

		sum = BigIntMathAdd{}.AddBigIntNums(sum, factor2)

	}

	sum.RoundToDecPlace(500)

	return sum, nil

}

// EPwrXFromTaylorSeriesFixedDecimal
func (bLog BigIntMathLogarithms) EPwrXFromTaylorSeriesFixedDecimal(
	exponentX BigIntFixedDecimal,
	a, nCycles uint) (BigIntFixedDecimal, error) {

	ePrefix := "BigIntMathLogarithms.EPwrXFromTaylorSeries() "

	e := eulersNumber1k.GetFixedDecimal()

	if e.IsZero() {
		return BigIntFixedDecimal{}.NewZero(0),
			errors.New(ePrefix + "EulersNumber1050 Constant is ZERO!")
	}

	internalMaxPrecision := uint(20000)

	outputMaxPrecision := uint(1500)

	ePwrBigInt, ePwrBigIntPrecision :=
		BigIntMathPower{}.BigIntPwrIteration(
			e.GetInteger(),
			e.GetPrecision(),
			a,
			internalMaxPrecision,
			outputMaxPrecision)

	eToPwr := BigIntFixedDecimal{}.New(ePwrBigInt, ePwrBigIntPrecision)

	fixedDecA := BigIntFixedDecimal{}.New(big.NewInt(int64(a)), 0)

	// sum = 0
	sum := BigIntFixedDecimal{}.NewZero(0)

	xNum := exponentX.CopyOut()

	xMinusA := BigIntMathSubtract{}.FixedDecimalSubtract(xNum, fixedDecA)

	xMinusANth := BigIntFixedDecimal{}.NewZero(0)

	nFact := BigIntFixedDecimal{}.NewZero(0)

	for n := uint(0); n < nCycles; n++ {

		if n == 0 {

			xMinusANth = BigIntFixedDecimal{}.NewInt(1, 0)
			nFact = BigIntFixedDecimal{}.NewInt(1, 0)

		} else if n == 1 {

			xMinusANth = xMinusA.CopyOut()
			nFact = BigIntFixedDecimal{}.NewUInt(n, 0)

		} else {

			xMinusANth = BigIntMathMultiply{}.FixedDecimalMultiply(xMinusANth, xMinusA)

			nFact =
				BigIntMathMultiply{}.FixedDecimalMultiply(
					nFact,
					BigIntFixedDecimal{}.New(big.NewInt(int64(n)), 0))

		}

		factor1, err := BigIntMathDivide{}.FixedDecimalFracQuotient(eToPwr, nFact, 500)

		if err != nil {
			return BigIntFixedDecimal{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient(eToPwr, nFact, 500) "+
					"eToPwr='%v' nFact='%v' Error='%v' ",
					eToPwr.GetNumStr(), nFact.GetNumStr(), err.Error())
		}

		factor2 := BigIntMathMultiply{}.FixedDecimalMultiply(factor1, xMinusANth)

		sum = BigIntMathAdd{}.FixedDecimalAdd(sum, factor2)

	}

	sum.RoundToDecPlace(500)

	return sum, nil
}

func (bLog BigIntMathLogarithms) EPwrXFromTaylorSeriesBigInt(
	exponentX,
	exponentXPrecision,
	a,
	nCycles,
	internalMaxPrecision,
	maxPrecision *big.Int) (xNum, xNumPrecision *big.Int, err error) {

	xNum = big.NewInt(0)
	xNumPrecision = big.NewInt(0)
	err = nil

	ePrefix := "BigIntMathLogarithms.EPwrXFromTaylorSeries() "

	e, ePrecision := eulersNumber1k.GetBigIntPrecision()

	bigZero := big.NewInt(0)

	if e.Cmp(bigZero) == 0 {
		err = errors.New(ePrefix + "EulersNumber1050 Constant is ZERO!")
		return xNum, xNumPrecision, err
	}

	eToPwr, eToPwrPrecision, errX :=
		BigIntMathPower{}.BigIntegerPwrIteration(
			e,
			ePrecision,
			a,
			internalMaxPrecision,
			maxPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return xNum, xNumPrecision, err
	}

	xNum = big.NewInt(0)
	xNumPrecision = big.NewInt(0)

	xMinusA, xMinusAPrecision, errX :=
		BigIntMathSubtract{}.BigIntSubtract(exponentX, exponentXPrecision, a, big.NewInt(0))

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return xNum, xNumPrecision, err
	}

	xMinusANth := big.NewInt(0)
	xMinusANthPrecision := big.NewInt(0)

	nFact := big.NewInt(0)
	nFactPrecision := big.NewInt(0)
	bigOne := big.NewInt(1)
	factor1 := big.NewInt(0)
	factor1Precision := big.NewInt(0)
	factor2 := big.NewInt(0)
	factor2Precision := big.NewInt(0)

	for n := big.NewInt(0); n.Cmp(nCycles) == -1; n.Add(n, bigOne) {

		if n.Cmp(bigZero) == 0 {
			xMinusANth = big.NewInt(1)
			xMinusANthPrecision = big.NewInt(0)
		} else if n.Cmp(bigOne) == 0 {
			xMinusANth = big.NewInt(0).Set(xMinusA)
			xMinusANth = big.NewInt(0).Set(xMinusAPrecision)
		} else {

			xMinusANth, xMinusANthPrecision, errX =
				BigIntMathMultiply{}.BigIntMultiply(
					xMinusANth,
					xMinusANthPrecision,
					xMinusA,
					xMinusAPrecision)

			if errX != nil {
				xNum = big.NewInt(0)
				xNumPrecision = big.NewInt(0)
				err = fmt.Errorf(ePrefix+"%v", errX.Error())
				return xNum, xNumPrecision, err
			}

			nFact, nFactPrecision, errX =
				BigIntMathMultiply{}.BigIntMultiply(
					nFact,
					nFactPrecision,
					n,
					big.NewInt(0))

		}

		factor1, factor1Precision, errX =
			BigIntMathDivide{}.BigIntFracQuotient(
				eToPwr,
				eToPwrPrecision,
				nFact,
				nFactPrecision,
				internalMaxPrecision)

		if errX != nil {
			xNum = big.NewInt(0)
			xNumPrecision = big.NewInt(0)
			err = fmt.Errorf(ePrefix+"%v", errX.Error())
			return xNum, xNumPrecision, err
		}

		factor2, factor2Precision, errX =
			BigIntMathMultiply{}.BigIntMultiply(
				factor1,
				factor1Precision,
				xMinusANth,
				xMinusANthPrecision)

		if errX != nil {
			xNum = big.NewInt(0)
			xNumPrecision = big.NewInt(0)
			err = fmt.Errorf(ePrefix+"%v", errX.Error())
			return xNum, xNumPrecision, err
		}

		xNum, xNumPrecision, errX = BigIntMathAdd{}.BigIntAdd(xNum, xNumPrecision, factor2, factor2Precision)

		if errX != nil {
			xNum = big.NewInt(0)
			xNumPrecision = big.NewInt(0)
			err = fmt.Errorf(ePrefix+"%v", errX.Error())
			return xNum, xNumPrecision, err
		}

	}

	xNum, xNumPrecision, errX = BigIntMath{}.RoundToMaxPrecision(xNum, xNumPrecision, maxPrecision, true)

	if errX != nil {
		xNum = big.NewInt(0)
		xNumPrecision = big.NewInt(0)
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return xNum, xNumPrecision, err
	}

	err = nil
	return xNum, xNumPrecision, err
}

func (bLog BigIntMathLogarithms) BigIntNumNatLogOfX(
	xNum BigIntNum,
	maxPrecision uint) (lnOfX BigIntNum, err error) {

	ePrefix := "BigIntMathLogarithms.BigIntNumNatLogOfX() "
	lnOfX = BigIntNum{}.NewZero(0)
	err = nil

	errX := xNum.IsValid(ePrefix)

	if errX != nil {
		err = fmt.Errorf("%v", errX.Error())
		return lnOfX, err
	}

	if xNum.IsZero() {
		err = errors.New(ePrefix + "Input Parameter 'xNum' is ZERO!")
		return lnOfX, err
	}

	biXNum := xNum.GetIntegerValue()
	biXNumPrecision := xNum.GetPrecisionBigInt()
	maxFinalPrecision := big.NewInt(0).SetUint64(uint64(maxPrecision))

	/*
			Assume maxPrecision = 50
		  m = 90
			agMeanMaxOutputPrecision 		= 190
			agMeanMaxInternalPrecision 	= 950
		  s4DivPrecisionMaxInternalPrecision = 2950
			piDivideMaxPrecision = 6950

	*/

	m, _, errX :=
		BigIntMath{}.RoundToMaxPrecision(
			big.NewInt(0).Mul(maxFinalPrecision, big.NewInt(18)),
			big.NewInt(1),
			big.NewInt(0),
			true)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", err.Error())
		return lnOfX, err
	}

	agMeanMaxOutputPrecision :=
		big.NewInt(0).Add(maxFinalPrecision, big.NewInt(100))

	agMeanMaxInternalPrecision :=
		big.NewInt(0).Mul(agMeanMaxOutputPrecision, big.NewInt(5))

	s4DivPrecisionMaxInternalPrecision :=
		big.NewInt(0).Add(agMeanMaxInternalPrecision, big.NewInt(2000))

	piDivideMaxPrecision :=
		big.NewInt(0).Add(s4DivPrecisionMaxInternalPrecision, big.NewInt(4000))

	biResult, biResultPrecision, errX :=
		BigIntMathLogarithms{}.SaskiKanadaNatLogOfX(
			biXNum,
			biXNumPrecision,
			m,
			s4DivPrecisionMaxInternalPrecision,
			agMeanMaxInternalPrecision,
			agMeanMaxOutputPrecision,
			piDivideMaxPrecision,
			maxFinalPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, err
	}

	lnOfX, errX = BigIntNum{}.NewBigIntPrecision(biResult, biResultPrecision)

	if errX != nil {
		lnOfX = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, err
	}

	err = nil

	return lnOfX, err
}

// https://en.wikipedia.org/wiki/Natural_logarithm
//
func (bLog BigIntMathLogarithms) SaskiKanadaNatLogOfX(
	xNum,
	xNumPrecision,
	m,
	s4DivPrecisionMaxInternalPrecision,
	agMeanMaxInternalPrecision,
	agMeanMaxOutputPrecision,
	piDivideMaxPrecision,
	maxFinalResultPrecision *big.Int) (lnOfX, lnOfXPrecision *big.Int, err error) {

	ePrefix := "BigIntMathLogarithms.SaskiKanadaNatLogOfX() "

	lnOfX = big.NewInt(0)
	lnOfXPrecision = big.NewInt(0)
	err = nil

	bigZero := big.NewInt(0)

	if xNumPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'xNumPrecision' is nil and INVALID!")
		return lnOfX, lnOfXPrecision, err
	}

	if xNum == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'xNum' is nil and INVALID!")
		return lnOfX, lnOfXPrecision, err
	}

	if m == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'm' is nil and INVALID!")
		return lnOfX, lnOfXPrecision, err
	}

	if piDivideMaxPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'piDivideMaxPrecision' is nil and INVALID!")
		return lnOfX, lnOfXPrecision, err
	}

	if maxFinalResultPrecision == nil {
		err = errors.New(ePrefix +
			"Error: Input parameter 'maxFinalResultPrecision' is nil and INVALID!")
		return lnOfX, lnOfXPrecision, err
	}

	if m.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'm' is negative and INVALID! "+
			"m='%v' ", m.Text(10))
		return lnOfX, lnOfXPrecision, err
	}

	if xNumPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'xNumPrecision' is negative and INVALID! "+
			"xNumPrecision='%v' ", xNumPrecision.Text(10))
		return lnOfX, lnOfXPrecision, err
	}

	if piDivideMaxPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'piDivideMaxPrecision' is negative and INVALID! "+
			"piDivideMaxPrecision='%v' ", piDivideMaxPrecision.Text(10))
		return lnOfX, lnOfXPrecision, err
	}

	if maxFinalResultPrecision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'maxFinalResultPrecision' is negative and INVALID! "+
			"maxFinalResultPrecision='%v' ", maxFinalResultPrecision.Text(10))
		return lnOfX, lnOfXPrecision, err
	}

	sFactor2Tom := big.NewInt(0).Exp(big.NewInt(2), m, nil)

	s := big.NewInt(0).Mul(xNum, sFactor2Tom)
	sPrecision := big.NewInt(0).Set(xNumPrecision)

	s, sPrecision, errX :=
		BigIntMath{}.RoundToMaxPrecision(
			s,
			sPrecision,
			big.NewInt(0).Add(s4DivPrecisionMaxInternalPrecision, big.NewInt(20)),
			true)

	fourDivS, fourDivSPrecision, errX :=
		BigIntMathDivide{}.BigIntFracQuotient(big.NewInt(4), big.NewInt(0), s, sPrecision, s4DivPrecisionMaxInternalPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, lnOfXPrecision, err
	}

	// Uses maxInternal Precision
	agMean, agMeanPrecision, _, _, _, errX :=
		BigIntMath{}.ArithmeticGeometricMean(
			big.NewInt(1),
			big.NewInt(0),
			fourDivS,
			fourDivSPrecision,
			agMeanMaxInternalPrecision,
			agMeanMaxOutputPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, lnOfXPrecision, err
	}

	denomFactorM := big.NewInt(0).Mul(big.NewInt(2), agMean)
	denomFactorMPrecision := big.NewInt(0).Set(agMeanPrecision)

	// Uses MaxInternal Precision
	factor1, factor1Precision, errX :=
		BigIntMathDivide{}.BigIntFracQuotient(
			piNumber20k.GetInteger(),
			piNumber20k.GetPrecisionBigInt(),
			denomFactorM,
			denomFactorMPrecision,
			piDivideMaxPrecision)

	if errX != nil {
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, lnOfXPrecision, err
	}

	factor2 := big.NewInt(0).Mul(m, natLogTwo20k.GetInteger())
	factor2Precision := big.NewInt(0).Set(natLogTwo20k.GetPrecisionBigInt())

	lnOfX, lnOfXPrecision, errX =
		BigIntMathSubtract{}.BigIntSubtract(
			factor1,
			factor1Precision,
			factor2,
			factor2Precision)

	if errX != nil {
		lnOfX = big.NewInt(0)
		lnOfXPrecision = big.NewInt(0)
		err = fmt.Errorf(ePrefix+"%v", errX.Error())
		return lnOfX, lnOfXPrecision, err
	}

	if lnOfXPrecision.Cmp(maxFinalResultPrecision) == 1 {
		lnOfX, lnOfXPrecision, errX =
			BigIntMath{}.RoundToMaxPrecision(
				lnOfX,
				lnOfXPrecision,
				maxFinalResultPrecision,
				true)

		if errX != nil {
			lnOfX = big.NewInt(0)
			lnOfXPrecision = big.NewInt(0)
			err = fmt.Errorf(ePrefix+"%v", errX.Error())
			return lnOfX, lnOfXPrecision, err
		}
	}

	err = nil

	return lnOfX, lnOfXPrecision, err
}

func (bLog BigIntMathLogarithms) codeDurationToStr(tDuration time.Duration) string {

	tMinutes := int64(0)
	tSeconds := int64(0)
	tMilliseconds := int64(0)
	tMicroseconds := int64(0)
	tNanoseconds := int64(0)

	i64TDur := int64(tDuration)
	outStr := ""
	totNanoSecs := int64(0)

	if i64TDur >= int64(time.Minute) {

		tMinutes = i64TDur / int64(time.Minute)
		outStr += fmt.Sprintf("%v-Minutes ", tMinutes)
		i64TDur -= tMinutes * int64(time.Minute)
		totNanoSecs = tMinutes * int64(time.Minute)
	}

	if i64TDur >= int64(time.Second) {
		tSeconds = i64TDur / int64(time.Second)
		outStr += fmt.Sprintf("%v-Seconds ", tSeconds)
		i64TDur -= tSeconds * int64(time.Second)
		totNanoSecs += tSeconds * int64(time.Second)
	}

	if i64TDur >= int64(time.Millisecond) {
		tMilliseconds = i64TDur / int64(time.Millisecond)
		i64TDur -= tMilliseconds * int64(time.Millisecond)
		totNanoSecs += tMilliseconds * int64(time.Millisecond)
	}

	if i64TDur >= int64(time.Microsecond) {
		tMicroseconds = i64TDur / int64(time.Microsecond)
		i64TDur -= tMicroseconds * int64(time.Microsecond)
		totNanoSecs += tMicroseconds * int64(time.Microsecond)
	}

	tNanoseconds = i64TDur
	totNanoSecs += tNanoseconds

	if totNanoSecs != int64(tDuration) {
		return fmt.Sprintf("Error: Total Calculated Duration= '%v'. Total Actual Duration= '%v'",
			totNanoSecs, int64(tDuration))
	}

	outStr += fmt.Sprintf("%v-Milliseconds ", tMilliseconds)

	outStr += fmt.Sprintf("%v-Microseconds ", tMicroseconds)

	outStr += fmt.Sprintf("%v-Nanoseconds ", tNanoseconds)

	return outStr
}
