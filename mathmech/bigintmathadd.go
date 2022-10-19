package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

type BigIntMathAdd struct {
	Input  BigIntPair
	Result BigIntNum
}

// AddBigInts - Adds two *big.Int numbers. Each *big.Int number
// is passed to the method with an associated decimal place precision
// specification.
//
// The BigIntNum 'result' returned by this addition operation will contain
// USA default numeric separators (decimal separator, thousands separator and
// currency symbol).
//
func (bAdd BigIntMathAdd) AddBigInts(
	b1 *big.Int,
	precision1 uint,
	b2 *big.Int,
	precision2 uint) BigIntNum {

	// No error is possible because both precision parameters
	// are by definition, greater than or equal to zero.

	result, resultPrecision, _ :=
		BigIntMathAdd{}.BigIntAdd(
			b1,
			big.NewInt(0).SetUint64(uint64(precision1)),
			b2,
			big.NewInt(0).SetUint64(uint64(precision2)))

	biNum, _ := BigIntNum{}.NewBigIntPrecision(result, resultPrecision)

	return biNum

}

// AddBigIntNums - Adds two BigIntNums and returns the result in a new
// BigIntNum instance
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'b1'.
//
func (bAdd BigIntMathAdd) AddBigIntNums(b1, b2 BigIntNum) BigIntNum {

	bPair := BigIntPair{}.NewBigIntNum(b1, b2)

	return bAdd.AddPair(bPair)
}

// AddBigIntNumArray - Adds an Array of 'BigIntNum' types and returns the result
// as type 'BigIntNum'
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from the first element of the bNums array (bNums[0]).
//
func (bAdd BigIntMathAdd) AddBigIntNumArray(bNums []BigIntNum) BigIntNum {

	finalResult := BigIntNum{}.NewZero(0)
	lenBNums := len(bNums)

	if lenBNums == 0 {
		return finalResult
	}

	numSeps := bNums[0].GetNumericSeparatorsDto()

	for i := 0; i < lenBNums; i++ {

		if i == 0 {
			finalResult = bNums[i].CopyOut()
			continue
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bNums[i])

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// AddBigIntNumOutputToArray - The first input parameter to this method
// is a BigIntNum Type labeled, 'addend'.  The second element is an
// array of BigIntNum types labeled 'bNums'. The 'addend' is added to
// each element of the 'bNums' array with the result output to another
// array of BigIntNum types ([]BigIntNum) which is returned to the calling
// function.
//
// Example
// =======
// 										Multiplicands												Output
//  Addend   				    	Array														Array
//
//		3			+					bNums[0] = 2			=				  outputarray[0] =  5
//		3			+					bNums[1] = 3			=				  outputarray[1] =  6
//		3			+					bNums[2] = 4			=				  outputarray[2] =  7
//		3			+					bNums[3] = 5			=				  outputarray[3] =  8
//		3			+					bNums[4] = 6			=				  outputarray[4] =  9
//		3			+					bNums[5] = 9			=				  outputarray[5] = 12
//
// Each element of the the []BigIntNum 'result' array returned by this addition
// operation will contain numeric separators (decimal separator, thousands
// separator and currency symbol) which were copied from input parameter
// 'addend'.
//
func (bAdd BigIntMathAdd) AddBigIntNumOutputToArray(
	addend BigIntNum,
	bNums []BigIntNum) []BigIntNum {

	lenBNums := len(bNums)

	if lenBNums == 0 {
		return []BigIntNum{}
	}

	numSeps := addend.GetNumericSeparatorsDto()

	resultArray := make([]BigIntNum, lenBNums)

	for i := 0; i < lenBNums; i++ {

		bPair := BigIntPair{}.NewBigIntNum(addend, bNums[i])

		result := bAdd.addPairNoNumSeps(bPair)

		_ = result.SetNumericSeparatorsDto(numSeps)

		resultArray[i] = result.CopyOut()

	}

	return resultArray
}

// AddBigIntNumSeries - Adds a series of BigIntNum types and returns the total in a
// BigIntNum instance.
//
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from the first element in input series 'bNums'.
//
func (bAdd BigIntMathAdd) AddBigIntNumSeries(bNums ...BigIntNum) BigIntNum {

	finalResult := BigIntNum{}.New()

	numSeps := NumericSeparatorDto{}

	for i, bNum := range bNums {

		if i == 0 {
			finalResult = bNum.CopyOut()
			numSeps = finalResult.GetNumericSeparatorsDto()
			continue
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bNum)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// AddDecimal - Receives two Decimal instances and adds their numeric values.
//
// The result is returned as type BigIntNum.
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'dec1'.
//

func (bAdd BigIntMathAdd) AddDecimal(dec1, dec2 Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrDto() "

	// This method tests the validity of dec1 and dec2
	bPair, err := BigIntPair{}.NewDecimal(dec1, dec2)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewDecimal(dec1, dec2). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bAdd.AddPair(bPair)

	return finalResult, nil
}

// AddDecimalArray - Adds an array of 'Decimal' types and returns the combined total
// as an instance of Type, 'BigIntNum'.
//
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from the first element of the input (decs[0]).
//
func (bAdd BigIntMathAdd) AddDecimalArray(decs []Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddDecimalArray() "

	finalResult := BigIntNum{}.New()
	var err error

	lenDecs := len(decs)

	if lenDecs == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: decs array is Empty!")
	}

	numSeps := decs[0].GetNumericSeparatorsDto()

	for i := 0; i < lenDecs; i++ {

		if i == 0 {

			// This method tests the validity of decs[i]
			finalResult, err = BigIntNum{}.NewDecimal(decs[i])

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by BigIntNum{}.NewDecimal(decs[i]). "+
						" i='%v' decs[i].GetNumStr()='%v' Error='%v' ",
						i, decs[i].GetNumStr(), err.Error())
			}

			continue
		}

		// This method tests the validity of decs[i]
		bigINumNextAddend, err := BigIntNum{}.NewDecimal(decs[i])

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewDecimal(decs[i]) "+
					"Error='%v' ", err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINumNextAddend)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// AddDecimalOutputToArray - The first input parameter to this method
// is a Decimal Type labeled, 'addend'.  The second element is an
// array of Decimal types labeled 'decs'. The 'addend' is added to
// each element of the 'decs' array with the result output to another
// array of Decimal types which is returned to the calling function.
//
// Example
// =======
// 										    decs										 Output
//  Addend   				    	Array											Array
//
//		3			+					decs[0] = 2			=				  outputarray[0] =  5
//		3			+					decs[1] = 3			=				  outputarray[1] =  6
//		3			+					decs[2] = 4			=				  outputarray[2] =  7
//		3			+					decs[3] = 5			=				  outputarray[3] =  8
//		3			+					decs[4] = 6			=				  outputarray[4] =  9
//		3			+					decs[5] = 9			=				  outputarray[5] = 12
//
//
// Each element in the []Decimal array 'result' returned by this addition
// operation will contain numeric separators (decimal separator, thousands
// separator and currency symbol) which were copied from input parameter
// 'addend'.
//
func (bAdd BigIntMathAdd) AddDecimalOutputToArray(
	addend Decimal,
	decs []Decimal) ([]Decimal, error) {

	ePrefix := "BigIntMathAdd.AddDecimalOutputToArray() "

	lenDecs := len(decs)

	if lenDecs == 0 {
		return []Decimal{},
			errors.New(ePrefix + "Error: decs array is Empty!")
	}

	numSeps := addend.GetNumericSeparatorsDto()

	resultsArray := make([]Decimal, lenDecs)

	for i := 0; i < lenDecs; i++ {

		// This method tests the validity of addend and decs[i]
		bPair, err := BigIntPair{}.NewDecimal(addend, decs[i])

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewDecimal(addend, decs[i]). "+
					" i='%v' dec[i].GetNumStr()='%v' Error='%v' ", i, decs[i].GetNumStr(), err.Error())
		}

		result := bAdd.addPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					" index='%v' Error='%v' ", i, err.Error())
		}

		resultsArray[i], err = result.GetDecimal()

		if err != nil {
			return []Decimal{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetDecimal(). "+
					" i='%v' decs[i].GetNumStr()='%v' Error='%v' ",
					i, decs[i].GetNumStr(), err.Error())
		}

	}

	return resultsArray, nil
}

