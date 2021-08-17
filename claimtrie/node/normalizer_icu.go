// +build use_icu_normalization

package node

// #cgo CFLAGS: -O2
// #cgo LDFLAGS: -licuio -licui18n -licuuc -licudata
// #include <unicode/unorm2.h>
// #include <unicode/ustring.h>
// #include <unicode/uversion.h>
// int icu_version() {
//    UVersionInfo info;
//    u_getVersion(info);
//    return ((int)(info[0]) << 16) + info[1];
// }
// int normalize(char* name, int length, char* result) {
//   UErrorCode ec = U_ZERO_ERROR;
//   static const UNormalizer2* normalizer = NULL;
//   if (normalizer == NULL) normalizer = unorm2_getNFDInstance(&ec);
//   UChar dest[256]; // maximum claim name size is 255; we won't have more UTF16 chars than bytes
//   int dest_len;
//   u_strFromUTF8(dest, 256, &dest_len, name, length, &ec);
//   if (U_FAILURE(ec) || dest_len == 0) return 0;
//   UChar normalized[256];
//   dest_len = unorm2_normalize(normalizer, dest, dest_len, normalized, 256, &ec);
//   if (U_FAILURE(ec) || dest_len == 0) return 0;
//   dest_len = u_strFoldCase(dest, 256, normalized, dest_len, U_FOLD_CASE_DEFAULT, &ec);
//   if (U_FAILURE(ec) || dest_len == 0) return 0;
//   u_strToUTF8(result, 512, &dest_len, dest, dest_len, &ec);
//   return dest_len;
// }
import "C"
import (
	"fmt"
	"unsafe"
)

func init() {
	Normalize = normalizeICU
}

func IcuVersion() string {
	// TODO: we probably need to explode if it's not 63.2 as it affects consensus
	result := C.icu_version()
	return fmt.Sprintf("%d.%d", result>>16, result&0xffff)
}

func normalizeICU(value []byte) []byte {
	if len(value) <= 0 {
		return value
	}
	name := (*C.char)(unsafe.Pointer(&value[0]))
	length := C.int(len(value))

	// hopefully this is a stack alloc (but it may be a bit large for that):
	var resultName [512]byte // inputs are restricted to 255 chars; it shouldn't expand too much past that
	result := unsafe.Pointer(&resultName[0])

	resultLength := C.normalize(name, length, (*C.char)(result))
	if resultLength == 0 {
		return value
	}

	// return resultName[0:resultLength] -- we want to shrink the result (not use a slice on 1024)
	return C.GoBytes(result, resultLength)
}
