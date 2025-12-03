import { type ReactNode } from 'react';
import { cn } from '../../lib/utils';

interface MainContentProps {
    children?: ReactNode;
    className?: string;
}

export function MainContent({ children, className }: MainContentProps) {
    return (
        <main
            data-testid="main-content"
            className={cn(
                "flex-1 relative overflow-hidden bg-gray-950 text-white",
                className
            )}
        >
            <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-blue-900/20 via-gray-900/0 to-gray-900/0 pointer-events-none" />
            <div className="relative h-full overflow-auto p-6 md:p-8">
                {children}
            </div>
        </main>
    );
}
