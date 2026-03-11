@echo off
setlocal enabledelayedexpansion

REM This script runs the Go server system from the data_labler_UI_production directory
REM with SQL database support WITHOUT vLLM server
REM Use this script when you don't need LLM functionality or want to run vLLM separately

echo Checking required commands...

REM Check for required commands
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH. Please install it first.
    echo Download from: https://go.dev/dl/
    exit /b 1
)

where ngrok >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: ngrok is not installed or not in PATH. Please install it first.
    echo Download from: https://ngrok.com/download
    exit /b 1
)

where sqlite3 >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: sqlite3 command-line tool not found. Database verification will be limited.
    echo Download from: https://www.sqlite.org/download.html
    set SQLITE3_AVAILABLE=0
) else (
    set SQLITE3_AVAILABLE=1
)

REM Configuration
set PORT=5002
set NGROK_HOST=wl2.studio
set DB_PATH=data\app.db

if "%PORT%"=="" (
    echo ERROR: PORT is empty. Please check the script configuration.
    exit /b 1
)

if "%NGROK_HOST%"=="" (
    echo ERROR: NGROK_HOST is empty. Please check the script configuration.
    exit /b 1
)

REM Get script directory and ensure we're in the right place
cd /d "%~dp0"
if not exist "cmd\server\main.go" (
    echo ERROR: Cannot find cmd\server\main.go
    echo This script must be run from the data_labler_UI_production directory
    exit /b 1
)

echo.
echo Cleaning up existing processes...

REM Kill processes using port 5002
echo Killing processes using port %PORT%...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%PORT%" ^| findstr "LISTENING"') do (
    echo Killing process %%a on port %PORT%
    taskkill /F /PID %%a >nul 2>&1
)

REM Kill processes using port 5003 (default server port)
echo Killing processes using port 5003...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":5003" ^| findstr "LISTENING"') do (
    echo Killing process %%a on port 5003
    taskkill /F /PID %%a >nul 2>&1
)

REM Kill existing ngrok processes
echo Killing existing ngrok processes...
tasklist /FI "IMAGENAME eq ngrok.exe" 2>NUL | find /I /N "ngrok.exe">NUL
if "%ERRORLEVEL%"=="0" (
    taskkill /F /IM ngrok.exe >nul 2>&1
    echo Killed existing ngrok processes
    timeout /t 1 >nul
) else (
    echo No ngrok processes found
)

echo.
echo Note: Running WITHOUT vLLM server. LLM features will not be available.

echo.
echo Checking database setup...

REM Check if database exists and is valid
if not exist "%DB_PATH%" (
    echo Database not found. Initializing database...
    
    echo Downloading Go dependencies ^(including SQLite driver^)...
    go mod download
    if !ERRORLEVEL! NEQ 0 (
        echo Failed to download dependencies. Exiting.
        exit /b 1
    )
    
    REM Check if there's JSON data to migrate
    if exist "data\users.json" (
        echo Found existing JSON data. Running migration...
        go run migrate_to_sql.go data
        if !ERRORLEVEL! NEQ 0 (
            echo Migration failed. Check migrate_to_sql.go for errors.
            exit /b 1
        )
        echo [OK] Migration completed successfully
    ) else (
        echo No existing JSON data found. Database will be initialized on first server start.
    )
) else (
    echo [OK] Database found at %DB_PATH%
    
    REM Verify database is valid (if sqlite3 is available)
    if !SQLITE3_AVAILABLE! EQU 1 (
        sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM users;" >nul 2>&1
        if !ERRORLEVEL! NEQ 0 (
            echo WARNING: Database appears to be corrupted or invalid.
            echo Backing up current database...
            set "BACKUP_NAME=%DB_PATH%.backup.%date:~-4%%date:~4,2%%date:~7,2%_%time:~0,2%%time:~3,2%%time:~6,2%"
            set "BACKUP_NAME=!BACKUP_NAME: =0!"
            move "%DB_PATH%" "!BACKUP_NAME!" >nul
            echo Database backed up. A new database will be created on server start.
        ) else (
            for /f %%i in ('sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM users;" 2^>nul') do set USER_COUNT=%%i
            for /f %%i in ('sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM lounge_posts;" 2^>nul') do set POST_COUNT=%%i
            echo [OK] Database is valid ^(Users: !USER_COUNT!, Posts: !POST_COUNT!^)
        )
    )
)

