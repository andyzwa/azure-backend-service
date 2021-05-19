#### build -----------------------------------------------------------------------------------------------
#FROM golang:alpine as build
#
#WORKDIR /go/src/app
#
#COPY . .
#
#RUN go build -o myservice.bin

### run -------------------------------------------------------------------------------------------------

# A minimal Docker image based on Alpine Linux with a complete package index and only 5 MB in size!
FROM alpine

# Add Maintainer Info
LABEL maintainer="André Zwahlen"

# Create appuser
RUN adduser -D -g '' appuser

# Set the Current Working Directory inside the container
WORKDIR /app

# expose port
EXPOSE 8000

#COPY --from=build /go/src/app/myservice.bin /app/myservice.bin
COPY ./app/myservice.bin /app/myservice.bin
COPY ./swaggerui /app/swaggerui

# Use an unprivileged user.
USER appuser

# execute command when docker launches / run
ENTRYPOINT ["/app/myservice.bin"]



