package nork

import (
	"encoding/binary"
	"time"
    "uuid"
    "fmt"
	"errors"
    SU "github.com/fbaube/stringutils"
)

func (p *Nork)  UUID() uuid.UUID  { return p.UUIDv7 }
func (p *Nork) TableName() string { return string(p.UUIDv7[13:16]) }
func TableNameOfUUID(u uuid.UUID) string {
     return string(u[13:16])
}

// TableNameInUUID returns the 3-letter DB table code.
// This is non-standard for UUIDs.
func TableNameInUUID(u *uuid.UUID) string {
     return string(u[13:16])
}

// SetTableNameInUUID sets the 3-letter DB table code.
// This	is non-standard for UUIDs. 
func SetTableNameInUUID(u *uuid.UUID, table string) {
     // copy(dst []byte, src string) int
     copy(u[13:16], (table+"   ")[0:3])
}

// String assumes v7, and table name in last three bytes.
// func (u uuid.UUID) String() string {
func StringOfUUID(u *uuid.UUID) string {
     // First convert the ms count to a Unix time.
     // Shift the elements into place
     var unix_ts_ms uint64
     unix_ts_ms = uint64(binary.BigEndian.Uint64(u[:6]) >> 16)
     // timestamp to YMDHM
     // sep. ts bits to sthg else
     // then tbl & rowID
     var unix_time time.Time
     unix_time = time.UnixMilli(int64(unix_ts_ms))
//   println(unix_time.String())
     var ymdhm string
     ymdhm = SU.TimeAsYMDHM(unix_time)
     return ymdhm + fmt.Sprintf("-%1x%2x-%c%c%c", // -%4d",
            u[6]&0xF, u[7], u[13], u[14], u[15]) // , u.RowID())
}

// =================================
//  These next func's were proposed
//  (most as methods) by Perplexity 
// =================================

func UUIDisValidKey(u uuid.UUID) error {
    if u == NilUUID {
        return errors.New("required UUID is Nil (all zeroes)")
    }
    return nil
}
func NullableUUIDisSet(u UUIDnullable) bool {
    return u.IsValid && u.UUID != NilUUID
}
func UUIDasBytes(u any) []byte {
    switch v := u.(type) {
    case uuid.UUID:
        return v[:]
    case UUIDnullable:
        if !v.IsValid {
            return nil
        }
        return v.UUID[:]
    // TODO: Accept a Nork ? 
    default:
        panic(fmt.Sprintf("UUIDasBytes: unsupported type %T", u))
    }
}
func UUIDfromBytes(b []byte) (uuid.UUID, error) { 
     var u uuid.UUID 
     if len(b) != 16 {
        return NilUUID, fmt.Errorf("UUID must be 16 bytes, got %d", len(b))
     }
     // copy(dst []byte, src string) int
     copy(u[:], b)
     return u, nil 
     }
func UUIDisNil(u any) bool {
   switch v := u.(type) {
    case uuid.UUID:
        return v == NilUUID
    case UUIDnullable:
        return !v.IsValid || v.UUID == NilUUID
    default:
        panic(fmt.Sprintf("UUIDisNil: unsupported type %T", u))
    }
}
func UUIDasHexString(u uuid.UUID) string { return "UUIDasHexString:TBS" }
func UUIDasString(u uuid.UUID) string { return "UUIDasString:TBS" }

func UUIDhasValidTableName(u uuid.UUID) bool {
    if u == NilUUID {
        return false // errors.New("required UUID is Nil (all zeroes)")
    }
    got := TableNameInUUID(&u)
    return SU.IsAsciiAlphanum(got) 
}

