package main

import (
    "log"
    "github.com/nats-io/nats"
    "os"
    "os/exec"
    "Radio/TrackData"
    "os/signal"
    "io/ioutil"
    "strconv"
    "strings"
    "syscall"
    "container/heap"
    "fmt"
)

const tracksDir string = "tracks/"
var trackIndex int = 0
var h *IntHeap

func main() {
    captureInterruptSignal()
    url := getRadioUrl()
    // Create server connection
    natsConnection, _ := nats.Connect(url)
    c, _ := nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER)
    log.Println("Connected to " + url)
    // init the heap
    h = &IntHeap{}
    heap.Init(h)
    // create directory for downloaded tracks
    os.Mkdir(tracksDir, 0744)

    // Subscribe to subject
    c.Subscribe("Radio-1", downloadTrack)

    for true {
        blockUntilAvailableTrack(h)
        playTrack(h)
    }
}

func getRadioUrl() string{
    var url string = nats.DefaultURL
    if(len(os.Args) > 1){
        url = strings.Replace(url, "localhost", os.Args[1], 1)
    }
    return url
}

func blockUntilAvailableTrack(h *IntHeap){
    if(len(*h) == 0){
        fmt.Println("waiting for tracks...")
        for(len(*h) == 0){
        }
    }
}

func playTrack(h *IntHeap){
    song_to_be_played := heap.Pop(h) //pop the next song to be played
    newTrack := tracksDir + "cancion" + strconv.Itoa(song_to_be_played.(int))+".mp3"
    log.Println("playing track: " + newTrack)
    cmd_wrap := exec.Command("mp3wrap","-v","album.mp3", newTrack)
    cmd_wrap.Start()
    cmd := exec.Command("mplayer", newTrack)
    err := cmd.Start()
    if err != nil {
            log.Fatal(err)
    }
    //wait until track finishes playback
    err = cmd.Wait()
}

func downloadTrack(data *TrackData.TrackData) {
    content := data.File
    trackIndex++ //updateIndex
    trackToWrite := tracksDir + "cancion"+strconv.Itoa(trackIndex)+".mp3"
    err := ioutil.WriteFile(trackToWrite, content, 0644)
    log.Println("downloaded file: "+ trackToWrite)
    if err != nil {
        log.Fatal(err)
    }
    heap.Push(h, trackIndex)
}

func cleanup(){
    fmt.Println("Deleting downloaded tracks...")
    os.RemoveAll(tracksDir)
}
func captureInterruptSignal(){
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

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
