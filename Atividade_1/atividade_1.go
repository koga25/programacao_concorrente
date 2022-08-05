package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	benchmarkTreads(files, dir)
}

func benchmarkTreads(files []fs.FileInfo, dir string) {
	start := time.Now()

	//1 thread
	times := make([]int64, 0)
	for i := 0; i < len(files); i++ {
		if strings.Contains(files[i].Name(), ".txt") {
			getSites(dir + "/" + files[i].Name())
		}
	}
	times = append(times, time.Since(start).Milliseconds())

	//2 threads
	start = time.Now()
	var wg sync.WaitGroup
	j := 0
	for i := 0; i < len(files); i++ {
		if strings.Contains(files[i].Name(), ".txt") {
			j++
			wg.Add(1)
			go getSitesThreads(dir+"/"+files[i].Name(), &wg)
		}
		if j == 2 {
			wg.Wait()
		}
	}
	wg.Wait()
	times = append(times, time.Since(start).Milliseconds())

	//4 threads
	start = time.Now()
	for i := 0; i < len(files); i++ {
		if strings.Contains(files[i].Name(), ".txt") {
			wg.Add(1)
			go getSitesThreads(dir+"/"+files[i].Name(), &wg)
		}
	}
	wg.Wait()
	times = append(times, time.Since(start).Milliseconds())

	for i := 0; i < len(times); i++ {
		duration := times[i]
		fmt.Printf("time since, %v° execution, is: %v\n", i, duration)
	}

}

func getSites(filePath string) {
	// Current working directory

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sites := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		sites = append(sites, scanner.Text())
	}

	scannerErr := scanner.Err()
	if scannerErr != nil {
		log.Fatal(scannerErr)
	}

	file.Close()

	for i := 0; i < len(sites); i++ {
		getRequest(sites[i])
	}
}

func getSitesThreads(filePath string, wg *sync.WaitGroup) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sites := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		sites = append(sites, scanner.Text())
	}

	scannerErr := scanner.Err()
	if scannerErr != nil {
		log.Fatal(scannerErr)
	}

	file.Close()

	for i := 0; i < len(sites); i++ {
		getRequest(sites[i])
	}

	wg.Done()
}

func getRequest(site string) {
	resp, err := http.Get(site)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
}
