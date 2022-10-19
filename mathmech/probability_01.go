package mathops

import (
	"errors"
	"fmt"
	"math/big"
)

/*

This file contains Probability source code for 'Combinations'.

For 'Permutations' see source code file:
	MikeAustin71\mathopsgo\mathops\probability_02.go


                                Combinations
																============

		Reference:
    ##########
		https://www.mathsisfun.com/combinatorics/combinations-permutations.html
		https://www.youtube.com/watch?v=bCxMhncR7PU&list=PL06A16C388F14E6FE&index=22
		https://www.youtube.com/watch?v=ZcSSI6VY1kM
		https://www.probabilitycourse.com/courses.php

$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
  *** This combination calculation assumes NO REPETITIONS! ***
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$

 	In the following combination formula, n= 'numOfItems'  and r = 'numOfItemsChosen'

  Formula WITHOUT Repetitions
  ---------------------------
																			 n!
											 nCr	 =		-----------
																	 (n-r)! r!

	Example WITHOUT Repetitions
  ---------------------------
  16 pool balls. How many ways to choose 3 pool balls with NO repetitions or No repeats.

                  16!            16 x 15 x 14
	16C3 =    -------------  =    --------------  =    560
            (16-3)! x 3!            3 x 2


$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
  *** This combination calculation assumes REPETITIONS ARE ALLOWED! ***
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$

 	In the following combination formula, n= 'numOfItems'  and r = 'numOfItemsChosen'

  Formula WITH Repetitions
  ------------------------
																	(r + n - 1)!
											nCr    =  -------------------
																	  r! (n-1)!

	Example WITH Repetitions:
  -------------------------
	Let us say there are five flavors of icecream:
  (1) banana, (2) chocolate, (3) lemon, (4) strawberry and (5) vanilla.

	We can have three scoops. How many variations will there be?
         nCr = (n==5 r==3) = 5C3 = Answer: 35


*/

type Probability struct {
	NoOfTrials                   BigIntNum
	PossibleOutcomesPerTrial     BigIntNum
	NoSuccessfulOutcomesPerTrial BigIntNum
	TotPossibleOutcomes          BigIntNum
	TotSuccessfulOutcomes        BigIntNum
	PercentCertainty             BigIntNum
}

