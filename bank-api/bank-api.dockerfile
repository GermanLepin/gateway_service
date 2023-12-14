FROM alpine:latest

RUN mkdir /app

COPY binary_file/bankAPIApp /app
COPY .env /

CMD [ "/app/bankAPIApp" ]
