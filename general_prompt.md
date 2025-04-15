## Prompt context

We are building a multi phase software development process around a collaborative human-AI development cycle.

For each project we first describe the overall aims.

We then break down project into testable libraries and microservices.

1) The libraries hold reusable code that can be shared across the microservices.

2) The microservices isolate parts of the project from each other = each microservice is independently testable.

To develop the actual code we have an iterative 3 step process.

1) We define our requirements for the code
2) Using these requirements we use machine learning prompts to create tests for the code - we DO NOT create the production code at this stage, only the tests
   a) To improve the performance of the prompts we may provide additional hints to the prompt context for any complex operations
   b) We review the tests and improve the prompts if necessary
3) We hand write initial code that passes the tests
4) We provide the machine learning context with our code and ask it to review the code and suggest improvements.
   a) if the machine learning review includes errors or misunderstanding of the problem we will adjust the prompt for the review.

Now we use do a human led review our requirements and the code - if the review finds changes are necessary we document the changes and repeat the process.

For the test process we have some overall guidelines to be used in all test generation.
For example numerical calculations should be tested for exception states or common errors: division by zero, rounding errors etc...
The machine learning system should identify any common mistakes likely to occur both in the tests and in their explanation.
Later when reviewing the code the AI agent should affirm these potential problem areas have been addressed.

To start with we want to create a library with generally reusable code through all our projects:

### Task definition:

1) Create in Golang generic implementations of functional programming tasks such as Map, Reduce, Fold etc...

2) Create tests for a GoLang implementation of Reduce, start with tests for Reduce using Int and a sum function - so we can pass the function a slice of type []int and
   sum function and the returned value is the sum of all ints in the slice

Create those test now and the programmer with then implement part 3, the initial code which passes the test before returning to the prompt 

