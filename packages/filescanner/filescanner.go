package filescanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func randomDuration() time.Duration {
	randomMs := rand.Intn(1000) + 500
	duration := time.Duration(randomMs) * time.Millisecond

	return duration
}

func randomBool() bool {
	random := rand.Intn(100)

	return random > 60
}

type workData struct {
	file    string
	pattern *regexp.Regexp
	result  chan Result
}
type Result struct {
	File       string
	LineNumber int
	Text       string
	WorkerId   int
}

type FileWorker struct {
	jobs          chan workData
	numberWorkers int
	allResults    chan chan Result
	wg            sync.WaitGroup
}

func NewFileWorker(number int) *FileWorker {
	return &FileWorker{
		numberWorkers: number,
		jobs:          make(chan workData),
		allResults:    make(chan chan Result),
		wg:            sync.WaitGroup{},
	}
}

func (fw *FileWorker) StartWorkers(filePath string, pattern *regexp.Regexp) {
	for i := 0; i < fw.numberWorkers; i++ {
		fw.wg.Add(1)
		go func(id int) {
			defer fw.wg.Done()
			fw.fileWorker(id)
		}(i + 1)
	}

	go fw.walkFilePath(filePath, pattern)
}

func (fw *FileWorker) Read() <-chan Result {
	returnedResult := make(chan Result)
	go func() {
		for resultCh := range fw.allResults {
			for result := range resultCh {
				returnedResult <- result
			}
		}

		close(returnedResult)
		close(fw.jobs)
		fw.wg.Wait()
	}()

	return returnedResult
}

func (fw *FileWorker) walkFilePath(filePath string, pattern *regexp.Regexp) {
	defer close(fw.allResults)
	filepath.Walk(filePath, func(path string, fileInfo fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip directories
		if fileInfo.IsDir() {
			return nil
		}

		ch := make(chan Result)
		fw.jobs <- workData{file: path, pattern: pattern, result: ch}
		fw.allResults <- ch

		return nil
	})
}

func (fw *FileWorker) fileWorker(id int) {
	for workData := range fw.jobs {
		// if randomBool() {
		// 	time.Sleep(randomDuration())
		// }
		file, err := os.Open(workData.file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		scan := bufio.NewScanner(file)
		lineNumber := 1
		for scan.Scan() {
			result := workData.pattern.Find(scan.Bytes())
			if len(result) > 0 {
				workData.result <- Result{
					File:       workData.file,
					LineNumber: lineNumber,
					Text:       string(result),
					WorkerId:   id,
				}
			}
			lineNumber++
		}
		file.Close()
		close(workData.result)
	}
}
