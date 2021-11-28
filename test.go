package main

import "log"

func main() {
	a := []int{1, 2, 3, 4, 8}

	index := 4

	a = append(a[:index+1], a[index:]...)
	a[index] = 6

	log.Println(a)
}
