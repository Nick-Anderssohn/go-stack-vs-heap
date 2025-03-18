package bench_test

import "testing"

const KiB = 1024
const MiB = 1048576

type small struct {
	foo [8]byte
}

func createSmall() small {
	return small{}
}

func newSmall() *small {
	return &small{}
}

func (s *small) increment() {
	s.foo[0] = (s.foo[0] + byte(1)) % byte(254)
}

// bigger than the small struct, but still not big enough
// to cause go to resize the stack yet (each goroutine gets a stack size of 2 KiB to start)
type med struct {
	foo [KiB]byte
}

func createMed() med {
	return med{}
}

func newMed() *med {
	return &med{}
}

func (m *med) increment() {
	m.foo[0] = (m.foo[0] + byte(1)) % byte(254)
}

// this one is large enough to cause the runtime to resize the stack
type large struct {
	foo [4 * KiB]byte
}

func createLarge() large {
	return large{}
}

func newLarge() *large {
	return &large{}
}

func (l *large) increment() {
	l.foo[0] = (l.foo[0] + byte(1)) % byte(254)
}

// this one is so huge, that go ends up leveraging the heap even when
// passing by value on my computer.
type huge struct {
	foo [MiB]byte
}

func createHuge() huge {
	return huge{}
}

func newHuge() *huge {
	return &huge{}
}

func (h *huge) increment() {
	h.foo[0] = (h.foo[0] + byte(1)) % byte(254)
}

// these benchmarks each use separate goroutines
// so that they each get a separate stack. Some of
// the struct sizes in these benchmarks are large
// enough to cause the runtime to resize the goroutine
// stack, so it's best to keep them separate rather than
// running them on the main goroutine.

func BenchmarkSmall(b *testing.B) {
	c := make(chan bool)

	b.Run("stack", func(b *testing.B) {
		go func() {
			for b.Loop() {
				s := createSmall()
				s.increment()
			}
			c <- true
		}()
		<-c
	})

	b.Run("heap", func(b *testing.B) {
		go func() {
			for b.Loop() {
				s := newSmall()
				s.increment()
			}
			c <- true
		}()
		<-c
	})
}

func BenchmarkMed(b *testing.B) {
	c := make(chan bool)

	b.Run("stack", func(b *testing.B) {
		go func() {
			for b.Loop() {
				m := createMed()
				m.increment()
			}
			c <- true
		}()
		<-c
	})

	b.Run("heap", func(b *testing.B) {
		go func() {
			for b.Loop() {
				m := newMed()
				m.increment()
			}
			c <- true
		}()
		<-c
	})
}

func BenchmarkLarge(b *testing.B) {
	c := make(chan bool)

	b.Run("stack", func(b *testing.B) {
		go func() {
			for b.Loop() {
				l := createLarge()
				l.increment()
			}
			c <- true
		}()
		<-c
	})

	b.Run("heap", func(b *testing.B) {
		go func() {
			for b.Loop() {
				l := newLarge()
				l.increment()
			}
			c <- true
		}()
		<-c
	})
}

func BenchmarkHuge(b *testing.B) {
	c := make(chan bool)

	b.Run("stack", func(b *testing.B) {
		go func() {
			for b.Loop() {
				h := createHuge()
				h.increment()
			}
			c <- true
		}()
		<-c
	})

	b.Run("heap", func(b *testing.B) {
		go func() {
			for b.Loop() {
				h := newHuge()
				h.increment()
			}
			c <- true
		}()
		<-c
	})
}
