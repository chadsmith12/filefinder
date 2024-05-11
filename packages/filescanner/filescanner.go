package filescanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type workData struct {
    file string
    pattern *regexp.Regexp
    result chan Result
    
}
type Result struct {
    File string
    LineNumber int 
    Text string
}

type FileWorker struct {
    Results chan Result
    jobs chan workData
    numberWorkers int
    allResults chan chan Result
}


func NewFileWorker(number int) *FileWorker {
    return &FileWorker{ numberWorkers: number, jobs: make(chan workData), allResults: make(chan chan Result), Results: make(chan Result) }
}

func (fw *FileWorker) StartWorkers(filePath string, pattern *regexp.Regexp) []Result {
    wg := sync.WaitGroup{}
    for i := 0; i < fw.numberWorkers; i++ {
        wg.Add(1)
        go func () { 
            defer wg.Done()
            fw.fileWorker() 
        }()
    }

    go fw.walkFilePath(filePath, pattern)
    
    var results []Result 
    for resultCh := range fw.allResults {
        for result := range resultCh {
            results = append(results, result)
        }
    }
    
    close(fw.jobs)
    wg.Wait()
    return results
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
        fw.jobs <- workData{ file: path, pattern: pattern, result: ch }
        fw.allResults <- ch

        return nil
    })
}

func (fw *FileWorker) fileWorker() {
    for workData := range fw.jobs {
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
                workData.result <- Result {
                    File: workData.file,
                    LineNumber: lineNumber,
                    Text: string(result),
                }
            }
        }
        file.Close()
        close(workData.result)
    }
}
