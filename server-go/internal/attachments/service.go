package attachments

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var AllowedCategories = map[string]string{
	".pdf":  "pdf",
	".doc":  "doc",
	".docx": "doc",
	".ppt":  "presentation",
	".pptx": "presentation",
	".png":  "image",
	".jpg":  "image",
	".jpeg": "image",
	".mp4":  "video",
}

type Attachment struct {
	ID           string `json:"id"`
	MessageID    string `json:"message_id,omitempty"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Category     string `json:"file_category"`
	FileSize     int64  `json:"file_size"`
	StoragePath  string `json:"storage_path"`
}

func SaveFile(ctx context.Context, db *pgxpool.Pool, uploadDir, uploaderID, messageID, originalName, mimeType string, fileSize int64, r io.Reader) (*Attachment, error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	cat, ok := AllowedCategories[ext]
	if !ok {
		return nil, fmt.Errorf("format de fichier non autorisé: %s", ext)
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, err
	}
	storedName := uuid.NewString() + ext
	dstPath := filepath.Join(uploadDir, storedName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(dst, hasher), r); err != nil {
		return nil, err
	}
	checksum := hex.EncodeToString(hasher.Sum(nil))

	var att Attachment
	err = db.QueryRow(ctx, `
		INSERT INTO attachments (message_id, uploader_id, file_name, original_name, mime_type, file_category, file_size, storage_path, checksum_sha256)
		VALUES (NULLIF($1, '')::uuid, $2::uuid, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text, COALESCE(message_id::text, ''), file_name, original_name, mime_type, file_category, file_size, storage_path
	`, messageID, uploaderID, storedName, originalName, mimeType, cat, fileSize, "/uploads/"+storedName, checksum).Scan(
		&att.ID, &att.MessageID, &att.FileName, &att.OriginalName, &att.MimeType, &att.Category, &att.FileSize, &att.StoragePath,
	)
	return &att, err
}
