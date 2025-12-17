# Build Scripts Documentation

## Overview

The build scripts (`scripts/build.sh` for Unix-like systems and `scripts/build.bat` for Windows) automate the process of building the Go application, managing backups, and archiving previous builds. These scripts ensure a clean deployment by stopping running processes, backing up existing files, creating compressed archives of backups, building the application, and copying required assets.

## Features

- **Process Management**: Automatically stops any running instances of the application
- **Backup Creation**: Moves old build files to a timestamped backup directory
- **Archive Compression**: Creates ZIP archives of backup folders for long-term storage
- **Clean Builds**: Builds the Go application from source
- **Asset Management**: Copies configuration files and web assets to the build directory
- **Cross-Platform**: Separate scripts optimized for Unix/Linux/macOS and Windows environments

## Prerequisites

### For Unix/Linux/macOS (`build.sh`)
- Bash shell
- Go compiler installed and in PATH
- `zip` utility for archive compression
- `pgrep` and `pkill` commands for process management

### For Windows (`build.bat`)
- Windows Command Prompt or PowerShell
- Go compiler installed and in PATH
- PowerShell (included with modern Windows) for timestamp generation and archiving

## Directory Structure

After running the build script, the project directory will contain:

```
project/
├── dist/                    # Build output directory
│   ├── bp-go-def           # Compiled binary (Unix/Linux/macOS)
│   ├── bp-go-def.exe       # Compiled binary (Windows)
│   ├── config.yaml         # Configuration file
│   ├── banner.txt          # Banner file
│   ├── monitoring_users.db # Database file
│   ├── web/               # Web assets
│   └── backups/           # Backup archives
│       ├── 20251217_002400.zip
│       └── 20251217_002500.zip
├── scripts/
│   ├── build.sh
│   └── build.bat
└── ... (other project files)
```

## Build Process

The scripts execute the following steps in order:

### 1. Timestamp Generation
- Creates a timestamp in `YYYYMMDD_HHMMSS` format
- Used for backup directory naming and archive filenames

### 2. Process Management
- Checks for running instances of the application
- Stops any running processes using the application name
- Waits briefly for processes to terminate

### 3. Backup Creation
- Creates a backup directory with the timestamp
- Moves existing build files to the backup directory:
  - Application binary (`bp-go-def` or `bp-go-def.exe`)
  - Configuration files (`config.yaml`, `banner.txt`)
  - Database files (`monitoring_users.db`)
  - Web assets directory (`web/`)

### 4. Backup Archiving
- Compresses the backup directory into a ZIP archive
- Removes the uncompressed backup directory to save space
- Archive is stored in `dist/backups/` with timestamp filename

### 5. Build Process
- Changes to the project root directory
- Executes `go build` to compile the application
- Outputs the binary to the `dist/` directory

### 6. Asset Copying
- Copies configuration files to `dist/`
- Copies web assets to `dist/web/`
- Copies database files to `dist/`

## Usage

### Unix/Linux/macOS

```bash
# Make the script executable (first time only)
chmod +x scripts/build.sh

# Run the build script
./scripts/build.sh
```

### Windows

```cmd
# Run the build script
scripts\build.bat
```

Or from PowerShell:

```powershell
.\scripts\build.bat
```

## Configuration

The scripts use the following configuration variables that can be modified at the top of each script:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `DIST_DIR` | Output directory for builds | `dist` |
| `APP_NAME` | Application binary name | `bp-go-def` (Unix) / `bp-go-def.exe` (Windows) |
| `MAIN_PATH` | Path to main Go file | `./cmd/app/main.go` |

## Output and Logging

The scripts provide colored console output indicating the progress:

- **Process Check**: Shows if the application is running and stops it if needed
- **Backup Status**: Displays backup creation and archiving progress
- **Build Status**: Shows Go compilation results
- **Asset Copying**: Indicates successful file copying operations

Example output:
```
   (\/)
   (o.o)   bp-go-def Builder by GabrielTenma
  c(")(")
----------------------------------------------------------------------
[1/5] Checking for running process...
   + App is not running.
[2/5] Backing up old files...
   + Backup created at: dist/backups/20251217_002400
[3/5] Archiving backup...
   + Backup archived: dist/backups/20251217_002400.zip
[4/5] Building Go binary...
   + Build successful: dist/bp-go-def
[5/5] Copying assets...
   + Copying web folder...
   + Copying config.yaml...
   + Copying banner.txt...
   + Copying monitoring_users.db...

======================================================================
 SUCCESS! Build ready at: dist/
======================================================================
```

## Error Handling

The scripts include error checking at critical points:

- **Build Failures**: If `go build` fails, the script exits with the error code
- **Permission Issues**: May require appropriate permissions for file operations
- **Missing Dependencies**: Will fail if Go compiler or required utilities are not available

## Backup Management

### Archive Format
- Backups are compressed using ZIP format
- Filenames include timestamps for easy identification
- Archives preserve directory structure

### Cleanup
- Uncompressed backup directories are automatically removed after archiving
- Old archives are not automatically deleted (manual cleanup recommended)

### Restoration
To restore from a backup archive:

```bash
# Extract archive (Unix/Linux/macOS)
unzip dist/backups/20251217_002400.zip -d temp_restore/

# Extract archive (Windows)
powershell Expand-Archive -Path dist/backups/20251217_002400.zip -DestinationPath temp_restore/
```

## Troubleshooting

### Common Issues

**"Permission denied" on Unix/Linux/macOS**
- Ensure the script has execute permissions: `chmod +x scripts/build.sh`
- Check file permissions in the project directory

**"go: command not found"**
- Ensure Go is installed and added to PATH
- Verify installation: `go version`

**"zip: command not found"**
- Install zip utility: `sudo apt-get install zip` (Ubuntu/Debian) or equivalent

**Build fails with compilation errors**
- Check Go code for syntax errors
- Ensure all dependencies are installed: `go mod tidy`

**Process not stopping**
- The script uses `pkill` (Unix) or `taskkill` (Windows) to stop processes
- May require manual intervention if processes don't respond

### Debug Mode

For additional logging, you can modify the scripts to enable verbose output or add debug statements.

## Best Practices

1. **Regular Backups**: Run the build script regularly to maintain backup archives
2. **Storage Management**: Periodically clean up old backup archives to save disk space
3. **Version Control**: Consider committing backup archives to version control for critical deployments
4. **Testing**: Test builds in a staging environment before production deployment
5. **Permissions**: Ensure the build user has appropriate permissions for all operations

## Integration with CI/CD

These build scripts can be integrated into CI/CD pipelines:

### GitHub Actions Example

```yaml
- name: Build Application
  run: |
    chmod +x scripts/build.sh
    ./scripts/build.sh

- name: Archive Build Artifacts
  uses: actions/upload-artifact@v2
  with:
    name: build-artifacts
    path: dist/
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'chmod +x scripts/build.sh'
                sh './scripts/build.sh'
            }
        }
        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'dist/**', fingerprint: true
            }
        }
    }
}
```

## Security Considerations

- **Sensitive Files**: Be cautious with backup archives containing sensitive configuration or database files
- **Permissions**: Ensure backup directories have appropriate access controls
- **Encryption**: Consider encrypting backup archives for sensitive deployments
- **Cleanup**: Remove temporary files and sensitive data from build artifacts
