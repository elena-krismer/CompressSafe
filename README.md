# CompressSafe
 
CompressSafe is a command-line tool for compressing files using gzip, decompressing them, verifying their integrity, and cleaning up temporary decompressed files. It supports processing individual files as well as entire directories, while ignoring already compressed files.

## Usage 
### Running on Windows

1. Download compresssafe.exe
2. Run in the commandline

```sh
compresssafe.exe -input <path>
```

### Running on macOS/Linux
```sh
compresssafe -input <path>
```

### Parameters
- `input `: Specifies the input file or directory path to be processed.


## Features

- Compresses files using gzip.
- Decompresses gzip-compressed files.
- Verifies the integrity of decompressed files by comparing checksums.
- Processes directories recursively.
- Skips files that are already compressed.
- Removes temporary decompressed files after processing.


## Installation advanced

### Requirements

- Go 1.16 or later

1. Clone the repository or download the source code.
2. Navigate to the directory containing the source code.
3. Build the executable for your target operating system.

### Building for Windows

```sh
GOOS=windows GOARCH=amd64 go build -o compresssafe.exe compresssafe.go
```

## How it Works
1. Compression: The tool compresses each file using gzip and appends the .gz extension.
2. Decompression: The tool decompresses the compressed file to a temporary directory named decompressed.
3. Verification: The tool calculates and compares the SHA-256 checksum of the original file and the decompressed file to ensure they are identical.
4. Cleanup: After processing all files, the tool removes the temporary decompressed directory.