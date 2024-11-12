export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <div className="max-w-7xl mx-auto">
        <div className="p-4">{children}</div>
      </div>
    </>
  );
}