echo.
echo Starting the Go backend server in DEBUG mode...

REM Set environment variables
set DEBUG=true
set SERVE_PORT=%PORT%

echo Removing old server binary...
if exist server_sql.exe del /f server_sql.exe

echo Building SQL-based server binary...
go build -o server_sql.exe .\cmd\server
if %ERRORLEVEL% NEQ 0 (
    echo Build failed. Exiting.
    echo Make sure you have run 'go mod tidy' to update dependencies.
    exit /b 1
)

echo [OK] Build completed

REM Start server as detached background process (no extra cmd window)
echo Starting server...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "Start-Process cmd.exe -ArgumentList '/c','server_sql.exe >> server.log 2>&1' -WorkingDirectory '%CD%' -WindowStyle Hidden"

REM Give the server a moment to start up
timeout /t 3 >nul

REM Verify the server is running and listening on the port
netstat -ano | findstr ":%PORT%" | findstr "LISTENING" >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Server does not appear to be listening on port %PORT%
    echo Check server.log for errors:
    echo.
    if exist server.log (
        powershell -Command "Get-Content server.log -Tail 20"
    )
    exit /b 1
)

echo [OK] Server is running and listening on port %PORT%

echo.
echo Starting ngrok...
REM Start ngrok as detached background process (no extra cmd window)
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "Start-Process cmd.exe -ArgumentList '/c','ngrok http --domain=%NGROK_HOST% http://127.0.0.1:%PORT% >> ngrok.log 2>&1' -WorkingDirectory '%CD%' -WindowStyle Hidden"

timeout /t 2 >nul

REM Verify ngrok is running
tasklist /FI "IMAGENAME eq ngrok.exe" 2>NUL | find /I /N "ngrok.exe">NUL
if "%ERRORLEVEL%" NEQ "0" (
    echo ERROR: ngrok does not appear to be running.
    if exist ngrok.log (
        echo Last ngrok log lines:
        powershell -Command "Get-Content ngrok.log -Tail 20"
    ) else (
        echo ngrok.log not found. Check ngrok auth or domain settings.
    )
    exit /b 1
)

echo [OK] ngrok tunnel is running at https://%NGROK_HOST%

echo.
echo ==========================================
echo Application and ngrok are running in the background.
echo ==========================================
echo.
echo Database: %CD%\%DB_PATH%
echo Logs:
echo   - Backend: %CD%\server.log
echo   - ngrok:   %CD%\ngrok.log
echo.
echo Access the application at:
echo   - Local:  http://localhost:%PORT%
echo   - Public: https://%NGROK_HOST%
echo.
echo [NOTE] IMPORTANT: Clear your browser cache to see updates.
echo   - Chrome/Edge: Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)
echo   - Firefox: Ctrl+F5 (Windows/Linux) or Cmd+Shift+R (Mac)
echo   - Or open in Incognito/Private mode
echo.
echo To stop the services, manually kill processes:
echo     taskkill /F /IM server_sql.exe
echo     taskkill /F /IM ngrok.exe
echo.
echo To view database contents:
echo   sqlite3 %CD%\%DB_PATH%
echo   .tables
echo   SELECT * FROM users;
echo   .quit
echo.
echo ==========================================
echo.
echo Watchdog is now monitoring services (every 60s).
echo Press Ctrl+C to stop monitoring (services will keep running).
echo Restart events are logged to: %CD%\watchdog.log
echo.

:watchdog_loop
timeout /t 60 /nobreak >nul

set SERVER_DOWN=0
set NGROK_DOWN=0

