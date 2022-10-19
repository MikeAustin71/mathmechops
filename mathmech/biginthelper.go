package mathops

import (
	"fmt"
	"math/big"
)

// BigIntPair - contains a pair of 'BitIntNum' types. This structure
// is used to set up calculations involving *big.Int types.
type BigIntPair struct {
	Big1 BigIntNum

	Big1Compare       int // 	1 = Big1 > Big2; 0 = Big1 == Big2; -1 = Big1 < Big2
	Big1AbsCompare    int // 	1 = Big1 > Big2; 0 = Big1 == Big2; -1 = Big1 < Big2
	Precision1Compare int // 	1 = Big1Precision > Big2Precision;
	//  0 = Big1Precision == Big2Precision;
	// -1 = Big1Precision < Big2Precision
	Big2         BigIntNum
	MaxPrecision uint // Used to control output from complex
	//  calculations
}

// CopyIn - Copies the values provided by incoming BigIntPair
// parameter into the current BigIntPair instance.
func (bPair *BigIntPair) CopyIn(bd2 BigIntPair) {

	bPair.Big1 = bd2.Big1.CopyOut()
	bPair.Big2 = bd2.Big2.CopyOut()
	bPair.Big1Compare = bd2.Big1Compare
	bPair.Big1AbsCompare = bd2.Big1AbsCompare
	bPair.Precision1Compare = bd2.Precision1Compare

}

// CopyOut - Makes a deep copy of the current BigIntPair
// instance and returns it as a new BigIntPair object.
//
func (bPair *BigIntPair) CopyOut() BigIntPair {

	bd2 := BigIntPair{}.NewBigIntNum(bPair.Big1, bPair.Big2)

	return bd2
}

// Empty - Sets all data fields for the current BigIntPair instance
// to their uninitialized or zero states.
func (bPair *BigIntPair) Empty() {
	bPair.Big1.Empty()
	bPair.Big2.Empty()
	bPair.Big1Compare = 0
	bPair.Big1AbsCompare = 0
	bPair.Precision1Compare = 0

}

// GetBig1BigInt Returns the *big.Int value of data field
// bPair.Big1
func (bPair *BigIntPair) GetBig1BigInt() *big.Int {

	b1, err := bPair.Big1.GetBigInt()

	if err != nil {
		return big.NewInt(0)
	}

	return b1
}

// GetBig1BigInt Returns the *big.Int value of data field
// bPair.Big1
func (bPair *BigIntPair) GetBig2BigInt() *big.Int {

	b2, err := bPair.Big2.GetBigInt()

	if err != nil {
		return big.NewInt(0)
	}

	return b2
}

// MakePrecisionsEqual - Analyzes the two component BigIntNum's, b1 and b2,
// and then converts the one with the smallest precision to a new value
// equivalent in precision to the other BigIntNum. When completed, this
// method insures that both component BigIntNum's are both formatted to
// the largest precision.
//
func (bPair *BigIntPair) MakePrecisionsEqual() {

	if bPair.Big1.precision == bPair.Big2.precision {
		// Nothing to do. Precisions are equal.
		return
	}

	base10 := big.NewInt(10)

	if bPair.Big1.precision > bPair.Big2.precision {
		deltaPrecision := big.NewInt(int64(bPair.Big1.precision - bPair.Big2.precision))
		deltaPrecisionScale := big.NewInt(0).Exp(base10, deltaPrecision, nil)
		newB2Int := big.NewInt(0).Mul(bPair.Big2.bigInt, deltaPrecisionScale)
		newB2Num := BigIntNum{}.NewBigInt(newB2Int, bPair.Big1.precision)
		newBPair := BigIntPair{}.NewBigIntNum(bPair.Big1, newB2Num)
		bPair.CopyIn(newBPair)
		return

	}

	// Must be bPair.Big2.precision > bPair.Big1.precision
	deltaPrecision := big.NewInt(int64(bPair.Big2.precision - bPair.Big1.precision))
	deltaPrecisionScale := big.NewInt(0).Exp(base10, deltaPrecision, nil)
	newB1Int := big.NewInt(0).Mul(bPair.Big1.bigInt, deltaPrecisionScale)
	newB1Num := BigIntNum{}.NewBigInt(newB1Int, bPair.Big2.precision)
	newBPair := BigIntPair{}.NewBigIntNum(newB1Num, bPair.Big2)

	bPair.CopyIn(newBPair)

	return
}

// New - Creates an Empty BigIntPair instance. Both
// 'Big1' and 'Big2' are set to zero.  Both precision
// values are also set to zero.
func (bPair BigIntPair) New() BigIntPair {
	base1Zero := big.NewInt(0)
	base2Zero := big.NewInt(0)
	b2Pair := BigIntPair{}.NewBase(base1Zero, 0, base2Zero, 0)
	return b2Pair
}

