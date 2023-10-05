% Die Programmiersprache Go - Kurz und knapp
% Martin Sulzmann

# Übersicht

* Grundlagen
* Multi-threading
* Kanal-basierte Kommunikation
    * Mit und ohne Puffer
    * Kanäle sind Werte
    * Nicht-deterministische Auswahl

# Grundlagen

* Packages
    * Qualifizierter Zugriff
    * Export falls "Grossbuchstabe"

* Deklarationen
    * Erst Name dann Typ
    * Return Type am Schluss

* Layout
    * Ein Statement pro Zeile
    * Falls mehrere Statements, Semikolon notwendig

* Einfache Typinferenz
    * `:=`

* Annonyme Funktionen ("lambdas")


~~~{.go}
package main

import "fmt"

func inc(x int) int {
	x = x + 1
	return x

}

func main() {
	var x int

	y := inc(x)
	fmt.Printf("%d", y)

	f := func(x int) int { return x + 1 }

	fmt.Printf("%d", f(y))

}
~~~~~~~~~~~

# Multi-threading

* Go-Routinen ("light-weight threads")

* "preemptive scheduling"

*  Go-Routinen dürfen dynamisch angelegt werden

* Falls "main" Thread terminiert, terminieren alle Threads (Verhalten verschieden von Java)

~~~{.go}
package main

import "fmt"
import "time"

func thread(s string) {
    for {
        fmt.Print(s)
        time.Sleep(1 * time.Second)
    }
}

func main() {

    go thread("A")
    go thread("B")
    thread("C")
}
~~~~~~~~

# Kanal mit Puffer

* Verhalten ähnlich with Semaphore

* Puffergröße wird bei Definition festgelegt

* Werte/Nachrichten können über einen Kanal übertragen werden ("indirekter Datenaustausch")

* Senden

    * Nicht-blockierend falls noch Platz im Puffer

* Empfangen

    * Wartet bis Element im Puffer verfügbar

    * Falls Element im Puffer, wird dieses Element entfernt

* Jeder Thread darf senden/empfangen und beliebig viele Kanäle anlegen

~~~~{.go}
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

}

// Kanal verhaelt sich wie eine "queue" (FIFO)
func channel2() {

	x := make(chan int, 2)

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
	time.Sleep(200 * time.Microsecond) // T
	snd()                              // E
	snd()                              // F

	time.Sleep(400 * time.Microsecond)

}

func main() {

	channel1()
	channel2()
	channel3()
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
~~~~~~~~~~


# Kanal ohne Puffer

* Sender synchronisiert sich mit Empfänger ("direkter Datenaustausch")

* Kanäle mit Puffer sind gleichmächtig wie ohne Puffer. Können gegenseitig emuliert werden. Siehe Unterlagen.

~~~{.go}
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

	go func() {
		rcv() // R1
	}()

	go func() {
		snd() // S1
	}()

	snd() // S2

}

func main() {

	channel1()
	channel2()
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
~~~~~~~~

# Kanäle sind Werte ("channel of channels")

* Kanal ist ein Wert

* Werte können überall vorkommen, z.B. als Nachricht in einem Kanal

~~~~{.go}
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
		s := <-clientCh
		fmt.Printf("\n Client %s receives %s", id, s)

	}

	go client("A")
	go client("B")

	time.Sleep(1 * time.Second)

	cl := <-ch
	cl <- "Hello"
	time.Sleep(1 * time.Second)
}

func main() {
	channelInsideChannelExample()

}
~~~~~~~~

# Nicht-deterministische Auswahl ("select")

Warte entweder auf

* Nachricht A, oder

* Nachricht B.


Mit dem bisher bekannten, folgender Versuch:

~~~{.go}
<- chA
...
~~~~~~~~

oder

~~~~{.go}
<- chB
...
~~~~~~~~

Was ist falls `<- chA` blockiert? Was ist falls `<- chB` blockiert?


In Go ist folgendes möglich:

~~~~{.go}
select {
  case <- chA: ...
  case <- chB: ...

}
~~~~~~~


* Falls keine Nachricht (via `chA` und `chB`) vorhanden blockiert "select"

* Falls Nachricht vorhanden via `chA` und `chB` wird ein "case" zufällig ausgewählt


~~~{.go}
package main

import "fmt"