// AddDecimalSeries - Adds a series of 'Decimal' types and returns the combined total
// as an instance of Type, 'BigIntNum'.
//
// The BigIntNum 'result' returned by this addition operation will contain numeric
// separators (decimal separator, thousands separator and currency symbol) which
// were copied from input series element 'decs[0]'.
//
func (bAdd BigIntMathAdd) AddDecimalSeries(decs ...Decimal) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddDecimalSeries() "

	finalResult := BigIntNum{}.New()
	var err error

	if len(decs) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: decs series is Empty!")
	}

	numSeps := NumericSeparatorDto{}

	for i, dec := range decs {

		if i == 0 {

			// This method tests the validity of 'dec'
			finalResult, err = BigIntNum{}.NewDecimal(dec)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by BigIntNum{}.NewDecimal(dec). "+
						" i='%v' dec.GetNumStr()='%v' Error='%v' ",
						i, dec.GetNumStr(), err.Error())
			}

			numSeps = dec.GetNumericSeparatorsDto()

			continue
		}

		bigINumNextAddend, err := dec.GetBigIntNum()

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by dec.GetBigIntNum(). "+
					" i='%v' dec='%v' Error='%v' ", i, dec.GetNumStr(), err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bigINumNextAddend)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())

	}

	return finalResult, nil
}

// AddIntAry - Receives two IntAry instances and adds their numeric values.
//
// The result is returned as type BigIntNum.
//
// The BigIntNum 'result' returned by this addition operation will contain numeric
// separators (decimal separator, thousands separator and currency symbol) which
// were copied from input parameter 'ia1'.
//
func (bAdd BigIntMathAdd) AddIntAry(ia1, ia2 IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrDto() "

	// This method will test the validity of ia1 and ia2
	bPair, err := BigIntPair{}.NewIntAry(ia1, ia2)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewIntAry(ia1, ia2). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bAdd.AddPair(bPair)

	return finalResult, nil
}

// AddIntAryArray - Receives an array of IntAry objects and totals their numeric values.
// The total numeric value is returned in a BigIntNum instance.
//
// The BigIntNum 'result' returned by this addition operation will contain numeric
// separators (decimal separator, thousands separator and currency symbol) which
// were copied from the first element of the input array, iarys[0].
//
func (bAdd BigIntMathAdd) AddIntAryArray(iarys []IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddIntAryArray() "

	finalResult := BigIntNum{}.New()

	lenIaArray := len(iarys)

	if lenIaArray == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: iarys array is Empty!")
	}

	numSeps := iarys[0].GetNumericSeparatorsDto()

	var err error

	for i := 0; i < lenIaArray; i++ {

		if i == 0 {

			// This method will test the validity of iarys[i]
			finalResult, err = BigIntNum{}.NewIntAry(iarys[i])

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(iarys[i]). "+
						" i='%v' NumStr='%v' Error='%v' ", i, iarys[i].GetNumStr(), err.Error())
			}

			continue
		}

		// This method will test the validity of iarys[i]
		bINumNextAddend, err := BigIntNum{}.NewIntAry(iarys[i])

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(iarys[%v]). "+
					"Error='%v'. ", i, err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bINumNextAddend)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// AddIntAryOutputToArray - The first input parameter to this method
