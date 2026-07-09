import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

type VersionSnapshot = {
  document_id: string;
  version: number;
  markdown: string;
  published_at: string;
  published_by: string;
};

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
    <main style={{ fontFamily: "system-ui, sans-serif", padding: "2rem", maxWidth: "48rem" }}>
      <p style={{ color: "#666" }}>
        Document {snapshot.document_id} · v{snapshot.version} · published by {snapshot.published_by}
      </p>
      <article className="markdown-preview">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{snapshot.markdown}</ReactMarkdown>
      </article>
    </main>
  );
}
