package mathops

import (
	"fmt"
	"math/big"
)

type FactorialDto struct {
	UpperLimit uint64
	LowerLimit uint64
}

type NFactorial struct {
	NumOfTrials uint64
	NFac        FactorialDto
}

func (nFac NFactorial) GetFactorialArray(nFactorial int) []int {

	limit := nFactorial

	result := make([]int, limit)

	for i := 0; i < limit; i++ {

		result[i] = nFactorial
		nFactorial--

	}

	return result
}

// CalcFactorialValueInt - Computes the value of nFactorial! and returns that value as a
// BigIntNum.
//
// Input Parameters:
// =================
//
// nFactorial	int		- The starting value in the factorial calculation. 'nFactorial' MUST BE
//                    a positive integer number.
//
//
// lowerLimit int		- The lower boundary for the factorial calculation. 'lowerLimit' MUST BE
//										a positive integer number.
//
// Examples:
// =========
//
//				1.	nFactorial = 7  and lowerLimit = 3
//
// 						The input parameter 'lowerLimit' specifies the lower boundary for the calculation.
// 						'nFactorial' = 7 and 'lowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	nFactorial = 7 and lowerLimit = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFactorial = 7 and lowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueInt(nFactorial, lowerLimit int) (BigIntNum, error) {

	ePrefix := "NFactorial.CalcFactorialValueInt() "

	if nFactorial < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'nFactorial' is less than zero! "+
				"'nFactorial' must be a positive integer. 'nFactorial'='%v' ",
				nFactorial)
	}

	if lowerLimit < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'lowerLimit' is less than zero! "+
				"'lowerLimit' must be a positive integer. 'lowerLimit'='%v' ",
				lowerLimit)
	}

	nFacBigInt := big.NewInt(int64(nFactorial))

	lowerLimitBigInt := big.NewInt(int64(lowerLimit))

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcFactorialValueInt32 - Computes the value of nFactorial! and returns that value as a
// BigIntNum.
//
// Input Parameters:
// =================
//
// nFactorial	int32		- The starting value in the factorial calculation. 'nFactorial' MUST BE
//                    	a positive integer number.
//
//
// lowerLimit int32		- The lower boundary for the factorial calculation. 'lowerLimit' MUST BE
//											a positive integer number.
//
// Examples:
// =========
//
//				1.	nFactorial = 7  and lowerLimit = 3
//
// 						The input parameter 'lowerLimit' specifies the lower boundary for the calculation.
// 						'nFactorial' = 7 and 'lowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	nFactorial = 7 and lowerLimit = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFactorial = 7 and lowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueInt32(nFactorial, lowerLimit int32) (BigIntNum, error) {

	ePrefix := "NFactorial.CalcFactorialValueInt32() "

	if nFactorial < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'nFactorial' is less than zero! "+
				"'nFactorial' must be a positive integer. 'nFactorial'='%v' ",
				nFactorial)
	}

	if lowerLimit < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'lowerLimit' is less than zero! "+
				"'lowerLimit' must be a positive integer. 'lowerLimit'='%v' ",
				lowerLimit)
	}

	nFacBigInt := big.NewInt(int64(nFactorial))

	lowerLimitBigInt := big.NewInt(int64(lowerLimit))

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcFactorialValueInt64 - Computes the value of nFactorial! and returns that value as a
// BigIntNum.
//
// Input Parameters:
// =================
//
// nFactorial	int64		- The starting value in the factorial calculation. 'nFactorial' MUST BE
//                    	a positive integer number.
//
//
// lowerLimit int64		- The lower boundary for the factorial calculation. 'lowerLimit' MUST BE
//											a positive integer number.
//
// Examples:
// =========
//
//				1.	nFactorial = 7  and lowerLimit = 3
//
// 						The input parameter 'lowerLimit' specifies the lower boundary for the calculation.
// 						'nFactorial' = 7 and 'lowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	nFactorial = 7 and lowerLimit = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFactorial = 7 and lowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueInt64(nFactorial, lowerLimit int64) (BigIntNum, error) {

	ePrefix := "NFactorial.CalcFactorialValueInt64() "

	if nFactorial < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'nFactorial' is less than zero! "+
				"'nFactorial' must be a positive integer. 'nFactorial'='%v' ",
				nFactorial)
	}

	if lowerLimit < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: Input parameter 'lowerLimit' is less than zero! "+
				"'lowerLimit' must be a positive integer. 'lowerLimit'='%v' ",
				lowerLimit)
	}

	nFacBigInt := big.NewInt(nFactorial)

	lowerLimitBigInt := big.NewInt(lowerLimit)

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcNFactorialValue - Computes the value of n factorial as expressed by the upper
// and lower limits of the input parameter, 'nFactorial'. 'nFactorial' is of type
// 'FactorialDto'. The result of the n factorial calculation is returned as a BigIntNum
// type.
//
// Input Parameters:
// =================
//
// nFactorial	FactorialDto		- This structure contains a value of 'UpperLimit' or initial starting value
// 															of the the factorial calculation. In addition, the structure contains a
//															data field, 'LowerLimit', which specifies the lower boundary for the factorial
//															calculation. Both 'UpperLimit' and 'LowerLimit' are uint64 types.
//
// Examples:
// =========
//
//				1.	nFactorial.UpperLimit = 7  and nFactorial.LowerLimit = 3
//
// 						The input parameter 'LowerLimit' specifies the lower boundary for the calculation.
// 						'UpperLimit' = 7 and 'LowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'UpperLimit' = 7 and 'LowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFactorial = 7 and lowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcNFactorialValue(nFactorial FactorialDto) (BigIntNum, error) {

	ePrefix := "NFactorial.CalcNFactorialValue() "

	if nFactorial.LowerLimit > nFactorial.UpperLimit {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: 'nFactorial' Upper Limit is LESS THAN Lower Limit!")
	}

	if nFactorial.LowerLimit < 1 {
		nFactorial.LowerLimit = 1
	}

	if nFactorial.UpperLimit < 2 {
		return BigIntNum{}.NewOne(0), nil
	}

	if nFactorial.UpperLimit == nFactorial.LowerLimit {
		return BigIntNum{}.NewOne(0), nil
	}

	nFacUpperLimit := big.NewInt(0).SetUint64(nFactorial.UpperLimit)

	nFacLowerLimit := big.NewInt(0).SetUint64(nFactorial.LowerLimit)

	result, err := nFac.CalcFactorialValueBigInt(nFacUpperLimit, nFacLowerLimit)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by nFac.CalcFactorialValueBigInt(nFacUpperLimit, nFacLowerLimit) "+
				"Error='%v'\n", err.Error())
	}

	return result, nil
}

