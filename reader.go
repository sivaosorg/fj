package fj

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// BufioRead reads all lines from an io.Reader and returns them as a single string.
//
// This function utilizes a buffered reader to read data line by line until EOF
// (end of file) is reached. It appends each line to a `strings.Builder`,
// which efficiently constructs the resulting string.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This could be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the lines read from the `io.Reader`, concatenated together.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function uses `bufio.NewReader` to create a buffered reader around the provided
//     `io.Reader`. This allows efficient reading of large input streams.
//   - It reads lines one by one using `ReadString('\n')`. If a line does not end with a
//     newline character before reaching EOF, the remaining text is still included in the
//     result.
//   - If an error occurs during reading (other than EOF), it stops reading and returns
//     the error along with an empty string.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := BufioRead(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text:")
//	content, err = BufioRead(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:", content)
func BufioRead(in io.Reader) (string, error) {
	var lines strings.Builder
	scanner := bufio.NewReader(in)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		lines.WriteString(line)
	}
	return lines.String(), nil
}

// BufioReadN reads all lines from an io.Reader and returns them as a slice of strings.
//
// This function utilizes a buffered reader to read data line by line until EOF
// (end of file) is reached. Each line is appended to a slice of strings.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A slice of strings, where each element corresponds to a line read from the `io.Reader`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function uses `bufio.NewReader` to create a buffered reader around the provided
//     `io.Reader`, allowing efficient reading of large input streams.
//   - It reads lines one by one using `ReadString('\n')`. Each line is appended to
//     a slice of strings. If a line does not end with a newline character before reaching EOF,
//     the remaining text is included as the last line in the slice.
//   - If an error occurs during reading (other than EOF), it stops reading and returns
//     the error along with the slice of lines read so far.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	lines, err := BufioReadN(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for i, line := range lines {
//	    fmt.Printf("Line %d: %s", i+1, line)
//	}
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	lines, err = BufioReadN(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	for _, line := range lines {
//	    fmt.Println(line)
//	}
func BufioReadN(in io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewReader(in)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return lines, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// IoRead reads all data from an io.Reader and returns it as a single string.
//
// This function reads raw bytes from the provided `io.Reader` in chunks and
// appends them to a buffer until EOF (end of file) is reached.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This could be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a buffer of size 1024 bytes to read chunks of data from
//     the `io.Reader`. This approach avoids reading the entire input into memory at once
//     and is suitable for handling large streams of data.
//   - The data read from each chunk is written to a `bytes.Buffer`, which efficiently
//     constructs the final string.
//   - If EOF is encountered during reading, the loop ends, and the accumulated data
//     in the buffer is returned as a string.
//   - Any other error during reading results in an immediate return of the error.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := IoRead(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	content, err = IoRead(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	fmt.Println(content)
func IoRead(in io.Reader) (string, error) {
	buf := make([]byte, 1024)
	var out bytes.Buffer
	for {
		n, err := in.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		out.Write(buf[:n])
	}
	return out.String(), nil
}

// BytesBufRead reads all data from an io.Reader and returns it as a single string.
//
// This function uses an `io.Copy` operation to efficiently read data from the provided
// `io.Reader` and write it to a `bytes.Buffer`. The resulting buffer content is then
// converted to a string and returned.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader`.
//   - An error if any I/O operation fails during the copy process. If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a `bytes.Buffer` to store the data read from the `io.Reader`.
//   - It uses the `io.Copy` function to transfer data from the `io.Reader` to the `bytes.Buffer`.
//     This approach is simple and efficient, leveraging built-in Go utilities for stream copying.
//   - After copying is complete, the data in the buffer is converted to a string and returned.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := BytesBufRead(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	content, err = BytesBufRead(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	fmt.Println(content)
func BytesBufRead(in io.Reader) (string, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, in); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TeeRead reads data from an io.Reader and logs the data simultaneously, returning the logged data as a string.
//
// This function utilizes `io.TeeReader` to split the data read from the provided `io.Reader`.
// The `io.TeeReader` writes a copy of all read data into a `bytes.Buffer` while continuing to pass
// the data through as it is read. The function consumes all data from the `io.Reader` but discards it,
// keeping only the logged copy.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any type that implements
//     the `io.Reader` interface, such as a file, standard input, or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader` as logged in the `bytes.Buffer`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is successfully
//     read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a `bytes.Buffer` to store a copy of the data being read from the `io.Reader`.
//   - An `io.TeeReader` is used to duplicate the data read: one copy is written to the buffer (`logs`),
//     and the other is read into a temporary buffer to consume the input stream.
//   - The function uses a loop to read from the `io.TeeReader` in chunks of 256 bytes until EOF is reached.
//   - The accumulated data in the `bytes.Buffer` is returned as a string after the reading is complete.
//
// Example Usage:
//
//	// Example: Reading from a file and logging data
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	loggedData, err := TeeRead(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Logged Data:")
//	fmt.Println(loggedData)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	loggedData, err = TeeRead(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered (logged):")
//	fmt.Println(loggedData)
func TeeRead(in io.Reader) (string, error) {
	var lines bytes.Buffer
	reader := io.TeeReader(in, &lines)
	buffer := make([]byte, 256)
	for {
		_, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return lines.String(), nil
}

// TeeProcess reads data from an io.Reader and writes the data simultaneously to an io.Writer,
// while also returning the logged data as a string.
//
// This function utilizes `io.TeeReader` to duplicate the data read from the provided `io.Reader`.
// A copy of the data is written to the provided `io.Writer`, while another copy is logged
// into a `bytes.Buffer`.
//
// Parameters:
//   - in: An `io.Reader` from which data will be read.
//   - out: An `io.Writer` to which the read data will be written.
//
// Returns:
//   - A string containing all the data read from the `io.Reader` as logged in the `bytes.Buffer`.
//   - An error if any I/O operation fails.
//
// Details:
//   - The function creates a `bytes.Buffer` to log a copy of the data.
//   - It uses `io.TeeReader` to split the data read from `in` between logging to the buffer
//     and passing the data to the consumer (`out`).
//   - A loop is used to read data in chunks and write it to the `io.Writer`.
//   - The function ensures that both reading and writing are handled efficiently.
//
// Example Usage:
//
//	// Example: Reading from a file and writing to another file while logging
//	inFile, err := os.Open("input.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer inFile.Close()
//
//	outFile, err := os.Create("output.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer outFile.Close()
//
//	loggedData, err := TeeProcess(inFile, outFile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Logged Data:")
//	fmt.Println(loggedData)
//
//	// Example: Reading from standard input and writing to standard output
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	loggedData, err = TeeProcess(os.Stdin, os.Stdout)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("\nLogged Data:")
//	fmt.Println(loggedData)
func TeeProcess(in io.Reader, out io.Writer) (string, error) {
	var lines bytes.Buffer
	reader := io.TeeReader(in, &lines)
	buffer := make([]byte, 256)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				return "", writeErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return lines.String(), nil
}
