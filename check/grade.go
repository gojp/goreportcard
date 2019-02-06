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

// BadgeFromGrade is a helper for getting the badge svg for a grade
func BadgeFromGrade(grade Grade) string {
	switch grade {
	case GradeAPlus:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="88" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="88" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#4c1" d="M61 0h27v20H61z"/><path fill="url(#b)" d="M0 0h88v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="73.5" y="15" fill="#010101" fill-opacity=".3">A+</text><text x="73.5" y="14">A+</text></g></svg>`
	case GradeA:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="78" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#4c1" d="M61 0h17v20H61z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="68.5" y="15" fill="#010101" fill-opacity=".3">A</text><text x="68.5" y="14">A</text></g></svg>`
	case GradeB:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="78" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#a4a61d" d="M61 0h17v20H61z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="68.5" y="15" fill="#010101" fill-opacity=".3">B</text><text x="68.5" y="14">B</text></g></svg>`
	case GradeC:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="78" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#dfb317" d="M61 0h17v20H61z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="68.5" y="15" fill="#010101" fill-opacity=".3">C</text><text x="68.5" y="14">C</text></g></svg>`
	case GradeD:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="80" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="80" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#fe7d37" d="M61 0h19v20H61z"/><path fill="url(#b)" d="M0 0h80v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="69.5" y="15" fill="#010101" fill-opacity=".3">D</text><text x="69.5" y="14">D</text></g></svg>`
	case GradeE:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="78" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#e05d44" d="M61 0h17v20H61z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="68.5" y="15" fill="#010101" fill-opacity=".3">E</text><text x="68.5" y="14">E</text></g></svg>`
	case GradeF:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="78" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h61v20H0z"/><path fill="#e05d44" d="M61 0h17v20H61z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="30.5" y="15" fill="#010101" fill-opacity=".3">go report</text><text x="30.5" y="14">go report</text><text x="68.5" y="15" fill="#010101" fill-opacity=".3">F</text><text x="68.5" y="14">F</text></g></svg>`
	default:
		return ""
	}
}
