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

func compressAndVerify(inputFile string) error {
	compressedFile := inputFile + ".gz"
	decompressedFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_decompressed" + filepath.Ext(inputFile)

	var wg sync.WaitGroup
	var compressErr, decompressErr error

	// Compress file
	wg.Add(1)
	go func() {
		defer wg.Done()
		compressErr = compressFile(inputFile, compressedFile)
	}()

	// Wait for compression to complete
	wg.Wait()

	if compressErr != nil {
		return fmt.Errorf("error compressing file: %v", compressErr)
	}

	// Decompress file
	wg.Add(1)
	go func() {
		defer wg.Done()
		decompressErr = decompressFile(compressedFile, decompressedFile)
	}()

	// Wait for decompression to complete
	wg.Wait()

	if decompressErr != nil {
		return fmt.Errorf("error decompressing file: %v", decompressErr)
	}

	// Verify integrity by comparing checksums
	valid, err := verifyIntegrity(inputFile, decompressedFile)
	if err != nil {
		return fmt.Errorf("error verifying integrity: %v", err)
	}

	if valid {
		fmt.Println("Verification successful: The decompressed file is identical to the original.")
	} else {
		fmt.Println("Verification failed: The decompressed file is not identical to the original.")
	}

	return nil
}

func main() {
	inputFile := flag.String("input", "", "Input file path")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Input file path is required.")
		flag.Usage()
		return
	}

	if err := compressAndVerify(*inputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
