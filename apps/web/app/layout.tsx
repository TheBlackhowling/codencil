import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Codencil",
  description: "Write in markdown. Review in the margin. Publish when it's ready.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
