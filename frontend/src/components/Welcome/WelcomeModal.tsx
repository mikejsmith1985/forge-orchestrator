import { useEffect } from 'react';
import { X, Sparkles, BookOpen, BarChart3, GitBranch, Key, Zap } from 'lucide-react';

interface WelcomeModalProps {
    isOpen: boolean;
    onClose: () => void;
    version: string;
}

interface Feature {
    icon: React.ReactNode;
    title: string;
    description: string;
}

const features: Feature[] = [
    {
        icon: <Sparkles className="w-6 h-6" />,
        title: 'Architect',
        description: 'Brain dump your ideas and get token estimates in real-time.'
    },
    {
        icon: <BookOpen className="w-6 h-6" />,
        title: 'Ledger',
        description: 'Track and analyze your token usage across all interactions.'
    },
    {
        icon: <GitBranch className="w-6 h-6" />,
        title: 'Flows',
        description: 'Create visual workflows connecting AI agents together.'
    },
    {
        icon: <Zap className="w-6 h-6" />,
        title: 'Commands',
        description: 'Save and execute frequently used commands with one click.'
    },
    {
        icon: <Key className="w-6 h-6" />,
        title: 'API Keys',
        description: 'Securely manage your AI provider credentials.'
    },
    {
        icon: <BarChart3 className="w-6 h-6" />,
        title: 'Optimization',
        description: 'Get suggestions to reduce token usage and costs.'
    },
];

export function WelcomeModal({ isOpen, onClose, version }: WelcomeModalProps) {
    useEffect(() => {
        if (!isOpen) return;

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape' || e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onClose();
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    const handleBackdropClick = (e: React.MouseEvent) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm"
            onClick={handleBackdropClick}
            data-testid="welcome-modal"
        >
            <div className="relative w-full max-w-2xl mx-4 bg-gray-900 rounded-2xl border border-gray-700 shadow-2xl overflow-hidden">
                {/* Close button */}
                <button
                    onClick={onClose}
                    className="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors z-10"
                    aria-label="Close"
                >
                    <X size={24} />
                </button>

                {/* Header */}
                <div className="pt-8 pb-6 px-8 text-center bg-gradient-to-b from-blue-900/30 to-transparent">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-600/20 mb-4">
                        <Sparkles className="w-8 h-8 text-blue-400" />
                    </div>
                    <h1 className="text-3xl font-bold text-white mb-2">
                        Forge Orchestrator
                    </h1>
                    <p className="text-gray-400 text-sm">
                        Version {version}
                    </p>
                </div>

                {/* Features Grid */}
                <div className="px-8 pb-6">
                    <p className="text-gray-300 text-center mb-6">
                        Your AI workflow command center. Here's what you can do:
                    </p>
                    
                    <div className="grid grid-cols-2 gap-4">
                        {features.map((feature) => (
                            <div
                                key={feature.title}
                                className="flex items-start gap-3 p-3 rounded-lg bg-gray-800/50 border border-gray-700/50"
                            >
                                <div className="flex-shrink-0 text-blue-400">
                                    {feature.icon}
                                </div>
                                <div>
                                    <h3 className="font-semibold text-white text-sm">
                                        {feature.title}
                                    </h3>
                                    <p className="text-gray-400 text-xs mt-0.5">
                                        {feature.description}
                                    </p>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Footer */}
                <div className="px-8 pb-8 text-center">
                    <button
                        onClick={onClose}
                        className="px-8 py-3 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg transition-colors shadow-lg shadow-blue-500/20"
                    >
                        Get Started
                    </button>
                    <p className="text-gray-500 text-xs mt-4">
                        Press any key to continue
                    </p>
                </div>
            </div>
        </div>
    );
}

export default WelcomeModal;
