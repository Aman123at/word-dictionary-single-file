# Word Dictionary using Single File in Go

This document outlines a Go program that implements a simple word dictionary using a single text file (`dictionary.txt`) to store both the words and their corresponding definitions. The dictionary allows users to look up word definitions quickly, leveraging file offset techniques for efficient access. The program supports both direct lookup and fast lookup methods.

## Overview

The program uses a single text file, `dictionary.txt`, to store a list of words and their corresponding definitions. The dictionary is written in two parts:

1. **Headers**: Contains word-to-offset mappings (i.e., the file positions where definitions are located).
2. **Body**: Contains the actual word and definition pairs.
Two methods are provided for looking up definitions:

- **Basic lookup**: Scans the file to find the word and then retrieves the definition.
- **Fast lookup**: Uses a preloaded map of word-to-offset values for quicker access.

## Dictionary Structure

### Writing the Dictionary
The function `createDictionary` writes predefined dictionary entries to the `dictionary.txt` file. The file contains two sections:

1. **Headers**: Each word is followed by the offset (i.e., the byte position in the file where its definition starts).
2. **Body**: The definitions, which are written after all the headers.

```go
func createDictionary() error {
    // Code to write headers and definitions to the file
}
```

### Header Calculation
The total offset for definitions is calculated to ensure correct placement in the file after the headers.

### Sample Entries
Sample dictionary entries are predefined in the code, such as:
```go
{"apple", "A round fruit with red or green skin and white flesh."},
{"book", "A written or printed work consisting of pages glued or sewn together."},
// More entries...
```

## Word Lookup
### Basic Lookup
The `lookupWord` function searches for a word in the file by scanning the header, then seeking to the offset and retrieving the corresponding definition.
```go
func lookupWord(word string) (string, error) {
    // Code to search for a word and fetch its definition
}
```

### Fast Lookup
The `lookupWordFast` function preloads the word-to-offset mappings into a `wordOffsets` map during initialization, allowing for quicker lookups by avoiding repeated file scans.
```go
func lookupWordFast(word string) (string, error) {
    // Code for fast word lookup using preloaded offsets
}
```

## Steps To Run:
1. **Create Dictionary**: The program first calls `createDictionary` to generate the dictionary file.
2. **Word Lookup**: It then looks up the definition of the word using fast lookup method (`lookupWordFast`) or basic lookup method (`lookupWord`).
3. **Timing**: The time taken for the lookup is printed to the console.



## Time Benchmark:
- Current Data Sample: ***50 entries***
- Basic Lookup: ***~0.00568 seconds***
- Fast Lookup: ***~0.000045 seconds***
- Result: Fast Lookup is almost ***10x*** fater then basic lookup, it will be much more faster for larger data set.


## Error Handling
The program handles various errors, including:

- **File Creation/Reading Errors**: Handled during file creation and access operations.
- **Lookup Errors**: If a word is not found or there is a problem reading from the file, an appropriate error message is returned.


## Conclusion

This Go program demonstrates how to create and manage a simple word dictionary stored in a single file, with support for both basic and fast lookup methods. By using file offsets, it provides efficient access to definitions, especially when using the fast lookup with preloaded offsets.




