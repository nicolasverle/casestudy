FROM golang:1.22 as gobuild

WORKDIR /home

COPY go.mod go.sum .

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -trimpath  -o linkextractor .

FROM busybox

COPY --from=gobuild /home/linkextractor /

USER 65535:65535

ENTRYPOINT [ "/linkextractor" ]