# Tablefy

Interactive terminal table formatter that transforms CLI command output into beautifully formatted, fully interactive tables. 

Tablefy takes the raw output from commands like `ps aux`, `docker ps`, `kubectl get pods`, or any space-delimited text and turns it into a rich, navigable interface with powerful features:

- **Interactive Navigation**: Navigate columns with arrow keys, scroll through rows with PgUp/PgDn
- **Fuzzy Filtering**: Instantly filter rows by searching column values with intelligent subsequence matching
- **Column Zoom**: Select and focus on specific columns to reduce visual noise and see details
- **Auto-Expand**: Automatically expand truncated columns to reveal full content when you navigate to them
- **Beautiful Formatting**: Rich borders, colors, and optimal column width calculations using lipgloss and bubbletea

Perfect for exploring large datasets from shell commands, system administration, container orchestration, and DevOps workflows.

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
- **↑ ↓ / j k**: Scroll through rows (one row at a time)
- **PgUp / Page Up**: Scroll up by page
- **PgDn / Page Down**: Scroll down by page
- **s**: Toggle selection of current column (can select multiple)
- **Enter / Space**: Zoom into selected columns (creates new table with only those columns)
- **f**: Fuzzy filter rows by current column values
- **c**: Clear active filter and show all rows
- **o**: Export and quit (prints the visible table with aligned columns, no borders)
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
- **Auto-expand** feature works with filtered data - expansions are maintained
- Column widths are automatically recalculated based on filtered data for optimal display
- A filter indicator badge appears above the table showing: `🔍 Filter active: [COLUMN] = "QUERY" (X results)`
- Press **c** to clear the filter and return to all rows

