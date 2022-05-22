// Package bin2memfd converts a Linux binary into a program which loads the
// binary into a memfd and runs it.
package bin2memfd

/*
 * bin2memfd.go
 * Convert a program to a script which runs the program from a memfd.
 * By J. Stuart McMurray
 * Created 20220521
 * Last Modified 20220522
 */

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

/* Program templates. */
var (
	//go:embed template.pl template.py
	tfs  embed.FS
	tmpl *template.Template
)

func init() {
	var err error
	tmpl, err = template.ParseFS(tfs, "*")
	if nil != err {
		panic(fmt.Sprintf("failed to init templates: %s", err))
	}
}

/* tArgs holds the arguments we pass to a template. */
type tArgs struct {
	Args []string /* Argv */
	File []string /* Contents of file to encode. */
	Name string   /* memfd name. */
}

// DefaultEncoder is the Encoder used by package-level functions.
var DefaultEncoder = Encoder{}

// Perl encodes the program read from r to a Perl script which loads it into
// a memfd and runs it.  It uses DefaultEncoder for settings.
func Perl(b []byte) ([]byte, error) {
	return DefaultEncoder.Perl(b)
}

// Python encodes the program read from r to a Python script which loads it
// into a memfd and runs it.  It uses DefaultEncoder for settings.
func Python(b []byte) ([]byte, error) {
	return DefaultEncoder.Python(b)
}

// Encoder is used to tune how a binary is encoded.  Encoder's fields must not
// be modified while any of encoder's methods are being called, though multiple
// methods may be called at once.
type Encoder struct {
	// Args is the executed program's argv.  If it is empty, the program's
	// argv[0] will be the path to the memfd and no other arguments will
	// be passed.
	Args []string

	// Name is the name passed to memfd_create.  This is the name visible
	// in the memfd's link target.
	Name string
}

/* execute executes the template with the given name. */
func (e Encoder) execute(name string, file []byte) ([]byte, error) {
	/* Roll sanitized arguments. */
	args := tArgs{
		Args: make([]string, len(e.Args)),
		File: encodeToStrings(file),
		Name: sanitize(e.Name),
	}
	for i, v := range e.Args {
		args.Args[i] = sanitize(v)
	}

	/* Encode with a template. */
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, name, args)
	return buf.Bytes(), err
}

// Perl encodes the program read from r to a Perl script which loads it into
// a memfd and runs it.  Currently this is only supported on amd64; adding
// other architectures is a very simple addition to template.pl.
func (e Encoder) Perl(b []byte) ([]byte, error) {
	return e.execute("template.pl", b)
}

// Python encodes the program read from r to a Python script which loads it
// into a memfd and runs it.
func (e Encoder) Python(b []byte) ([]byte, error) {
	return e.execute("template.py", b)
}
