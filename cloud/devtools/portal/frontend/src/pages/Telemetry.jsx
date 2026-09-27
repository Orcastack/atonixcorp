import { useEffect, useState } from "react";
import { apiGet } from "../api/client";

export default function Telemetry() {
  const [snap, setSnap] = useState(null);

  useEffect(() => {
    apiGet("/telemetry")
      .then(setSnap)
      .catch(console.error);
  }, []);

  if (!snap) return <div>Loading...</div>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Telemetry</h2>

      <div className="grid grid-cols-3 gap-6">
        <div className="p-4 bg-gray-100 rounded">CPU: {snap.cpu_usage}%</div>
        <div className="p-4 bg-gray-100 rounded">Memory: {snap.memory_usage}%</div>
        <div className="p-4 bg-gray-100 rounded">Disk: {snap.disk_usage}%</div>
      </div>
    </div>
  );
}
