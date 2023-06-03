package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func roundUp(i, m int) int {
	return (i + m - 1) & ^(m - 1)
}

type Asar struct {
	path       string
	fp         *os.File
	header     map[string]interface{}
	baseOffset int
}

func OpenAsar(path string) (*Asar, error) {
	fp, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			fp.Close()
		}
	}()

	var data_size, header_size, header_object_size, header_string_size uint32
	if err := binary.Read(fp, binary.LittleEndian, &data_size); err != nil {
		return nil, err
	}
	if err := binary.Read(fp, binary.LittleEndian, &header_size); err != nil {
		return nil, err
	}
	if err := binary.Read(fp, binary.LittleEndian, &header_object_size); err != nil {
		return nil, err
	}
	if err := binary.Read(fp, binary.LittleEndian, &header_string_size); err != nil {
		return nil, err
	}

	headerJson := make([]byte, header_string_size)
	if _, err := fp.Read(headerJson); err != nil {
		return nil, err
	}

	header := make(map[string]interface{})
	if err := json.Unmarshal(headerJson, &header); err != nil {
		return nil, err
	}

	baseOffset := roundUp(int(16+header_string_size), 4)

	return &Asar{
		path:       path,
		fp:         fp,
		header:     header,
		baseOffset: baseOffset,
	}, nil
}

func (a *Asar) Close() {
	a.fp.Close()
}

func ReadFileFromAsar(path, filePath string) (string, error) {
	asar, err := OpenAsar(path)
	if err != nil {
		return "", err
	}
	defer asar.Close()

	content, err := readFileFromAsarHelper(asar.header["files"].(map[string]interface{}), filePath, asar.fp, asar.baseOffset)
	if err != nil {
		return "", err
	}

	return content, nil
}

func PrintFilePaths(asarPath string) error {
	asar, err := OpenAsar(asarPath)
	if err != nil {
		return err
	}
	defer asar.Close()

	fmt.Println("File Paths:")
	printFilePaths(asar.header["files"].(map[string]interface{}), "")

	return nil
}

func printFilePaths(files map[string]interface{}, prefix string) {
	for name, info := range files {
		fullPath := filepath.Join(prefix, name)
		if subFiles, ok := info.(map[string]interface{}); ok {
			fmt.Println(fullPath + "/")
			printFilePaths(subFiles, fullPath)
		} else {
			fmt.Println(fullPath)
		}
	}
}

func readFileFromAsarHelper(files map[string]interface{}, filePath string, fp *os.File, baseOffset int) (string, error) {
	for name, info := range files {
		fullPath := filepath.Join(name)

		if fullPath == filePath {
			if fileData, ok := info.(map[string]interface{}); ok {
				if offset, ok := fileData["offset"].(string); ok {
					offsetVal, err := strconv.Atoi(offset)
					if err != nil {
						return "", err
					}

					size, ok := fileData["size"].(float64)
					if !ok {
						return "", fmt.Errorf("file '%s' has invalid size in the Asar archive", filePath)
					}

					fp.Seek(int64(baseOffset+offsetVal), 0)
					data := make([]byte, int(size))
					if _, err := fp.Read(data); err != nil {
						return "", err
					}

					return string(data), nil
				}
			}

			return "", fmt.Errorf("file '%s' does not contain data in the Asar archive", filePath)
		}

		if subFiles, ok := info.(map[string]interface{}); ok {
			subData, err := readFileFromAsarHelper(subFiles, filePath, fp, baseOffset)
			if err == nil {
				return subData, nil
			}
		}
	}

	return "", fmt.Errorf("file '%s' not found in the Asar archive", filePath)
}

func getFileOffset(files map[string]interface{}, filePath string) (map[string]interface{}, bool) {
	segments := strings.Split(filePath, "/")
	currentFiles := files

	for _, segment := range segments {
		if fileInfo, ok := currentFiles[segment].(map[string]interface{}); ok {
			currentFiles = fileInfo
		} else {
			return nil, false
		}
	}

	return currentFiles, true
}

func WriteFileToAsar(path, filePath string, content []byte) error {
	asar, err := OpenAsar(path)
	if err != nil {
		return err
	}
	defer asar.Close()

	// Find the offset where the new file should be written
	fileInfo, ok := getFileOffset(asar.header["files"].(map[string]interface{}), filePath)
	if !ok {
		return fmt.Errorf("file '%s' not found in the Asar archive", filePath)
	}

	offset, ok := fileInfo["offset"].(string)
	if !ok {
		return fmt.Errorf("invalid offset for file '%s' in the Asar archive", filePath)
	}

	fileOffset, err := strconv.Atoi(offset)
	if err != nil {
		return err
	}

	// Seek to the file offset
	_, err = asar.fp.Seek(int64(asar.baseOffset+fileOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Write the file content
	_, err = asar.fp.Write(content)
	if err != nil {
		return err
	}

	return nil
}
