package main

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
)

func getHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func getUIDAndHash(s string) (uint64, string, error) {
	r := regexp.MustCompile(`\[WZH_(\d+)_(.*)\]`)

	match := r.FindStringSubmatch(s)
	fmt.Println(match)
	if len(match) != 3 {
		return 0, "", fmt.Errorf("no uid or hash found")
	}

	uidStr := match[1]
	hashStr := match[2]

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return 0, "", err
	}

	return uid, hashStr, nil
}

func main() {
	s := "[WZH_10_62370639-5c49-4855-aca3-fa67fa1823e2]"

	uidStr, hashStr, err := getUIDAndHash(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("uid:", uidStr)
	fmt.Println("hash:", hashStr)

}
