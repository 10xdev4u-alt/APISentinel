package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// AuditEvent represents a single blocked security event.
type AuditEvent struct {
	Timestamp     time.Time `json:"timestamp"`
	SourceIP      string    `json:"source_ip"`
	Method        string    `json:"method"`
	Path          string    `json:"path"`
	ViolationType string    `json:"violation_type"`
	Details       string    `json:"details"`
}

// AuditLogger handles thread-safe writing of audit events to a file.
type AuditLogger struct {
	file *os.File
	mu   sync.Mutex
}

var globalAuditLogger *AuditLogger

// InitAuditLogger initializes the global audit logger.
func InitAuditLogger(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	globalAuditLogger = &AuditLogger{
		file: f,
	}
	log.Printf("📜 Audit logging initialized: %s", filepath)
	return nil
}

// LogEvent writes a structured security event to the audit log.
func LogEvent(ip, method, path, violation, details string) {
	if globalAuditLogger == nil {
		return // Silently skip if not initialized
	}

	event := AuditEvent{
		Timestamp:     time.Now().UTC(),
		SourceIP:      ip,
		Method:        method,
		Path:          path,
		ViolationType: violation,
		Details:       details,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Failed to marshal audit event: %v", err)
		return
	}

	globalAuditLogger.mu.Lock()
	defer globalAuditLogger.mu.Unlock()

	if _, err := globalAuditLogger.file.Write(append(jsonData, '\n')); err != nil {
		log.Printf("❌ Failed to write audit log: %v", err)
	}
}

// Close closes the audit log file.
func Close() {
	if globalAuditLogger != nil && globalAuditLogger.file != nil {
		globalAuditLogger.file.Close()
	}
}
