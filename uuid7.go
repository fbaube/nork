package nork

// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uuid supports generating and manipulating UUIDs.
//
// See [RFC 9562] for details.
//
// Random components of new UUIDs are generated with 
// a cryptographically secure random number generator.
//
// UUIDs may be generated using various algorithms.
// The [New] function returns a new UUID generated 
// using an algorithm suitable for most purposes.
//
// [RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562.html

// This file also works with SQLite RowIDs. Note 
// that altho a rowID is officially 64 bits, here
// it is 30 bits, because we do not expect more
// than one gigabyte of rows in any table.

// This creates a new UUID by buildig it up field
// by field. Another way would be to just call the
// stdlib to create a new one and then set custom
// fields. But this would drain random unnecessarily.
// As written hear, a new UUDI uses MUCH less random
// than the original stdlib version. However that
// approach can still be taken if problems are
// encountered with the approach here. 

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"
    "uuid"
    /*
    "bytes"
    "cmp"
	"encoding/hex"
	"errors"
    */
)

// A UUID is a Universally Unique Identifier as specified in RFC 9562.
//
// UUIDs are comparable, such as with the == operator.
// .
type Uuid7 uuid.UUID // [16]byte

var (
	v7mu            sync.Mutex
	v7lastSecs      uint64
	v7lastTimestamp uint64
)

// NewV7 returns a new version 7 UUID.
//
// Version 7 UUIDs have a timestamp in the most significant 
// 48 bits, and at least 62 bits of random data.
//
// NewV7 always returns UUIDs which sort in increasing 
// order, except when the system clock moves backwards.
//
// Note that the same effect can be achieved by using
// the stdlib to create a new V7 UUID:
//
//     u := uuid.NewV7()
//
// and then calling [SetRowID] and [SetTable]. Note tho
// that this is not the default because it consumes much
// more random than is stricly necessary. 
// 
// Local Modifications
//
// "TBL" is (all or part of) a DB table short name, 
// to speed and simplfy UUID lookup in a database.
// NOTE: To favor UC or LC (or neither) is TBD. 
//
// row_id is copied from SQLite. It might 
// be unnecessary, and replaced by random.
// It is 30 bits, for one gigabyte.
//
// unix_ts_sub_ms (formnerly: rand_a) is 12 bits
// of sub-ms timestamp, 
// as optionally permitted by the RFC.
// rand_b is 8 bits of random. 
// unix_ts_ms is 48 bits i.e. bytes [0..5]
// ver is hi-nibl[6]
// var is hi-bits[8]
//
//   bits 0                   1                   2                   3
//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  bytes +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  0..3  |                           unix_ts_ms                          |
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  4..7  |          unix_ts_ms           |  ver  |    unix_ts_sub_ms     |
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  8..11 |var|                        row_id                             |
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// 12..15 |    rand_b     |      'T'      |      'B'      |      'L'      |
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// .
func NewV7() Uuid7 {
	// UUIDv7 is defined in RFC 9562 section 5.7 as:
	//
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                           unix_ts_ms                          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |          unix_ts_ms           |  ver  |       rand_a          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |var|                        rand_b                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                            rand_b                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	// We store a 12 bit sub-millisecond timestamp fraction in 
	// the rand_a section, as optionally permitted by the RFC.
	v7mu.Lock()

	// 60 bits in our timestamp:
    // 48 bits of millisecond-resolution, followed by
    // 12 bits of 1/4096-millisecond resolution.
    
    // 64 bits: flag + 33(seconds) + 32(nanoseconds) 
    now := time.Now()
    // func (t Time) Unix() int64 ::
    // The number of seconds elapsed
    // since (start of) 1970.01.01 UTC .
    // Currently 1.8 billion (needs 31 bits) 
	secs := uint64(now.Unix())
    // func (t Time) Nanosecond() int :: (ns mod s) :: 
    // The nanosecond offset within the second specified
    // by t, in the range [0, (10**9)-1] (needs 30 bits) 
	nanos := uint64(now.Nanosecond())
    // (ms mod s) :: in the range [0, 999] (10 bits)
	msecs := nanos / 1000000
    // (ns mod ms): The nanosecond offset within the MILLIsecond 
    // specified by t, in the range [0, (10**6)-1] (needs 30 bits)
	fracns   :=  nanos - (1000000 * msecs)
    fracns4k := (fracns * 4096) / 1000000 
    // So now we have:
    // secs (31 bits) , msecs (10 bits) , fracns , fracns4k (12 bits)
    // This result needs (31+10+12) bits = 53 bits 
	timestamp := ((1000*secs + msecs) << 12) + fracns4k 

	if v7lastSecs > secs {
		// Time has gone backwards.
		// This presumably indicates the system clock has changed.
		// Ignore previously-generated UUIDs.
        // TODO: Put a grumpy message on os.Stderr 
	} else if timestamp <= v7lastTimestamp {
		// This timestamp is the same as a previously-generated UUID.
		// To preserve the property that we generate UUIDs in order,
		// use a timestamp 1/4096 millisecond later than the most 
		// recently generated UUID.
        // NOTE: This places a soft limit on sustained thruput of
        // 4000+ UUIDs per second. 
		timestamp = v7lastTimestamp + 1
	}

	v7lastSecs = secs
	v7lastTimestamp = timestamp
	v7mu.Unlock()

	// Insert a gap for the 4 bits of the ver field into the timestamp.
	hibits := ((timestamp << 4) &        // room for one hex digit 
                0xffff_ffff_ffff_0000) | // 12 hex digits = 48 bits 
               (timestamp & 0x0ffff)     //  4 hex digits = 16 bits
    // hibits in hex is now TTTT TTTT TTTT vvvv TTTT

    var u Uuid7
    // var  BigEndian bigEndian
    // type bigEndian struct{}
    // func PutUint64 stores v into b[0:8].
    // func (bigEndian) PutUint64(b []byte, v uint64) {
    // BigEndian is the big-endian implementation
    // of ByteOrder and AppendByteOrder.
    // So this loads 8 bytes i.e. [0..7] 
	binary.BigEndian.PutUint64(u[0:8], hibits)
    // rand.Read generates len(p) random bytes from 
    // the default [Source] and writes them into p. 
    // It always returns len(p) and a nil error.
    // It is safe for concurrent use.
    // func Read(p []byte) (n int, err error) {
    //      return globalRand().Read(p) }
	// rand.Read(u[8:])
    // Now instead of 8 bytes of random, we load only one.
	rand.Read(u[8:8])
	u.setVersion(7)
	u.setVariant(0b10)

	return u
}

