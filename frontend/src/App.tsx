import { useState, useEffect, useCallback } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';
import { ArchitectView } from './components/Architect/ArchitectView';
import { LedgerView } from './components/Ledger/LedgerView';
import { CommandDeck } from './components/Commands/CommandDeck';
import { KeyManagement } from './components/Settings/KeyManagement';
import FlowList from './components/Flows/FlowList';
import FlowEditor from './components/Flows/FlowEditor';
import { UpdateModal, UpdateToast } from './components/Update';

interface UpdateInfo {
  available: boolean;
  currentVersion: string;
  latestVersion?: string;
  releaseNotes?: string;
  downloadUrl?: string;
  assetName?: string;
  assetSize?: number;
  error?: string;
}

function App() {
  const [currentVersion, setCurrentVersion] = useState('');
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [showUpdateToast, setShowUpdateToast] = useState(false);
  const [showUpdateModal, setShowUpdateModal] = useState(false);

  const checkForUpdates = useCallback(async () => {
    try {
      // Get current version
      const versionRes = await fetch('/api/version');
      const versionData = await versionRes.json();
      setCurrentVersion(versionData.version || '');
      
      // Check for updates
      const res = await fetch('/api/update/check');
      const data: UpdateInfo = await res.json();
      
      setUpdateInfo(data);
      
      if (data.available && data.latestVersion) {
        // Check if user dismissed this version recently (within 24 hours)
        const dismissedAt = localStorage.getItem('updateDismissedAt');
        const dismissedVersion = localStorage.getItem('updateDismissedVersion');
        const dayInMs = 24 * 60 * 60 * 1000;
        
        const wasRecentlyDismissed = dismissedAt && 
          dismissedVersion === data.latestVersion &&
          (Date.now() - parseInt(dismissedAt, 10)) < dayInMs;
        
        if (!wasRecentlyDismissed) {
          setShowUpdateToast(true);
        }
      }
    } catch (err) {
      console.error('Failed to check for updates:', err);
    }
  }, []);

  useEffect(() => {
    // Check for updates on mount
    checkForUpdates();
    
    // Check every 30 minutes
    const interval = setInterval(checkForUpdates, 30 * 60 * 1000);
    
    return () => clearInterval(interval);
  }, [checkForUpdates]);

  const handleViewUpdate = () => {
    setShowUpdateToast(false);
    setShowUpdateModal(true);
  };

  return (
    <div className="flex h-screen bg-gray-950 overflow-hidden">
      <Sidebar 
        currentVersion={currentVersion} 
        hasUpdate={updateInfo?.available || false}
        onUpdateClick={() => setShowUpdateModal(true)}
      />
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

      {/* Update Toast Notification */}
      {showUpdateToast && updateInfo?.latestVersion && (
        <UpdateToast
          version={updateInfo.latestVersion}
          onViewUpdate={handleViewUpdate}
          onDismiss={() => setShowUpdateToast(false)}
        />
      )}

      {/* Update Modal */}
      <UpdateModal
        isOpen={showUpdateModal}
        onClose={() => setShowUpdateModal(false)}
        updateInfo={updateInfo}
        currentVersion={currentVersion}
      />
    </div>
  );
}

export default App;