// NewBase - Creates a BigIntPair instance using two sets of
// *big.Int integers and their associated precision specifications.
//
// Before using BigIntPair in a math operation, it may be necessary
// to specifically set 'BigIntPair.maxPrecision'
//
// The BigIntNum instances created by this method will contain default
// numeric separators (decimal separator, thousands separator and
// currency symbol)
//
func (bPair BigIntPair) NewBase(
	b1 *big.Int,
	b1Precision uint,
	b2 *big.Int,
	b2Precision uint) BigIntPair {

	b1BigIntNum := BigIntNum{}.NewBigInt(b1, b1Precision)
	b2BigIntNum := BigIntNum{}.NewBigInt(b2, b2Precision)

	return BigIntPair{}.NewBigIntNum(b1BigIntNum, b2BigIntNum)

}

// NewBigIntNum - Creates a new BigIntPair instance from input parameters
// consisting of two 'BigIntNum' types.
//
// Before using BigIntPair in a math operation, it may be necessary
// to specifically set 'BigIntPair.maxPrecision'
//
func (bPair BigIntPair) NewBigIntNum(b1, b2 BigIntNum) BigIntPair {

	bd2 := BigIntPair{}

	bd2.SetBigIntPair(b1, b2)

	return bd2
}

// NewIntAry - Creates a new BigIntPair instance from two
// Decimal instances passed as input parameters.
//
// The returned BigIntNum BigIntPair.Big1 will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'dec1'. BigIntPair.Big2 will contain numeric separators copied from 'dec2'.
//
func (bPair BigIntPair) NewDecimal(dec1, dec2 Decimal) (BigIntPair, error) {

	ePrefix := "BigIntPair.NewDecimal() "

	err := dec1.IsValid(ePrefix)

	if err != nil {
		return BigIntPair{}, err
	}

	err = dec2.IsValid(ePrefix)

	if err != nil {
		return BigIntPair{}, err
	}

	b1Num, err := dec1.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+
				"Error returned by dec1.GetBigIntNum(). Error='%v'. ", err.Error())
	}

	b2Num, err := dec2.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+
				"Error returned by dec2.GetBigIntNum(). Error='%v'. ", err.Error())
	}

	bd2 := BigIntPair{}

	bd2.SetBigIntPair(b1Num, b2Num)

	return bd2, nil
}

// NewIntAry - Creates a new BigIntPair instance from two
// IntAry instances passed as input parameters.
//
// This method will test the validity of input parameters ia1 and
// ia2.
//
// Be careful, IntAry's can accommodate very, very large numbers.
//
// The returned BigIntNum BigIntPair.Big1 will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'ia1'. BigIntPair.Big2 will contain numeric separators copied from 'ia2'.
//
func (bPair BigIntPair) NewIntAry(ia1, ia2 IntAry) (BigIntPair, error) {

	ePrefix := "BigIntPair.NewIntAry() "

	// Method NewIntAry will test the validity of ia1
	b1Num, err := ia1.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by ia1.GetBigIntNum(). "+
				"Error='%v' ", err.Error())
	}

	// Method NewIntAry will test the validity of ia2
	b2Num, err := ia2.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by ia2.GetBigIntNum(). "+
				"Error='%v' ", err.Error())
	}

	bd2 := BigIntPair{}

	bd2.SetBigIntPair(b1Num, b2Num)

	return bd2, nil
}

// NewINumMgr - Creates a new BigIntPair instance from two objects implementing the
// INumMgr interface.
//
// This method will test the validity of num1 and num2
//
// The returned BigIntNum BigIntPair.Big1 will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'num1'. BigIntPair.Big2 will contain numeric separators copied from 'num2'.
//
func (bPair BigIntPair) NewINumMgr(num1, num2 INumMgr) (BigIntPair, error) {

	ePrefix := "BigIntPair.NewINumMgr() "

	err := num1.IsValid(ePrefix + "'num1' INVALID! ")

	if err != nil {
		return BigIntPair{}, err
	}

	err = num2.IsValid(ePrefix + "'num2' INVALID! ")

	if err != nil {
		return BigIntPair{}, err
	}

	b1Num, err := num1.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by num1.GetBigIntNum(). "+
				"Error='%v'. ", err.Error())
	}

	b2Num, err := num2.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by num2.GetBigIntNum(). "+
				"Error='%v'. ", err.Error())
	}

	bPair2 := BigIntPair{}.NewBigIntNum(b1Num, b2Num)

	return bPair2, nil
}

