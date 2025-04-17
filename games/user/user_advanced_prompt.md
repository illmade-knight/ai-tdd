### Response and JSON

Now if a player does not exist, instead of a 0 response, we want the server to response with a status of 404

we also want a standard JSON response for User so we use the following JSON form for user

```json
{
  "name": "fooname",
  "score": 1
}
```

this JSON response is marshalled from a Golang struct

```go
type User struct {
	Name string
	Score int
}
```

to distinguish between *real* 0 scores and a non-existent User the store needs to return a standardized UserNotFound error

update the tests to reflect this