FROM golang:alpine
COPY bin/gok8s .
EXPOSE 8080
CMD ["./gok8s"]