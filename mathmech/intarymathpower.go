package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

type IntAryMathPower struct {
	Input  IntAryPair
	Result IntAry
}

// MinimumRequiredPrecision - designed to be used with the power function
// below. This method will compute the minimum number of decimal places
// required to support the result of raising a 'base' value to a specified
// 'exponent'. Both the 'base' and the 'exponent' are passed to this function
// as type *IntAry.
//
// For example, raising the value 3.12 to the power of 4 means that the
// result will require at least 8-decimal places to the right of the
// decimal in order to display a correct result. In the following example
// with base ='3.12' and exponent = '4', this method will return '8'.
//
// 		Example: 3.12^4 = 94.75854336 (8-digits to the right of the decimal)
//
// The calculated minimum required precision is returned as a positive
// value of type 'int'.
//
// If the minimum required precision exceeds the maximum positive value for
// IntAry precision ( +2,147,483,646, which equals 2^31 − 2), an error
// message is returned in addition to the maximum positive value for IntAry
// precision (+2,147,483,646).
//
func (iaPwr IntAryMathPower) MinimumRequiredPrecision(
	base, exponent *IntAry) (int, error) {

	ePrefix := "IntAryMathPower.MinimumRequiredPrecision() "

	maxValue := 2147483646

	basePrecision := IntAry{}.NewUint(base.GetPrecisionUint(), 0)

	tExponent := exponent.CopyOut()

	if tExponent.GetSign() == -1 {
		tExponent.ChangeSign()
	}

	iaResult := IntAry{}

	maxPrecision := tExponent.GetPrecision() + 5

	err := IntAryMathMultiply{}.Multiply(&basePrecision, &tExponent, &iaResult, 0, maxPrecision)

	if err != nil {
		return maxValue, err
	}

	if iaResult.GetPrecisionUint() > 0 {

		iaResult.RoundToPrecision(0)

	}

	intVal, err := iaResult.GetInt()

	if err != nil {

		return maxValue,
			fmt.Errorf(ePrefix+
				"Error: Minimum Required Precision exceeded maximum value of %v",
				maxValue)
	}

	return intVal, nil
}

