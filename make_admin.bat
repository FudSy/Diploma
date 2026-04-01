@echo off
setlocal enabledelayedexpansion

set "LOGIN=%~1"
if "%LOGIN%"=="" (
    echo Usage: %~nx0 ^<login^> [--docker]
    exit /b 1
)

set "MODE=%~2"

if not defined DB_HOST     set "DB_HOST=localhost"
if not defined DB_PORT     set "DB_PORT=5433"
if not defined DB_NAME     set "DB_NAME=test"
if not defined DB_USER     set "DB_USER=postgres"
if not defined DB_PASSWORD set "DB_PASSWORD=postgres"
if not defined DOCKER_CONTAINER set "DOCKER_CONTAINER=diploma-postgres"

if /i "%MODE%"=="--docker" (
    echo Connecting via Docker container: %DOCKER_CONTAINER%
    docker exec %DOCKER_CONTAINER% psql -U %DB_USER% -d %DB_NAME% -c "UPDATE users SET role = 'ADMIN' WHERE login = '%LOGIN%';" -c "SELECT login, role FROM users WHERE login = '%LOGIN%';"
) else (
    echo Connecting to local DB %DB_HOST%:%DB_PORT%/%DB_NAME%
    set "PGPASSWORD=%DB_PASSWORD%"
    psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "UPDATE users SET role = 'ADMIN' WHERE login = '%LOGIN%';" -c "SELECT login, role FROM users WHERE login = '%LOGIN%';"
)

if %ERRORLEVEL% neq 0 (
    echo ERROR: Could not connect to database.
    exit /b 1
)

echo Done. User '%LOGIN%' is now ADMIN. Re-login required.

endlocal
