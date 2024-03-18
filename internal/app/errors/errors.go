package errors

import "errors"

var ErrURLConflict = errors.New(`url уже существует`)

var ErrNoDBConnection = errors.New(`нет подключения к бд`)
