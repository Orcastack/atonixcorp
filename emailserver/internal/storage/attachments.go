package storage

import (
	"context"
	"fmt"

	"atonixcorp/emailserver/internal/blob"
)

type Attachment struct {
	ID        string
	MessageID string
	UserID    string
	Filename  string
	Size      int
}

type AttachmentStore struct {
	blob blob.Store
}

func NewAttachmentStore(blobStore blob.Store) *AttachmentStore {
	return &AttachmentStore{blob: blobStore}
}

func (s *AttachmentStore) SaveAttachment(ctx context.Context, userID, msgID, attachID string, data []byte) error {
	key := fmt.Sprintf("users/%s/messages/%s/attachments/%s", userID, msgID, attachID)
	return s.blob.Put(ctx, key, data)
}

func (s *AttachmentStore) GetAttachment(ctx context.Context, userID, msgID, attachID string) ([]byte, error) {
	key := fmt.Sprintf("users/%s/messages/%s/attachments/%s", userID, msgID, attachID)
	return s.blob.Get(ctx, key)
}

func (s *AttachmentStore) DeleteAttachment(ctx context.Context, userID, msgID, attachID string) error {
	key := fmt.Sprintf("users/%s/messages/%s/attachments/%s", userID, msgID, attachID)
	return s.blob.Delete(ctx, key)
}
