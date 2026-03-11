@echo off
setlocal

REM This script stops all running services (Go server and ngrok)

echo Stopping all services...
echo.

set STOPPED_SOMETHING=0

REM Stop processes on port 5002
echo Checking for processes on port 5002...
set FOUND_5002=0
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":5002" ^| findstr "LISTENING" 2^>nul') do (
    set FOUND_5002=1
    echo Stopping process %%a on port 5002...
    taskkill /F /PID %%a >nul 2>&1
    set STOPPED_SOMETHING=1
)
if %FOUND_5002%==1 (
    echo ? Stopped processes on port 5002
    timeout /t 1 >nul
) else (
    echo   No processes found on port 5002
)

REM Stop processes on port 5003
echo Checking for processes on port 5003...
set FOUND_5003=0
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":5003" ^| findstr "LISTENING" 2^>nul') do (
    set FOUND_5003=1
    echo Stopping process %%a on port 5003...
    taskkill /F /PID %%a >nul 2>&1
    set STOPPED_SOMETHING=1
)
if %FOUND_5003%==1 (
    echo ? Stopped processes on port 5003
    timeout /t 1 >nul
) else (
    echo   No processes found on port 5003
)

REM Stop ngrok processes
echo Checking for ngrok processes...
tasklist /FI "IMAGENAME eq ngrok.exe" 2>NUL | find /I /N "ngrok.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo Stopping ngrok...
    taskkill /F /IM ngrok.exe >nul 2>&1
    echo ? Stopped ngrok
    set STOPPED_SOMETHING=1
    timeout /t 1 >nul
) else (
    echo   No ngrok processes found
)

REM Stop server_sql.exe processes
echo Checking for server_sql.exe processes...
tasklist /FI "IMAGENAME eq server_sql.exe" 2>NUL | find /I /N "server_sql.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo Stopping server_sql.exe...
    taskkill /F /IM server_sql.exe >nul 2>&1
    echo ? Stopped server_sql.exe
    set STOPPED_SOMETHING=1
    timeout /t 1 >nul
) else (
    echo   No server_sql.exe processes found
)

REM Stop any cmd.exe processes running server logs (background processes from run_go3.bat)
echo Checking for background server processes...
set FOUND_CMD=0
for /f "tokens=2" %%a in ('tasklist /FI "IMAGENAME eq cmd.exe" /FO LIST ^| findstr /C:"PID:"') do (
    REM Check if this cmd process has server_sql as a child
    for /f %%b in ('wmic process where "ParentProcessId=%%a" get Name 2^>nul ^| findstr /I "server_sql"') do (
        set FOUND_CMD=1
        echo Stopping background process %%a...
        taskkill /F /PID %%a >nul 2>&1
        set STOPPED_SOMETHING=1
    )
)
if %FOUND_CMD%==1 (
    echo ? Stopped background processes
)

echo.
echo ==========================================
if %STOPPED_SOMETHING%==1 (
    echo ? All services stopped successfully
) else (
    echo ? No services were running
)
echo ==========================================
echo.
echo You can now safely start services again with:
echo   - Windows: run_go3.bat (from this directory^)
echo   - Linux/Mac: ./run_go3.sh (from parent directory^)
echo.
pause

