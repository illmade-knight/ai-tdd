we have a golang library which generated svg analog clock faces.
Write tests for the underlying implementation that give a time returns points on
a unit circle for the hour hand, the minute hand and the second hand

the clock's unit circle origin is at 0,0 

in our implementation we have a helper that is used by all the hands
we specify the value and the hand and return the point, the hand value is the amount of divisions 
if the clock face i.e. minute is 60, hour 12 etc.

func unitPoint(t float64, hand Hand) Point,

the input to unitPoint is not a radian value, instead if it's half six the minute value is just 30

write tests for this