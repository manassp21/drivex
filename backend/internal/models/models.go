
package models
 
import "time"
 
type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
 
type Folder struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	OwnerID   string    `json:"owner_id" db:"owner_id"`
	ParentID  *string   `json:"parent_id,omitempty" db:"parent_id"` // nil = root folder
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
 
type File struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	OwnerID    string    `json:"owner_id" db:"owner_id"`
	FolderID   *string   `json:"folder_id,omitempty" db:"folder_id"` // nil = root
	SizeBytes  int64     `json:"size_bytes" db:"size_bytes"`
	MimeType   string    `json:"mime_type" db:"mime_type"`
	StorageKey string    `json:"-" db:"storage_key"` // path/key in your blob storage (disk, S3, etc.)
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
 