**How fuzzy matching works:**
The filter uses subsequence matching where query characters must appear in order (case-insensitive):
- Query "run" matches: running, runner, runtime (all have r, u, n in order)
- Query "rung" matches: running (has r, u, n, g in order)
- Query "runn" matches: running, runner (both have two n's)
- More specific queries return fewer results - exactly what you'd expect!
- Characters can be non-consecutive but must maintain order: "jv" matches "javascript" but not "java"

**Column width optimization:**
When you filter data, column widths are recalculated based on the filtered subset rather than the entire dataset. This means:
- Narrower columns with less data will use less space
- Long text that couldn't fit before might now display without truncation
- The terminal space is used more efficiently with your filtered results
- Widths adjust dynamically as you refine your filter query

**Example workflow:**
```bash
ps aux | tablefy
# Navigate to the STAT column (use ← →)
# Press 'f' to filter
# Type "S" to show only sleeping processes (S, Ss, S+, etc.)
# Notice the filter indicator and see matches update in real-time
# Press Enter to apply filter
# Now navigate, scroll, and zoom through only the matching rows
# Column widths are now optimized for the filtered data
# Press 'c' to clear filter and go back to all processes
```

**Use cases:**
- Filter docker containers by STATUS to see only running/stopped containers
- Filter Kubernetes pods to show only those in Error or Pending state
- Filter process list to show only Java processes, Bash shells, or specific users
- Quickly reduce large datasets to find what you're looking for
- Combine with auto-expand to inspect detailed fields in filtered results

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

### Export Data

Press **o** to export the currently visible table and quit the application. The exported table will be printed to stdout with:
- Aligned columns with proper spacing
- No borders or styling (plain text format)
- No header row (only data rows)
- Support for all view modes (normal, filtered, zoomed)

**Use cases:**
- Export filtered results for further processing: `ps aux | tablefy | grep something`
- Save zoomed view output to a file: `docker ps | tablefy > containers.txt`
- Pipe filtered data to another command: `kubectl get pods | tablefy | xargs kubectl describe`
- Copy cleaned data for documentation or reporting

**Example workflow:**
```bash
# 1. Filter data interactively
ps aux | tablefy
# 2. Navigate and filter to show only interesting processes
# 3. Press 'o' to export and quit
# 4. The filtered, formatted table appears on stdout

# Or combine in a pipeline:
ps aux | tablefy > process_output.txt
# Then use the formatted table for further analysis
```

## Advanced Features

### Combining Auto-expand with Fuzzy Filter

When both **auto-expand** and **fuzzy filter** are used together, they work seamlessly:

**Workflow:**
```bash
ps aux | tablefy --auto-expand
# Navigate to COMMAND column (which may have truncated content)
# The column expands to show full command
# Press 'f' to filter by command
# Type a search term (e.g., "node" or "python")
# The expansion is MAINTAINED while you filter
# Filtered results are shown with optimized column widths
# You can still navigate and zoom on filtered data
# Column expansion persists until you navigate away
```

**Key behaviors:**
- When you press `f` to enter filter mode, **auto-expand state is preserved**
- While typing in the filter, column widths dynamically adjust based on the filtered subset
- If your filter significantly reduces the dataset, columns may shrink (no longer need the full expansion width)
- When you apply the filter and return to normal view, the expanded column is re-expanded based on the filtered data
- Pressing `c` to clear the filter re-expands columns for the full dataset

### Zoom with Filtered Data

The zoom feature works perfectly with filtered data:

```bash
ps aux | tablefy
# Filter by STATUS to show only sleeping processes
# Select USER and COMMAND columns with 's'
# Press Enter to zoom - zoomed view shows only filtered rows
# Zoom view also shows the filter indicator
# Press 'q' to return to filtered view (filter still active)
```


## Example Output

### Real-World Example: Kubernetes Pods

Here's a practical example showing tablefy's power with real data:

**Input: `kubectl get pods -A`**
```
NAMESPACE     NAME                                    READY  STATUS     RESTARTS  AGE
kube-system   coredns-5d78c0869f-7xk9q                1/1     Running   0         45d
kube-system   etcd-master                             1/1     Running   1         45d
kube-system   kube-apiserver-master                   1/1     Running   2         45d
default       my-app-deployment-7c5f4b8c9d-k2x5n      1/1     Running   0         3d
default       my-app-deployment-7c5f4b8c9d-m7q3p      1/1     Running   5         3d
ingress-nginx nginx-ingress-controller-9f8c7b6d5e     1/1     Running   1         10d
monitoring    prometheus-server-7b8c4d3e9f-a1s2k      0/1     Pending   0         2h
```

**Output Table:**
```
┌───────────────┬────────────────────────────────────────┬───────┬─────────┬──────────┬──────┐
│ NAMESPACE     │ NAME                                   │ READY │ STATUS  │ RESTARTS │ AGE  │
├───────────────┼────────────────────────────────────────┼───────┼─────────┼──────────┼──────┤
│ kube-system   │ coredns-5d78c0869f-7xk9q               │ 1/1   │ Running │ 0        │ 45d  │
│ kube-system   │ etcd-master                            │ 1/1   │ Running │ 1        │ 45d  │
│ kube-system   │ kube-apiserver-master                  │ 1/1   │ Running │ 2        │ 45d  │
│ default       │ my-app-deployment-7c5f4b8c9d-k2x5n     │ 1/1   │ Running │ 0        │ 3d   │
│ default       │ my-app-deployment-7c5f4b8c9d-m7q3p     │ 1/1   │ Running │ 5        │ 3d   │
│ ingress-nginx │ nginx-ingress-controller-9f8c7b6d5e    │ 1/1   │ Running │ 1        │ 10d  │
│ monitoring    │ prometheus-server-7b8c4d3e9f-a1s2k     │ 0/1   │ Pending │ 0        │ 2h   │
└───────────────┴────────────────────────────────────────┴───────┴─────────┴──────────┴──────┘

← → / h l: Navigate | s: Toggle select (0 selected) | Enter: Zoom | f: Filter | c: Clear | q: Quit
```

**Interactive Features in Action:**

1. **Navigate & Explore:**
   ```
   Press → to move to NAME column
   Press → again to see READY column
   The table scrolls to keep your focused column visible
   ```

2. **Fuzzy Filter by Status:**
   ```
   Press → → → → to navigate to STATUS column
   Press 'f' to start filtering
   Type "Pend" to find pending pods
   See live updates: "🔍 Filter [STATUS]: Pend (1 matches)"
   Press Enter to apply
   Now only the prometheus pod is shown (1/7 rows)
   ```

3. **Select & Zoom Columns:**
   ```
   Press 's' to select NAMESPACE column (turns purple)
   Press → → to navigate to NAME column
   Press 's' to select NAME as well (2 columns selected)
   Press Enter to zoom - see only NAMESPACE and NAME in focused view
   Press 'q' to return to full table
   ```

4. **Fast Navigation:**
   ```
   Press PgDn to scroll down one page (4-5 rows at once)
   Press PgUp to scroll back up quickly
   Great for exploring large datasets
   ```

5. **Auto-Expand for Details:**
   ```
   Run: kubectl get pods -A | tablefy --auto-expand
   Press → to navigate to NAME column
   The column auto-expands to show full pod names (no truncation)
   You can now see complete information without "..." placeholders
   ```

