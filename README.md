# dotfiler

[![Static Checks - Go](https://github.com/jakewan/dotfiler/actions/workflows/static-checks-go.yml/badge.svg)](https://github.com/jakewan/dotfiler/actions/workflows/static-checks-go.yml)

A tool for managing dotfiles.

Create a private Git repository containing source dotfiles and a `dotfiler.yml` manifest file.

Example:

```shell
git/
├─ .gitconfig
├─ .gitignore_global
gnupg/
├─ mac/
│  ├─ gpg-agent.conf
│  ├─ gpg-agent-arm.conf
neovim/
├─ init.vim
vscode/
├─ mac/
│  ├─ settings.json
│  ├─ keybindings.json
zsh/
├─ .zshrc
dotfiler.yml
```

Example dotfiler.yml:

```yaml
---
- op: symlink
  srcFilePath: zsh/.zshrc
  dstFilePath: .zshrc
  targetOS:
    - darwin
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: git/.gitconfig
  dstFilePath: .gitconfig
  targetOS:
    - darwin
    - linux
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: git/.gitignore_global
  dstFilePath: .gitignore_global
  targetOS:
    - darwin
    - linux
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: neovim/init.vim
  dstFilePath: .config/nvim/init.vim
  targetOS:
    - darwin
    - linux
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: vscode/mac/settings.json
  dstFilePath: Library/Application Support/Code/User/settings.json
  targetOS:
    - darwin
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: vscode/mac/keybindings.json
  dstFilePath: Library/Application Support/Code/User/keybindings.json
  targetOS:
    - darwin
  targetArch:
    - amd64
    - arm
- op: symlink
  srcFilePath: gnupg/mac/gpg-agent.conf
  dstFilePath: .gnupg/gpg-agent.conf
  targetOS:
    - darwin
  targetArch:
    - amd64
- op: symlink
  srcFilePath: gnupg/mac/gpg-agent-arm.conf
  dstFilePath: .gnupg/gpg-agent.conf
  targetOS:
    - darwin
  targetArch:
    - arm
```

The values for `targetOS` and `targetArch` should match those indicated by Go's `GOOS` and `GOARCH` [runtime Constants](https://pkg.go.dev/runtime#pkg-constants).

You can view a list of possible combinations by running `go tool dist list`.

For example:

```shell
go tool dist list | egrep 'darwin|linux'
darwin/amd64
darwin/arm64
linux/386
linux/amd64
linux/arm
linux/arm64
...
```

Finally, run the command to update your system's dotfiles:

```shell
dotfiler files update -m <path to dotfiler.yml>
```
