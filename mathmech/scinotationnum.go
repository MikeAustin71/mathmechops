package mathops

import (
	"errors"
	"fmt"
	"strings"
)

// SciNotationNum
//
// # Example Scientific Notation
//
//	 2.652e+8
//
//	 significand 		= '2.652'
//	 significand integer digits = '2'
//		mantissa				= significand factional digits = '.652'
//	 exponent    		= '8'  (10^8)
//		mantissaLength	= length of fractional digits displayed in scientific notation.
type SciNotationNum struct {
	significand BigIntNum // The significand consists of the leading integer and
	//	fractional digits of the scientific notation.
	exponent         BigIntNum // The exponent portion of the scientific notation string
	exponentChar     rune      // 	defaults to 'e'. May be customized to 'E'
	decimalSeparator rune      // The decimal separator used to separate integer and
	// 	fractional digits in the significand. The default is
	// 	the standard USA decimal separator, the decimal point ('.').
	mantissaLength uint // The length of the fractional digits in
	// 	the significand which will be displayed
	// 	when SciNotationNum.GetSciNotationStr()
	// 	is called.
	exponentUsesLeadingPlus bool // If true, positive exponent values are
	// 	prefixed with a leading plus (+) sign.
	//  '2.652e+8'
}

// GetDecimalSeparator - returns the current decimal separator
// character as a string.
//
// In a default scientific notation display, '2.652e+8',
// the decimal separator character is presented as a period
// ('.') separating integer and fractional digits in the
// significand ('2.652'). However, the user has the option to
// customize this decimal separator character through method
// SciNotationNum.SetDecimalSeparatorChar().
func (sciNotan *SciNotationNum) GetDecimalSeparator() rune {

	if sciNotan.decimalSeparator == 0 {
		sciNotan.SetDecimalSeparatorIfEmpty()
	}

	return sciNotan.decimalSeparator
}

// GetExponent - Returns the exponent element of the scientific
// notation as a BigIntNum.
// Example Scientific Notation
// ===========================
//
//	 2.652e+8
//
//	 significand 		= '2.652'
//	 significand integer digits = '2'
//		mantissa				= significand factional digits = '.652'
//	 exponent    		= '8'  (10^8)
//		mantissaLength	= length of fractional digits displayed in scientific notation.
func (sciNotan *SciNotationNum) GetExponent() BigIntNum {

	return sciNotan.exponent.CopyOut()
}

// GetExponentChar - Returns the current exponent char
// as a string.
//
// In a default scientific notation display, '2.652e+8',
// the exponent char is presented as an 'e'. However, the
// user has the option to customize this exponent character
// through method SciNotationNum.SetExponentChar().
//
// This method returns the current exponent character which
// will be used in formatting scientific notation strings.
func (sciNotan *SciNotationNum) GetExponentChar() string {

	if sciNotan.exponentChar == 0 {
		sciNotan.SetExponentCharIfEmpty()
	}

	return string(sciNotan.exponentChar)
}

// GetExponentUsesLeadingPlus - returns the internal 'exponentUsesLeadingPlus'
// flag. By default this is set to 'true'.  When set to 'true' the scientific
// notation is displayed with positive exponents prefixed with a plus sign (+).
//
// Default Example
// ---------------
// 2.652e+8
func (sciNotan *SciNotationNum) GetExponentUsesLeadingPlus() bool {

	return sciNotan.exponentUsesLeadingPlus
}

// GetNumStr - Returns a scientific notation string representing
// the underlying numeric value. The number of decimals in the
// mantissa will default to the current value of sciNotan.mantissaLength.
func (sciNotan *SciNotationNum) GetNumStr() string {

	sciNotan.SetMantissaLengthIfEmpty()

	result, _ := sciNotan.GetSciNotationStr(sciNotan.mantissaLength)

	return result
}

