users
- id
- name
- email
- password_hash
- created_at
- updated_at

refresh_tokens
- id
- user_id (foreign key → users.id)
- token
- expires_at
- created_at

folders
- id
- user_id (foreign key → users.id)
- parent_id (foreign key → folders.id, NULL for root folder)
- name
- created_at
- updated_at

files
- id
- user_id (foreign key → users.id)
- folder_id (foreign key → folders.id)
- name
- path
- size
- mime_type
- created_at
- updated_at

trash
- id
- user_id (foreign key → users.id)
- file_id (foreign key → files.id, NULL if folder)
- folder_id (foreign key → folders.id, NULL if file)
- deleted_at
- expires_at (deleted_at + 30 days)

shared_links
- id
- user_id (foreign key → users.id)
- file_id (foreign key → files.id)
- token (unique random string, e.g. abc123xyz)
- created_at