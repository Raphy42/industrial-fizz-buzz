FROM denoland/deno:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY ./client .
CMD deno task run
