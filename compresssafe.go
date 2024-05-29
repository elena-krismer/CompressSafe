package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileStatus struct {
	Path    string
	Success bool
	Error   error
}

func compressFile(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := gzip.NewWriter(out)
	defer writer.Close()

	_, err = io.Copy(writer, in)
	return err
}

func decompressFile(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	reader, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(out, reader)
	return err
}

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyIntegrity(originalFile, decompressedFile string) (bool, error) {
	originalChecksum, err := calculateChecksum(originalFile)
	if err != nil {
		return false, err
	}

	decompressedChecksum, err := calculateChecksum(decompressedFile)
	if err != nil {
		return false, err
	}

	return originalChecksum == decompressedChecksum, nil
}

func processFile(path string, relativePath string, statuses *[]FileStatus, wg *sync.WaitGroup) {
	defer wg.Done()

	compressedFile := path + ".gz"
	decompressedFile := filepath.Join("decompressed", relativePath)

	compressErr := compressFile(path, compressedFile)
	if compressErr != nil {
		*statuses = append(*statuses, FileStatus{Path: path, Success: false, Error: compressErr})
		return
	}

	if err := os.MkdirAll(filepath.Dir(decompressedFile), os.ModePerm); err != nil {
		*statuses = append(*statuses, FileStatus{Path: path, Success: false, Error: fmt.Errorf("error creating directories for %s: %v", decompressedFile, err)})
		return
	}

	decompressErr := decompressFile(compressedFile, decompressedFile)
	if decompressErr != nil {
		*statuses = append(*statuses, FileStatus{Path: path, Success: false, Error: decompressErr})
		return
	}

	valid, verifyErr := verifyIntegrity(path, decompressedFile)
	if verifyErr != nil {
		*statuses = append(*statuses, FileStatus{Path: path, Success: false, Error: verifyErr})
		return
	}

	if valid {
		*statuses = append(*statuses, FileStatus{Path: path, Success: true, Error: nil})
	} else {
		*statuses = append(*statuses, FileStatus{Path: path, Success: false, Error: fmt.Errorf("verification failed")})
	}
}

func compressAndVerify(inputPath string) error {
	var wg sync.WaitGroup
	var statuses []FileStatus

	if fileInfo, err := os.Stat(inputPath); err == nil && fileInfo.IsDir() {
		err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && !strings.HasSuffix(path, ".gz") {
				relativePath := strings.TrimPrefix(path, inputPath)
				wg.Add(1)
				go processFile(path, relativePath, &statuses, &wg)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking the path %s: %v", inputPath, err)
		}
	} else if !strings.HasSuffix(inputPath, ".gz") {
		wg.Add(1)
		go processFile(inputPath, filepath.Base(inputPath), &statuses, &wg)
	} else {
		fmt.Printf("Skipping already compressed file: %s\n", inputPath)
	}

	wg.Wait()

	// Print summary
	successCount := 0
	for _, status := range statuses {
		if status.Success {
			successCount++
		} else {
			fmt.Printf("Error processing file %s: %v\n", status.Path, status.Error)
		}
	}

	fmt.Printf("Processed %d files: %d successful, %d failed\n", len(statuses), successCount, len(statuses)-successCount)

	// Remove decompressed directory
	if err := os.RemoveAll(decompressedDir); err != nil {
		return fmt.Errorf("error removing decompressed directory: %v", err)
	}

	return nil
}

func main() {
	inputPath := flag.String("input", "", "Input file or directory path")

	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Input file or directory path is required.")
		flag.Usage()
		return
	}

	if err := compressAndVerify(*inputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
