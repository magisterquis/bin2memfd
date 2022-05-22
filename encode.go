package bin2memfd

/*
 * sanitize.go
 * Perl-specific code
 * By J. Stuart McMurray
 * Created 20220522
 * Last Modified 20220522
 */

import (
	"fmt"
	"strings"
)

/* stringChunkLen is the number of bytes to put in a string returned by
encodeToStrings. */
const stringChunkLen = 32

/* sanitize sanitizes a string to something all languages can handle. */
func sanitize(s string) string {
	var sb strings.Builder
	for _, v := range []byte(s) {
		if strings.ContainsRune(`@$\"`, rune(v)) ||
			v < ' ' || '~' < v {
			fmt.Fprintf(&sb, "\\x%02x", v)
			continue
		}
		sb.WriteByte(v)
	}
	return sb.String()
}

/* encodeToStrings encodes the bytes in b to chunks encoded in escaped string
format. */
func encodeToStrings(b []byte) []string {
	ss := make([]string, 0, len(b)/stringChunkLen+1) /* Close enough. */
	/* Grab each chunk and encode it. */
	for i := 0; i < len(b); i += stringChunkLen {
		/* Get hold of a chunk of b. */
		end := i + 32
		if len(b) < end {
			end = len(b)
		}
		ss = append(ss, sanitize(string(b[i:end])))
	}
	return ss
}
