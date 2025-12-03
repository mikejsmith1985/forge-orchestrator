import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';

function App() {
  return (
    <div className="flex h-screen bg-gray-950 overflow-hidden">
      <Sidebar />
      <MainContent>
        <div className="flex flex-col items-center justify-center h-full text-center space-y-4">
          <div className="p-4 rounded-full bg-blue-500/10 border border-blue-500/20">
            <div className="w-16 h-16 rounded-full bg-blue-500 animate-pulse" />
          </div>
          <h2 className="text-2xl font-bold text-white">Welcome to Forge Vision</h2>
          <p className="text-gray-400 max-w-md">
            Select a flow from the sidebar to begin orchestrating your agents.
          </p>
        </div>
      </MainContent>
    </div>
  );
}

export default App;
