# Time Tracker CLI

A personal time tracking tool built in Go that helps you manage and track your work hours. This is a lightweight, terminal-based application that stores time records locally using SQLite and provides both a command-line interface (REPL) and a system tray GUI.

## Features

- 📊 **Time Tracking**: Start/stop time recording with precise tracking
- 💾 **Local Storage**: SQLite database for reliable local data storage
- 🖥️ **Dual Interface**: Terminal REPL + System tray GUI
- ⏰ **Focus Timer**: Built-in Pomodoro-style focus sessions
- 📅 **Date Management**: Work across different dates
- 🪣 **Pool System**: Manage overflow hours for flexible time allocation
- 📈 **Debt Tracking**: Monitor accumulated work debt
- 🎯 **Status Management**: Track pending, committed, and pooled time

## Installation

### Prerequisites
- Go 1.21 or higher
- A terminal that supports ANSI colors

### Building from Source

```bash
# Clone or download the repository
cd time-tracker-cli

# Build the application
sh build.sh
```

This creates a binary called `tt` in the `build/` directory.

### Running the Application

```bash
./build/tt
```

On first run, the application will:
1. Create necessary directory structure
2. Initialize SQLite database
3. Create a configuration file (`config.json`)
4. Start both the terminal interface and system tray GUI

## Configuration

The application creates a `config.json` file with the following settings:

```json
{
    "logLevel": "info",
    "workingTime": 8.0
}
```

- `logLevel`: Controls logging verbosity (error, info, debug)
- `workingTime`: Your daily working hours (used for debt calculation)

## Command Line Interface

The tool provides a REPL (Read-Eval-Print Loop) interface. Type `help` to see available commands:

### Basic Usage

```bash
tt > help
```

### Quick Command Syntax

Commands support both interactive and quick syntax:

**Interactive:**
```bash
tt > rec
- At: 09:00
**** Record started at 09:00! ****
```

**Quick syntax:**
```bash
tt > rec at; 09:00
**** Record started at 09:00! ****
```

## Available Commands

### Time Recording
- **`rec`** - Start a new time recorder
- **`rec at`** - Start recording at a specific time (format: HH:MM)
- **`end`** - End current time recording
- **`end at`** - End recording at a specific time
- **`drop`** - Drop current recording without saving
- **`add`** - Manually add a time record with specified hours

### Record Management
- **`list`** - Show all records for current date
- **`commit`** - Mark records as committed (accepts amount parameter)
- **`send pool`** - Send pending records to the pool
- **`poure`** - Pour pool time to current date

### Navigation & Utilities
- **`change date`** - Change working date (formats: `yy-mm-dd`, `yesterday`, `now`, `±N` days)
- **`debt`** - Show accumulated work debt
- **`help`** - Show command list

## System Tray GUI

The application also provides a system tray interface with:

- 🔥 **Focus Timer**: 2-minute focus sessions with visual feedback
- 🍅 **Pomodoro Timer**: 25-minute work sessions
- ⏯️ **Start/Stop**: Quick recording controls
- 📊 **Status Display**: Real-time work statistics in tooltip

### GUI States
- **💤 Idle**: Not currently recording
- **🔥 Focus**: In a focus session (2 min)
- **🍅 Pomodoro**: In a Pomodoro session (25 min)
- **⌛ Working**: Recording without timers

## Workflow

### Basic Daily Workflow

1. **Start Recording**: Use `rec` to begin tracking time
2. **Work Sessions**: Use focus/Pomodoro timers for productivity
3. **End Recording**: Use `end` to stop and save your work session
4. **Review**: Use `list` to see your recorded time
5. **Commit**: Use `commit` to mark time as final
6. **Manage Overflow**: Use pool system for time that exceeds daily limits

### Example Session

