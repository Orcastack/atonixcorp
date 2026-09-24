"use client";

import { FormEvent, useState } from "react";

type SubmissionState = "idle" | "submitting" | "success" | "error";

export default function ContactForm() {
  const [state, setState] = useState<SubmissionState>("idle");

  async function submitContact(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setState("submitting");

    const form = event.currentTarget;
    const payload = Object.fromEntries(new FormData(form).entries());
    const response = await fetch("/api/contact", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (response.ok) {
      form.reset();
      setState("success");
      return;
    }

    setState("error");
  }

  return <form className="contact-form" onSubmit={submitContact}>
    <div className="contact-form-row">
      <label>Name<input name="name" required autoComplete="name" /></label>
      <label>Email<input name="email" type="email" required autoComplete="email" /></label>
    </div>
    <label>Subject<input name="subject" required /></label>
    <label>Message<textarea name="message" required rows={4} /></label>
    <button type="submit" disabled={state === "submitting"}>{state === "submitting" ? "Submitting..." : "Submit inquiry"}<span aria-hidden="true">→</span></button>
    {state === "success" && <p className="form-status success" role="status">Your inquiry has been received.</p>}
    {state === "error" && <p className="form-status error" role="alert">Unable to submit your inquiry. Please try again.</p>}
  </form>;
}
