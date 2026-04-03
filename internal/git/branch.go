package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
)

// CurrentBranchName 返回仓库当前 HEAD 所在的分支名称。
// 如果 HEAD 处于 detached 状态，返回错误。
func CurrentBranchName(repo *gogit.Repository) (string, error) {
	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("获取 HEAD 失败：%w", err)
	}

	if !headRef.Name().IsBranch() {
		return "", fmt.Errorf("HEAD 处于 detached 状态，不在任何分支上")
	}

	return headRef.Name().Short(), nil
}
