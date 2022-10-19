package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

// BigIntMathNthRoot - Used to extract square roots and nth roots of positive and negative
// real numbers. Nth Roots may be either integer or fractional numbers. In addition, Nth
// Roots may be either positive or integer numbers.
//
// The technique employed to calculate nth roots is known as the
// "shifting nth-root algorithm".
//
// This source file is located in source code repository:
//
// https://github.com/MikeAustin71/mathopsgo.git
//
// See: https://en.wikipedia.org/wiki/Shifting_nth_root_algorithm
//
type BigIntMathNthRoot struct {
	NthRoot               BigIntNum
	OriginalRadicand      BigIntNum
	SetupRadicand         BigIntNum
	IntBundleRadicand     BigIntNum
	FracBundleRadicand    BigIntNum
	BundleAddOnPrecision  *big.Int
	FracBundleLength      *big.Int
	TotalBundleLength     *big.Int
	ResultBInt            *big.Int
	ActualResultPrecision *big.Int
	FracPrecision         *big.Int
	//ResultPrecision    int
	ResultBINum        BigIntNum
	RequestedPrecision uint
	BigOne             *big.Int
	Big10              *big.Int
	Big10ToNthPower    *big.Int
	BigZero            *big.Int
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

func (nthrt *BigIntMathNthRoot) Empty() {
	nthrt.NthRoot = BigIntNum{}.NewZero(0)
	nthrt.OriginalRadicand = BigIntNum{}.NewZero(0)
	nthrt.IntBundleRadicand = BigIntNum{}.NewZero(0)
	nthrt.FracBundleRadicand = BigIntNum{}.NewZero(0)
	nthrt.BundleAddOnPrecision = big.NewInt(0)
	nthrt.FracBundleLength = big.NewInt(0)
	nthrt.TotalBundleLength = big.NewInt(0)
	nthrt.ResultBInt = big.NewInt(0)
	nthrt.ResultBINum = BigIntNum{}.NewZero(0)
	nthrt.RequestedPrecision = 0
	nthrt.BigOne = big.NewInt(0)
	nthrt.Big10 = big.NewInt(0)
	nthrt.Big10ToNthPower = big.NewInt(0)
	nthrt.BigZero = big.NewInt(0)
	nthrt.YPrime = big.NewInt(0)
	nthrt.Minuend = big.NewInt(0)
	nthrt.Subtrahend = big.NewInt(0)
	nthrt.R = big.NewInt(0)
	nthrt.RPrime = big.NewInt(0)
	nthrt.BaseNum = big.NewInt(0)
	nthrt.Alpha = big.NewInt(0)
	nthrt.Beta = big.NewInt(0)
}

// OriginalNthRoot  - Calculates the Nth Root of a real number ('radicand')
// passed to the method as Type BigIntNum.  The calling function must supply
// input parameters for 'radicand', 'nthRoot' and 'maxPrecision'.
//
// Input Parameters
// ================
//
// radicand	- The radicand value from which the nth Root will be taken.
//            nthRootResult^nthRoot = radicand
//
// nthRoot  - Specifies the root which will be calculated for parameter,
// 						'radicand'. Example, square root, cube root, 4th root,
// 						9th root etc. 'nthRoot' is a BigIntNum Type which may be
// 						a positive or negative number. In addition, the nthRoot may
// 						be either an integer number or a fractional number.
// 					  The nthRoot must be with a numeric value greater than one ('1')
//            or less than minus one (-1). nthRoots with a value of zero will
// 						always return an nthRoot result of zero. Nth Root values of +1 or
//						-1 will generate an error.
//
//						If the radicand is negative and the nthRoot value is an even number
//						(evenly divisible by 2 with no remainder), an error will be returned
//						since the result of such a calculation is an imaginary number.
//
// maxPrecision	- Specifies the maximum number of decimals to the right of the
// 								decimal point to which the Nth root result will be calculated.
//
//	Returns
//	=======
//
//	BigIntNum - If successful the nth root result will be returned as a
//							BigIntNum type. This returned BigIntNum nth root will contain
//							numeric separators (decimal separator, thousands separator
//							and currency symbol) copied from input parameter,'radicand'.
//
//	error			- If the calculation completes successfully, the 'error' type
//							returned will be set equal to 'nil'. If an error is encountered
//							the returned 'error' type will contain an error message.
//
func (nthrt BigIntMathNthRoot) GetNthRoot(
	radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.OriginalNthRoot() "

	if radicand.GetSign() == -1 {

		isEvenNum, err := nthRoot.IsEvenNumber()

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by nthRoot.IsEvenNumber() "+
					"nthRoot='%v' Error='%v'\n", nthRoot.GetNumStr(), err.Error())
		}

		if isEvenNum {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"INVALID ENTRY - Cannot calculate nthRoot of a negative radicand when nthRoot is even. "+
					"Original Number= %v  nthRoot= %v\n", radicand.GetNumStr(), nthRoot.GetNumStr())
		}

	}

	// If the radicand is zero, the result will always be zero
	if radicand.IsZero() {
		return radicand, nil
	}

	numSeps := radicand.GetNumericSeparatorsDto()
	bigINumOne := BigIntNum{}.NewOne(0)

	var err error

	err = bigINumOne.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error returned by bigINumOne.SetNumericSeparatorsDto(numSeps).")
	}

	// If nthRoot is zero, the result will always be '1'
	if nthRoot.IsZero() {
		return bigINumOne, nil
	}

	// Error if nthRoot == 1
	if nthRoot.Cmp(bigINumOne) == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix +
				"Error - Input Parameter 'nthRoot' INVALID! 'nthRoot' cannot equal 1.\n")
	}

	var nthRootResult BigIntNum

	if nthRoot.GetSign() == -1 {

		nthRootResult, err = nthrt.calcNegativeNthRoot(radicand, nthRoot, maxPrecision)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+"Error returned by nthrt.calcNegativeNthRoot(...). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		nthRootResult, err = nthrt.calcPositiveNthRoot(radicand, nthRoot, maxPrecision)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+"Error returned by nthrt.calcPositiveNthRoot(...). "+
					"Error='%v' \n", err.Error())
		}

	}

	err = nthRootResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error returned by nthRootResult.SetNumericSeparatorsDto(numSeps).")
	}

	return nthRootResult, nil
}

