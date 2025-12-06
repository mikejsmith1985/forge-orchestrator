import { useState, useEffect, useCallback } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';
import { ArchitectView } from './components/Architect/ArchitectView';
import { LedgerView } from './components/Ledger/LedgerView';
import { CommandDeck } from './components/Commands/CommandDeck';
import { Settings } from './components/Settings';
import { Terminal } from './components/Terminal';
import FlowList from './components/Flows/FlowList';
import FlowEditor from './components/Flows/FlowEditor';
import { UpdateModal, UpdateToast } from './components/Update';
import { WelcomeModal } from './components/Welcome';
import { FeedbackModal } from './components/Feedback';
import { ToastProvider } from './lib/ToastContext';

// Initialize logger to capture console output
import './utils/logger';

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

interface WelcomeInfo {
  shown: boolean;
  currentVersion: string;
  lastVersion: string;
}

function App() {
  const [currentVersion, setCurrentVersion] = useState('');
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [showUpdateToast, setShowUpdateToast] = useState(false);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [showWelcomeModal, setShowWelcomeModal] = useState(false);
  const [showFeedbackModal, setShowFeedbackModal] = useState(false);

  const checkWelcome = useCallback(async () => {
    try {
      const res = await fetch('/api/welcome');
      if (!res.ok) return;
      
      const data: WelcomeInfo = await res.json();
      
      if (!data.shown) {
        setShowWelcomeModal(true);
        setCurrentVersion(data.currentVersion);
      }
    } catch (err) {
      console.error('Failed to check welcome status:', err);
    }
  }, []);

  const dismissWelcome = useCallback(async () => {
    setShowWelcomeModal(false);
    
    try {
      await fetch('/api/welcome', { method: 'POST' });
    } catch (err) {
      console.error('Failed to mark welcome as shown:', err);
    }
  }, []);

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
    // Use setTimeout to avoid synchronous setState in effect body
    const welcomeTimeout = setTimeout(() => {
      checkWelcome();
    }, 0);
    
    const updateTimeout = setTimeout(() => {
      checkForUpdates();
    }, 0);
    
    // Check for updates every 30 minutes
    const interval = setInterval(checkForUpdates, 30 * 60 * 1000);
    
    return () => {
      clearTimeout(welcomeTimeout);
      clearTimeout(updateTimeout);
      clearInterval(interval);
    };
  }, [checkWelcome, checkForUpdates]);

  const handleViewUpdate = () => {
    setShowUpdateToast(false);
    setShowUpdateModal(true);
  };

  return (
    <ToastProvider>
      <div className="flex h-screen bg-gray-950 overflow-hidden">
        <Sidebar 
          currentVersion={currentVersion} 
          hasUpdate={updateInfo?.available || false}
          onUpdateClick={() => setShowUpdateModal(true)}
          onFeedbackClick={() => setShowFeedbackModal(true)}
        />
        <MainContent>
          <Routes>
            <Route path="/" element={<Navigate to="/terminal" replace />} />
            <Route path="/terminal" element={<Terminal />} />
            <Route path="/architect" element={<ArchitectView />} />
            <Route path="/ledger" element={<LedgerView />} />
            <Route path="/commands" element={<CommandDeck />} />
            <Route path="/flows" element={<FlowList />} />
            <Route path="/flows/new" element={<FlowEditor />} />
            <Route path="/flows/:id" element={<FlowEditor />} />
            <Route path="/settings" element={<Settings />} />
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

        {/* Welcome Modal */}
        <WelcomeModal
          isOpen={showWelcomeModal}
          onClose={dismissWelcome}
          version={currentVersion || '1.1.0'}
        />

        {/* Feedback Modal */}
        <FeedbackModal
          isOpen={showFeedbackModal}
          onClose={() => setShowFeedbackModal(false)}
        />
      </div>
    </ToastProvider>
  );
}

export default App;