// is an IntAry Type labeled, 'addend'.  The second element is an
// array of IntAry types labeled 'iarys'. The 'addend' is added to
// each element of the 'iarys' array with the result output to another
// array of IntAry types which is returned to the calling function.
//
// Example
// =======
// 										    decs										 Output
//  Addend   				    	Array											Array
//
//		3			+					iarys[0] = 2			=				  outputarray[0] =  5
//		3			+					iarys[1] = 3			=				  outputarray[1] =  6
//		3			+					iarys[2] = 4			=				  outputarray[2] =  7
//		3			+					iarys[3] = 5			=				  outputarray[3] =  8
//		3			+					iarys[4] = 6			=				  outputarray[4] =  9
//		3			+					iarys[5] = 9			=				  outputarray[5] = 12
//
// Each element of the []IntAry returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency symbol)
// which were copied from input parameter 'addend'.
//
func (bAdd BigIntMathAdd) AddIntAryOutputToArray(
	addend IntAry,
	iarys []IntAry) ([]IntAry, error) {

	ePrefix := "BigIntMathAdd.AddIntAryOutputToArray() "

	// This method will test the validity of 'addend'
	bINumAddend, err := BigIntNum{}.NewIntAry(addend)

	if err != nil {
		return []IntAry{},
			fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(addend). "+
				"Error='%v' ", err.Error())
	}

	lenIaArray := len(iarys)

	if lenIaArray == 0 {
		return []IntAry{},
			errors.New(ePrefix + "Error: iarys array is Empty!")
	}

	numSeps := addend.GetNumericSeparatorsDto()

	resultsArray := make([]IntAry, lenIaArray)

	for i := 0; i < lenIaArray; i++ {

		// This method tests the validity of iarys[i]
		bINumNextAddend, err := BigIntNum{}.NewIntAry(iarys[i])

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(iarys[%v]). "+
					"Error='%v' ", i, err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(bINumAddend, bINumNextAddend)

		result := bAdd.addPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					"index='%v' Error='%v' ", i, err.Error())
		}

		resultsArray[i], err = result.GetIntAry()

		if err != nil {
			return []IntAry{},
				fmt.Errorf(ePrefix+
					"Error returned by result.GetIntAry(). "+
					"index='%v', result='%v', Error='%v' ",
					i, result.GetNumStr(), err.Error())
		}
	}

	return resultsArray, nil
}

// AddIntArySeries - Adds a series of IntAry objects and returns the total in a
// BigIntNum.
//
// The BigIntNum result of this addition operation will contain numeric separators
// (decimal separator, thousands separator and currency symbol) which were copied
// from the first element input series 'iarys'.
//
func (bAdd BigIntMathAdd) AddIntArySeries(iarys ...IntAry) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddIntArySeries() "

	finalResult := BigIntNum{}.New()

	if len(iarys) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'iarys' series is Empty!")
	}

	var err error
	numSeps := NumericSeparatorDto{}

	for i, ia := range iarys {

		if i == 0 {

			numSeps = ia.GetNumericSeparatorsDto()

			// This method tests the validity of 'ia'
			finalResult, err = BigIntNum{}.NewIntAry(ia)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewIntAry(ia). "+
						" i='%v' NumStr='%v' Error='%v' ", i, ia.GetNumStr(), err.Error())
			}

			continue
		}

		bINumNextAddend, err := BigIntNum{}.NewIntAry(ia)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewIntAry(ia). "+
					"index='%v' Error='%v' ", i, err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bINumNextAddend)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// AddINumMgr - Receives two objects which implement the INumMgr Interface
// and adds their numeric values.
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry.
//
// The result is returned as an instance of Type, 'BigIntNum'.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) which were copied from the input parameter 'num1'.
//
func (bAdd BigIntMathAdd) AddINumMgr(num1, num2 INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddINumMgr() "

	// This method will test the validity of num1 and num2
	bPair, err := BigIntPair{}.NewINumMgr(num1, num2)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntPair{}.NewINumMgr(num1, num2). "+
				"num1.GetNumStr()='%v', num2.GetNumStr()='%v' Error='%v' ",
				num1.GetNumStr(), num2.GetNumStr(), err.Error())
	}

	finalResult := bAdd.AddPair(bPair)

	return finalResult, nil
}

// AddINumMgrArray - Adds an array of objects which implement the 'INumMgr'
// interface. The combined total of the numeric values from these objects
// is returned as an instance of Type, 'BigIntNum'.
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry. This allows the user to mix different types in
// a single array and add their numeric values.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) which were copied from the first element of the input parameter
// array, nums (nums[0]).
//
func (bAdd BigIntMathAdd) AddINumMgrArray(nums []INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddINumMgrArray() "

	finalResult := BigIntNum{}.New()
	var err error

	lenDecs := len(nums)

	if lenDecs == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'nums' array is Empty!")
	}

	numSeps := nums[0].GetNumericSeparatorsDto()

	for i := 0; i < lenDecs; i++ {

		if i == 0 {
			// This method will test the validity of nums[i]
			finalResult, err = BigIntNum{}.NewINumMgr(nums[i])

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by BigIntNum{}.NewINumMgr(nums[i]). "+
						" i='%v' nums[i].GetNumStr()='%v' Error='%v' ",
						i, nums[i].GetNumStr(), err.Error())
			}

			continue
		}

		bPair, err := BigIntPair{}.NewINumMgr(&finalResult, nums[i])

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(&finalResult.Result, &nums[i]). "+
					" i='%v' nums[i].GetNumStr()='%v' Error='%v' ", i, nums[i].GetNumStr(), err.Error())
		}

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return finalResult, nil
}

