import React from 'react';
import { Link } from 'react-router-dom';
import { FaApple } from 'react-icons/fa';
import { FiBookOpen, FiBox, FiCode, FiCommand, FiDownload, FiExternalLink, FiMonitor, FiServer, FiShield } from 'react-icons/fi';
import { SiGoogleplay, SiLinux } from 'react-icons/si';

import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import './Developer.css';

const artifacts = [
  { name: 'atonixcorp CLI', version: 'v1.8.0', detail: 'Governance module deployment and environment management.', command: 'npm install -g @atonixcorp/cli', icon: FiCommand },
  { name: 'JavaScript SDK', version: 'v2.4.0', detail: 'Typed API client for governance, finance, and workflows.', command: 'npm install @atonixcorp/sdk', icon: FiCode },
  { name: 'Governance Toolbox', version: 'v1.3.2', detail: 'Reviewed templates, workflow scripts, and sample modules.', command: 'atonixcorp toolbox install governance-core', icon: FiBox },
];

const channels = [
  { label: 'Mac', detail: 'CLI package', icon: FaApple },
  { label: 'Windows', detail: 'CLI installer', icon: FiMonitor },
  { label: 'Linux', detail: 'Package archive', icon: SiLinux },
  { label: 'Google Play', detail: 'Mobile console', icon: SiGoogleplay },
  { label: 'App Store', detail: 'Mobile console', icon: FaApple },
];

const examples = [
  { label: 'Governance', code: "const amendment = await atonix.governance.amendments.create({\n  policy: 'GOV-001',\n  impactAnalysis: review\n});" },
  { label: 'Finance', code: "const ledger = await atonix.finance.ledger.list({\n  entity: 'entity_123',\n  period: '2026-Q3'\n});" },
  { label: 'Automation', code: "await atonix.workflows.deploy({\n  template: 'policy-attestation',\n  environment: 'sandbox'\n});" },
];

export default function DeveloperPortal() {
  return (
    <div className="developer-page">
      <Header />
      <main>
        <section className="developer-hero">
          <div className="developer-shell developer-hero__layout">
            <div>
              <p className="developer-eyebrow">AtonixCorp Developer Platform</p>
              <h1>Build governed systems with a platform your institution can trust.</h1>
              <p>Ship governance, finance, and policy automation with versioned tooling, documented APIs, sandbox-ready environments, and deployment records that remain visible to the organization.</p>
              <div className="developer-actions">
                <Link to="/register" className="developer-button developer-button--primary">Create developer workspace</Link>
                <Link to="/v1/docs" className="developer-button developer-button--secondary">Read API reference <FiExternalLink aria-hidden="true" /></Link>
              </div>
            </div>
            <aside className="developer-terminal" aria-label="Atonix CLI quick start">
              <div className="developer-terminal__bar"><span /><span /><span /> <strong>Quick start</strong></div>
              <code><span>$</span> npm install -g @atonixcorp/cli<br /><span>$</span> atonixcorp login<br /><span>$</span> atonixcorp sandbox create governance-lab<br /><span className="developer-terminal__success">Ready: sandbox governance-lab</span></code>
            </aside>
          </div>
        </section>

        <section className="developer-section">
          <div className="developer-shell">
            <div className="developer-section__heading"><p className="developer-eyebrow">Versioned tools</p><h2>Start with a controlled, reproducible toolchain.</h2></div>
            <div className="developer-artifacts">
              {artifacts.map(({ name, version, detail, command, icon: Icon }) => <article className="developer-artifact" key={name}><Icon aria-hidden="true" /><div><div className="developer-artifact__title"><h3>{name}</h3><span>{version}</span></div><p>{detail}</p><code>{command}</code></div><a href="#downloads" aria-label={`Download ${name}`}><FiDownload aria-hidden="true" /></a></article>)}
            </div>
          </div>
        </section>

        <section className="developer-section developer-section--silver" id="downloads">
          <div className="developer-shell"><div className="developer-section__heading"><p className="developer-eyebrow">Platform downloads</p><h2>Use the same governed platform from every operating environment.</h2></div><div className="developer-downloads">{channels.map(({ label, detail, icon: Icon }) => <a href="#download-request" className="developer-download" key={label}><Icon aria-hidden="true" /><span><strong>{label}</strong><small>{detail}</small></span><FiDownload aria-hidden="true" /></a>)}</div></div>
        </section>

        <section className="developer-section">
            <div className="developer-shell developer-code-layout"><div><p className="developer-eyebrow">Integration examples</p><h2>Bring governance controls into the systems teams already operate.</h2><p className="developer-body-copy">Each API is versioned and scoped. Start in a sandbox, validate your integration, and promote a reviewed module to an organization when it is ready.</p><div className="developer-resource-links"><Link to="/v1/docs"><FiBookOpen aria-hidden="true" /> API guides</Link><Link to="/contact"><FiShield aria-hidden="true" /> Security review</Link><Link to="/app/marketplace"><FiServer aria-hidden="true" /> Module marketplace</Link><Link to="/support"><FiBookOpen aria-hidden="true" /> Developer support and community</Link><Link to="/contact"><FiExternalLink aria-hidden="true" /> Partner and integration contact</Link></div></div><div className="developer-examples">{examples.map(({ label, code }) => <article key={label}><header>{label}<span>JavaScript</span></header><pre>{code}</pre></article>)}</div></div>
        </section>

        <section className="developer-section developer-section--navy"><div className="developer-shell developer-sandbox"><div><p className="developer-eyebrow">Sandbox environments</p><h2>Test with representative rules, without touching production records.</h2><p>Sandbox API keys, sample governance templates, and module manifests let teams validate their deployment before requesting organizational access.</p></div><Link to="/register" className="developer-button developer-button--light">Open a sandbox <FiExternalLink aria-hidden="true" /></Link></div></section>
      </main>
      <Footer />
    </div>
  );
}