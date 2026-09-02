package main

import (
	"fmt"
	"geometry/cmd/app/internal/calculations"
	"geometry/cmd/app/internal/validation"
	"os"
)

func main() {
	flags, err := validation.ParseFlags()
	if err != nil {
		fmt.Println("error parsing flags:", err)
		os.Exit(1)
	}
	figures := calculations.CreateFigures(&flags)

	fmt.Println(figures.Calculate(&flags))
}
