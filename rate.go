package mymain

import (
    "bytes"
    "fmt"
    "io"
    "time"
    //"context"
    "golang.org/x/time/rate"
)

type reader struct {
    r      io.Reader
    limiter *rate.Limiter
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func NewReader(r io.Reader, l *rate.Limiter) io.Reader {
    return &reader{
        r:      r,
        limiter:l,
    }
}

func (r *reader) Read(buf []byte) (int, error) {
    now := time.Now()
    rv := r.limiter.AllowN(now, len(buf))
    if !rv{
        return 0, nil
    }
    n, err := r.r.Read(buf)
    if n <= 0 {
        return n, err
    }
    //fmt.Printf("Read %d bytes\n", n)
    return n, err
}

func mymain() {
    // Source holding 1MB
    counter := 0
    timeslot := 0
    src := bytes.NewReader(make([]byte, 1024*1024))
    // Destination
    dst := &bytes.Buffer{}

    // Bucket adding 100KB every second, holding max 100KB
    limit := rate.NewLimiter(100*1024, 100*1024)
    limit.SetBurst(100*1024)

    start := time.Now()

    buf := make([]byte, 10*1024)
    // Copy source to destination, but wrap our reader with rate limited one
    //io.CopyBuffer(dst, NewReader(src, limit), buf)
    r := NewReader(src, limit)
    rateChecker := time.NewTicker(time.Second * 1)
L1:
    for{
        select{
        case <- rateChecker.C:
            fmt.Printf("time slot %d : rate is %d\n", timeslot, counter)
            timeslot++
            counter = 0
        default:
            if n, err := r.Read(buf); err == nil {
                dst.Write(buf[0:n])
                if n != 0{
                    counter += n
                }
            }else{
                break L1
            }
        }
    }

    fmt.Printf("Copied %d bytes in %s\n", dst.Len(), time.Since(start))
}