// CombinationsNoRepsBigInt - Calculates the number of combinations from
// 'numOfItems' and 'numOfItemsChosen' or 'n' things chosen 'r' at a time
// with NO repetitions and order does NOT matter. This is also referred to
// as an unordered sampling WITHOUT replacement. The calculation result is
// returned as a BigIntNum type.
//
// Input Parameters
// ================
//
// numOfItems 				*big.Int	- Must be a positive integer number greater than zero.
//                                'numOfItems' must be greater than or equal to
// 																'numOfItemsChosen'.
//
// numOfItemsChosen 	*big.Int	- Must be a positive integer number greater than zero.
// 																'numOfItemsChosen' must be less than or equal to
// 																'numOfItems'.
//
// Returns
// =======
//
// BigIntNum			- If the calculation is successful, the result is returned as a
//                  BigIntNum type. If the calculation fails, the error return is
//                  populated.
//
// error					- If the calculation is successful, this return value is 'nil'. If
//                  the calculation fails, 'error' is populated with an appropriate
//                  error message.
//
// Calculation
// ===========
//
// The calculation performed by this method uses the following combinations formula,
// n= 'numOfItems' and r = 'numOfItemsChosen'.
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, repetition is NOT allowed,
// 						and order does NOT matter.
//
// *** This calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
// Example WITHOUT Repetitions
// ===========================
//
// There are 16 pool balls. How many ways to choose 3 pool balls with NO repetitions
// or No repeats.
//
//                  16!            16 x 15 x 14
//	16C3 =    -------------  =    --------------  =    560
//            (16-3)! x 3!            3 x 2
//
//
// Note: 0! = 1
//
func (prob Probability) CombinationsNoRepsBigInt(
	numOfItems, numOfItemsChosen *big.Int) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsNoRepsBigInt() "
	bigZero := big.NewInt(0)

	if numOfItems.Cmp(bigZero) == 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.Cmp(bigZero) < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.Cmp(bigZero) < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItems.Cmp(numOfItemsChosen) < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: 'numOfItems' is LESS THAN 'numOfItemsChosen'! ")

	}

	if numOfItemsChosen.Cmp(bigZero) == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: 'numOfItemsChosen' is ZERO! ")
	}

	bigOne := big.NewInt(1)

	if numOfItemsChosen.Cmp(bigOne) == 0 {
		return BigIntNum{}.NewBigInt(numOfItems, 0), nil
	}

	numeratorNUpperLimit := big.NewInt(0).Set(numOfItems)
	numeratorNLowerLimit := big.NewInt(0).Set(bigOne)

	nMinusR := big.NewInt(0).Sub(numOfItems, numOfItemsChosen)

	if nMinusR.Cmp(bigZero) == 0 {
		nMinusR = big.NewInt(0).Set(bigOne)
	}

	var nFactorial BigIntNum
	var rFactorial BigIntNum
	var err error

	if nMinusR.Cmp(numOfItemsChosen) == 1 {
		// nMinusR is Greater Than numOfItemsChosen
		numeratorNLowerLimit = big.NewInt(0).Set(nMinusR)
		nFactorial, err = NFactorial{}.CalcFactorialValueBigInt(
			numeratorNUpperLimit, numeratorNLowerLimit)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
					"numeratorNUpperLimit, numeratorNLowerLimit). "+
					"numeratorNUpperLimit='%v'  numeratorNLowerLimit=nMinusR='%v' Error='%v'",
					numeratorNUpperLimit.Text(10), numeratorNLowerLimit.Text(10),
					err.Error())
		}

		rUpperLimit := big.NewInt(0).Set(numOfItemsChosen)
		rLowerLimit := big.NewInt(0).Set(bigOne)

		rFactorial, err = NFactorial{}.CalcFactorialValueBigInt(
			rUpperLimit, rLowerLimit)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
					"rUpperLimit, rLowerLimit). "+
					"rUpperLimit='%v'  rLowerLimit='%v' Error='%v'",
					rUpperLimit.Text(10), rLowerLimit.Text(10),
					err.Error())
		}

	} else {
		// nMinusR is Less Than OR Equal to numOfItemsChosen

		numeratorNLowerLimit = big.NewInt(0).Set(numOfItemsChosen)

		nFactorial, err = NFactorial{}.CalcFactorialValueBigInt(
			numeratorNUpperLimit, numeratorNLowerLimit)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
					"numeratorNUpperLimit, numeratorNLowerLimit). "+
					"numeratorNUpperLimit='%v'  numeratorNLowerLimit=numOfItemsChosen='%v' Error='%v'",
					numeratorNUpperLimit.Text(10), numeratorNLowerLimit.Text(10),
					err.Error())
		}

		rFactorial, err = NFactorial{}.CalcFactorialValueBigInt(
			nMinusR, big.NewInt(0).Set(bigOne))

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
					"nMinusR, 1). "+
					"UpperLimit=nMinusR='%v'  Error='%v'",
					nMinusR.Text(10), err.Error())
		}
	}

	combinationsResult, err :=
		BigIntMathDivide{}.BigIntNumFracQuotient(nFactorial, rFactorial, 10)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient("+
				"nFactorial, rFactorial, 10)). "+
				"nFactorial='%v'  rFactorial='%v' Error='%v'",
				nFactorial.GetNumStr(), rFactorial.GetNumStr(),
				err.Error())
	}

	return combinationsResult, nil
}

