@echo off
setlocal EnableDelayedExpansion
set "DIST_DIR=dist"
set "APP_NAME=bp-go-def.exe"
set "MAIN_PATH=./cmd/app/main.go"

:: Define ANSI Colors
for /F "tokens=1,2 delims=#" %%a in ('"prompt #$H#$E# & echo on & for %%b in (1) do rem"') do (
  set "ESC=%%b"
)
set "RESET=%ESC%[0m"
set "CYAN=%ESC%[36m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "RED=%ESC%[31m"
set "MAGENTA=%ESC%[35m"
set "GRAY=%ESC%[90m"

:: Robustly switch to project root
cd /d "%~dp0.."

echo.
echo    %MAGENTA%(\_/)%RESET%
echo    %MAGENTA%(o.o)%RESET%   %CYAN%%APP_NAME% Builder by GabrielTenma %RESET%
echo   %MAGENTA%c(")(")%RESET%
echo %GRAY%----------------------------------------------------------------------%RESET%

REM 1. Generate Timestamp from PowerShell
for /f "usebackq tokens=*" %%a in (`powershell -NoProfile -Command "Get-Date -Format 'yyyyMMdd_HHmmss'"`) do set "TIMESTAMP=%%a"
set "BACKUP_ROOT=%DIST_DIR%\backups"
set "BACKUP_PATH=%BACKUP_ROOT%\%TIMESTAMP%"

REM 2. Stop running process
echo %CYAN%[1/4] Checking for running process...%RESET%
tasklist /FI "IMAGENAME eq %APP_NAME%" 2>NUL | find /I /N "%APP_NAME%">NUL
if "%ERRORLEVEL%"=="0" (
    echo    %YELLOW%! App is running. Stopping...%RESET%
    taskkill /F /IM %APP_NAME% >NUL
    timeout /t 1 /nobreak >NUL
) else (
    echo    %GREEN%+ App is not running.%RESET%
)

REM 3. Backup Old Files
echo %CYAN%[2/4] Backing up old files...%RESET%
if exist "%DIST_DIR%" (
    if not exist "%BACKUP_PATH%" mkdir "%BACKUP_PATH%"
    
    if exist "%DIST_DIR%\%APP_NAME%" (
        echo    %GRAY%- Moving old binary...%RESET%
        move "%DIST_DIR%\%APP_NAME%" "%BACKUP_PATH%\" >NUL
    )
    if exist "%DIST_DIR%\config.yaml" (
        move "%DIST_DIR%\config.yaml" "%BACKUP_PATH%\" >NUL
    )
    if exist "%DIST_DIR%\banner.txt" (
        move "%DIST_DIR%\banner.txt" "%BACKUP_PATH%\" >NUL
    )
    if exist "%DIST_DIR%\monitoring_users.db" (
        echo    %GRAY%- Backing up database...%RESET%
        move "%DIST_DIR%\monitoring_users.db" "%BACKUP_PATH%\" >NUL
    )
    if exist "%DIST_DIR%\web" (
        echo    %GRAY%- Moving old web assets...%RESET%
        move "%DIST_DIR%\web" "%BACKUP_PATH%\" >NUL
    )
    
    echo    %GREEN%+ Backup created at: %BACKUP_PATH%%RESET%
) else (
    echo    %GRAY%+ No existing dist directory. Skipping backup.%RESET%
    mkdir "%DIST_DIR%"
)

REM Ensure dist directory
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"

REM 4. Build
echo %CYAN%[3/4] Building Go binary...%RESET%
go build -o "%DIST_DIR%\%APP_NAME%" %MAIN_PATH%
if %ERRORLEVEL% NEQ 0 (
    echo    %RED%x Build FAILED! Exit code: %ERRORLEVEL%%RESET%
    exit /b %ERRORLEVEL%
)
echo    %GREEN%+ Build successful: %DIST_DIR%\%APP_NAME%%RESET%

REM 5. Copy Assets
echo %CYAN%[4/4] Copying assets...%RESET%

if exist "web" (
    echo    %GREEN%+ Copying web folder...%RESET%
    xcopy /E /I /Y /Q "web" "%DIST_DIR%\web" >NUL
)

if exist "config.yaml" (
    echo    %GREEN%+ Copying config.yaml...%RESET%
    copy /Y "config.yaml" "%DIST_DIR%" >NUL
)

if exist "banner.txt" (
    echo    %GREEN%+ Copying banner.txt...%RESET%
    copy /Y "banner.txt" "%DIST_DIR%" >NUL
)

if exist "monitoring_users.db" (
    echo    %GREEN%+ Copying monitoring_users.db...%RESET%
    copy /Y "monitoring_users.db" "%DIST_DIR%" >NUL
)

echo.
echo %GRAY%======================================================================%RESET%
echo  %GREEN%SUCCESS! Build ready at: %DIST_DIR%\%RESET%
echo %GRAY%======================================================================%RESET%
endlocal
