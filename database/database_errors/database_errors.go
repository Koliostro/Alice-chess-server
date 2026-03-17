package database_errors

import "errors"

var SQLErrObjDup = errors.New("Object duplication")
var SQLErrUnexp = errors.New("Unexpected sql error")
