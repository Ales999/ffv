package main

import "fmt"

type FreeRange struct {
	start int
	stop  int
}

// Конструктор - создаем объект
func NewFreeRange(startid int, stopid int) *FreeRange {
	return &FreeRange{start: startid, stop: stopid}
}

// Распечатать по простому диапазоны.
func (a *FreeRange) PrintData() {
	// Нам не нужны нулевые значение
	if (a.start == a.stop) && ((a.start + a.stop) == 0) {
		return
	}
	// Если диапазон из одного числа состоит
	if a.start == a.stop {
		fmt.Println(a.start)
	} else {
		fmt.Printf("%d - %d\n", a.start, a.stop)
	}
}
