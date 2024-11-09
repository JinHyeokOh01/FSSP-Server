package db

import (
    "sync"
)

var (
    CurrentUser map[string]interface{}
    Mu          sync.Mutex
)