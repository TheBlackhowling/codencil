import { DocumentReviewView } from "./document-review-view";
import type { VersionSnapshot } from "./review-api";

async function fetchVersion(id: string, version: string): Promise<VersionSnapshot> {
  const base = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  const res = await fetch(`${base}/documents/${id}/versions/${version}`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`version not found (${res.status})`);
  }
  return res.json();
}

export default async function DocumentVersionPage({
  params,
}: {
  params: Promise<{ id: string; version: string }>;
}) {
  const { id, version } = await params;
  const snapshot = await fetchVersion(id, version);

  return (
    <main style={{ fontFamily: "system-ui, sans-serif", padding: "2rem", maxWidth: "960px", margin: "0 auto" }}>
      <DocumentReviewView snapshot={snapshot} />
    </main>
  );
}
