export type Comment = {
  id: string;
  thread_id: string;
  author_id: string;
  body: string;
  created_at: string;
};

export type Anchor = {
  id: string;
  anchor_id: string;
  document_id: string;
  version: number;
  thread_id: string;
  start_line: number;
  end_line: number;
  quoted_text: string;
  anchor_status: string;
  review_state: "open" | "resolved";
  resolved_by?: string | null;
  resolved_at?: string | null;
  created_at: string;
  comments: Comment[];
};

export type VersionSnapshot = {
  document_id: string;
  version: number;
  markdown: string;
  published_at: string;
  published_by: string;
};

const apiBase = () => process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function fetchAnchors(documentId: string, version: number): Promise<Anchor[]> {
  const res = await fetch(`${apiBase()}/documents/${documentId}/versions/${version}/anchors`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`fetch anchors failed (${res.status})`);
  }
  return res.json();
}

export async function createAnchor(
  documentId: string,
  version: number,
  payload: {
    start_line: number;
    end_line: number;
    quoted_text: string;
    body: string;
    author_id?: string;
  }
): Promise<Anchor> {
  const res = await fetch(`${apiBase()}/documents/${documentId}/versions/${version}/anchors`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Dev-User-Id": "web-reviewer" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    throw new Error(`create anchor failed (${res.status})`);
  }
  return res.json();
}

export async function addComment(
  threadId: string,
  body: string,
  authorId = "web-reviewer"
): Promise<Comment> {
  const res = await fetch(`${apiBase()}/threads/${threadId}/comments`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Dev-User-Id": authorId },
    body: JSON.stringify({ body, author_id: authorId }),
  });
  if (!res.ok) {
    throw new Error(`add comment failed (${res.status})`);
  }
  return res.json();
}

export async function resolveAnchor(anchorId: string): Promise<void> {
  const res = await fetch(`${apiBase()}/anchors/${anchorId}/resolve`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Dev-User-Id": "web-reviewer" },
    body: JSON.stringify({ resolved_by: "web-reviewer" }),
  });
  if (!res.ok) {
    throw new Error(`resolve anchor failed (${res.status})`);
  }
}

export async function reopenAnchor(anchorId: string): Promise<void> {
  const res = await fetch(`${apiBase()}/anchors/${anchorId}/reopen`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });
  if (!res.ok) {
    throw new Error(`reopen anchor failed (${res.status})`);
  }
}

export function lineRangeFromSelection(container: HTMLElement): { start: number; end: number; text: string } | null {
  const selection = window.getSelection();
  if (!selection || selection.isCollapsed || !selection.rangeCount) {
    return null;
  }

  const range = selection.getRangeAt(0);
  if (!container.contains(range.commonAncestorContainer)) {
    return null;
  }

  const text = selection.toString().trim();
  if (!text) {
    return null;
  }

  const lines = new Set<number>();
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ELEMENT);
  let node: Node | null = walker.currentNode;
  while (node) {
    if (node instanceof HTMLElement && node.dataset.line && range.intersectsNode(node)) {
      lines.add(Number(node.dataset.line));
    }
    node = walker.nextNode();
  }

  if (lines.size === 0) {
    return null;
  }

  const sorted = [...lines].sort((a, b) => a - b);
  return { start: sorted[0], end: sorted[sorted.length - 1], text };
}
