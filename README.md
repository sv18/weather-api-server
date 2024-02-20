# Open Weather API Coding challenge

**Technical Stack**

Language: Golang

API: RESTful

## How to Run

```
go run main.go
```

## how to test

There is one endpoint in this service we can test it using postman or cURL 

1) Get handler postman URL path

```
http://localhost:8080/weather?lat=24.761681&long=-81.191788&appid=<API KEY>

```

2) Get handler cURL

```
curl --location 'http://localhost:8080/weather?lat=24.761681&long=-81.191788&appid=<API KEY>'
```




