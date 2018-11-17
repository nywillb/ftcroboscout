package main

import (
	"fmt"

	"github.com/nywillb/ftcroboscout/ftc"
)

func main() {
	resp, err := ftc.FetchSeasons()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
