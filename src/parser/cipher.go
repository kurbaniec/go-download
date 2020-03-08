package parser

type cipherOperations struct {
	operations []cipherOperation
}

func newCipherOperations() *cipherOperations {
	return &cipherOperations{operations: nil}
}

func (c *cipherOperations) addOperation(operation cipherOperation) {
	_ = append(c.operations, operation)
}

func (c *cipherOperations) decipher(input string) string {
	output := input
	for _, operation := range c.operations {
		output = operation.decipher(output)
	}
	return output
}

type cipherOperation interface {
	decipher(input string) string
}

type cipherReverse struct {
}

func newCipherReverse() *cipherReverse {
	return &cipherReverse{}
}

func (*cipherReverse) decipher(input string) string {
	return reverse(input)
}

type cipherSlice struct {
	index int
}

func newCipherSlice(index int) *cipherSlice {
	return &cipherSlice{index: index}
}

func (c *cipherSlice) decipher(input string) string {
	return input[c.index:]
}

type cipherSwap struct {
	index int
}

func newCipherSwap(index int) *cipherSwap {
	return &cipherSwap{index: index}
}

func (c *cipherSwap) decipher(input string) string {
	a, b := input[0], input[c.index]
	a, b = b, a
	return input
}

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
