package stringtool

import (
	"testing"
)

func TestToExportedName(t *testing.T) {
	for i, c := range []struct {
		src    string
		expect string
	}{
		{"IGotInternAtGeeksForGeeks", "IGotInternAtGeeksForGeeks"},
		{"iGotInternAtGeeksForGeeks", "IGotInternAtGeeksForGeeks"},
		{"i-got-intern-at-geeks-for-geeks", "IGotInternAtGeeksForGeeks"},
		{"i_got_intern_at_geeks_for_geeks", "IGotInternAtGeeksForGeeks"},
	} {
		actual := toExportedName(c.src)
		if actual != c.expect {
			t.Fatalf("%5d. expert %q -> %q but got %q", i, c.src, c.expect, actual)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	for i, c := range []struct {
		src    string
		expect string
	}{
		{"IGotInternAtGeeksForGeeks", "IGotInternAtGeeksForGeeks"},
		{"iGotInternAtGeeksForGeeks", "IGotInternAtGeeksForGeeks"},
		{"i-got-intern-at-geeks-for-geeks", "IGotInternAtGeeksForGeeks"},
		{"i_got_intern_at_geeks_for_geeks", "IGotInternAtGeeksForGeeks"},
	} {
		actual := ToCamelCase(c.src)
		if actual != c.expect {
			t.Fatalf("%5d. expert %q -> %q but got %q", i, c.src, c.expect, actual)
		}
	}
}

func TestToSmallCamelCase(t *testing.T) {
	for i, c := range []struct {
		src    string
		expect string
	}{
		{"IGotInternAtGeeksForGeeks", "iGotInternAtGeeksForGeeks"},
		{"iGotInternAtGeeksForGeeks", "iGotInternAtGeeksForGeeks"},
		{"i-got-intern-at-geeks-for-geeks", "iGotInternAtGeeksForGeeks"},
		{"i_got_intern_at_geeks_for_geeks", "iGotInternAtGeeksForGeeks"},
	} {
		actual := ToSmallCamelCase(c.src)
		if actual != c.expect {
			t.Fatalf("%5d. expert %q -> %q but got %q", i, c.src, c.expect, actual)
		}
	}
}

func TestToKebabCase(t *testing.T) {
	for i, c := range []struct {
		src    string
		expect string
	}{
		{"IGotInternAtGeeksForGeeks", "i-got-intern-at-geeks-for-geeks"},
		{"iGotInternAtGeeksForGeeks", "i-got-intern-at-geeks-for-geeks"},
		{"i-got-intern-at-geeks-for-geeks", "i-got-intern-at-geeks-for-geeks"},
		{"i_got_intern_at_geeks_for_geeks", "i-got-intern-at-geeks-for-geeks"},
	} {
		actual := ToKebabCase(c.src)
		if actual != c.expect {
			t.Fatalf("%5d. expert %q -> %q but got %q", i, c.src, c.expect, actual)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	for i, c := range []struct {
		src    string
		expect string
	}{
		{"IGotInternAtGeeksForGeeks", "i_got_intern_at_geeks_for_geeks"},
		{"iGotInternAtGeeksForGeeks", "i_got_intern_at_geeks_for_geeks"},
		{"i-got-intern-at-geeks-for-geeks", "i_got_intern_at_geeks_for_geeks"},
		{"i_got_intern_at_geeks_for_geeks", "i_got_intern_at_geeks_for_geeks"},
	} {
		actual := ToSnakeCase(c.src)
		if actual != c.expect {
			t.Fatalf("%5d. expert %q -> %q but got %q", i, c.src, c.expect, actual)
		}
	}
}