```bash
# Start recording
tt > rec
**** Record started! ****

# Check current status
[Rec:0.25]['] tt > list
Result
1. 0.25

# End recording after work
tt > end
**** 2.50 hours inserted! ****

# Commit your time
tt > commit; 8
**** Records committed! ****
```

## Status Bar

The terminal prompt shows real-time status information:

```bash
[Debt:2.0][Worked:1.25][Commited:6.00][Pool:0.50][Rec:0.25]['] tt >
```

- **Debt**: Accumulated work time debt (hours behind)
- **Worked**: Time worked today (not yet committed)
- **Committed**: Time marked as committed today
- **Pool**: Available time in the pool
- **Rec**: Currently recording time
- **Date**: Shows date if not today (format: yy-mm-dd)

## Pool System

The pool manages time that exceeds your configured daily working hours:

- **Overflow Protection**: When committing more than your daily limit, excess goes to pool
- **Weekend Work**: Weekend recordings automatically go to pool
- **Flexible Allocation**: Pour pool time to any date when needed

```bash
# Send today's work to pool
tt > send pool
**** Records saved to pool! ****

# Pour pool time to current date
tt > poure
**** Pool poured! ****
```

## Date Management

The tool works with today's date by default, but supports flexible date navigation:

### Date Commands

```bash
# Change to specific date
tt > change date; 24-05-23

# Quick date shortcuts
tt > change date; yesterday
tt > change date; now        # back to today
tt > change date; -3         # 3 days ago
tt > change date; 5          # 5 days in future
```

### Supported Date Formats
- `yy-mm-dd`: Specific date (e.g., `24-05-23`)
- `yesterday`: Previous day
- `now`, `today`, or empty: Current day
- `±N`: Relative days (e.g., `-3` for 3 days ago, `5` for 5 days ahead)

## Time Recording Details

### Time Precision
- Time is recorded in 1-minute precision
- Display rounds to nearest minute for readability
- Internal calculations maintain full precision

### Recording States
- **Pending**: New records waiting to be committed
- **Committed**: Records marked as final work
- **Pool**: Overflow or weekend work available for allocation

### Smart Features
- **Weekend Detection**: Weekend work automatically goes to pool
- **Overflow Management**: Excess daily time automatically pooled
- **Debt Tracking**: Monitors accumulated work debt across weekdays

## File Structure

The application creates the following structure in the executable's directory:

```
build/
├── tt              # Main executable
├── config.json     # Configuration file
├── tt.db          # SQLite database
└── tt.log         # Application logs
```

### Configuration File (`config.json`)
```json
{
    "logLevel": "info",
    "workingTime": 8.0
}
```

### Database Schema
The SQLite database contains:
- **records**: Time entries with ID, date, status, and hours
- **state_variables**: Application state (e.g., current recording)

## Dependencies

The application uses the following Go modules:

- `github.com/getlantern/systray` - System tray GUI
- `github.com/jmoiron/sqlx` - Database operations
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/sirupsen/logrus` - Logging
- `golang.org/x/term` - Terminal interface
- `github.com/google/uuid` - UUID generation

## Troubleshooting

### Common Issues

1. **Terminal not supported**: Ensure your terminal supports ANSI colors and terminal mode
2. **Permissions**: Make sure the build directory is writable
3. **System tray not showing**: Some desktop environments require additional setup

### Logs

Check `tt.log` in the build directory for detailed error information:

```bash
tail -f build/tt.log
```

### Reset Configuration

To reset the application:
```bash
rm build/config.json build/tt.db
```

## Development

### Project Structure

```
├── cmd/           # Application entry point
├── tt/
│   ├── app/       # Application layer
│   ├── domain/    # Business logic
│   └── infrastructure/
│       ├── cmd/   # Command handling
│       ├── config/ # Configuration
│       ├── display/ # GUI components
│       └── repositories/ # Data access
└── build.sh       # Build script
```

### Building for Development

```bash
go mod download
go build -o build/tt cmd/*.go
```

## License

MIT License - see LICENSE file for details.
