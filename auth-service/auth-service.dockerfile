FROM alpine:latest

RUN mkdir /app

COPY binary_file/authServiceApp /app
COPY .env /

CMD [ "/app/authServiceApp" ]
