package mathops

import (
	"errors"
	"fmt"
)

/*
	StrMathOp
	=========

	The source code repository for decimal.go is located at:
			https://github.com/MikeAustin71/mathopsgo.git

	The source file decimal.go is located in directory:
		MikeAustin71/mathopsgo/mathops/strmathop.go

*/

// StrMathOp - Used to perform string based math operations.
// The are stored internally in a series of data fields.
//
// Dependencies:
// 	IntAry
//
type StrMathOp struct {
	N1       IntAry
	N2       IntAry
	N3       IntAry
	IntMAry  [][]int
	IFinal   IntAry
	Dividend IntAry
	Divisor  IntAry
	Quotient IntAry
	Modulo   IntAry
}

func (sMathOp StrMathOp) New() StrMathOp {
	iAry := StrMathOp{}
	iAry.N1 = IntAry{}.New()
	iAry.N2 = IntAry{}.New()
	iAry.N3 = IntAry{}.New()
	iAry.IntMAry = make([][]int, 0)
	iAry.IFinal = IntAry{}.New()
	iAry.Dividend = IntAry{}.New()
	iAry.Divisor = IntAry{}.New()
	iAry.Quotient = IntAry{}.New()
	iAry.Modulo = IntAry{}.New()
	return iAry
}

func (sMathOp *StrMathOp) Empty() {

	sMathOp.N1 = IntAry{}.New()
	sMathOp.N2 = IntAry{}.New()
	sMathOp.N3 = IntAry{}.New()
	sMathOp.IntMAry = make([][]int, 0)
	sMathOp.IFinal = IntAry{}.New()
	sMathOp.Dividend = IntAry{}.New()
	sMathOp.Divisor = IntAry{}.New()
	sMathOp.Quotient = IntAry{}.New()
	sMathOp.Modulo = IntAry{}.New()

}

// AddN1N2 - Adds the values in the N1 and N2 arrays
// and returns the sum in the IFinal Array. Arrays
// N1 and N2 must be first correctly populated before
// calling this method.
func (sMathOp *StrMathOp) AddN1N2() error {

	sMathOp.IFinal = sMathOp.N1.CopyOut()

	err := sMathOp.IFinal.AddToThis(&sMathOp.N2)

	if err != nil {
		return fmt.Errorf("StrMathOp.AddN1N2() Error returned by "+
			"sMathOp.IFinal.AddToThis(&sMathOp.N2). Error='%v'",
			err.Error())
	}

	return nil
}

// RaiseThisToPower - Raises the value of sMathOp.N1 Int Array
// to the power specified by the 'power' parameter passed
// to this method.
//								N1 ^ power
// Before calling this method, sMathOp.N1 Int Array must
// be set to the desired value.
func (sMathOp *StrMathOp) RaiseToPower(power int) error {

	if power < 0 {
		return fmt.Errorf("error: power is less than zero - power= '%v'", power)
	}

	sMathOp.N1.SetInternalFlags()

	sMathOp.IFinal = sMathOp.N1.CopyOut()

	resultPrecision := sMathOp.N1.GetPrecision()

	if sMathOp.N1.IsZero() {
		sMathOp.IFinal.SetIntAryToZero(uint(resultPrecision))
		return nil
	}

	if power == 0 {

		sMathOp.IFinal.SetIntAryToOne(resultPrecision)
		return nil
	}

	if power == 1 {

		return nil
	}

	resultPrecision = resultPrecision * power

	sMathOp.IFinal.Pow(power, resultPrecision, resultPrecision+1500)

	return nil
}

// MultiplyN1N2 - Mulitplies Array N1 by Array N2 and
// places the result in sMathOp.IFinal.
func (sMathOp *StrMathOp) MultiplyN1N2() error {

	sMathOp.N1.SetInternalFlags()
	sMathOp.N2.SetInternalFlags()

	sMathOp.IFinal = sMathOp.N1.CopyOut()
	sMathOp.IFinal.SetInternalFlags()

	n1Precision := sMathOp.N1.GetPrecision()
	n2Precision := sMathOp.N2.GetPrecision()

	greatestPrecision := n1Precision

	if n2Precision > n1Precision {
		greatestPrecision = n2Precision
	}

	maxPrecision := n1Precision + n2Precision

	err := sMathOp.IFinal.MultiplyThisBy(&sMathOp.N2, greatestPrecision, maxPrecision)

	if err != nil {
		return fmt.Errorf("StrMathOp.MultiplyN1N2() Error returned by "+
			"sMathOp.IFinal.MultiplyThisBy(&sMathOp.N2, maxPrecision) "+
			"Error='%v' ", err.Error())
	}

	return nil
}

