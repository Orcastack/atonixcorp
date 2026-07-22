import React from 'react';
import { useEquity } from '../../../context/EquityContext';
import { buildBalancedMetricOrder } from '../../../utils/dashboardMetrics';
import '../../../styles/premiumDashboards.css';
import './EquityModuleScreen.css';

const formatNumber = (value) => new Intl.NumberFormat('en-US').format(Number(value || 0));
const formatCurrency = (value) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(Number(value || 0));

const EquityModuleScreen = ({
  title,
  description,
  metrics,
  records,
  columns,
  emptyTitle,
  emptyBody,
}) => {
  const { loading, error } = useEquity();
  const orderedMetrics = buildBalancedMetricOrder(metrics, metrics.length);

  return (
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
        {orderedMetrics.map((metric) => (
          <article key={metric.label} className="eq-metric-card premium-metric-card">
            <span className="eq-metric-label premium-metric-label">{metric.label}</span>
            <strong className="eq-metric-value premium-metric-value">{metric.format === 'currency' ? formatCurrency(metric.value) : formatNumber(metric.value)}</strong>
            <span className="eq-metric-note premium-metric-note">{metric.note}</span>
          </article>
        ))}
      </div>

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
                </tr>
              </thead>
              <tbody>
                {records.map((record) => (
                  <tr key={record.id}>
                    {columns.map((column) => (
                      <td key={column.key}>
                        {typeof column.render === 'function' ? column.render(record) : record[column.key] || '—'}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
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
};

export default EquityModuleScreen;
