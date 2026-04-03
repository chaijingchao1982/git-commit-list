package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/chaijingchao1982/git-commit-list/internal/excel"
	internalgit "github.com/chaijingchao1982/git-commit-list/internal/git"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/mattn/go-runewidth"
	"github.com/xuri/excelize/v2"
)

const (
	outputFile     = "output.xlsx"
	sheetName      = "Sheet1"
	mainBranchName = "main"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <repo-path>", os.Args[0])
	}

	repoPath := os.Args[1]
	os.Exit(run(repoPath))
}

func run(repoPath string) int {
	lastCommitHash, err := excel.ReadLastCommitHash(outputFile, sheetName)
	if err != nil {
		log.Printf("[ERROR] failed to read last commit: %v", err)
		return 1
	}
	hashPrefixLen := len(lastCommitHash)
	log.Printf("[INFO] last commit: %s (len %d)", lastCommitHash, hashPrefixLen)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Printf("[ERROR] Unable to open the Git repository: %s", repoPath)
		return 1
	}

	branchName, err := internalgit.CurrentBranchName(repo)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return 1
	}
	log.Printf("[INFO] current branch: %s", branchName)
	if branchName != mainBranchName {
		log.Printf("[ERROR] current branch is '%s', please switch to '%s' branch before running", branchName, mainBranchName)
		return 1
	}

	startHash, err := repo.ResolveRevision(plumbing.Revision(lastCommitHash))
	if err != nil {
		log.Printf("[ERROR] invalid commit ID '%s'", lastCommitHash)
		return 1
	}

	head, err := repo.Head()
	if err != nil {
		log.Printf("[ERROR] unable to get HEAD")
		return 1
	}
	headHash := head.Hash()
	log.Printf("[INFO] HEAD: %s", headHash.String()[:hashPrefixLen])

	commits, err := collectCommits(repo, headHash, *startHash)
	if err != nil {
		log.Printf("[ERROR] failed to collect commits: %v", err)
		return 1
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Author.When.Before(commits[j].Author.When)
	})

	if err := writeCommitsToExcel(outputFile, sheetName, commits, hashPrefixLen); err != nil {
		log.Printf("[ERROR] failed to write commits to Excel: %v", err)
		return 1
	}

	return 0
}

// collectCommits collects all commits from HEAD to stopHash (excluding stopHash).
func collectCommits(repo *git.Repository, headHash, stopHash plumbing.Hash) ([]*object.Commit, error) {
	iter, err := repo.Log(&git.LogOptions{From: headHash})
	if err != nil {
		return nil, fmt.Errorf("failed to read log: %w", err)
	}
	defer iter.Close()

	var commits []*object.Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == stopHash {
			return storer.ErrStop
		}
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commit history: %w", err)
	}

	return commits, nil
}

// writeCommitsToExcel appends commit records to an Excel file and prints aligned summaries to the terminal.
func writeCommitsToExcel(filename, sheet string, commits []*object.Commit, hashPrefixLen int) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer f.Close()

	maxWidth := 0
	for _, c := range commits {
		subject := strings.SplitN(c.Message, "\n", 2)[0]
		if w := runewidth.StringWidth(subject); w > maxWidth {
			maxWidth = w
		}
	}
	alignCol := maxWidth + 4

	for _, c := range commits {
		subject := strings.SplitN(c.Message, "\n", 2)[0]
		shortHash := c.Hash.String()[:hashPrefixLen]

		if err := excel.AppendRow(f, sheet, subject, shortHash); err != nil {
			return fmt.Errorf("failed to append row to Excel: %w", err)
		}

		printAligned(subject, shortHash, alignCol)
	}

	return nil
}

// printAligned prints aligned subject and commit hash to the terminal.
func printAligned(subject, hash string, alignCol int) {
	w := runewidth.StringWidth(subject)
	if w < alignCol {
		padding := strings.Repeat(" ", alignCol-w)
		fmt.Printf("%s%s%s\n", subject, padding, hash)
	} else {
		fmt.Printf("%s %s\n", subject, hash)
	}
}