//                             Stage 1                                             //
// ******************************************************************************* //

// calcPositiveNthRoot - Calculates the nth root of a radicand where the nth root
// is a positive value.
func (nthrt *BigIntMathNthRoot) calcPositiveNthRoot(radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.calcPositiveNthRoot() "

	if nthRoot.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error: This method only calculates nthRoot results for positive "+
				"nthRoot values. The entry for nthRoot is negative. nthRoot='%v'\n", nthRoot.GetNumStr())
	}

	if nthRoot.GetPrecisionUint() > 0 {
		return nthrt.calcPositiveFractionalNthRoot(radicand, nthRoot, maxPrecision)
	}

	return nthrt.calcPositiveIntegerNthRoot(radicand, nthRoot, maxPrecision)

}

// calcNegativeNthRoot - calculates the nth root result of a radicand where the
// nth root is a negative value.
//
func (nthrt *BigIntMathNthRoot) calcNegativeNthRoot(radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.calcNegativeNthRoot() "

	if nthRoot.GetSign() != -1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error: This method only calculates nthRoot results for positive "+
				"nthRoot values. The entry for nthRoot is negative. nthRoot='%v'\n", nthRoot.GetNumStr())
	}

	newNthRoot := nthRoot.CopyOut()

	newNthRoot.ChangeSign()

	var nthRootResult BigIntNum
	var err error

	if newNthRoot.GetPrecisionUint() == 0 {
		nthRootResult, err = nthrt.calcPositiveIntegerNthRoot(radicand, newNthRoot, maxPrecision)
		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+"Error returned by calcPositiveIntegerNthRoot(...) "+
					"Error='%v' \n", err.Error())
		}
	} else {
		nthRootResult, err = nthrt.calcPositiveFractionalNthRoot(radicand, newNthRoot, maxPrecision)
		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+"Error returned by calcPositiveFractionalNthRoot(...) "+
					"Error='%v' \n", err.Error())
		}
	}

	inverse, err := nthRootResult.Inverse(maxPrecision)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by nthRootResult.Inverse(maxPrecision) "+
				"maxPrecision='%v' Error='%v' \n", maxPrecision, err.Error())

	}

	return inverse, nil
}

//                             Stage 2                                             //
// ******************************************************************************* //

func (nthrt *BigIntMathNthRoot) calcPositiveIntegerNthRoot(radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.calcPositiveIntegerNthRoot() "

	if nthRoot.GetSign() != 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Expected 'nthRoot' to be positive. nthRoot is negative! "+
				"nthRoot='%v' ", nthRoot.GetNumStr())
	}

	if nthRoot.GetPrecision() != 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Expected 'nthRoot' to be integer value. nthRoot is fractional value! "+
				"nthRoot='%v' ", nthRoot.GetNumStr())
	}

	if radicand.GetSign() == -1 {

		isEvenNum, err := nthRoot.IsEvenNumber()

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by nthRoot.IsEvenNumber() "+
					"nthRoot='%v' Error='%v'\n", nthRoot.GetNumStr(), err.Error())
		}

		if isEvenNum {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even. "+
					"Original Number= %v  nthRoot= %v\n", radicand.GetNumStr(), nthRoot.GetNumStr())
		}

	}

	bINumResult, err := nthrt.calcNthRootGateway(radicand, nthRoot, maxPrecision)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.calcNthRootGateway(radicand, nthRoot, maxPrecision) "+
				"Error='%v' ", err.Error())
	}

	return bINumResult, nil
}

