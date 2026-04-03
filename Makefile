$(VERBOSE).SILENT:

# COLORS
RED		:= $(shell tput -Txterm setaf 1)
CYAN	:= $(shell tput -Txterm setaf 6)
RESET	:= $(shell tput -Txterm sgr0)

update: 
	echo "$(CYAN)show git commit list$(RESET)"
	go run cmd/git-commit-list/main.go "{your-git-repo-path}"
