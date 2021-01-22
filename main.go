package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type artistStruct struct {
	Name         string
	TimesPlayed  int
	TimeListened int
}

type mainStruct []struct {
	EndTime    string `json:"endTime"`
	ArtistName string `json:"artistName"`
	TrackName  string `json:"trackName"`
	MsPlayed   int    `json:"msPlayed"`
}

func convert(thing string) int {
	kek, err := strconv.Atoi(thing)
	if err != nil {
		return 1
	}
	return kek
}

var (
	red    = color("\033[1;31m%s\033[0m")
	green  = color("\033[1;32m%s\033[0m")
	yellow = color("\033[1;33m%s\033[0m")
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//got this from stackoverflow, cannot remember which one though. mb.
func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func makeMainSlice(tosort map[string]*artistStruct) []*artistStruct {
	var mainSlice []*artistStruct
	for _, data := range tosort {
		mainSlice = append(mainSlice, data)
	}
	sort.Slice(mainSlice, func(i, j int) bool {
		return mainSlice[i].TimeListened > mainSlice[j].TimeListened
	})
	return mainSlice
}

func makePrintAndLogStrings(tomake []*artistStruct) (string, string) {
	var return1, return2 string
	//return1 = color return2 = nocolor
	for pos, artist := range tomake {
		return1 += fmt.Sprintf("%s.) %s\n	Times played: %s\n	Time listened total (seconds): %s\n", fmt.Sprint(pos+1), artist.Name, fmt.Sprint(artist.TimesPlayed), fmt.Sprint(artist.TimeListened))
		return2 += green(fmt.Sprintf("%s.) %s\n", fmt.Sprint(pos), artist.Name)) + red(fmt.Sprintf("	Times played: %s\n	Time listened total (seconds): %s\n", fmt.Sprint(artist.TimesPlayed), fmt.Sprint(artist.TimeListened)))
	}
	return return2, return1
}

//https://stackoverflow.com/a/10510783
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	exist := exists("./spotify data")
	if !exist {
		fmt.Println(red("You need a folder named \"spotify data\" with your streaming history in it!"))
		return
	}
	var totaltimelistened = 0
	var totalSongsPlayed = 0
	toPrint := map[string]*artistStruct{}
	var num = 0
	var tmp string
	fmt.Print(green("How many \"StreamingHistory\" json files are in your \"spotify data\" folder? "))
	fmt.Scanln(&tmp)
	num = convert(tmp)
	if num < 1 {
		fmt.Println(red("You need 1 or more files to read!"))
		return
	}
	for i := 0; i < num; i++ {
		if !exists("./spotify data/StreamingHistory" + fmt.Sprint(i) + ".json") {
			fmt.Println(red("StreamingHistory" + fmt.Sprint(i) + ".json doesn't exist!\n"))
			return
		}
		input, err := ioutil.ReadFile("./spotify data/StreamingHistory" + fmt.Sprint(i) + ".json")
		if err != nil {
			fmt.Println("Error reading file StreamingHistory"+fmt.Sprint(i)+".json with error:", err)
			return
		}
		var allArtists mainStruct
		err = json.Unmarshal(input, &allArtists)
		if err != nil {
			fmt.Println("Error parsing file StreamingHistory"+fmt.Sprint(i)+".json with error:", err)
			return
		}
		for _, song := range allArtists {
			totaltimelistened += song.MsPlayed
			totalSongsPlayed++
			if _, ok := toPrint[song.ArtistName]; ok {
				toPrint[song.ArtistName].TimeListened += song.MsPlayed / 1000
				toPrint[song.ArtistName].TimesPlayed++
			} else {
				toPrint[song.ArtistName] = &artistStruct{
					Name:         song.ArtistName,
					TimeListened: song.MsPlayed / 1000,
					TimesPlayed:  1,
				}
			}
		}
	}
	if len(toPrint) < 1 {
		fmt.Println(red("No spotify data detected! Make sure you have a foled named \"spotify data\" containing your streaming history!\n"))
	}
	fmt.Println(yellow("\nYou've listened to " + fmt.Sprint(len(toPrint)) + " artists total!\n"))
	fmt.Println(yellow("You've listened to a whopping " + fmt.Sprint(totalSongsPlayed) + " songs elapsing a total:\n	MilliSeconds: " + fmt.Sprint(totaltimelistened) + "ms\n	Seconds: " + fmt.Sprint(totaltimelistened/1000) + "s\n	Minutes: " + fmt.Sprint((totaltimelistened/1000)/60) + "min\n	Hours: " + fmt.Sprint(((totaltimelistened/1000)/60)/60) + "h\n	Days: " + fmt.Sprint((((totaltimelistened/1000)/60)/60)/24) + "D"))
	fmt.Println(red("\nI am now sorting all your data! This may take a few seconds depending on how much data you have!\n"))
	start := makeTimestamp()
	mainSlice := makeMainSlice(toPrint)
	printdata, logdata := makePrintAndLogStrings(mainSlice)
	fmt.Println(yellow("\nDone storting in " + fmt.Sprint(makeTimestamp()-start) + "ms!\n"))
	var option1 string
	fmt.Print(green("Do you wish to print the data to console (y / yes): "))
	fmt.Scanln(&option1)
	option1 = strings.ToLower(option1)
	if option1 == "y" || option1 == "yes" {
		fmt.Println(printdata)
	}
	var option string
	fmt.Print("\n" + green("Log all data to file (log.txt needs to exist) (y or yes for yes): "))
	fmt.Scanln(&option)
	option = strings.ToLower(option)
	if option == "y" || option == "yes" {
		if !exists("./log.txt") {
			fmt.Println(red("\nFile named \"log.txt\" does not exist, please create one and try again!"))
			return
		}
		f, err := os.Create("log.txt")
		if err != nil {
			fmt.Println("File with name log.txt doesn't exist or cannot be written to!")
			f.Close()
			return
		}
		f.WriteString(logdata)
		f.Close()
		fmt.Println(green("\nData written to log.txt!"))
	}
}