// AddINumMgrOutputToArray - The first input parameter to this method
// is an object that implements the  INumMgr interface and is labeled,
// 'addend'.  The second element is an array of objects which implement
// the INumMgr interface. This array is labeled labeled, 'numMgrs'. The
// 'addend' is added to each element of the 'numMgrs' array with the result
// output to another INumMgr array which is returned to the calling function.
//
// Example
// =======
// 										    numMgrs										 Output
//  Addend   				    	Array											Array
//
//		3			+					numMgrs[0] = 2			=				  outputarray[0] =  5
//		3			+					numMgrs[1] = 3			=				  outputarray[1] =  6
//		3			+					numMgrs[2] = 4			=				  outputarray[2] =  7
//		3			+					numMgrs[3] = 5			=				  outputarray[3] =  8
//		3			+					numMgrs[4] = 6			=				  outputarray[4] =  9
//		3			+					numMgrs[5] = 9			=				  outputarray[5] = 12
//
// Each element of the returned []INumMgr array will contain numeric separators
// (decimal separator, thousands separator and currency symbol) which were
// copied from the input parameter 'addend'.
//
func (bAdd BigIntMathAdd) AddINumMgrOutputToArray(
	addend INumMgr,
	numMgrs []INumMgr) ([]INumMgr, error) {

	ePrefix := "BigIntMathAdd.AddINumMgrOutputToArray() "

	lenDecs := len(numMgrs)

	if lenDecs == 0 {
		return []INumMgr{},
			errors.New(ePrefix + "Error: 'numMgrs' array is Empty!")
	}

	numSeps := addend.GetNumericSeparatorsDto()

	resultsArray := make([]INumMgr, lenDecs)

	for i := 0; i < lenDecs; i++ {

		// This method will test the validity of numMgrs[i]
		bPair, err := BigIntPair{}.NewINumMgr(addend, numMgrs[i])

		if err != nil {
			return []INumMgr{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(addend, numMgrs[i]). "+
					" i='%v' numMgrs[i].GetNumStr()='%v' Error='%v' ", i, numMgrs[i].GetNumStr(), err.Error())
		}

		bIntNum := bAdd.addPairNoNumSeps(bPair)

		err = bIntNum.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []INumMgr{},
				fmt.Errorf(ePrefix+
					"Error returned by bIntNum.SetNumericSeparatorsDto(numSeps). "+
					"Error='%v' ", err.Error())
		}

		resultsArray[i] = &bIntNum
	}

	return resultsArray, nil
}

// AddINumMgrSeries - Adds a series of objects which implement the 'INumMgr'
// interface. The combined total of the numeric values from these objects
// is returned as an instance of Type, 'BigIntNum'.
//
// The INumMgr interface is implemented by types, BigIntNum, Decimal,
// NumStrDto and IntAry. This allows the user to mix different types in
// a single array and add their numeric values.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) which were copied from the first element of the input parameter
// 'nums'.
//
func (bAdd BigIntMathAdd) AddINumMgrSeries(nums ...INumMgr) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddINumMgrSeries() "

	finalResult := BigIntNum{}.New()

	if len(nums) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'nums' series is Empty!")
	}

	var err error

	numSeps := NumericSeparatorDto{}

	for i, num := range nums {

		if i == 0 {

			numSeps = num.GetNumericSeparatorsDto()

			// This method will test the validity of 'num'
			finalResult, err = BigIntNum{}.NewINumMgr(num)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by BigIntNum{}.NewINumMgr(num). "+
						" i='%v' num.GetNumStr()='%v' Error='%v' ",
						i, num.GetNumStr(), err.Error())
			}

			continue
		}

		bPair, err := BigIntPair{}.NewINumMgr(&finalResult, num)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntPair{}.NewINumMgr(&finalResult.Result, &num). "+
					" i='%v' num.GetNumStr()='%v' Error='%v' ",
					i, num.GetNumStr(), err.Error())
		}

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return finalResult, nil
}

// AddNumStr - Receives two number strings and adds their numeric values.
//
// The first two input parameters, 'n1NumStr' and 'n2NumStr' must be formatted
// as strings of numeric digits, or number strings. Number strings may have
// a leading minus sign ('-') to indicate the numeric sign value. In addition,
// the string of numeric digits may include a delimiting decimal separator
// to identify fractional digits. The number strings are parsed based on the
// decimal separator character specified by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to
// parse the number strings 'n1NumStr' and 'n2NumStr'. 'numSeps' represents the
// applicable decimal separator, thousands separator and currency symbol. In
// addition, 'numSeps' is also used in configuring the return value for this
// addition operation.
//
// The result is returned as type BigIntNum.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) specified by input parameter, 'numSeps'.
//
func (bAdd BigIntMathAdd) AddNumStr(n1NumStr, n2NumStr string, numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStr() "

	numSeps.SetDefaultsIfEmpty()

	bPair, err := BigIntPair{}.NewNumStrWithNumSeps(n1NumStr, n2NumStr, numSeps)

	if err != nil {
		return BigIntNum{},
			fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewNumStr(n1NumStr, n2NumStr). "+
				"Error='%v' ", err.Error())
	}

	return bAdd.AddPair(bPair), nil
}

