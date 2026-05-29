// App.tsx
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Landing from './pages/Landing';
import VerifierDashboard from './pages/VerifierDashboard';
import IssuerDashboard from './pages/IssuerDashboard';
import { Toaster } from 'react-hot-toast';

function App() {
  return (
    <Router>
      <Toaster position="top-right" />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/verify" element={<VerifierDashboard />} />
        <Route path="/issue" element={<IssuerDashboard />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;
