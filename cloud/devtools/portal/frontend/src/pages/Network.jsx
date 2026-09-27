import { useEffect, useState } from "react";
import { apiGet } from "../api/client";

export default function Network() {
  const [nets, setNets] = useState([]);

  useEffect(() => {
    apiGet("/network")
      .then(setNets)
      .catch(console.error);
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Networks</h2>

      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th>ID</th>
            <th>Name</th>
            <th>CIDR</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          {nets.map(n => (
            <tr key={n.id}>
              <td>{n.id}</td>
              <td>{n.name}</td>
              <td>{n.cidr}</td>
              <td>{n.state}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