// CombinationsWithRepsBigInt - Calculates the number of combinations from
// 'numOfItems' and 'numOfItemsChosen' or 'n' things chosen 'r' at a time
// WITH repetitions and order does NOT matter. This is also referred to
// as an unordered sampling WITH replacement. The calculation result is
// returned as a BigIntNum type.
//
// Input Parameters
// ================
//
// numOfItems 				*big.Int	- Must be a positive integer number greater than zero.
//
// numOfItemsChosen 	*big.Int	- Must be a positive integer number greater than zero.
//
// Returns
// =======
//
// BigIntNum			- If the calculation is successful, the result is returned as a
//                  BigIntNum type. If the calculation fails, the error return is
//                  populated.
//
// error					- If the calculation is successful, this return value is 'nil'. If
//                  the calculation fails, 'error' is populated with an appropriate
//                  error message.
//
// Calculation
// ===========
//
// The calculation performed by this method uses the following combinations formula,
// n= 'numOfItems' and r = 'numOfItemsChosen'. Since this calculation applies to
// unordered sampling WITH replacement, 'numOfItemsChosen' may exceed 'numOfItems'.
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// Example WITH Repetitions:
// -------------------------
// Let us say there are five flavors of icecream:
// (1) banana, (2) chocolate, (3) lemon, (4) strawberry and (5) vanilla.
//
// We can have three scoops. How many variations will there be?
// nCr = (n==5 r==3) = 5C3 = Answer: 35
//
// Note: 0! = 1
//
func (prob Probability) CombinationsWithRepsBigInt(
	numOfItems, numOfItemsChosen *big.Int) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsWithRepsBigInt() "
	bigZero := big.NewInt(0)

	if numOfItems.Cmp(bigZero) == 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.Cmp(bigZero) < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.Cmp(bigZero) < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.Cmp(bigZero) == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: 'numOfItemsChosen' is ZERO! ")
	}

	bigOne := big.NewInt(1)

	// If 'numOfItemsChosen' == 1, result is always equal to 'numOfItems'.
	if numOfItemsChosen.Cmp(bigOne) == 0 {
		return BigIntNum{}.NewBigInt(numOfItems, 0), nil
	}

	temp1 := big.NewInt(0).Add(numOfItems, numOfItemsChosen)

	numeratorUpperLimit := big.NewInt(0).Sub(temp1, bigOne)
	numeratorLowerLimit := big.NewInt(0).Sub(numOfItems, bigOne)

	nFactorial, err := NFactorial{}.CalcFactorialValueBigInt(
		numeratorUpperLimit, numeratorLowerLimit)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
				"numeratorUpperLimit, numeratorLowerLimit). "+
				"numeratorNUpperLimit='%v'  numeratorNLowerLimit='%v' Error='%v'",
				numeratorUpperLimit.Text(10), numeratorLowerLimit.Text(10),
				err.Error())
	}

	rUpperLimit := big.NewInt(0).Set(numOfItemsChosen)
	rLowerLimit := big.NewInt(1)

	rFactorial, err := NFactorial{}.CalcFactorialValueBigInt(
		rUpperLimit, rLowerLimit)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by NFactorial{}.CalcFactorialValueBigInt("+
				"rUpperLimit, rLowerLimit). "+
				"rUpperLimit='%v'  rLowerLimit='%v' Error='%v'",
				rUpperLimit.Text(10), rLowerLimit.Text(10),
				err.Error())
	}

	combinationsResult, err :=
		BigIntMathDivide{}.BigIntNumFracQuotient(nFactorial, rFactorial, 10)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntMathDivide{}.BigIntNumFracQuotient("+
				"nFactorial, rFactorial, 10)). "+
				"nFactorial='%v'  rFactorial='%v' Error='%v'",
				nFactorial.GetNumStr(), rFactorial.GetNumStr(),
				err.Error())
	}

	return combinationsResult, nil
}

// CombinationsBigIntNum - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'BigIntNum'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsBigIntNum(
	numOfItems, numOfItemsChosen BigIntNum, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsBigIntNum() "

	if numOfItems.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsPicked' is ZERO!")
	}

	if numOfItems.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsPicked' is LESS THAN ZERO!")
	}

	if numOfItems.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if numOfItemsChosen.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if !allowRepetitions && numOfItems.Cmp(numOfItemsChosen) < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: 'numOfItems' is LESS THAN 'numOfItemsPicked'! ")

	}

	if !allowRepetitions {
		return Probability{}.CombinationsNoRepsBigInt(numOfItems.bigInt, numOfItemsChosen.bigInt)
	}

	return Probability{}.CombinationsWithRepsBigInt(numOfItems.bigInt, numOfItemsChosen.bigInt)
}

