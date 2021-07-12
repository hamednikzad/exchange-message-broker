FROM ubuntu:20.04
ENV GO111MODULE=on

WORKDIR /app/exchangeMessageBroker
COPY ./build/exchMsgBroker .
ENV PORT 8080
EXPOSE 8080

WORKDIR /app/exchangeMessageBroker
ENTRYPOINT ./exchMsgBroker

#ENV SOURCES /home/hamed/go/src/finobot.ir/exchangeMessageBroker
#COPY . ${SOURCES}
#//RUN cd ${SOURCES}/exchMsgBroker && CGO_ENABLED=0 go install