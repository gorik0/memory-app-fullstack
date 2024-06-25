package main

import "log"

type Rio struct {
	name string
}

func main() {
	rio1 := Rio{"egor"}
	rio2 := Rio{"bugor"}
	println(&rio1, "p to p")
	log.Println(rio1)
	rio1 = rio2
	log.Println(rio1)
	println(&rio1, "p to p")
	println(&rio2, "p to p")

}

func copys(rio1 *Rio, rio2 *Rio) {
	println(&rio1, "p to p")
	println(rio1)
	rio1 = rio2
	println(rio1)
	println(&rio1, "p to p")
}
