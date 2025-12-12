#!/bin/bash

# Configuration
DIST_DIR="dist"
APP_NAME="bp-go-def"
MAIN_PATH="./cmd/app/main.go"

# Define ANSI Colors
RESET="\033[0m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
MAGENTA="\033[35m"
GRAY="\033[90m"

# Robustly switch to project root (one level up from this script)
cd "$(dirname "$0")/.." || exit 1

echo ""
echo -e "   ${MAGENTA}(\_/)${RESET}"
echo -e "   ${MAGENTA}(o.o)${RESET}   ${CYAN}${APP_NAME} Builder by GabrielTenma ${RESET}"
echo -e "  ${MAGENTA}c(\")(\")${RESET}"
echo -e "${GRAY}----------------------------------------------------------------------${RESET}"

# 1. Generate Timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_ROOT="${DIST_DIR}/backups"
BACKUP_PATH="${BACKUP_ROOT}/${TIMESTAMP}"

# 2. Stop running process
echo -e "${CYAN}[1/4] Checking for running process...${RESET}"
if pgrep -x "$APP_NAME" >/dev/null; then
    echo -e "   ${YELLOW}! App is running. Stopping...${RESET}"
    pkill -x "$APP_NAME"
    sleep 1
else
    echo -e "   ${GREEN}+ App is not running.${RESET}"
fi

# 3. Backup Old Files
echo -e "${CYAN}[2/4] Backing up old files...${RESET}"
if [ -d "$DIST_DIR" ]; then
    mkdir -p "$BACKUP_PATH"

    # Move old binary (check for both plain and .exe just in case)
    if [ -f "$DIST_DIR/$APP_NAME" ]; then
        echo -e "   ${GRAY}- Moving old binary...${RESET}"
        mv "$DIST_DIR/$APP_NAME" "$BACKUP_PATH/"
    elif [ -f "$DIST_DIR/$APP_NAME.exe" ]; then
        echo -e "   ${GRAY}- Moving old binary (.exe)...${RESET}"
        mv "$DIST_DIR/$APP_NAME.exe" "$BACKUP_PATH/"
    fi

    if [ -f "$DIST_DIR/config.yaml" ]; then
        mv "$DIST_DIR/config.yaml" "$BACKUP_PATH/"
    fi
    if [ -f "$DIST_DIR/banner.txt" ]; then
        mv "$DIST_DIR/banner.txt" "$BACKUP_PATH/"
    fi
    if [ -f "$DIST_DIR/monitoring_users.db" ]; then
        echo -e "   ${GRAY}- Backing up database...${RESET}"
        mv "$DIST_DIR/monitoring_users.db" "$BACKUP_PATH/"
    fi
    if [ -d "$DIST_DIR/web" ]; then
        echo -e "   ${GRAY}- Moving old web assets...${RESET}"
        mv "$DIST_DIR/web" "$BACKUP_PATH/"
    fi
    
    echo -e "   ${GREEN}+ Backup created at: ${BACKUP_PATH}${RESET}"
else
    echo -e "   ${GRAY}+ No existing dist directory. Skipping backup.${RESET}"
    mkdir -p "$DIST_DIR"
fi

# Ensure dist directory
mkdir -p "$DIST_DIR"

# 4. Build
echo -e "${CYAN}[3/4] Building Go binary...${RESET}"
go build -o "$DIST_DIR/$APP_NAME" "$MAIN_PATH"
if [ $? -ne 0 ]; then
    echo -e "   ${RED}x Build FAILED! Exit code: $?${RESET}"
    exit $?
fi
echo -e "   ${GREEN}+ Build successful: ${DIST_DIR}/${APP_NAME}${RESET}"

# 5. Copy Assets
echo -e "${CYAN}[4/4] Copying assets...${RESET}"

if [ -d "web" ]; then
    echo -e "   ${GREEN}+ Copying web folder...${RESET}"
    cp -r "web" "$DIST_DIR/web"
fi

if [ -f "config.yaml" ]; then
    echo -e "   ${GREEN}+ Copying config.yaml...${RESET}"
    cp "config.yaml" "$DIST_DIR/"
fi

if [ -f "banner.txt" ]; then
    echo -e "   ${GREEN}+ Copying banner.txt...${RESET}"
    cp "banner.txt" "$DIST_DIR/"
fi

if [ -f "monitoring_users.db" ]; then
    echo -e "   ${GREEN}+ Copying monitoring_users.db...${RESET}"
    cp "monitoring_users.db" "$DIST_DIR/"
fi

echo ""
echo -e "${GRAY}======================================================================${RESET}"
echo -e " ${GREEN}SUCCESS! Build ready at: ${DIST_DIR}/${RESET}"
echo -e "${GRAY}======================================================================${RESET}"
