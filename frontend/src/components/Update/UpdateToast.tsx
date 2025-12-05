import { X, Download, Clock } from 'lucide-react';

interface UpdateToastProps {
  version: string;
  onViewUpdate: () => void;
  onDismiss: () => void;
}

export function UpdateToast({ version, onViewUpdate, onDismiss }: UpdateToastProps) {
  const handleRemindLater = () => {
    localStorage.setItem('updateDismissedAt', Date.now().toString());
    localStorage.setItem('updateDismissedVersion', version);
    onDismiss();
  };

  return (
    <div className="fixed bottom-4 right-4 z-50 animate-slide-up">
      <div className="bg-purple-900/90 border border-purple-500 rounded-lg shadow-lg p-4 max-w-sm">
        <div className="flex items-start gap-3">
          <Download size={20} className="text-purple-300 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <p className="text-purple-100 font-medium">Update available: {version}</p>
            <p className="text-purple-300 text-sm mt-1">A new version is ready to install.</p>
            <div className="flex gap-2 mt-3">
              <button
                onClick={onViewUpdate}
                className="flex items-center gap-1 px-3 py-1.5 bg-purple-600 hover:bg-purple-500 rounded text-sm transition-colors"
              >
                View Update
              </button>
              <button
                onClick={handleRemindLater}
                className="flex items-center gap-1 px-3 py-1.5 text-purple-300 hover:text-white text-sm transition-colors"
              >
                <Clock size={14} />
                Later
              </button>
            </div>
          </div>
          <button
            onClick={onDismiss}
            className="text-purple-400 hover:text-white"
          >
            <X size={18} />
          </button>
        </div>
      </div>
    </div>
  );
}

export default UpdateToast;
