"use client";

export default function RoadmapSection() {
  const roadmap = [
    ["I", "Foundation & Stabilization", "Core infrastructure, cloud systems, HPC environments, and governance frameworks."],
    ["II", "Infrastructure Expansion", "Multi-region cloud, HPC, AI, and energy systems at scale."],
    ["III", "Global Integration", "International enterprise, research, and institutional collaboration."],
    ["IV", "International Positioning", "Expanded developer ecosystems and media intelligence capabilities."],
    ["V", "Autonomous Infrastructure Era", "Self-governing cloud, HPC, and AI operational frameworks."],
    ["VI", "Sovereign Technology Era", "Global operational independence and enduring technological leadership."],
  ];

  return (
    <section className="roadmap" id="roadmap">
      <div className="wrap">
        <div className="roadmap-head">
          <div>
            <p className="eyebrow light">06 / Strategic roadmap</p>
            <h2>
              Thirty years.<br />
              <em>Six phases.</em>
            </h2>
          </div>

          <p>
            A future built in deliberate phases, advancing independence and infrastructure
            stability at each horizon.
          </p>
        </div>

        <div className="timeline">
          {roadmap.map(([phase, title, detail]) => (
            <article key={phase}>
              <span>Phase {phase}</span>
              <h3>{title}</h3>
              <p>{detail}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