// calcPositiveFractionalNthRoot - Calculates the nth root result for positive fractional
// nthRoot values.
//
func (nthrt *BigIntMathNthRoot) calcPositiveFractionalNthRoot(radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.calcPositiveFractionalNthRoot() "

	if nthRoot.GetSign() < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: This method only processes Positive NthRootInt Values. "+
				"The current NthRootInt is a NEGATIVE Value. nthRoot='%v'",
				nthRoot.GetNumStr())
	}

	if nthRoot.GetPrecisionUint() == 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: This method only processes Fractional NthRoots. "+
				"The current NthRootInt is an Integer Value. nthRoot='%v'",
				nthRoot.GetNumStr())
	}

	modXZero := big.NewInt(0)

	scaleFactor := big.NewInt(0).Exp(
		big.NewInt(10),
		big.NewInt(int64(nthRoot.GetPrecisionUint())),
		modXZero)

	ratFraction := big.NewRat(1, 1).SetFrac(nthRoot.GetAbsoluteBigIntValue(), scaleFactor)

	// Numerator of ratFraction is new nthRoot
	newNthRoot := BigIntNum{}.NewBigInt(ratFraction.Num(), 0)

	// If nthRoot is zero, the result will always be '1'
	if newNthRoot.IsZero() {
		return BigIntNum{}.NewOne(maxPrecision), nil
	}

	bigINumOne := BigIntNum{}.NewOne(0)

	// if nthRoot == 1
	if newNthRoot.Cmp(bigINumOne) == 0 {
		return radicand.CopyOut(), nil
	}

	// Denominator of ratFraction is exponent for current radicand.
	exponent := BigIntNum{}.NewBigInt(ratFraction.Denom(), 0)

	radicandPrecision := BigIntNum{}.NewBigInt(big.NewInt(int64(radicand.GetPrecisionUint())), 0)

	newMaxPrecision := BigIntMathMultiply{}.MultiplyBigIntNums(exponent, radicandPrecision)

	exponentMaxPrecision, _ := newMaxPrecision.GetUInt()

	newRadicand, _ := BigIntMathPower{}.Pwr(radicand, exponent, exponentMaxPrecision)

	if newRadicand.GetSign() == -1 {

		isEvenNum, err := newNthRoot.IsEvenNumber()

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by newNthRoot.IsEvenNumber() "+
					"newNthRoot='%v' Error='%v'\n", newNthRoot.GetNumStr(), err.Error())
		}

		if isEvenNum {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"INVALID ENTRY - Cannot calculate nthRoot of a negative number when nthRoot is even. "+
					"Original Number= %v  newNthRoot= %v\n", radicand.GetNumStr(), newNthRoot.GetNumStr())
		}

	}

	nthRootResult, err := nthrt.calcNthRootGateway(newRadicand, newNthRoot, maxPrecision)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by nthrt.calcNthRootGateway(newRadicand, newNthRoot, maxPrecision). "+
				"Error='%v' ", err.Error())
	}

	return nthRootResult, nil
}

// calcNthRootGateway - Calculate the nth root of a radicand. The result is returned
// as a BigIntNum.
//
// Input Parameters
// ================
//
// radicand - BigIntNum 	A 'radicand' is the number under a radical symbol (√)
//
// nthRoot 	- BigIntNum 	n√ The index or root. This must be a positive integer value.
//
//maxPrecision	- uint		The maximum precision specified for the result.
//
//
func (nthrt *BigIntMathNthRoot) calcNthRootGateway(
	radicand, nthRoot BigIntNum,
	maxPrecision uint) (BigIntNum, error) {

	ePrefix := "BigIntMathNthRoot.calcNthRootGateway() "

	if nthRoot.GetSign() < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error: Cannot process a negative Nth Root Value. "+
				"Input parameter, nthRoot='%v' ",
				nthRoot.GetNumStr())
	}

	_, err := nthrt.initialize(radicand, nthRoot, maxPrecision)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"-Error returned from initialization. Error= %v", err.Error())
	}

	err = nthrt.doRootExtraction()

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"-Error returned from nthrt.doRootExtraction() - %v", err.Error())
	}

	// fmt.Println("FracPrecision Count: ", nthrt.FracPrecision.Text(10))

	return nthrt.ResultBINum, nil
}

