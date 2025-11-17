# Tablefy

Interactive table formatter for bash command output. Converts text with table format into beautiful tables using lipgloss and bubbletea.

## Installation

### Option 1: Download from GitHub Releases (Recommended)

Download the latest release:

```bash
curl -L https://github.com/janorga/tablefy/releases/latest/download/tablefy-linux-amd64 -o tablefy
chmod +x tablefy
sudo mv tablefy /usr/local/bin/
```

### Option 2: Build from Source

Requirements:
- Go 1.25.4 or higher

```bash
git clone https://github.com/janorga/tablefy.git
cd tablefy

# Using make (recommended)
make build                 # Build with version from git tags
make build-dev            # Quick build for development

# Or using go directly
go build -o bin/tablefy cmd/tablefy/main.go

# Optional: Install globally
sudo cp bin/tablefy /usr/local/bin/
```

#### Make targets

- `make build`: Builds the application with version information from git tags
- `make build-dev`: Quick build for development (version: dev)
- `make test`: Run all tests
- `make version`: Show current version from git
- `make clean`: Remove build artifacts

## Usage

Tablefy reads from stdin and provides an interactive interface to explore your data:

```bash
# Example with ps
ps aux | tablefy

# Example with docker ps
docker ps | tablefy

# Example with df
df -h | tablefy

# Example with custom output
echo "NAME  AGE  CITY
John  25  London
Alice 30  Paris
Bob   28  Berlin" | tablefy
```

## Command Line Flags

### Version
```bash
tablefy --version    # Long form
tablefy -v           # Short form
```

Displays the version of tablefy and exits.

### Auto-expand mode
```bash
tablefy --auto-expand    # Long form
tablefy -a               # Short form
```

Enables auto-expand mode for focused columns with truncated content. When navigating to a column that contains truncated cells (marked with "..."), that column automatically expands to show the full content. Other columns shrink proportionally to maintain the terminal width. The column shrinks back to normal when you navigate away.

**Example:**
```bash
helm list -A | tablefy --auto-expand
```

Navigate with arrow keys to any column with truncated content and watch it expand to reveal the full data.

## Features

### Interactive Navigation
- **← → / h l**: Navigate between columns
- **↑ ↓ / j k**: Scroll through rows
- **s**: Toggle selection of current column (can select multiple)
- **Enter / Space**: Zoom into selected columns (creates new table with only those columns)
- **f**: Fuzzy filter rows by current column values
- **c**: Clear active filter and show all rows
- **q**: Exit zoom mode or quit the application
- **Esc / Ctrl+C**: Quit the application

### Fuzzy Filter

Press **f** to enter filter mode for the currently focused column. The fuzzy filter allows you to quickly narrow down rows by searching for values in that column using intelligent subsequence matching.

**How it works:**
- Navigate to a column you want to filter by
- Press **f** to activate the fuzzy finder
- Type a search query (e.g., "run" will match "running", "runner", "runtime")
- See live previews of matching rows with a count (e.g., "Filter [STATUS]: run (3 matches)")
- Press **Enter** to apply the filter and work with filtered data
- Press **Esc** to cancel without applying

**After applying filter:**
- You can **navigate** between columns normally (← →, h/l)
- You can **scroll** through filtered rows (↑ ↓, j/k)
- You can **select columns** and **zoom** into them (s, Enter)
- **Auto-expand** feature works with filtered data
- Press **c** to clear the filter and return to all rows

**How fuzzy matching works:**
The filter uses subsequence matching where query characters must appear in order (case-insensitive):
- Query "run" matches: running, runner, runtime (all have r, u, n in order)
- Query "rung" matches: running (has r, u, n, g in order)
- Query "runn" matches: running, runner (both have two n's)
- More specific queries return fewer results - exactly what you'd expect!

**Example workflow:**
```bash
ps aux | tablefy
# Navigate to the STAT column (use ← →)
# Press 'f' to filter
# Type "S" to show only sleeping processes (S, Ss, S+, etc.)
# Press Enter to apply filter
# Now navigate, scroll, and zoom through only the matching rows
# Press 'c' to clear filter and go back to all processes
```

**Use cases:**
- Filter docker containers by STATUS to see only running/stopped containers
- Filter Kubernetes pods to show only those in Error or Pending state
- Filter process list to show only Java processes, Bash shells, or specific users
- Quickly reduce large datasets to find what you're looking for

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
1. Run `ps aux | tablefy`
2. Use arrow keys to navigate to the "USER" column
3. Press **s** to select it (it turns purple)
4. Navigate to "COMMAND" column
5. Press **s** to select it too
6. Press **Enter** to zoom - now you see only USER and COMMAND columns in a new table
7. Press **q** to return to the full table view

## Example Output

Input:
```
NAME    AGE    CITY
John    25     London
Alice   30     Paris
```

Output (with colors and borders):
```
┌───────┬─────┬────────┐
│ NAME  │ AGE │ CITY   │
├───────┼─────┼────────┤
│ John  │ 25  │ London │
│ Alice │ 30  │ Paris  │
└───────┴─────┴────────┘

← → / h l: Navigate | s: Toggle select (0 selected) | Enter: Zoom | q: Quit
```