// Pwr - Raises input parameter 'base' to the power of 'exponent'.
// This method uses the 'Power By Twos' technique.
// See:
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring#Computation_by_powers_of_2
// This method is based on revised code taken in part from Ye Lin Aung.
// https://stackoverflow.com/questions/30182129/calculating-large-exponentiation-in-golang
// This algorithm modified by Mike Rapp to achieve improved performance.
// The result of raising 'base' to the power of 'exponent' will return
// the result in 'base'. As such the original value of 'base' will be
// overwritten. The returned value 'base' will contain the same numeric
// separators (decimal separator, thousands separator and currency symbol)
// as that of the original 'base instance.'base' numeric separators will
// therefore remain unchanged.
//
//										base = base^exponent
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
func (iaPwr IntAryMathPower) Pwr(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) error {

	ePrefix := "IntAryMathPower.Pwr() "

	err := base.IsValid(ePrefix + "- base IntAry Error: ")

	if err != nil {
		return err
	}

	err = exponent.IsValid(ePrefix + "- exponent IntAry Error: ")

	if err != nil {
		return err
	}

	if base.IsZero() {
		return errors.New(ePrefix + "Error: Input parameter 'base' is zero value. INVALID INPUT!")
	}

	if exponent.IsZero() {
		base.SetIntAryToOne(minResultPrecision)
		return nil
	}

	iaOne := IntAry{}.NewOne(exponent.GetPrecision())

	if exponent.Equals(&iaOne) {
		return nil
	}

	exponentPrecision := exponent.GetPrecision()
	exponentSign := exponent.GetSign()

	if exponentPrecision == 0 && exponentSign == 1 {
		return iaPwr.pwrTwoPositiveIntegerExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision == 0 && exponentSign == -1 {
		return iaPwr.pwrTwoNegativeIntegerExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision > 0 && exponentSign == 1 {
		return iaPwr.pwrTwoPositiveFractionalExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision > 0 && exponentSign == -1 {
		return iaPwr.pwrTwoNegativeFractionalExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	return errors.New(ePrefix + "Error: input parameters failed to match valid calculation types!")
}

// PwrByMultiplication - raises base to the power of exponent using
// repetitive multiplication. This method may be slower than the method
// IntAryMathPower.Pwr(); however, this method is capable of handling
// very large exponents. Effectively, input parameter 'base' is raised
// to the power of exponent.
//
//								result = base^exponent
//
// The result of this operation is returned as pointer to an IntAry
// instance. This returned IntAry instance will contain numeric
// separators (decimal separator, thousands separator and currency
// symbol) copied from input parameter 'base'.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
func (iaPwr IntAryMathPower) PwrByMultiplication(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) (*IntAry, error) {

	ePrefix := "IntAryMathPower.PwrByMultiplication() "
	iaReturn := IntAry{}.NewZero(0)

	iaReturn.SetNumericSeparatorsDto(base.GetNumericSeparatorsDto())

	err := base.IsValid(ePrefix + "Invalid 'base' - ")

	if err != nil {
		return &iaReturn, err
	}

	err = exponent.IsValid(ePrefix + "Invalid 'exponent' - ")

	if err != nil {
		return &iaReturn, err
	}

	if base.IsZero() {
		return &iaReturn,
			errors.New(ePrefix + "'base' is zero value. INVALID INPUT!")
	}

	if exponent.IsZero() {
		iaReturn.SetIntAryToOne(minResultPrecision)
		return &iaReturn, nil
	}

	iaOne := IntAry{}.NewOne(exponent.GetPrecision())

	if exponent.Equals(&iaOne) {
		iaReturn = base.CopyOut()
		return &iaReturn, nil
	}

	exponentPrecision := exponent.GetPrecision()
	exponentSign := exponent.GetSign()

	if exponentPrecision == 0 && exponentSign == 1 {
		return iaPwr.pwrMultiplyPositiveIntegerExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision == 0 && exponentSign == -1 {
		return iaPwr.pwrMultiplyNegativeIntegerExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision > 0 && exponentSign == 1 {
		return iaPwr.pwrMultiplyPositiveFractionalExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	if exponentPrecision > 0 && exponentSign == -1 {
		return iaPwr.pwrMultiplyNegativeFractionalExponent(
			base,
			exponent,
			minResultPrecision,
			maxResultPrecision)
	}

	return &iaReturn,
		errors.New(ePrefix + "Error: input parameters failed to match valid calculation types!")
}

// pwrByTwos - Raises a *big.Int 'base', to the specified 'power'
// using the Exponentiation by squaring algorithm.
//
// See:
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
// https://en.wikipedia.org/wiki/Exponentiation_by_squaring#Computation_by_powers_of_2
// This method is based on revised code taken in part from Ye Lin Aung.
// https://stackoverflow.com/questions/30182129/calculating-large-exponentiation-in-golang
// This algorithm modified by Mike Rapp to achieve improved performance.
//
// Input Parameters:
// =================
//
// 'ia' *IntAry -
//			The base which will be raised to an exponent specified by input
//			parameter 'power'.
//
// 'power' *big.Int -
// 				The input parameter 'power' may be either
// 				a positive or negative integer.
//
// 'maxResultPrecision' int -
//				'maxResultPrecision' will determine the maximum
// 				number of digits to the right of the decimal
//				place in the result.
//
//				Valid values are -1 and values >= zero ('0')
//        Values less than -1 will trigger an error.
//
//				A value of -1 signals that no limit will be placed on
//				the number of decimals places to right of the decimal
//				point in the result.
//
//	'internalPrecision' int -
// 				'internalPrecision' will control the number of digits of
//				accuracy to the right of the decimal point maintained by
//				internal multiplication operations used in raising the intAry
//				value to the designated power.
//
//				Valid values are -1 and values >= zero ('0')
//        Values less than -1 will trigger an error.
//
//				A value of -1 signals that no limit will be placed on
//				the number of decimals places to right of the decimal
//				point during internal multiplication operations.
//
//  Return Values
//  =============
//
//  ia			The result of this operation is returned through the pointer to
//					input parameter, 'ia'. As such, the original value of 'ia' will
//					be overwritten. The The returned value 'ia' will contain the same
// 					numeric separators (decimal separator, thousands separator and
// 					currency symbol) as that of the original 'ia' instance. 'ia'
// 					numeric separators will therefore remain unchanged.
//
//  error		If, during the execution of this method, an error is identified,
//					method will set the returned error to a non nil value and return
//					an error message.
//
func (iaPwr *IntAryMathPower) pwrByTwos(
	ia *IntAry,
	power *big.Int,
	maxResultPrecision,
	internalPrecision int) error {

	ePrefix := "IntAryMathPower.PwrByTwos() "

	if maxResultPrecision < -1 {
		return fmt.Errorf(ePrefix+
			"Error: Parameter maxResultPrecision is less than -1. "+
			"maxResultPrecision= %v", maxResultPrecision)
	}

	if internalPrecision < -1 {
		return fmt.Errorf(ePrefix+
			"Error: Parameter internalPrecision is less than -1. "+
			"internalPrecision= %v", internalPrecision)
	}

	numSeps := ia.GetNumericSeparatorsDto()

	if maxResultPrecision == internalPrecision {
		internalPrecision += 20
	}

	ia.SetInternalFlags()
	tPower := big.NewInt(0).Set(power)
	one := big.NewInt(1)
	zero := big.NewInt(0)
	two := big.NewInt(2)

	if ia.isZeroValue {
		ia.SetIntAryToZero(0)
		ia.SetNumericSeparatorsDto(numSeps)
		return nil
	}

	if tPower.Cmp(two) == -1 {

		if tPower.Cmp(zero) == -1 {

			ia2, err := ia.Inverse(internalPrecision)

			if err != nil {
				return fmt.Errorf(ePrefix+
					"Error returned from ia.Inverse(-1): Error= %v", err.Error())
			}

			ia.CopyIn(&ia2, false)
			tPower = big.NewInt(0).Mul(tPower, big.NewInt(-1))

		} else if tPower.Cmp(one) == 0 {
			// no change in value. x^1 == x
			ia.SetNumericSeparatorsDto(numSeps)
			return nil
		} else if tPower.Cmp(zero) == 0 {
			ia.SetIntAryToOne(0)
			ia.SetNumericSeparatorsDto(numSeps)
			return nil
		}

	}

	tBase := ia.CopyOut()

	ia.SetIntAryToOne(0)

	for tPower.Cmp(zero) == 1 {
		//temp, _:= intAry{}.NewNumStr("0")

		if big.NewInt(0).Mod(tPower, two).Cmp(one) == 0 {
			//temp = big.NewInt(0).Mul(result, tBase)
			//result = big.NewInt(0).Set(temp)
			err := ia.MultiplyThisBy(&tBase, -1, internalPrecision)
			//fmt.Println("ia precision = ", ia.GetPrecision())

			if err != nil {
				return fmt.Errorf(ePrefix+
					"Error From result.MultiplyThisBy(&tBase, true). Error= %v", err.Error())
			}

			if tPower.Cmp(one) == 0 {
				if maxResultPrecision > -1 && maxResultPrecision < ia.GetPrecision() {
					ia.SetPrecision(maxResultPrecision, true)
				}

				ia.SetNumericSeparatorsDto(numSeps)
				return nil
			}
		}

		err := tBase.MultiplyThisBy(&tBase, -1, internalPrecision)

		if err != nil {
			return fmt.Errorf(ePrefix+
				"Error From tBase.MultiplyThisBy(&temp, true). Error= %v", err.Error())
		}

		tPower = big.NewInt(0).Div(tPower, two)
	}

	ia.SetNumericSeparatorsDto(numSeps)

	if maxResultPrecision == -1 {
		return nil
	}

	if maxResultPrecision > -1 && maxResultPrecision < ia.GetPrecision() {
		ia.SetPrecision(maxResultPrecision, true)
	}

	return nil
}

// pwrMultiplyPositiveIntegerExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a positive integer.
// If 'exponent' is NOT a positive integer, an error will be thrown.
//
// If 'exponent' is a fractional value, i.e., it has digits to the right of the
// decimal place, it is by definition NOT an integer value and an error will
// be thrown.
//
// This method uses simple multiplication to generate the result.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of the power operation is returned as a pointer to a new
// 'result' IntAry. None of the input parameters are altered by this
// operation. The returned 'result' IntAry will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter 'base'.
//
// Note: This method does not perform tests for base==0, exponent==0 or exponent==1.
// It is assumed that these tests were performed before calling this method.
//
func (iaPwr *IntAryMathPower) pwrMultiplyPositiveIntegerExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) (*IntAry, error) {

	ePrefix := "IntAryMathPower.pwrMultiplyPositiveIntegerExponent() "
	iaErrReturn := IntAry{}.NewZero(0)

	if exponent.GetSign() != 1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a positive integer. "+
				"Instead, 'exponent' is negative! exponent='%v'",
				exponent.GetNumStr())
	}

	if exponent.GetPrecision() != 0 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be an integer value. "+
				"Instead, 'exponent' is a fractional value! exponent='%v'",
				exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	result := IntAry{}.NewOne(0)

	internalMaxPrecision := maxResultPrecision + 100

	opExponent := exponent.CopyOut()

	for !opExponent.IsZero() {

		IntAryMathMultiply{}.MultiplyInPlace(&result, base, minResultPrecision, internalMaxPrecision)

		err := opExponent.DecrementIntegerOne()

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by opExponent.DecrementIntegerOne() "+
					"Error='%v' ", err.Error())
		}

	}

	if result.GetPrecision() > maxResultPrecision {
		result.RoundToPrecision(maxResultPrecision)
	}

	if result.GetPrecision() < minResultPrecision {
		result.SetPrecision(minResultPrecision, false)
	}

	err := result.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by result.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' ", err.Error())
	}

	return &result, nil
}

// pwrMultiplyNegativeIntegerExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a negative integer value.
// If 'exponent' is a positive integer, an error will be thrown.
//
// If 'exponent' is a fractional value, i.e., it has digits to the right of the
// decimal place, it is by definition NOT an integer value and therefore, an
// error will be thrown.
//
// This method uses simple multiplication to generate the result.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of this power operation is returned as a pointer to a new
// 'result' IntAry. None of the input parameters are altered by this
// operation. The returned 'result' IntAry will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter 'base'.
//
// Note: This method does not perform tests for base==0 .
// It is assumed that this test was performed before calling this method.
//
func (iaPwr *IntAryMathPower) pwrMultiplyNegativeIntegerExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) (*IntAry, error) {

	ePrefix := "IntAryMathPower.pwrMultiplyNegativeIntegerExponent() "
	iaErrReturn := IntAry{}.NewZero(0)

	if exponent.GetSign() != -1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a negative integer. "+
				"Instead, 'exponent' is positive! exponent='%v'",
				exponent.GetNumStr())
	}

	if exponent.GetPrecision() > 0 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be an integer value. "+
				"Instead, 'exponent' is a fractional value! exponent='%v'",
				exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	internalMaxPrecision := maxResultPrecision + 100

	newBase, err := base.Inverse(internalMaxPrecision)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix +
				"Error returned by base.Inverse(internalMaxPrecision)")
	}

	result := IntAry{}.NewOne(0)
	opExponent := exponent.CopyOut()
	opExponent.ChangeSign()
	iaZero := IntAry{}.NewZero(0)
	internalMaxPrecision += 5

	for !opExponent.Equals(&iaZero) {

		IntAryMathMultiply{}.MultiplyInPlace(&result, &newBase, minResultPrecision, internalMaxPrecision)

		err := opExponent.DecrementIntegerOne()

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by opExponent.DecrementIntegerOne() "+
					"Error='%v' ", err.Error())
		}

	}

	if result.GetPrecision() > maxResultPrecision {
		result.RoundToPrecision(maxResultPrecision)
	}

	if result.GetPrecision() < minResultPrecision {
		result.SetPrecision(minResultPrecision, false)
	}

	err = result.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' ", err.Error())
	}

	return &result, nil
}

