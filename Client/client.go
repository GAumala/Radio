package main

import (

        "runtime"
        "log"
        "github.com/nats-io/nats"
        "os/exec"
        "io/ioutil"
        "strconv"
        "container/heap"
      
)

func main() {
        // Create server connection
        natsConnection, _ := nats.Connect(nats.DefaultURL)
        log.Println("Connected to " + nats.DefaultURL)
        i := 1
        n := strconv.Itoa(i)

        // Subscribe to subject
        natsConnection.Subscribe("Radio-1", func(msg *nats.Msg) {
                log.Println("se recibio audio...")
                content := msg.Data
                name_song := &IntHeap{i}
                heap.Init(name_song)

                err := ioutil.WriteFile("cancion"+n+".mp3",content,0777)
                if err != nil {
                        log.Fatal(err)
                }
        })
        cmd := exec.Command("mplayer", "cancion"+n+".mp3")
        err := cmd.Start()
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("Waiting for command to finish...")
        err = cmd.Wait()

        runtime.Goexit()
}
// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
        // Push and Pop use pointer receivers because they modify the slice's length,
        // not just its contents.
        *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
        old := *h
        n := len(old)
        x := old[n-1]
        *h = old[0 : n-1]
        return x
}