// GetSciNotationStr - Returns a string containing the scientific
// notation display.  This method differs from method SciNotationNum.GetNumStr()
// in that this method requires the user to provide input parameter 'mantissaLen'
// which sets the number of decimal places in the significand of the returned
// science notation string.
//
// Default Example
// ---------------
// 2.652e+8
func (sciNotan *SciNotationNum) GetSciNotationStr(mantissaLen uint) (string, error) {

	outStr := ""

	if mantissaLen == 0 {
		mantissaLen = 2
	}

	sciNotan.SetMantissaLength(mantissaLen)

	sciNotan.SetDecimalSeparatorIfEmpty()

	sciNotan.significand.SetDecimalSeparator(sciNotan.decimalSeparator)

	if sciNotan.significand.GetPrecisionUint() != sciNotan.mantissaLength {
		bINum := sciNotan.significand.CopyOut()

		bINum.SetPrecision(sciNotan.mantissaLength)

		outStr += bINum.GetNumStr()

	} else {

		outStr += sciNotan.significand.GetNumStr()
	}

	outStr += sciNotan.GetExponentChar()

	if sciNotan.exponent.GetSign() == 1 &&
		sciNotan.exponentUsesLeadingPlus {
		outStr += "+"
	}

	outStr += sciNotan.exponent.GetNumStr()

	return outStr, nil
}

// GetSignificand - Returns the significand component of the
// scientific notation as a BigIntNum.
//
// Example Scientific Notation
// ===========================
//
//	 2.652e+8
//
//	 significand 		= '2.652'
//	 significand integer digits = '2'
//		mantissa				= significand factional digits = '.652'
//	 exponent    		= '8'  (10^8)
//		mantissaLength	= length of fractional digits displayed in scientific notation.
func (sciNotan *SciNotationNum) GetSignificand() BigIntNum {

	return sciNotan.significand.CopyOut()
}

// New() - Creates and returns an empty SciNotationNum
// structure. It is a good idea to call this method
// in order to initialize default settings.
func (sciNotan SciNotationNum) New() SciNotationNum {

	s2 := SciNotationNum{}

	s2.SetExponentCharIfEmpty()
	s2.SetDecimalSeparatorIfEmpty()

	s2.significand = BigIntNum{}.NewZero(1)
	s2.exponent = BigIntNum{}.NewZero(0)
	s2.exponentUsesLeadingPlus = true

	return s2
}

// NewNumStr - Creates and returns a new SciNotationNum instance initialized from
// an input string.
//
// Input parameter 'sciNotationStr' should be properly formatted as a valid scientific
// notation string. Invalid input strings will trigger an error.
func (sciNotan SciNotationNum) NewNumStr(sciNotationStr string) (SciNotationNum, error) {

	s2 := SciNotationNum{}.New()

	err := s2.SetNumStr(sciNotationStr)

	if err != nil {
		ePrefix := "SciNotationNum.NewNumStr() "
		return SciNotationNum{}.New(),
			fmt.Errorf(ePrefix+
				"Error returned by s2.SetNumStr(sciNotationStr). Error='%v'\n", err.Error())
	}

	return s2, nil
}

// SetDecimalSeparatorIfEmpty() - Sets the default decimal separator
// if the SciNotationNum.decimalSeparator rune is empty or equal to zero.
//
// In the example scientific notation '2.652e+8', the decimal separator is
// the period or decimal point ('.'). If SciNotationNum.decimalSeparator rune
// is empty, this method will set the decimal separator to the USA default
// decimal separator, the period or decimal point ('.').
//
// The decimal separator character may be customized to support characters
// used by other cultures or nations. See method SciNotationNum.SetDecimalSeparatorChar()
// below.
func (sciNotan *SciNotationNum) SetDecimalSeparatorIfEmpty() {

	if sciNotan.decimalSeparator == 0 {
		sciNotan.decimalSeparator = '.'
		sciNotan.significand.SetDecimalSeparator('.')
	}

}

// SetDecimalSeparatorChar - Sets the value of the Scientific Notation decimal
// separator used to separate integer and fractional digits in the significand.
//
// In a default scientific notation display, '2.652e+8', the decimal separator
// character is presented as a period ('.') separating integer and fractional
// digits in the significand ('2.652'). However, the user has the option to
// customize this decimal separator character by calling this method.
func (sciNotan *SciNotationNum) SetDecimalSeparatorChar(decimalChar rune) {

	sciNotan.decimalSeparator = decimalChar
	sciNotan.significand.SetDecimalSeparator(decimalChar)
}

