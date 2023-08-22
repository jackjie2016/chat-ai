package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "[WZH_10_276b4736-ff6e-4431-a0ee-01c241356cac] A Chinese baby, OC effectbest quality, 8K, bright, super detai ~~ --s 750 --niji 5 - <@1020788970310881290> (31%) (fast)"
	re := regexp.MustCompile(`\((\d+)%\)`) // 匹配括号中的数字
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		fmt.Println(match[1]) // 输出匹配到的数字（百分比）
	}
}
