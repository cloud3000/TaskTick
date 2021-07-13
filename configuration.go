package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// {"tick":1, "script":"./scripts/tick1.sh"},
type task struct {
	Tick   int    `json:"tick"`
	Script string `json:"script"`
}

type que struct {
	Title string `json:"title"`
	Tasks []task `json:"tasks"`
}

type clockConf struct {
	Secque que `json:"sec"`
	Minque que `json:"min"`
	Hrque  que `json:"hr"`
}

func readConf(fname string) clockConf {
	// json data
	var obj clockConf

	// read file
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Print(err)
		return obj
	}

	// Unmarshal json.
	// NOTE: When we edit clock.json with VsCode we get a json error, but the Unmarshal is successful.
	// ALSO NOTE: Some editors do not write to the real file until you exit.
	// For best results use VI to edit the clock.json file.
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		return obj
	}

	return obj
}

func configuration(wg *sync.WaitGroup, dir string, fname string, done <-chan bool, chg chan bool) {
	defer wg.Done()
	configuration, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer configuration.Close()

	go func() {
		for {
			select {
			case event, ok := <-configuration.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if fname == event.Name {
						chg <- true
					}
				}
			case err, ok := <-configuration.Errors:
				if !ok {
					return
				}
				log.Println("configuration error:", err)
			}
		}
	}()
	err = configuration.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration supports realtime config changes in %s.\n", fname)
	<-done
	log.Println("configuration is shuting down.")
}