REM Check if server_sql.exe process is alive
tasklist /FI "IMAGENAME eq server_sql.exe" 2>NUL | find /I /N "server_sql.exe">NUL
if !ERRORLEVEL! NEQ 0 (
    set SERVER_DOWN=1
) else (
    REM Process exists -- verify it's actually listening on the port
    netstat -ano | findstr ":%PORT%" | findstr "LISTENING" >nul 2>&1
    if !ERRORLEVEL! NEQ 0 (
        set SERVER_DOWN=1
    )
)

REM Check if ngrok.exe process is alive
tasklist /FI "IMAGENAME eq ngrok.exe" 2>NUL | find /I /N "ngrok.exe">NUL
if !ERRORLEVEL! NEQ 0 set NGROK_DOWN=1

if !SERVER_DOWN! EQU 0 if !NGROK_DOWN! EQU 0 goto watchdog_loop

REM --- Something is down, take action ---
for /f "tokens=*" %%t in ('powershell -NoProfile -Command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'"') do set TIMESTAMP=%%t

if !SERVER_DOWN! EQU 1 (
    echo [!TIMESTAMP!] WATCHDOG: server_sql.exe is DOWN - restarting...
    echo [!TIMESTAMP!] WATCHDOG: server_sql.exe is DOWN - restarting >> watchdog.log

    REM Kill any zombie server process that might still linger
    taskkill /F /IM server_sql.exe >nul 2>&1
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%PORT%" ^| findstr "LISTENING"') do (
        taskkill /F /PID %%a >nul 2>&1
    )
    timeout /t 1 >nul

    REM Restart server (append to log so we don't lose crash context)
    powershell -NoProfile -ExecutionPolicy Bypass -Command ^
        "Start-Process cmd.exe -ArgumentList '/c','server_sql.exe >> server.log 2>&1' -WorkingDirectory '%CD%' -WindowStyle Hidden"
    timeout /t 3 >nul

    REM Verify it came back
    netstat -ano | findstr ":%PORT%" | findstr "LISTENING" >nul 2>&1
    if !ERRORLEVEL! EQU 0 (
        echo [!TIMESTAMP!] WATCHDOG: server_sql.exe restarted OK
        echo [!TIMESTAMP!] WATCHDOG: server_sql.exe restarted OK >> watchdog.log
    ) else (
        echo [!TIMESTAMP!] WATCHDOG: server_sql.exe FAILED to restart - will retry next cycle
        echo [!TIMESTAMP!] WATCHDOG: server_sql.exe FAILED to restart >> watchdog.log
    )
)

if !NGROK_DOWN! EQU 1 (
    echo [!TIMESTAMP!] WATCHDOG: ngrok.exe is DOWN - restarting...
    echo [!TIMESTAMP!] WATCHDOG: ngrok.exe is DOWN - restarting >> watchdog.log

    taskkill /F /IM ngrok.exe >nul 2>&1
    timeout /t 1 >nul

    powershell -NoProfile -ExecutionPolicy Bypass -Command ^
        "Start-Process cmd.exe -ArgumentList '/c','ngrok http --domain=%NGROK_HOST% http://127.0.0.1:%PORT% >> ngrok.log 2>&1' -WorkingDirectory '%CD%' -WindowStyle Hidden"
    timeout /t 3 >nul

    tasklist /FI "IMAGENAME eq ngrok.exe" 2>NUL | find /I /N "ngrok.exe">NUL
    if !ERRORLEVEL! EQU 0 (
        echo [!TIMESTAMP!] WATCHDOG: ngrok.exe restarted OK
        echo [!TIMESTAMP!] WATCHDOG: ngrok.exe restarted OK >> watchdog.log
    ) else (
        echo [!TIMESTAMP!] WATCHDOG: ngrok.exe FAILED to restart - will retry next cycle
        echo [!TIMESTAMP!] WATCHDOG: ngrok.exe FAILED to restart >> watchdog.log
    )
)

goto watchdog_loop

