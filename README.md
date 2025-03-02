# ai-tdd
experiment with LLMs to power better test driven development

### why?

I've been enjoying the gitbook [learn go with tests](https://quii.gitbook.io/learn-go-with-tests)

There's lots of good stuff in there (some issues too of course - we wouldn't be programmers if we didn't disagree...)

It's based around test driven development 

now, I've not been a huge fan of test driven development, I find the initial definition of  
the aim/problem within the test a little artificial. I've also found IDE support for this poor:
creation of tests from existing files seems better supported.

anyway, the book got me thinking - could the test development phase be improved by using prompts to an LLM?

### the idea
if we can clearly state our aims and have test code generated for us
1) we're already documenting our process
2) we can reproduce our process for different target languages 
3) we can iterate within our prompt 'documentation' so this stays up to date
4) we get a kind of peer programming experience - the LLMs solution may be unexpected - 
by analysing it we may reinforce our own ideas, or change them

### content
* [basic](basic), try out some basic prompts
* [shapes](shapes), move onto something more complex
* [clock](clock), some basic math, things start to go wrong...

### branches
things were going well but we've encountered some problems - 
we'll start tracking results on different platforms 

so lets start with our initial llm, gemini - we'll record raw outcomes, warts and all
