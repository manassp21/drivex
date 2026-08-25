import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function Home() {
  const [folders, setFolders] = useState([]);
  const [files, setFiles] = useState([]);
  const [error, setError] = useState('');
  const [newFolderName, setNewFolderName] = useState('');
  const navigate = useNavigate();

  const token = localStorage.getItem('token');

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }
    loadData();
  }, []);

  async function loadData() {
    try {
      const [foldersRes, filesRes] = await Promise.all([
        fetch('/api/folders', { headers: { Authorization: `Bearer ${token}` } }),
        fetch('/api/files', { headers: { Authorization: `Bearer ${token}` } }),
      ]);

      if (foldersRes.status === 401 || filesRes.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
        return;
      }

      const foldersData = await foldersRes.json();
      const filesData = await filesRes.json();
      setFolders(foldersData || []);
      setFiles(filesData || []);
    } catch (err) {
      setError('Failed to load data');
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
        body: JSON.stringify({ name: newFolderName }),
      });
      if (!res.ok) throw new Error('Failed to create folder');
      setNewFolderName('');
      loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleDeleteFolder(id) {
    try {
      await fetch(`/api/folders/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      loadData();
    } catch (err) {
      setError('Failed to delete folder');
    }
  }

  async function handleUploadFile(e) {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await fetch('/api/files', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      if (!res.ok) throw new Error('Upload failed');
      loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleDeleteFile(id) {
    try {
      await fetch(`/api/files/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      loadData();
    } catch (err) {
      setError('Failed to delete file');
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

  return (
    <div style={{ maxWidth: 700, margin: '40px auto', fontFamily: 'sans-serif' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>DriveX</h1>
        <button onClick={handleLogout}>Log Out</button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <section style={{ marginTop: 30 }}>
        <h2>Folders</h2>
        <form onSubmit={handleCreateFolder} style={{ marginBottom: 10 }}>
          <input
            type="text"
            placeholder="New folder name"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
          />
          <button type="submit">Create Folder</button>
        </form>
        <ul>
          {folders.map((folder) => (
            <li key={folder.id}>
              {folder.name}{' '}
              <button onClick={() => handleDeleteFolder(folder.id)}>Delete</button>
            </li>
          ))}
        </ul>
      </section>

      <section style={{ marginTop: 30 }}>
        <h2>Files</h2>
        <input type="file" onChange={handleUploadFile} />
        <ul>
          {files.map((file) => (
            <li key={file.id}>
              {file.name} ({(file.size_bytes / 1024).toFixed(1)} KB){' '}
              <button onClick={() => handleDownloadFile(file.id, file.name)}>Download</button>{' '}
              <button onClick={() => handleDeleteFile(file.id)}>Delete</button>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default Home;