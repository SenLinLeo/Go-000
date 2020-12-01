package main

import (
	"Go-000/Week02/service"
	"fmt"
)

func main() {
	err := service.Biz()
	fmt.Printf("service: %+v\n", err)
}
