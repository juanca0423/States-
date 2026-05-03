FROM golang:alpine
WORKDIR /app
RUN apk add --no-cache gcc musl-dev git
ENV GOSUMDB=off
ENV GOPROXY=https://proxy.golang.org,direct
RUN go install github.com/air-verse/air@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]
