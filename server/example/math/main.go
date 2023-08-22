package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 生成一个 19 位的随机数字字符串
	max := int64(math.Pow10(10)) - 1
	randNum := rand.Int63n(max)
	randNum2 := rand.Int63n(int64(math.Pow10(8)))

	randStr := fmt.Sprintf("%d%010d%08d", rand.Intn(8)+1, randNum, randNum2)

	fmt.Println(rand.Intn(9))
	fmt.Println(randStr)
}
