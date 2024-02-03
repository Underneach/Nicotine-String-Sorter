@echo off
echo "Update packages"
go get -u all
echo "Make resources"
go-winres make
IF EXIST "cpu.pprof" (
    echo "Build with PGO"
    go build -o "Nicotine String Sorter.exe" -pgo="cpu.pprof"
) ELSE (
    echo "Build without PGO"
    go build -o "Nicotine String Sorter.exe"
)
echo "Done"
pause