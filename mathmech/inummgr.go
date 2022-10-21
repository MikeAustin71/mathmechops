package mathops

import "math/big"

type INumMgr interface {
	GetNumStr() string
	GetBigInt() (*big.Int, error)
	GetBigIntNum() (BigIntNum, error)
	GetDecimal() (Decimal, error)
	GetIntAry() (IntAry, error)
	GetNumericSeparatorsDto() NumericSeparatorDto
	GetNumStrDto() (NumStrDto, error)
	GetPrecisionUint() uint
	GetSign() int
	IsValid(errName string) error
	IsZero() bool
	SetNumericSeparatorsToDefaultIfEmpty()
	SetNumericSeparatorsDto(customSeparators NumericSeparatorDto) error
}
