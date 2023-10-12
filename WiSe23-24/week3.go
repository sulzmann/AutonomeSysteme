
package main

import "fmt"
import "time"

func channelInsideChannelExample() {

    // Channel of channels.
    var ch chan (chan string)

    ch = make(chan chan string) // 'chan' ist rechts-assoziativ

    client := func(id string) {
        clientCh := make(chan string)

        // Client sends request.
        ch <- clientCh

	// Waits for acknowledgment.
	// s has type String
        s := <-clientCh
        fmt.Printf("\n Client %s receives %s", id, s)

    }

    go client("A")
    go client("B")

    time.Sleep(1 * time.Second)

    cl := <-ch  // cl is type chan string
    cl <- "Hello"
    time.Sleep(1 * time.Second)

    cl2 := <-ch  // cl is type chan string
    cl2 <- "HaHa"
    time.Sleep(1 * time.Second)
}

func main() {
    channelInsideChannelExample()

}