// setVersion zeroes the high hex digit and then
// sets it to the version (left-shifted ofc). 
func (u *Uuid7) setVersion(version byte) {
	u[6] = (u[6] & 0b0000_1111) | (version << 4)
}

// setVariant zeroes the high two bits and then
// sets themn to the variant (left-shifted ofc).
func (u *Uuid7) setVariant(variant byte) {
	u[8] = (u[8] & 0b0011_1111) | (variant << 6)
}

// RowID
func (u *Uuid7) RowID() uint32 {
     var rowid uint32
     // middle of byte slice -> uint32
     // y := binary.BigEndian.Uint32(buf[offset : offset+4])
     rowid = binary.BigEndian.Uint32(u[8:16])
     rowid &= 0x3FFF
     return rowid
}

// SetRowID
func (u *Uuid7) SetRowID(rowid uint32) {
     var variant byte
     variant = u[8] & 0b1100_0000
     // uint32 -> middle of byte slice
     // binary.BigEndian.PutUint32(buf[offset:offset+4], x)
     binary.BigEndian.PutUint32(u[8:12], rowid)
     u[8] |= variant
}

// Table
func (u *Uuid7) Table() string {
     return string(u[13:16])
}

// SetTable
func (u *Uuid7) SetTable(table string) {
     // table += "   "
     // copy(dst []byte, src string) int
     copy(u[13:16], (table+"   ")[0:3])
}

/*
// String
func (u Uuid7) String() string {
     // First convert the ms count to a Unix time.
     // Shift the elements into place
     var unix_ts_ms uint64
     unix_ts_ms = (u[] << n) | (u[] << n) | (u[] << n) | (u[] << n) | (u[] << n) | (u[] << n) | 
     
}
*/

// =========================================================

