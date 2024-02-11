@echo off
set build_name="Nicotine String Sorter.exe"

IF EXIST %build_name% (
    echo "Delete old build"
    del %build_name%
)


echo "Update packages"
go get -u all


IF EXIST "rsrc_windows_amd64.syso" (
    echo "Resources already make"
) ELSE (
    echo "Make resources"
    go-winres make
)


IF EXIST "cpu.pprof" (
    echo "Build with PGO"
    go build -o %build_name% -pgo="cpu.pprof"
) ELSE (
    echo "Build without PGO"
    go build -o %build_name%
)


echo "Done"
pause