// NewNumStr - Creates a new BigIntPair instance from two integerNum strings
// passed as input parameters.
//
// The returned BigIntNum's BigIntPair.Big1 and BigIntPair.Big2 will contain default
// numeric separators (decimal separator, thousands separator and currency symbol).
//
func (bPair BigIntPair) NewNumStr(n1NumStr, n2NumStr string) (BigIntPair, error) {
	ePrefix := "BigIntPair.NewNumStrDto() "
	b1Num, err := BigIntNum{}.NewNumStr(n1NumStr)

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStr(n1NumStr). "+
				"numStr='%v' Error='%v' ", n1NumStr, err.Error())
	}

	b2Num, err := BigIntNum{}.NewNumStr(n2NumStr)

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStr(n2NumStr). "+
				"numStr='%v' Error='%v' ", n2NumStr, err.Error())
	}

	b2Pair := BigIntPair{}

	b2Pair.SetBigIntPair(b1Num, b2Num)

	return b2Pair, nil

}

// NewNumStr - Creates a new BigIntPair instance from two integerNum strings
// passed as input parameters.
//
// The returned BigIntNum's BigIntPair.Big1 and BigIntPair.Big2 will contain default
// numeric separators (decimal separator, thousands separator and currency symbol).
//
func (bPair BigIntPair) NewNumStrWithNumSeps(
	n1NumStr, n2NumStr string, numSeps NumericSeparatorDto) (BigIntPair, error) {

	ePrefix := "BigIntPair.NewNumStrDto() "

	numSeps.SetDefaultsIfEmpty()

	b1Num, err := BigIntNum{}.NewNumStrWithNumSeps(n1NumStr, numSeps)

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps("+
				"n1NumStr, numSeps). "+
				"n1NumStr='%v' Error='%v' ", n1NumStr, err.Error())
	}

	b2Num, err := BigIntNum{}.NewNumStrWithNumSeps(n2NumStr, numSeps)

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStrWithNumSeps(n2NumStr, numSeps). "+
				"n2NumStr='%v' Error='%v' ", n2NumStr, err.Error())
	}

	b2Pair := BigIntPair{}

	b2Pair.SetBigIntPair(b1Num, b2Num)

	return b2Pair, nil
}

// NewNumStrDto - Creates a new BigIntPair instance from two NumStrDto
// instances passed as input parameters.
//
// This method will test the validity of n1Dto and n2Dto
//
// The returned BigIntNum BigIntPair.Big1 will contain numeric separators (decimal
// separator, thousands separator and currency symbol) copied from input parameter,
// 'n1Dto'. BigIntPair.Big2 will contain numeric separators copied from 'n2Dto'.
//
func (bPair BigIntPair) NewNumStrDto(n1Dto, n2Dto NumStrDto) (BigIntPair, error) {

	ePrefix := "BigIntPair.NewNumStrDto() "

	// This method will test the validity of n1Dto
	b1Num, err := n1Dto.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by n1Dto.GetBigIntNum(). "+
				"numStr='%v' Error='%v' ", n1Dto.GetNumStr(), err.Error())
	}

	// This method will test the validity of n2Dto
	b2Num, err := n2Dto.GetBigIntNum()

	if err != nil {
		return BigIntPair{},
			fmt.Errorf(ePrefix+"Error returned by n2Dto.GetBigIntNum(). "+
				"numStr='%v' Error='%v' ", n2Dto.GetNumStr(), err.Error())
	}

	b2Pair := BigIntPair{}

	b2Pair.SetBigIntPair(b1Num, b2Num)

	return b2Pair, nil
}

// SetBigIntPairSetBigIntPair -Sets the values of the current
// BigIntPair instance to the input values of b1 and b2
// respectively.
//
// This method will set BigIntNum BigIntPair.Big1 with numeric separators (decimal
// separator, thousands separator and currency symbol) copied from BigIntNum input
// parameter, 'b1'. BigIntPair.Big2 will contain numeric separators copied from 'b2'.
//
func (bPair *BigIntPair) SetBigIntPair(b1, b2 BigIntNum) {

	bPair.Empty()

	bPair.Big1.CopyIn(b1)
	bPair.Big2.CopyIn(b2)

	bPair.Big1Compare = bPair.Big1.bigInt.Cmp(bPair.Big2.bigInt)
	bPair.Big1AbsCompare = bPair.Big1.absBigInt.Cmp(bPair.Big2.absBigInt)

	if bPair.Big1.precision == bPair.Big2.precision {
		bPair.Precision1Compare = 0
	} else if bPair.Big1.precision > bPair.Big2.precision {
		bPair.Precision1Compare = 1
	} else {
		// Must be bPair.Big1.precision < bPair.Big2.precision
		bPair.Precision1Compare = -1
	}
}
