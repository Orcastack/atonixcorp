import { useEffect, useState } from "react";
import { apiGet } from "../api/client";

export default function Compute() {
  const [items, setItems] = useState([]);

  useEffect(() => {
    apiGet("/compute")
      .then(setItems)
      .catch(console.error);
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Compute Instances</h2>

      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th>ID</th>
            <th>Name</th>
            <th>Plan</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          {items.map(i => (
            <tr key={i.id}>
              <td>{i.id}</td>
              <td>{i.name}</td>
              <td>{i.plan}</td>
              <td>{i.state}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
