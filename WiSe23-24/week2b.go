package main

import "fmt"
import "time"

// Kanal ohne Buffer.
// Entweder A oder B blockiert.
func channel1() {
    var x chan int

    x = make(chan int)

    go func() {
        x <- 1 // A
    }()

    go func() {
        x <- 2 // B
    }()

    y := <-x // C
    fmt.Printf("%d", y)

}

/*
Immer noch chaotisches Verhalten.
Z.B.

    S1 <-> R1  und S2 <-> R2

oder

    S1 <-> R2 und S3 <-> R1
*/
func channel2() {

    x := make(chan int)

    snd := func() { x <- 1 }

    rcv := func() { <-x }

    go snd() // S1
    go snd() // S2
    go snd() // S3
    go rcv() // R1
    go rcv() // R2
    time.Sleep(1 * time.Second)

}

// Deadlock moeglich. Wieso?
func channel3() {
    x := make(chan int)

    snd := func() { x <- 1 }

    rcv := func() { <-x }

    go rcv() // R1


    go snd() // S1


    snd() // S2

}

func main() {

//    channel1()
//    channel2()
    channel3()

}

/*

channel3:

 Falls

     R1 <-> S1

 dann blockiert S2.

 Alle (verbleibenden) Threads sind blockiert.

 =>

 Deadlock.

*/