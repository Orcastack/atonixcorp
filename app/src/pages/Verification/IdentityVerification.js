import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const IdentityVerification = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [verification, setVerification] = useState(null);
  const [idDocument, setIdDocument] = useState(null);
  const [selfie, setSelfie] = useState(null);
  const [message, setMessage] = useState('');
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    fetch('/api/auth/identity-verification/', { headers: { Authorization: `Bearer ${localStorage.getItem('token')}` } })
      .then(async (response) => ({ ok: response.ok, data: await response.json() }))
      .then(({ ok, data }) => {
        if (!ok) throw new Error(data?.error?.message || 'Unable to load identity verification.');
        setVerification(data);
        if (data.status === 'verified') navigate('/app/console', { replace: true });
      })
      .catch((error) => setMessage(error.message));
  }, [navigate]);

  const submit = async (event) => {
    event.preventDefault();
    if (!user?.email_verified) {
      setMessage('Please verify your email first.');
      return;
    }
    if (!idDocument || !selfie) {
      setMessage('Upload both an ID document and a selfie.');
      return;
    }
    setBusy(true);
    setMessage('');
    const data = new FormData();
    data.append('id_document', idDocument);
    data.append('selfie', selfie);
    try {
      const response = await fetch('/api/auth/identity-verification/', {
        method: 'POST',
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        body: data,
      });
      const payload = await response.json();
      if (!response.ok) throw new Error(payload?.error?.message || 'Unable to submit identity verification.');
      setVerification(payload);
      setMessage('Identity documents submitted for review.');
    } catch (error) {
      setMessage(error.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="auth-page verification-page">
      <div className="auth-container">
        <div className="auth-card">
          <h1>Identity Verification</h1>
          <p className="auth-subtitle">Email verified. Upload your government-issued ID and a current selfie to complete verification.</p>
          {message && <p className="auth-error" role="alert">{message}</p>}
          {verification?.status === 'submitted' ? <p className="auth-subtitle">Your documents are under review.</p> : (
            <form className="auth-form" onSubmit={submit}>
              <div className="form-group"><label htmlFor="id-document">Government-issued ID</label><input id="id-document" type="file" accept="image/*,.pdf" onChange={(event) => setIdDocument(event.target.files?.[0] || null)} required /></div>
              <div className="form-group"><label htmlFor="selfie">Selfie</label><input id="selfie" type="file" accept="image/*" capture="user" onChange={(event) => setSelfie(event.target.files?.[0] || null)} required /></div>
              <button className="btn-primary btn-full" disabled={busy}>{busy ? 'Submitting...' : 'Submit for review'}</button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};

export default IdentityVerification;