// CalcFactorialValueBigIntNum - Computes the value of n factorial as expressed by the
// upper and lower limits of the input parameters, 'nFacUpperLimit' and 'nFacLowerLimit'.
// Both input parameters are passed to this method as type, BigIntNum.
//
// The result of the n factorial calculation is returned as a BigIntNum type.
//
// Input Parameters:
// =================
//
// nFacUpperLimit	BigIntNum		- This BigIntNum type must be a positive integer value. It
// 															represents the upper limit or starting value of the n-factorial
//															calculation.
//
// nFacLowerLimit	BigIntNum		- This BigIntNum type must be a positive integer value. It
//															calculation.  This value is NOT multiplied by the previous
//                              n-factorial value and is therefore NOT included in the
// 															calculation. It's sole purpose is to signal the end of
// 															n-factorial calculations.
//
// Return Values
// =============
//
// BigIntNum type						  - The value of the n-factorial calculation is returned as
// 															a type BigIntNum. The returned BigIntNum type will always
// 															be a positive integer value.
//
// error											- For a successful calculation, both 'nFacUpperLimit' and
//                              'nFacLowerLimit' must be positive integer values. In
// 															addition, 'nFacUpperLimit' must be greater than
// 															'nFacLowerLimit'. If these conditions are not met, an
//															error will be returned.
//
// Examples:
// =========
//
//				1.	nFacUpperLimit = 7  and nFacLowerLimit = 3
//
// 						The input parameter 'nFacLowerLimit' specifies the lower boundary
// 						for the calculation. Values of 'nFacUpperLimit' = 7 and
// 						'nFacLowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'nFacUpperLimit' = 7 and 'nFacLowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFacUpperLimit = 7 and nFacLowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueBigIntNum(
	nFacUpperLimit, nFacLowerLimit BigIntNum) (BigIntNum, error) {

	return nFac.CalcFactorialValueBigInt(nFacUpperLimit.bigInt, nFacLowerLimit.bigInt)
}