// SetExponentCharIfEmpty - If the exponent rune is empty,
// this method will set the exponentChar value to 'e'.
func (sciNotan *SciNotationNum) SetExponentCharIfEmpty() {

	if sciNotan.exponentChar == 0 {
		sciNotan.exponentChar = 'e'
	}

}

// SetExponentChar - The exponent char used in displaying scientific notation
// may be customized through this method.
//
// By default, scientific notation is displayed in the following format using the
// 'e' character as the exponentChar: 2.652e+8
//
// The user may find it preferable to display the scientific notation using capital
// 'E': 2.652E+8
//
// This method is used to customized the exponent character used when calling
// SciNotationNum.GetSciNotationStr()
func (sciNotan *SciNotationNum) SetExponentChar(exponentChar rune) {
	sciNotan.exponentChar = exponentChar
}

func (sciNotan *SciNotationNum) SetExponentUsesLeadingPlus(useLeadingPlus bool) {

	sciNotan.exponentUsesLeadingPlus = useLeadingPlus

}

// SetMantissaLength - This method sets the length of the mantissa or
// fractional digits which will be displayed in the the significand
// when SciNotationNum.GetSciNotationStr() is called.
//
// If input parameter mantissaLen is set equal to zero, this method will
// automatically set the value to two ('2').
//
// Example Scientific Notation
// ===========================
//
//	 2.652e+8
//
//	 significand = '2.652'
//	 significand integer digit = '2'
//		mantissa		= significand factional digits = '.652'
//	 exponent    = '8'  (10^8)
func (sciNotan *SciNotationNum) SetMantissaLength(mantissaLen uint) {

	if mantissaLen == 0 {
		mantissaLen = 2
	}

	sciNotan.mantissaLength = mantissaLen

}

// SetMantissaLengthIfEmpty - If mantissa length is zero, this
// method attempts to set mantissa length equal to the precision
// of 'significand'.
//
// Mantissa is defined as the length or number of fractional digits
// which will be displayed in the the significand when
// SciNotationNum.GetSciNotationStr() is called to produce scientific
// notation as a string.
//
// In the example scientific notation '2.652e+8', the mantissa is '.652'
func (sciNotan *SciNotationNum) SetMantissaLengthIfEmpty() {

	if sciNotan.mantissaLength == 0 {

		if sciNotan.significand.GetPrecisionUint() == 0 {
			sciNotan.mantissaLength = 2
		} else {
			sciNotan.mantissaLength = sciNotan.significand.GetPrecisionUint()
		}

	}

}

// SetBigIntNumElements - Sets the components of the current SciNotationNum
// instance based on two BigIntNum input parameters.
//
// Input Parameters
// ================
//
// significand 	BigIntNum 	- In the example scientific notation '2.652e+8',
//
//	the significand is represented by '2.652'.
//
// exponent  		BigIntNum		- In the example '2.652e+8' the exponent component
//
//	is represented by the integer value, '8'. Note:
//	If exponent is NOT an integer value and contains
//	fractional digits, an error will be triggered.
func (sciNotan *SciNotationNum) SetBigIntNumElements(
	significand, exponent BigIntNum) error {

	ePrefix := "SciNotationNum.SetBigIntNumElements() "

	if exponent.GetPrecisionUint() > 0 {
		return fmt.Errorf(ePrefix +
			"Error: Input parameter 'exponent' contains fractional digits!")
	}

	sciNotan.SetExponentCharIfEmpty()

	sciNotan.significand.CopyIn(significand)
	sciNotan.decimalSeparator = significand.GetDecimalSeparator()
	sciNotan.mantissaLength = sciNotan.significand.GetPrecisionUint()
	sciNotan.exponent.CopyIn(exponent)

	return nil
}

