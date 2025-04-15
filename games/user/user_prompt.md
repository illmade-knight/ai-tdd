# A Game server

OK we're going to start building a more complex example using our AI driven test development.

This is going to involve microservices and to begin with we are going to concentrate on the User

## The user

we will begin very simply and add a REST api for a User where we show the number of wins a given user has.

the microservice will have 2 routes initially

1) GET /user/[name]/score
2) PUT /user/[name]/score

where [name] is an UID for the player (in future GET /user/[name] should return the full public representation of the given user)

Note that we are ignoring the User create steps at this point.
Our user will initialized with a score of 0, so an initial call to GET /user/score will result in 0 being returned
The score will be an Int

For this we're going to introduce *mocking*

1) GET will access a PlayerStore to get scores for a player. We use an interface for PlayerStore so when we test we can create a simple stub to test we can also have flexibility in future on how our user is stored.

2) For PUT we can spy on its calls to PlayerStore to make sure it stores players correctly. Our implementation of saving won't be coupled to retrieval.

The PlayerStore is accessed by a PlayerServer.

All our microservices follow the same pattern a [foobar]Server struct with a public Start() method which calls a 
private startHttp method defining the paths

so here the tests should first setup a PlayerServer call it's Start() method to initiate the path handlers and then the paths
using mocks