package components

import "fmt"

const USERNAME_REGEXP = "^[a-zA-Z0-9_-]{8,16}$"
const ID_REGEXP = "^[a-z0-9]{64}$"
const DATETIME_REGEXP = "^([0-9]{4})-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[01])$"
const COMMENT_REGEXP = "^[a-zA-Z0-9.,!?@#%^&*()_+-=:;'\"<>/[\\]{}`~\\s]{1,150}$"

const StatusInternalServerError = "{\"ErrorCode\": 500, \"Description\": \"Internal Server Error: %s\"}"
const StatusBadRequest = "{\"ErrorCode\": 400, \"Description\": \"Bad Request: %s\"}"
const StatusUnauthorized = "{\"ErrorCode\": 401, \"Description\": \"Unauthorized: %s\"}"
const StatusForbidden = "{\"ErrorCode\": 403, \"Description\": \"Unauthenticated: %s\"}"
const StatusNotFound = "{\"ErrorCode\": 404, \"Description\": \"Resource Not Found: %s\"}"
const StatusUnsupportedMediaType = "{\"ErrorCode\": 415, \"Description\": \"Unsupported media type\"}"

var ErrIDNotValid = fmt.Errorf("provided ID not valid")
var ErrUsernameNotValid = fmt.Errorf("provided username not valid")
var ErrCommentNotValid = fmt.Errorf("provided comment not valid")
var ErrForeignKeyConstraint = fmt.Errorf("FOREIGN KEY constraint failed")
var ErrUniqueConstraintUsername = fmt.Errorf("UNIQUE constraint failed: User.Username")
