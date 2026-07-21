import React, { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import AtonixCorpLogo from '../../components/branding/AtonixCorpLogo';

const EmailVerification = () => {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { completeEmailVerification } = useAuth();
  const [status, setStatus] = useState(searchParams.get('token') ? 'verifying' : 'sent');
  const [message, setMessage] = useState('');

  useEffect(() => {
    const token = searchParams.get('token');
    if (!token) return;

    let active = true;
    completeEmailVerification(token).then((result) => {
      if (!active) return;
      if (result.success) {
        setStatus('verified');
        setMessage('Your email is verified. Continuing to identity verification.');
        window.setTimeout(() => navigate(result.nextPath, { replace: true }), 750);
      } else {
        setStatus('failed');
        setMessage(result.error || 'This verification link is invalid or has expired.');
      }
    });
    return () => { active = false; };
  }, [completeEmailVerification, navigate, searchParams]);

  const email = location.state?.email;
  return (
    <div className="auth-page">
      <div className="auth-container">
        <div className="auth-header"><Link to="/" className="auth-logo-link"><AtonixCorpLogo size="medium" withText /></Link></div>
        <div className="auth-card">
          {status === 'verifying' && <><h1>Verifying your email</h1><p className="auth-subtitle">Confirming your secure verification link.</p></>}
          {status === 'sent' && <><h1>Check your inbox</h1><p className="auth-subtitle">We sent a verification link{email ? ` to ${email}` : ''}. Open it to continue with identity verification.</p><p className="auth-footer"><Link to="/login">Return to sign in</Link></p></>}
          {status === 'verified' && <><h1>Email verified</h1><p className="auth-subtitle">{message}</p></>}
          {status === 'failed' && <><h1>Verification unavailable</h1><p className="auth-error" role="alert">{message}</p><p className="auth-footer"><Link to="/login">Return to sign in to request a new link</Link></p></>}
        </div>
      </div>
    </div>
  );
};

export default EmailVerification;
