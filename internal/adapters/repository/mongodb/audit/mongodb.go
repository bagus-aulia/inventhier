package audit

import (
	"context"
	"time"

	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAuditLogger struct {
	db *mongo.Database
}

type AuditLog struct {
	Action    string    `bson:"action"`
	ProductID string    `bson:"product_id"`
	Details   string    `bson:"details"`
	Timestamp time.Time `bson:"timestamp"`
}

// NewMongoAuditLogger creates a new MongoDB-based audit logger.
func NewMongoAuditLogger(db *mongo.Database) ports.ProductAuditLogger {
	return &mongoAuditLogger{
		db: db,
	}
}

func (m *mongoAuditLogger) LogAction(ctx context.Context, action string, productID string, details string) error {
	collection := m.db.Collection("audit_logs")
	logEntry := AuditLog{
		Action:    action,
		ProductID: productID,
		Details:   details,
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(ctx, logEntry)
	return err
}
