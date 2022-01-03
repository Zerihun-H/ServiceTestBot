package main

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(data ...interface{}) {
	fmt.Println("[")
	for i, d := range data {

		var p []byte
		//    var err := error
		p, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s \n", p)
		if i+1 != len(data) {
			fmt.Println(",")
		}
	}
	fmt.Println("]")
}

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %s", name, elapsed)
// }
