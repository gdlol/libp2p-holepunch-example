FROM golang:1.18-alpine AS build
WORKDIR /root/project
COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /bin/app

FROM alpine
RUN apk add --no-cache iptables
WORKDIR /bin
COPY --from=build /bin/app ./
CMD [ "./app" ]
