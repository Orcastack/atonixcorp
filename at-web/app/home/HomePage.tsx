"use client";

import "../styles/HomePage.css";

// Top Ribbon
import TopHeader from "../sections/header/TopHeader";

// Main Header
import GlobalHeader from "../sections/header/GlobalHeader";

// Sections
import HeroSection from "../sections/hero/HeroSection";
import OperationalRail from "../sections/operational/OperationalRail";
import HeroCorporate from "../sections/corporate/HeroCorporate";
import Footer from "../sections/footer/Footer";

export default function HomePage() {
  return (
    <main className="home-page">

      <TopHeader />
      <GlobalHeader />

      <HeroSection />
      <OperationalRail />
      <HeroCorporate />

      <Footer />
    </main>
  );
}