// initialize - Initializes the data fields of the BigIntMathNthRoot structure and validates the
// original number passed to the Nth Root Calculation. This method assumes that radicand and
// nthRoot have already been validated.
//
func (nthrt *BigIntMathNthRoot) initialize(
	radicand, nthRoot BigIntNum,
	maxPrecision uint) (setupRadicand BigIntNum, err error) {

	ePrefix := "BigIntMathNthRoot.initialize() "
	var errx error
	precisionAdjustment := big.NewInt(0)
	nthrt.OriginalRadicand = radicand.CopyOut()
	nthrt.NthRoot = nthRoot.CopyOut()
	nthrt.RequestedPrecision = maxPrecision

	nthrt.SetupRadicand, nthrt.IntBundleRadicand, nthrt.FracBundleRadicand, precisionAdjustment, errx =
		nthrt.setupBundles(nthrt.OriginalRadicand,
			nthrt.NthRoot)
	if errx != nil {
		setupRadicand = BigIntNum{}.NewZero(0)
		err = fmt.Errorf(ePrefix+
			"Error returned from nthrt.setupBundles(OriginalRadicand, NthRootInt). "+
			"OriginalRadicand='%v' NthRootInt='%v' Error= %v",
			nthrt.OriginalRadicand.GetNumStr(), nthrt.NthRoot.GetNumStr(), errx.Error())
		return setupRadicand, err
	}

	nthrt.BundleAddOnPrecision =
		big.NewInt(int64(maxPrecision + uint(2)))

	nthrt.FracBundleLength = big.NewInt(0)
	nthrt.TotalBundleLength = big.NewInt(0)

	nthrt.TotalBundleLength, nthrt.FracBundleLength, errx = nthrt.calcBundleLength(
		nthrt.SetupRadicand,
		nthrt.NthRoot,
		nthrt.BundleAddOnPrecision)

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned from nthrt.calcBundleLength(SetupRadicand, NthRootInt, AddOnPrecision). "+
			"SetupRadicand='%v' NthRootInt='%v' AddOnPrecision='%v' Error= %v",
			nthrt.SetupRadicand.GetNumStr(), nthrt.NthRoot.GetNumStr(),
			nthrt.BundleAddOnPrecision.Text(10), errx.Error())
		setupRadicand = BigIntNum{}.NewZero(0)
		return setupRadicand, err
	}

	radicandPrecision := big.NewInt(int64(nthrt.SetupRadicand.GetPrecisionUint()))

	nthrt.ActualResultPrecision, errx =
		nthrt.calcPrecision(radicandPrecision,
			nthrt.FracBundleLength,
			precisionAdjustment,
			nthrt.BundleAddOnPrecision,
			nthrt.NthRoot.GetAbsoluteBigIntValue())

	if errx != nil {
		err = fmt.Errorf(ePrefix+"Error returned from nthrt.calcPrecision(). Error= %v", errx)
		setupRadicand = BigIntNum{}.NewZero(0)
		return setupRadicand, err
	}

	// Set constants for calculations
	nthrt.FracPrecision = big.NewInt(0)
	nthrt.BigOne = big.NewInt(1)
	nthrt.Big10 = big.NewInt(10)
	nthrt.Big10ToNthPower = big.NewInt(0).Exp(
		nthrt.Big10,
		nthrt.NthRoot.GetAbsoluteBigIntValue(),
		nil)
	nthrt.BigZero = big.NewInt(0)
	nthrt.ResultBInt = big.NewInt(0)
	nthrt.ResultBINum = BigIntNum{}.NewZero(0)
	nthrt.Y = big.NewInt(0)
	nthrt.YPrime = big.NewInt(0)
	nthrt.R = big.NewInt(0)
	nthrt.RPrime = big.NewInt(0)
	nthrt.BaseNum = big.NewInt(10)
	nthrt.Alpha = big.NewInt(0)
	nthrt.Beta = big.NewInt(0)
	return nthrt.SetupRadicand, nil

}

// Experimental 2
func (nthrt *BigIntMathNthRoot) setupBundles(
	radicand, nthRoot BigIntNum) (setupRadicand,
	intBundleRadicand,
	fracBundleRadicand BigIntNum,
	precisionAdjustment *big.Int,
	err error) {

	ePrefix := "BigIntMathNthRoot.setupBundles() "

	intBundleRadicand = BigIntNum{}.NewZero(0)
	fracBundleRadicand = BigIntNum{}.NewZero(0)
	precisionAdjustment = big.NewInt(0)
	err = nil

	var errx error

	modX := big.NewInt(0)
	bigTen := big.NewInt(10)
	scaleValue := big.NewInt(0)

	setupRadicand =
		BigIntNum{}.NewBigInt(
			radicand.GetAbsoluteBigIntValue(),
			radicand.GetPrecisionUint())

	setupRadicand.TrimTrailingFracZeros()
	setupRadicand.SetExpectedToActualNumberOfDigits()

	if setupRadicand.GetPrecisionUint() == 0 {
		intBundleRadicand = setupRadicand.CopyOut()

		/*
		 fmt.Println("setupRadicand: ", setupRadicand.GetNumStr())
		 fmt.Println("intBundleRadicand:", intBundleRadicand.GetNumStr())
		 fmt.Println("fracBundleRadicand:", fracBundleRadicand.GetNumStr())
		*/

		return setupRadicand, intBundleRadicand, fracBundleRadicand, precisionAdjustment, err
	}

	setupRadicandTotalDigits, _, errx := setupRadicand.GetActualNumberOfDigits()

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by setupRadicand.GetActualNumberOfDigits(). "+
			"Error='%v' ", errx.Error())

		return BigIntNum{}.NewZero(0), intBundleRadicand, fracBundleRadicand, precisionAdjustment, err
	}

	// Precision must be greater than zero
	radicandPrecision := setupRadicand.GetPrecisionBigInt()
	setupRadicandIntegerDigits :=
		big.NewInt(0).Sub(setupRadicandTotalDigits, radicandPrecision)

	expectedFractionalDigits := big.NewInt(0).Set(radicandPrecision)

	scaleValue = big.NewInt(0).Exp(bigTen, radicandPrecision, nil)

	intBundleRadicandBigInt, fracBundleRadicandBigInt :=
		big.NewInt(0).QuoRem(setupRadicand.GetAbsoluteBigIntValue(), scaleValue, modX)

	intBundleRadicand = BigIntNum{}.NewBigInt(intBundleRadicandBigInt, 0)
	intBundleRadicand.SetExpectedNumberOfDigits(setupRadicandIntegerDigits)

	mod := big.NewInt(0).Rem(radicandPrecision, nthRoot.GetAbsoluteBigIntValue())

	if mod.Cmp(big.NewInt(0)) == 1 {
		delta := big.NewInt(0).Sub(nthRoot.GetAbsoluteBigIntValue(), mod)
		scaleVal := big.NewInt(0).Exp(big.NewInt(10), delta, nil)
		fracBundleRadicandBigInt = big.NewInt(0).Mul(fracBundleRadicandBigInt, scaleVal)
		expectedFractionalDigits = big.NewInt(0).Add(expectedFractionalDigits, delta)
	}

	fracBundleRadicand = BigIntNum{}.NewBigInt(fracBundleRadicandBigInt, 0)

	if intBundleRadicand.IsZero() &&
		fracBundleRadicand.IsZero() {
		err = errors.New(ePrefix + "Error: Both intBundleRadicand and fracBundleRadicand are ZERO!")
		return BigIntNum{}.NewZero(0), intBundleRadicand, fracBundleRadicand, precisionAdjustment, err
	}

	fracBundleRadicand.SetExpectedNumberOfDigits(expectedFractionalDigits)

	err = nil

	return setupRadicand, intBundleRadicand, fracBundleRadicand, precisionAdjustment, err
}

