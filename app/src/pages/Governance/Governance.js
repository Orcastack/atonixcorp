import React from 'react';
import { Link } from 'react-router-dom';
import { FiArrowRight, FiCheckCircle, FiFileText, FiShield } from 'react-icons/fi';

import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import './Governance.css';

const pillars = [
  ['Sovereignty and independence', 'Decision-making, technology, and institutional continuity are governed with a long-term responsibility standard.'],
  ['Ethics and accountability', 'Material decisions receive review, evidence, and a durable record that can be inspected after the moment has passed.'],
  ['Security and information stewardship', 'Sensitive information is protected through disciplined access, lifecycle control, and incident response.'],
];

const amendmentRules = [
  ['Standard', '60%', '72 hours'],
  ['Operational', '70%', '72 hours'],
  ['Ethical or security', '75%', '72 hours'],
  ['Constitutional', '80%', '7 days'],
  ['Sovereignty', '90%', '7 days'],
  ['Emergency', '75%', '24 hours'],
];

export default function Governance() {
  return (
    <div className="governance-page">
      <Header />
      <main>
        <section className="governance-hero">
          <div className="governance-shell governance-hero__layout">
            <div>
              <p className="governance-eyebrow">AtonixCorp governance charter</p>
              <h1>Governance that remains accountable over time.</h1>
              <p className="governance-lead">
                AtonixCorp turns governing principles into controlled policy editions, reviewable amendments,
                verified decisions, and an enduring institutional record.
              </p>
              <div className="governance-actions">
                <Link className="governance-button governance-button--primary" to="/register">Establish a workspace <FiArrowRight aria-hidden="true" /></Link>
                <a className="governance-button governance-button--secondary" href="#edition">Review edition notes</a>
              </div>
            </div>
            <div className="governance-hero__seal" aria-label="Governance charter edition 1.1 editorial correction">
              <FiShield aria-hidden="true" />
              <span>Controlled charter</span>
              <strong>Edition 1.1</strong>
              <small>Editorial correction register</small>
            </div>
          </div>
        </section>

        <section className="governance-section">
          <div className="governance-shell">
            <div className="governance-section__heading">
              <p className="governance-eyebrow">Operating commitments</p>
              <h2>Designed for decisions that deserve evidence.</h2>
            </div>
            <div className="governance-pillars">
              {pillars.map(([title, text], index) => (
                <article className="governance-pillar" key={title}>
                  <span className="governance-pillar__number">0{index + 1}</span>
                  <h3>{title}</h3>
                  <p>{text}</p>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="governance-section governance-section--silver" id="edition">
          <div className="governance-shell governance-edition">
            <div className="governance-edition__copy">
              <p className="governance-eyebrow">Edition control</p>
              <h2>Every amendment is a complete case file.</h2>
              <p>Proposals record their rationale, impact analysis, ethical review, sovereignty check, operational feasibility, security implications, and delivery timeline before they reach a vote.</p>
              <ul className="governance-checklist">
                <li><FiCheckCircle aria-hidden="true" /> Internal division, council, and final-authority review stages</li>
                <li><FiCheckCircle aria-hidden="true" /> Identity-bound voting and permanent result evidence</li>
                <li><FiCheckCircle aria-hidden="true" /> Edition, approval, and publication records retained together</li>
              </ul>
            </div>
            <aside className="governance-edition__note">
              <FiFileText aria-hidden="true" />
              <h3>Current editorial register</h3>
              <p>Edition 1.1 resolves duplicated leadership-tenure text, restores chapter sequencing for the affected succession and council clauses, and corrects typographic artifacts without changing policy intent.</p>
              <span>Source preserved: AtonixCorp-Gorvernance-policy.pdf</span>
            </aside>
          </div>
        </section>

        <section className="governance-section">
          <div className="governance-shell">
            <div className="governance-section__heading governance-section__heading--split">
              <div>
                <p className="governance-eyebrow">Amendment protocol</p>
                <h2>Thresholds reflect the consequence of a change.</h2>
              </div>
              <Link className="governance-text-link" to="/login">Open Governance Center <FiArrowRight aria-hidden="true" /></Link>
            </div>
            <div className="governance-thresholds" role="table" aria-label="Governance amendment voting thresholds">
              <div className="governance-thresholds__head" role="row"><span>Amendment class</span><span>Approval required</span><span>Voting window</span></div>
              {amendmentRules.map(([kind, threshold, window]) => (
                <div className="governance-thresholds__row" role="row" key={kind}><span>{kind}</span><strong>{threshold}</strong><span>{window}</span></div>
              ))}
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}