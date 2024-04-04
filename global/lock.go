package globaldata

import "sync"

var OpcReadLock = &sync.Mutex{}
var OpcWriteLock = &sync.Mutex{}

var WebHookWriteLock = &sync.Mutex{}