// Divide - Divides the Dividend IntAry
// field by the Divisor IntAry field. The results are
// stored int IntAry fields Quotient and Modulo
func (sMathOp *StrMathOp) Divide(maxPrecision int) error {

	sMathOp.Quotient.SetIntAryToZero(0)
	sMathOp.Modulo.SetIntAryToZero(0)
	tensCount := IntAry{}.New()
	tensCount.SetIntAryToOne(0)

	newSignVal := 1

	if sMathOp.Divisor.GetSign() != sMathOp.Dividend.GetSign() {
		newSignVal = -1
	}

	if sMathOp.Divisor.GetSign() == -1 {
		sMathOp.Divisor.SetSign(1)
	}

	if sMathOp.Dividend.GetSign() == -1 {
		sMathOp.Dividend.SetSign(1)
	}

	sMathOp.Divisor.SetIsZeroValue()
	if sMathOp.Divisor.IsZero() {
		return errors.New("divisor is zero - divide by zero error")
	}

	sMathOp.Dividend.SetIsZeroValue()

	if sMathOp.Dividend.IsZero() {
		return nil
	}

	trialDividend := sMathOp.Dividend.CopyOut()

	dividendMag := sMathOp.Dividend.GetMagnitudeDigits()
	divisorMag := sMathOp.Divisor.GetMagnitudeDigits()
	deltaMag := uint(0)
	incrementVal := IntAry{}.New()
	incrementVal.SetIntAryWithNumStr(sMathOp.Divisor.GetNumStr())

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

	for true {

		if sMathOp.Quotient.GetPrecision() == precisionCutOff {
			sMathOp.Quotient.SetSign(newSignVal)
			sMathOp.Quotient.RoundToPrecision(maxPrecision)
			return nil
		}

		compare = incrementVal.CompareAbsoluteValues(&trialDividend)

		if compare == 0 {
			// incrementalVal is equal to trialDividend
			sMathOp.Quotient.AddToThis(&tensCount)
			sMathOp.Quotient.SetSign(newSignVal)
			return nil

		} else if compare == -1 {
			// incrementalVal < trialDividend
			sMathOp.Quotient.AddToThis(&tensCount)

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

	return nil
}

// DivideDividendByDivisor - Divides the Dividend IntAry
// field by the Divisor IntAry field. The results are
// stored int IntAry fields Quotient and
func (sMathOp *StrMathOp) DivideDividendByDivisor() error {

	sMathOp.N1 = sMathOp.Divisor
	compare := -1
	quotient := 0
	sN := fmt.Sprintf("%v", quotient)
	sMathOp.N2.SetIntAryWithNumStr(sN)

	for compare < 1 {
		quotient++
		sMathOp.N2.IncrementIntegerOne()
		sMathOp.N3 = sMathOp.IFinal
		sMathOp.MultiplyN1N2()
		compare = sMathOp.IFinal.CompareSignedValues(&sMathOp.Dividend)

	}
	quotient = quotient - 1
	sN = fmt.Sprintf("%v", quotient)
	sMathOp.Quotient.SetIntAryWithNumStr(sN)

	if compare == 0 {
		sMathOp.Modulo.SetIntAryToZero(0)
		return nil
	}

	sMathOp.N1 = sMathOp.Dividend
	sMathOp.N2 = sMathOp.N3
	sMathOp.SubtractN1N2()
	sMathOp.Modulo = sMathOp.IFinal
	return nil
}

func (sMathOp *StrMathOp) DivideBySubtraction() {

	sMathOp.Modulo = sMathOp.Dividend.CopyOut()
	sMathOp.Quotient.SetIntAryToZero(0)
	compare := 1

	for compare >= 0 {
		sMathOp.Modulo.SubtractFromThis(&sMathOp.Divisor)
		sMathOp.Quotient.IncrementIntegerOne()
		compare = sMathOp.Modulo.CompareSignedValues(&sMathOp.Divisor)
	}

	return
}

func (sMathOp *StrMathOp) SubtractDivArys() {

}

// SubtractN1N2 - Subtracts N2 IntAry from
// N1 IntAry. The result is returned in the IFinal
// IntAry.
//
func (sMathOp *StrMathOp) SubtractN1N2() error {

	sMathOp.IFinal = sMathOp.N1.CopyOut()

	err := sMathOp.IFinal.SubtractFromThis(&sMathOp.N2)

	if err != nil {
		return fmt.Errorf("StrMathOp.SubtractN1N2() Error returned by "+
			"sMathOp.IFinal.SubtractFromThis(&sMathOp.N2). Error='%v'",
			err.Error())
	}

	return nil

}
