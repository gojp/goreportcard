package check

// Grade represents a grade returned by the server, which is normally
// somewhere between A+ (highest) and F (lowest).
type Grade string

// The Grade constants below indicate the current available
// grades.
const (
	GradeAPlus Grade = "A+"
	GradeA           = "A"
	GradeB           = "B"
	GradeC           = "C"
	GradeD           = "D"
	GradeE           = "E"
	GradeF           = "F"
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
