FROM alpine:latest

RUN mkdir /app

COPY binary_file/paymentServiceApp /app
COPY .env /

CMD [ "/app/paymentServiceApp" ]
