"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {
  Anchor,
  VersionSnapshot,
  createAnchor,
  fetchAnchors,
  lineRangeFromSelection,
} from "./review-api";

type Props = {
  snapshot: VersionSnapshot;
};

export function DocumentReviewView({ snapshot }: Props) {
  const previewRef = useRef<HTMLDivElement>(null);
  const [anchors, setAnchors] = useState<Anchor[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [pendingSelection, setPendingSelection] = useState<{
    start: number;
    end: number;
    text: string;
  } | null>(null);
  const [commentBody, setCommentBody] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const loadAnchors = useCallback(async () => {
    try {
      const data = await fetchAnchors(snapshot.document_id, snapshot.version);
      setAnchors(data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load anchors");
    }
  }, [snapshot.document_id, snapshot.version]);

  useEffect(() => {
    void loadAnchors();
  }, [loadAnchors]);

  const handleMouseUp = () => {
    if (!previewRef.current) {
      return;
    }
    const range = lineRangeFromSelection(previewRef.current);
    if (range) {
      setPendingSelection(range);
      setCommentBody("");
    }
  };

  const submitAnchor = async () => {
    if (!pendingSelection || !commentBody.trim()) {
      return;
    }
    setSubmitting(true);
    try {
      await createAnchor(snapshot.document_id, snapshot.version, {
        start_line: pendingSelection.start,
        end_line: pendingSelection.end,
        quoted_text: pendingSelection.text,
        body: commentBody.trim(),
      });
      setPendingSelection(null);
      setCommentBody("");
      await loadAnchors();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to create anchor");
    } finally {
      setSubmitting(false);
    }
  };

  const lines = snapshot.markdown.split("\n");

  return (
    <div style={{ display: "grid", gridTemplateColumns: "1fr 280px", gap: "1.5rem", alignItems: "start" }}>
      <section>
        <p style={{ color: "#666", marginBottom: "1rem" }}>
          Document {snapshot.document_id} · v{snapshot.version} · published by {snapshot.published_by}
        </p>
        <p style={{ fontSize: "0.875rem", color: "#888", marginBottom: "1rem" }}>
          Select text in the preview to add a margin comment.
        </p>

        {pendingSelection && (
          <div
            style={{
              marginBottom: "1rem",
              padding: "0.75rem 1rem",
              border: "1px solid #f59e0b",
              borderRadius: "8px",
              background: "#fffbeb",
            }}
          >
            <p style={{ margin: "0 0 0.5rem", fontWeight: 600 }}>New comment on lines {pendingSelection.start}–{pendingSelection.end}</p>
            <blockquote style={{ margin: "0 0 0.75rem", color: "#444", borderLeft: "3px solid #f59e0b", paddingLeft: "0.75rem" }}>
              {pendingSelection.text}
            </blockquote>
            <textarea
              value={commentBody}
              onChange={(e) => setCommentBody(e.target.value)}
              placeholder="Write your comment…"
              rows={3}
              style={{ width: "100%", marginBottom: "0.5rem", padding: "0.5rem", fontFamily: "inherit" }}
            />
            <div style={{ display: "flex", gap: "0.5rem" }}>
              <button type="button" onClick={() => void submitAnchor()} disabled={submitting || !commentBody.trim()}>
                Add comment
              </button>
              <button type="button" onClick={() => setPendingSelection(null)} disabled={submitting}>
                Cancel
              </button>
            </div>
          </div>
        )}

        <div
          ref={previewRef}
          onMouseUp={handleMouseUp}
          className="markdown-preview"
          style={{
            border: "1px solid #e5e7eb",
            borderRadius: "8px",
            padding: "1rem 1.25rem",
            background: "#fff",
            userSelect: "text",
          }}
        >
          {lines.map((line, index) => (
            <div key={index} data-line={index + 1} style={{ minHeight: "1.25rem" }}>
              <ReactMarkdown remarkPlugins={[remarkGfm]} components={{ p: ({ children }) => <span>{children}</span> }}>
                {line || " "}
              </ReactMarkdown>
            </div>
          ))}
        </div>
      </section>

      <aside>
        <h2 style={{ fontSize: "1rem", marginTop: 0 }}>Comments</h2>
        {error && <p style={{ color: "#b91c1c" }}>{error}</p>}
        {anchors.length === 0 && !error && (
          <p style={{ color: "#888", fontSize: "0.875rem" }}>No comments yet.</p>
        )}
        <div style={{ display: "flex", flexDirection: "column", gap: "0.75rem" }}>
          {anchors.map((anchor) => (
            <article
              key={anchor.id}
              style={{
                border: "1px solid #e5e7eb",
                borderRadius: "8px",
                padding: "0.75rem",
                background: anchor.review_state === "resolved" ? "#f3f4f6" : "#fff",
              }}
            >
              <p style={{ margin: "0 0 0.25rem", fontSize: "0.75rem", color: "#6b7280" }}>
                Lines {anchor.start_line}–{anchor.end_line}
                {anchor.review_state === "resolved" ? " · resolved" : ""}
              </p>
              <blockquote style={{ margin: "0 0 0.5rem", fontSize: "0.8125rem", color: "#374151", borderLeft: "2px solid #d1d5db", paddingLeft: "0.5rem" }}>
                {anchor.quoted_text}
              </blockquote>
              {anchor.comments[0] && (
                <p style={{ margin: 0, fontSize: "0.875rem" }}>{anchor.comments[0].body}</p>
              )}
            </article>
          ))}
        </div>
      </aside>
    </div>
  );
}
