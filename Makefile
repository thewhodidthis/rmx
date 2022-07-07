go/run:
	@go run main.go
go:
	env GOOS=openbsd GOARCH=amd64 go build
post:
	rmx | fmt -c | fmt -n -w 80 > $$(date +"%Y%m%d".txt)
list:
	@ls -A1r *.{txt,xml} | listing -x gophermap.tmpl > gophermap
feed: x
	@ls -A1r *.txt | head -15 | listing -x feed.xml.tmpl -t "R*M*X" > feed.xml
a:
	@echo $$(find *.txt -exec sh -c 'sed -n 1p {} | xargs | sed "s/~ //" | sed "s/ ~//" | sed "s/\(.*\)/\"\1\"\ /g"' \; | tr -ds "\n" " " | sed "s/\(.*\)/{{- $$\a := list \1\-}}/") > a
b:
	@echo $$(find *.txt -exec sh -c 'sed -n 5,6p {} | tr -ds "\n" " " | sed "s/\(.*\)/\"\1...\"\ /g"' \; | sed "s/\(.*\)/{{- $$\b := list \1\-}}/") > b
x: a b
	cat a b feed.tmpl > c; mv c feed.xml.tmpl; rm -rf a b c
