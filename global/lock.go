package globaldata

import "sync"

var OpcWriteLock = &sync.Mutex{}

var WebHookWriteLock = &sync.Mutex{}
