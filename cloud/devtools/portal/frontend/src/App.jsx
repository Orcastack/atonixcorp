import { BrowserRouter, Routes, Route } from "react-router-dom";
import Sidebar from "./components/Sidebar";
import Topbar from "./components/Topbar";

import Dashboard from "./pages/Dashboard";
import Compute from "./pages/Compute";
import Network from "./pages/Network";
import Storage from "./pages/Storage";
import Projects from "./pages/Projects";
import Telemetry from "./pages/Telemetry";

export default function App() {
  return (
    <BrowserRouter>
      <div className="flex h-screen">
        <Sidebar />
        <div className="flex-1 flex flex-col">
          <Topbar />
          <div className="p-6 overflow-auto">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/compute" element={<Compute />} />
              <Route path="/network" element={<Network />} />
              <Route path="/storage" element={<Storage />} />
              <Route path="/projects" element={<Projects />} />
              <Route path="/telemetry" element={<Telemetry />} />
            </Routes>
          </div>
        </div>
      </div>
    </BrowserRouter>
  );
}
