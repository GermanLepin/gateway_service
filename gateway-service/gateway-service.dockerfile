FROM alpine:latest

RUN mkdir /app

COPY binary_file/gatewayServiceApp /app
COPY .env /

CMD [ "/app/gatewayServiceApp" ]
