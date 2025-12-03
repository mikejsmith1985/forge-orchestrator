import { useState } from 'react';
import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';
import { ArchitectView } from './components/Architect/ArchitectView';
import { LedgerView } from './components/Ledger/LedgerView';
import { CommandDeck } from './components/Commands/CommandDeck';

function App() {
  // Educational Comment: We use state to manage client-side routing.
  // This is a simple implementation; for larger apps, use react-router.
  const [view, setView] = useState<'architect' | 'ledger' | 'commands'>('architect');

  return (
    <div className="flex h-screen bg-gray-950 overflow-hidden">
      <Sidebar currentView={view} onViewChange={setView} />
      <MainContent>
        {view === 'architect' && <ArchitectView />}
        {view === 'ledger' && <LedgerView />}
        {view === 'commands' && <CommandDeck />}
      </MainContent>
    </div>
  );
}

export default App;
