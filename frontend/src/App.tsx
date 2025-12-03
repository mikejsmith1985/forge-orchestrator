import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';
import { ArchitectView } from './components/Architect/ArchitectView';

function App() {
  return (
    <div className="flex h-screen bg-gray-950 overflow-hidden">
      <Sidebar />
      <MainContent>
        <ArchitectView />
      </MainContent>
    </div>
  );
}

export default App;
