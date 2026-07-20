import React from 'react';
import { Link, Navigate } from 'react-router-dom';

import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import { useAuth } from '../../context/AuthContext';
import './Landing.css';

const systemSteps = [
  {
    step: '01',
    title: 'Define the governance model',
    text: 'Map organizational hierarchy, roles, permissions, and policy boundaries into one controlled environment.',
  },
  {
    step: '02',
    title: 'Automate policy execution',
    text: 'Move approvals, compliance checks, workflows, and exception handling through clear, reviewable stages.',
  },
  {
    step: '03',
    title: 'Stay inspection-ready',
    text: 'Surface decisions, audit trails, subscription access, and supporting evidence in a form that withstands scrutiny.',
  },
];

const operatingPrinciples = [
  'One environment for governance, finance, equity, and workflows — without losing enterprise visibility.',
  'Controllers, administrators, compliance leads, and finance teams share one source of operational truth.',
  'Execution stays structured enough for audits, board reviews, policy enforcement, and regulated reporting cycles.',
  'Subscription tiers control access to the right modules, dashboards, and integrations for each organization.',
];

const operationalPillars = [
  {
    title: 'Governance engine',
    text: 'Maintain organizational hierarchy, role-based access, policy rules, and board-ready evidence in one controlled environment.',
  },
  {
    title: 'Finance and equity',
    text: 'Keep ledger activity, allocations, ownership, reconciliations, and reporting aligned across entities and divisions.',
  },
  {
    title: 'Policy and workflow enforcement',
    text: 'Track obligations, filing calendars, subscription entitlements, approvals, and review checkpoints without losing traceability.',
  },
];

const heroSignals = ['Entity governance', 'Approval discipline', 'Filing readiness'];

const heroSummary = [
  { label: 'Coverage', value: 'Governance + Finance + Equity' },
  { label: 'Execution mode', value: 'Structured, subscribed, and accountable' },
];

const heroDomains = ['Governance', 'Finance', 'Equity', 'Policy', 'Analytics'];

const heroTimeline = [
  {
    title: 'Governance records',
    detail: 'Ownership, approvals, permissions',
    active: true,
  },
  {
    title: 'Operational books',
    detail: 'Ledger activity, controls, reporting',
  },
  {
    title: 'Compliance execution',
    detail: 'Tax obligations, filings, evidence',
  },
];

const Landing = () => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return (
      <div className="landing-page landing-page--loading">
        Loading...
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/app/console" replace />;
  }

  return (
    <div className="landing-page">
      <Header />

      <main className="landing-main">
        <section className="landing-hero" aria-label="AtonixCorp overview">
          <div className="landing-shell landing-hero__layout">
            <div className="landing-hero__copy">
              <p className="landing-kicker">Governance-as-a-Service for modern institutions</p>
              <h1>One operating system for governance, finance, equity, subscriptions, and controlled execution.</h1>
              <p className="landing-lead">
                AtonixCorp gives organizations a single platform for policy enforcement, finance and equity management,
                workflow approvals, analytics, and subscription-based access across the whole institution.
              </p>

              <div className="landing-actions">
                <Link to="/register" className="landing-button landing-button--primary">Open Account</Link>
                <Link to="/features" className="landing-button landing-button--secondary">Review Modules</Link>
              </div>

              <div className="landing-hero__trustband" aria-label="Operational trust indicators">
                {heroSignals.map((signal) => (
                  <span key={signal}>{signal}</span>
                ))}
              </div>
              <p className="landing-hero__microcopy">
                Built to align with OrcaOS, OrcaCLI, OrcaSDK, and OrcaCompute for end-to-end institutional operations.
              </p>
            </div>

            <aside className="landing-hero__frame" aria-label="Platform operating frame">
              <div className="landing-hero__frame-bar" />
              <p className="landing-hero__frame-title">Operating Domains</p>
              <div className="landing-hero__domain-list">
                {heroDomains.map((domain) => (
                  <span key={domain}>{domain}</span>
                ))}
              </div>
              <div className="landing-hero__summary-grid">
                {heroSummary.map((item) => (
                  <div key={item.label} className="landing-hero__summary-card">
                    <span>{item.label}</span>
                    <strong>{item.value}</strong>
                  </div>
                ))}
              </div>
              <div className="landing-hero__frame-copy">
                {heroTimeline.map((item) => (
                  <div key={item.title} className={`landing-hero__timeline-item${item.active ? ' is-active' : ''}`}>
                    <span className="landing-hero__timeline-dot" aria-hidden="true" />
                    <div>
                      <p>{item.title}</p>
                      <small>{item.detail}</small>
                    </div>
                  </div>
                ))}
              </div>
            </aside>
          </div>
        </section>

        <section className="landing-section landing-section--system">
          <div className="landing-shell">
            <div className="landing-section__intro">
              <p className="landing-kicker">What the platform does</p>
              <h2>It turns governance into a digital operating framework.</h2>
              <p>
                Organizations subscribe to AtonixCorp to manage policy, finance, equity, approvals, and records from one dashboard.
                Each tier unlocks the right modules and integrations for the organization’s size and complexity.
              </p>
            </div>
            <div className="landing-system-grid">
              {systemSteps.map((lane) => (
                <article key={lane.step} className="landing-system-card">
                  <span className="landing-system-card__step">{lane.step}</span>
                  <h3>{lane.title}</h3>
                  <p>{lane.text}</p>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="landing-section">
          <div className="landing-shell landing-philosophy">
            <ul className="landing-principles">
              {operatingPrinciples.map((principle) => (
                <li key={principle}>{principle}</li>
              ))}
            </ul>
          </div>
        </section>

        <section className="landing-section landing-section--pillars">
          <div className="landing-shell">
            <div className="landing-pillars-grid">
              {operationalPillars.map((pillar) => (
                <article key={pillar.title} className="landing-pillar-card">
                  <h3>{pillar.title}</h3>
                  <p>{pillar.text}</p>
                </article>
              ))}
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
};

export default Landing;
