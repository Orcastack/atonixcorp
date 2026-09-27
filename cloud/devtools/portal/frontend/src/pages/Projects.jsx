import { useEffect, useState } from "react";
import { apiGet } from "../api/client";

export default function Projects() {
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    apiGet("/projects")
      .then(setProjects)
      .catch(console.error);
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Projects</h2>

      <ul className="list-disc pl-6">
        {projects.map(p => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
    </div>
  );
}
