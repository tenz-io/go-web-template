# Use base golang image from Docker Hub
FROM golang:1.21 AS build

WORKDIR /webapp

# Avoid dynamic linking of libc, since we are using a different deployment image
# that might have a different version of libc.
ENV CGO_ENABLED=0

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY . ./

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /app ./cmd

# Now create separate deployment image
FROM gcr.io/distroless/static-debian11

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

# Copy template & assets
WORKDIR /webapp
COPY --from=build /app ./app
COPY web web/
COPY config config/

ENTRYPOINT ["./app", "-a", "8085", "-p", "8080", "-vv"]
