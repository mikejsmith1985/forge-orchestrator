import { Routes, Route, Navigate } from 'react-router-dom';
import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';
import { ArchitectView } from './components/Architect/ArchitectView';
import { LedgerView } from './components/Ledger/LedgerView';
import { CommandDeck } from './components/Commands/CommandDeck';
import { KeyManagement } from './components/Settings/KeyManagement';
import FlowList from './components/Flows/FlowList';
import FlowEditor from './components/Flows/FlowEditor';

function App() {
  return (
    <div className="flex h-screen bg-gray-950 overflow-hidden">
      <Sidebar />
      <MainContent>
        <Routes>
          <Route path="/" element={<Navigate to="/architect" replace />} />
          <Route path="/architect" element={<ArchitectView />} />
          <Route path="/ledger" element={<LedgerView />} />
          <Route path="/commands" element={<CommandDeck />} />
          <Route path="/flows" element={<FlowList />} />
          <Route path="/flows/new" element={<FlowEditor />} />
          <Route path="/flows/:id" element={<FlowEditor />} />
          <Route path="/settings" element={<KeyManagement />} />
        </Routes>
      </MainContent>
    </div>
  );
}

export default App;
