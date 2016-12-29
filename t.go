package main

import "fmt"

func test(){
	number := 68523
	min := 0
	max := 100000
	for {
		if (min + max) / 2 < number {
			min = (min + max) / 2
		}else if (min + max) / 2 > number{
			max = (min + max) / 2
		}else{
			fmt.Println((min + max) /2, "END")
			break;
		}
		fmt.Println(max, min)
	}
}
