# git-commit-list

## 简介

当我们使用主干开发方式，面对生产环境分支较多的情况，就需要一份“git 提交记录文件”来帮助我们记录哪些提交被合并到了哪些分支。

本工具的功能就是：将指定 git 仓库的主干 main 分支的 commit 信息，持续的追加到“git 提交记录文件”的最后，以便用其来记录管理各分支合并情况。

## 用法

- 首先要在 output.xlsx 中填写一个起始的 commit id，例如：xxx 功能  e56a1234，之后每次运行本工具，都会将最新的 commit 信息追加到文件最后
- 执行命令之前确保 git 仓库已经切换到了 main 分支，且 git pull 拉取了最新代码
- 命令：go run main.go [git-repo-path]
- 推荐使用 Makefile 简化使用

    在 Makefile 中添加：

    ```shell
    go run main.go "/home/chaijingchao/xxx"
    ````

    使用 Makefile 命令：

    ```shell
    make update
    ````

## 输出格式

```shell
提交日志1  commit-id1
提交日志2  commit-id2
提交日志3  commit-id3
````
