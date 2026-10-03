package update

import (
	"strconv"
	"strings"
)

// Semver representa una versión parseada major.minor.patch. No se modelan
// pre-releases: las releases de goreleaser son siempre vX.Y.Z.
type Semver struct {
	Major int
	Minor int
	Patch int
}

// isDev indica si la versión corresponde a una compilación local sin release.
func isDev(v string) bool {
	return strings.TrimSpace(v) == "dev"
}

// parseSemver parsea "v?major.minor.patch" (la "v" es opcional). Devuelve
// ok=false para "dev", cadenas vacías o cualquier formato no reconocido.
func parseSemver(v string) (Semver, bool) {
	s := strings.TrimSpace(v)
	s = strings.TrimPrefix(s, "v")
	if s == "" {
		return Semver{}, false
	}
	// Ignora cualquier sufijo de pre-release/metadata tras el patch.
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	fields := strings.Split(s, ".")
	if len(fields) != 3 {
		return Semver{}, false
	}
	var nums [3]int
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil || n < 0 {
			return Semver{}, false
		}
		nums[i] = n
	}
	return Semver{Major: nums[0], Minor: nums[1], Patch: nums[2]}, true
}

// compareSemver compara a con b y devuelve -1, 0 o 1.
func compareSemver(a, b Semver) int {
	if c := cmpInt(a.Major, b.Major); c != 0 {
		return c
	}
	if c := cmpInt(a.Minor, b.Minor); c != 0 {
		return c
	}
	return cmpInt(a.Patch, b.Patch)
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