// AddNumStrArray - Adds a series of number strings and returns the combined total
// as an instance of Type 'BigIntNum'.
//
// All the elements of the 'numStrs' array must be formatted as strings of numeric
// digits or number strings. Number strings may have a leading minus sign ('-')
// to indicate the numeric sign value.  In addition, the string of numeric digits
// may include a delimiting decimal separator to identify fractional digits. The
// number strings are parsed based on the decimal separator character specified
// by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to parse
// the number strings contained in the 'numStrs' array. Input parameter, 'numSeps',
// represents the applicable decimal separator, thousands separator and currency
// symbol. 'numSeps' is also used in configuring the BigIntNum return value for
// this addition operation.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) as specified by input parameter 'numSeps'.
//
func (bAdd BigIntMathAdd) AddNumStrArray(
	numStrs []string, numSeps NumericSeparatorDto) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrArray() "

	numSeps.SetDefaultsIfEmpty()

	finalResult := BigIntNum{}.New()
	var err error

	lenNumStrs := len(numStrs)

	if lenNumStrs == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'numStrs' array is Empty!")
	}

	for i := 0; i < lenNumStrs; i++ {

		if i == 0 {

			finalResult, err = BigIntNum{}.NewNumStrWithNumSeps(numStrs[i], numSeps)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by NewNumStrWithNumSeps(numStrs[i], numSeps). "+
						" i='%v' NumStr='%v' Error='%v' ", i, numStrs[i], err.Error())
			}

			continue
		}

		b2Num, err := BigIntNum{}.NewNumStrWithNumSeps(numStrs[i], numSeps)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(numStrs[i], numSeps). "+
					" i='%v' NumStr='%v' Error='%v' ", i, numStrs[i], err.Error())
		}

		finalResult = bAdd.AddBigIntNums(finalResult, b2Num)
	}

	return finalResult, nil
}

// AddNumStrOutputToArray - The first input parameter to this method
// is a string Type labeled, 'addend'.  The second element is an
// array of string types labeled 'numStrs'. The 'addend' is added to
// each element of the 'numStrs' array with the result output to another
// array of string types. This output array of strings is then returned
// to the calling function.
//
// Input parameters 'addend' and all elements of the 'numStrs' array must
// be formatted as strings of numeric digits or number strings. Number strings
// may have a leading minus sign ('-') to indicate the numeric sign value. In
// addition, the string of numeric digits may include a delimiting decimal
// separator to identify fractional digits. The number strings are parsed based
// on the decimal separator character specified by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to parse
// the number strings contained in 'addend' and the 'numStrs' array . Input
// parameter, 'numSeps', represents the applicable decimal separator, thousands
// separator and currency symbol. 'numSeps' is also used in configuring the output
// string array returned by this addition operation.
//
// Example
// =======
// 										    numStrs										 Output
//  Addend   				    	Array											Array
//
//		3			+					numStrs[0] = 2			=				  outputarray[0] =  5
//		3			+					numStrs[1] = 3			=				  outputarray[1] =  6
//		3			+					numStrs[2] = 4			=				  outputarray[2] =  7
//		3			+					numStrs[3] = 5			=				  outputarray[3] =  8
//		3			+					numStrs[4] = 6			=				  outputarray[4] =  9
//		3			+					numStrs[5] = 9			=				  outputarray[5] = 12
//
//
func (bAdd BigIntMathAdd) AddNumStrOutputToArray(
	addend string,
	numStrs []string,
	numSeps NumericSeparatorDto) ([]string, error) {

	ePrefix := "BigIntMathAdd.AddNumStrOutputToArray() "

	numSeps.SetDefaultsIfEmpty()

	lenNumStrs := len(numStrs)

	if lenNumStrs == 0 {
		return []string{},
			errors.New(ePrefix + "Error: 'numStrs' array is Empty!")
	}

	bINumAddend, err := BigIntNum{}.NewNumStrWithNumSeps(addend, numSeps)

	if err != nil {
		return []string{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(addend, numSeps). "+
				" addend='%v' Error='%v' ", addend, err.Error())
	}

	resultsArray := make([]string, lenNumStrs)

	for i := 0; i < lenNumStrs; i++ {

		b2Num, err := BigIntNum{}.NewNumStrWithNumSeps(numStrs[i], numSeps)

		if err != nil {
			return []string{},
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrWithNumSeps(numStrs[i], numSeps). "+
					" i='%v' NumStr='%v' Error='%v' ", i, numStrs[i], err.Error())
		}

		result := bAdd.AddPair(BigIntPair{}.NewBigIntNum(bINumAddend, b2Num))

		resultsArray[i] = result.GetNumStr()

	}

	return resultsArray, nil
}

// AddNumStrSeries - Adds a series of number strings and returns
// the combined total as an instance of Type, 'BigIntNum'.
//
// The second Input Parameter ,'numStrs', is series strings of numeric digits,
// or number strings. Number strings may have a leading minus sign ('-')
// to indicate the numeric sign value. In addition, the string of numeric digits
// may include a delimiting decimal separator to identify fractional digits. The
// number strings are parsed based on the decimal separator character specified
// by input parameter 'numSeps'.
//
// Input parameter 'numSeps' is a type NumericSeparatorDto and is used to parse
// the number strings contained in the 'numStrs' series . Input parameter, 'numSeps',
// represents the applicable decimal separator, thousands separator and currency symbol.
// 'numSeps' is also used in configuring the output string array returned by this addition
// operation.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) specified by input parameter 'numSeps'.
//
func (bAdd BigIntMathAdd) AddNumStrSeries(
	numSeps NumericSeparatorDto, numStrs ...string) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrSeries() "

	numSeps.SetDefaultsIfEmpty()

	finalResult := BigIntNum{}.NewZero(0)

	if len(numStrs) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'numStrs' series is Empty!")
	}

	var err error

	for i, numStr := range numStrs {

		if i == 0 {

			finalResult, err = BigIntNum{}.NewNumStrWithNumSeps(numStr, numSeps)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+
						"Error returned by BigIntNum{}.NewNumStrWithNumSeps(numStr, numSeps). "+
						" i='%v' NumStr='%v' Error='%v' ", i, numStr, err.Error())
			}

			continue
		}

		b2Num, err := BigIntNum{}.NewNumStrWithNumSeps(numStr, numSeps)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStrWithNumSeps(numStr, numSeps). "+
					" i='%v' NumStr='%v' Error='%v' ", i, numStr, err.Error())
		}

		finalResult = bAdd.AddBigIntNums(finalResult, b2Num)

	}

	finalResult.SetNumericSeparatorsToDefaultIfEmpty()

	return finalResult, nil
}

