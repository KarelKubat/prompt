# prompt

<!-- toc -->
<!-- /toc -->

Just my favorite console prompt. If you like it, use it for your benefit. The installation is very simple:

```shell
# Install as a binary
go install prompt.go

# Then edit your ~/.bashrc and set (as an example):
PS1=""
PROMPT_COMMAND="prompt -host $(hostname)"
```

To see settings, try:

```shell
prompt --help
```

This shows how to include in the prompt:

- The hostname, 
- Date, 
- Current working directory and git status (if inside a git repository),
- Etc.

`prompt` shows its output as colored text, suitable for regular terminals. There is currently no possibility to change the colors.
