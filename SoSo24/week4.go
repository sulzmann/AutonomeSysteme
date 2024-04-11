package main

import "fmt"
import "time"

func chooseNaive() {
	booking := make(chan int)
	agoda := make(chan int)

	// Modellierung der Umgebung
	go func() {
		time.Sleep(3 * time.Second) // "down time"
		booking <- 1
	}()

	go func() {
		agoda <- 1
	}()

	// Ziel. Erhalte Info via booking.com oder agoda.com
	// Die erste Anfrage reicht.

	// Versuch 1. Fixe Reihenfolge.
	// Problem: Was ist falls der booking.com server "down" ist
	// Umkehrung der Reihenfolge fuehrt zum gleichen Problem.

	x := <-booking
	y := <-agoda
	fmt.Printf("%d %d", x, y)
}

func chooseMoreClever() {
	booking := make(chan int)
	agoda := make(chan int)

	// Modellierung der Umgebung
	go func() {
		time.Sleep(1 * time.Second) // "down time"
		booking <- 1
	}()

	go func() {
		time.Sleep(1 * time.Second) // "down time"
		agoda <- 2
	}()

	// Ziel. Erhalte Info via booking.com oder agoda.com
	// Die erste Anfrage reicht.

	// Beachte:
	//  Wir koennen nicht "testen", ob
	// "send" und "receive" blockiert!!!!

	// Idee
	// Nebenlaeufig warten wir auf "booking" und "agoda"

	waitForBookingOrAgoda := func() {
		ch := make(chan int)

		// Warte bis Nachricht verfuegbar via "booking" und
		// leite diese weiter via "ch"
		// Dieser Teil findet "asynchron" statt.
		// H1.
		go func() {
			ch <- (<-booking)

		}()

		// H2.
		go func() {
			ch <- (<-agoda)

		}()
		x := <-ch
		fmt.Printf("%d", x)
	}

	waitForBookingOrAgoda()
	// Kann passieren, der lokale Kanael "ch" erhalte beide Nachrichten
	// von "booking" und "agoda"
	// D.h. eine dieser Hilfsroutinen, z.B. H1 bleibt stecken.

	// Der folgende Aufruf blockiert, weil
	// beide Nachrichten (in "booking" und "agoda") wurden
	// schon konsumiert.
	waitForBookingOrAgoda()

	// Was tun?
	// Es wird kompliziert.
	// Z.B. wir koennten versuchen, weiter via "ch" zu empfangen, und
	// diese Werte an den Absender zurueckschicken.

}

func chooseReallyClever() {
	booking := make(chan int)
	agoda := make(chan int)

	// Modellierung der Umgebung
	go func() {
		time.Sleep(10 * time.Second) // "down time"
		booking <- 1
	}()

	go func() {
		time.Sleep(15 * time.Second) // "down time"
		agoda <- 2
	}()

	// Ziel. Erhalte Info via booking.com oder agoda.com
	// Die erste Anfrage reicht.

	// Beachte:
	//  Wir koennen nicht "testen", ob
	// "send" und "receive" blockiert!!!!

	// Idee
	// Nebenlaeufig warten wir auf "booking" und "agoda"

	waitForBookingOrAgoda := func() {

		// Magie.
		// Falls beide "cases" moeglich sind,
		// wir nur ein "case" ausgefuehrt.
		// Die Auswahl ist zufaellig.
		// Falls beide "case" nicht moeglich sind,
		// blockiert das "select" statment.
		select {

		case x := <-booking:
			fmt.Printf("Hooray booking %d", x)
		case x := <-agoda:
			fmt.Printf("Hooray agoda %d", x)
		}

	}

	waitForBookingOrAgoda()

	waitForBookingOrAgoda()

}

func main() {

	// chooseNaive()

	// chooseMoreClever()

	chooseReallyClever()
}
