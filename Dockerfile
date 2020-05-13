FROM golang:alpine AS build
WORKDIR /scratch
COPY . .
RUN go get
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /scratch/go-ising go-ising
CMD ["./go-ising"]
