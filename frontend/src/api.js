export function listFolders(parentId) {
  const query = parentId ? `?parent_id=${parentId}` : '';
  return apiFetch(`/api/folders${query}`);
}

export function listFiles(parentId) {
  const query = parentId ? `?parent_id=${parentId}` : '';
  return apiFetch(`/api/files${query}`);
}