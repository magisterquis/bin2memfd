Bin2MemFD
=========
Encodes a program (which can be a script, despite the name) to a Perl or Python
script which sticks it in a Linux memfd and runs it.  The goal is to enable
staged implants to be run with `curl | perl`, or something similar.

The perl version only works on 64-bit Intel linux (i.e. `uname -m` returns
`x86_64`, but updating [`template.pl`](./template.pl) to accomodate other
architectures is very easy.  Pull requests welcome.

Example
-------
```go
b := slurpImplant()
e, err := bin2memfd.Encoder{
        Args: []string{"cron", "hourly-cleanup"},
        Name: "cron-tmp48240",
}.Perl(b)
if nil != err {
        log.Fatalf("Perlizing: %s", err)
}
serveToStager(e)
```

Standalone Program
------------------
A standalone program to encode programs is included in
[cmd/bin2memfd](./cmd/bin2memfd).
