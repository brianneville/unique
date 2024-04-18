### unique counter values

package provides functions to help with concurrent
generation of unique, incremental counters.

`Uint64Counter.Get` returns a uint64 counter value.
each call to the function will cause the counter to increment until it rolls 
over to 0 at the unint64 max value.


`BigIntCounter.Get` returns a big int counter value.
each call to the function will cause the counter to increment.
This counter will never roll over, and increases infinitely