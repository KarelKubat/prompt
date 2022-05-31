# prompt

Just my favorite console prompt. If you like it, use it for your benefit. The installation is very simple:

```shell
# Install as a binary
go build prompt.go

# Then edit your ~/.bashrc and set:
  PS1=""
  PROMPT_COMMAND="prompt -host $(hostname)"
```