// calcBundleLength - Calculates the final bundle length. These are bundles
// of integers that are packaged for submission to the Nth Root calculation
// routine.
//
// The final bundle length is equal to the actual number of bundles associated
// with the radicand plus bundles added on for additional precision in the
// final nthRoot result.
//
func (nthrt *BigIntMathNthRoot) calcBundleLength(
	radicand,
	nthRoot BigIntNum,
	bundleAddonPrecision *big.Int) (totalBundleLength, fracBundleLength *big.Int, err error) {

	ePrefix := "BigIntMathNthRoot.calcBundleLength() "
	totalBundleLength = big.NewInt(0)
	fracBundleLength = big.NewInt(0)

	err = nil
	var errx error

	modX := big.NewInt(0)
	bigZero := big.NewInt(0)
	bigOne := big.NewInt(1)
	numOfDigits := big.NewInt(1)
	intBundleLength := big.NewInt(0)
	intBundleLengthMod := big.NewInt(0)
	bigIntNthRoot := nthRoot.GetAbsoluteBigIntValue()

	scaleValue := big.NewInt(0).Exp(
		big.NewInt(10),
		big.NewInt(int64(radicand.GetPrecisionUint())),
		nil)

	radicandBigInt := radicand.GetAbsoluteBigIntValue()

	// break radicand into integer digits
	intRadicand :=
		big.NewInt(0).Quo(
			radicandBigInt,
			scaleValue)

	if intRadicand.Cmp(bigZero) == 1 {

		numOfDigits, errx = BigIntMath{}.GetMagnitude(intRadicand)

		if errx != nil {
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMath{}.GetMagnitudeDigits(intRadicand). "+
				"intRadicand='%v' Error='%v'",
				intRadicand.Text(10), errx.Error())

			return totalBundleLength, fracBundleLength, err
		}

		// Convert integer radicand to number of digits
		numOfDigits = big.NewInt(0).Add(numOfDigits, bigOne)

		intBundleLength, intBundleLengthMod = big.NewInt(0).QuoRem(
			numOfDigits,
			bigIntNthRoot,
			modX)

		// For the integer bundle calculation add one to bundle length if
		// there are any remaining digits.
		// Example  5 / 3 = Quotient of 1 Mod of 2. Add 1 to bundle size
		if intBundleLengthMod.Cmp(bigZero) == 1 {
			intBundleLength = big.NewInt(0).Add(intBundleLength, bigOne)
		}

	}

	radicandPrecision := big.NewInt(int64(radicand.GetPrecisionUint()))

	if radicandPrecision.Cmp(bigZero) == 1 {
		fracBundleLength = big.NewInt(0).Quo(radicandPrecision, bigIntNthRoot)
	}

	/*
		fmt.Println("        intBundleLength: ", intBundleLength.Text(10))
		fmt.Println("       fracBundleLength: ", fracBundleLength.Text(10))
		fmt.Println( "  bundleAddonPrecision: ", bundleAddonPrecision.Text(10))
	*/

	totalBundleLength = big.NewInt(0).Add(intBundleLength, fracBundleLength)

	totalBundleLength = big.NewInt(0).Add(totalBundleLength, bundleAddonPrecision)

	err = nil

	return totalBundleLength, fracBundleLength, err
}