func selectExample() {
	apples := make(chan int)
	oranges := make(chan int)

	go func() {
		apples <- 1
	}()

	go func() {
		oranges <- 1
	}()

	f := func() {
		// Nur ein "case" wird getriggert.
		select {
		case <-apples:
			fmt.Printf("\n apples")
		case <-oranges:
			fmt.Printf("\n oranges")
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

	selectExample()
	// selectExampleEmulate()
	fmt.Printf("done")
}

// Problem: selectExampleEmulate()

/*
 Beobachtung:

In der Regel erhalten wir folgende Ausgabe

oranges
 f once
 applesfatal error: all goroutines are asleep - deadlock!


Was passiert?

1. Beim ersten Aufruf von f werden zwei Hilfsthreads
   gestartet welche auf "apples" und "oranges" empfangen.

   Auf Grund der Aussage scheint es so zu sein,
   dass wir "oranges" empfangen.

2. Beim zweiten Aufruf von f werden wieder
   zwei Hilfsthreads gestartet.

   Der eine Hilfsthreads vom ersten Aufruf von f ist
   aber immer noch aktiv und kann "apples" empfangen.

   Deshalb blockieren dann beide Hilfsthreads des
   zweiten Aufrufs von F.
   Wir geraten in einen Deadlock!


*/

/*

Weiteres Beispiel mit "Challenge".

select {
  case x = <-ch1: ...
  case y = <-ch2: ...
  case ch3 <- 1:
  // default and timeout possible
}

 Koennen wir nicht select emulieren ("nachbauen")?

 Idee.
 Fuer jeden "case" einen "Hilfsthread".
 Falls Ereignis eintrifft, Hilfsthread sendet "notify".

 Skizze in Go.

 notify := make(chan int)

 // T1
 go func() {
      <-ch1       // Teste Ereignis
      notify <- 1  // Sende notify
 }()

 // T2
 go func() {
      <-ch2       // Teste Ereignis
      notify <- 1  // Sende notify
 }()

 // T3
 go func() {
      ch3 <- 1    // Teste Ereignis
      notify <- 1  // Sende notify
 }()

 <-notify     // Warte bis eins der Ereignisse eintrifft.

 select fuehrt immer nur ein Ereignis aus,
 auch wenn alle 3 Ereignisse verfuegbar sind.
 Z.B. Werte sind vorhanden auf ch1 und ch2.
 Deshalb wird entweder <-ch1 oder <-ch2 ausgefuehrt.
 Z.B. wir fuehren <-ch1 aus, d.h. erhalten Wert von
 ch1, aber Wert in ch2 bleibt erhalten.

 Problem der Emulation.

 Betrachte, Werte sind vorhanden auf ch1 und ch2.
 T1 und T2 senden notify.
 D.h. Wert aus ch1 und Wert aus ch2 wird empfangen!
 D.h. zwei "cases" werden ausgefuehrt.
 Dies entspricht nicht der Semantik von select.

*/
~~~~~~~~~~~~


# Nicht-deterministische Auswahl ("select") - "timeout" und "default"


* "select" blockiert falls kein Fall verfügbar ist

* Mittels eines "timeout" oder "default" Fall kann dies verhindert werden


~~~~{.go}
package main

import "fmt"
import "time"

// 1. Select blocks if none of the cases is available.
// 2. We can include a timeout to guarantee that select unblocks after a certain time.
func selectWithTimeout() {
	apples := make(chan int)
	oranges := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		apples <- 1
	}()

	go func() {
		time.Sleep(1 * time.Second)
		oranges <- 1
	}()

	select {
	case <-apples:
		fmt.Printf("\n apples")
	case <-oranges:
		fmt.Printf("\n oranges")
	case <-time.After(1 * time.Second):
		fmt.Printf("\n Nothing")
	}

}

// 1. Select with default case never blocks.
func selectWithDefault() {
	apples := make(chan int)
	oranges := make(chan int)

	go func() {
		apples <- 1
	}()

	go func() {
		oranges <- 1
	}()

	select {
	case <-apples:
		fmt.Printf("\n apples")
	case <-oranges:
		fmt.Printf("\n oranges")
	default:
		fmt.Printf("\n Nothing")
	}
}

func main() {

	selectWithTimeout()
	selectWithDefault()
	fmt.Printf("done")
}
~~~~~~~~~~~
