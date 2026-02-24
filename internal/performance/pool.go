package performance

import "sync"

var BytesPool = sync.Pool{New: func() any { b := make([]byte, 0, 4096); return &b }}

func GetBytes() *[]byte  { return BytesPool.Get().(*[]byte) }
func PutBytes(b *[]byte) { *b = (*b)[:0]; BytesPool.Put(b) }
