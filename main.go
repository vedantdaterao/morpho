package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/morpho/bencoding"
)

func main() {
    data, _ := os.ReadFile("another.torrent")
    result, _ := bencoding.Decode(data)

    // JSON -
    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
	    fmt.Println("Error marshaling to JSON:", err)
	    return
    }
    fmt.Println(string(jsonData))

    

}
