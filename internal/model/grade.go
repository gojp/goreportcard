package model

// Grade represents a grade returned by the server, which is normally
// somewhere between A+ (highest) and F (lowest).
type Grade string

// The Grade constants below indicate the current available
// grades.
const (
	// GradeAPlus  means "A+"
	GradeAPlus Grade = "A+"
	// GradeA means "A"
	GradeA = "A"
	// GradeB means "B"
	GradeB = "B"
	// GradeC means "C"
	GradeC = "C"
	// GradeD means "D"
	GradeD = "D"
	// GradeE means "E"
	GradeE = "E"
	// GradeF means "F"
	GradeF = "F"
)

// GradeFromPercentage is a helper for getting the GradeFromPercentage for a percentage
func GradeFromPercentage(percentage float64) Grade {
	switch {
	case percentage > 90:
		return GradeAPlus
	case percentage > 80:
		return GradeA
	case percentage > 70:
		return GradeB
	case percentage > 60:
		return GradeC
	case percentage > 50:
		return GradeD
	case percentage > 40:
		return GradeE
	default:
		return GradeF
	}
}
