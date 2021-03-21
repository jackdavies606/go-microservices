FROM golang:1.15

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/microservices/item

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY ./Item.go .
COPY ./main.go .
COPY items.csv .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["item"]