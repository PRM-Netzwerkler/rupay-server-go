package util

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func PtrToNullString(ptr *string) sql.NullString {
    if ptr != nil {
        return sql.NullString{
            String: *ptr,
            Valid:  true,
        }
    }
    return sql.NullString{Valid: false}
}

func PtrToNullInt64(ptr *int64) sql.NullInt64 {
    if ptr != nil {
        return sql.NullInt64{
            Int64: *ptr,
            Valid: true,
        }
    }
    return sql.NullInt64{Valid: false}
}

func PtrToNullInt32(ptr *int32) sql.NullInt32 {
    if ptr != nil {
        return sql.NullInt32{
            Int32: *ptr,
            Valid: true,
        }
    }
    return sql.NullInt32{Valid: false}
}

func PtrToNullFloat64(ptr *float64) sql.NullFloat64 {
    if ptr != nil {
        return sql.NullFloat64{
            Float64: *ptr,
            Valid:   true,
        }
    }
    return sql.NullFloat64{Valid: false}
}

func PtrToNullUUID(ptr *uuid.UUID) uuid.NullUUID {
    if ptr != nil {
        return uuid.NullUUID{
            UUID: *ptr,
            Valid:   true,
        }
    }
    return uuid.NullUUID{Valid: false}
}

func PtrToNullTime(ptr *time.Time) sql.NullTime {
    if ptr != nil {
        return sql.NullTime{
            Time: *ptr,
            Valid:   true,
        }
    }
	return sql.NullTime{Valid: false}
}
