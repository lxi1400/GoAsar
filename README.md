## Functions Documentation

### roundUp
```go
func roundUp(i, m int) int
```
The `roundUp` function takes an integer `i` and a modulus `m` and returns the smallest multiple of `m` that is greater than or equal to `i`.

### Asar struct
```go
type Asar struct {
    path       string
    fp         *os.File
    header     map[string]interface{}
    baseOffset int
}
```
The `Asar` struct represents an Asar archive. It contains the path to the archive file (`path`), a file pointer to the opened archive (`fp`), the parsed header information (`header`), and the base offset of the archive (`baseOffset`).

### OpenAsar
```go
func OpenAsar(path string) (*Asar, error)
```
The `OpenAsar` function opens an Asar archive located at the specified `path` and returns an `Asar` object. It reads and parses the archive header information, initializes the necessary fields in the `Asar` struct, and returns any errors encountered during the process.

### Close
```go
func (a *Asar) Close()
```
The `Close` method closes the file pointer associated with the `Asar` object, releasing any resources.

### ReadFileFromAsar
```go
func ReadFileFromAsar(path, filePath string) (string, error)
```
The `ReadFileFromAsar` function reads the content of a file specified by `filePath` from the Asar archive located at `path`. It searches for the file recursively in the archive's file structure and returns the content as a string. If the file is not found or any errors occur during the process, an error is returned.

### PrintFilePaths
```go
func PrintFilePaths(asarPath string) error
```
The `PrintFilePaths` function prints all file paths contained in the Asar archive located at `asarPath`. It recursively traverses the file structure of the archive and prints each file path to the standard output. If any errors occur during the process, an error is returned.

### printFilePaths (helper function)
```go
func printFilePaths(files map[string]interface{}, prefix string)
```
The `printFilePaths` function is a helper function used by `PrintFilePaths`. It recursively traverses the file structure represented by the `files` map and prints the full file paths to the standard output, prepending the `prefix` to each file path.

### readFileFromAsarHelper (helper function)
```go
func readFileFromAsarHelper(files map[string]interface{}, filePath string, fp *os.File, baseOffset int) (string, error)
```
The `readFileFromAsarHelper` function is a helper function used by `ReadFileFromAsar`. It recursively searches for the file specified by `filePath` in the file structure represented by the `files` map. It reads and returns the content of the file as a string. It takes the open file pointer (`fp`) and the base offset of the archive (`baseOffset`) to correctly read the file data.

### getFileOffset
```go
func getFileOffset(files map[string]interface{}, filePath string) (map[string]interface{}, bool)
```
The `getFileOffset` function searches for the file specified by `filePath` in the file structure represented by the `files` map. It returns the file information (as a map) and a boolean indicating whether the file was found or not.

### WriteFileToAsar
```go
func WriteFileToAsar(path, filePath string, content []byte) error
```
The `WriteFileTo

Asar` function writes the `content` byte array to the file specified by `filePath` in the Asar archive located at `path`. It searches for the file offset in the archive's file structure, seeks to the offset, and overwrites the file content with the provided content. If the file is not found or any errors occur during the process, an error is returned.