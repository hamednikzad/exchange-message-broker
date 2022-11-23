echo Creating build directory
mkdir -p build
cd exchMsgBroker
echo Building...
go build -o ../build/exchange-message-broker
echo Building successful
cp ./appsettings.json ../build/
