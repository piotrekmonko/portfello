FROM golang:1.21.1-alpine AS base
RUN apk add git
WORKDIR /usr/src/portfello
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS build
WORKDIR /usr/src/portfello
COPY . .
RUN GIT_COMMIT=$(git rev-list --abbrev-commit --abbrev=4 -1 HEAD) && time go build -ldflags "-X github.com/piotrekmonko/portfello/cmd.buildNumber=$GIT_COMMIT" -v -o /bin/portfello main.go
RUN /bin/portfello --version

FROM base AS test
WORKDIR /usr/src/portfello
COPY . .
RUN time go test ./... -v

FROM golang:1.21.1-alpine AS binary
COPY --from=build /bin/portfello /bin/portfello
EXPOSE 8080
ENTRYPOINT ["/bin/portfello"]
