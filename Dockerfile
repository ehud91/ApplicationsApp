FROM golang:1.10

ENV APP_SERVICE_PORT=8123
ENV PG_DBNAME=application_manager
ENV PG_HOST=localhost
ENV PG_PASSWORD=ehud@91
ENV PG_PORT=5432
ENV PG_URL=host=%s port=%d user=%s password=%s dbname=%s sslmode=disable
ENV PG_USER=postgres

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/appManagerApi

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -u github.com/google/uuid
RUN go get -u github.com/gorilla/handlers
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/lib/pq


# Install the package
RUN go install -v github.com/google/uuid
RUN go install -v github.com/gorilla/handlers
RUN go install -v github.com/gorilla/mux
RUN go install -v github.com/lib/pq

# This container exposes port 8123 to the outside world
EXPOSE 8123

# Run the executable
CMD ["go", "run", "main.go"]