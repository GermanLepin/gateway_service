FROM alpine:latest

RUN mkdir /app

COPY binary_file/authenticationServiceApp /app
COPY .env /

CMD [ "/app/authenticationServiceApp" ]
