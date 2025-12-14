@echo off
cls

setlocal EnableDelayedExpansion
set "DIST_DIR=dist"
set "APP_NAME=bp-go-def.exe"
set "MAIN_PATH=./cmd/app/main.go"

:: Define ANSI Colors
for /F "tokens=1,2 delims=#" %%a in ('"prompt #$H#$E# & echo on & for %%b in (1) do rem"') do (
  set "ESC=%%b"
)
set "RESET=%ESC%[0m"
set "BOLD=%ESC%[1m"
set "DIM=%ESC%[2m"
set "UNDERLINE=%ESC%[4m"

:: Fancy Pastel Palette
set "P_PURPLE=%ESC%[95m"
set "B_PURPLE=%ESC%[1;95m"
set "P_CYAN=%ESC%[96m"
set "B_CYAN=%ESC%[1;96m"
set "P_GREEN=%ESC%[92m"
set "B_GREEN=%ESC%[1;92m"
set "P_YELLOW=%ESC%[93m"
set "B_YELLOW=%ESC%[1;93m"
set "P_RED=%ESC%[91m"
set "B_RED=%ESC%[1;91m"
set "GRAY=%ESC%[90m"
set "WHITE=%ESC%[97m"
set "B_WHITE=%ESC%[1;97m"

:: Robustly switch to project root
cd /d "%~dp0.."

echo.
echo    %P_PURPLE%(\_/)%RESET%
echo    %P_PURPLE%(o.o)%RESET%   %B_PURPLE%%APP_NAME% Builder%RESET% %GRAY%by%RESET% %B_WHITE%GabrielTenma%RESET%
echo   %P_PURPLE%c(")(")%RESET%
echo %GRAY%----------------------------------------------------------------------%RESET%

REM 1. Generate Timestamp from PowerShell
for /f "usebackq tokens=*" %%a in (`powershell -NoProfile -Command "Get-Date -Format 'yyyyMMdd_HHmmss'"`) do set "TIMESTAMP=%%a"
set "BACKUP_ROOT=%DIST_DIR%\backups"
set "BACKUP_PATH=%BACKUP_ROOT%\%TIMESTAMP%"

REM 2. Stop running process
echo %B_PURPLE%[1/4]%RESET% %P_CYAN%Checking for running process...%RESET%
tasklist /FI "IMAGENAME eq %APP_NAME%" 2>NUL | find /I /N "%APP_NAME%">NUL
if "%ERRORLEVEL%"=="0" (
    echo    %B_YELLOW%! App is running. Stopping...%RESET%
    taskkill /F /IM %APP_NAME% >NUL
    timeout /t 1 /nobreak >NUL
) else (
    echo    %B_GREEN%+ App is not running.%RESET%
)

REM 3. Backup Old Files
echo %B_PURPLE%[2/4]%RESET% %P_CYAN%Backing up old files...%RESET%
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
    
    echo    %B_GREEN%+ Backup created at:%RESET% %B_WHITE%%BACKUP_PATH%%RESET%
) else (
    echo    %GRAY%+ No existing dist directory. Skipping backup.%RESET%
    mkdir "%DIST_DIR%"
)

REM Ensure dist directory
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"

REM 4. Build
echo %B_PURPLE%[3/4]%RESET% %P_CYAN%Building Go binary...%RESET%
go build -o "%DIST_DIR%\%APP_NAME%" %MAIN_PATH%
if %ERRORLEVEL% NEQ 0 (
    echo    %B_RED%x Build FAILED! Exit code: %ERRORLEVEL%%RESET%
    exit /b %ERRORLEVEL%
)
echo    %B_GREEN%+ Build successful:%RESET% %B_WHITE%%DIST_DIR%\%APP_NAME%%RESET%

REM 5. Copy Assets
echo %B_PURPLE%[4/4]%RESET% %P_CYAN%Copying assets...%RESET%

if exist "web" (
    echo    %B_GREEN%+ Copying web folder...%RESET%
    xcopy /E /I /Y /Q "web" "%DIST_DIR%\web" >NUL
)

if exist "config.yaml" (
    echo    %B_GREEN%+ Copying config.yaml...%RESET%
    copy /Y "config.yaml" "%DIST_DIR%" >NUL
)

if exist "banner.txt" (
    echo    %B_GREEN%+ Copying banner.txt...%RESET%
    copy /Y "banner.txt" "%DIST_DIR%" >NUL
)

if exist "monitoring_users.db" (
    echo    %B_GREEN%+ Copying monitoring_users.db...%RESET%
    copy /Y "monitoring_users.db" "%DIST_DIR%" >NUL
)

echo.
echo %GRAY%======================================================================%RESET%
echo  %B_PURPLE%SUCCESS!%RESET% %GREEN%Build ready at:%RESET% %UNDERLINE%%B_WHITE%%DIST_DIR%\%RESET%
echo %GRAY%======================================================================%RESET%
endlocal
