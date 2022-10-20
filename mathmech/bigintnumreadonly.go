package mathops

import (
	"fmt"
	"math/big"
)

type BigIntNumReadOnly struct {
	bigIntNum BigIntNum
}

// GetIntegerValue
//
// Returns all the numeric digits in the base BigIntNum
// as a *big.Int integer value
func (birO *BigIntNumReadOnly) GetIntegerValue() *big.Int {

	return birO.bigIntNum.GetIntegerValue()
}

// GetPrecisionUint
//
// Returns the precision of the underlying BigIntNum
// as a type uint.
func (birO *BigIntNumReadOnly) GetPrecisionUint() uint {

	return birO.bigIntNum.GetPrecisionUint()
}

// GetBigIntNum
//
// Returns a deep copy of the underlying BigIntNum.
func (birO *BigIntNumReadOnly) GetBigIntNum() BigIntNum {

	return birO.bigIntNum.CopyOut()
}

// GetFixedDecimal
//
// Returns a deep copy of the underlying BigIntNum
// numeric value as a type BigIntFixedDecimal
func (birO *BigIntNumReadOnly) GetFixedDecimal() BigIntFixedDecimal {

	return birO.bigIntNum.GetBigIntFixedDecimal()
}

// NewBigIntNum
//
// Receives a BigIntNum parameter and returns a new
// BigIntNumReadOnly instance.
func (birO *BigIntNumReadOnly) NewBigIntNum(biNum BigIntNum) BigIntNumReadOnly {

	birO2 := BigIntNumReadOnly{}

	birO2.bigIntNum = new(BigIntNum).NewZero(0)

	birO2.bigIntNum.CopyIn(biNum)

	return birO2
}

// NewBigIntNum
//
// Receives a BigIntFixedDecimal parameter and returns a new
// BigIntNumReadOnly instance.
func (birO *BigIntNumReadOnly) NewFixedDecimal(fixedDec BigIntFixedDecimal) BigIntNumReadOnly {

	birO2 := BigIntNumReadOnly{}

	birO2.bigIntNum = new(BigIntNum).NewBigIntFixedDecimal(fixedDec)

	return birO2
}

// NewNumStr
//
//	Receives a number string as input and returns a new
//	BigIntNumReadOnly instance.
//
//	This method assumes that the input parameter 'numStr'
//	is a pure number string of numeric digits which may
//	be delimited by default United States (US) numeric
//	separators. Default US numeric separators are defined
//	as:
//
//		decimal separator = '.'
//		integer separator = ','
//		integer grouping = 3 (thousands)
//		currency symbol = '$'
func (birO *BigIntNumReadOnly) NewNumStr(numStr string) (BigIntNumReadOnly, error) {

	ePrefix := "BigIntNumReadOnly.NewNumStr() "

	readOnly, err := new(BigIntNum).NewNumStr(numStr)

	if err != nil {
		return BigIntNumReadOnly{},
			fmt.Errorf(ePrefix+"Error returned by BigIntNum{}.NewNumStr(numStr). "+
				"Error='%v' ", err.Error())
	}

	biRo := BigIntNumReadOnly{}

	biRo.bigIntNum = new(BigIntNum).NewZero(0)
	biRo.bigIntNum.CopyIn(readOnly)

	return biRo, nil
}
