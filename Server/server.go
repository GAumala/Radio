package main

// Import packages
import (
        "log"
        "time"
        "github.com/nats-io/nats"
        "os"
        "io/ioutil"
        "fmt"
        "os/exec"
	"strconv"
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

        //obtener longitud del archivo mp3
        out, err := exec.Command("/bin/mp3info","-p","%S","test.mp3").Output()
        if err != nil {
                log.Fatal(err)
        }

	// obtenemos la longitud en segundos.
	s := string(out[:3])
	length,err := strconv.Atoi(s)
	print(length)
	
        if err != nil {

                log.Fatal(err)
        }

        if(file != nil){
                log.Println("success loading file")
        }

        natsConnection, _ := nats.Connect(nats.DefaultURL)
        defer natsConnection.Close()
        log.Println("Connected to " + nats.DefaultURL)
        for true {
                subject := "Radio-1"

                err := natsConnection.Publish(subject,[]byte("no vale"))
		//err := natsConnection.Publish(subject,file)
		if err!= nil {
			log.Fatal(err)
		}
                fmt.Println("enviando")		
                time.Sleep(time.Duration(length)*time.Second + 1) //sleep time = song's length + 1
        }
}
