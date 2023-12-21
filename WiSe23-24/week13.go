
package main

import "fmt"
import "time"
import "net/http"
import "math/rand"

////////////////////
// Simple futures with generics

// A future, once available, will be transmitted via a channel.
// The Boolean parameter indicates if the (future) computation succeeded or failed.

// Java uses "<...>".
// In Go we write "[...]".

type Comp[T any] struct {
    val    T
    status bool
}


// type Future[T any] chan (T, bool)
// The above is not allowed in Go.
// Unfortunately, tuples are not "first-class" in Go.
type Future[T any] chan Comp[T]

// "Server-side" approach
func future[T any](f func() (T, bool)) Future[T] {
    ch := make(chan Comp[T])

    // Execute f "asynchronously"
    go func() {
	r, s := f()
	// Must provide the type instance "[T]".
	// Currently, Go supports limited type inference.
	v := Comp[T]{r, s}
	// FIXME:
	// Add timeout, or some "terminate" signal.
        for {
            ch <- v
        }
    }()
    return ch

}


// Method definition in Go.
// Methods are functions in Go where
// the "leading" argument is the object.
// The get method is blocking!
func (f Future[T]) get() (T, bool) {
    v := <-f
    return v.val, v.status
}

// Non-blocking!
func (ft Future[T]) onSuccess(cb func(T)) {
	// FIXME: code below can be improved
	// If future is available, can call cb immediately.
	// Otherwise, register cb in some "wait" queue.
	// Saves us some "go-routine".
    go func() {
        v, o := ft.get()
        if o {
            cb(v)
        }
    }()

}

// "_" don't care pattern
// Computation resulted in failure, so we don't care about the value.
func (ft Future[T]) onFailure(cb func()) {
    go func() {
        _, o := ft.get()
        if !o {
		cb()    // called if not successful ("false")
        }
    }()

}

///////////////////////////////
// Adding more functionality

// Pick first available future
func (ft Future[T]) first(ft2 Future[T]) Future[T] {

    return future(func() (T, bool) {

        var v T
        var o bool

        // check for any result to become available
        select {
        case x := <-ft:
            v = x.val
            o = x.status

        case x2 := <-ft2:
            v = x2.val
            o = x2.status

        }

        return v, o
    })
}

// Pick first successful future
func (ft Future[T]) firstSucc(ft2 Future[T]) Future[T] {

    return future(func() (T, bool) {

        var v T
        var o bool

	// Start a race among ft and ft2.
	// Pick the first successful future!
        select {
        case x := <-ft:
            if x.status {
                v = x.val
                o = x.status
            } else {
                v, o = ft2.get()
            }

        case x2 := <-ft2:
            if x2.status {
                v = x2.val
                o = x2.status
            } else {
                v, o = ft.get()
            }

        }

        return v, o
    })
}

// Impose some guard function p.
func (ft Future[T]) when(p func(T) bool) Future[T] {

    return future(func() (T, bool) {
        v, o := ft.get()

	// Short-circuit evaluation.
	// First check if the result is successful (o == true)
        if o && p(v) {
            return v, o
        } else {
            return v, false
        }
    })

}

func (ft Future[T]) then(f func(T) (T, bool)) Future[T] {

    return future(func() (T, bool) {
        v, o := ft.get()
        if o {
            return f(v)
        } else {
            return v, o
        }
    })

}

///////////////////////
// Examples

func getSite(url string) Future[*http.Response] {
    return future(func() (*http.Response, bool) {
        resp, err := http.Get(url)
        if err == nil {
            return resp, true
        }
        return resp, false // ignore err, we only report "false"
    })
}

func printResponse(response *http.Response) {
    fmt.Println(response.Request.URL)
    header := response.Header
    // fmt.Println(header)
    date := header.Get("Date")
    fmt.Println(date)

}

func example1() {

    stern := getSite("http://www.stern.de")

    stern.onSuccess(func(response *http.Response) {
        printResponse(response)

    })

    stern.onFailure(func() {
        fmt.Printf("failure \n")
    })

    fmt.Printf("do something else \n")

    time.Sleep(2 * time.Second)

}

func example2() {

    spiegel := getSite("http://www.spiegel.de")
    stern := getSite("http://www.stern.de")
    welt := getSite("http://www.welt.com")

    req := spiegel.first(stern.first(welt))

    req.onSuccess(func(response *http.Response) {
        printResponse(response)

    })

    req.onFailure(func() {
        fmt.Printf("failure \n")
    })

    fmt.Printf("do something else \n")

    time.Sleep(2 * time.Second)

}

// Holiday booking
func example3() {

    // Book some Hotel. Report price (int) and some potential failure (bool).
    booking := func() (int, bool) {
        time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
        return rand.Intn(50), true
    }

    f1 := future[int](booking)

    // Another booking request.
    f2 := future[int](booking)

    f3 := f1.firstSucc(f2)

    f3.onSuccess(func(quote int) {

        fmt.Printf("\n Hotel asks for %d Euros", quote)
    })

    time.Sleep(2 * time.Second)
}

// Flight booking
func example4() {

    rnd := func() int {
        return rand.Intn(300) + 500
    }

    flightLH := func() (int, bool) {
        time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
        return rnd(), true
    }

    flightTH := func() (int, bool) {
        time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
        return rnd(), true
    }

    f1 := future[int](flightLH) // Flight Lufthansa
    f2 := future[int](flightTH) // Flight Thai Airways

    // 1. Check with Lufthansa and Thai Airways.
    // 2. Set some ticket limit

    // Looks nice thanks to method chaining.
		// Can remove (...)
    // f3 := (f1.firstSucc(f2)).when(func(x int) bool { return x < 800 })
    f3 := f1.firstSucc(f2).when(func(x int) bool { return x < 800 })



    f3.onSuccess(func(overall int) {

        fmt.Printf("\n Flight %d Euros", overall)
    })

    f3.onFailure(func() {

        fmt.Printf("\n Booking failed")
    })

    time.Sleep(2 * time.Second)

}

// Composition of several "future" operations: Flight+Hotel booking
func example5() {

    rnd := func() int {
        return rand.Intn(300) + 500
    }

    flightLH := func() (int, bool) {
        time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
        return rnd(), true
    }

    flightTH := func() (int, bool) {
        time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
        return rnd(), true
    }

    stopOverHotel := func() int {
        return 50
    }

    f1 := future[int](flightLH) // Flight Lufthansa
    f2 := future[int](flightTH) // Flight Thai Airways

    // 1. Check with Lufthansa and Thai Airways.
    // 2. Set some ticket limit
    // 3. If okay proceed booking some stop-over Hotel.
    f3 := f1.firstSucc(f2).when(func(x int) bool { return x < 800 }).then(func(flight int) (int, bool) {
        hotel := stopOverHotel()
        return flight + hotel, true  // flight + hotel is the combined price
    })

    f3.onSuccess(func(overall int) {

        fmt.Printf("\n Flight+Stop-over Hotel %d Euros", overall)
    })

    f3.onFailure(func() {

        fmt.Printf("\n booking failed")
    })

    time.Sleep(2 * time.Second)

}

func main() {

    // example1()

    // example2()

    // example3()

    // example4()

    // example5()
}