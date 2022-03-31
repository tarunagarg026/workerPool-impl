package main

import (
	"fmt"
	"workerPool/workerPool"
)

func main() {
	fmt.Println("Hello World")
	pool := workerPool.GetPool(10, 2, 1)
	workerPool.InitializeWorkerPool("127.0.0.1:8000", pool)
}
