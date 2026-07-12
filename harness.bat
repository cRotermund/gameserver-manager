@echo off
setlocal
cd /d "%~dp0harness"
go run . %*