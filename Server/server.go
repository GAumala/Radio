package main

// Import packages
import (
    "log"
    "time"
    "github.com/nats-io/nats"
    "io/ioutil"
    "fmt"
    "os/exec"
    "bytes"
    "os"
    // "strconv"
    dll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

const splitsDir string = "Server/Songs_to_Send/Splits/"
const songsDir string = "Server/Songs_to_Send/"

func main() {
    // create directory for splits
    list_songs_splitted := getfiles_fromdir(splitsDir)
    list_songs := getfiles_fromdir(songsDir) /* Check what files have to be played  */
    fmt.Println("Number of songs to play: ", list_songs.Size())
    fmt.Println("Number of songs splitted: ", list_songs_splitted.Size())
    createMissingSplits(list_songs, list_songs_splitted)

    list_to_send := getListOfSplits(list_songs)

    it := list_to_send.Iterator()
    subject := "Radio-1"
    for true {

    	if(!it.Next()){
    		it.First()
    	}

        file, err := ioutil.ReadFile(it.Value().(string))
        if err != nil {
                log.Fatal(err)
        }

        if(file != nil){
            log.Println("success loading file")
        }

        natsConnection, _ := nats.Connect(nats.DefaultURL)
        defer natsConnection.Close()
        log.Println("Connected to " + nats.DefaultURL)

        err_nats := natsConnection.Publish(subject,file)
        if err_nats!= nil {
            log.Fatal(err)
        }

        fmt.Println("enviando archivo "+ it.Value().(string))
        time.Sleep(time.Duration(10)*time.Second + 1) //sleep time = song's length + 1 (10 seconds)

    }
}

/* function that splits the mp3 file and returns the directory where the files are splitted. */
func splitmp3(fileName string){
    os.Chdir("./Server/Songs_to_Send/Splits")
    // dir,_:= os.Getwd()
    cmd := exec.Command("/bin/mp3splt","-S","10","-d",fileName[0:len(fileName)-4],"../"+fileName)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
            fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
            return
    }
    os.Chdir("../../../") // trace back to root directory
}

/* function that lists files(mp3 files) in a directory and adds them to a list*/
func getfiles_fromdir(dir string) (*dll.List){
    list := dll.New()
    files, err := ioutil.ReadDir(dir)
    if err != nil {
            log.Fatal(err)
    }
    for i :=0; i<len(files); i++ {
            if files[i].Name() != "Splits"{
                    list.Add(files[i].Name())
            }

    }
    return list
}

/*function: check if file is splitted already
listFiles -> list that has the filenames in the directory Songs_to_Send/Splits/
filename -> name of the file that has to be found e.g. test.mp3*/
func isfileSplitted(listFiles *dll.List,fileName string)(bool){
    for i :=0; i<listFiles.Size(); i++ {
        temp_fileName, _ := listFiles.Get(i)
        if temp_fileName.(string) == fileName[0:len(fileName)-4]{
                return true
        }
    }
    return false
}

/*   */
func createMissingSplits(list_songs *dll.List, list_songs_splitted *dll.List){
    for i := 0; i<list_songs.Size(); i++ {
        song, _:= list_songs.Get(i)
        if isfileSplitted(list_songs, song.(string)) == false {
            splitmp3(song.(string)) // create the new split folder
        }

    }
}

func getListOfSplits(list_songs *dll.List)(*dll.List){
    list := dll.New()
    list_Splits := getfiles_fromdir("Server/Songs_to_Send/Splits") //directories in the splits directory
    fmt.Println("Size: ", list_Splits.Size())
    for i :=0; i<list_Splits.Size(); i++ {
        temp_fileName, _ := list_Splits.Get(i)
        temp_list := getfiles_fromdir("Server/Songs_to_Send/Splits/"+temp_fileName.(string)) //directories in the splits directory
        it := temp_list.Iterator()
        for it.Next() {
                list.Add("Server/Songs_to_Send/Splits/"+temp_fileName.(string)+"/"+it.Value().(string))
        }
    }

    return list
}
