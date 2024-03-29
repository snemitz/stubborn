package stub

import "fmt"

var (
	ErrIncompleteCode  = fmt.Errorf("incomplete")
	ErrUnknownProperty = fmt.Errorf("unknown property")
)