// CalcFactorialValueUint - Computes the value of n factorial as expressed by the
// upper and lower limits of the input parameters, 'nFacUpperLimit' and 'nFacLowerLimit'.
// Both input parameters are passed to this method as type, uint
//
// The result of the n factorial calculation is returned as a BigIntNum type.
//
// Input Parameters:
// =================
//
// nFacUpperLimit	uint		- This uint type represents the upper limit or starting value
// 															of the n-factorial calculation.
//
// nFacLowerLimit	uint		- This uint type  represents the lower limit of the n-factorial
// 															calculation.  This value is NOT multiplied by the previous
//                              n-factorial value and is therefore NOT included in the
// 															calculation results. It's sole purpose is to signal the end
// 															of n-factorial calculations.
//
// Return Values
// =============
//
// BigIntNum type						  - The value of the n-factorial calculation is returned as
// 															a type BigIntNum. The returned BigIntNum type will always
// 															be a positive integer value.
//
// error											- For a successful calculation, 'nFacUpperLimit' must be greater
// 															than 'nFacLowerLimit'. If these conditions are not met, an
//															error will be returned.
//
// Examples:
// =========
//
//				1.	nFacUpperLimit = 7  and nFacLowerLimit = 3
//
// 						The input parameter 'nFacLowerLimit' specifies the lower boundary
// 						for the calculation. Values of 'nFacUpperLimit' = 7 and
// 						'nFacLowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'nFacUpperLimit' = 7 and 'nFacLowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFacUpperLimit = 7 and nFacLowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueUint(nFactorial, lowerLimit uint) (BigIntNum, error) {

	nFacBigInt := big.NewInt(int64(nFactorial))

	lowerLimitBigInt := big.NewInt(int64(lowerLimit))

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcFactorialValueUint32 - Computes the value of n factorial as expressed by the
// upper and lower limits of the input parameters, 'nFacUpperLimit' and 'nFacLowerLimit'.
// Both input parameters are passed to this method as type, uint32
//
// The result of the n factorial calculation is returned as a BigIntNum type.
//
// Input Parameters:
// =================
//
// nFacUpperLimit	uint32		- This uint32 type represents the upper limit or starting value
// 															of the n-factorial calculation.
//
// nFacLowerLimit	uint32		- This uint32 type  represents the lower limit of the n-factorial
// 															calculation.  This value is NOT multiplied by the previous
//                              n-factorial value and is therefore NOT included in the
// 															calculation results. It's sole purpose is to signal the end
// 															of n-factorial calculations.
//
// Return Values
// =============
//
// BigIntNum type						  - The value of the n-factorial calculation is returned as
// 															a type BigIntNum. The returned BigIntNum type will always
// 															be a positive integer value.
//
// error											- For a successful calculation, 'nFacUpperLimit' must be greater
// 															than 'nFacLowerLimit'. If these conditions are not met, an
//															error will be returned.
//
// Examples:
// =========
//
//				1.	nFacUpperLimit = 7  and nFacLowerLimit = 3
//
// 						The input parameter 'nFacLowerLimit' specifies the lower boundary
// 						for the calculation. Values of 'nFacUpperLimit' = 7 and
// 						'nFacLowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'nFacUpperLimit' = 7 and 'nFacLowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFacUpperLimit = 7 and nFacLowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueUint32(nFactorial, lowerLimit uint32) (BigIntNum, error) {

	nFacBigInt := big.NewInt(int64(nFactorial))

	lowerLimitBigInt := big.NewInt(int64(lowerLimit))

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcFactorialValueUint64 - Computes the value of n factorial as expressed by the
// upper and lower limits of the input parameters, 'nFacUpperLimit' and 'nFacLowerLimit'.
// Both input parameters are passed to this method as type, uint64
//
// The result of the n factorial calculation is returned as a BigIntNum type.
//
// Input Parameters:
// =================
//
// nFacUpperLimit	uint64		- This uint64 type represents the upper limit or starting value
// 															of the n-factorial calculation.
//
// nFacLowerLimit	uint64		- This uint64 type  represents the lower limit of the n-factorial
// 															calculation.  This value is NOT multiplied by the previous
//                              n-factorial value and is therefore NOT included in the
// 															calculation results. It's sole purpose is to signal the end
// 															of n-factorial calculations.
//
// Return Values
// =============
//
// BigIntNum type						  - The value of the n-factorial calculation is returned as
// 															a type BigIntNum. The returned BigIntNum type will always
// 															be a positive integer value.
//
// error											- For a successful calculation, 'nFacUpperLimit' must be greater
// 															than 'nFacLowerLimit'. If these conditions are not met, an
//															error will be returned.
//
// Examples:
// =========
//
//				1.	nFacUpperLimit = 7  and nFacLowerLimit = 3
//
// 						The input parameter 'nFacLowerLimit' specifies the lower boundary
// 						for the calculation. Values of 'nFacUpperLimit' = 7 and
// 						'nFacLowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'nFacUpperLimit' = 7 and 'nFacLowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFacUpperLimit = 7 and nFacLowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueUint64(nFactorial, lowerLimit uint64) (BigIntNum, error) {

	nFacBigInt := big.NewInt(0).SetUint64(nFactorial)

	lowerLimitBigInt := big.NewInt(0).SetUint64(lowerLimit)

	return nFac.CalcFactorialValueBigInt(nFacBigInt, lowerLimitBigInt)
}

