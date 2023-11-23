
// Siehe
// https://sulzmann.github.io/AutonomeSysteme/lec-data-race-examples.html#(3)

package main

import "sync"
import "fmt"
import "time"


func detected() {
    // Main = T1
    x := 1
    var m sync.RWMutex

    // T2
    go func() {
        x = 2
        m.Lock()
        m.Unlock()

    }()

    m.Lock()
    x = 3
    m.Unlock()

    fmt.Printf("%d", x)
    time.Sleep(1 * time.Second)
}

/*

       T1       T2

e1.   acq
e2.   w(x)
e3.   rel
e4.            w(x)
e5.            acq
e6.            rel

Wir wissen

  e2 und e4 sind nicht geordnet unter HB.

Was waere wenn

       T1       T2

e1.            w(x)
e2.            acq
e3.            rel
e4.   acq
e5.   w(x)
e6.   rel

Aber dann e1 <HB e5.

Dieser Trace wird von folgendem Programm produziert.


*/

// false negative
// because critical sections are not reordered
func notDetected() {

    x := 1
    var m sync.RWMutex

    go func() {
        x = 2
        m.Lock()
        m.Unlock()

    }()

    time.Sleep(1 * time.Second)

    m.Lock()
    x = 3
    m.Unlock()

    fmt.Printf("%d", x)

}

func main() {
//    detected()
    notDetected()

}