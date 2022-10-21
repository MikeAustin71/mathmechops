package mathops

type IntAryPair struct {
	IAry1             IntAry
	IAry1Compare      int // 	1 = IAry1 > IAry2; 0 = IAry1 == IAry2; -1 = IAry1 < IAry2
	IAry1AbsCompare   int // 	1 = IAry1 > IAry2; 0 = IAry1 == IAry2; -1 = IAry1 < IAry2
	Precision1Compare int // 	1 = IAry1Precision > IAry2Precision;
	//  0 = Big1Precision == Big2Precision;
	// -1 = Big1Precision < Big2Precision
	IAry2        IntAry
	MaxPrecision int // Used to control output from complex
	//  calculations
}