// CalcFactorialValueBigInt - Calculates the value of n factorial using an upper and lower limit.
// Both the upper limit, 'nFacUpperLimit', and the lower limit, 'nFacLowerLimit', are passed as
// *big.Int types.
//
// The resulting value of the n factorial calculation is returned as a 'BigIntNum' type.
//
// Input Parameters:
// =================
//
// nFacUpperLimit	*big.Int	- This *big.Int type represents the upper limit or starting value
// 															of the n-factorial calculation.
//
// nFacLowerLimit	*big.Int	- This *big.Int type  represents the lower limit of the n-factorial
// 															calculation.  This value is NOT multiplied by the previous
//                              n-factorial value and is therefore NOT included in the
// 															calculation results. It's sole purpose is to signal the end
// 															of n-factorial calculations.
//
// Return Values
// =============
//
// BigIntNum type						  - The value of the n-factorial calculation is returned as
// 															a type BigIntNum. The returned BigIntNum type will always
// 															be a positive integer value.
//
// error											- For a successful calculation, 'nFacUpperLimit' must be greater
// 															than 'nFacLowerLimit'. If these conditions are not met, an
//															error will be returned.
//
// Examples:
// =========
//
//				1.	nFacUpperLimit = 7  and nFacLowerLimit = 3
//
// 						The input parameter 'nFacLowerLimit' specifies the lower boundary
// 						for the calculation. Values of 'nFacUpperLimit' = 7 and
// 						'nFacLowerLimit' = 3  will yield a calculation of:
//														7x6x5x4 = 840
//
//				2.	'nFacUpperLimit' = 7 and 'nFacLowerLimit' = 7
// 						Equivalent of 0! and 0! = 1
//
//        3.	nFacUpperLimit = 7 and nFacLowerLimit = 1
//														7x6x5x4x3x2 = 5040
//
func (nFac NFactorial) CalcFactorialValueBigInt(
	nFacUpperLimit, nFacLowerLimit *big.Int) (result BigIntNum, err error) {

	ePrefix := "NFactorial.CalcFactorialValueBigInt() "

	result = BigIntNum{}.NewZero(0)
	err = nil

	cmpResult := nFacLowerLimit.Cmp(big.NewInt(0))

	// lower limit Must be at least '1'
	if cmpResult == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: Lower Limit is less than '0'. nFacLowerLimit='%v'", nFacLowerLimit.Text(10))
		return result, err
	}

	if cmpResult == 0 {
		nFacLowerLimit = big.NewInt(1)
	}

	cmpResult = nFacUpperLimit.Cmp(nFacLowerLimit)

	if cmpResult == -1 {
		err = fmt.Errorf(ePrefix+
			"Error: 'nFacUpperLimit' is less than 'nFacLowerLimit'. "+
			"nFacUpperLimit='%v' nFacLowerLimit='%v'", nFacUpperLimit.Text(10), nFacLowerLimit.Text(10))
		return result, err
	}

	// nFacUpperLimit = nFacLowerLimit
	// This is equivalent to 0!
	// 0! = 1
	if cmpResult == 0 {
		err = nil
		result = BigIntNum{}.NewBigInt(big.NewInt(1), 0)
		return result, err
	}

	total := big.NewInt(0).Set(nFacUpperLimit)

	bigOne := big.NewInt(1)
	count := big.NewInt(0).Sub(total, bigOne)
	cmpResult = count.Cmp(nFacLowerLimit)

	for cmpResult == 1 {
		total = big.NewInt(0).Mul(total, count)
		count = big.NewInt(0).Sub(count, bigOne)
		cmpResult = count.Cmp(nFacLowerLimit)
	}

	result = BigIntNum{}.NewBigInt(total, 0)
	err = nil

	return result, err
}
