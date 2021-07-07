package normalization

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizationGo(t *testing.T) {
	testNormalization(t, normalizeGo)
}

func testNormalization(t *testing.T, normalize func(value []byte) []byte) {

	r := require.New(t)

	r.Equal("test", string(normalize([]byte("TESt"))))
	r.Equal("test 23", string(normalize([]byte("tesT 23"))))
	r.Equal("\xFF", string(normalize([]byte("\xFF"))))
	r.Equal("\xC3\x28", string(normalize([]byte("\xC3\x28"))))
	r.Equal("\xCF\x89", string(normalize([]byte("\xE2\x84\xA6"))))
	r.Equal("\xD1\x84", string(normalize([]byte("\xD0\xA4"))))
	r.Equal("\xD5\xA2", string(normalize([]byte("\xD4\xB2"))))
	r.Equal("\xE3\x81\xB5\xE3\x82\x99", string(normalize([]byte("\xE3\x81\xB6"))))
	r.Equal("\xE1\x84\x81\xE1\x85\xAA\xE1\x86\xB0", string(normalize([]byte("\xEA\xBD\x91"))))
}

func randSeq(n int) []byte {
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyz̃ABCDEFGHIJKLMNOPQRSTUVWXYZ̃")

	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return []byte(string(b))
}

func BenchmarkNormalize(b *testing.B) {
	benchmarkNormalize(b, normalizeGo)
}

func benchmarkNormalize(b *testing.B, normalize func(value []byte) []byte) {
	rand.Seed(42)
	strings := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		strings[i] = randSeq(32)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := normalize(strings[i])
		require.True(b, len(s) >= 8)
	}
}
