# Health check
GET http://localhost:8080/health
HTTP 200
[Asserts]
jsonpath "$.status" == "ok"
