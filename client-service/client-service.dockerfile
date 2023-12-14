FROM alpine:latest

RUN mkdir /app

COPY binary_file/clientServiceApp /app
COPY .env /

CMD [ "/app/clientServiceApp" ]
