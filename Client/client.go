package main

import (

        "log"
        "github.com/nats-io/nats"
        "os/exec"
        "io/ioutil"
        "strconv"
        "container/heap"
        "fmt"
)

func main() {
        // Create server connection
        natsConnection, _ := nats.Connect(nats.DefaultURL)
        log.Println("Connected to " + nats.DefaultURL)
        // init the heap
        h := &IntHeap{}
        heap.Init(h)
        i := 1
        // n := strconv.Itoa(i)

        // Subscribe to subject
        natsConnection.Subscribe("Radio-1", func(msg *nats.Msg) {
                log.Println("se recibio audio...")
                fmt.Println("length of heap: ", len(*h) )
                fmt.Println("actual index: ", i)
                content := msg.Data
                i++
                name_song := i //number of the slice put in the heap
                heap.Push(h, name_song)

                err := ioutil.WriteFile("cancion"+strconv.Itoa(name_song)+".mp3",content,0777)
                log.Println("se creo  archivo cancion"+ strconv.Itoa(name_song)+".mp3")
                if err != nil {
                        log.Fatal(err)
                }
        })
        for true {
                fmt.Println("Songs playing")
                for(len(*h) == 0){
                        fmt.Println("waiting for songs")
                }
                song_to_be_played := heap.Pop(h) //pop the next song to be played
                cmd := exec.Command("mplayer", "cancion"+strconv.Itoa(song_to_be_played.(int))+".mp3")
                err := cmd.Start()
                if err != nil {
                        log.Fatal(err)
                }
                log.Printf("Waiting for command to finish...")
                err = cmd.Wait()
        }
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
