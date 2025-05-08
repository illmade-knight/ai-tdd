## A Game server

OK we're going to start building a more complex example using our AI driven test development.

This is going to involve microservices and to begin with we are going to concentrate on the User

we're going to structure the game server in a directory structure:

/games

the user part is in 

/games/user

we will have reused components in /games/user/lib

a server in 

/games/user/server

and store implementations in 

/games/user/store

we keep test files in the same directory as the files they are testing