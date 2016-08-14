package main

import (

        "runtime"
        "log"
        "github.com/nats-io/nats"
        "os/exec"
        "io/ioutil"
)

func main() {
        // Create server connection
        natsConnection, _ := nats.Connect(nats.DefaultURL)
        log.Println("Connected to " + nats.DefaultURL)

        // Subscribe to subject
        natsConnection.Subscribe("Radio-1", func(msg *nats.Msg) {
                log.Println("se recibio audio...")
                content := msg.Data
                err := ioutil.WriteFile("example.mp3",content,0777)
                if err != nil {
                        log.Fatal(err)
                }
        })
        cmd := exec.Command("/bin/mplayer", "test.mp3")
        err := cmd.Start()
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("Waiting for command to finish...")
        err = cmd.Wait()

        runtime.Goexit()
}
