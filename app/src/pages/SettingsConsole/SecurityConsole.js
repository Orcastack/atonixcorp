import React, { useEffect, useMemo, useState } from 'react';
import { useEnterprise } from '../../context/EnterpriseContext';
import { securityApiKeysAPI } from '../../services/api';
import './SecurityConsole.css';

const DEFAULT_SECURITY = {
  mandatory_2fa: true,
  allowed_2fa_methods: ['totp', 'email', 'sms'],
  encryption_standard: 'AES-256-GCM',
  transport_standard: 'TLS 1.3',
  compliance_frameworks: ['ISO 27001', 'SOC 2', 'GDPR'],
};

const SecurityConsole = () => {
  const { currentOrganization, updateOrganization, hasPermission } = useEnterprise();
  const canManage = hasPermission ? hasPermission('manage_org_settings') : false;
  const organizationId = currentOrganization?.enterprise_code || currentOrganization?.public_id || currentOrganization?.id;
  const [security, setSecurity] = useState(DEFAULT_SECURITY);
  const [keys, setKeys] = useState([]);
  const [keyName, setKeyName] = useState('');
  const [keyScope, setKeyScope] = useState('workspace:read');
  const [newSecret, setNewSecret] = useState('');
  const [status, setStatus] = useState('');
  const [loadingKeys, setLoadingKeys] = useState(false);

  useEffect(() => {
    const saved = currentOrganization?.settings?.security || {};
    setSecurity({ ...DEFAULT_SECURITY, ...saved });
  }, [currentOrganization]);

  useEffect(() => {
    if (!organizationId) return;
    setLoadingKeys(true);
    securityApiKeysAPI.list(organizationId)
      .then((response) => setKeys(Array.isArray(response.data) ? response.data : []))
      .catch(() => setKeys([]))
      .finally(() => setLoadingKeys(false));
  }, [organizationId]);

  const roleLabel = useMemo(() => {
    if (currentOrganization?.owner_id) return 'Organization owner';
    return 'Scoped administrator';
  }, [currentOrganization?.owner_id]);

  const saveSecurity = async () => {
    if (!canManage || !currentOrganization?.id) return;
    setStatus('Saving security policy...');
    try {
      await updateOrganization(currentOrganization.id, {
        settings: { ...(currentOrganization.settings || {}), security },
      });
      setStatus('Security policy saved and audit recorded.');
    } catch (error) {
      setStatus(error?.message || 'Unable to save security policy.');
    }
  };

  const createKey = async (event) => {
    event.preventDefault();
    if (!keyName.trim() || !organizationId || !canManage) return;
    setStatus('Generating secure API key...');
    try {
      const response = await securityApiKeysAPI.create(organizationId, {
        name: keyName.trim(),
        scopes: [keyScope],
      });
      setKeys((current) => [response.data, ...current]);
      setNewSecret(response.data.client_secret || '');
      setKeyName('');
      setStatus('Key created. Copy the secret now; it will not be shown again.');
    } catch (error) {
      setStatus(error?.response?.data?.message || error?.response?.data?.detail || 'Unable to create API key.');
    }
  };

  const revokeKey = async (key) => {
    if (!window.confirm(`Revoke ${key.name}?`)) return;
    await securityApiKeysAPI.revoke(organizationId, key.id);
    setKeys((current) => current.map((item) => item.id === key.id ? { ...item, is_active: false, status: 'revoked' } : item));
    setStatus('API key revoked and audit recorded.');
  };

  return (
    <div className="security-console">
      <header className="security-console-hero">
        <div>
          <span className="security-console-kicker">Security control plane</span>
          <h1>Security Settings Console</h1>
          <p>Mandatory access controls, scoped credentials, encryption posture, and compliance evidence in one governed surface.</p>
        </div>
        <span className="security-role-badge">{roleLabel}</span>
      </header>

      {status && <div className="security-console-status" role="status">{status}</div>}

      <section className="security-posture-grid" aria-label="Security posture">
        <article><span>Encryption at rest</span><strong>AES-256-GCM</strong><small>Enforced by platform services</small></article>
        <article><span>Transport security</span><strong>TLS 1.3</strong><small>Production transport baseline</small></article>
        <article><span>Compliance posture</span><strong>ISO 27001 · SOC 2</strong><small>GDPR controls mapped</small></article>
        <article><span>Cryptographic acceleration</span><strong>Service-managed</strong><small>GPU telemetry is infrastructure scoped</small></article>
      </section>

      <div className="security-console-grid">
        <section className="security-panel">
          <div className="security-panel-heading"><div><span className="security-console-kicker">Identity assurance</span><h2>Mandatory two-factor authentication</h2></div><span className="security-enforced">Required</span></div>
          <p className="security-panel-copy">All users must complete a second factor before login and sensitive actions such as billing and equity changes.</p>
          <div className="security-methods">
            {['totp', 'email', 'sms'].map((method) => <label key={method}><input type="checkbox" checked={security.allowed_2fa_methods.includes(method)} disabled={!canManage} onChange={() => setSecurity((current) => ({ ...current, allowed_2fa_methods: current.allowed_2fa_methods.includes(method) ? current.allowed_2fa_methods.filter((item) => item !== method) : [...current.allowed_2fa_methods, method] }))} />{method === 'totp' ? 'Authenticator app (TOTP)' : method.toUpperCase()}</label>)}
          </div>
          <button className="security-primary-button" type="button" disabled={!canManage} onClick={saveSecurity}>Save enforced policy</button>
        </section>

        <section className="security-panel">
          <div className="security-panel-heading"><div><span className="security-console-kicker">Credential governance</span><h2>API keys</h2></div><span className="security-enforced">Role scoped</span></div>
          <form className="security-key-form" onSubmit={createKey}><input value={keyName} onChange={(event) => setKeyName(event.target.value)} placeholder="Key name" aria-label="Key name" disabled={!canManage} /><select value={keyScope} onChange={(event) => setKeyScope(event.target.value)} aria-label="API key scope" disabled={!canManage}><option value="workspace:read">Workspace read</option><option value="enterprise:read">Enterprise read</option><option value="enterprise:write">Enterprise write</option></select><button className="security-primary-button" type="submit" disabled={!canManage || !keyName.trim()}>Generate</button></form>
          {newSecret && <div className="security-secret"><strong>One-time secret</strong><code>{newSecret}</code><button type="button" onClick={() => navigator.clipboard?.writeText(newSecret)}>Copy</button></div>}
          <div className="security-key-list">{loadingKeys ? <span>Loading keys...</span> : keys.length === 0 ? <span>No API keys registered.</span> : keys.map((key) => <div className="security-key-row" key={key.id}><span><strong>{key.name}</strong><small>{key.client_id || key.token_prefix || 'Scoped credential'} · {(key.scopes || []).join(', ')}</small></span><button type="button" onClick={() => revokeKey(key)} disabled={!canManage || key.status === 'revoked' || key.is_active === false}>{key.status === 'revoked' || key.is_active === false ? 'Revoked' : 'Revoke'}</button></div>)}</div>
        </section>
      </div>
    </div>
  );
};

export default SecurityConsole;