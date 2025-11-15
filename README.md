# Tablefy

Table formatter for bash command output. Converts text with table format into beautiful tables using lipgloss.

## Installation

```bash
go build -o tablefy
```

## Usage

Tablefy reads from stdin and automatically detects columns based on the header structure:

```bash
# Example with ls -l
ls -l | tablefy

# Example with ps
ps aux | tablefy

# Example with df
df -h | tablefy

# Example with custom output
echo "NAME  AGE  CITY
John  25  Madrid
Maria  30  Barcelona
Pedro  28  Valencia" | tablefy
```

## How it works

- Reads input from stdin
- Detects columns based on whitespace patterns in the header
- The first line is used as the header
- Automatically adjusts column widths to fit terminal width
- Truncates data when necessary with "..." to fit the screen
- Formats the result with borders and colors using lipgloss
- Gives priority to the last column (typically COMMAND in ps output)

## Example output

Input:
```
NAME    AGE    CITY
John    25     Madrid
Maria   30     Barcelona
```

Output (with colors and borders):
```
┌───────┬─────┬───────────┐
│ NAME  │ AGE │ CITY      │
├───────┼─────┼───────────┤
│ John  │ 25  │ Madrid    │
│ Maria │ 30  │ Barcelona │
└───────┴─────┴───────────┘
```
