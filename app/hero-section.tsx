"use client";

import { useEffect, useState } from "react";

const slides = [
  { label: "OrcaCompute", title: "Compute systems for the work that defines tomorrow.", description: "High-performance environments engineered for scientific workloads, AI training, and uninterrupted enterprise operations.", image: "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1800&q=85", meta: "01 / High-performance compute", action: "Explore compute" },
  { label: "OrcaStorage", title: "Storage engineered for enduring intelligence.", description: "Resilient, governed data foundations that preserve access, performance, and sovereignty across global operations.", image: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1800&q=85", meta: "02 / Sovereign data systems", action: "Explore storage" },
  { label: "OrcaNetwork", title: "Networks built for absolute operational confidence.", description: "Secure connectivity and distributed coordination for systems where availability is a permanent requirement.", image: "https://images.unsplash.com/photo-1497366754035-f200968a6e72?auto=format&fit=crop&w=1800&q=85", meta: "03 / Global network infrastructure", action: "Explore networking" },
  { label: "OrcaHardware", title: "Hardware aligned to long-term infrastructure control.", description: "Energy-conscious physical systems, autonomous operations, and advanced cooling built for enterprise scale.", image: "https://images.unsplash.com/photo-1518432031352-d6fc5c10da5a?auto=format&fit=crop&w=1800&q=85", meta: "04 / Integrated hardware platforms", action: "Explore hardware" },
  { label: "RIDD", title: "Research that advances the infrastructure frontier.", description: "Applied research environments for artificial intelligence, energy systems, and the next era of autonomous engineering.", image: "https://images.unsplash.com/photo-1532094349884-543bc11b234d?auto=format&fit=crop&w=1800&q=85", meta: "05 / Research and innovation", action: "Explore research" },
];

export default function HeroSection() {
  const [activeSlide, setActiveSlide] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const slide = slides[activeSlide];

  useEffect(() => {
    if (isPaused) return;
    const interval = window.setInterval(() => setActiveSlide((current) => (current + 1) % slides.length), 7000);
    return () => window.clearInterval(interval);
  }, [isPaused]);

  return <section className="hero hero-slider" id="top" aria-roledescription="carousel" aria-label="AtonixCorp infrastructure capabilities" onMouseEnter={() => setIsPaused(true)} onMouseLeave={() => setIsPaused(false)}>
    <div className="hero-slider-image" style={{ backgroundImage: `url(${slide.image})` }} />
    <div className="hero-slider-overlay" />
    <div className="hero-slider-content">
      <p className="hero-slider-label">{slide.label}</p>
      <p className="hero-slider-count">{String(activeSlide + 1).padStart(2, "0")} <span>/ 05</span></p>
      <h1>{slide.title}</h1>
      <p className="hero-slider-description">{slide.description}</p>
      <a className="hero-slider-action" href="#divisions">{slide.action} <span aria-hidden="true">→</span></a>
    </div>
    <div className="hero-slider-footer">
      <p>{slide.meta}</p>
      <div className="hero-slider-controls" aria-label="Select infrastructure capability">
        {slides.map((item, index) => <button key={item.label} type="button" onClick={() => { setActiveSlide(index); setIsPaused(true); }} className={index === activeSlide ? "active" : ""} aria-label={`Show ${item.label}`} aria-current={index === activeSlide ? "true" : undefined}><span>{String(index + 1).padStart(2, "0")}</span>{item.label.replace("Orca", "")}</button>)}
      </div>
    </div>
  </section>;
}
