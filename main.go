package main

import (
	"fmt"
	"os"
)

func main() {
	myInput := os.Getenv("INPUT_DD_API_KEY")

	output := fmt.Sprintf("Hello %s", myInput)

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
}
