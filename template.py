"python3"
import ctypes,os
f = os.fdopen(ctypes.CDLL("libc.so.6").memfd_create("{{.Name}}",0),mode="wb",buffering=0)
{{range .File}}f.write(b"{{.}}")
{{end}}f.flush()
os.execl("/proc/%d/fd/%d"%(os.getpid(),f.fileno()){{if eq 0 (len .Args)}},"/proc/%d/fd/%d"%(os.getpid(),f.fileno()){{else}}{{range .Args}},b"{{.}}"{{end}}{{end}})

