"use client";

import { useState } from "react";
import type { Anchor } from "./review-api";
import { addComment, reopenAnchor, resolveAnchor } from "./review-api";

type Props = {
  anchor: Anchor;
  onUpdated: () => Promise<void>;
};

export function AnchorThreadCard({ anchor, onUpdated }: Props) {
  const isResolved = anchor.review_state === "resolved";
  const [expanded, setExpanded] = useState(!isResolved);
  const [replyBody, setReplyBody] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const runAction = async (action: () => Promise<void>) => {
    setBusy(true);
    setError(null);
    try {
      await action();
      await onUpdated();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Action failed");
    } finally {
      setBusy(false);
    }
  };

  const firstComment = anchor.comments[0];
  const replyCount = Math.max(0, anchor.comments.length - 1);

  return (
    <article
      style={{
        border: "1px solid #e5e7eb",
        borderRadius: "8px",
        padding: "0.75rem",
        background: isResolved ? "#f9fafb" : "#fff",
        opacity: isResolved && !expanded ? 0.85 : 1,
      }}
    >
      <button
        type="button"
        onClick={() => setExpanded((v) => !v)}
        style={{
          all: "unset",
          cursor: "pointer",
          display: "block",
          width: "100%",
        }}
      >
        <p style={{ margin: "0 0 0.25rem", fontSize: "0.75rem", color: "#6b7280" }}>
          Lines {anchor.start_line}–{anchor.end_line}
          {isResolved ? " · resolved" : " · open"}
          {replyCount > 0 ? ` · ${replyCount} repl${replyCount === 1 ? "y" : "ies"}` : ""}
        </p>
        <blockquote
          style={{
            margin: "0 0 0.5rem",
            fontSize: "0.8125rem",
            color: "#374151",
            borderLeft: "2px solid #d1d5db",
            paddingLeft: "0.5rem",
            whiteSpace: "nowrap",
            overflow: "hidden",
            textOverflow: "ellipsis",
          }}
        >
          {anchor.quoted_text}
        </blockquote>
        {!expanded && firstComment && (
          <p style={{ margin: 0, fontSize: "0.875rem", color: "#4b5563" }}>{firstComment.body}</p>
        )}
      </button>

      {expanded && (
        <div style={{ marginTop: "0.5rem" }}>
          <ul style={{ listStyle: "none", margin: "0 0 0.75rem", padding: 0, display: "flex", flexDirection: "column", gap: "0.5rem" }}>
            {anchor.comments.map((comment) => (
              <li key={comment.id} style={{ fontSize: "0.875rem" }}>
                <strong style={{ color: "#374151" }}>{comment.author_id}</strong>
                <span style={{ color: "#9ca3af", marginLeft: "0.5rem", fontSize: "0.75rem" }}>
                  {new Date(comment.created_at).toLocaleString()}
                </span>
                <p style={{ margin: "0.25rem 0 0" }}>{comment.body}</p>
              </li>
            ))}
          </ul>

          <textarea
            value={replyBody}
            onChange={(e) => setReplyBody(e.target.value)}
            placeholder="Reply…"
            rows={2}
            disabled={busy}
            style={{ width: "100%", marginBottom: "0.5rem", padding: "0.5rem", fontFamily: "inherit", fontSize: "0.875rem" }}
          />

          <div style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
            <button
              type="button"
              disabled={busy || !replyBody.trim()}
              onClick={() =>
                void runAction(async () => {
                  await addComment(anchor.thread_id, replyBody.trim());
                  setReplyBody("");
                })
              }
            >
              Reply
            </button>
            {isResolved ? (
              <button
                type="button"
                disabled={busy}
                onClick={() => void runAction(async () => reopenAnchor(anchor.id))}
              >
                Reopen
              </button>
            ) : (
              <button
                type="button"
                disabled={busy}
                onClick={() => void runAction(async () => resolveAnchor(anchor.id))}
              >
                Resolve
              </button>
            )}
            <button type="button" disabled={busy} onClick={() => setExpanded(false)}>
              Collapse
            </button>
          </div>
        </div>
      )}

      {error && <p style={{ color: "#b91c1c", fontSize: "0.8125rem", margin: "0.5rem 0 0" }}>{error}</p>}
    </article>
  );
}
