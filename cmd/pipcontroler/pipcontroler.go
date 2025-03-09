package main

import (
    "fmt"
    "os"
)

func main() {
    version := "0.1.0 Beta"
    command := map[string]string{
        "pipcontroler": "main command",
        "ipp":          "Scan the code and Install the library",
    }

    args := os.Args
    if len(args) > 1 {
        switch args[1] {
        case "-v", "--version":
            fmt.Printf("PipControler Version %s\n", version)
        case "-h", "--help":
            if len(args) > 2 {
                help_message, found := command[args[2]]
                if found {
                    fmt.Printf("%s: %s\n", args[2], help_message)
                } else {
                    fmt.Printf("%s Not Found\n", args[2])
                }
            } else {
                for com, mes := range command {
                    fmt.Printf("%s: %s\n", com, mes)
                }
            }
        default:
            fmt.Println("Unknown command. Use -h or --help for help.")
        }
    } else {
        fmt.Println("No command provided. Use -h or --help for help.")
    }
}