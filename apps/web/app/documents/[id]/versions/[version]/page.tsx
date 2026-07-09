import { DocumentReviewView } from "./document-review-view";
import type { VersionSnapshot, VersionSummary } from "./review-api";

const base = () => process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

async function fetchVersion(id: string, version: string): Promise<VersionSnapshot> {
  const res = await fetch(`${base()}/documents/${id}/versions/${version}`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`version not found (${res.status})`);
  }
  return res.json();
}

async function fetchVersions(id: string): Promise<VersionSummary[]> {
  const res = await fetch(`${base()}/documents/${id}/versions`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`versions not found (${res.status})`);
  }
  return res.json();
}

export default async function DocumentVersionPage({
  params,
}: {
  params: Promise<{ id: string; version: string }>;
}) {
  const { id, version } = await params;
  const [snapshot, versions] = await Promise.all([
    fetchVersion(id, version),
    fetchVersions(id),
  ]);

  return (
    <main style={{ fontFamily: "system-ui, sans-serif", padding: "2rem", maxWidth: "960px", margin: "0 auto" }}>
      <DocumentReviewView snapshot={snapshot} versions={versions} />
    </main>
  );
}
