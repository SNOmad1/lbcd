package normalization

import (
	"github.com/lbryio/lbcd/claimtrie/param"
	"golang.org/x/text/unicode/norm"
)

var Normalize = normalizeGo
var NormalizeTitle = "Normalizing strings via Go. Casefold table version = 11.0.0, NFD version = " + norm.Version

func NormalizeIfNecessary(name []byte, height int32) []byte {
	if height < param.ActiveParams.NormalizedNameForkHeight {
		return name
	}
	return Normalize(name)
}

func normalizeGo(value []byte) []byte {

	normalized := norm.NFD.Bytes(value) // may need to hard-code the version on this
	// not using x/text/cases because it does too good of a job; it seems to use v14 tables even when it claims v13
	return CaseFold(normalized)
}
