const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export default function HomePage() {
  return (
    <main style={{ fontFamily: "system-ui, sans-serif", padding: "2rem" }}>
      <h1>Codencil</h1>
      <p>Write in markdown. Review in the margin. Publish when it&apos;s ready.</p>
      <p>
        API: <code>{apiUrl}</code>
      </p>
    </main>
  );
}
