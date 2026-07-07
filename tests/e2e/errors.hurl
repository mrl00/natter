# Error cases

# Missing space
POST http://localhost:8080/spaces/nonexistent/messages
{
  "author": "alice",
  "content": "hello"
}
HTTP 404
[Asserts]
jsonpath "$.error" == "space not found"

# List messages in nonexistent space
GET http://localhost:8080/spaces/nonexistent/messages
HTTP 404
[Asserts]
jsonpath "$.error" == "space not found"

# Get message in nonexistent space
GET http://localhost:8080/spaces/nonexistent/messages/abc123
HTTP 404
[Asserts]
jsonpath "$.error" == "space or message not found"

# Create space with missing name
POST http://localhost:8080/spaces
{
  "name": "",
  "owner": "alice"
}
HTTP 400
[Asserts]
jsonpath "$.error" == "name and owner are required"

# Create space with missing owner
POST http://localhost:8080/spaces
{
  "name": "Test",
  "owner": ""
}
HTTP 400
[Asserts]
jsonpath "$.error" == "name and owner are required"

# Add message with missing author
POST http://localhost:8080/spaces/abc/messages
{
  "author": "",
  "content": "hello"
}
HTTP 400
[Asserts]
jsonpath "$.error" == "author and content are required"

# Create space with missing fields (empty JSON)
POST http://localhost:8080/spaces
{}
HTTP 400
[Asserts]
jsonpath "$.error" == "name and owner are required"
