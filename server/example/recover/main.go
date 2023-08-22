package main

import (
	"fmt"
	"path"
	"strings"
)

func main() {
	url := "https://cdn.discordapp.com/attachments/1118171402680410114/1125342844295598101/twise2805_WZH_1_dfd2dff6-1926-45b6-b698-feebe24372cd__baf6145f-dc0a-42ab-b470-e720525f864d.png"

	fileName := path.Base(url)
	fileName = strings.Split(fileName, "?")[0]

	fmt.Println(fileName)
}
