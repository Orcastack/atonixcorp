import React, { useMemo, useState } from 'react';
import { useEquity } from '../../../context/EquityContext';
import { buildBalancedMetricOrder } from '../../../utils/dashboardMetrics';
import '../components/EquityModuleScreen.css';
import '../components/EquityCrudModuleScreen.css';

const EMPTY_FORM = {
  grant_number: '',
  shareholder: '',
  share_class: '',
  grant_type: 'stock_option',
  total_units: '0',
  exercise_price: '0',
  grant_date: '',
  vesting_start_date: '',
  cliff_months: '12',
  vesting_months: '48',
  vesting_interval: 'monthly',
  acceleration_type: 'none',
  termination_treatment: 'forfeit_unvested',
  expiration_date: '',
  notes: '',
};

const VestingGrants = () => {
  const {
    grants,
    vestingEvents,
    shareholders,
    shareClasses,
    summary,
    loading,
    error,
    saving,
    createGrant,
    updateGrant,
    deleteGrant,
    downloadGrantPackage,
    regenerateGrantPackage,
    rebuildGrantSchedule,
    triggerSingleAcceleration,
    triggerDoubleAcceleration,
    terminateGrant,
  } = useEquity();
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);

  const metrics = useMemo(() => {
    const exercisableUnits = grants.reduce(
      (total, grant) => total + Number(grant.vesting_summary?.available_to_exercise || 0),
      0,
    );
    const pendingVestEvents = vestingEvents.filter((event) => event.status === 'pending').length;
    return [
      { label: 'Active grants', value: summary.totalGrants, note: 'Founders, options, RSUs, ESOPs, and warrants' },
      { label: 'Exercisable units', value: exercisableUnits, note: 'Units already vested and available to exercise' },
      { label: 'Pending vesting events', value: pendingVestEvents, note: 'Future vesting releases queued by the engine' },
    ];
  }, [grants, summary.totalGrants, vestingEvents]);

  const vestingTimeline = useMemo(() => {
    return [...vestingEvents]
      .sort((left, right) => new Date(left.vest_date) - new Date(right.vest_date))
      .slice(0, 6);
  }, [vestingEvents]);

  const timelineMax = useMemo(() => Math.max(...vestingTimeline.map((event) => Number(event.units || 0)), 1), [vestingTimeline]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const resetForm = () => {
    setForm(EMPTY_FORM);
    setEditingId(null);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const payload = {
      ...form,
      total_units: Number(form.total_units || 0),
      exercise_price: Number(form.exercise_price || 0),
      cliff_months: Number(form.cliff_months || 0),
      vesting_months: Number(form.vesting_months || 0),
    };
    if (editingId) {
      await updateGrant(editingId, payload);
    } else {
      await createGrant(payload);
    }
    resetForm();
  };

  const handleEdit = (grant) => {
    setEditingId(grant.id);
    setForm({
      grant_number: grant.grant_number || '',
      shareholder: grant.shareholder || '',
      share_class: grant.share_class || '',
      grant_type: grant.grant_type || 'stock_option',
      total_units: String(grant.total_units ?? 0),
      exercise_price: String(grant.exercise_price ?? 0),
      grant_date: grant.grant_date || '',
      vesting_start_date: grant.vesting_start_date || '',
      cliff_months: String(grant.cliff_months ?? 12),
      vesting_months: String(grant.vesting_months ?? 48),
      vesting_interval: grant.vesting_interval || 'monthly',
      acceleration_type: grant.acceleration_type || 'none',
      termination_treatment: grant.termination_treatment || 'forfeit_unvested',
      expiration_date: grant.expiration_date || '',
      notes: grant.notes || '',
    });
  };

  const handleDelete = async (grantId) => {
    await deleteGrant(grantId);
    if (editingId === grantId) {
      resetForm();
    }
  };

  const today = new Date().toISOString().slice(0, 10);

  return (
    <section className="eq-screen">
      <div className="eq-screen-hero">
        <div>
          <h2>Vesting Engine</h2>
          <p>Configure grant schedules, cliffs, accelerations, termination treatment, and automatically track units through grant, vest, forfeiture, and exercise.</p>
        </div>
        <div className="eq-screen-banner">
          The vesting engine turns grants into dated vesting events, acceleration records, payroll-tax sync data, and exercise-ready units.
        </div>
      </div>

      <div className="eq-metric-grid">
        {buildBalancedMetricOrder(metrics, metrics.length).map((metric) => (
          <article key={metric.label} className="eq-metric-card">
            <span className="eq-metric-label">{metric.label}</span>
            <strong className="eq-metric-value">{metric.value}</strong>
            <span className="eq-metric-note">{metric.note}</span>
          </article>
        ))}
      </div>

      <div className="eq-crud-layout">
        <aside className="eq-data-card eq-form-card">
          <div className="eq-data-card-head">
            <h3>{editingId ? 'Edit grant' : 'Create grant'}</h3>
            {editingId && (
              <button type="button" className="eq-inline-btn secondary" onClick={resetForm}>
                Cancel edit
              </button>
            )}
          </div>
          <p className="eq-form-copy">Set the legal terms of a grant once, and the engine will maintain vesting releases, trigger acceleration, and prepare exercise eligibility.</p>
          <form className="eq-form-grid" onSubmit={handleSubmit}>
            <label className="eq-form-field">
              <span className="eq-form-label">Grant number</span>
              <input className="eq-form-input" name="grant_number" value={form.grant_number} onChange={handleChange} placeholder="GRANT-2026-001" />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Grant type</span>
              <select className="eq-form-select" name="grant_type" value={form.grant_type} onChange={handleChange}>
                <option value="stock_option">Stock option</option>
                <option value="esop">ESOP</option>
                <option value="rsu">RSU</option>
                <option value="restricted_stock">Restricted stock</option>
                <option value="founder">Founder vesting</option>
                <option value="warrant">Warrant</option>
              </select>
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Shareholder</span>
              <select className="eq-form-select" name="shareholder" value={form.shareholder} onChange={handleChange}>
                <option value="">Select shareholder</option>
                {shareholders.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
              </select>
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Share class</span>
              <select className="eq-form-select" name="share_class" value={form.share_class} onChange={handleChange}>
                <option value="">Select share class</option>
                {shareClasses.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
              </select>
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Total units</span>
              <input className="eq-form-input" type="number" min="0" step="1" name="total_units" value={form.total_units} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Exercise price</span>
              <input className="eq-form-input" type="number" min="0" step="0.0001" name="exercise_price" value={form.exercise_price} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Grant date</span>
              <input className="eq-form-input" type="date" name="grant_date" value={form.grant_date} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Vesting start</span>
              <input className="eq-form-input" type="date" name="vesting_start_date" value={form.vesting_start_date} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Cliff months</span>
              <input className="eq-form-input" type="number" min="0" step="1" name="cliff_months" value={form.cliff_months} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Vesting months</span>
              <input className="eq-form-input" type="number" min="1" step="1" name="vesting_months" value={form.vesting_months} onChange={handleChange} />
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Interval</span>
              <select className="eq-form-select" name="vesting_interval" value={form.vesting_interval} onChange={handleChange}>
                <option value="monthly">Monthly</option>
                <option value="quarterly">Quarterly</option>
                <option value="annual">Annual</option>
                <option value="custom">Custom</option>
              </select>
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Acceleration</span>
              <select className="eq-form-select" name="acceleration_type" value={form.acceleration_type} onChange={handleChange}>
                <option value="none">None</option>
                <option value="single">Single trigger</option>
                <option value="double">Double trigger</option>
              </select>
            </label>
            <label className="eq-form-field">
              <span className="eq-form-label">Termination treatment</span>
              <select className="eq-form-select" name="termination_treatment" value={form.termination_treatment} onChange={handleChange}>
                <option value="forfeit_unvested">Forfeit unvested</option>
                <option value="continue_vesting">Continue vesting</option>
                <option value="accelerate_to_cliff">Accelerate to cliff</option>
                <option value="full_acceleration">Full acceleration</option>
              </select>
            </label>
            <label className="eq-form-field full">
              <span className="eq-form-label">Notes</span>
              <textarea className="eq-form-textarea" rows="4" name="notes" value={form.notes} onChange={handleChange} placeholder="Plan terms, board approval notes, and termination handling details" />
            </label>
            <div className="eq-form-actions">
              <button type="submit" className="eq-inline-btn primary" disabled={saving}>
                {saving ? 'Saving…' : editingId ? 'Save grant' : 'Create grant'}
              </button>
            </div>
          </form>
        </aside>

        <div className="eq-data-card">
          <div className="eq-data-card-head">
            <h3>Grant Register</h3>
            {!loading && !error && <span className="eq-status-chip success">Live</span>}
            {loading && <span className="eq-status-chip">Syncing</span>}
            {error && <span className="eq-status-chip danger">Attention</span>}
          </div>
          {error && <div className="eq-error-banner">{error}</div>}
          <div className="eq-table-wrap">
            <table className="eq-table">
              <thead>
                <tr>
                  <th>Grant</th>
                  <th>Holder</th>
                  <th>Type</th>
                  <th>Total</th>
                  <th>Exercisable</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {grants.map((grant) => (
                  <tr key={grant.id}>
                    <td>{grant.grant_number}</td>
                    <td>{grant.shareholder_name || '—'}</td>
                    <td>{grant.grant_type}</td>
                    <td>{grant.total_units}</td>
                    <td>{grant.vesting_summary?.available_to_exercise || 0}</td>
                    <td>{grant.lifecycle_status}</td>
                    <td>
                      <div className="eq-table-actions">
                        <button type="button" className="eq-inline-btn" onClick={() => handleEdit(grant)}>Edit</button>
                        <button type="button" className="eq-inline-btn" onClick={() => downloadGrantPackage(grant.id, `grant-package-${grant.grant_number}.pdf`)}>PDF</button>
                        <button type="button" className="eq-inline-btn" onClick={() => regenerateGrantPackage(grant.id)}>Regen PDF</button>
                        <button type="button" className="eq-inline-btn" onClick={() => rebuildGrantSchedule(grant.id)}>Rebuild</button>
                        <button type="button" className="eq-inline-btn" onClick={() => triggerSingleAcceleration(grant.id, { trigger_date: today })}>Single</button>
                        <button type="button" className="eq-inline-btn" onClick={() => triggerDoubleAcceleration(grant.id, { trigger_date: today })}>Double</button>
                        <button type="button" className="eq-inline-btn danger" onClick={() => terminateGrant(grant.id, { termination_date: today })}>Terminate</button>
                        <button type="button" className="eq-inline-btn danger" onClick={() => handleDelete(grant.id)}>Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
                {!loading && grants.length === 0 && (
                  <tr>
                    <td colSpan="7">No grants created yet.</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="eq-data-card" style={{ marginTop: 18, padding: 0, boxShadow: 'none', border: '0', background: 'transparent' }}>
            <div className="eq-data-card-head" style={{ marginBottom: 10 }}>
              <h3>Upcoming Vesting Events</h3>
            </div>
            <div className="eq-table-wrap">
              <table className="eq-table">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Grant</th>
                    <th>Type</th>
                    <th>Units</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {vestingEvents.slice(0, 12).map((event) => (
                    <tr key={event.id}>
                      <td>{event.vest_date}</td>
                      <td>{event.grant_number || '—'}</td>
                      <td>{event.event_type}</td>
                      <td>{event.units}</td>
                      <td>{event.status}</td>
                    </tr>
                  ))}
                  {!loading && vestingEvents.length === 0 && (
                    <tr>
                      <td colSpan="5">No vesting events generated yet.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <section className="eq-screen premium-dashboard-shell">
        <div className="premium-shell-body">
          <div className="eq-data-card premium-panel">
            <div className="eq-data-card-head">
              <h3>Vesting Timeline</h3>
              <span className="eq-status-chip success">Schedule view</span>
            </div>
            <div className="eq-table-wrap">
              <table className="eq-table">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Grant</th>
                    <th>Units</th>
                    <th>Progress</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {vestingTimeline.map((event) => {
                    const width = Math.max(8, Math.round((Number(event.units || 0) / timelineMax) * 100));
                    return (
                      <tr key={event.id}>
                        <td>{event.vest_date}</td>
                        <td>{event.grant_number || '—'}</td>
                        <td>{event.units}</td>
                        <td>
                          <div className="premium-progress-track">
                            <div className="premium-progress-fill" style={{ width: `${width}%` }} />
                          </div>
                        </td>
                        <td>{event.status}</td>
                      </tr>
                    );
                  })}
                  {!loading && vestingTimeline.length === 0 && (
                    <tr>
                      <td colSpan="5">No vesting events generated yet.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>

          <div className="eq-data-card premium-panel">
            <div className="eq-data-card-head">
              <h3>Vesting Engine Notes</h3>
              <span className="eq-status-chip">Audit view</span>
            </div>
            <div className="eq-empty-state">
              <h4>Acceleration, termination, and exercise stay synchronized</h4>
              <p>The roadmap view now makes the vesting runway readable while preserving the grant engine and schedule controls below.</p>
            </div>
          </div>
        </div>
      </section>
    </section>
  );
};

export default VestingGrants;
