package main

// Dynamische Deadlock Analyse.
// Aufschrieb aus der Vorlesung.
// Details siehe Vorlesungsunterlagen.

import "time"
import "sync"

func test() {
	x := make(chan int)

	rcv := func() {
		<-x
	}

	snd := func() {
		x <- 1
	}

	go snd() // S1
	go rcv() // R1

	// time.Sleep(1 * time.Second)
	// Dadurch kommuniziert S1 <-> R1 und S2 blockiert.
	// => Kommunikationsdeadlock
	// aber im Programmablauf ist dieser nicht ersichtlich
	snd() // S2

	time.Sleep(1 * time.Second)

}

/*

Go Laufzeitsystem erkennt nur wenn ein tatsaechlicher Deadlock vorliegt.

Wir brauchen eine Deadlock Analyse!
Dynamische Analyse!

Vorgehen: Aufzeichnung der Programmspur und deren Analyse.


Viele Dinge sind aus dieser Programmspur.
Z.B. dass S1 "verhungert" falls S2 und R1 kommunizieren.

Zum grossen Teil ist die Beobachtung und Analyse von Go Kommunikationsprimitiven
noch Gegenstand der aktuellen Forschung.

Deshalb betrachten wir hier "nur" Deadlocks in welchen Locks involviert.
Dies wird auch Resourcen Deadlock genannt.


*/

func testMutex() {
	var x sync.Mutex
	var y sync.Mutex

	go func() {
		y.Lock()
		x.Lock()
		// do something
		x.Unlock()
		y.Unlock()
	}()

	x.Lock()
	y.Lock()
	// do something
	y.Unlock()
	x.Unlock()

}

/*

 Wir betrachten den Trace (Programmpfad)

    T1         T2

1.  acq(y)
2.  acq(x)
3.  rel(x)
4.  rel(y)
5.             acq(x)
6.             acq(y)
7.             rel(y)
8.             rel(x)


Gibt es eine Umordnung in welche wir in einen Deadlock laufen?
Ja.

    T1         T2

1.  acq(y)
5.             acq(x)
2.  acq(x) BLOCKED
6.             acq(y)  BLOCKED

weil

T1 hat lock y und will x und
T2 hat lock x und will y.


 Idee: Koennen wir Abhaengigkeiten zwischen "locks" aus dem Trace generieren?

Am Beispiel.
Wir bauen einen Lock Graphen!

 Resourcen = Locks = Knoten im Graph

 Konten y -> x bedeutet ein thread hat y in seinem lockset und holt (aquired) lock x.


    T1         T2                LS(T1)     LS(T2)    Lock graph

1.  acq(y)                       {y}
2.  acq(x)                       {y,x}                 y -> x
3.  rel(x)                       {y}
4.  rel(y)                       {}
5.             acq(x)                       {x}
6.             acq(y)                                 x -> y
7.             rel(y)
8.             rel(x)


Deadlock Warnung falls Zyklus im Lock Graphen!

In unserem Fall.

  y -> x -> y    Zyklus!!!!


 Annahme:

 Lock acquired in Thread T wird auch in thread T released!



*/

// Go Mutexe sind "anders", siehe Vorlesungsunterlagen
// Deshalb liefert die Standard Lock Graphen Konstruktion "zu viele" false positives.
func testMutex2() {
	var x sync.Mutex
	var y sync.Mutex

	go func() {
		// lock und unlock sind nicht synchroniziert,
		// siehe "Go-style mutexes behave like semaphores"
		go func() {
			y.Lock()
			x.Lock()
		}()
		// do something

		go func() {
			x.Unlock()
			y.Unlock()
		}()
		time.Sleep(1 * time.Second)

	}()

	x.Lock()
	y.Lock()
	// do something
	y.Unlock()
	x.Unlock()

}

// Beispiel mit "guard" (aka gate) lock.
func testMutex3() {
	var x sync.Mutex
	var y sync.Mutex
	var z sync.Mutex // guard lock

	go func() {
		z.Lock()
		y.Lock()
		x.Lock()
		// do something
		x.Unlock()
		y.Unlock()
		z.Unlock()
	}()

	z.Lock()
	x.Lock()
	y.Lock()
	// do something
	y.Unlock()
	x.Unlock()
	z.Unlock()

}

// Basierend lock dependencies signalisieren wir hier einen Deadlock.
// Dies ist ein false postive weil
func testMutex4() {
	var x sync.Mutex
	var y sync.Mutex

	// T1
	x.Lock()
	y.Lock()          // Acq_y_T1
	// D1=(T1,y,{x})
	// do something
	y.Unlock()
	x.Unlock()

	// T2
	go func() {
		y.Lock()   // Acq_x_T2
		x.Lock()
		// D2 = (T2,x,{y})
		// do something
		x.Unlock()
		y.Unlock()
	}()

}

/*

Zyklische lock dependencies

D1=(T1,y,{x})

D2=(T2,x,{y})

- Verschiedene Threads
- In D1, y ist in {y} von D2
- In D2, x is int {x} von D1
 - Alle locksets sind disjunkt, dh {x} cap {y} = {}

Welche Information fehlt?

Acq_y_T1 "must happen" immer vor Acq_x_T2

*/

func main() {
	//	test()
	// testMutex()

	testMutex2()
	testMutex3()
}
