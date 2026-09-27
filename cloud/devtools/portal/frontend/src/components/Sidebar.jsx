export default function Sidebar() {
  return (
    <div className="w-64 bg-gray-900 text-white p-4">
      <h1 className="text-xl font-bold mb-6">ATCloud Portal</h1>

      <nav className="flex flex-col gap-4">
        <a href="/" className="hover:text-blue-400">Dashboard</a>
        <a href="/compute" className="hover:text-blue-400">Compute</a>
        <a href="/network" className="hover:text-blue-400">Network</a>
        <a href="/storage" className="hover:text-blue-400">Storage</a>
        <a href="/projects" className="hover:text-blue-400">Projects</a>
        <a href="/telemetry" className="hover:text-blue-400">Telemetry</a>
      </nav>
    </div>
  );
}