// AddNumStrDto - Receives two NumStrDto instances and adds their numeric values.
//
// The result is returned as type BigIntNum.
//
// The returned BigIntNum result of this addition operation will contain
// the numeric separators (decimal separator, thousands separator and currency
// symbol) copied from input parameter 'n1Dto'.
//
func (bAdd BigIntMathAdd) AddNumStrDto(n1Dto, n2Dto NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrDto() "

	// This method will test the validity of n1Dto and n2Dto
	bPair, err := BigIntPair{}.NewNumStrDto(n1Dto, n2Dto)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewNumStrDto(n1Dto, n2Dto). "+
				"Error='%v' ", err.Error())
	}

	finalResult := bAdd.AddPair(bPair)

	return finalResult, nil
}

// AddNumStrDtoArray - Adds an array of 'NumStrDto' types and returns the combined total
// as an instance of Type, 'BigIntNum'.
//
// The returned BigIntNum result of this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) copied from the first element of the input array, nDtos[0].
//
func (bAdd BigIntMathAdd) AddNumStrDtoArray(nDtos []NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddDecimalArray() "
	finalResult := BigIntNum{}.New()
	var err error

	lenDecs := len(nDtos)

	if lenDecs == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'nDtos' array is Empty!")
	}

	numSeps := NumericSeparatorDto{}

	for i := 0; i < lenDecs; i++ {

		if i == 0 {

			// This method will test the validity of 'nDtos[i]'
			finalResult, err = BigIntNum{}.NewNumStrDto(nDtos[i])

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStrDto(nDtos[i]). "+
						" i='%v' nDtos[i].GetNumStr()='%v' Error='%v' ", i, nDtos[i].GetNumStr(), err.Error())
			}

			numSeps = nDtos[0].GetNumericSeparatorsDto()

			continue
		}

		bPair, err := BigIntPair{}.NewINumMgr(&finalResult, &nDtos[i])

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewINumMgr(&finalResult.Result, &nDtos[i]). "+
					" i='%v' dec[i].GetNumStr()='%v' Error='%v' ", i, nDtos[i].GetNumStr(), err.Error())
		}

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' \n", err.Error())
	}

	return finalResult, nil
}

// AddNumStrDtoOutputToArray - The first input parameter to this method
// is a NumStrDto Type labeled, 'addend'.  The second element is an
// array of NumStrDto types labeled 'nDtos'. The 'addend' is added to
// each element of the 'nDtos' array with the result output to another
// array of NumStrDto types which is returned to the calling function.
//
// Example
// =======
// 										    nDtos										 Output
//  Addend   				    	Array											Array
//
//		3			+					nDtos[0] = 2			=				  outputarray[0] =  5
//		3			+					nDtos[1] = 3			=				  outputarray[1] =  6
//		3			+					nDtos[2] = 4			=				  outputarray[2] =  7
//		3			+					nDtos[3] = 5			=				  outputarray[3] =  8
//		3			+					nDtos[4] = 6			=				  outputarray[4] =  9
//		3			+					nDtos[5] = 9			=				  outputarray[5] = 12
//
// Each element in the returned array of []NumStrDto resulting of this addition
// operation will contain numeric separators (decimal separator, thousands separator
// and currency symbol) copied from input parameter 'addend'.
//
func (bAdd BigIntMathAdd) AddNumStrDtoOutputToArray(
	addend NumStrDto,
	nDtos []NumStrDto) ([]NumStrDto, error) {

	ePrefix := "BigIntMathAdd.AddNumStrDtoOutputToArray() "

	lenDecs := len(nDtos)

	if lenDecs == 0 {
		return []NumStrDto{},
			errors.New(ePrefix + "Error: 'nDtos' array is Empty!")
	}

	numSeps := addend.GetNumericSeparatorsDto()

	resultsArray := make([]NumStrDto, lenDecs)

	for i := 0; i < lenDecs; i++ {

		err := nDtos[i].IsValid(ePrefix +
			fmt.Sprintf("nDtos[%v] INVALID! ", i))

		if err != nil {
			return []NumStrDto{}, err
		}

		// This method tests the validity of addend and nDtos[i]
		bPair, err := BigIntPair{}.NewNumStrDto(addend, nDtos[i])

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+"Error returned by BigIntPair{}.NewNumStrDto(addend, nDtos[i]). "+
					" i='%v' dec[i].GetNumStr()='%v' Error='%v' ", i, nDtos[i].GetNumStr(), err.Error())
		}

		result := bAdd.addPairNoNumSeps(bPair)

		err = result.SetNumericSeparatorsDto(numSeps)

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by result.SetNumericSeparatorsDto(numSeps). "+
					"Error='%v' \n", err.Error())
		}

		resultsArray[i], err = result.GetNumStrDto()

		if err != nil {
			return []NumStrDto{},
				fmt.Errorf(ePrefix+
					"Error returned by result.Result.GetNumStrDto(). "+
					" i='%v' nDtos[i].GetNumStr()='%v' Error='%v' ",
					i, nDtos[i].GetNumStr(), err.Error())
		}

	}

	return resultsArray, nil
}

