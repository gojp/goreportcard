package handlers

import (
	"net/url"
	"testing"

	"github.com/gojp/goreportcard/check"
)

func TestBadgeURL(t *testing.T) {
	cases := []struct {
		name   string
		grade  check.Grade
		params url.Values
		want   string
	}{
		{
			check.GradeAPlus,
			check.GradeAPlus,
			url.Values{},
			"https://img.shields.io/static/v1?color=brightgreen&label=go+report&message=A%2B&style=flat",
		},
		{
			check.GradeA,
			check.GradeA,
			url.Values{},
			"https://img.shields.io/static/v1?color=green&label=go+report&message=A&style=flat",
		},
		{
			check.GradeB,
			check.GradeB,
			url.Values{},
			"https://img.shields.io/static/v1?color=yellowgreen&label=go+report&message=B&style=flat",
		},
		{
			check.GradeC,
			check.GradeC,
			url.Values{},
			"https://img.shields.io/static/v1?color=yellow&label=go+report&message=C&style=flat",
		},
		{
			check.GradeD,
			check.GradeD,
			url.Values{},
			"https://img.shields.io/static/v1?color=orange&label=go+report&message=D&style=flat",
		},
		{
			check.GradeE,
			check.GradeE,
			url.Values{},
			"https://img.shields.io/static/v1?color=red&label=go+report&message=E&style=flat",
		},
		{
			check.GradeF,
			check.GradeF,
			url.Values{},
			"https://img.shields.io/static/v1?color=red&label=go+report&message=F&style=flat",
		},
		{
			"override style",
			check.GradeAPlus,
			url.Values{"style": []string{"for-the-badge"}},
			"https://img.shields.io/static/v1?color=brightgreen&label=go+report&message=A%2B&style=for-the-badge",
		},
		{
			"override color",
			check.GradeAPlus,
			url.Values{"color": []string{"ff69b4"}},
			"https://img.shields.io/static/v1?color=ff69b4&label=go+report&message=A%2B&style=flat",
		},
		{
			"override logo params",
			check.GradeAPlus,
			url.Values{"logo": []string{"go"}, "logoWidth": []string{"100"}, "logoColor": []string{"ff69b4"}},
			"https://img.shields.io/static/v1?color=brightgreen&label=go+report&logo=go&logoColor=ff69b4&logoWidth=100&message=A%2B&style=flat",
		},
		{
			"override label params",
			check.GradeAPlus,
			url.Values{"label": []string{"code quality"}, "labelColor": []string{"0080ff"}},
			"https://img.shields.io/static/v1?color=brightgreen&label=code+quality&labelColor=0080ff&message=A%2B&style=flat",
		},
	}

	for _, tt := range cases {
		tt := tt

		t.Run(string(tt.name), func(t *testing.T) {
			t.Parallel()

			got := badgeURL(tt.grade, tt.params)
			if got != tt.want {
				t.Errorf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
