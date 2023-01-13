@ECHO OFF

REM 1) Build unsigned binary
cd ..\\app
SETLOCAL
SET GOOS=windows
SET GOARCH=amd64
go build -o ../dist/pluralith_cli_windows_amd64.exe
ENDLOCAL

REM 2) Sign binary
cd ..\\dist
signtool sign /debug /n "Pluralith Industries Inc." /t http://time.certum.pl/ /fd sha256 /v ./pluralith_cli_windows_amd64.exe
