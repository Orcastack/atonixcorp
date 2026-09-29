"use client";
import "../../styles/OperationalRail.css";

export default function OperationalRail() {
  return (
    <section className="operational-rail">
      <div className="rail-wrap">

        <div className="rail-item">
          <h4>Cloud Control</h4>
          <p>Unified orchestration, identity, policy, and automation services.</p>
        </div>

        <div className="rail-item">
          <h4>Compute & HPC</h4>
          <p>High‑performance execution environments for demanding workloads.</p>
        </div>

        <div className="rail-item">
          <h4>Network & Data</h4>
          <p>Secure connectivity, governed storage, and durable data services.</p>
        </div>

        <div className="rail-item">
          <h4>Secure Operations</h4>
          <p>Delivery, protection, analytics, and operational intelligence.</p>
        </div>

      </div>
    </section>
  );
}