// pwrMultiplyPositiveFractionalExponent - raises 'base' to the power of
// 'exponent'. Input parameter 'exponent' is expected to represent a positive
// fractional value, i.e., a numeric value with digits to the right of the
// decimal place.
//
// If 'exponent' is NOT a positive value, i.e., a value less than zero, an
// error will be thrown.
//
// If 'exponent' is a numeric value with no digits to the right of the
// decimal place, it is by definition an integer value and NOT a fractional
// value. In that case, an error will be thrown.
//
// This method uses simple multiplication to generate the result.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of this power operation is returned as a pointer to a new
// 'result' IntAry. None of the input parameters are altered by this
// operation. The returned 'result' IntAry will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter 'base'.
//
// Note: This method does not perform tests for base==0, exponent==0
//       or exponent==1. Is is assumed that these tests were performed
//       before calling this method.
//
func (iaPwr *IntAryMathPower) pwrMultiplyPositiveFractionalExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) (*IntAry, error) {

	ePrefix := "IntAryMathPower.pwrMultiplyPositiveFractionalExponent() "
	iaErrReturn := IntAry{}.NewZero(0)

	if exponent.GetSign() != 1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a positive fractional value. "+
				"Instead, 'exponent' is a negative value! exponent='%v'",
				exponent.GetNumStr())
	}

	if exponent.GetPrecision() < 1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a fractional value. "+
				"Instead, 'exponent' is an integer value! exponent='%v'",
				exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	internalMaxPrecision := maxResultPrecision + 100

	fracIntAry := FracIntAry{}.NewFracIntAry(exponent)

	err := fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision) "+
				"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	newBase, err := NthRootOp{}.NewNthRoot(base, &fracIntAry.Denominator, internalMaxPrecision)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by NthRootOp{}.NewNthRoot(...) "+
				"Error='%v' ", err.Error())
	}

	result := IntAry{}.NewOne(0)
	internalMaxPrecision += 5

	if fracIntAry.Numerator.GetPrecision() != 0 ||
		fracIntAry.Numerator.GetSign() == -1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+"Error: fracIntAry.Numerator (new exponent) INVALID!  "+
				"fracIntAry.Numerator='%v' ", fracIntAry.Numerator.GetNumStr())
	}

	// fracIntAry.Numerator == new exponent
	for !fracIntAry.Numerator.IsZero() {

		err = IntAryMathMultiply{}.MultiplyInPlace(&result, &newBase, minResultPrecision, internalMaxPrecision)

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by IntAryMathMultiply{}.MultiplyInPlace() "+
					"Error='%v' ", err.Error())
		}

		err = fracIntAry.Numerator.DecrementIntegerOne()

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by fracIntAry.Numerator.DecrementIntegerOne() "+
					"Error='%v' ", err.Error())
		}

	}

	if result.GetPrecision() > maxResultPrecision {
		result.RoundToPrecision(maxResultPrecision)
	}

	if result.GetPrecision() < minResultPrecision {
		result.SetPrecision(minResultPrecision, false)
	}

	err = result.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by fracIntAry.Numerator.DecrementIntegerOne() "+
				"Error='%v' ", err.Error())
	}

	return &result, nil
}

