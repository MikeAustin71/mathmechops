package mathops

import (
	"errors"
	"fmt"
)

type IntAryMathDivide struct {
	Input  IntAryPair
	Result IntAry
}

// Divide - Divides the parameter 'dividend' by the parameter 'divisor'.
// The result of this division is a 'quotient' which is returned as an
// IntAry type.
//
// The dividend is the number being divided, the divisor is
// the number by which the dividend is divided, and the quotient
// is the result of the division.
//
// Given a ÷ b = c, a is the dividend, b is the divisor and c is
// the quotient.
//
// Maximum precision of the division result is controlled by the input
// parameter, 'maxPrecision'.
//
// If 'maxPrecision' is greater than or equal to zero ('0'),
// the number of digits to the right of the decimal place will
// not exceed 'maxPrecision'.
//
// If 'maxPrecision' is set equal to minus one ('-1'), 'maxPrecision'
// will be automatically set to a maximum of 4,096 digits to the right
// of the decimal point.
//
// 'minPrecision' specifies the minimum precision of the final result.
// If 'minPrecision' is less than zero, it is automatically set to zero.
//
// Returns
// =======
// IntAry - If the division calculation completes successfully, the
//						computed quotient is returned as a type IntAry.
// 						The return value 'IntAry' will contain the numeric separators
// 						(decimal separator, thousands separator and currency symbol)
// 						copied from input parameter 'dividend'.
//
/// error	- If the division calculation completes successfully this
//						error type is set equal to 'nil'. Otherwise, it will
//						contain an error message.
//
func (iaDivide IntAryMathDivide) Divide(
	dividend, divisor *IntAry,
	minPrecision, maxPrecision int) (IntAry, error) {

	ePrefix := "IntAryMathDivide.Divide() "

	dividend.SetInternalFlags()
	divisor.SetInternalFlags()

	if divisor.isZeroValue {
		return IntAry{}.New(), errors.New(ePrefix + "Error: divide by zero")
	}

	if maxPrecision < -1 {
		return IntAry{}.New(),
			errors.New(ePrefix +
				"Error: Input parameter 'maxPrecision' is INVALID. 'maxPrecision' is less than -1")
	}

	if maxPrecision == -1 {
		maxPrecision = 4096
	}

	if minPrecision < 0 {
		minPrecision = 0
	}

	if minPrecision > maxPrecision {
		minPrecision = maxPrecision
	}

	quotient := IntAry{}.New()
	quotient.SetIntAryToZero(0)

	if dividend.isZeroValue {
		return quotient, nil
	}

	numSeps := dividend.GetNumericSeparatorsDto()

	trialDividend := dividend.CopyOut()

	tempDivisor := divisor.CopyOut()

	tensCount := IntAry{}.New()
	tensCount.SetIntAryToOne(0)

	newSignVal := 1

	if trialDividend.signVal != tempDivisor.signVal {
		newSignVal = -1
	}

	if trialDividend.signVal == -1 {
		trialDividend.signVal = 1
	}

	if tempDivisor.signVal == -1 {
		tempDivisor.signVal = 1
	}

	dividendMag := trialDividend.GetMagnitudeDigits()
	divisorMag := tempDivisor.GetMagnitudeDigits()
	deltaMag := uint(0)
	incrementVal := IntAry{}.New()
	incrementVal = tempDivisor.CopyOut()

	if dividendMag > divisorMag {
		deltaMag = uint(dividendMag - divisorMag)
		tensCount.MultiplyByTenToPower(deltaMag)
		incrementVal.MultiplyThisBy(&tensCount, -1, -1)

	} else if divisorMag > dividendMag {
		deltaMag = uint(divisorMag - dividendMag)
		trialDividend.MultiplyByTenToPower(deltaMag)
		tensCount.DivideByTenToPower(deltaMag)

	}

	compare := 0
	precisionCutOff := maxPrecision + dividendMag + 1
	var err error

	for true {

		if quotient.precision >= precisionCutOff {
			quotient.RoundToPrecision(maxPrecision)
			quotient.OptimizeIntArrayLen(true)
			quotient.signVal = newSignVal

			if quotient.GetPrecision() > maxPrecision {
				quotient.RoundToPrecision(maxPrecision)
			}

			if quotient.GetPrecision() < minPrecision {
				quotient.SetPrecision(minPrecision, false)
			}

			err = quotient.SetNumericSeparatorsDto(numSeps)

			if err != nil {
				return IntAry{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by quotient.SetNumericSeparatorsDto(numSeps). "+
						"Error='%v' ", err.Error())
			}

			return quotient, nil
		}

		compare = incrementVal.CompareAbsoluteValues(&trialDividend)

		if compare == 0 {
			// incrementalVal is equal to trialDividend
			quotient.AddToThis(&tensCount)

			quotient.OptimizeIntArrayLen(true)
			quotient.signVal = newSignVal

			if quotient.GetPrecision() > maxPrecision {
				quotient.RoundToPrecision(maxPrecision)
			}

			if quotient.GetPrecision() < minPrecision {
				quotient.SetPrecision(minPrecision, false)
			}

			err = quotient.SetNumericSeparatorsDto(numSeps)

			if err != nil {
				return IntAry{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by quotient.SetNumericSeparatorsDto(numSeps). "+
						"Error='%v' ", err.Error())
			}

			return quotient, nil

		} else if compare == -1 {
			// incrementalVal < trialDividend
			quotient.AddToThis(&tensCount)

			// Calc Remainder
			trialDividend.SubtractFromThis(&incrementVal)

			continue

		} else {
			// Must Be compare == 1
			// incrementalVal > trialDividend

			tensCount.DivideByTenToPower(1)
			incrementVal.DivideByTenToPower(1)
		}

	}

	if quotient.GetPrecision() > maxPrecision {
		quotient.RoundToPrecision(maxPrecision)
	}

	if quotient.GetPrecision() < minPrecision {
		quotient.SetPrecision(minPrecision, false)
	}

	err = quotient.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return IntAry{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by quotient.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return quotient, nil
}

// DivideByInt64 - Divide the IntAry parameter 'ia' by an int64 divisor and returns
// the quotient in the pointer to the 'ia' parameter. Consequently, the original value
// of 'ia' will be overwritten and replaced by the resulting quotient. The returned
// value 'ia' will contain the same numeric separators (decimal separator, thousands
// separator and currency symbol) as that of the original 'ia' instance. 'ia' numeric
// separators will therefore remain unchanged.
//
// If the quotient has a number of decimal places to the right of the decimal point which is
// greater than 'maxPrecision', the result is rounded to 'maxPrecision' decimal places.
//
// If 'maxPrecision' is set equal to -1, 'maxPrecision' is automatically set to 4,096.
//
// If 'maxPrecision' is less than -1, an error will be returned.
//
func (iaDivide IntAryMathDivide) DivideByInt64(
	ia *IntAry, divisor int64, maxPrecision int) error {

	ePrefix := "iaDivide IntAryMathDivide() "

	if divisor == 0 {
		return errors.New(ePrefix + "'divisor' Equals zero. Cannot divide by zero! \n")
	}

	if maxPrecision < -1 {
		return fmt.Errorf(ePrefix+"Error: Input Parameter 'maxPrecision' is less than -1! "+
			"maxPrecision='%v' \n", maxPrecision)
	}

	if maxPrecision == -1 {
		maxPrecision = 4096
	}

	ia.OptimizeIntArrayLen(false)

	if ia.isZeroValue {

		ia.SetIntAryToZero(ia.GetPrecisionUint())
		return nil
	}

	dSignVal := 1

	if divisor < 0 {
		dSignVal = -1
		divisor = divisor * -1
	}

	ia.signVal = dSignVal * ia.signVal

	n1 := int64(0)
	n2 := int64(0)
	carry := int64(0)
	iMaxPrecision := int(maxPrecision) + 1
	newAryLen := ia.intAryLen
	intAryLen := ia.intAryLen - ia.precision
	precisionCnt := 0

	for i := 0; i < newAryLen; i++ {

		if i >= intAryLen {
			precisionCnt++
		}

		if i < ia.intAryLen {
			n1 = int64(ia.intAry[i]) + carry
		} else {
			n1 = int64(0) + carry
		}

		n2 = n1 / divisor
		carry = (n1 - (n2 * divisor)) * 10

		if i < ia.intAryLen {
			ia.intAry[i] = uint8(n2)
		} else {
			ia.intAry = append(ia.intAry, uint8(n2))
		}

		if i == newAryLen-1 &&
			carry > 0 && precisionCnt <= iMaxPrecision {

			newAryLen++

		}

	}

	ia.precision = precisionCnt

	ia.intAryLen = newAryLen

	if precisionCnt >= iMaxPrecision {
		iMaxPrecision--
		ia.RoundToPrecision(iMaxPrecision)
	}

	if ia.intAry[0] == 0 {
		ia.SetSignificantDigitIdxs()
		ia.intAry = ia.intAry[ia.firstDigitIdx:]
		ia.SetIntAryLength()
	}

	return nil

}

// DivideByTenToPower - Divide input parameter 'ia' of type IntAry by 10 raised
// to the power of input parameter 'power'.
//
// The result, or quotient, is returned via the pointer to input parameter 'ia'.
// Consequently, the original value of 'ia' will be overwritten and replaced by the
// resulting quotient. The returned value 'ia' will contain the same numeric
// separators (decimal separator, thousands separator and currency symbol)
// as that of the original 'ia' instance. 'ia' numeric separators will
// therefore remain unchanged.
//
func (iaDivide IntAryMathDivide) DivideByTenToPower(ia *IntAry, power uint) {

	if power == 0 {
		return
	}

	ia.precision += int(power)
	ia.intAryLen = len(ia.intAry)
	newLen := ia.precision + 1

	if ia.intAryLen < newLen {

		t := make([]uint8, newLen)

		deltaLen := newLen - ia.intAryLen

		for i := 0; i < newLen; i++ {

			if i < deltaLen {
				t[i] = 0
			} else {
				t[i] = ia.intAry[i-deltaLen]
			}

		}

		ia.intAry = make([]uint8, newLen)
		for i := 0; i < newLen; i++ {
			ia.intAry[i] = t[i]
		}

		ia.intAryLen = newLen
	}

	ia.OptimizeIntArrayLen(false)
}

// DivideByTwoQuoMod - Receives an input parameter of pointer to
// an IntAry instance ('ia'). 'ia' is then divided by two (2). The
// result, or quotient, is returned via the pointer to input parameter
// 'ia'. Consequently, the original value of 'ia' will be overwritten
// and replaced by the resulting quotient. The returned value 'ia' will
// contain the same numeric separators (decimal separator, thousands
// separator and currency symbol) as that of the original 'ia' instance.
// 'ia' numeric separators will therefore remain unchanged.
//
func (iaDivide IntAryMathDivide) DivideByTwo(ia *IntAry) {

	ia.OptimizeIntArrayLen(false)

	if ia.isZeroValue {

		ia.SetIntAryToZero(ia.GetPrecisionUint())

		return
	}

	n1 := uint8(0)
	n2 := uint8(0)
	carry := uint8(0)

	for i := 0; i < ia.intAryLen; i++ {

		n1 = ia.intAry[i] + carry
		n2 = n1 / 2
		carry = (n1 - (n2 * 2)) * 10
		ia.intAry[i] = n2

	}

	if carry > 0 {
		ia.intAry = append(ia.intAry, 5)
		ia.intAryLen++
		ia.precision++
	}

	if ia.intAry[0] == 0 {
		ia.SetSignificantDigitIdxs()
		ia.intAry = ia.intAry[ia.firstDigitIdx:]
		ia.SetIntAryLength()
	}

}
