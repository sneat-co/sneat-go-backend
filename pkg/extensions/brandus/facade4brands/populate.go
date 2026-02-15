package facade4brands

import "fmt"

func Populate() {

	for _, maker := range autoMakers {
		fmt.Println(maker.Title)
	}
}
