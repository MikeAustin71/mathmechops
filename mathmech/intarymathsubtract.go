package mathops

type IntAryMathSubtract struct {
	Input  IntAryPair
	Result IntAry
}

// Subtract - Receives two input parameters of type *IntAry. The second parameter
// 'subtrahend' is subtracted from the first parameter, 'minuend'. The result, or
// difference, is returned as a new IntAry instance.
//
// 					'minuend' - 'subtrahend' = difference or result
//
// The two input parameters, 'miunuend' and 'subtrahend', are assumed to be valid
// IntAry instances properly initialized. No validation is performed on 'minuend'
// or 'subtrahend'.
//
// The result of this subtraction operation is returned as an IntAry instance. This
// IntAry will contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter 'minuend'.

func (iaSubtract IntAryMathSubtract) Subtract(minuend, subtrahend *IntAry) IntAry {

	ia3 := minuend.CopyOut()

	iaSubtract.SubtractTotal(&ia3, subtrahend)

	return ia3
}

// SubtractTotal - This method performs a subtraction operation subtracting
// input parameter 'ia2' from input parameter 'ia1'. The result, or difference,
// is returned through use of a pointer in 'ia1'. This means that the original
// value of 'ia1' will be overwritten and destroyed by the subtraction operation.
//
// The returned 'ia1' IntAry will contain numeric separators (decimal separator,
// thousands separator and currency symbol) from the original 'ia1' IntAry
// instance. This means that the numeric separators contained in the original
// ia1 IntAry will remain unchanged.
//
func (iaSubtract IntAryMathSubtract) SubtractTotal(ia1, ia2 *IntAry) {

	numSeps := ia1.GetNumericSeparatorsDto()

	ia1.SetEqualArrayLengths(ia2)

	if ia1.isZeroValue && ia2.isZeroValue {
		ia1.SetIntAryToZero(ia1.GetPrecisionUint())
		return
	}

	compare := ia1.CompareAbsoluteValues(ia2)
	isZeroResult := false

	// Largest Value in now in N1 slot
	newSignVal := ia1.signVal
	doAdd := false
	doReverseNums := false

	if compare == 1 {
		// compare == + 1
		// Absolute Value: N1 > N2

		if ia1.signVal == 1 && ia2.signVal == 1 {
			doAdd = false
			newSignVal = 1
		} else if ia1.signVal == -1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = -1
		} else if ia1.signVal == -1 && ia2.signVal == -1 {
			doAdd = false
			newSignVal = -1
		} else {
			// Must Be ia1.signVal == 1 && ia2.signVal == -1
			doAdd = true
			newSignVal = 1
		}

	} else if compare == -1 {
		// Absolute Values: N2 > N1
		if ia1.signVal == 1 && ia2.signVal == 1 {
			doAdd = false
			doReverseNums = true
			newSignVal = -1
		} else if ia1.signVal == -1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = -1
		} else if ia1.signVal == -1 && ia2.signVal == -1 {
			doAdd = false
			doReverseNums = true
			newSignVal = 1
		} else {
			// Must Be ia1.signVal == 1 && ia2.signVal == -1
			doAdd = true
			newSignVal = 1
		}

	} else {
		// Must be compare == 0
		// Absolute Values: N1==N2
		if ia1.signVal == 1 && ia2.signVal == 1 {
			doAdd = false
			newSignVal = 1
			isZeroResult = true
		} else if ia1.signVal == -1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = -1
		} else if ia1.signVal == -1 && ia2.signVal == -1 {
			doAdd = false
			newSignVal = 1
			isZeroResult = true
		} else {
			// Must Be ia1.signVal == 1 && ia2.signVal == -1
			doAdd = true
			newSignVal = 1
		}

	}

	iaSubtract.addToSubtract(ia1, ia2, newSignVal, doAdd, isZeroResult, doReverseNums)

	ia1.SetNumericSeparatorsDto(numSeps)

	return
}

// addToSubtract - Adds or subtracts two IntAry instances and returns the result
// in the first IntAry parameter,'ia1'.
//
func (iaSubtract *IntAryMathSubtract) addToSubtract(
	ia1, ia2 *IntAry,
	newSignVal int,
	doAdd bool,
	isZeroResult bool,
	doReverseNums bool) {

	if isZeroResult {
		ia1.SetIntAryToZero(ia1.GetPrecisionUint())
		return
	}

	ia1.signVal = newSignVal

	carry := 0
	n1 := 0
	n2 := 0
	n3 := 0

	for j := ia1.intAryLen - 1; j >= 0; j-- {

		if doReverseNums {

			n2 = int(ia1.intAry[j])
			n1 = int(ia2.intAry[j])

		} else {
			n1 = int(ia1.intAry[j])
			n2 = int(ia2.intAry[j])

		}

		if doAdd {
			// doAdd == true
			// Do Addition

			n3 = n1 + n2 + carry

			if n3 > 9 {
				n3 = n1 + n2 + carry - 10
				carry = 1

			} else {
				carry = 0
			}

		} else {
			// doAdd == false
			// Do Subtraction
			n3 = n1 - n2 - carry

			if n3 < 0 {
				n3 = n1 + 10 - n2 - carry
				carry = 1
			} else {
				carry = 0
			}
		}

		ia1.intAry[j] = uint8(n3)

	}

	if carry > 0 {
		ia1.intAry = append([]uint8{1}, ia1.intAry...)
		ia1.intAryLen++
	}

	if ia1.intAry[0] == 0 {
		ia1.SetSignificantDigitIdxs()
		ia1.intAry = ia1.intAry[ia1.firstDigitIdx:]
	}

	ia1.SetInternalFlags()

	return
}
