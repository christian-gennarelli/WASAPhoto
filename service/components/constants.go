package components

const USERNAME_REGEXP = "" //"/^[a-zA-Z0-9]([._-](?![._-])|[a-zA-Z0-9]){6,14}[a-zA-Z0-9]$/"
const ID_REGEXP = ""       //"/^[A-Z0-9]{64}$/"

const StatusInternalServerError = "{\"ErrorCode\": 500, \"Description\": \"Internal Server Error: %s\"}"
const StatusBadRequest = "{\"ErrorCode\": 400, \"Description\": \"Bad Request: %s\"}"
const StatusUnauthorized = "{\"ErrorCode\": 401, \"Description\": \"Unauthorized: %s\"}"
const StatusNotFound = "{\"ErrorCode\": 404, \"Description\": \"Resource Not Found: %s\"}"
