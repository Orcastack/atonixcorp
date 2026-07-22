import React from 'react';
import { buildBalancedMetricOrder } from '../../../utils/dashboardMetrics';
import '../../../styles/premiumDashboards.css';
import './EquityModuleScreen.css';
import './EquityCrudModuleScreen.css';

const displayValue = (value) => {
  if (value === null || value === undefined || value === '') {
    return '—';
  }
  if (typeof value === 'boolean') {
    return value ? 'Yes' : 'No';
  }
  return value;
};

const renderField = (field, value, onChange) => {
  if (field.type === 'textarea') {
    return (
      <textarea
        className="eq-form-textarea"
        rows={field.rows || 4}
        value={value ?? ''}
        onChange={(event) => onChange(field.key, event.target.value, field.type)}
        placeholder={field.placeholder || ''}
      />
    );
  }

  if (field.type === 'select') {
    return (
      <select
        className="eq-form-select"
        value={value ?? ''}
        onChange={(event) => onChange(field.key, event.target.value, field.type)}
      >
        <option value="">{field.placeholder || 'Select an option'}</option>
        {(field.options || []).map((option) => (
          <option key={option.value} value={option.value}>{option.label}</option>
        ))}
      </select>
    );
  }

  if (field.type === 'checkbox') {
    return (
      <label className="eq-form-checkbox">
        <input
          type="checkbox"
          checked={Boolean(value)}
          onChange={(event) => onChange(field.key, event.target.checked, field.type)}
        />
        <span>{field.checkboxLabel || field.label}</span>
      </label>
    );
  }

  return (
    <input
      className="eq-form-input"
      type={field.type || 'text'}
      value={value ?? ''}
      onChange={(event) => onChange(field.key, event.target.value, field.type)}
      placeholder={field.placeholder || ''}
      min={field.min}
      step={field.step}
    />
  );
};

const EquityCrudModuleScreen = ({
  title,
  description,
  metrics,
  records,
  columns,
  emptyTitle,
  emptyBody,
  formTitle,
  formDescription,
  formFields,
  formState,
  onFieldChange,
  onSubmit,
  onCancel,
  onEdit,
  onDelete,
  editingLabel,
  saving,
  loading,
  error,
}) => (
  <section className="eq-screen premium-dashboard-shell">
    <div className="premium-shell-body">
    <div className="eq-screen-hero premium-hero">
      <div>
        <div className="premium-hero-kicker">Equity operating surface</div>
        <h2 className="premium-hero-title">{title}</h2>
        <p className="premium-hero-text">{description}</p>
      </div>
      <div className="eq-screen-banner premium-metric-card">
        Every action in this module is designed for auditability, ownership clarity, and institutional reporting.
      </div>
    </div>

    <div className="eq-metric-grid premium-grid-3">
      {buildBalancedMetricOrder(metrics, metrics.length).map((metric) => (
        <article key={metric.label} className="eq-metric-card premium-metric-card">
          <span className="eq-metric-label premium-metric-label">{metric.label}</span>
          <strong className="eq-metric-value premium-metric-value">{metric.value}</strong>
          <span className="eq-metric-note premium-metric-note">{metric.note}</span>
        </article>
      ))}
    </div>

    <div className="eq-crud-layout">
      <aside className="eq-data-card eq-form-card premium-panel">
        <div className="eq-data-card-head">
          <h3>{editingLabel || formTitle}</h3>
          {editingLabel && (
            <button type="button" className="eq-inline-btn secondary" onClick={onCancel}>
              Cancel edit
            </button>
          )}
        </div>
        <p className="eq-form-copy">{formDescription}</p>
        <form className="eq-form-grid" onSubmit={onSubmit}>
          {formFields.map((field) => (
            <label key={field.key} className={`eq-form-field${field.fullWidth ? ' full' : ''}`}>
              {field.type !== 'checkbox' && <span className="eq-form-label">{field.label}</span>}
              {renderField(field, formState[field.key], onFieldChange)}
            </label>
          ))}
          <div className="eq-form-actions">
            <button type="submit" className="eq-inline-btn primary" disabled={saving}>
              {saving ? 'Saving…' : editingLabel ? 'Save changes' : 'Create record'}
            </button>
          </div>
        </form>
      </aside>

      <div className="eq-data-card premium-panel">
        <div className="eq-data-card-head">
          <h3>Live Register</h3>
          {loading && <span className="eq-status-chip">Syncing</span>}
          {!loading && !error && <span className="eq-status-chip success">Live</span>}
          {error && <span className="eq-status-chip danger">Attention</span>}
        </div>

        {error && <div className="eq-error-banner">{error}</div>}

        {!loading && records.length === 0 ? (
          <div className="eq-empty-state">
            <h4>{emptyTitle}</h4>
            <p>{emptyBody}</p>
          </div>
        ) : (
          <div className="eq-table-wrap">
            <table className="eq-table">
              <thead>
                <tr>
                  {columns.map((column) => <th key={column.key}>{column.label}</th>)}
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {records.map((record) => (
                  <tr key={record.id}>
                    {columns.map((column) => (
                      <td key={column.key}>
                        {typeof column.render === 'function' ? column.render(record) : displayValue(record[column.key])}
                      </td>
                    ))}
                    <td>
                      <div className="eq-table-actions">
                        <button type="button" className="eq-inline-btn" onClick={() => onEdit(record)}>Edit</button>
                        <button type="button" className="eq-inline-btn danger" onClick={() => onDelete(record)}>Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
    <footer className="premium-footer">
      <div className="premium-footer-group">
        <span className="premium-status-pill">Compliance current</span>
        <span className="premium-footer-note">AtonixCorp equity dashboard</span>
      </div>
      <div className="premium-footer-group">
        <span className="premium-footer-note">Ownership, vesting, governance, and auditability in one surface</span>
      </div>
    </footer>
    </div>
  </section>
);

export default EquityCrudModuleScreen;
