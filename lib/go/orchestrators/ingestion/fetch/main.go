package main

import (
    "os"
    "fmt"
    "bufio"
)


func main() {
    fmt.Println("Enter current location: ")
    inputReader := bufio.NewReader(os.Stdin)
    cur, _ := inputReader.ReadString('\n')

    fmt.Println("Similar or different (1-5): ")
    var dest string
    fmt.Scanln(&dest)

    fmt.Print("Your source/dest compatibility ratio is: ", cur + " " + dest + "\n")
}
