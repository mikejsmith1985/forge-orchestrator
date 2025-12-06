import { useState, useEffect } from 'react';
import { X, Camera, Github, Settings, Loader2, ExternalLink, Trash2, Check } from 'lucide-react';
import { getLogs } from '../../utils/logger';

interface FeedbackModalProps {
    isOpen: boolean;
    onClose: () => void;
}

interface ScreenshotData {
    dataUrl: string;
    timestamp: Date;
}

export function FeedbackModal({ isOpen, onClose }: FeedbackModalProps) {
    const [description, setDescription] = useState('');
    const [screenshots, setScreenshots] = useState<ScreenshotData[]>([]);
    const [isCapturing, setIsCapturing] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [showSetup, setShowSetup] = useState(false);
    const [githubToken, setGithubToken] = useState('');
    const [status, setStatus] = useState<{ type: 'info' | 'error' | 'success'; message: string } | null>(null);
    const [createdIssueUrl, setCreatedIssueUrl] = useState<string | null>(null);

    useEffect(() => {
        const savedToken = localStorage.getItem('forge_github_token');
        if (savedToken) {
            setGithubToken(savedToken);
            setShowSetup(false);
        } else {
            setShowSetup(true);
        }
    }, []);

    useEffect(() => {
        if (!isOpen) return;

        // Reset state when modal opens
        const savedToken = localStorage.getItem('forge_github_token');
        if (savedToken) {
            setGithubToken(savedToken);
            setShowSetup(false);
        } else {
            setShowSetup(true);
        }

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') {
                onClose();
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    const handleSaveSettings = () => {
        if (!githubToken.trim()) {
            setStatus({ type: 'error', message: 'GitHub Token is required' });
            return;
        }
        localStorage.setItem('forge_github_token', githubToken.trim());
        setShowSetup(false);
        setStatus(null);
    };

    const handleCapture = async () => {
        setIsCapturing(true);
        try {
            // Hide the modal temporarily
            const modal = document.querySelector('[data-testid="feedback-modal"]');
            if (modal instanceof HTMLElement) {
                modal.style.visibility = 'hidden';
            }

            // Dynamic import of html2canvas
            const { default: html2canvas } = await import('html2canvas');
            
            const canvas = await html2canvas(document.body, {
                allowTaint: true,
                useCORS: true,
                logging: false,
            });

            // Show the modal again
            if (modal instanceof HTMLElement) {
                modal.style.visibility = 'visible';
            }

            const dataUrl = canvas.toDataURL('image/png');
            setScreenshots(prev => [...prev, { dataUrl, timestamp: new Date() }]);

            // Also copy to clipboard
            canvas.toBlob(async (blob) => {
                if (blob) {
                    try {
                        await navigator.clipboard.write([
                            new ClipboardItem({ 'image/png': blob })
                        ]);
                    } catch (err) {
                        console.error('Clipboard copy failed:', err);
                    }
                }
            });
        } catch (err) {
            console.error('Screenshot failed:', err);
            setStatus({ type: 'error', message: 'Failed to capture screenshot' });
            
            // Make sure modal is visible again
            const modal = document.querySelector('[data-testid="feedback-modal"]');
            if (modal instanceof HTMLElement) {
                modal.style.visibility = 'visible';
            }
        } finally {
            setIsCapturing(false);
        }
    };

    const handleRemoveScreenshot = (index: number) => {
        setScreenshots(prev => prev.filter((_, i) => i !== index));
    };

    const uploadToGithub = async (base64Image: string): Promise<string> => {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const filename = `feedback-${timestamp}.png`;
        const content = base64Image.split(',')[1];

        const res = await fetch(`https://api.github.com/repos/mikejsmith1985/forge-orchestrator/contents/feedback-screenshots/${filename}`, {
            method: 'PUT',
            headers: {
                'Authorization': `token ${githubToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: 'upload feedback screenshot',
                content: content,
            })
        });

        if (!res.ok) {
            if (res.status === 401) {
                throw new Error('Invalid GitHub token - please check your token and try again');
            }
            if (res.status === 403) {
                throw new Error('GitHub token lacks permissions - ensure it has "public_repo" scope');
            }
            if (res.status === 404) {
                throw new Error('Repository not found - check your token has access to public repositories');
            }
            throw new Error(`Failed to upload screenshot (HTTP ${res.status})`);
        }

        const data = await res.json();
        return data.content.download_url;
    };

    const createIssue = async (imageUrls: string[]): Promise<{ html_url: string; number: number }> => {
        const title = `Feedback: ${description.substring(0, 50)}${description.length > 50 ? '...' : ''}`;
        let body = `**Description**\n${description}\n\n`;

        if (imageUrls.length > 0) {
            body += `**Screenshot${imageUrls.length > 1 ? 's' : ''}**\n`;
            imageUrls.forEach((url) => {
                body += `<img src="${url}">\n\n`;
            });
        }

        body += `**Environment**\n- User Agent: ${navigator.userAgent}\n- Time: ${new Date().toISOString()}\n\n`;

        const logs = getLogs();
        if (logs) {
            body += `<details>\n<summary>Application Logs</summary>\n\n\`\`\`\n${logs}\n\`\`\`\n</details>`;
        }

        const res = await fetch(`https://api.github.com/repos/mikejsmith1985/forge-orchestrator/issues`, {
            method: 'POST',
            headers: {
                'Authorization': `token ${githubToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title, body, labels: ['feedback', 'user-submitted'] })
        });

        if (!res.ok) {
            if (res.status === 401) {
                throw new Error('Invalid GitHub token - please update your token in settings');
            }
            if (res.status === 403) {
                throw new Error('GitHub token lacks permissions - ensure it has "public_repo" scope');
            }
            throw new Error(`Failed to create issue (HTTP ${res.status})`);
        }

        return await res.json();
    };

    const handleSubmit = async () => {
        if (!description.trim()) return;

        setIsSubmitting(true);
        setStatus({ type: 'info', message: 'Processing...' });

        try {
            const imageUrls: string[] = [];

            // Upload screenshots
            for (let i = 0; i < screenshots.length; i++) {
                setStatus({ type: 'info', message: `Uploading screenshot ${i + 1}/${screenshots.length}...` });
                try {
                    const url = await uploadToGithub(screenshots[i].dataUrl);
                    imageUrls.push(url);
                } catch (err) {
                    console.error('Screenshot upload failed:', err);
                    // Continue without this screenshot
                }
            }

            setStatus({ type: 'info', message: 'Creating GitHub issue...' });
            const issue = await createIssue(imageUrls);

            setCreatedIssueUrl(issue.html_url);
            setStatus({ type: 'success', message: `Issue #${issue.number} created!` });

            // Close after delay
            setTimeout(() => {
                onClose();
                setDescription('');
                setScreenshots([]);
                setStatus(null);
                setCreatedIssueUrl(null);
            }, 2000);

        } catch (err) {
            console.error('Feedback submission failed:', err);
            setStatus({ 
                type: 'error', 
                message: err instanceof Error ? err.message : 'Failed to submit feedback' 
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleBackdropClick = (e: React.MouseEvent) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm"
            onClick={handleBackdropClick}
            data-testid="feedback-modal"
        >
            <div className="relative w-full max-w-xl mx-4 bg-gray-900 rounded-xl border border-gray-700 shadow-2xl">
                {/* Header */}
                <div className="flex items-center justify-between p-4 border-b border-gray-700">
                    <h2 className="text-lg font-semibold text-white flex items-center gap-2">
                        <Github size={20} />
                        Send Feedback
                    </h2>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-white transition-colors"
                        aria-label="Close"
                    >
                        <X size={20} />
                    </button>
                </div>

                {/* Content */}
                <div className="p-6">
                    {showSetup ? (
                        <div className="space-y-4">
                            <div className="p-4 bg-blue-900/20 border border-blue-500/30 rounded-lg">
                                <h3 className="font-medium text-blue-300 flex items-center gap-2 mb-2">
                                    <Settings size={16} />
                                    Setup Required
                                </h3>
                                <p className="text-sm text-gray-400 mb-3">
                                    To submit feedback directly to GitHub, you need a Personal Access Token (PAT).
                                </p>
                                <div className="text-xs text-gray-400 space-y-2">
                                    <p>A PAT allows this app to:</p>
                                    <ul className="list-disc list-inside ml-2 space-y-1">
                                        <li>Create issues in the forge-orchestrator repository</li>
                                        <li>Upload screenshots to help developers debug</li>
                                    </ul>
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-300 mb-2">
                                    GitHub Personal Access Token (Required)
                                </label>
                                <div className="p-3 bg-gray-800/50 border border-gray-700 rounded-lg mb-3">
                                    <p className="text-xs text-gray-400 mb-2">
                                        Click the button below to generate a token with the correct permissions:
                                    </p>
                                    <a 
                                        href="https://github.com/settings/tokens/new?scopes=public_repo&description=Forge+Orchestrator+Feedback" 
                                        target="_blank" 
                                        rel="noreferrer"
                                        className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-sm rounded-lg transition-colors"
                                        data-testid="generate-token-link"
                                    >
                                        <Github size={16} />
                                        Generate Token on GitHub
                                        <ExternalLink size={14} />
                                    </a>
                                    <p className="text-xs text-gray-500 mt-2">
                                        Scope needed: <code className="px-1 py-0.5 bg-gray-900 rounded text-blue-400">public_repo</code>
                                    </p>
                                </div>
                                <input
                                    type="password"
                                    value={githubToken}
                                    onChange={(e) => setGithubToken(e.target.value)}
                                    placeholder="ghp_xxxxxxxxxxxx"
                                    className="w-full px-3 py-2 bg-gray-950 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    data-testid="github-token-input"
                                />
                            </div>

                            {status && (
                                <div className={`p-3 rounded-lg text-sm ${
                                    status.type === 'error' ? 'bg-red-900/20 text-red-300' : 'bg-blue-900/20 text-blue-300'
                                }`}>
                                    {status.message}
                                </div>
                            )}

                            <div className="flex justify-end">
                                <button
                                    onClick={handleSaveSettings}
                                    className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors"
                                >
                                    Save Settings
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div className="space-y-4">
                            <div className="flex justify-between items-center">
                                <label className="block text-sm font-medium text-gray-300">
                                    Issue Description
                                </label>
                                <button
                                    onClick={() => setShowSetup(true)}
                                    className="text-xs text-gray-500 hover:text-gray-300"
                                >
                                    Update Settings
                                </button>
                            </div>
                            <textarea
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                placeholder="Describe the issue or suggestion..."
                                rows={4}
                                className="w-full px-3 py-2 bg-gray-950 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                            />

                            {/* Screenshots */}
                            <div>
                                <div className="flex justify-between items-center mb-2">
                                    <label className="text-sm font-medium text-gray-300">
                                        Screenshots {screenshots.length > 0 && `(${screenshots.length})`}
                                    </label>
                                    <button
                                        onClick={handleCapture}
                                        disabled={isCapturing}
                                        className="flex items-center gap-1 px-3 py-1.5 text-sm bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg transition-colors disabled:opacity-50"
                                    >
                                        <Camera size={14} />
                                        {isCapturing ? 'Capturing...' : 'Capture Screen'}
                                    </button>
                                </div>

                                {screenshots.length > 0 && (
                                    <div className="space-y-2">
                                        {screenshots.map((screenshot, index) => (
                                            <div key={index} className="relative p-2 bg-gray-800 rounded-lg">
                                                <img 
                                                    src={screenshot.dataUrl} 
                                                    alt={`Screenshot ${index + 1}`}
                                                    className="max-w-full rounded"
                                                />
                                                <button
                                                    onClick={() => handleRemoveScreenshot(index)}
                                                    className="absolute top-2 right-2 p-1 bg-red-600 hover:bg-red-500 rounded text-white text-xs flex items-center gap-1"
                                                >
                                                    <Trash2 size={12} />
                                                    Remove
                                                </button>
                                            </div>
                                        ))}
                                        <p className="text-xs text-gray-500 flex items-center gap-1">
                                            <Check size={12} className="text-green-400" />
                                            {screenshots.length} screenshot{screenshots.length > 1 ? 's' : ''} ready
                                        </p>
                                    </div>
                                )}
                            </div>

                            {/* Status */}
                            {status && (
                                <div className={`p-3 rounded-lg text-sm ${
                                    status.type === 'error' ? 'bg-red-900/20 text-red-300' :
                                    status.type === 'success' ? 'bg-green-900/20 text-green-300' :
                                    'bg-blue-900/20 text-blue-300'
                                }`}>
                                    {createdIssueUrl ? (
                                        <a 
                                            href={createdIssueUrl}
                                            target="_blank"
                                            rel="noreferrer"
                                            className="flex items-center gap-1 hover:underline"
                                        >
                                            {status.message} <ExternalLink size={12} />
                                        </a>
                                    ) : status.message}
                                </div>
                            )}

                            {/* Actions */}
                            <div className="flex justify-end gap-2 pt-2 border-t border-gray-700">
                                <button
                                    onClick={onClose}
                                    className="px-4 py-2 text-gray-400 hover:text-white transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleSubmit}
                                    disabled={!description.trim() || isSubmitting}
                                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 disabled:text-gray-500 text-white rounded-lg transition-colors"
                                >
                                    {isSubmitting && <Loader2 size={16} className="animate-spin" />}
                                    {isSubmitting ? 'Submitting...' : 'Submit Feedback'}
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export default FeedbackModal;
