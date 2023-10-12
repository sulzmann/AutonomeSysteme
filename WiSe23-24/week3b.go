package main

import "fmt"
import "time"

func selectExample() {
    apples := make(chan int)
    oranges := make(chan int)

    go func() {
        apples <- 1
    }()

    go func() {
        oranges <- 1
    }()

    time.Sleep(1 * time.Second)

    var f func()
//        ^^^^^^ Funktions typ

    f = func() {
        // Nur ein "case" wird getriggert.
        select {
        case <-oranges:
		fmt.Printf("\n oranges")
        case <-apples:
            fmt.Printf("\n apples")
        }
    }

    f()
    f()
}

// Versuch einer Emulation von "select" via Hilfthreads.
// Was ist das Problem der Emulation?
func selectExampleEmulate() {
    apples := make(chan int)
    oranges := make(chan int)

    go func() {
        apples <- 1
    }()

    go func() {
        oranges <- 1
    }()

    f := func() {
        ch := make(chan int)
        // Wait for apples
        go func() {
            <-apples
            fmt.Printf("\n apples")
            ch <- 1
        }()

        // Wait for oranges
        go func() {
	    // Put "organges" to sleep, so that only one case triggers.
	    time.Sleep(1 * time.Second)
            <-oranges
            fmt.Printf("\n oranges")
            ch <- 1
        }()

        <-ch // wait for either apples or oranges or we block
    }

    f()
    fmt.Printf("\n f once")
    f()
    fmt.Printf("\n f twice")

}

func main() {

    // selectExample()
    selectExampleEmulate()
    fmt.Printf("done")
}