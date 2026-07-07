# Full flow: create space, add messages, list, get

# 1. Create space
POST http://localhost:8080/spaces
{
  "name": "Book Club",
  "owner": "bob"
}
HTTP 201
[Captures]
space_id: jsonpath "$.id"

# 2. Add first message
POST http://localhost:8080/spaces/{{space_id}}/messages
{
  "author": "bob",
  "content": "Hey everyone, what book should we read next?"
}
HTTP 201
[Captures]
message_id: jsonpath "$.id"
[Asserts]
jsonpath "$.author" == "bob"
jsonpath "$.content" == "Hey everyone, what book should we read next?"
jsonpath "$.spaceId" == "{{space_id}}"

# 3. Add second message
POST http://localhost:8080/spaces/{{space_id}}/messages
{
  "author": "carol",
  "content": "How about Dune?"
}
HTTP 201

# 4. List all messages
GET http://localhost:8080/spaces/{{space_id}}/messages
HTTP 200
[Asserts]
jsonpath "$[0].author" == "bob"
jsonpath "$[1].author" == "carol"

# 5. List messages with since filter
GET http://localhost:8080/spaces/{{space_id}}/messages?since=2020-01-01T00:00:00Z
HTTP 200

# 6. Get single message
GET http://localhost:8080/spaces/{{space_id}}/messages/{{message_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == "{{message_id}}"
jsonpath "$.author" == "bob"
jsonpath "$.content" == "Hey everyone, what book should we read next?"