// calcPrecision - This method will calculate the actual precision of the
// result output by the Nth Root calculation.  Actual result precision is
// equal to the Bundle Add On Precision + the quotient of radicand precision
// divided by nthRoot.
//			bundleAddOnPrecision + Quotient(radicandPrecision/nthRoot)
//
func (nthrt *BigIntMathNthRoot) calcPrecision(
	radicandPrecision,
	fracBundleLength,
	precisionAdjustment,
	bundleAddOnPrecision,
	nthRoot *big.Int) (actualPrecision *big.Int, err error) {

	actualPrecision = big.NewInt(0)
	err = nil
	bigZero := big.NewInt(0)

	if nthRoot.Cmp(bigZero) == 0 {
		ePrefix := "BigIntMathNthRoot.calcPrecision() "
		err = errors.New(ePrefix +
			"Error: nthRoot=0. Division by nthRoot is Division by zero and will FAIL!")
		return actualPrecision, err

	}

	actualPrecision = big.NewInt(0).Add(bundleAddOnPrecision, fracBundleLength)
	actualPrecision = big.NewInt(0).Add(actualPrecision, precisionAdjustment)

	if actualPrecision.Cmp(bigZero) < 0 {
		actualPrecision = big.NewInt(0)
	}

	/*
		fmt.Println("    fracBundleLength: ", fracBundleLength.Text(10))
		fmt.Println("bundleAddOnPrecision: ", bundleAddOnPrecision.Text(10))
		fmt.Println("precision Adjustment: ", precisionAdjustment.Text(10))
		fmt.Println("    actual precision: ", actualPrecision.Text(10))
	*/
	err = nil

	return actualPrecision, err
}

// doRootExtraction - Manages the root extraction process
//
func (nthrt *BigIntMathNthRoot) doRootExtraction() error {

	ePrefix := "BigIntMathNthRoot.doRootExtraction() "
	nthrt.Y = big.NewInt(0)
	nthrt.Minuend = big.NewInt(0)
	nthrt.Subtrahend = big.NewInt(0)
	nthrt.R = big.NewInt(0)
	bigOne := big.NewInt(1)

	for i := big.NewInt(0); i.Cmp(nthrt.TotalBundleLength) == -1; i = big.NewInt(0).Add(i, bigOne) {

		err := nthrt.findNextRoot()

		if err != nil {
			return fmt.Errorf(ePrefix+
				"Error returned by nthrt.findNextRoot(i). "+
				"Error='%v' ", err.Error())
		}

	}

	if nthrt.OriginalRadicand.GetSign() < 1 {
		nthrt.ResultBInt = big.NewInt(0).Neg(nthrt.ResultBInt)
	}

	actualPrecision := uint(nthrt.ActualResultPrecision.Int64())

	nthrt.ResultBINum = BigIntNum{}.NewBigInt(nthrt.ResultBInt, actualPrecision)

	nthrt.ResultBINum.RoundToDecPlace(nthrt.RequestedPrecision)

	return nil
}

// findNextRoot - Called by doRootExtraction() to find the next
// root value.
//
func (nthrt *BigIntMathNthRoot) findNextRoot() error {

	ePrefix := "BigIntMathNthRoot.findNextRoot() "
	var err error
	var bundle *big.Int

	bundle, nthrt.IntBundleRadicand, nthrt.FracBundleRadicand, err =
		nthrt.getNextBundleBigIntValue(nthrt.IntBundleRadicand,
			nthrt.FracBundleRadicand,
			nthrt.NthRoot)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by nthrt.getNextBundleBigIntValue(nthrt.IntBundleRadicand). "+
			"Error='%v' ", err.Error())
	}

	// alpha = next n-digits of radicand
	nthrt.Alpha = big.NewInt(0).Set(bundle)

	nthrt.RPrime = big.NewInt(-1)

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

		term_2a = big.NewInt(0).Exp(term_2a2,
			big.NewInt(0).Set(nthrt.NthRoot.GetAbsoluteBigIntValue()),
			nil)

		term_2b1 = big.NewInt(0).Set(nthrt.Big10ToNthPower)

		term_2b2 = big.NewInt(0).Exp(nthrt.Y,
			big.NewInt(0).Set(nthrt.NthRoot.GetAbsoluteBigIntValue()),
			nil)

		term_2b = big.NewInt(0).Mul(term_2b1, term_2b2)

		nthrt.Subtrahend = big.NewInt(0).Sub(term_2a, term_2b)

		nthrt.RPrime = big.NewInt(0).Sub(nthrt.Minuend, nthrt.Subtrahend)

		itatr = big.NewInt(0).Sub(itatr, nthrt.BigOne)
	}

	nthrt.R = big.NewInt(0).Set(nthrt.RPrime)
	nthrt.Y = big.NewInt(0).Set(nthrt.YPrime)
	nthrt.ResultBInt = big.NewInt(0).Mul(nthrt.ResultBInt, nthrt.Big10)
	nthrt.ResultBInt = big.NewInt(0).Add(nthrt.ResultBInt, nthrt.Beta)

	return nil
}

