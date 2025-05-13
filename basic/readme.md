# Basic prompts

lets start with some of the basic examples from  [learn go with tests](https://quii.gitbook.io/learn-go-with-tests)

### repeat prompt

```aiprompt
Write a test in golang for a function that repeats a character, 
in the function to be tested you should be able to specify how many times 
the character repeats. Only create the test, no implementation.
```

I'm going to try google gemini first.

the code created is in [repeat_test](repeat_test.go)

the generated test has a few test cases and covers 0, -1 repeats non-ascii inputs etc

not bad - except for an odd quirk of gemini where it doesn't include golang's square brackets, slightly odd for a google
developed programming language in a google developed LLM environment...
(no longer the case in Gemini 2.5)


Gemini also provides an explanation of what it's providing:

```aiexplanation
This test function, TestRepeatCharacter, is designed to test a function named RepeatCharacter 
(which you'll need to implement separately). 
It includes several test cases, each with a different character and repeat count, 
and checks if the output of RepeatCharacter matches the expected result. 
This setup allows for comprehensive testing of the RepeatCharacter function by covering various scenarios and 
ensuring it produces the correct output for different inputs.
```

### sum prompt

now go onto a sum of ints method, here's the prompt

```aiprompt
Write a test in golang for a function that repeats a character, 
in the function to be tested you should be able to specify how many times 
the character repeats. Only create the test, no implementation.
```

the code created is in [sum_test](sum_test.go)

as I hoped the test expects the requirement to be implemented using: Sum(ints ...)

result := Sum(tc.input...)

looking promising, we'll move onto more complex [shapes](../shapes)