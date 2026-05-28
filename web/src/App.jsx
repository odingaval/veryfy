import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Landing from './pages/Landing';
import VerifierDashboard from './pages/VerifierDashboard';
import IssuerDashboard from './pages/IssuerDashboard';
import { TopAppBar } from './components/ui/TopAppBar';
import { BottomNavBar } from './components/ui/BottomNavBar';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-surface flex flex-col">
        <TopAppBar />
        <main className="flex-1 pb-24 pt-16">
          <Routes>
            <Route path="/" element={<Landing />} />
            <Route path="/verify" element={<VerifierDashboard />} />
            <Route path="/issue" element={<IssuerDashboard />} />
          </Routes>
        </main>
        <BottomNavBar />
      </div>
    </Router>
  );
}

export default App;
