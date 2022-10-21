package mathops

type IntAryMathAdd struct {
	Input  IntAryPair
	Result IntAry
}

// AddManyArray - Adds the contents of the []IntAry, 'iaMany' to the
// the value of 'total'.
//
// The original 'total' numeric separators will remain unchanged.
// Numeric separators consist of the decimal separator, thousands
// separator and currency symbol.
//
func (iaAdd IntAryMathAdd) AddManyArray(total *IntAry, iaMany []IntAry) {

	lAry := len(iaMany)

	for i := 0; i < lAry; i++ {

		iaAdd.RunTotal(total, &iaMany[i])

	}

}

// AddMany - Adds multiple IntAry instances and returns the result in
// input parameter 'total'.
//
// The original 'total' numeric separators will remain unchanged.
// Numeric separators consist of the decimal separator, thousands
// separator and currency symbol.
//
func (iaAdd IntAryMathAdd) AddMany(total *IntAry, iaMany ...*IntAry) {

	for _, iAry := range iaMany {

		iaAdd.RunTotal(total, iAry)

	}

	return
}

// Adds two IntAry instances and returns the result as a new
// IntAry Instance.
//
// The return value 'IntAry' will contain the numeric separators
// (decimal separator, thousands separator and currency symbol)
// copied from input parameter 'ia1'.
//
func (iaAdd IntAryMathAdd) Add(ia1, ia2 *IntAry) IntAry {

	ia3 := ia1.CopyOut()

	iaAdd.RunTotal(&ia3, ia2)

	return ia3
}

// RunTotal- Adds to IntAry input parameters and returns the results
// in the first parameter, 'ia'.
//
// The returned addition 'result' in 'ia' will contain numeric separators
// (decimal separator, thousands separator and currency symbol) copied from
// the original 'ia' structure. In other words, the 'ia' numeric separators
// will remain unchanged.
//
func (iaAdd IntAryMathAdd) RunTotal(ia, ia2 *IntAry) {

	ia.SetEqualArrayLengths(ia2)

	if ia2.isZeroValue {
		return
	}

	numSeps := ia.GetNumericSeparatorsDto()

	compare := ia.CompareAbsoluteValues(ia2)

	newSignVal := ia.signVal
	doAdd := true
	isZeroResult := false
	doReverseNums := false

	if compare == 1 {
		// compare == + 1
		// Absolute Value: N1 > N2

		if ia.signVal == 1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = 1
		} else if ia.signVal == -1 && ia2.signVal == 1 {
			doAdd = false
			newSignVal = -1
		} else if ia.signVal == -1 && ia2.signVal == -1 {
			doAdd = true
			newSignVal = -1
		} else {
			// Must Be ia.signVal == 1 && ia2.signVal == -1
			doAdd = false
			newSignVal = 1
		}

	} else if compare == -1 {
		// Absolute Values: N2 > N1
		if ia.signVal == 1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = 1
		} else if ia.signVal == -1 && ia2.signVal == 1 {
			doAdd = false
			doReverseNums = true
			newSignVal = 1
		} else if ia.signVal == -1 && ia2.signVal == -1 {
			doAdd = true
			newSignVal = -1
		} else {
			// Must Be ia.signVal == 1 && ia2.signVal == -1
			doAdd = false
			doReverseNums = true
			newSignVal = -1
		}

	} else {
		// Must be compare == 0
		// Absolute Values: N1==N2
		if ia.signVal == 1 && ia2.signVal == 1 {
			doAdd = true
			newSignVal = 1
		} else if ia.signVal == -1 && ia2.signVal == 1 {
			doAdd = false
			newSignVal = 1
			isZeroResult = true
		} else if ia.signVal == -1 && ia2.signVal == -1 {
			doAdd = true
			newSignVal = -1
		} else {
			// Must Be ia.signVal == 1 && ia2.signVal == -1
			doAdd = false
			newSignVal = 1
			isZeroResult = true
		}

	}

	iaAdd.addToSubtract(ia, ia2, newSignVal, doAdd, isZeroResult, doReverseNums)

	ia.SetNumericSeparatorsDto(numSeps)

	return
}

// addToSubtract - Adds or subtracts two IntAry instances
// and returns the result in the first IntAry parameter.
func (iaAdd *IntAryMathAdd) addToSubtract(
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
