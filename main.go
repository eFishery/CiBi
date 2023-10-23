package main

import (
	"fmt"

	"github.com/eFishery/CiBi/cmd"
	figure "github.com/common-nighthawk/go-figure"
)

var (
    buildTime string
    version   string
)

const ciBiLogo = `

	   .'|_.-
         .'  '  /_
      .-"    -.   '>
   .- -. -.    '. /    /|_
  .-.--.-.       ' >  /  /
 (o( o( o )       \_."  <
  '-'-''-'            ) <
(       _.-'-.   ._\.  _\
 '----"/--.__.-) _-  \|
	`

func main() {
	figure.NewFigure("CiBi", "", false).Print()
	fmt.Printf("version %s (%s)", version, buildTime)
	fmt.Println(ciBiLogo)
	fmt.Print("The CI/CD created with love from efisherian\n\n")
	cmd.Execute()
}
