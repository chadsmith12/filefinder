package worker

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Result struct {
	File       string
	LineNumber int
	Path       string
	Text       string
}

type workerData struct {
	path    string
	pattern *regexp.Regexp
}

func (wd workerData) Execute() interface{} {
	file, err := os.Open(wd.path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	scan := bufio.NewScanner(file)
	lineNumber := 1
	results := make([]Result, 0)
	for scan.Scan() {
		result := wd.pattern.Find(scan.Bytes())
		if len(result) > 0 {
			workResult := Result{
				File:       wd.path,
				LineNumber: lineNumber,
				Text:       string(result),
			}
			results = append(results, workResult)
		}
		lineNumber++
	}
	file.Close()

	return results
}

type FileWorker struct {
	worker  *WorkerPool
	path    string
	pattern *regexp.Regexp
}

func NewFileWorker(path, pattern string) *FileWorker {
	regexPattern := regexp.MustCompile(pattern)
	return &FileWorker{
		worker:  NewWorkerPool(3),
		pattern: regexPattern,
		path:    path,
	}
}

func (fw *FileWorker) Start() {
	go fw.walk()
}

func (fw *FileWorker) Result() <-chan Result {
	resultChan := make(chan Result)
	go func() {
		defer close(resultChan)
		for results := range fw.worker.Result() {
			if results == nil {
				continue
			}
			fileResults := results.([]Result)
			for _, fileResult := range fileResults {
				resultChan <- fileResult
			}
		}
	}()

	return resultChan
}

func (fw *FileWorker) walk() {
	filepath.WalkDir(fw.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}
		data := workerData{
			path:    path,
			pattern: fw.pattern,
		}
		fw.worker.Add(data)

		return nil
	})

	fw.worker.Stop()
}
