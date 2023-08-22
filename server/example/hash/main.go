package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	str := "62947013-97e6-4df5-a5fa-acab83246624/mccormickjeremy_an_advertisement_for_a_game_called_midjourney_i_63af6fc1-72ad-4b4e-9e9a-0936f128d3e1_upscayl_4x_RealESRGAN_General_x4_v3.png"
	filename := path.Base(str)

	fmt.Println(filename)

	gopath := os.Getenv("DEBUG")
	//gopath2 := os.Getenv("PATH")
	fmt.Println(gopath)
	//fmt.Println(gopath2)
}
