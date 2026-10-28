import './Auth.css';
import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');

    try {
      const res = await fetch('/api/users/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        throw new Error('Login failed');
      }

      const data = await res.json();
      localStorage.setItem('token', data.token);
      navigate('/home');
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="dx-auth">
      <div className="dx-auth-card">
        <div className="dx-auth-brand">
          <span className="dx-auth-logo">D</span>
          DriveX
        </div>
        <h1 className="dx-auth-title">Log in</h1>
        <p className="dx-auth-sub">Welcome back to your drive.</p>

        {error && (
          <p className="dx-auth-error" role="alert">
            {error}
          </p>
        )}

        <form className="dx-auth-form" onSubmit={handleSubmit}>
          <label className="dx-field">
            Email
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              required
            />
          </label>
          <label className="dx-field">
            Password
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />
          </label>
          <button type="submit" className="dx-auth-btn">
            Log in
          </button>
        </form>

        <p className="dx-auth-switch">
          New here? <Link to="/register">Create an account</Link>
        </p>
      </div>
    </div>
  );
}

export default Login;