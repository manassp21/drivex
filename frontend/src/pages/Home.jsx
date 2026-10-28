import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './Home.css';

/* ---------- helpers ---------- */

function getFileKind(mime = '') {
  if (mime.startsWith('image/')) return { label: 'Image', cls: 'image' };
  if (mime.startsWith('video/')) return { label: 'Video', cls: 'video' };
  if (mime.startsWith('audio/')) return { label: 'Audio', cls: 'audio' };
  if (mime === 'application/pdf') return { label: 'PDF', cls: 'pdf' };
  if (mime.startsWith('text/')) return { label: 'Text', cls: 'text' };
  if (mime.includes('zip') || mime.includes('compressed')) return { label: 'Archive', cls: 'archive' };
  return { label: 'File', cls: 'other' };
}

function formatSize(bytes) {
  if (bytes == null) return '';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
}

function formatDate(iso) {
  if (!iso) return '';
  return new Date(iso).toLocaleDateString(undefined, {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

/* ---------- icons ---------- */

const svgProps = {
  viewBox: '0 0 24 24',
  fill: 'none',
  stroke: 'currentColor',
  strokeWidth: 2,
  strokeLinecap: 'round',
  strokeLinejoin: 'round',
  'aria-hidden': true,
};

const FolderIcon = () => (
  <svg viewBox="0 0 24 24" width="26" height="26" aria-hidden="true">
    <path
      d="M3 6.5A1.5 1.5 0 0 1 4.5 5h4.38a1.5 1.5 0 0 1 1.06.44L11.5 7h8A1.5 1.5 0 0 1 21 8.5v9a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 17.5v-11Z"
      fill="currentColor"
    />
  </svg>
);

const FileGlyph = () => (
  <svg {...svgProps} width="36" height="36">
    <path d="M7 3h7l5 5v12a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1Z" />
    <path d="M14 3v5h5" />
  </svg>
);

const UploadIcon = () => (
  <svg {...svgProps} width="16" height="16">
    <path d="M12 16V6m0 0-4 4m4-4 4 4M5 19h14" />
  </svg>
);

const DownloadIcon = () => (
  <svg {...svgProps} width="15" height="15">
    <path d="M12 4v10m0 0-4-4m4 4 4-4M5 19h14" />
  </svg>
);

const TrashIcon = () => (
  <svg {...svgProps} width="17" height="17">
    <path d="M4 7h16M10 11v6M14 11v6M6 7l1 12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1l1-12M9 7V4h6v3" />
  </svg>
);

/* ---------- component ---------- */

function Home() {
  const { folderID } = useParams();
  const [folders, setFolders] = useState([]);
  const [files, setFiles] = useState([]);
  const [currentFolder, setCurrentFolder] = useState(null);
  const [error, setError] = useState('');
  const [newFolderName, setNewFolderName] = useState('');
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const navigate = useNavigate();

  const token = localStorage.getItem('token');

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }
    loadData();
  }, [folderID]);

  async function loadData() {
    setLoading(true);
    try {
      const folderQuery = folderID ? `?parent_id=${folderID}` : '';
      const fileQuery = folderID ? `?folder_id=${folderID}` : '';
      const headers = { Authorization: `Bearer ${token}` };

      const [foldersRes, filesRes, currentRes] = await Promise.all([
        fetch(`/api/folders${folderQuery}`, { headers }),
        fetch(`/api/files${fileQuery}`, { headers }),
        folderID ? fetch(`/api/folders/${folderID}`, { headers }) : null,
      ]);

      if (foldersRes.status === 401 || filesRes.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
        return;
      }

      if (!foldersRes.ok) {
        const text = await foldersRes.text();
        throw new Error(`Folders: ${text || foldersRes.status}`);
      }
      if (!filesRes.ok) {
        const text = await filesRes.text();
        throw new Error(`Files: ${text || filesRes.status}`);
      }

      setFolders((await foldersRes.json()) || []);
      setFiles((await filesRes.json()) || []);
      setCurrentFolder(currentRes && currentRes.ok ? await currentRes.json() : null);
      setError('');
    } catch (err) {
      console.error('loadData error:', err);
      setError(err.message || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateFolder(e) {
    e.preventDefault();
    if (!newFolderName.trim()) return;

    try {
      const res = await fetch('/api/folders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: newFolderName.trim(), parent_id: folderID || null }),
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Failed to create folder');
      }
      setNewFolderName('');
      loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleDeleteFolder(id, name) {
    if (!window.confirm(`Delete the folder "${name}"?`)) return;
    try {
      const res = await fetch(`/api/folders/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Failed to delete folder');
      }
      loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleUploadFile(e) {
    const input = e.target;
    const file = input.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);
    if (folderID) formData.append('folder_id', folderID);

    setUploading(true);
    try {
      const res = await fetch('/api/files', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      if (!res.ok) throw new Error('Upload failed. Files must be under 25 MB.');
      loadData();
    } catch (err) {
      setError(err.message);
    } finally {
      input.value = '';
      setUploading(false);
    }
  }

  async function handleDeleteFile(id, name) {
    if (!window.confirm(`Delete "${name}"?`)) return;
    try {
      const res = await fetch(`/api/files/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error('Failed to delete file');
      loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleDownloadFile(id, filename) {
    try {
      const res = await fetch(`/api/files/${id}/download`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error('Download failed');

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError(err.message);
    }
  }

  function handleLogout() {
    localStorage.removeItem('token');
    navigate('/login');
  }

  function goUp() {
    const parent = currentFolder?.parent_id;
    navigate(parent ? `/home/${parent}` : '/home');
  }

  const isEmpty = !loading && folders.length === 0 && files.length === 0;

  return (
    <div className="dx-page">
      <header className="dx-topbar">
        <div className="dx-brand">
          <span className="dx-logo">D</span>
          DriveX
        </div>
        <button className="dx-btn dx-btn-ghost" onClick={handleLogout}>
          Log out
        </button>
      </header>

      <main className="dx-main">
        <div className="dx-toolbar">
          <nav className="dx-crumbs" aria-label="Current location">
            {folderID && (
              <button className="dx-icon-btn" onClick={goUp} aria-label="Go up one folder">
                ←
              </button>
            )}
            <button className="dx-crumb" onClick={() => navigate('/home')}>
              My Drive
            </button>
            {folderID && (
              <>
                <span className="dx-sep">/</span>
                <span className="dx-crumb-current">{currentFolder?.name || 'Folder'}</span>
              </>
            )}
          </nav>

          <div className="dx-actions">
            <form className="dx-newfolder" onSubmit={handleCreateFolder}>
              <input
                type="text"
                placeholder="Folder name"
                value={newFolderName}
                onChange={(e) => setNewFolderName(e.target.value)}
                aria-label="New folder name"
              />
              <button type="submit" className="dx-btn dx-btn-secondary">
                New folder
              </button>
            </form>
            <label className={`dx-btn dx-btn-primary${uploading ? ' is-busy' : ''}`}>
              <UploadIcon />
              {uploading ? 'Uploading…' : 'Upload file'}
              <input
                type="file"
                className="dx-visually-hidden"
                onChange={handleUploadFile}
                disabled={uploading}
              />
            </label>
          </div>
        </div>

        {error && (
          <div className="dx-error" role="alert">
            <span>{error}</span>
            <button onClick={() => setError('')} aria-label="Dismiss error">
              ×
            </button>
          </div>
        )}

        {loading && folders.length === 0 && files.length === 0 && (
          <p className="dx-muted">Loading…</p>
        )}

        {isEmpty && (
          <div className="dx-empty">
            <h2>This folder is empty</h2>
            <p>Upload a file or create a folder to get started.</p>
          </div>
        )}

        {folders.length > 0 && (
          <section>
            <h2 className="dx-section-title">Folders</h2>
            <div className="dx-grid dx-grid-folders">
              {folders.map((folder) => (
                <div className="dx-folder" key={folder.id}>
                  <button
                    className="dx-folder-open"
                    onClick={() => navigate(`/home/${folder.id}`)}
                  >
                    <span className="dx-folder-icon">
                      <FolderIcon />
                    </span>
                    <span className="dx-name" title={folder.name}>
                      {folder.name}
                    </span>
                  </button>
                  <button
                    className="dx-icon-btn dx-danger"
                    onClick={() => handleDeleteFolder(folder.id, folder.name)}
                    aria-label={`Delete folder ${folder.name}`}
                  >
                    <TrashIcon />
                  </button>
                </div>
              ))}
            </div>
          </section>
        )}

        {files.length > 0 && (
          <section className={folders.length > 0 ? 'dx-section-gap' : ''}>
            <h2 className="dx-section-title">Files</h2>
            <div className="dx-grid dx-grid-files">
              {files.map((file) => {
                const kind = getFileKind(file.mime_type);
                return (
                  <article className="dx-file" key={file.id}>
                    <div className={`dx-file-preview kind-${kind.cls}`}>
                      <FileGlyph />
                      <span className="dx-kind">{kind.label}</span>
                    </div>
                    <div className="dx-file-body">
                      <div className="dx-name" title={file.name}>
                        {file.name}
                      </div>
                      <div className="dx-meta">
                        <span>{formatSize(file.size_bytes)}</span>
                        <span>{formatDate(file.created_at)}</span>
                      </div>
                    </div>
                    <div className="dx-file-actions">
                      <button
                        className="dx-btn dx-btn-small"
                        onClick={() => handleDownloadFile(file.id, file.name)}
                      >
                        <DownloadIcon />
                        Download
                      </button>
                      <button
                        className="dx-icon-btn dx-danger"
                        onClick={() => handleDeleteFile(file.id, file.name)}
                        aria-label={`Delete ${file.name}`}
                      >
                        <TrashIcon />
                      </button>
                    </div>
                  </article>
                );
              })}
            </div>
          </section>
        )}
      </main>
    </div>
  );
}

export default Home;