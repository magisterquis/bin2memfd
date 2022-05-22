package bin2memfd

/*
 * perl_test.go
 * Perl-specific code tests
 * By J. Stuart McMurray
 * Created 20220521
 * Last Modified 20220522
 */

import (
	"reflect"
	"testing"
)

func TestEncodeToStrings(t *testing.T) {
	for _, c := range []struct {
		Have string
		Want []string
	}{{
		Have: "foo",
		Want: []string{
			"foo",
		},
	}, {
		Have: string([]byte{
			0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01,
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x03, 0x00,
			0x3e, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x80, 0xcf, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x18, 0xfd,
			0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
		}),
		Want: []string{
			"\\x7fELF\\x02\\x01\\x01\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x03\\x00>\\x00\\x01\\x00\\x00\\x00\\x80\\xcf\\x00\\x00\\x00\\x00\\x00\\x00",
			"\\x40\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x18\\xfd\\x04\\x00\\x00\\x00\\x00\\x00",
		},
	}} {
		got := encodeToStrings([]byte(c.Have))
		if reflect.DeepEqual(got, c.Want) {
			continue
		}
		t.Errorf(
			"encodeToStrings\n\thave:%q\n\t got:%#v\n\twant:%#v",
			c.Have,
			got,
			c.Want,
		)
	}
}

func TestSanitize(t *testing.T) {
	for _, c := range []struct {
		Have string
		Want string
	}{{
		Have: "abc123",
		Want: "abc123",
	}, {
		Have: "abc123@$%\n",
		Want: "abc123\\x40\\x24%\\x0a",
	}, {
		Have: `foo"bar`,
		Want: "foo\\x22bar",
	}} {
		got := sanitize(c.Have)
		if c.Want == got {
			continue
		}
		t.Errorf(
			"sanitize have:%q got:%q want:%q",
			c.Have,
			got,
			c.Want,
		)
	}
}
