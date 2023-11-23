

// Aus der Vorlesung

/*

Trace A.

     T1          T2

1.   w(x)
2.   acq(y)
3.   rel(y)
4.               acq(y)
5.               w(x)
6.               rel(y)

In obigem Trace steck ein Data Race. Wieso?

Betrachte Trace B.

     T1          T2

2.   acq(y)
3.   rel(y)
4.               acq(y)
1.   w(x)
5.               w(x)
6.               rel(y)

In Trace B steckt ein Data Race.
Zwei writes auf x, direkt nebeneinander.

Aber? Trace B ist keine gueltige Umordnung von Trace A,
weil Programm Ordnung in Thread T1 ist nicht mehr die gleiche.


Betrachte Trace C.

     T1          T2

4.               acq(y)
5.               w(x)
1.   w(x)
6.               rel(y)
2.   acq(y)
3.   rel(y)

In Trace C steckt ein Data Race.
Zwei writes auf x, direkt nebeneinander.

Trace C ist eine gueltige Umordnung von Trace A.
 Alle drei Kriterien sind erfuellt:
 - Programm Ordnung
 - Lock Semantik
 - Last Write

Aber. Die Ordnung zwischen den kritischen Sektionen hat sich geaendert.

In der Lamport HB Relation ist dies nicht moeglich.

 =>

False negatives.


///////////////////////////////////////
LOCKSET

   T1          T2               LS

e1.   acq(y1)
e2.   acq(y2)
e3.   rel(y2)
e4.   w(x)                       {y1}
e5.   rel(y1)
e6.               acq(y2)
e7.               acq(y1)
e8.               rel(y1)
e9                w(x)            {y2}
e10.              rel(y2)


 LS(e4) und LS(e9) sind disjunkt!
 Deren Schnittmenge ist leer.

 => Warnung

Dies aber ein false positive.
 Es gibt keine (gueltige) Umordnung in welcher e4 und e9 nebeneinander sind.


"Einfacheres" Beispiel fuer einen false positive

     T1         T2       LS

 e1. w(x)                {}
 e2. fork(T2)
 e3.           w(x)      {}


 LS(e1) und LS(e3) ist disjunkt!

 => Warnung

Aber sicher false positive, weil wegen "fork" e1 immer vor e3 passiert.

*/