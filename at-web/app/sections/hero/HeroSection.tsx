"use client";
import "../../styles/HeroSection.css";

export default function HeroSection() {
  return (
    <section className="enterprise-hero">
      <div className="enterprise-hero-wrap">

        {/* Left Content */}
        <div className="enterprise-hero-content">
          <p className="enterprise-hero-index">00 / Introduction</p>

          <h1 className="enterprise-hero-title">
            The Enterprise Cloud Platform<br />
            for Modern Infrastructure
          </h1>

          <p className="enterprise-hero-subtitle">
            Built for scale. Engineered for reliability. Designed for mission‑critical systems.
          </p>

          <p className="enterprise-hero-description">
            AtonixCorp delivers cloud, compute, network, data, and secure‑operations
            capabilities that power enterprise workloads, scientific research,
            critical infrastructure, and next‑generation digital systems.
          </p>

          <a href="#divisions" className="enterprise-hero-action">
            Explore AtonixCorp Platform →
          </a>
        </div>

        {/* Right Illustration */}
        <div className="enterprise-hero-graphic">
          <img
            src="/heroimage/atheroimage.png"
            alt="AtonixCorp orbital diagram"
            className="hero-image"
          />
        </div>

      </div>
    </section>
  );
}
