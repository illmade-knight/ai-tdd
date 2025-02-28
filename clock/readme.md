# A new idea

we're going to be pursuing a rather quixotic aim - 
generating an svg 'analog' clock with hands for hour, minute, second

quixotic? well are you really going to use golang for svgs and especially ones that update 
regularly?

the new idea to start is using [Juypter](https://jupyter.org/) to sketch out ideas in advance.

we'll use the python kernel (why?: it's good, 
we're trying to avoid copying and pasting our implementation across and thinking instead).

### sketch ideas

we use the [notebook](clocks.ipynb) and then return here with ideas for the prompt

lets try something very general to start and see how we get on

```aiprompt
we have a golang library which generated svg analog clock faces.
Write tests for the underlying implementation that give a time returns points on 
a unit circle for the hour hand, the minute hand and the second hand

(needed a bit more work)
the clocks origin is at 0,0 
```

```
in our implementation we have a helper that is used by hourHandPoint and minuteHandPoint - 
we specify the value and the hand and return the point - 
func unitPoint(t float64, hand Hand) Point, 

(need a bit more work)
the input to unitPoint is not a radian value, instead if it's half six the minute value is just 30

write tests for this
```

this was one of our first big fails - Gemini constantly got the clock face reversed and wrote
tests that said 45 minutes should be at point (1, 0)