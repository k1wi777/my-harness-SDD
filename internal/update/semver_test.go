package update

import "testing"

func TestParseSemver(t *testing.T) {
	cases := []struct {
		in   string
		want Semver
		ok   bool
	}{
		{"v0.1.0", Semver{0, 1, 0}, true},
		{"0.1.0", Semver{0, 1, 0}, true},
		{"1.2.3", Semver{1, 2, 3}, true},
		{"v1.2.3-rc1", Semver{1, 2, 3}, true},
		{"  v2.0.0 ", Semver{2, 0, 0}, true},
		{"dev", Semver{}, false},
		{"", Semver{}, false},
		{"1.2", Semver{}, false},
		{"v1.2.x", Semver{}, false},
		{"v1.-2.3", Semver{}, false},
	}
	for _, c := range cases {
		got, ok := parseSemver(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("parseSemver(%q) = (%v, %v), want (%v, %v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestCompareSemver(t *testing.T) {
	cases := []struct {
		a, b Semver
		want int
	}{
		{Semver{0, 1, 0}, Semver{0, 1, 0}, 0},
		{Semver{0, 2, 0}, Semver{0, 1, 9}, 1},
		{Semver{0, 1, 0}, Semver{1, 0, 0}, -1},
		{Semver{1, 0, 1}, Semver{1, 0, 0}, 1},
	}
	for _, c := range cases {
		if got := compareSemver(c.a, c.b); got != c.want {
			t.Errorf("compareSemver(%v, %v) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestIsDev(t *testing.T) {
	if !isDev("dev") || !isDev(" dev ") {
		t.Fatal("isDev debe reconocer 'dev'")
	}
	if isDev("v0.1.0") || isDev("") {
		t.Fatal("isDev no debe aceptar versiones publicadas")
	}
}
