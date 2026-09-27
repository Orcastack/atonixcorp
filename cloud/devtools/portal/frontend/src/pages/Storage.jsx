import { useEffect, useState } from "react";
import { apiGet } from "../api/client";

export default function Storage() {
  const [vols, setVols] = useState([]);

  useEffect(() => {
    apiGet("/storage")
      .then(setVols)
      .catch(console.error);
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Volumes</h2>

      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th>ID</th>
            <th>Name</th>
            <th>Size</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          {vols.map(v => (
            <tr key={v.id}>
              <td>{v.id}</td>
              <td>{v.name}</td>
              <td>{v.size_gb} GB</td>
              <td>{v.state}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
