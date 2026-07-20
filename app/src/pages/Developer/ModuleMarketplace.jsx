import React, { useEffect, useState } from 'react';
import { FiCheckCircle, FiLock, FiPackage, FiRefreshCw } from 'react-icons/fi';
import { Card, PageHeader, Button } from '../../components/ui';
import { useEnterprise } from '../../context/EnterpriseContext';
import { developerModuleInstallationsAPI } from '../../services/api';
import './Developer.css';

const modules = [
  { key: 'governance-core', name: 'Governance Core', category: 'Governance engine', tier: 'basic', version: '1.8.0', description: 'Policy register, amendment protocols, and governed evidence records.' },
  { key: 'finance-controls', name: 'Finance Control Pack', category: 'Finance & equity', tier: 'professional', version: '2.1.0', description: 'Approval controls, entity finance policies, and audit-ready workflows.' },
  { key: 'policy-attestation', name: 'Policy Attestation', category: 'Policy enforcement', tier: 'professional', version: '1.6.3', description: 'Scheduled attestations, exception routing, and completion evidence.' },
  { key: 'institutional-analytics', name: 'Institutional Analytics', category: 'Analytics dashboard', tier: 'enterprise', version: '3.0.1', description: 'Executive control coverage, risk signals, and board-ready reporting.' },
  { key: 'sovereignty-controls', name: 'Sovereignty Controls', category: 'Governance engine', tier: 'institutional', version: '1.2.0', description: 'Enhanced sovereign review, continuity controls, and protected decision records.' },
];
const rank = { basic: 0, professional: 1, enterprise: 2, institutional: 3 };
const normalize = (response) => response.data.results || response.data || [];

export default function ModuleMarketplace() {
  const { currentOrganization } = useEnterprise();
  const [installed, setInstalled] = useState([]);
  const [deploying, setDeploying] = useState('');
  const [notice, setNotice] = useState('');
  const [error, setError] = useState('');
  const tier = (currentOrganization?.settings?.subscription_tier || 'basic').toLowerCase();

  useEffect(() => {
    if (!currentOrganization?.id) return;
    developerModuleInstallationsAPI.getAll({ organization: currentOrganization.id })
      .then((response) => setInstalled(normalize(response)))
      .catch(() => setError('Installed modules could not be loaded.'));
  }, [currentOrganization?.id]);

  const deploy = async (module) => {
    if (!currentOrganization?.id) { setError('Choose an organization before deploying a module.'); return; }
    setDeploying(module.key); setError(''); setNotice('');
    try {
      const response = await developerModuleInstallationsAPI.deploy({ organization: currentOrganization.id, module_key: module.key, module_name: module.name, category: module.category, version: module.version, required_tier: module.tier });
      setInstalled((current) => [...current.filter((item) => item.module_key !== module.key), response.data]);
      setNotice(`${module.name} is active in ${currentOrganization.name}. The deployment was added to the audit trail.`);
    } catch (requestError) {
      setError(requestError.response?.data?.required_tier?.[0] || 'Module deployment could not be completed.');
    } finally { setDeploying(''); }
  };

  return <div className="module-page developer-marketplace"><PageHeader title="Module Marketplace" subtitle="Deploy governed capabilities to your organization with subscription-aware access and a complete installation record." />
    {notice && <div className="developer-notice developer-notice--success"><FiCheckCircle aria-hidden="true" />{notice}</div>}{error && <div className="developer-notice developer-notice--error">{error}</div>}
    <Card className="developer-marketplace__summary"><div><span>Current entitlement</span><strong>{tier}</strong><small>{currentOrganization?.name || 'No organization selected'}</small></div><div><span>Installed modules</span><strong>{installed.filter((item) => item.status === 'active').length}</strong><small>Deployment events are retained in Platform Audit</small></div><div><span>Release channel</span><strong>Stable</strong><small>Versioned updates require a recorded change</small></div></Card>
    <div className="developer-module-grid">{modules.map((module) => { const record = installed.find((item) => item.module_key === module.key); const eligible = rank[tier] >= rank[module.tier]; return <Card className="developer-module-card" key={module.key}><div className="developer-module-card__top"><span>{module.category}</span><strong>v{module.version}</strong></div><h2>{module.name}</h2><p>{module.description}</p><div className="developer-module-card__footer"><span className={eligible ? 'developer-tier' : 'developer-tier developer-tier--locked'}>{eligible ? <FiCheckCircle aria-hidden="true" /> : <FiLock aria-hidden="true" />}{module.tier}</span>{record ? <Button variant="secondary" disabled icon={FiPackage}>Active</Button> : <Button variant="primary" disabled={!eligible || deploying === module.key} onClick={() => deploy(module)} icon={deploying === module.key ? FiRefreshCw : FiPackage}>{deploying === module.key ? 'Deploying' : eligible ? 'Deploy' : 'Upgrade required'}</Button>}</div></Card>; })}</div>
  </div>;
}