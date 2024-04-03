default:
	ls -l
# Genrate
docs:
	hype export -format=markdown -f README.md.tpl > README.md
