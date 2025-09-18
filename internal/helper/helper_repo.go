package helper


import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

// Helper functions for UUID conversion

func ConvertStringToUUID(s string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	
	if s == "" {
		return pgUUID, fmt.Errorf("UUID string cannot be empty")
	}
	
	// Parse the string as UUID first to validate it
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return pgUUID, fmt.Errorf("invalid UUID format: %w", err)
	}
	
	// Convert to pgtype.UUID
	pgUUID.Bytes = parsedUUID
	pgUUID.Valid = true
	
	return pgUUID, nil
}

func ConvertUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}
	
	// Convert pgtype.UUID to uuid.UUID and then to string
	u := uuid.UUID(pgUUID.Bytes)
	return u.String()
}