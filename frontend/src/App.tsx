import "./App.css";
import { Home } from "./home";
import { Stats } from "./stats";
import { Routes, Route } from "react-router-dom";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/stats" element={<Stats />} />
    </Routes>
  );
}

export default App;