// pwrMultiplyNegativeFractionalExponent - raises 'base' to the power of
// 'exponent'. Input parameter 'exponent' is expected to represent a negative
// fractional value, i.e., a numeric value with digits to the right of the
// decimal place.
//
// Input parameter must be negative numeric value with a value less than zero.
// If 'exponent' is a positive value, an error will be thrown.
//
// If 'exponent' is a numeric value with no digits to the right of the
// decimal place, it is by definition an integer value and NOT a fractional
// value. In that case, an error will be thrown.
//
// This method uses simple multiplication to generate the result.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of this power operation is returned as a pointer to a new
// 'result' IntAry. None of the input parameters are altered by this
// operation. The returned 'result' IntAry will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// input parameter 'base'.
//
// Note: This method does not perform a test for base==0. It is assumed
//       this test was performed before calling this method.
//
func (iaPwr *IntAryMathPower) pwrMultiplyNegativeFractionalExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) (*IntAry, error) {

	ePrefix := "IntAryMathPower.pwrMultiplyNegativeFractionalExponent() "
	iaErrReturn := IntAry{}.NewZero(0)

	if exponent.GetSign() != -1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a negative fractional value. "+
				"Instead, 'exponent' is a positive value! exponent='%v'",
				exponent.GetNumStr())
	}

	if exponent.GetPrecision() < 1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error: 'exponent' is expected to be a fractional value. "+
				"Instead, 'exponent' is an integer value! exponent='%v'",
				exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	internalMaxPrecision := maxResultPrecision + 100

	newExponent := exponent.CopyOut()

	// Change sign from minus to plus
	newExponent.ChangeSign()

	newInverseBase, err := base.Inverse(internalMaxPrecision)

	if err != nil {
		return &iaErrReturn, fmt.Errorf(ePrefix+
			"Error returned by base.Inverse(internalMaxPrecision) "+
			"Error='%v' ", err.Error())
	}

	fracIntAry := FracIntAry{}.NewFracIntAry(&newExponent)

	internalMaxPrecision += 5

	err = fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision) "+
				"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	newBase, err := NthRootOp{}.NewNthRoot(&newInverseBase, &fracIntAry.Denominator, internalMaxPrecision)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by NthRootOp{}.NewNthRoot(...) "+
				"Error='%v' ", err.Error())
	}

	result := IntAry{}.NewOne(0)

	internalMaxPrecision += 5

	if fracIntAry.Numerator.GetPrecision() != 0 ||
		fracIntAry.Numerator.GetSign() == -1 {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+"Error: fracIntAry.Numerator (new exponent) INVALID!  "+
				"fracIntAry.Numerator='%v' ", fracIntAry.Numerator.GetNumStr())
	}

	// fracIntAry.Numerator == new exponent
	for !fracIntAry.Numerator.IsZero() {

		err := IntAryMathMultiply{}.MultiplyInPlace(&result, &newBase, minResultPrecision, internalMaxPrecision)

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by IntAryMathMultiply{}.MultiplyInPlace(...) "+
					"Error='%v' ", err.Error())
		}

		err = fracIntAry.Numerator.DecrementIntegerOne()

		if err != nil {
			return &iaErrReturn,
				fmt.Errorf(ePrefix+
					"Error returned by fracIntAry.Numerator.DecrementIntegerOne() "+
					"Error='%v' ", err.Error())
		}
	}

	if result.GetPrecision() > maxResultPrecision {
		result.RoundToPrecision(maxResultPrecision)
	}

	if result.GetPrecision() < minResultPrecision {
		result.SetPrecision(minResultPrecision, false)
	}

	err = result.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return &iaErrReturn,
			fmt.Errorf(ePrefix+
				"Error returned by result.SetNumericSeparatorsDto(numSeps) "+
				"Error='%v' ", err.Error())
	}

	return &result, nil
}

// pwrTwoPositiveIntegerExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a positive integer.
// If 'exponent' is NOT a positive integer, an error will be thrown.
//
// If 'exponent' is a fractional value, i.e., it has digits to the right of the
// decimal place, it is by definition NOT an integer value and therefore, an
// error will be thrown.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of the power operation is returned in the input parameter
// 'base'. During this procedure the original value of 'base' is destroyed.
// The returned IntAry object 'base' will contain same numeric separators
// (decimal separator, thousands separator and currency symbol) as those
// in the original 'base' instance. As such, the 'base' numeric separators
// will remain unchanged.
//
func (iaPwr *IntAryMathPower) pwrTwoPositiveIntegerExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) error {

	ePrefix := "IntAryMathPower.pwrTwoPositiveIntegerExponent() "

	if exponent.GetSign() != 1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a positive integer. "+
			"Instead, 'exponent' is negative! exponent='%v'",
			exponent.GetNumStr())
	}

	if exponent.GetPrecision() != 0 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be an integer value. "+
			"Instead, 'exponent' is a fractional value! exponent='%v'",
			exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	// This is an integer Exponent!
	bInt, err := exponent.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by exponent.GetBigInt() Error='%v'\n",
			err.Error())
	}

	internalPrecision := maxResultPrecision + 100

	err = iaPwr.pwrByTwos(base, bInt, maxResultPrecision, internalPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by iaPwr.PwrByTwos(...) Error='%v' ", err.Error())
	}

	if base.GetPrecision() > maxResultPrecision {
		base.RoundToPrecision(maxResultPrecision)
	}

	if base.GetPrecision() < minResultPrecision {
		base.SetPrecision(minResultPrecision, false)
	}

	err = base.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by base.SetNumericSeparatorsDto(numSeps) Error='%v' ",
			err.Error())
	}

	return nil
}

// pwrTwoNegativeIntegerExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a negative integer
// value, i.e., a value less than zero.
//
// If 'exponent' is a positive integer (value greater than -1), an error
// will be thrown.
//
// If 'exponent' is a fractional value, i.e., it has digits to the right of the
// decimal place, it is by definition NOT an integer value and therefore, an
// error will be thrown.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of the power operation is returned in the input parameter
// 'base'. During this procedure the original value of 'base' is destroyed.
// The returned IntAry object 'base' will contain same numeric separators
// (decimal separator, thousands separator and currency symbol) as those
// in the original 'base' instance. As such, the 'base' numeric separators
// will remain unchanged.
//
func (iaPwr *IntAryMathPower) pwrTwoNegativeIntegerExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) error {

	ePrefix := "IntAryMathPower.pwrTwoNegativeIntegerExponent() "

	if exponent.GetSign() != -1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a negative integer. "+
			"Instead, 'exponent' is positive! exponent='%v'",
			exponent.GetNumStr())
	}

	if exponent.GetPrecision() != 0 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be an integer value. "+
			"Instead, 'exponent' is a fractional value! exponent='%v'",
			exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	// This is an integer Exponent!
	bInt, err := exponent.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by exponent.GetBigInt() Error='%v'\n",
			err.Error())
	}

	internalPrecision := maxResultPrecision + 100

	err = iaPwr.pwrByTwos(base, bInt, maxResultPrecision, internalPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by iaPwr.PwrByTwos(...) Error='%v' ", err.Error())
	}

	if base.GetPrecision() > maxResultPrecision {
		base.RoundToPrecision(maxResultPrecision)
	}

	if base.GetPrecision() < minResultPrecision {
		base.SetPrecision(minResultPrecision, false)
	}

	err = base.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by base.SetNumericSeparatorsDto(numSeps) Error='%v' ",
			err.Error())
	}

	return nil
}

// pwrTwoPositiveFractionalExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a positive fractional
// value, i.e., a value greater than -1 with digits to the right of the
// decimal point.
//
// If 'exponent' is a negative value, i.e., a value less than zero, an error
// will be thrown.
//
// Also, if 'exponent' has a precision value less than '1', an error will be
// thrown.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of the power operation is returned in the input parameter
// 'base'. During this procedure the original value of 'base' is destroyed.
// The returned IntAry object 'base' will contain same numeric separators
// (decimal separator, thousands separator and currency symbol) as those
// in the original 'base' instance. As such, the 'base' numeric separators
// will remain unchanged.
//
func (iaPwr *IntAryMathPower) pwrTwoPositiveFractionalExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) error {

	ePrefix := "IntAryMathPower.pwrTwoPositiveFractionalExponent() "

	if exponent.GetSign() != 1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a positive fractional value. "+
			"Instead, 'exponent' is negative! exponent='%v'",
			exponent.GetNumStr())
	}

	if exponent.GetPrecision() < 1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a fractional value. "+
			"Instead, 'exponent' is an integer value! exponent='%v'",
			exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	internalMaxPrecision := maxResultPrecision + 100

	fracIntAry := FracIntAry{}.NewFracIntAry(exponent)

	err := fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision) "+
			"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	newBase, err := NthRootOp{}.NewNthRoot(base, &fracIntAry.Denominator, internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by NthRootOp{}.NewNthRoot(...) "+
			"Error='%v' ", err.Error())
	}

	base.CopyIn(&newBase, false)

	expBigInt, err := fracIntAry.Numerator.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by fracIntAry.Numerator.GetBigInt() "+
			"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	err = iaPwr.pwrByTwos(base, expBigInt, maxResultPrecision, internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by iaPwr.PwrByTwos(...) "+
			"Error='%v' ", err.Error())
	}

	if base.GetPrecision() > maxResultPrecision {
		base.RoundToPrecision(maxResultPrecision)
	}

	if base.GetPrecision() < minResultPrecision {
		base.SetPrecision(minResultPrecision, false)
	}

	err = base.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by base.SetNumericSeparatorsDto(numSeps) "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// pwrTwoNegativeFractionalExponent - raises 'base' to the power of 'exponent'.
