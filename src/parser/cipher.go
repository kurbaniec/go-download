package parser

type CipherOperations struct {
	operations []cipherOperation
}

func newCipherOperations() *CipherOperations {
	return &CipherOperations{operations: []cipherOperation{}}
}

func (c *CipherOperations) addOperation(operation cipherOperation) {
	c.operations = append(c.operations, operation)
}

func (c *CipherOperations) decipher(input string) string {
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
	tmp := reverse(input)
	return tmp
}

type cipherSlice struct {
	index int
}

func newCipherSlice(index int) *cipherSlice {
	return &cipherSlice{index: index}
}

func (c *cipherSlice) decipher(input string) string {
	// return input[c.index:]
	//return string([]rune(input)[c.index:])
	output := []rune(input)
	tmp := string(output[:c.index]) + string(output[c.index*2:])
	return tmp
}

type cipherSwap struct {
	index int
}

func newCipherSwap(index int) *cipherSwap {
	return &cipherSwap{index: index}
}

func (c *cipherSwap) decipher(input string) string {
	out := []rune(input)
	out[0], out[c.index] = out[c.index%len(out)], out[0]
	tmp := string(out)
	return tmp
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
