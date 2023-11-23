// Siehe
// https://sulzmann.github.io/AutonomeSysteme/lec-data-race-examples.html#(4)

package main

import "sync"
import "fmt"
import "time"

// TSan vector clock size is limited to 256
// If there are more active threads,
// we might not be able to detect the race.
func manyThreads() {
    x := 1
    var m sync.RWMutex

    threadNo := 500 // for 200 we'll encounter the race

    protectedWrite := func() {

        time.Sleep(1 * time.Second)
        m.Lock()
        x = 3
        m.Unlock()
        time.Sleep(4 * time.Second)
    }

    for i := 0; i < threadNo; i++ {

        go protectedWrite()

    }
    time.Sleep(3 * time.Second)
    fmt.Printf("%d", x) // unprotected read

    time.Sleep(5 * time.Second)

}

// TSan only records four "most recent" reads.
// If there's an unprotected read followed by many protected reads,
// we might not be able to detect the race.
func manyConcurrentReads() {
    x := 1
    var m sync.RWMutex

    threadNo := 100 // for 1 we'll encounter the race

    f := func(x int) {}

    protectedRead := func() {
        m.Lock()
        f(x)
        m.Unlock()

    }

    // unprotected read
    go func() {
        f(x)
    }()

    for i := 0; i < threadNo; i++ {

        go protectedRead()

    }

    time.Sleep(1 * time.Second)
    m.Lock()
    x = 1 // protected write
    m.Unlock()

}

// Writes don't seem to replace reads.
// We always run into a write-read race here.
func manyConcurrentWrites() {
    x := 1
    var m sync.RWMutex

    threadNo := 100 // for 1 we'll encounter the race

    f := func(x int) {}

    protectedWrite := func() {
        m.Lock()
        x = 3
        m.Unlock()

    }

    // unprotected read
    go func() {
        f(x)
    }()

    for i := 0; i < threadNo; i++ {

        go protectedWrite()

    }

    time.Sleep(1 * time.Second)
    m.Lock()
    x = 1 // protected write
    m.Unlock()

}

// TSan reports two data race (pairs).
func subsequentRW() {
    x := 1
    y := 2

    go func() {
        fmt.Printf("%d", x)
        fmt.Printf("%d", y)

    }()

    x = 3
    y = 3
    time.Sleep(1 * time.Second)

}

// There are two data races: (R1,W) and (R2,W).
// TSan reports (R2,W).
// The later R2 "overwrites" the earlier R1.
func subsequentRW2() {
    x := 1

    go func() {
        fmt.Printf("%d", x)   // R1
        fmt.Printf("%d", x+1) // R2

    }()

    x = 3 // W

    time.Sleep(1 * time.Second)

}

/////////////////////////////
// Playing

/*

There are threadNo read-write races.
TSan only keeps the four most recent reads.
TSan only reports at most one read-write race.

 From a diagnosis point of view:
   - Would it help to report all read-write races?


*/

func concReadsRaceWithWrite() {
    x := 1

    threadNo := 100

    f := func(x int) {}

    unprotectedRead := func() {
        f(x)
        time.Sleep(1 * time.Second)
    }

    for i := 0; i < threadNo; i++ {

        go unprotectedRead()

    }

    time.Sleep(1 * time.Second)
    x = 2 // unprotected write
    fmt.Printf("%d", x)

}

func main() {

    // manyThreads()
    manyConcurrentReads()
    // manyConcurrentWrites()
    // subsequentRW()
    // subsequentRW2()
    // concReadsRaceWithWrite()

}