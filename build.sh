echo Creating build directory
mkdir -p build
cd exchMsgBroker
echo Building...
go build -o ../build/exchangeMsgBroker
echo Building successful
cp ./appsettings.json ../build/
