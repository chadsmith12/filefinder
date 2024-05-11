# filefinder

filefinder is a simple grep like utility tool. The main goal of this tool is to search for a pattern inside files.

The tool will use a number of worker threads specified by the user, defaulting to 3, to find the pattern specified. You give the tool the path to start from and the pattern.

The main goal is to learn concurrency inside of Go so it uses a worker pool. The scanner will send the files off to the jobs channel with the workers listening to the jobs channel for more work to preform.

## Building from source

Simply run `go build` to build from source.

## Usage

* Searching for a pattern starting from the current directory: `filefinder ./ "func main"`
    - Will attempt to find the string "func main" starting from the current directory.
* Change number worker threads: `filefinder ./ "func main" -w 10`
    - Uses 10 workers to attempt to find the string "func main" starting from the current directory.

