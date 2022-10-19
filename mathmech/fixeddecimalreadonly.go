package mathops

import (
	"fmt"
	"math/big"
	"sync"
)

// FixedDecimalReadOnly - encapsulates a Read Only
// BigIntFixedDecimal objects. This type is designed
// to be used with constant numeric values that should
// not be changed.
//
// A series of related functions provides the means to
// read the numeric value encapsulated by the FixedDecimalReadOnly
// type.
type FixedDecimalReadOnly struct {
	fixedDecimal BigIntFixedDecimal

	lock *sync.Mutex
}

// GetFixedDecimal - Returns a deep copy of the underlying
// BigIntFixedDecimal value.
func (fDecRO *FixedDecimalReadOnly) GetFixedDecimal() BigIntFixedDecimal {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	fDecRO.fixedDecimal.IsValid()

	return fDecRO.fixedDecimal.CopyOut()
}

// GetInteger - Returns the *big.Int integer value from the
// underlying BigIntFixedDecimal
func (fDecRO *FixedDecimalReadOnly) GetInteger() *big.Int {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	fDecRO.fixedDecimal.IsValid()

	fDecRO.fixedDecimal.IsValid()

	return fDecRO.fixedDecimal.GetInteger()

}

// GetNumStr - Converts the underlying BigIntFixedDecimal
// value to a string of numbers which includes the decimal
// place and decimal digits if they exist. Note that
// the period ('.') is the decimal separator character
// always used in the returned number string to separate
// integer and fractional digits.
func (fDecRO *FixedDecimalReadOnly) GetNumStr() string {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	_ = fDecRO.fixedDecimal.IsValid()

	return fDecRO.fixedDecimal.GetNumStr()

}

// GetPrecision - returns an unsigned integer specifying
// the number of digits to the right of the decimal place.
func (fDecRO *FixedDecimalReadOnly) GetPrecision() uint {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	_ = fDecRO.fixedDecimal.IsValid()

	return fDecRO.fixedDecimal.GetPrecision()

}

// GetPrecisionBigInt - returns an unsigned integer specifying
// the number of digits to the right of the decimal place.
func (fDecRO *FixedDecimalReadOnly) GetPrecisionBigInt() *big.Int {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	_ = fDecRO.fixedDecimal.IsValid()

	return fDecRO.fixedDecimal.GetPrecisionBigInt()

}

func (fDecRO *FixedDecimalReadOnly) GetBigIntPrecision() (integer, precision *big.Int) {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	integer = fDecRO.GetInteger()

	if integer == nil {
		integer = big.NewInt(0)
	}

	precision = fDecRO.GetPrecisionBigInt()

	if precision == nil {
		precision = big.NewInt(0)
	}

	return integer, precision
}

// IsValid - Performs diagnostic and remedial actions on
// the underlying BigIntFixedDecimal instance. Returns
// 'false' if the instance is uninitialized.
func (fDecRO *FixedDecimalReadOnly) IsValid() bool {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	return fDecRO.fixedDecimal.IsValid()
}

func (fDecRO *FixedDecimalReadOnly) IsZero() bool {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	return fDecRO.fixedDecimal.IsZero()
}

// NewNumStr - Creates and returns a new FixedDecimalReadOnly
// instance initialized to the value of the number string
// input parameter, 'numStr'.
//
// A number string is a string of numeric digits which may,
// or may not, be prefixed with a minus sign ('-') indicating
// a negative number. If the numeric string of digits is prefixed
// by a left parenthesis "(" and suffixed by a corresponding
// right parenthesis ")", this also indicates a negative value.
//
// The numeric string of digits may also contain a period
// ('.') which is treated as a decimal separator and used to
// separate integer and fractional digits within the number
// string.
//
// The only decimal separator recognized by this method is the
// period ('.').
func (fDecRO *FixedDecimalReadOnly) NewNumStr(numStr string) (FixedDecimalReadOnly, error) {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	ePrefix := "FixedDecimalReadOnly.NewNumStr() "

	fo2 := FixedDecimalReadOnly{}

	fo2.fixedDecimal = new(BigIntFixedDecimal).NewZero(0)

	err := fo2.fixedDecimal.SetNumStr(numStr)

	if err != nil {
		return new(FixedDecimalReadOnly).NewZero(0),
			fmt.Errorf(ePrefix)
	}

	return fo2, nil
}

// NewFixedDecimal - Receives a BigIntFixedDecimal instance
// as an input parameter and returns a new FixedDecimalReadOnly
// object.
func (fDecRO *FixedDecimalReadOnly) NewFixedDecimal(
	fixedDecimal BigIntFixedDecimal) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	_ = fixedDecimal.IsValid()

	f2 := new(FixedDecimalReadOnly).NewZero(0)

	f2.fixedDecimal.CopyIn(fixedDecimal)

	return f2
}

// NewInt - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'intValue' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewInt(
	intValue int,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewInt(intValue, precision)

	return f2

}

// NewInt32 - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'int32Value' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewInt32(
	int32Value int32,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewInt32(int32Value, precision)

	return f2

}

// NewInt64 - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'int64Value' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewInt64(
	int64Value int64,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewInt64(int64Value, precision)

	return f2

}

// NewUInt - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'uintValue' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewUInt(
	uintValue,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewUInt(uintValue, precision)

	return f2

}

// NewUInt32 - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'uintValue' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewUInt32(
	uint32Value uint32,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewUInt32(uint32Value, precision)

	return f2

}

// NewUInt64 - Creates and returns a new FixedDecimalReadOnly instance initialized
// to by the input parameters 'uintValue' and 'precision'. 'precision' specifies
// the number of digits to the right of the decimal place in 'intValue'.
func (fDecRO *FixedDecimalReadOnly) NewUInt64(
	uint64Value uint64,
	precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewUInt64(uint64Value, precision)

	return f2
}

// NewZero - Creates and returns a new FixedDecimalReadOnly instance initialized
// to zero. The input parameter 'precision' specifies the number of zero digits
// to the right of the decimal place.
func (fDecRO *FixedDecimalReadOnly) NewZero(precision uint) FixedDecimalReadOnly {

	if fDecRO.lock == nil {
		fDecRO.lock = new(sync.Mutex)
	}

	fDecRO.lock.Lock()

	defer fDecRO.lock.Unlock()

	f2 := FixedDecimalReadOnly{}

	f2.fixedDecimal = new(BigIntFixedDecimal).NewZero(precision)

	return f2

}
