FROM golang
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/jamespearly/loggly
RUN go build -o main .
CMD ["./main"]
