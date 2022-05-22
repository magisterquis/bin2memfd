package bin2memfd

/*
 * perl_test.go
 * Perl-specific code tests
 * By J. Stuart McMurray
 * Created 20220522
 * Last Modified 20220522
 */

import "testing"

func TestEncoderPerl(t *testing.T) {
	for _, c := range []struct {
		Args []string
		Have string
		Want string
	}{{
		Have: "abc\n123",
		Want: "\"perl\";\n$a=`uname -m`;chomp$a;$s={\"x86_64\"=>319}->{$a} or die \"unknown arch $a\";\n$n=syscall(319,my$a=\"\",0) or die\"s:$!\";\nopen(F,\">&=$n\") or die\"o:$!\";\nselect F;$|=1;\nprint \"abc\\x0a123\";\nexec \"/proc/$$/fd/$n\" or die\"e:$!\";\n",
	}, {
		Args: []string{"foo", "bar", "trid\nge$x\\n"},
		Have: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Want: "\"perl\";\n$a=`uname -m`;chomp$a;$s={\"x86_64\"=>319}->{$a} or die \"unknown arch $a\";\n$n=syscall(319,my$a=\"\",0) or die\"s:$!\";\nopen(F,\">&=$n\") or die\"o:$!\";\nselect F;$|=1;\nprint \"abcdefghijklmnopqrstuvwxyzABCDEF\";\nprint \"GHIJKLMNOPQRSTUVWXYZ\";\nmy@a=(\"foo\",\"bar\",\"trid\\x0age\\x24x\\x5cn\",);\nexec{\"/proc/$$/fd/$n\"}@a or die\"e:$!\";\n",
	}} {
		got, err := Encoder{Args: c.Args}.Perl([]byte(c.Have))
		if nil != err {
			t.Errorf(
				"Perl\n\targs:%q\n\thave:%q\n\terr:%s",
				c.Args,
				c.Have,
				err,
			)
			continue
		}
		if string(got) == c.Want {
			continue
		}
		t.Errorf(
			"Perl\n\targs:%q\n\thave:%q\n\t got:%q\n\twant:%q",
			c.Args,
			c.Have,
			got,
			c.Want,
		)
	}
}
