# Tablefy

Interactive table formatter for bash command output. Converts text with table format into beautiful tables using lipgloss and bubbletea.

## Requirements

- Go 1.25.4 or higher

## Installation

```bash
go build -o bin/tablefy cmd/tablefy/main.go
```

## Usage

Tablefy reads from stdin and provides an interactive interface to explore your data:

```bash
# Example with docker ps 
docker ps | ./bin/tablefy

# Example with custom output
echo "NAME  AGE  CITY
John  25  Madrid
Maria  30  Barcelona
Pedro  28  Valencia" | ./bin/tablefy
```

## Features

### Interactive Navigation
- **← → / h l**: Navigate between columns
- **↑ ↓ / j k**: Scroll through rows
- **s**: Toggle selection of current column (can select multiple)
- **Enter / Space**: Zoom into selected columns (creates new table with only those columns)
- **q**: Exit zoom mode or quit the application
- **Esc / Ctrl+C**: Quit the application

### Automatic Formatting
- Reads input from stdin
- Detects columns based on whitespace patterns in the header
- The first line is used as the header
- Automatically adjusts column widths to fit terminal width
- Truncates data when necessary with "..." to fit the screen
- Formats the result with borders and colors using lipgloss
- Gives priority to the last column (typically COMMAND in ps output)

### Column Zoom
- Navigate with arrow keys or h/l to highlight a column
- Press **s** to toggle selection (selected columns are highlighted in purple)
- You can select one or multiple columns
- Press **Enter** or **Space** to zoom into the selected columns
- The zoomed view creates a new table with only the selected columns
- This new table applies all the same formatting rules (width calculation, truncation, etc.)
- Press **q** to exit zoom and return to the normal view

### Workflow example:
1. Run `ps aux | ./bin/tablefy`
2. Use arrow keys to navigate to the "USER" column
3. Press **s** to select it (it turns purple)
4. Navigate to "COMMAND" column
5. Press **s** to select it too
6. Press **Enter** to zoom - now you see only USER and COMMAND columns in a new table
7. Press **q** to return to the full table view

