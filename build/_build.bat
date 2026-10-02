@echo off
rem Builds altping for one platform into bin\.
rem Usage: build\_build.bat <goos> <goarch>
setlocal

set "T_OS=%~1"
set "T_ARCH=%~2"
set "APP=altping"
set "NFPM_VERSION=v2.47.0"
pushd "%~dp0.."

set "VERSION=dev"
for /f "delims=" %%v in ('git describe --tags --always --dirty 2^>nul') do set "VERSION=%%v"

set "LDFLAGS=-s -w -X github.com/ipoluianov/altping/app.Version=%VERSION%"
set "EXT="
if "%T_OS%"=="windows" (
  set "EXT=.exe"
  set "LDFLAGS=%LDFLAGS% -H=windowsgui"
)

set "OUT=bin\%APP%-%T_OS%-%T_ARCH%%EXT%"
if not exist bin mkdir bin
echo Building %OUT% (%VERSION%)
set "CGO_ENABLED=0"
set "GOOS=%T_OS%"
set "GOARCH=%T_ARCH%"
go build -trimpath -ldflags="%LDFLAGS%" -o "%OUT%" .
set "RC=%ERRORLEVEL%"
rem Tools below (go run nfpm) must be built for the host
set "GOOS="
set "GOARCH="
set "CGO_ENABLED="
if not "%RC%"=="0" goto done

if "%T_OS%"=="linux" goto linux
if "%T_OS%"=="darwin" goto darwin
goto done

:linux
rem nfpm is pure Go, so packages can be built on any OS
if not exist bin\.pkg mkdir bin\.pkg
copy /y "%OUT%" bin\.pkg\altping >nul
set "PKG_ARCH=%T_ARCH%"
set "PKG_VERSION=%VERSION%"
for %%f in (deb rpm) do (
  go run "github.com/goreleaser/nfpm/v2/cmd/nfpm@%NFPM_VERSION%" pkg --config build/nfpm.yaml --packager %%f --target "%OUT%.%%f" || set "RC=1"
)
rmdir /s /q bin\.pkg
goto done

:darwin
rem The .dmg is packed by _dmg.sh, so it needs Git Bash
rem (bash.exe from System32 is the WSL launcher, not used here)
set "GITBASH="
for /f "delims=" %%g in ('where git 2^>nul') do (
  if not defined GITBASH if exist "%%~dpg..\bin\bash.exe" set "GITBASH=%%~dpg..\bin\bash.exe"
)
if not defined GITBASH if exist "%ProgramFiles%\Git\bin\bash.exe" set "GITBASH=%ProgramFiles%\Git\bin\bash.exe"
if not defined GITBASH (
  echo Warning: Git Bash not found, skipping %OUT%.dmg
  goto done
)
"%GITBASH%" build/_dmg.sh "%OUT:\=/%" "%VERSION%"
set "RC=%ERRORLEVEL%"

:done
popd
exit /b %RC%
