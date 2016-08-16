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
var currentTrack int = 0
var latestTrack int = 0
var nowPlayingTrack string = ""
var h *TrackHeap

func main() {
    captureInterruptSignal()
    url := getRadioUrl()
    // Create server connection
    natsConnection, _ := nats.Connect(url)
    c, _ := nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER)
    log.Println("Connected to " + url)
    // init the heap
    h = &TrackHeap{}
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

func blockUntilAvailableTrack(h *TrackHeap){
    if(len(*h) == 0){
        fmt.Println("waiting for tracks...")
        for(len(*h) == 0){
        }
    }
}

func playTrack(h *TrackHeap){
    song_to_be_played := heap.Pop(h).(*TrackData.TrackData) //pop the next song to be played
    trackTitle := song_to_be_played.Name
    newTrack := getTrackFile(trackTitle, song_to_be_played.Index)
    updateNowPlayingStatus(trackTitle)
    cmd := exec.Command("mplayer", newTrack)
    err := cmd.Start()
    if err != nil {
            log.Fatal(err)
    }
    //wait until track finishes playback
    err = cmd.Wait()
    rm := exec.Command("rm", newTrack)
    rm.Start()
}

func updateNowPlayingStatus(nowPlaying string){
    if(strings.Compare(nowPlayingTrack, nowPlaying) != 0){
        nowPlayingTrack = nowPlaying
        log.Println("Now playing: " + nowPlayingTrack)
    }
}

func downloadTrack(data *TrackData.TrackData) {
    content := data.File
    trackName := data.Name
    trackToWrite := getTrackFile(trackName, latestTrack)
    err := ioutil.WriteFile(trackToWrite, content, 0644)
    //log.Println("downloaded file: "+ data.Name)
    if err != nil {
        log.Fatal(err)
    }
    data.File = nil
    data.Index = latestTrack
    heap.Push(h, data)
    latestTrack = latestTrack + 1
}

func getTrackFile(name string, index int) string {
    return tracksDir + name +strconv.Itoa(index)+".mp3"
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


// An TrackHeap is a min-heap of ints.
type TrackHeap []*TrackData.TrackData

func (h TrackHeap) Len() int           { return len(h) }
func (h TrackHeap) Less(i, j int) bool { return h[i].Index < h[j].Index }
func (h TrackHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TrackHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *h = append(*h, x.(*TrackData.TrackData))
}

func (h *TrackHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
