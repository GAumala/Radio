package main

// Import packages
import (
        "log"
        "time"
        "github.com/nats-io/nats"
        "os"
	"io/ioutil"
	"fmt"
)

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Escribir el nombre del archivo mp3.")
                return
        }
        //obteniendo el nombre del archivo.
        fileName := os.Args[1]
	// obtenemos el array de bytes.
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if(file != nil){
		log.Println("success")
	}
        natsConnection, _ := nats.Connect(nats.DefaultURL)
        defer natsConnection.Close()
        log.Println("Connected to " + nats.DefaultURL)
        for true {
                subject := "Radio-1"
                natsConnection.Publish(subject,[]byte("test"))
		fmt.Println("enviando")		
		time.Sleep(10000*time.Millisecond)
        }
}
