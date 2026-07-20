import React, { useEffect, useState } from 'react';
import { FiAlertCircle, FiCheckCircle, FiFileText, FiShield } from 'react-icons/fi';

import { Card, PageHeader, Table } from '../../components/ui';
import { governanceAmendmentsAPI, governancePoliciesAPI } from '../../services/api';

const parseList = (response) => response.data.results || response.data || [];
const thresholds = [
  ['Standard', '60%'], ['Operational', '70%'], ['Ethical / security', '75%'],
  ['Constitutional', '80%'], ['Sovereignty', '90%'], ['Emergency', '75%'],
];

export default function GovernanceCenter() {
  const [policies, setPolicies] = useState([]);
  const [amendments, setAmendments] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [policyResponse, amendmentResponse] = await Promise.all([
          governancePoliciesAPI.getAll(),
          governanceAmendmentsAPI.getAll(),
        ]);
        setPolicies(parseList(policyResponse));
        setAmendments(parseList(amendmentResponse));
      } catch (requestError) {
        setError(requestError.response?.data?.detail || 'Governance records could not be loaded.');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const openVotes = amendments.filter((item) => item.status === 'voting').length;
  const reviewedPolicies = policies.filter((item) => item.status === 'active').length;

  return (
    <div className="module-page governance-center">
      <PageHeader
        title="Governance Center"
        subtitle="Control policy editions, amendment evidence, and verified voting records across the organization."
      />
      {error ? <div className="governance-center__error"><FiAlertCircle aria-hidden="true" />{error}</div> : null}
      <div className="governance-center__metrics">
        <Card className="stat-card"><FiFileText aria-hidden="true" /><div><span>Policy editions</span><strong>{policies.length}</strong></div></Card>
        <Card className="stat-card"><FiCheckCircle aria-hidden="true" /><div><span>Active policies</span><strong>{reviewedPolicies}</strong></div></Card>
        <Card className="stat-card"><FiShield aria-hidden="true" /><div><span>Open votes</span><strong>{openVotes}</strong></div></Card>
      </div>
      <div className="governance-center__grid">
        <Card title="Policy Register">
          {loading ? <p className="governance-center__empty">Loading policy editions...</p> : (
            <Table
              columns={[
                { key: 'policy_code', label: 'Code' },
                { key: 'title', label: 'Policy' },
                { key: 'edition', label: 'Edition' },
                { key: 'status', label: 'Status' },
                { key: 'next_review_date', label: 'Next review', render: (value) => value || 'Not scheduled' },
              ]}
              data={policies}
              emptyMessage="No policy editions have been registered for this organization."
            />
          )}
        </Card>
        <Card title="Voting Thresholds">
          <div className="governance-center__thresholds">
            {thresholds.map(([label, value]) => <div key={label}><span>{label}</span><strong>{value}</strong></div>)}
          </div>
        </Card>
      </div>
      <Card title="Amendment Docket">
        {loading ? <p className="governance-center__empty">Loading amendment dossiers...</p> : (
          <Table
            columns={[
              { key: 'amendment_number', label: 'Reference' },
              { key: 'title', label: 'Amendment' },
              { key: 'amendment_type', label: 'Class' },
              { key: 'status', label: 'Stage' },
              { key: 'required_approval_percent', label: 'Threshold', render: (value) => `${value}%` },
              { key: 'vote_summary', label: 'Approval', render: (value) => value ? `${value.approval_percent}%` : '0%' },
            ]}
            data={amendments}
            emptyMessage="No amendment dossiers are currently registered."
          />
        )}
      </Card>
      <Card title="Edition 1.1 Correction Register">
        <ul className="governance-center__register">
          <li>Correct the document label to “Governance” while preserving the original supplied filename as source evidence.</li>
          <li>Retain the leadership-tenure clause once; remove the duplicated chapter in the next published edition.</li>
          <li>Restore sequential clause references for the succession and presidential-council material.</li>
          <li>Correct the “161.5?Lifetime Salary” typographic artifact and related spacing defects.</li>
        </ul>
      </Card>
    </div>
  );
}