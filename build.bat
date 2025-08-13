@echo off
echo Building TerraDrift Watcher...
echo.

REM Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)

REM Download dependencies
echo Downloading dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo Error: Failed to download dependencies
    exit /b 1
)

REM Build the binary
echo Building binary...
go build -o terradrift-watcher.exe .
if %errorlevel% neq 0 (
    echo Error: Build failed
    exit /b 1
)

echo.
echo Build successful!
echo Binary created: terradrift-watcher.exe
echo.
echo Run with: terradrift-watcher.exe run --config config.yml 