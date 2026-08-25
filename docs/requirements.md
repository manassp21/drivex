## Authentication

### Signup
- Fields: name, email, password
- Password rules: minimum 8 characters, must include at least one number, 
  one letter, and one special character
- On successful signup, user is automatically logged in (no separate login step)

### Login
- Fields: email, password
- Returns: access token + refresh token

### Tokens
- Access token: short-lived (15 min), sent with every API request, used to 
  authenticate the user
- Refresh token: long-lived (7 days), used only to obtain a new access token 
  when it expires
- A user can be logged in on multiple devices simultaneously — each login 
  session gets its own refresh token
- Refresh tokens are stored in the database (one row per session/device), 
  not just trusted blindly — this lets us revoke a single session later

### Logout
- Invalidates the refresh token for that specific session/device
- Access token isn't stored server-side, so it just naturally expires

## Files

### Upload
- Accepts all file types
- Maximum file size: 500MB per file
- Files are private to the owner by default

### Download
- File downloads with its original filename

### Rename
- A file cannot be renamed to a name that already exists 
  in the same folder

### Delete
- Files are moved to Trash on delete
- Trash items are permanently deleted after 30 days
- Users can manually empty trash or restore files within 
  the 30 day window

### Ownership & Privacy
- Files are strictly private to their owner
- No other user can access, view, or modify another 
  user's files without explicit sharing (covered in Phase 13)	

## Folders

### Structure
- Folders can be nested to any depth (subfolders inside subfolders)
- Every user gets a default root folder ("My Drive") automatically 
  created on signup

### Naming
- Duplicate folder names are not allowed within the same parent folder

### Delete
- Deleting a folder moves the folder and everything inside it 
  (files and subfolders) to Trash
- Trash rules apply: permanent deletion after 30 days

### Move
- Folders can be moved into other folders
- Files can also be moved between folders

## Sharing

### Public Links
- A user can generate a public link for any file
- Anyone with the link can download the file without logging in
- Links stay valid forever until the owner manually revokes them
- Owner can revoke (delete) a shared link at any time

## Search

### Scope
- Search works on both filenames and folder names
- Two modes:
  - Current folder search: searches only within the folder the user 
    is currently in
  - Global search: searches across all of the user's files and folders, 
    accessible via a separate search bar

## Dashboard

### Stats Shown
- Total storage used
- Total file count
- Total folder count
- Recently uploaded/modified files

## Settings

### Options
- Change name
- Change password (requires current password confirmation)
- Delete account (permanently deletes account, all files, 
  folders, and shared links after confirmation)