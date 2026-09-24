"use client";

import Image from "next/image";
import { useEffect, useRef, useState } from "react";

type NavItem = {
  label: string;
  href?: string;
  items?: { label: string; href: string }[];
};

const navigation: NavItem[] = [
  { label: "Home", href: "#top" },
  { label: "About", items: [
    { label: "About AtonixCorp", href: "#about" }, { label: "Mission & Vision", href: "#mission" }, { label: "Leadership", href: "#about" }, { label: "Corporate Governance", href: "#sovereignty" }, { label: "Global Presence", href: "#roadmap" }, { label: "Operational Mandate", href: "#mission" },
  ] },
  { label: "Divisions", items: [
    { label: "OrcaStack - Cloud & AI", href: "#divisions" }, { label: "OrcaCompute - HPC & Energy", href: "#divisions" }, { label: "RIDD - Research & Innovation", href: "#divisions" }, { label: "PrimeSourceDaily - Media Intelligence", href: "#divisions" }, { label: "Technology Sovereignty", href: "#sovereignty" }, { label: "Cybersecurity & Threat Defense", href: "#systems" },
  ] },
  { label: "Technology", items: [
    { label: "Cloud Infrastructure", href: "#divisions" }, { label: "Artificial Intelligence", href: "#divisions" }, { label: "High-Performance Computing", href: "#divisions" }, { label: "Energy Systems", href: "#systems" }, { label: "Future Systems (10-30 Years)", href: "#roadmap" }, { label: "Autonomous Infrastructure", href: "#mission" },
  ] },
  { label: "Research", items: [
    { label: "Research Divisions", href: "#divisions" }, { label: "Innovation Pipeline", href: "#roadmap" }, { label: "Scientific Facilities", href: "#systems" }, { label: "Technical Publications", href: "#contact" }, { label: "Engineering Standards", href: "#sovereignty" },
  ] },
  { label: "Developers", items: [
    { label: "Developer Hub", href: "#systems" }, { label: "Documentation", href: "#contact" }, { label: "API Reference", href: "#contact" }, { label: "SDK Downloads", href: "#contact" }, { label: "Developer Tools", href: "#systems" }, { label: "Open Innovation Programs", href: "#contact" },
  ] },
  { label: "Media Center", items: [
    { label: "Newsroom", href: "#contact" }, { label: "Press Releases", href: "#contact" }, { label: "Corporate Announcements", href: "#contact" }, { label: "Media Kit", href: "#contact" }, { label: "PrimeSourceDaily Portal", href: "#divisions" },
  ] },
  { label: "Contact", items: [
    { label: "Contact Us", href: "#contact" }, { label: "Support Center", href: "mailto:dev@atonixcorp.com" }, { label: "Enterprise Solutions", href: "mailto:partners@atonixcorp.com" }, { label: "Partnerships", href: "mailto:partners@atonixcorp.com" }, { label: "Careers", href: "mailto:info@atonixcorp.com" },
  ] },
];

export default function GlobalHeader() {
  const [activeMenu, setActiveMenu] = useState<string | null>(null);
  const [mobileOpen, setMobileOpen] = useState(false);
  const headerRef = useRef<HTMLElement>(null);

  useEffect(() => {
    const closeOnOutsideClick = (event: MouseEvent) => {
      if (!headerRef.current?.contains(event.target as Node)) setActiveMenu(null);
    };
    document.addEventListener("mousedown", closeOnOutsideClick);
    return () => document.removeEventListener("mousedown", closeOnOutsideClick);
  }, []);

  const closeAll = () => { setActiveMenu(null); setMobileOpen(false); };

  return <header className="global-header" ref={headerRef}>
    <div className="global-header-inner">
      <a className="global-brand" href="#top" aria-label="AtonixCorp home" onClick={closeAll}><Image src="/atonixcorp-icon.png" alt="" width={34} height={34} priority /><span>ATONIX<span>CORP</span></span></a>
      <nav className="global-nav" aria-label="Global navigation">
        {navigation.map((item) => item.items ? <div className="nav-menu" key={item.label}>
          <button className={activeMenu === item.label ? "nav-trigger open" : "nav-trigger"} type="button" aria-expanded={activeMenu === item.label} onClick={() => setActiveMenu(activeMenu === item.label ? null : item.label)}>{item.label}<span aria-hidden="true">⌄</span></button>
          <div className={activeMenu === item.label ? "dropdown visible" : "dropdown"}>{item.items.map((subItem) => <a href={subItem.href} key={subItem.label} onClick={closeAll}>{subItem.label}<span aria-hidden="true">↗</span></a>)}</div>
        </div> : <a className="nav-trigger home-link" href={item.href} key={item.label}>{item.label}</a>)}
      </nav>
      <button className={mobileOpen ? "menu-toggle open" : "menu-toggle"} type="button" aria-expanded={mobileOpen} aria-label="Toggle navigation menu" onClick={() => setMobileOpen(!mobileOpen)}><span /><span /></button>
    </div>
    <nav className={mobileOpen ? "mobile-nav visible" : "mobile-nav"} aria-label="Mobile global navigation">
      {navigation.map((item) => <div key={item.label} className="mobile-nav-group">{item.href ? <a href={item.href} onClick={closeAll}>{item.label}</a> : <><span>{item.label}</span>{item.items?.map((subItem) => <a href={subItem.href} key={subItem.label} onClick={closeAll}>{subItem.label}</a>)}</>}</div>)}
    </nav>
  </header>;
}
