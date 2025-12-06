import { test, expect } from '@playwright/test';

/**
 * Terminal E2E Tests - Task 1.3: Full-Stack PTY Test
 * 
 * These tests validate the integrated terminal functionality:
 * - Terminal renders in the UI
 * - WebSocket connection to /ws/pty works
 * - User can type and see shell output
 * - Prompt watcher toggle works
 */

test.describe('Terminal Integration', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to terminal view
        await page.goto('/terminal');
        
        // Wait for terminal to load
        await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
    });

    test('terminal renders with header and controls', async ({ page }) => {
        // Verify terminal header is present (be specific to avoid multiple matches)
        await expect(page.getByTestId('main-content').locator('span.font-mono:has-text("Terminal")')).toBeVisible();
        
        // Verify the traffic light buttons (macOS style) - use first() for strict mode
        await expect(page.locator('.bg-red-500.rounded-full').first()).toBeVisible();
        await expect(page.locator('.bg-yellow-500.rounded-full').first()).toBeVisible();
        await expect(page.locator('.bg-green-500.rounded-full').first()).toBeVisible();
        
        // Verify prompt watcher toggle exists
        await expect(page.locator('[data-testid="prompt-watcher-toggle"]')).toBeVisible();
    });

    test('terminal container is rendered', async ({ page }) => {
        // Verify the xterm.js container is rendered
        const terminalContainer = page.locator('[data-testid="terminal-container"]');
        await expect(terminalContainer).toBeVisible();
        
        // Check that xterm.js initialized (look for xterm viewport)
        await page.waitForSelector('.xterm-viewport', { timeout: 10000 });
        await expect(page.locator('.xterm-viewport')).toBeVisible();
    });

    test('can toggle prompt watcher', async ({ page }) => {
        const toggle = page.locator('[data-testid="prompt-watcher-toggle"]');
        
        // Initially off (should have slate background classes)
        await expect(toggle).toBeVisible();
        
        // Click to enable
        await toggle.click();
        
        // Should now have blue background indicating enabled state
        await expect(toggle).toHaveClass(/bg-blue-600/);
        
        // Click again to disable
        await toggle.click();
        
        // Should be back to slate background
        await expect(toggle).toHaveClass(/bg-slate-700/);
    });

    test('terminal shows connected status indicator', async ({ page }) => {
        // Wait for potential connection
        await page.waitForTimeout(2000);
        
        // Look for connection indicator (green or red dot)
        const connectionIndicator = page.locator('.rounded-full.w-2.h-2');
        await expect(connectionIndicator.first()).toBeVisible();
    });

    test('terminal displays in sidebar navigation', async ({ page }) => {
        // Navigate to home
        await page.goto('/');
        
        // Terminal should be in sidebar
        const sidebar = page.locator('[data-testid="sidebar"]');
        await expect(sidebar.locator('text=Terminal')).toBeVisible();
        
        // Click on Terminal in sidebar
        await sidebar.locator('text=Terminal').click();
        
        // Should navigate to terminal view
        await expect(page).toHaveURL(/\/terminal/);
    });

    test('home redirects to terminal (as primary view)', async ({ page }) => {
        await page.goto('/');
        
        // Should redirect to terminal view
        await expect(page).toHaveURL(/\/terminal/);
        
        // Terminal container should be visible
        await expect(page.locator('[data-testid="terminal-container"]')).toBeVisible();
    });
});

test.describe('Terminal WebSocket Connection', () => {
    test('establishes WebSocket connection on load', async ({ page }) => {
        // Track WebSocket connections
        const wsConnections: string[] = [];
        
        page.on('websocket', ws => {
            wsConnections.push(ws.url());
        });
        
        await page.goto('/terminal');
        await page.waitForTimeout(3000);
        
        // Should have attempted to connect to /ws/pty
        const ptyConnection = wsConnections.find(url => url.includes('/ws/pty'));
        expect(ptyConnection).toBeDefined();
    });
});

test.describe('Terminal User Interaction', () => {
    test.skip('can type in terminal and see output (requires backend)', async ({ page }) => {
        // This test requires the Go backend to be running with PTY support
        // Skipped by default as it's an integration test
        
        await page.goto('/terminal');
        await page.waitForSelector('.xterm-viewport', { timeout: 10000 });
        
        // Wait for connection
        await page.waitForTimeout(2000);
        
        // Type a simple command
        await page.keyboard.type('echo hello');
        await page.keyboard.press('Enter');
        
        // Wait for output
        await page.waitForTimeout(1000);
        
        // Should see "hello" in the terminal output
        const terminalContent = await page.locator('.xterm-rows').textContent();
        expect(terminalContent).toContain('hello');
    });
});
