package main

import (
	"fmt"
	"time"

	"github.com/Woodfyn/inMemoryCache/cache/cache"
)

func main() {
	myCache := cache.NewCache(1)
	myCache.Set(1, "qwewqeqwewqe", 10*time.Second)

	valueOne, err := myCache.Get(1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(valueOne)

	time.Sleep(time.Second * 12)

	valueTwo, err := myCache.Get(1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(valueTwo)
}
