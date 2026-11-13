FROM public.ecr.aws/d3j8x8q7/olympus-base-go:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/jstemmer/go-junit-report/v2@latest
ENV PATH="${PATH}:/root/go/bin"

COPY . .

RUN go build -o /app/fsql ./cmd/fsql

CMD ["/bin/bash"]