// CombinationsDecimal - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'Decimal'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'Decimal'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsDecimal(
	numOfItems, numOfItemsChosen Decimal, allowRepetitions bool) (Decimal, error) {

	ePrefix := "Probability.CombinationsDecimal() "

	if numOfItems.GetSign() == -1 {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems.IsZero() {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.GetPrecisionUint() > 0 {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if numOfItemsChosen.GetSign() == -1 {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.IsZero() {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	if numOfItemsChosen.GetPrecisionUint() > 0 {
		return Decimal{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is NOT an Integer!")
	}

	n, err := numOfItems.GetBigInt()

	if err != nil {
		return Decimal{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItems.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	r, err := numOfItemsChosen.GetBigInt()

	if err != nil {
		return Decimal{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItemsChosen.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	if !allowRepetitions && n.Cmp(r) < 0 {
		return Decimal{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItems' is LESS THAN 'numOfItemsChosen'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems.GetNumStr(), numOfItemsChosen.GetNumStr())
	}

	var result BigIntNum

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return Decimal{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return Decimal{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	resultDecimal, err := result.GetDecimal()

	if err != nil {
		return Decimal{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by result.GetDecimal(). "+
				"Error='%v' \n", err.Error())
	}

	return resultDecimal, nil
}

// CombinationsIntAry - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'IntAry'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'IntAry'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsIntAry(
	numOfItems, numOfItemsChosen IntAry, allowRepetitions bool) (IntAry, error) {

	ePrefix := "Probability.CombinationsIntAry() "

	if numOfItems.GetSign() == -1 {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems.IsZero() {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.GetPrecisionUint() > 0 {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if numOfItemsChosen.GetSign() == -1 {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.IsZero() {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	if numOfItemsChosen.GetPrecisionUint() > 0 {
		return IntAry{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is NOT an Integer!")
	}

	n, err := numOfItems.GetBigInt()

	if err != nil {
		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItems.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	r, err := numOfItemsChosen.GetBigInt()

	if err != nil {
		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItemsChosen.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	if !allowRepetitions && r.Cmp(n) == 1 {
		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems.GetNumStr(), numOfItemsChosen.GetNumStr())
	}

	var result BigIntNum

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return IntAry{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return IntAry{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	resultIntAry, err := result.GetIntAry()

	if err != nil {
		return IntAry{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by result.GetIntAry(). "+
				"Error='%v' \n", err.Error())
	}

	return resultIntAry, nil
}

// CombinationsINumMgr - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'INumMgr'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'INumMgr'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsINumMgr(
	numOfItems, numOfItemsChosen INumMgr, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsINumMgr() "

	if numOfItems.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if numOfItemsChosen.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	if numOfItemsChosen.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is NOT an Integer!")
	}

	n, err := numOfItems.GetBigInt()

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItems.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	r, err := numOfItemsChosen.GetBigInt()

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItemsChosen.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems.GetNumStr(), numOfItemsChosen.GetNumStr())
	}

	var result BigIntNum

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsInt - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'int'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsInt(
	numOfItems, numOfItemsChosen int, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsInt() "

	if numOfItems < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(int64(numOfItems))

	r := big.NewInt(int64(numOfItemsChosen))

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsInt32 - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'int32'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsInt32(
	numOfItems, numOfItemsChosen int32, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsInt32() "

	if numOfItems < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(int64(numOfItems))

	r := big.NewInt(int64(numOfItemsChosen))

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt32(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsInt64 - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'int32'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsInt64(
	numOfItems, numOfItemsChosen int64, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsInt64() "

	if numOfItems < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen < 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(numOfItems)

	r := big.NewInt(numOfItemsChosen)

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt64(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsNumStrDto - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'NumStrDto'. Both input parameters must
// be non-zero, positive integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'NumStrDto'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsNumStrDto(
	numOfItems, numOfItemsChosen NumStrDto, allowRepetitions bool) (NumStrDto, error) {

	ePrefix := "Probability.CombinationsNumStrDto() "

	if numOfItems.GetSign() == -1 {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if numOfItems.IsZero() {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItems.GetPrecisionUint() > 0 {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	if numOfItemsChosen.GetSign() == -1 {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if numOfItemsChosen.IsZero() {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	if numOfItemsChosen.GetPrecisionUint() > 0 {
		return NumStrDto{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is NOT an Integer!")
	}

	n, err := numOfItems.GetBigInt()

	if err != nil {
		return NumStrDto{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItems.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	r, err := numOfItemsChosen.GetBigInt()

	if err != nil {
		return NumStrDto{}.NewZero(0),
			fmt.Errorf(ePrefix+"Error returned by numOfItemsChosen.GetBigInt(). "+
				"Error='%v' \n", err.Error())
	}

	if !allowRepetitions && r.Cmp(n) == 1 {
		return NumStrDto{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems.GetNumStr(), numOfItemsChosen.GetNumStr())
	}

	var result BigIntNum

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return NumStrDto{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return NumStrDto{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	resultNumStrDto, err := result.GetNumStrDto()

	if err != nil {
		return NumStrDto{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by result.GetNumStrDto(). "+
				"Error='%v' \n", err.Error())
	}

	return resultNumStrDto, nil
}

// CombinationsBigIntNum - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant.
//
// Input parameters 'numOfItems' and 'numOfItemsChosen' are passed as strings. These strings
// must be formatted as valid number strings. Number strings may be prefixed by a plus (+) or
// minus (-) and must consist of a string of numeric digits which may be delimited by a thousands
// separator. If the numeric value is a fractional value, the fractional digits must be
// preceded by a period ('.') or decimal separator. However, for purposes of this calculation,
// both input parameters must be passed as positive. integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be positive integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsNumberStr(
	numOfItems, numOfItemsChosen string, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsNumberStr() "

	if numOfItems == "" {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is an EMPTY string!")
	}

	if numOfItemsChosen == "" {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is an EMPTY string!")
	}

	nBigIntNum, err := BigIntNum{}.NewNumStr(numOfItems)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStr(numOfItems). "+
				"Error='%v' ", err.Error())
	}

	if nBigIntNum.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if nBigIntNum.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is LESS THAN ZERO!")
	}

	if nBigIntNum.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is NOT an Integer!")
	}

	rBigIntNum, err := BigIntNum{}.NewNumStr(numOfItemsChosen)

	if err != nil {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error returned by BigIntNum{}.NewNumStr(numOfItemsChosen). "+
				"Error='%v' ", err.Error())
	}

	if rBigIntNum.IsZero() {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	if rBigIntNum.GetSign() == -1 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is LESS THAN ZERO!")
	}

	if rBigIntNum.GetPrecisionUint() > 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is NOT an Integer!")
	}

	if !allowRepetitions && nBigIntNum.Cmp(rBigIntNum) < 0 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItems' is LESS THAN 'numOfItemsChosen'! "+
				"numOfItems='%v' numOfItemsChosen='%v' \n", numOfItems, numOfItemsChosen)
	}

	if !allowRepetitions {
		return Probability{}.CombinationsNoRepsBigInt(nBigIntNum.bigInt, rBigIntNum.bigInt)
	}

	return Probability{}.CombinationsWithRepsBigInt(nBigIntNum.bigInt, rBigIntNum.bigInt)
}

// CombinationsUint - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'uint'. Both input parameters must
// be non-zero, integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsUint(
	numOfItems, numOfItemsChosen uint, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsUint() "

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(0).SetUint64(uint64(numOfItems))

	r := big.NewInt(0).SetUint64(uint64(numOfItemsChosen))

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsUint32 - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'uint32'. Both input parameters must
// be non-zero, integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsUint32(
	numOfItems, numOfItemsChosen uint32, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsUint32() "

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(0).SetUint64(uint64(numOfItems))

	r := big.NewInt(0).SetUint64(uint64(numOfItemsChosen))

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}

// CombinationsUint64 - Calculates the number of combinations associated with a collection of
// 'numOfItems' from which one chooses 'numOfItemsChosen'. Order IS NOT significant. Input parameters
// 'numOfItems' and 'numOfItemsChosen' are passed as type 'uint64'. Both input parameters must
// be non-zero, integer numbers.
//
// The input parameter 'allowRepetitions' is a boolean value which will determine whether the
// calculation results will allow repetitions or not. The formula for the calculation of combinations
// will therefore vary depending on whether repetitions are allowed. 'allowRepetitions' == false signals
// unordered sampling WITHOUT replacement. 'allowRepetitions' == true signals unordered sampling WITH
// replacement.
//
// The result of this combination calculation is returned as a type 'BigIntNum'.
//
// In the following combination formulas, n= 'numOfItems'  and r = 'numOfItemsChosen'. The actual
// formula applied depends on whether input parameter 'allowRepetitions' is 'true' or 'false'.
//
//      ====================================================================
// 			'allowRepetitions' = false
//      ====================================================================
//
// 																			 n!
//											 nCr	 =		-----------
//																	 (n-r)! r!
//
//        --------------------------------------------------
//
// 						Where n is the number of things to choose from,
//						and we choose r of them, order does NOT matter
//            and repetition is NOT allowed.
//
// *** This version of the combinations calculation assumes NO REPETITIONS! ***
//    (a.k.a. as unordered sampling WITHOUT replacement)
//
//      When 'allowRepetitions' = false, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. In addition, 'numOfItems' MUST be greater than or
// 			equal to 'numOfItemsChosen'.
//
//      ====================================================================
// 			'allowRepetitions' = true
//      ====================================================================
//
//
// 															(r + n - 1)!
// 								nCr    =  -------------------
// 											 				 r! (n-1)!
//
// 				Where n is the number of things to choose from
//				and we choose r of them, order does NOT matter
// 				and repetition IS NOT allowed.
//
// *** This version of the combinations calculation assumes REPETITIONS ARE ALLOWED! ***
//    (a.k.a. as unordered sampling WITH replacement)
//
//      When 'allowRepetitions' = true, 'numOfItems' and 'numOfItemsChosen' must both
// 			be integer numbers. 'numOfItems' can be greater than, equal to or less
// 			than 'numOfItemsChosen'.
//
func (prob Probability) CombinationsUint64(
	numOfItems, numOfItemsChosen uint64, allowRepetitions bool) (BigIntNum, error) {

	ePrefix := "Probability.CombinationsUint64() "

	if numOfItems == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItems' is ZERO!")
	}

	if numOfItemsChosen == 0 {
		return BigIntNum{}.NewZero(0),
			errors.New(ePrefix + "Error: Input parameter 'numOfItemsChosen' is ZERO!")
	}

	n := big.NewInt(0).SetUint64(numOfItems)

	r := big.NewInt(0).SetUint64(numOfItemsChosen)

	if !allowRepetitions && r.Cmp(n) == 1 {
		return BigIntNum{}.NewZero(0),
			fmt.Errorf(ePrefix+
				"Error: 'numOfItemsChosen' is GREATER THAN 'numOfItems'.  "+
				"numOfItems='%v' numOfItemsChosen='%v' \n",
				numOfItems, numOfItemsChosen)
	}

	var result BigIntNum
	var err error

	if !allowRepetitions {

		result, err = Probability{}.CombinationsNoRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsNoRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	} else {

		result, err = Probability{}.CombinationsWithRepsBigInt(n, r)

		if err != nil {
			return BigIntNum{}.NewZero(0),
				fmt.Errorf(ePrefix+
					"Error returned by Probability{}.CombinationsWithRepsBigInt(numOfItems, numOfItemsChosen). "+
					"Error='%v' \n", err.Error())
		}

	}

	return result, nil
}
