"perl";
$a=`uname -m`;chomp$a;$s={"x86_64"=>319}->{$a} or die "unknown arch $a";
$n=syscall(319,my$a="{{.Name}}",0) or die"s:$!";
open(F,">&=$n") or die"o:$!";
select F;$|=1;
{{range .File}}print "{{.}}";
{{end}}{{if eq 0 (len .Args)}}exec "/proc/$$/fd/$n" or die"e:$!";{{else}}my@a=({{range .Args}}"{{.}}",{{end}});
exec{"/proc/$$/fd/$n"}@a or die"e:$!";{{end}}