// getNextBundleBigIntValue - Calculate and return the next bundle of
// numeric digits for nth Root extraction calculations.
//

func (nthrt *BigIntMathNthRoot) getNextBundleBigIntValue(

	intBundleRadicand, fracBundleRadicand, nthRoot BigIntNum) (
	nextBundleValue *big.Int, newIntBundleRadicand, newFracBundleRadicand BigIntNum, err error) {

	ePrefix := "BigIntMathNthRoot.getNextBundleBigIntValue() "
	nextBundleValue = big.NewInt(0)
	newIntBundleRadicand = BigIntNum{}.NewZero(0)
	newIntBundleRadicand.SetExpectedToActualNumberOfDigits()
	newFracBundleRadicand = BigIntNum{}.NewZero(0)
	newFracBundleRadicand.SetExpectedToActualNumberOfDigits()
	err = nil
	modX := big.NewInt(0)

	var errx error
	var magnitude *big.Int
	var numberOfDigits *big.Int
	var tempRadicand *big.Int
	var digitsQuotient *big.Int
	var exponent *big.Int
	var divisor *big.Int

	intExpectedNumOfDigits := intBundleRadicand.GetExpectedNumberOfDigits()
	intActualNumOfDigits, _, errx := intBundleRadicand.GetActualNumberOfDigits()
	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by intBundleRadicand.GetActualNumberOfDigits(). "+
			"Error='%v' ", errx.Error())

		return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
	}

	intExpectedActualDelta := big.NewInt(0).Sub(intExpectedNumOfDigits, intActualNumOfDigits)
	intModExpectedNthRoot := big.NewInt(0).Rem(intExpectedNumOfDigits, nthRoot.GetAbsoluteBigIntValue())
	intIsZeroValue := intBundleRadicand.IsZero()

	bigZero := big.NewInt(0)

	fracExpectedNumOfDigits := fracBundleRadicand.GetExpectedNumberOfDigits()
	fracActualNumOfDigits, _, errx := fracBundleRadicand.GetActualNumberOfDigits()

	if errx != nil {
		err = fmt.Errorf(ePrefix+
			"Error returned by fracBundleRadicand.GetActualNumberOfDigits(). "+
			"Error='%v' ", errx.Error())

		return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
	}

	fracExpectedActualDelta := big.NewInt(0).Sub(fracExpectedNumOfDigits, fracActualNumOfDigits)
	fracIsZeroValue := fracBundleRadicand.IsZero()

	if (intExpectedActualDelta.Cmp(bigZero) > 0 &&
		intExpectedActualDelta.Cmp(nthRoot.GetAbsoluteBigIntValue()) >= 0) ||
		(intExpectedActualDelta.Cmp(bigZero) > 0 &&
			intIsZeroValue) {

		// fmt.Println("Calc nthRoot IntegerDeltaDigits >=0")

		nextBundleValue = big.NewInt(0)
		newIntBundleRadicand = intBundleRadicand.CopyOut()
		newFracBundleRadicand = fracBundleRadicand.CopyOut()
		intExpectedNumOfDigits = big.NewInt(0).Sub(intExpectedNumOfDigits, nthRoot.GetAbsoluteBigIntValue())

		if intExpectedNumOfDigits.Cmp(intActualNumOfDigits) == -1 {
			intExpectedNumOfDigits = big.NewInt(0).Set(intActualNumOfDigits)
		}

		newIntBundleRadicand.SetExpectedNumberOfDigits(intExpectedNumOfDigits)

		err = nil

	} else if !intBundleRadicand.IsZero() {

		// fmt.Println("Calc inBundleRadicand is NOT zero")

		magnitude, errx = BigIntMath{}.GetMagnitude(intBundleRadicand.GetAbsoluteBigIntValue())

		if errx != nil {
			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMath{}.GetMagnitudeDigits(intBundleRadicand.GetAbsoluteBigIntValue()). "+
				"intBundleRadicand='%v' Error='%v'",
				intBundleRadicand.GetNumStr(), errx.Error())

			nextBundleValue = big.NewInt(0)
			newIntBundleRadicand = BigIntNum{}.NewZero(0)
			newFracBundleRadicand = BigIntNum{}.NewZero(0)

			return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
		}

		intActualNumOfDigits = big.NewInt(0).Add(magnitude, big.NewInt(1))
		digitsQuotient = big.NewInt(0).Quo(magnitude, nthrt.NthRoot.GetAbsoluteBigIntValue())
		exponent = big.NewInt(0).Mul(digitsQuotient, nthrt.NthRoot.GetAbsoluteBigIntValue())
		divisor = big.NewInt(0).Exp(nthrt.Big10, exponent, nil)

		nextBundleValue, tempRadicand =
			big.NewInt(0).QuoRem(intBundleRadicand.GetAbsoluteBigIntValue(), divisor, modX)

		numberOfDigits, errx = BigIntMath{}.GetMagnitude(nextBundleValue)

		if errx != nil {
			err = fmt.Errorf(ePrefix+"Error returned by BigIntMath{}.GetMagnitudeDigits(nextBundleValue). "+
				"nextBundleValue='%v' Error='%v'",
				nextBundleValue.Text(10), errx.Error())

			nextBundleValue = big.NewInt(0)
			newIntBundleRadicand = BigIntNum{}.NewZero(0)
			newFracBundleRadicand = BigIntNum{}.NewZero(0)

			return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
		}

		numberOfDigits = big.NewInt(0).Add(numberOfDigits, big.NewInt(1))

		if intModExpectedNthRoot.Cmp(bigZero) > 0 {
			intExpectedNumOfDigits = big.NewInt(0).Sub(intExpectedNumOfDigits, numberOfDigits)
		} else {
			intExpectedNumOfDigits = big.NewInt(0).Sub(intExpectedNumOfDigits, nthrt.NthRoot.GetAbsoluteBigIntValue())
		}

		newIntBundleRadicand = BigIntNum{}.NewBigInt(tempRadicand, 0)
		newIntBundleRadicand.SetExpectedNumberOfDigits(intExpectedNumOfDigits)
		newFracBundleRadicand = fracBundleRadicand.CopyOut()
		err = nil

	} else if (fracExpectedActualDelta.Cmp(bigZero) > 0 &&
		fracExpectedActualDelta.Cmp(nthRoot.GetAbsoluteBigIntValue()) >= 0) ||
		(fracExpectedActualDelta.Cmp(bigZero) > 0 && fracIsZeroValue) {

		// fmt.Println("Calc nthRoot FracDeltaDigits >= 0")

		nextBundleValue = big.NewInt(0)
		newIntBundleRadicand = intBundleRadicand.CopyOut()
		newFracBundleRadicand = fracBundleRadicand.CopyOut()
		fracExpectedNumOfDigits = big.NewInt(0).Sub(fracExpectedNumOfDigits, nthRoot.GetAbsoluteBigIntValue())

		if fracExpectedNumOfDigits.Cmp(fracActualNumOfDigits) == -1 {
			fracExpectedNumOfDigits = big.NewInt(0).Set(fracActualNumOfDigits)
		}

		newFracBundleRadicand.SetExpectedNumberOfDigits(fracExpectedNumOfDigits)
		nthrt.FracPrecision = big.NewInt(0).Add(nthrt.FracPrecision, big.NewInt(1))
		err = nil

	} else if !fracBundleRadicand.IsZero() {

		// fmt.Println("Calc fracBundleRadicand Is NOT zero")

		magnitude, errx = BigIntMath{}.GetMagnitude(fracBundleRadicand.GetAbsoluteBigIntValue())

		if errx != nil {

			err = fmt.Errorf(ePrefix+
				"Error returned by BigIntMath{}.GetMagnitudeDigits(fracBundleRadicand.GetAbsoluteBigIntValue()). "+
				"fracBundleRadicand='%v' Error='%v' ",
				fracBundleRadicand.GetNumStr(), errx.Error())

			nextBundleValue = big.NewInt(0)
			newIntBundleRadicand = BigIntNum{}.NewZero(0)
			newFracBundleRadicand = BigIntNum{}.NewZero(0)
			return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
		}

		digitsQuotient = big.NewInt(0).Quo(magnitude, nthrt.NthRoot.GetAbsoluteBigIntValue())
		exponent = big.NewInt(0).Mul(digitsQuotient, nthrt.NthRoot.GetAbsoluteBigIntValue())
		divisor = big.NewInt(0).Exp(nthrt.Big10, exponent, nil)

		nextBundleValue, tempRadicand =
			big.NewInt(0).QuoRem(fracBundleRadicand.GetAbsoluteBigIntValue(), divisor, modX)

		newIntBundleRadicand = BigIntNum{}.NewZero(0)
		newIntBundleRadicand.SetExpectedToActualNumberOfDigits()
		fracExpectedNumOfDigits = big.NewInt(0).Sub(fracExpectedNumOfDigits, nthrt.NthRoot.GetAbsoluteBigIntValue())
		newFracBundleRadicand = BigIntNum{}.NewBigInt(tempRadicand, 0)
		newFracBundleRadicand.SetExpectedNumberOfDigits(fracExpectedNumOfDigits)
		nthrt.FracPrecision = big.NewInt(0).Add(nthrt.FracPrecision, big.NewInt(1))
		err = nil

	} else {

		// fmt.Println("Calc else zero")

		nextBundleValue = big.NewInt(0)
		newIntBundleRadicand = BigIntNum{}.NewZero(0)
		newIntBundleRadicand.SetExpectedToActualNumberOfDigits()
		newFracBundleRadicand = BigIntNum{}.NewZero(0)
		newFracBundleRadicand.SetExpectedToActualNumberOfDigits()
		nthrt.FracPrecision = big.NewInt(0).Add(nthrt.FracPrecision, big.NewInt(1))
		err = nil
	}

	return nextBundleValue, newIntBundleRadicand, newFracBundleRadicand, err
}
