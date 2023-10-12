package main

import "fmt"
import "time"

// Kanaele sind getypt.
func channel1() {
    var x chan int

    x = make(chan int, 1) // Auf Heap allokiert, automatische Speicherverwaltung

    go func() {
        x <- 1 // Senden
    }()

    go func() {
        x <- 2 // Senden
    }()

    y := <-x // Empfangen
    fmt.Printf("%d", y)

    y2 := <-x // Empfangen
    fmt.Printf("%d", y2)

}

// Kanal verhaelt sich wie eine "queue" (FIFO)
func channel2() {

    x := make(chan int, 3)

    x <- 1 // A
    x <- 2 // B

    go func() {
        fmt.Printf("%d", <-x) // Kommuniziert mit A
        fmt.Printf("%d", <-x) // Kommuniziert mit B
    }()

    time.Sleep(1 * time.Second)

}

/*
Chaotisches Verhalten.
Jeder Empfaenger koennte von jedem Sender den Wert erhalten.

Einfluss der Puffergroesse.
Je groesser der Puffer, desto wahrscheinlich kann
der Sender den Wert in den Puffer ablegen.

Puffergroesse hat keinen Einfluss auf das chaotische Verhalten.
*/
func channel3() {

    x := make(chan int, 1)

    snd := func() { x <- 1 }

    rcv := func() { <-x }

    go snd() // A
    go snd() // B
    go snd() // C
    go rcv() // D
    go rcv() // E
    time.Sleep(1 * time.Second)

}

// Falls T auskommentiert wird kommt es zu einem Deadlock. Wieso?
func channel4() {
    x := make(chan int, 2)

    snd := func() { x <- 1 }

    rcv := func() { <-x }

    go func() {
        snd() // A
        rcv() // B
        rcv() // C
    }()

    snd()                              // D
    // time.Sleep(200 * time.Microsecond) // T
    snd()                              // E
    snd()                              // F

    time.Sleep(400 * time.Microsecond)

}

func main() {

//    channel1()
//    channel2()
//    channel3()
    channel4()

}

/*

channel4

 (1)
 Folgende Ausfuehrungsreihenfolge ist moeglich:

     Main Thread    |   Other Thread

1.   D
2.   E
3.                     A
4.   F

     F und A sind blockiert

 => All Threads blockiert

 => Deadlock



 (1)
 Folgende Ausfuehrungsreihenfolge ist moeglich:

     Main Thread    |   Other Thread

1.   D
2.                     A
3.                     B
4.                     C
5.   E
6.   F



 Programm terminiert.

*/