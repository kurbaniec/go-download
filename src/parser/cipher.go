package parser

// Holds all operations neede to decipher URL signature.
type CipherOperations struct {
	operations []cipherOperation
}

// Returns a new CipherOperations struct.
func newCipherOperations() *CipherOperations {
	return &CipherOperations{operations: []cipherOperation{}}
}

// Adds a cipher operation to the list.
func (c *CipherOperations) addOperation(operation cipherOperation) {
	c.operations = append(c.operations, operation)
}

// Returns the deciphered output of an input.
func (c *CipherOperations) decipher(input string) string {
	output := input
	for _, operation := range c.operations {
		output = operation.decipher(output)
	}
	return output
}

// Interface for cipher operations.
type cipherOperation interface {
	decipher(input string) string
}

// Represents a reverse operation.
type cipherReverse struct {
}

// Returns a reverse operation.
func newCipherReverse() *cipherReverse {
	return &cipherReverse{}
}

// Performs a reverse operation.
func (*cipherReverse) decipher(input string) string {
	tmp := reverse(input)
	return tmp
}

// Represents a slice operation.
type cipherSlice struct {
	index int
}

// Returns a slice operation.
func newCipherSlice(index int) *cipherSlice {
	return &cipherSlice{index: index}
}

// Performs a slice operation.
func (c *cipherSlice) decipher(input string) string {
	// return input[c.index:]
	//return string([]rune(input)[c.index:])
	output := []rune(input)
	tmp := string(output[:c.index]) + string(output[c.index*2:])
	return tmp
}

// Represents a swap operation.
type cipherSwap struct {
	index int
}

// Returns a swap operation.
func newCipherSwap(index int) *cipherSwap {
	return &cipherSwap{index: index}
}

// Performs a swap operation.
func (c *cipherSwap) decipher(input string) string {
	out := []rune(input)
	out[0], out[c.index] = out[c.index%len(out)], out[0]
	tmp := string(out)
	return tmp
}

// Reverses a string input and returns it.
func reverse(input string) string {
	n := 0
	theRune := make([]rune, len(input))
	for _, r := range input {
		theRune[n] = r
		n++
	}
	theRune = theRune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		theRune[i], theRune[n-1-i] = theRune[n-1-i], theRune[i]
	}
	// Convert back to UTF-8.
	output := string(theRune)
	return output
}