// AddNumStrDtoSeries - Adds a series of 'NumStrDto' types and returns the combined total
// as an instance of Type, 'BigIntNum'.
//
// The returned BigIntNum resulting of this addition operation will contain numeric
// separators (decimal separator, thousands separator and currency symbol) copied
// from the first element of the input series 'nDtos'.
//
func (bAdd BigIntMathAdd) AddNumStrDtoSeries(nDtos ...NumStrDto) (BigIntNum, error) {

	ePrefix := "BigIntMathAdd.AddNumStrDtoSeries() "

	finalResult := BigIntNum{}.New()

	if len(nDtos) == 0 {
		return finalResult,
			errors.New(ePrefix + "Error: 'nDtos' series is Empty!")
	}

	var err error
	numSeps := NumericSeparatorDto{}

	for i, nDto := range nDtos {

		if i == 0 {

			// This method will test the validity of 'nDto'
			finalResult, err = BigIntNum{}.NewNumStrDto(nDto)

			if err != nil {
				return BigIntNum{}.New(),
					fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStrDto(nDto). "+
						" i='%v' NumStr='%v' Error='%v' ", i, nDto.GetNumStr(), err.Error())
			}

			numSeps = nDto.GetNumericSeparatorsDto()

			continue
		}

		// This method will test the validity of 'nDto'
		bINumNextAddend, err := BigIntNum{}.NewNumStrDto(nDto)

		if err != nil {
			return BigIntNum{}.New(),
				fmt.Errorf(ePrefix+
					"Error returned by BigIntNum{}.NewNumStrDto(nDto). "+
					" i='%v' Error='%v' ", i, err.Error())
		}

		bPair := BigIntPair{}.NewBigIntNum(finalResult, bINumNextAddend)

		finalResult = bAdd.addPairNoNumSeps(bPair)
	}

	err = finalResult.SetNumericSeparatorsDto(numSeps)

	if err != nil {
		return BigIntNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by finalResult.SetNumericSeparatorsDto(numSeps). "+
				"Error='%v' ", err.Error())
	}

	return finalResult, nil
}

// AddPair - Receives a BigIntPair instance and proceeds to add b1.BigIntNum
// to b2.BigIntNum.
//
// The result is returned as type BigIntNum.
//
// The BigIntNum 'result' returned by this addition operation will contain
// numeric separators (decimal separator, thousands separator and currency
// symbol) copied from b1.BigIntNum.
//
func (bAdd BigIntMathAdd) AddPair(bPair BigIntPair) BigIntNum {

	numSeps := bPair.Big1.GetNumericSeparatorsDto()

	finalResult := bAdd.addPairNoNumSeps(bPair)

	_ = finalResult.SetNumericSeparatorsDto(numSeps)

	return finalResult
}

// BigIntAdd - Adds two fixed length floating point numbers and generates
// the sum or total resulting from that addition. The two numbers added
// together are configured as pairs of *big.Int integer numbers and precision
// specifications. Each integer precision pair is used to define a fixed
// length floating point number.
//
// Examples:
// =========
//
// In the addition operation:
// 						b1 + b2 = total or sum
//
// This method provides for the addition of fixed length
// floating point values by means of integer and precision
// specification pairs.
//
// As an example, consider the following addition operation
// 						752.314 + 21.67894 = 773.99294 = total
//              b1    +     b2   =   total
//
// In this case 'b1', 'b2' and 'total' would be configured as integer
// precision pairs:
//									b1 							= 752314
//                  b1Precision			= 3
//                  b2 							= 2167894
//                  b2Precision 		= 5
//
//                  total						= 77399294
//                  totalPrecision 	= 5
//
// In this way, the method uses integer, precision pairs to define fixed
// length floating point numbers.
//
// Input Parameters
// ================
//
//	b1 					*big.Int	- The the first number which will be added to 'b2' to
//                    			generate a total.
//
//	b1Precision	*big.Int	- Specifies the precision for input parameter 'b1'.
// 													Precision defines the number of fractional digits
//													after the decimal place. 'b1Precision' must be equal
//                          to or greater than zero.
//
//	b2 					*big.Int	- The second number which is added to 'b1' in order to
//                 					generate a total.
//
//	b2Precision	*big.Int  - The 'b2' precision or the number of fractional digits
//													after the decimal place. 'b2Precision' must be equal
//                          to or greater than zero.
//
// Return Values
// =============
//
// total 					*big.Int		- The sum or total of 'b1' and 'b2' input values.
//
// totalPrecision *big.Int   	- The 'total' precision or the number of fractional
// 															digits after the decimal place.
//
// err						error				- If input parameters 'b1Precision' or 'b2Precision'
// 															are less than zero, an error will be returned.
//
// Taken together, 'total' and 'totalPrecision' can define a fixed length floating point number.
//
func (bAdd BigIntMathAdd) BigIntAdd(
	b1,
	b1Precision,
	b2,
	b2Precision *big.Int) (total *big.Int, totalPrecision *big.Int, err error) {

	ePrefix := "BigIntMathAdd.BigIntAdd() "

	total = big.NewInt(0)
	totalPrecision = big.NewInt(0)
	err = nil

	if b1 == nil {
		b1 = big.NewInt(0)
	}

	if b2 == nil {
		b2 = big.NewInt(0)
	}

	bigZero := big.NewInt(0)

	if b1Precision.Cmp(bigZero) == -1 {

		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'b1Precision' is LESS THAN ZERO! "+
			"b1Precision='%v' ", b1Precision.Text(10))

		return total, totalPrecision, err
	}

	if b2Precision.Cmp(bigZero) == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Input parameter 'b2Precision' is LESS THAN ZERO! "+
			"b2Precision='%v' ", b2Precision.Text(10))

		return total, totalPrecision, err
	}

	if b1.Cmp(bigZero) == 0 &&
		b2.Cmp(bigZero) == 0 {
		total = big.NewInt(0)
		totalPrecision = big.NewInt(0)
		return total, totalPrecision, nil
	}

	bigTen := big.NewInt(10)
	delta := big.NewInt(0)
	scale := big.NewInt(0)

	if b1Precision.Cmp(b2Precision) == 0 {
		total = big.NewInt(0).Add(b1, b2)
		totalPrecision = big.NewInt(0).Set(b1Precision)

	} else if b1Precision.Cmp(b2Precision) == 1 {
		// b1Precision > b2Precision
		delta = big.NewInt(0).Sub(b1Precision, b2Precision)
		scale = big.NewInt(0).Exp(bigTen, delta, nil)

		b2ToScale := big.NewInt(0).Mul(b2, scale)

		total = big.NewInt(0).Add(b1, b2ToScale)
		totalPrecision = big.NewInt(0).Set(b1Precision)

	} else {
		// b2Precision must be GREATER than b1Precision
		delta = big.NewInt(0).Sub(b2Precision, b1Precision)

		scale = big.NewInt(0).Exp(bigTen, delta, nil)

		b1ToScale := big.NewInt(0).Mul(b1, scale)

		total = big.NewInt(0).Add(b1ToScale, b2)

		totalPrecision = big.NewInt(0).Set(b2Precision)

	}

	if total.Cmp(bigZero) == 0 {
		totalPrecision = big.NewInt(0)
	}

	// Delete trailing fractional zeros
	if totalPrecision.Cmp(bigZero) == 1 {
		//totalPrecision > 0
		scrap := big.NewInt(0)
		biBase10 := big.NewInt(10)
		biBaseZero := big.NewInt(0)
		newTotal, mod10 := big.NewInt(0).QuoRem(total, biBase10, scrap)
		bigOne := big.NewInt(1)
		for mod10.Cmp(biBaseZero) == 0 && totalPrecision.Cmp(bigZero) == 1 {
			total.Set(newTotal)
			totalPrecision.Sub(totalPrecision, bigOne)
			newTotal, mod10 = big.NewInt(0).QuoRem(total, biBase10, scrap)
		}
	}

	return total, totalPrecision, nil
}

