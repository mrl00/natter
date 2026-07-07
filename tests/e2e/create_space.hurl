# Create a space and capture the URI
POST http://localhost:8080/spaces
{
  "name": "Coffee Morning",
  "owner": "alice"
}
HTTP 201
[Asserts]
jsonpath "$.name" == "Coffee Morning"
jsonpath "$.uri" matches "/spaces/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"