/*

// Parse returns the UUID represented by s.
//
// It accepts strings in the following forms:
//
//	f81d4fae-7dec-11d0-a765-00a0c91e6bf6
//	{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}
//	urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6
//	f81d4fae7dec11d0a76500a0c91e6bf6
//
// Alphabetic characters in the hexadecimal digits may be any case.
func Parse(s string) (Uuid7, error) {
	var u Uuid7
	err := u.UnmarshalText([]byte(s))
	return u, err
}

// MustParse returns the UUID represented by s.
//
// It panics if s is not a valid string representation of a UUID as defined by [Parse].
func MustParse(s string) UUID {
	u, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// New returns a new UUID.
//
// Programs which do not have a need for a specific UUID generation algorithm should use New.
// At this time, New is equivalent to [NewV4].
func New() UUID {
	return NewV4()
}

// Nil returns the Nil UUID 00000000-0000-0000-0000-000000000000.
//
// The Nil UUID is defined in [Section 5.9 of RFC 9562].
// Note that this is not the same as the Go value nil.
//
// [Section 5.9 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-5.9
func Nil() UUID {
	return UUID{}
}

// Max returns the Max UUID ffffffff-ffff-ffff-ffff-ffffffffffff.
//
// The Max UUID is defined in [Section 5.10 of RFC 9562].
//
// [Section 5.10 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-5.10
func Max() UUID {
	return UUID{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
}

// String returns the string representation of u.
//
// It uses the lowercase hex-and-dash representation defined in RFC 9562.
func (u UUID) String() string {
	b, _ := u.MarshalText()
	return string(b)
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// The encoding is the same as returned by [UUID.String]
func (u UUID) MarshalText() ([]byte, error) {
	return u.AppendText(make([]byte, 0, 36))
}

// AppendText implements the [encoding.TextAppender] interface.
// The encoding is the same as returned by [UUID.String]
func (u UUID) AppendText(b []byte) ([]byte, error) {
	off := len(b)
	b = append(b, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"...)
	dst := b[off:]
	hex.Encode(dst[0:8], u[0:4])
	hex.Encode(dst[9:13], u[4:6])
	hex.Encode(dst[14:18], u[6:8])
	hex.Encode(dst[19:23], u[8:10])
	hex.Encode(dst[24:36], u[10:16])
	return b, nil
}

var errInvalid = errors.New("invalid uuid")

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// The UUID is expected in a form accepted by [Parse].
func (u *UUID) UnmarshalText(b []byte) error {
	var dst UUID
	switch len(b) {
	case len("urn:uuid:") + 36:
		// urn:uuid:00000000-0000-0000-0000-000000000000
		b = bytes.TrimPrefix(b, []byte("urn:uuid:"))
	case 2 + 36:
		// {00000000-0000-0000-0000-000000000000}
		b = bytes.TrimPrefix(b, []byte("{"))
		b = bytes.TrimSuffix(b, []byte("}"))
	case 32:
		// 00000000000000000000000000000000
		if _, err := hex.Decode(dst[:], b); err != nil {
			return errInvalid
		}
		*u = dst
		return nil
	}
	// 00000000-0000-0000-0000-000000000000
	if len(b) != 36 {
		return errInvalid
	}
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return errInvalid
	}
	if _, err := hex.Decode(dst[0:4], b[0:8]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(dst[4:6], b[9:13]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(dst[6:8], b[14:18]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(dst[8:10], b[19:23]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(dst[10:16], b[24:36]); err != nil {
		return errInvalid
	}
	*u = dst
	return nil
}

// Compare compares the UUID u with v.
// If u is before v, it returns -1.
// If u is after v, it returns +1.
// If they are the same, it returns 0.
//
// Compare uses the big-endian byte order defined in
// [Section 6.11 of RFC 9562] for sorting.
//
// [Section 6.11 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-6.11
func (u UUID) Compare(v UUID) int {
	for i := range u {
		if c := cmp.Compare(u[i], v[i]); c != 0 {
			return c
		}
	}
	return 0
}

// NewV4 returns a new version 4 UUID.
//
// Version 4 UUIDs contain 122 bits of random data.
func NewV4() UUID {
	var u UUID
	rand.Read(u[:])
	u.setVersion(4)
	u.setVariant(0b10)
	return u
}

*/

