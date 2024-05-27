default:
	ls -l
# Genrate
docs:
	hype export -format=markdown -f README.md.tpl > README.md
install:
	go install cmd/cleura/cleura.go
