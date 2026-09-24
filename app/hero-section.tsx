export default function HeroSection() {
  return <section className="hero hero-corporate" id="top">
    <div className="hero-campus" role="img" aria-label="Solar-powered AtonixCorp engineering campus">
      <div className="hero-campus-overlay" />
      <div className="hero-corporate-copy">
        <p className="hero-label">AtonixCorp global engineering systems</p>
        <h1>Engineering the Future of <span>Global Infrastructure</span></h1>
        <p className="hero-corporate-intro">Building secure, efficient, and scalable systems that power the world&apos;s most critical digital and physical environments.</p>
        <a className="hero-corporate-button" href="#divisions">Explore AtonixCorp Technologies <span aria-hidden="true">→</span></a>
      </div>
      <div className="hero-campus-caption"><span>ATX-01</span><strong>Solar integrated systems campus</strong><span>Global operations</span></div>
    </div>
  </section>;
}
