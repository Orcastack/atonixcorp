import React, { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import AtonixCorpLogo from '../../components/branding/AtonixCorpLogo';

const EmailVerification = () => {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { completeEmailVerification, resendEmailVerification } = useAuth();
  const [status, setStatus] = useState(searchParams.get('token') ? 'verifying' : 'sent');
  const [message, setMessage] = useState('');
  const [resendStatus, setResendStatus] = useState('');
  const [isResending, setIsResending] = useState(false);

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
  const requestNewLink = async () => {
    if (!email) {
      setResendStatus('Return to registration and enter your email address to request a new link.');
      return;
    }

    setIsResending(true);
    setResendStatus('');
    const result = await resendEmailVerification(email);
    setIsResending(false);
    setResendStatus(result.success
      ? 'If this account requires verification, a new link has been sent.'
      : (result.error || 'Unable to request a new verification link.'));
  };

  return (
    <div className="auth-page">
      <div className="auth-container">
        <div className="auth-header"><Link to="/" className="auth-logo-link"><AtonixCorpLogo size="medium" withText /></Link></div>
        <div className="auth-card">
          {status === 'verifying' && <><h1>Verifying your email</h1><p className="auth-subtitle">Confirming your secure verification link.</p></>}
          {status === 'sent' && <><h1>Check your inbox</h1><p className="auth-subtitle">{location.state?.existingAccount ? 'Your account already exists and still needs email verification.' : 'We sent a verification link'}{email ? ` to ${email}` : ''}. Open it to continue with identity verification.</p><button type="button" className="btn-link" onClick={requestNewLink} disabled={isResending}>{isResending ? 'Sending verification link...' : 'Send a new verification link'}</button>{resendStatus && <p className="auth-subtitle" role="status">{resendStatus}</p>}<p className="auth-footer"><Link to="/login">Return to sign in</Link></p></>}
          {status === 'verified' && <><h1>Email verified</h1><p className="auth-subtitle">{message}</p></>}
          {status === 'failed' && <><h1>Verification unavailable</h1><p className="auth-error" role="alert">{message}</p><p className="auth-footer"><Link to="/login">Return to sign in to request a new link</Link></p></>}
        </div>
      </div>
    </div>
  );
};

export default EmailVerification;
