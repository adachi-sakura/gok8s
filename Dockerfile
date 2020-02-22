FROM golang:alpine
COPY gok8s .
EXPOSE 8080
CMD ["./gok8s"]