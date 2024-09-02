# dotfiler

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

Run the command to update dotfiles:

```shell
dotfiler files update -m <path to dotfiler.yml>
```
