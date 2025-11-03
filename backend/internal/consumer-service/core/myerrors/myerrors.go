package myerrors

import "fmt"

var ErrDBConnClosed = fmt.Errorf("database connection is closed")
