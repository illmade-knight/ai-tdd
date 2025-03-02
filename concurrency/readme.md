# Concurrency

we'll just cover concurreny to keep up with the gitbook - but I'll probably start skipping sections
as this is primarily to evaluate the use of LLM prompts not an attempt to replicate all the book

````aiprompt
write a golang test for a simple counter, the counter can be called concurently by different threads. 
Do not include an implemenation, just the test
````

gemini responded with the tests in [counter_test](counter_test.go) and also an interface for Counter which 
I use as the base for the implementation

next is a [clock](../clock) where I think I'll start diverging a bit more from my source material