// Input parameter 'exponent' is expected to represent a negative fractional
// value, i.e., a value less than zero with digits to the right of the decimal
// point.
//
// If 'exponent' is a positive value, i.e., a greater than than -1, an error
// will be thrown.
//
// Also, if 'exponent' has a precision value less than '1', an error will be
// thrown.
//
// Input parameter 'maxResultPrecision' will round the result to this
// number of decimal places after the decimal point if the result is
// greater than 'maxResultPrecision'.  If the value of 'maxResultPrecision'
// is less than zero, it will be automatically set to a value of '4096'.
//
// Input parameter 'minResultPrecision' signals that if the result precision
// is less than 'minResultPrecision', zeros will be added to the right of
// the decimal place in order to implement the 'minResultPrecision'
// specification. If the value of 'minResultPrecision' is less than zero,
// 'minResultPrecision' will be automatically set to a value of zero.
//
// The result of the power operation is returned in the input parameter
// 'base'. During this procedure the original value of 'base' is destroyed.
// The returned IntAry object 'base' will contain same numeric separators
// (decimal separator, thousands separator and currency symbol) as those
// in the original 'base' instance. As such, the 'base' numeric separators
// will remain unchanged.
//
func (iaPwr *IntAryMathPower) pwrTwoNegativeFractionalExponent(
	base, exponent *IntAry,
	minResultPrecision, maxResultPrecision int) error {

	ePrefix := "IntAryMathPower.pwrTwoNegativeFractionalExponent() "

	if exponent.GetSign() != -1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a negative fractional value. "+
			"Instead, 'exponent' is positive! exponent='%v'",
			exponent.GetNumStr())
	}

	if exponent.GetPrecision() < 1 {
		return fmt.Errorf(ePrefix+
			"Error: 'exponent' is expected to be a fractional value. "+
			"Instead, 'exponent' is an integer value! exponent='%v'",
			exponent.GetNumStr())
	}

	numSeps := base.GetNumericSeparatorsDto()

	if maxResultPrecision < 0 {
		maxResultPrecision = 4096
	}

	if minResultPrecision < 0 {
		minResultPrecision = 0
	}

	if minResultPrecision > maxResultPrecision {
		minResultPrecision = maxResultPrecision
	}

	internalMaxPrecision := maxResultPrecision + 100

	// Set exponent to a positive value
	newExponent := exponent.CopyOut()

	newExponent.ChangeSign()

	newInverseBase, err := base.Inverse(internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by base.Inverse(internalMaxPrecision) "+
			"Error='%v' ", err.Error())
	}

	fracIntAry := FracIntAry{}.NewFracIntAry(&newExponent)

	internalMaxPrecision += 5

	err = fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by fracIntAry.ReduceToLowestCommonDenom(internalMaxPrecision) "+
			"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	newBase, err := NthRootOp{}.NewNthRoot(&newInverseBase, &fracIntAry.Denominator, internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by NthRootOp{}.NewNthRoot(...) "+
			"Error='%v' ", err.Error())
	}

	expBigInt, err := fracIntAry.Numerator.GetBigInt()

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by newExponent.GetBigInt() "+
			"Error='%v' ", err.Error())
	}

	internalMaxPrecision += 5

	err = iaPwr.pwrByTwos(&newBase, expBigInt, maxResultPrecision, internalMaxPrecision)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by iaPwr.PwrByTwos(...) "+
			"Error='%v' ", err.Error())
	}

	base.CopyIn(&newBase, false)

	if base.GetPrecision() > maxResultPrecision {
		base.RoundToPrecision(maxResultPrecision)
	}

	if base.GetPrecision() < minResultPrecision {
		base.SetPrecision(minResultPrecision, false)
	}

	base.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by base.SetNumericSeparatorsDto(numSeps) "+
			"Error='%v' ", err.Error())
	}

	return nil
}
