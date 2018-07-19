package utils

import (
	"io"
	"os"
	"unicode/utf8"
)

// \n
func ReadStdin() string {
	r := &readRune{reader: os.Stdin, peekRune: -1}
	var atEOF bool
	rrs := make([]rune, 0, 64)
	for {
		if atEOF {
			break
		}
		rr, _, err := r.ReadRune()
		if err == nil {
			if rr == '\n' {
				atEOF = true
				continue
			}
			rrs = append(rrs, rr)
		} else if err == io.EOF {
			atEOF = true
		} else {
			atEOF = true
			// todo
		}
	}
	return string(rrs)
}

// copy from fmt.pkg
type readRune struct {
	reader   io.Reader
	buf      [utf8.UTFMax]byte // used only inside ReadRune
	pending  int               // number of bytes in pendBuf; only >0 for bad UTF-8
	pendBuf  [utf8.UTFMax]byte // bytes left over
	peekRune rune              // if >=0 next rune; when <0 is ^(previous Rune)
}

// readByte returns the next byte from the input, which may be
// left over from a previous read if the UTF-8 was ill-formed.
func (r *readRune) readByte() (b byte, err error) {
	if r.pending > 0 {
		b = r.pendBuf[0]
		copy(r.pendBuf[0:], r.pendBuf[1:])
		r.pending--
		return
	}
	n, err := io.ReadFull(r.reader, r.pendBuf[:1])
	if n != 1 {
		return 0, err
	}
	return r.pendBuf[0], err
}

// ReadRune returns the next UTF-8 encoded code point from the
// io.Reader inside r.
func (r *readRune) ReadRune() (rr rune, size int, err error) {
	if r.peekRune >= 0 {
		rr = r.peekRune
		r.peekRune = ^r.peekRune
		size = utf8.RuneLen(rr)
		return
	}
	r.buf[0], err = r.readByte()
	if err != nil {
		return
	}
	if r.buf[0] < utf8.RuneSelf { // fast check for common ASCII case
		rr = rune(r.buf[0])
		size = 1 // Known to be 1.
		// Flip the bits of the rune so it's available to UnreadRune.
		r.peekRune = ^rr
		return
	}
	var n int
	for n = 1; !utf8.FullRune(r.buf[:n]); n++ {
		r.buf[n], err = r.readByte()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
	}
	rr, size = utf8.DecodeRune(r.buf[:n])
	if size < n { // an error, save the bytes for the next read
		copy(r.pendBuf[r.pending:], r.buf[size:n])
		r.pending += n - size
	}
	// Flip the bits of the rune so it's available to UnreadRune.
	r.peekRune = ^rr
	return
}
