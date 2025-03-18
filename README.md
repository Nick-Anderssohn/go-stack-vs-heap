# golang stack-vs-heap
Small benchmarks to see how much overhead dynamic memory allocation has
in golang.

## Results
As expected, keeping things on the stack is significantly faster. Let's look at each
benchmark one-by-one. All tests were run on my 2017 MacBook Pro.
```
goos: darwin
goarch: amd64
pkg: github.com/Nick-Anderssohn/go-stack-vs-heap/bench
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
```

### 8 byte struct
Here were the results for the 8 byte struct:
```
BenchmarkSmall/stack-8         	280214240	         4.273 ns/op	       0 B/op	       0 allocs/op
BenchmarkSmall/heap-8          	62014394	        17.61 ns/op	       8 B/op	       1 allocs/op
```
Each operation was about **4.12x faster** when using the stack. We can confirm that the
first benchmark was correctly using the stack since there were 0 allocs/op, and the
second benchmark was correctly using the heap since there was 1 allocs/op.

### 1 KiB struct
Here were the results for the 1 KiB struct:
```
BenchmarkMed/stack-8         	29585313	        39.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkMed/heap-8          	 7303219	       154.4 ns/op	    1024 B/op	       1 allocs/op
```
Each operation was **3.93x faster** when using the stack. We can confirm that the
first benchmark was correctly using the stack since there were 0 allocs/op, and the
second benchmark was correctly using the heap since there was 1 allocs/op.

### 4 KiB struct
Here were the results for the 4 KiB struct:
```
BenchmarkLarge/stack-8         	12217219	        94.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkLarge/heap-8          	 1965657	       605.9 ns/op	    4096 B/op	       1 allocs/op
```
Each operation was **6.43x faster** when using the stack. We can confirm that the
first benchmark was correctly using the stack since there were 0 allocs/op, and the
second benchmark was correctly using the heap since there was 1 allocs/op.

### 1 MiB struct
Here were the results for the 1 MiB struct:
```
BenchmarkHuge/stack-8         	    2444	    421463 ns/op	 1048598 B/op	       1 allocs/op
BenchmarkHuge/heap-8          	   12634	     94800 ns/op	 1048585 B/op	       1 allocs/op
```
Something interesting happened! You'll notice both benchmarks had 1 allocs/op. This
means the runtime put the struct on the heap in both benchmarks, even though the
first benchmark returns by value, not by pointer. That makes sense since the starting
stack size of a goroutine is 2 KiB. Growing it to 1 MiB would be wasteful. You'll
also notice that this is the first time where the return-by-pointer code is faster!

## The takeaway
The results show you should keep things on the stack as much as you can. Not only is it
faster, but your program will take up less memory, and the garbage collector won't have to do as much.
One small thing you can do to improve your code, is to default to factory functions that return
by value, instead of by pointer. Personally, I follow the `CreateFoo()` convention
for factory functions that return by value, and use `NewFoo()` if it returns by pointer.
I do this because the `new` keyword indicates dynamic memory allocation in many languages.
In fact, the built-in `new` function in go follows this convention as well.
Note that you should still keep function receivers and parameters as pointers. That way you
don't copy the entire struct along the stack for every function call, and instead just pass
the memory address :)