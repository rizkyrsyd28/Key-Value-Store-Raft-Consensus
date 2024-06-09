@echo off

:: VARIABLE
set PROTO_DIR=protos
set PROTO_DIST=.
set TEST_DIR=.\test

set SERVER=server
set CLIENT=client
set CLIENT=.\web\main

:: Generate Protobuf
if "%1" == "protobuf" (
    echo Generate gRPC Interface from Protobuf...
    protoc %PROTO_DIR%/*.proto --go_out=%PROTO_DIST% --go-grpc_out=%PROTO_DIST%
    echo Generate Complete...
    exit /b 0
)

:: Start Server
if "%1" == "server" (
    if "%2" == "" (
        echo Error: Server IP Not Found, Usage: ./run.bat server ^<server ip^> ^<server port^> 
        exit /b 0
    )
    if "%3" == "" (
        echo Error: Server Port Not Found, Error Usage: ./run.bat server ^<server ip^> ^<server port^>  
        exit /b 0
    )
    echo Starting server on %2:%3
    go run %SERVER%.go %2 %3 %4 %5
    exit /b 0
)

:: Start Client
if "%1" == "client" (
    if "%2" == "" (
        echo Error: Server IP Not Found, Usage: ./run.bat client ^<server ip^> ^<server port^> time?
        exit /b 0
    )
    if "%3" == "" (
        echo Error: Server Port Not Found, Error Usage: ./run.bat client ^<server ip^> ^<server port^> time?
        exit /b 0
    )
    
    echo Starting client connect to server %2:%3
    if "%4" == "time" (
        go run %CLIENT%.go %2 %3 %4
    ) else (
        go run %CLIENT%.go %2 %3
    )
    exit /b 0
)

:: Start Web Client
if "%1" == "web" (
    if "%2" == "" (
        echo Error: Server IP Not Found, Usage: ./run.bat client ^<server ip^> ^<server port^>
        exit /b 0
    )
    if "%3" == "" (
        echo Error: Server Port Not Found, Error Usage: ./run.bat client ^<server ip^> ^<server port^>
        exit /b 0
    )
    
    echo Starting client connect to server %2:%3
    go run .\web\main.go %2 %3

    exit /b 0
)

:: Start Testing
if "%1" == "test" (
    if "%2" == "" (
        echo Start Test All Unit Test
        go test %TEST_DIR% -v
    ) else (
        echo Start Test %2 Unit Test
        go test %TEST_DIR% -v -run %2
    )   
    exit /b 0
)

echo Unknown command: %1
echo Usage: run.bat ^[ protobuf ^| server ^| client ^| test ^] ...options
exit /b 1