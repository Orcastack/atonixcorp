import React from 'react';
import { Link } from 'react-router-dom';

import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import './About.css';

const About = () => {
  return (
    <div className="about-page">
      <Header />

      {/*  HERO  */}
      <section className="about-hero">
        <div className="about-hero-bg" />
        <div className="container">
          <div className="about-hero-inner">
            <p className="about-eyebrow">Governance and Enterprise Management Platform</p>
            <h1>AtonixCorp</h1>
            <p className="about-hero-sub">AtonixCorp is the digital framework for modern institutions.<br />It unifies governance, finance, equity, workflows, subscriptions, and analytics into one controlled platform.
            </p>
            <div className="about-hero-cta">
              <Link to="/register" className="btn-primary btn-large">Get Started
              </Link>
              <Link to="/features" className="btn-outline-hero btn-large">Explore Features
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/*  THE PROBLEM  */}
      <section className="about-problem-section">
        <div className="container">
          <div className="about-section-header">
            <p className="about-eyebrow-dark">The Core Problem</p>
            <h2>Why Organizations Need a Governance Operating System</h2>
            <p className="about-section-sub">Most institutions still manage policy, finance, HR, compliance, and approvals in separate tools.
              That fragmentation creates blind spots, slows execution, and weakens accountability.
            </p>
          </div>
          <div className="problem-grid">
            {[
              'Policy and approvals live in disconnected tools',
              'Governance rules are enforced manually or inconsistently',
              'Finance, equity, and compliance data are split across systems',
              'Leadership lacks real-time visibility into organizational health',
              'Teams lose time reconciling records and chasing sign-offs',
              'Subscription access and module control are handled ad hoc',
              'External integrations are bolted on instead of governed centrally',
              'Audit readiness is delayed by incomplete or inconsistent records',
            ].map((p) => (
              <div className="problem-row" key={p}>
                <span className="problem-x"></span>
                <span>{p}</span>
              </div>
            ))}
          </div>
          <div className="problem-resolve-banner">AtonixCorp exists to turn those fragmented processes into a governed digital system.
          </div>
        </div>
      </section>

      {/*  SOLUTION  */}
      <section className="about-solution-section">
        <div className="container">
          <div className="solution-split">
            <div className="solution-text">
              <p className="about-eyebrow-purple">The Solution</p>
              <h2>One platform. All governance operations. Fully connected.</h2>
              <p>AtonixCorp is a unified enterprise management environment — a single platform where organizations
                manage governance, finance, equity, policy enforcement, workflow automation, analytics, and subscriptions.
              </p>
              <p className="solution-manifesto">Everything in one place.<br />Everything connected.<br />Everything real-time.
              </p>
              <p>This is the new standard for governance-led enterprise operations.</p>
            </div>
            <div className="solution-cards">
              {[
                { label: 'Governance Engine' },
                { label: 'Policy Enforcement' },
                { label: 'Finance & Equity' },
                { label: 'Workflow Automation' },
                { label: 'Compliance' },
                { label: 'Analytics' },
                { label: 'Subscription Management' },
                { label: 'OrcaStack Integrations' },
                { label: 'Audit Trails' },
              ].map((c) => (
                <div className="solution-chip" key={c.label}>
                  <span className="solution-chip-icon">{c.icon}</span>
                  <span>{c.label}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/*  CORE PILLARS  */}
      <section className="about-pillars-section">
        <div className="container">
          <div className="about-section-header">
            <p className="about-eyebrow-dark">Core Architecture</p>
            <h2>Six Platform Layers</h2>
            <p className="about-section-sub">AtonixCorp is built as a layered operating system so every institution can subscribe to the right capabilities
              without losing governance, accountability, or auditability.
            </p>
          </div>
          <div className="about-pillars-grid">
            {[
              { n: '01', title: 'Governance Engine', desc: 'Define organizational hierarchy, roles, permissions, and policy boundaries.' },
              { n: '02', title: 'Policy Enforcement', desc: 'Automate compliance rules, approvals, and operational controls across the platform.' },
              { n: '03', title: 'Finance and Equity System', desc: 'Manage budgets, transactions, ownership, and capital records in one place.' },
              { n: '04', title: 'Workflow Automation', desc: 'Run tasks, reviews, and institutional processes through structured workflows.' },
              { n: '05', title: 'Analytics Dashboard', desc: 'Surface real-time insights for leadership, compliance, and decision-making.' },
              { n: '06', title: 'Subscription Layer', desc: 'Control access, billing, and organizational tiers with subscription-based entitlements.' },
            ].map((p) => (
              <div className="about-pillar-card" key={p.n}>
                <div className="about-pillar-num">{p.n}</div>
                <div className="about-pillar-icon">{p.icon}</div>
                <h3>{p.title}</h3>
                <p>{p.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/*  USER EXPERIENCE  */}
      <section className="about-ux-section">
        <div className="container">
          <div className="ux-split">
            <div className="ux-text">
              <p className="about-eyebrow-purple">The User Experience</p>
              <h2>What AtonixCorp Feels Like</h2>
              <p>When a user logs into AtonixCorp, they should immediately understand what the platform does:
                it organizes the institution, enforces policy, and gives leadership a clear operating view.
              </p>
              <div className="ux-qualities">
                {[
                  'A governance-first operating system',
                  'A premium, disciplined enterprise interface',
                  'A system built for serious institutional work',
                  'A platform that respects the user\'s time and accountability',
                ].map((q) => (
                  <div className="ux-quality" key={q}>

                    <span>{q}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="ux-feelings">
              <p className="ux-feelings-title">When you use AtonixCorp, you feel:</p>
              {['Informed', 'In Control', 'Structured', 'Efficient', 'Confident', 'Secure'].map((f) => (
                <div className="ux-feeling-chip" key={f}>{f}</div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/*  BRAND IDENTITY  */}
      <section className="about-brand-section">
        <div className="container">
          <div className="about-section-header">
            <p className="about-eyebrow-dark">Brand Identity</p>
            <h2>Built on Seven Commitments</h2>
          </div>
          <div className="brand-values-row">
            {['Precision', 'Security', 'Automation', 'Clarity', 'Governance', 'Scalability', 'Trust'].map((v) => (
              <div className="brand-value-tile" key={v}>{v}</div>
            ))}
          </div>
          <div className="brand-voice-grid">
            <div className="brand-voice-card">
              <h3>Brand Voice</h3>
              <ul>
                {['Confident', 'Clear', 'Professional', 'Modern', 'Authoritative', 'Vision-driven'].map((b) => (
                  <li key={b}> {b}</li>
                ))}
              </ul>
            </div>
            <div className="brand-promise-card">
              <h3>Brand Promise</h3>
              <blockquote>
                "AtonixCorp gives organizations the governance, automation, and visibility they need to
                operate with structure, transparency, and control."
              </blockquote>
            </div>
            <div className="brand-taglines-card">
              <h3>Brand Taglines</h3>
              <ul>
                {[
                  'The Financial Operating System for Modern Accounting Firms.',
                  'Where Accounting Meets Automation.',
                  'Real-Time Finance. Real-Time Control.',
                  'Built for Firms That Refuse to Fall Behind.',
                  'Your Entire Financial World. Unified.',
                ].map((t) => (
                  <li key={t}><span className="tagline-dash">—</span> {t}</li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/*  FUTURE VISION  */}
      <section className="about-future-section">
        <div className="container">
          <div className="future-inner">
            <p className="about-eyebrow about-eyebrow--inverse">Long-Term Vision</p>
            <h2>The Future of AtonixCorp</h2>
            <p className="future-sub">AtonixCorp is not just a platform — it is a movement. A transformation.
              A new standard for how the world manages financial operations.
            </p>
            <div className="future-grid">
              {[
                { label: 'Global Banking Integrations' },
                { label: 'AI-Driven Financial Forecasting' },
                { label: 'Automated Compliance Engines' },
                { label: 'Full Tax Automation' },
                { label: 'Enterprise-Grade Analytics' },
                { label: 'Cross-Border Financial Intelligence' },
                { label: 'Global Marketplace of Financial Tools' },
              ].map((f) => (
                <div className="future-item" key={f.label}>
                  <div className="future-icon">{f.icon}</div>
                  <span>{f.label}</span>
                </div>
              ))}
            </div>
            <p className="future-closing">AtonixCorp will become the platform that powers the world's financial operations.
            </p>
          </div>
        </div>
      </section>

      {/*  CTA  */}
      <section className="about-cta-section">
        <div className="container">
          <div className="about-cta-inner">
            <h2>This is AtonixCorp.<br />
              <span>The future of financial operations begins here.</span>
            </h2>
            <p>Built for the future. Built for firms that demand excellence. Built for businesses
              that want clarity. Built for financial institutions that require precision.
            </p>
            <div className="about-cta-buttons">
              <Link to="/register" className="btn-primary btn-large">Get Started Today
              </Link>
              <Link to="/contact" className="btn-outline btn-large">Talk to Us
              </Link>
            </div>
          </div>
        </div>
      </section>

      <Footer />
    </div>
  );
};

export default About;
