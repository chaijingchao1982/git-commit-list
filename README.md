# git-commit-list

## 简介

当我们使用主干开发方式，面对生产环境分支较多的情况，就需要一份“git 提交记录文件”来帮助我们记录哪些提交被合并到了哪些分支。

本工具的功能就是：将指定 git 仓库的主干 main 分支的 commit 信息，持续的追加到“git 提交记录文件”的最后，以便用其来记录管理各分支合并情况。

## 用法

- 使用 go install 安装本工具

    ```shell
    go install github.com/chaijingchao1982/git-commit-list/cmd/git-commit-list@latest
    ```

- 创建一个工作目录，新建 output.xlsx 文件，并填写起始的 commit id，[参考例子](./output.xlsx)。

  之后的每次执行命令，都会将最新的 commit 信息追加到 output.xlsx 文件的最后。

- 运行工具

    ```shell
    git-commit-list {your-git-repo-path}
    ```
