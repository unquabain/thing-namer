package main

import (
	"fmt"
	"os"
	wb "render_test/wizardbacon"
)

func main() {
	r, err := wb.Fetch()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	fmt.Println(r.Intro)
	fmt.Println(r.ProjectName)
	fmt.Println(r.Outro)

}
