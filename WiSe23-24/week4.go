
package main

import "fmt"
import "time"


// Erstes Beispiel

func philo1(id int, forks chan int) {

    for {
        <-forks
        <-forks
        fmt.Printf("%d eats \n", id)
        time.Sleep(1 * 1e9)
        forks <- 1
        forks <- 1

        time.Sleep(1 * 1e9) // think

    }

}

func main1() {
    var forks = make(chan int, 3)
    forks <- 1
    forks <- 1
    forks <- 1
    go philo1(1, forks)     // P1
    go philo1(2, forks)     // P2
    philo1(3, forks)        // P3
}


func philo2(id int, forks chan int) {
    for {
        <-forks
        select {
        case <-forks:
            fmt.Printf("%d eats \n", id)
            time.Sleep(1 * 1e9)
            forks <- 1
            forks <- 1

            time.Sleep(1 * 1e9) // think
        default:
            forks <- 1
        }
    }

}

func main2() {
    var forks = make(chan int, 3)
    forks <- 1
    forks <- 1
    forks <- 1
    go philo2(1, forks)
    go philo2(2, forks)
    philo2(3, forks)
}

func philo3(id int, forks chan int) {
    for {
        <-forks
        select {
        case <-forks:
            fmt.Printf("%d eats \n", id)
            time.Sleep(1 * 1e9)
            forks <- 1
            forks <- 1

            time.Sleep(1 * 1e9) // think
        default:
            // forks <- 1  // (LOC)
        }
    }

}

func main3() {
    var forks = make(chan int, 3)
    forks <- 1
    forks <- 1
    forks <- 1
    go philo3(1, forks)   // P1
    go philo3(2, forks)   // P2
    philo3(3, forks)      // P3
}


func main() {
//	main1()
	//	main2()
	main3()
}