// SetIntAryElements - Sets the components of the current SciNotationNum
// instance based on two IntAry input parameters.
//
// Input Parameters
// ================
//
// significand 	IntAry 			- In the example scientific notation '2.652e+8',
//
//	the significand is represented by '2.652'.
//
// exponent  		IntAry			- In the example '2.652e+8' the exponent component
//
//	is represented by the integer value, '8'. Note:
//	If exponent is NOT an integer value and contains
//	fractional digits, an error will be triggered.
func (sciNotan *SciNotationNum) SetIntAryElements(
	significand, exponent IntAry) error {

	ePrefix := "SciNotationNum.SetBigIntNumElements() "

	if exponent.GetPrecisionUint() > 0 {
		return fmt.Errorf(ePrefix +
			"Error: Input parameter 'exponent' contains fractional digits!")
	}

	sciNotan.SetExponentCharIfEmpty()

	biNumSignificand, err := significand.GetBigIntNum()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by significand.GetBigIntNum(). "+
			"Error='%v'", err.Error())
	}

	sciNotan.significand = biNumSignificand.CopyOut()
	sciNotan.decimalSeparator = significand.GetDecimalSeparator()
	sciNotan.mantissaLength = sciNotan.significand.GetPrecisionUint()

	biNumExponent, err := exponent.GetBigIntNum()

	if err != nil {
		return fmt.Errorf(ePrefix+"Error returned by exponent.GetBigIntNum(). "+
			"Error='%v'", err.Error())
	}

	sciNotan.exponent.CopyIn(biNumExponent)

	return nil
}

// SetNumStr - Receives a properly formatted scientific notation string as
// an input parameter and converts to the data fields of the current
// SciNotationNum instance.
//
// Examples of properly formatted scientific notation strings are provided
// below:
//
//	2.652e+8
//	2.652E+8
//	2.652e8
//	2.652E8
//
// The use of the decimal point may be customize by first setting the desired
// decimal separator using method, SciNotationNum.SetDecimalSeparatorChar().
//
// Also, note that exponent digits must be integer numbers. Used of fractional
// digits in the exponent will trigger an error. Example:
//
//	2.652E9.24 = ERROR fractional digits in exponent!
func (sciNotan *SciNotationNum) SetNumStr(sciNotationStr string) error {

	ePrefix := "BigIntNum.SetNumStr() "

	if len(sciNotationStr) == 0 {
		return errors.New(ePrefix +
			"Error: Input parameter 'sciNotationStr' is an EMPTY string!")
	}

	sciNotan.SetDecimalSeparatorIfEmpty()
	sciNotan.SetExponentCharIfEmpty()

	i := strings.Index(sciNotationStr, "e")

	if i == -1 {
		i = strings.Index(sciNotationStr, "E")
	}

	if i == -1 {
		return errors.New(ePrefix +
			"Error: Input parameter 'sciNotationStr' does NOT contain an " +
			"Exponent Character ('e' or 'E') ")
	}

	significandStr := sciNotationStr[:i]

	if significandStr == "" {
		return errors.New(ePrefix +
			"Error: Input parameter 'sciNotationStr' does NOT contain any " +
			"digits in the significand! ")

	}

	exponentStr := sciNotationStr[i+1:]

	if exponentStr == "" {
		return errors.New(ePrefix +
			"Error: Input parameter 'sciNotationStr' does NOT contain any " +
			"digits in the exponent! ")

	}

	sciNotan.SetDecimalSeparatorIfEmpty()

	bINumSignificand := BigIntNum{}.New()
	bINumSignificand.SetDecimalSeparator(sciNotan.decimalSeparator)

	err := bINumSignificand.SetNumStr(significandStr)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by bINumSignificand.SetNumStr(significandStr). "+
			"significand='%v' Error='%v'\n", significandStr, err.Error())
	}

	bINumExponent, err := BigIntNum{}.NewNumStr(exponentStr)

	if err != nil {
		return fmt.Errorf(ePrefix+
			"Error returned by BigIntNum{}.NewNumStr(exponentStr). "+
			"exponentStr='%v' Error='%v'\n", exponentStr, err.Error())
	}

	if bINumExponent.GetPrecisionUint() > 0 {
		return errors.New(ePrefix +
			"Error: The exponent component of the input parameter 'sciNotationStr' " +
			"contains fractional digits!")

	}

	sciNotan.significand.CopyIn(bINumSignificand)
	sciNotan.SetMantissaLength(bINumSignificand.GetPrecisionUint())
	sciNotan.exponent.CopyIn(bINumExponent)

	return nil
}