// FixedDecimalAdd - Performs an addition operation using two BigIntFixedDecimal
// types. The addition result or total is also returned as a BigIntFixedDecimal
// type.
//
// Examples:
// =========
//
// In the addition operation:
//
// 						b1 + b2 = total or sum
//
// For this method 'b1', 'b2' and 'total' are all configured as BigIntFixedDecimal
// types.
//
// The BigIntFixedDecimal type is used to defined fixed length floating point
// numbers and is defined as follows:
//
// type BigIntFixedDecimal struct {
//
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
// 	To represent the floating point number 52.459	a BigIntDecimal Structure
// 	would be configured as follows:
//
// 			BigIntFixedDecimal.integerNum	= 52459
// 			BigIntFixedDecimal.precision	= 3
//
// As an example consider the following addition operation:
// 						752.314 + 21.67894 = 773.99294 = total
//              b1    +     b2   =   total
//
// In this case 'b1', 'b2' and 'total' would be configured as BigIntDecimal
// types:
//									b1.integerNum			= 752314
//                  b1.precision			= 3
//                  b2.integerNum			= 2167894
//                  b2.precision 			= 5
//
//                  total.integerNum	= 77399294
//                  total.precision 	= 5
//
// In this way, the method uses BigIntFixedDecimal types to define fixed
// length floating point numbers.
//
// Input Parameters
// ================
//
//	b1 BigIntFixedDecimal	- The the first number which will be added to 'b2' to
//                    			generate a total.
//
//	b2 BigIntFixedDecimal	- The second number which is added to 'b1' in order to
//                 					generate a total.
//
// Return Values
// =============
//
// total BigIntFixedDecimal	- The sum or total of 'b1' and 'b2' input values.
//
func (bAdd BigIntMathAdd) FixedDecimalAdd(
	b1,
	b2 BigIntFixedDecimal) (total BigIntFixedDecimal) {

	total = BigIntFixedDecimal{}.NewZero(0)

	b1.IsValid()

	b2.IsValid()

	// No error is possible because by definition, both precision
	// values must be equal to or greater than zero.

	bIResult, bIPrecision, _ :=
		BigIntMathAdd{}.BigIntAdd(
			b1.GetInteger(),
			b1.GetPrecisionBigInt(),
			b2.GetInteger(),
			b2.GetPrecisionBigInt())

	total.SetNumericValue(bIResult, uint(bIPrecision.Uint64()))

	return total
}

// addPairNoNumSeps - Receives a BigIntPair and proceeds to add b1.BigIntNum to
// b2.BigIntNum.
//
// The result is returned as type BigIntNum.
//
// The BigIntNum 'result' returned by this subtraction operation will contain
// default numeric separators (decimal separator, thousands separator and
// currency symbol).
//
func (bAdd BigIntMathAdd) addPairNoNumSeps(bPair BigIntPair) BigIntNum {

	bPair.MakePrecisionsEqual()

	b3 := big.NewInt(0).Add(bPair.GetBig1BigInt(), bPair.GetBig2BigInt())

	bResult := BigIntNum{}.NewBigInt(b3, bPair.Big2.GetPrecisionUint())

	return bResult
}
