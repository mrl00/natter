# Create a space and capture the ID
POST http://localhost:8080/spaces
{
  "name": "Coffee Morning",
  "owner": "alice"
}
HTTP 201
[Captures]
space_id: jsonpath "$.id"
[Asserts]
jsonpath "$.name" == "Coffee Morning"
jsonpath "$.owner" == "alice"
jsonpath "$.id" matches "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"
jsonpath "$.createdAt" exists
