import { useState, useEffect } from 'react';
import { Download, RefreshCw, ExternalLink, AlertTriangle, CheckCircle, Clock, ChevronDown, ChevronUp, History, X } from 'lucide-react';

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

interface ReleaseInfo {
  version: string;
  name: string;
  publishedAt: string;
  releaseNotes: string;
  downloadUrl: string;
  isCurrent: boolean;
}

interface UpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
  updateInfo: UpdateInfo | null;
  currentVersion: string;
}

export function UpdateModal({ isOpen, onClose, updateInfo, currentVersion }: UpdateModalProps) {
  const [isUpdating, setIsUpdating] = useState(false);
  const [updateStatus, setUpdateStatus] = useState<'downloading' | 'applying' | 'success' | 'error' | null>(null);
  const [errorMessage, setErrorMessage] = useState('');
  const [showVersions, setShowVersions] = useState(false);
  const [versions, setVersions] = useState<ReleaseInfo[]>([]);
  const [loadingVersions, setLoadingVersions] = useState(false);

  useEffect(() => {
    if (!isOpen) {
      setIsUpdating(false);
      setUpdateStatus(null);
      setErrorMessage('');
      setShowVersions(false);
    }
  }, [isOpen]);

  const fetchVersions = async () => {
    if (versions.length > 0) {
      setShowVersions(!showVersions);
      return;
    }
    
    setLoadingVersions(true);
    try {
      const res = await fetch('/api/update/versions');
      const data = await res.json();
      if (data.releases) {
        setVersions(data.releases);
        setShowVersions(true);
      }
    } catch (err) {
      console.error('Failed to fetch versions:', err);
    } finally {
      setLoadingVersions(false);
    }
  };

  const handleUpdate = async () => {
    setIsUpdating(true);
    setUpdateStatus('downloading');
    setErrorMessage('');

    try {
      const res = await fetch('/api/update/apply', { method: 'POST' });
      const data = await res.json();
      
      if (data.success) {
        setUpdateStatus('success');
        setTimeout(() => window.location.reload(), 2000);
      } else {
        setUpdateStatus('error');
        setErrorMessage(data.error || 'Unknown error occurred');
        setIsUpdating(false);
      }
    } catch (err) {
      setUpdateStatus('error');
      setErrorMessage(err instanceof Error ? err.message : 'Failed to connect to server');
      setIsUpdating(false);
    }
  };

  const handleRemindLater = () => {
    localStorage.setItem('updateDismissedAt', Date.now().toString());
    localStorage.setItem('updateDismissedVersion', updateInfo?.latestVersion || '');
    onClose();
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  };

  if (!isOpen) return null;

  const hasUpdate = updateInfo?.available;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-zinc-900 rounded-lg shadow-xl max-w-xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-4 border-b border-zinc-700">
          <h3 className="text-lg font-semibold flex items-center gap-2">
            <Download size={20} />
            Software Update
          </h3>
          <button 
            onClick={onClose} 
            disabled={isUpdating}
            className="text-zinc-400 hover:text-white disabled:opacity-50"
          >
            <X size={20} />
          </button>
        </div>

        <div className="p-5 space-y-5">
          {/* Current Version */}
          <div className="flex justify-between items-center p-3 bg-zinc-800 rounded-lg">
            <span className="text-zinc-400">Current Version</span>
            <span className="font-mono font-semibold">v{currentVersion}</span>
          </div>

          {hasUpdate ? (
            <>
              {/* Available Update */}
              <div className="flex justify-between items-center p-3 bg-purple-900/30 border border-purple-500 rounded-lg">
                <span className="text-purple-300">Available Update</span>
                <span className="font-mono font-semibold text-purple-200">
                  {updateInfo.latestVersion}
                </span>
              </div>

              {/* Warning Message */}
              <div className="flex gap-3 p-3 bg-orange-900/30 border border-orange-500 rounded-lg">
                <AlertTriangle size={20} className="text-orange-400 flex-shrink-0 mt-0.5" />
                <div className="text-sm text-orange-200">
                  <strong>Warning:</strong> Updating will restart the application. 
                  Save any unsaved work before updating.
                </div>
              </div>

              {/* Release Notes */}
              {updateInfo.releaseNotes && (
                <div>
                  <h4 className="text-sm text-zinc-400 mb-2">Release Notes</h4>
                  <div className="bg-zinc-950 border border-zinc-700 rounded-lg p-3 max-h-40 overflow-y-auto text-sm whitespace-pre-wrap text-zinc-300">
                    {updateInfo.releaseNotes}
                  </div>
                </div>
              )}

              {/* Update Status */}
              {updateStatus && (
                <div className={`flex items-center gap-3 p-3 rounded-lg ${
                  updateStatus === 'error' ? 'bg-red-900/30 border border-red-500' :
                  updateStatus === 'success' ? 'bg-green-900/30 border border-green-500' :
                  'bg-blue-900/30 border border-blue-500'
                }`}>
                  {updateStatus === 'downloading' && (
                    <>
                      <RefreshCw size={18} className="text-blue-400 animate-spin" />
                      <span className="text-blue-200">Downloading update...</span>
                    </>
                  )}
                  {updateStatus === 'success' && (
                    <>
                      <CheckCircle size={18} className="text-green-400" />
                      <span className="text-green-200">Update applied! Restarting...</span>
                    </>
                  )}
                  {updateStatus === 'error' && (
                    <>
                      <AlertTriangle size={18} className="text-red-400" />
                      <span className="text-red-200">{errorMessage}</span>
                    </>
                  )}
                </div>
              )}
            </>
          ) : (
            <div className="text-center py-8 text-zinc-400">
              <CheckCircle size={48} className="text-green-500 mx-auto mb-4" />
              <p>You're running the latest version!</p>
            </div>
          )}

          {/* Version History Section */}
          <div className="border-t border-zinc-700 pt-4">
            <button
              onClick={fetchVersions}
              disabled={loadingVersions}
              className="w-full flex items-center justify-center gap-2 p-3 border border-zinc-600 rounded-lg text-zinc-400 hover:text-white hover:border-zinc-500 transition-colors"
            >
              {loadingVersions ? (
                <RefreshCw size={16} className="animate-spin" />
              ) : (
                <History size={16} />
              )}
              {showVersions ? 'Hide' : 'Show'} Previous Versions
              {showVersions ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
            </button>

            {showVersions && versions.length > 0 && (
              <div className="mt-4 max-h-48 overflow-y-auto border border-zinc-700 rounded-lg">
                {versions.map((release, idx) => (
                  <div 
                    key={release.version}
                    className={`flex justify-between items-center p-3 ${
                      idx < versions.length - 1 ? 'border-b border-zinc-700' : ''
                    } ${release.isCurrent ? 'bg-green-900/20' : ''}`}
                  >
                    <div>
                      <span className={`font-mono font-semibold ${release.isCurrent ? 'text-green-400' : 'text-zinc-300'}`}>
                        {release.version}
                        {release.isCurrent && (
                          <span className="ml-2 text-xs bg-green-500 text-black px-2 py-0.5 rounded">
                            Current
                          </span>
                        )}
                      </span>
                      <div className="text-xs text-zinc-500 mt-1">
                        {formatDate(release.publishedAt)}
                      </div>
                    </div>
                    {release.downloadUrl && !release.isCurrent && (
                      <a
                        href={release.downloadUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm text-blue-400 hover:text-blue-300 flex items-center gap-1"
                      >
                        <Download size={14} />
                        Download
                      </a>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* GitHub Releases Link */}
          <div className="text-center pt-4 border-t border-zinc-700">
            <a 
              href="https://github.com/mikejsmith1985/forge-orchestrator/releases" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-sm text-blue-400 hover:text-blue-300 inline-flex items-center gap-1"
            >
              <ExternalLink size={14} />
              View all releases on GitHub
            </a>
          </div>
        </div>

        <div className="flex justify-end gap-3 p-4 border-t border-zinc-700">
          {hasUpdate ? (
            <>
              <button 
                onClick={handleRemindLater}
                disabled={isUpdating}
                className="flex items-center gap-2 px-4 py-2 bg-zinc-700 hover:bg-zinc-600 rounded-lg disabled:opacity-50 transition-colors"
              >
                <Clock size={16} />
                Remind Me Later
              </button>
              <button 
                onClick={handleUpdate}
                disabled={isUpdating}
                className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg disabled:opacity-50 transition-colors"
              >
                {isUpdating ? (
                  <>
                    <RefreshCw size={16} className="animate-spin" />
                    Updating...
                  </>
                ) : (
                  <>
                    <Download size={16} />
                    Update Now
                  </>
                )}
              </button>
            </>
          ) : (
            <button 
              onClick={onClose}
              className="px-4 py-2 bg-zinc-700 hover:bg-zinc-600 rounded-lg transition-colors"
            >
              Close
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

export default UpdateModal;
