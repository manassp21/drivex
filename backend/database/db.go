package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manassp21/drivex/backend/internal/models"
)

var DB *pgxpool.Pool

func CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at",
		email, passwordHash,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}


func CreateFolder(ctx context.Context, name, ownerID string, parentID *string) (*models.Folder, error) {
	folder := &models.Folder{}
	err := DB.QueryRow(ctx,
		"INSERT INTO folders (name, owner_id, parent_id) VALUES ($1, $2, $3) RETURNING id, name, owner_id, parent_id, created_at",
		name, ownerID, parentID,
	).Scan(&folder.ID, &folder.Name, &folder.OwnerID, &folder.ParentID, &folder.CreatedAt)

	if err != nil {
		return nil, err
	}
	return folder, nil
}

func GetFolder(ctx context.Context, id string) (*models.Folder, error) {
	folder := &models.Folder{}
	err := DB.QueryRow(ctx,
		"SELECT id, name, owner_id, parent_id, created_at FROM folders WHERE id = $1 AND deleted_at = NULL",
		id,
	).Scan(&folder.ID, &folder.Name, &folder.OwnerID, &folder.ParentID, &folder.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("folder not found")
		}
		return nil, err
	}
	return folder, nil
}

func ListFolders(ctx context.Context, ownerID string, parentID *string) ([]*models.Folder, error) {
	query := "SELECT id, name, owner_id, parent_id, created_at FROM folders WHERE owner_id = $1 AND deleted_at IS NULL"
	var args []interface{}
	args = append(args, ownerID)

	if parentID != nil {
		query += " AND parent_id = $2"
		args = append(args, *parentID)
	} else {
		query += " AND parent_id IS NULL"
	}

	query += " ORDER BY created_at DESC"

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		folder := &models.Folder{}
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.OwnerID, &folder.ParentID, &folder.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, rows.Err()
}

func DeleteFolder(ctx context.Context, id string) error {
	result, err := DB.Exec(ctx, "UPDATE folders SET deleted_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("folder not found")
	}

	return nil
}


func CreateFile(ctx context.Context, name, ownerID string, folderID *string, sizeBytes int64, mimeType, storageKey string) (*models.File, error) {
	fmt.Printf("createfiles is running")
	file := &models.File{}
	err := DB.QueryRow(ctx,
		"INSERT INTO files (name, owner_id, folder_id, size_bytes, mime_type, storage_key) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, owner_id, folder_id, size_bytes, mime_type, storage_key, created_at, updated_at",
		name, ownerID, folderID, sizeBytes, mimeType, storageKey,
	).Scan(&file.ID, &file.Name, &file.OwnerID, &file.FolderID, &file.SizeBytes, &file.MimeType, &file.StorageKey, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFile(ctx context.Context, id string) (*models.File, error) {
	file := &models.File{}
	err := DB.QueryRow(ctx,
		"SELECT id, name, owner_id, folder_id, size_bytes, mime_type, storage_key, created_at, updated_at FROM files WHERE id = $1 AND deleted_at = NULL",
		id,
	).Scan(&file.ID, &file.Name, &file.OwnerID, &file.FolderID, &file.SizeBytes, &file.MimeType, &file.StorageKey, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return file, nil
}

func ListFiles(ctx context.Context, ownerID string, folderID *string) ([]*models.File, error) {
	query := "SELECT id, name, owner_id, folder_id, size_bytes, mime_type, storage_key, created_at, updated_at FROM files WHERE owner_id = $1 AND deleted_at IS NULL"
	var args []interface{}
	args = append(args, ownerID)

	if folderID != nil {
		query += " AND folder_id = $2"
		args = append(args, *folderID)
	} else {
		query += " AND folder_id IS NULL"
	}

	query += " ORDER BY created_at DESC"

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		file := &models.File{}
		if err := rows.Scan(&file.ID, &file.Name, &file.OwnerID, &file.FolderID, &file.SizeBytes, &file.MimeType, &file.StorageKey, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

func DeleteFile(ctx context.Context, id string) error {
	result, err := DB.Exec(ctx, "DELETE FROM files WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("file not found")
	}

	return nil
}