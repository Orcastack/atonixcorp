import { NextResponse } from "next/server";

const apiUrl = process.env.API_URL ?? "http://localhost:8080";

export async function POST(request: Request) {
  const payload = await request.json();
  const response = await fetch(`${apiUrl}/api/v1/contact/submit`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const data = await response.json().catch(() => ({ error: "The API returned an invalid response." }));
  return NextResponse.json(data, { status: response.status });
}
