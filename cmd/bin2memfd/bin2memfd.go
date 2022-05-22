// Program bin2memfd turns a program into a script which will run it in a
// memfd.
package main

/*
 * bin2memfd.go
 * Run a program in a memfd
 * By J. stuart McMurray
 * Created 20220521
 * Last Modified 20220522
 */

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/magisterquis/bin2memfd"
)

func main() {
	var (
		progF = flag.String(
			"prog",
			"",
			"Program file to encode, or stdin if not specified",
		)
		name = flag.String(
			"name",
			"",
			"Memfd symlink target `name`",
		)
		lang = flag.String(
			"language",
			"perl",
			"Generated script `language`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options] [arg [arg...]]

Encodes the program to a script which, which executed, places the program in
a memfd.  If no program is specified, it is read from stdin.  If arguments
are specified, they will be used as the program's argv, starting with argv[0].

This can be used something like

%s ./implant nmap -A -v -n -p- 10.0.0.0/16 | ssh target perl

The currently-supported -language options are perl and python.

Options:
`,
			os.Args[0],
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Add the name and args, if we have them. */
	var enc bin2memfd.Encoder
	if 0 < flag.NArg() {
		enc.Args = flag.Args()
	}
	if "" != *name {
		enc.Name = *name
	}

	/* Work out which language to use. */
	var encf func([]byte) ([]byte, error)
	switch *lang {
	case "python":
		encf = enc.Python
	case "perl":
		fallthrough
	default:
		encf = enc.Perl
	}

	/* Open the script or program. */
	var f *os.File
	if "" != *progF {
		var err error
		f, err = os.Open(*progF)
		if nil != err {
			log.Fatalf("Error opening %q: %s", *progF, err)
		}
	} else {
		f = os.Stdin
	}
	prog, err := io.ReadAll(f)
	if nil != err {
		log.Fatalf("Reading in program: %s", err)
	}

	/* Encode the program and spit it out. */
	s, err := encf(prog)
	if nil != err {
		log.Fatalf("Encoding program: %s", err)
	}
	if _, err := os.Stdout.Write(s); nil != err {
		log.Fatalf("Writing program: %s", err)
	}
}
