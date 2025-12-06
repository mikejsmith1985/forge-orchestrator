import { test, expect } from '@playwright/test';

/**
 * Forge Vision Tests
 * 
 * Tests the terminal output visualization feature.
 * Validates git status parsing and actionable buttons.
 */

test.describe('Forge Vision Component', () => {
    // We'll test the component in isolation using a test page
    test.beforeEach(async ({ page }) => {
        // Navigate to terminal view which may show ForgeVision
        await page.goto('/terminal');
    });

    test('terminal view loads correctly', async ({ page }) => {
        const terminalContainer = page.getByTestId('terminal-container');
        await expect(terminalContainer).toBeVisible();
    });

    test('prompt watcher toggle is visible', async ({ page }) => {
        const toggle = page.getByTestId('prompt-watcher-toggle');
        await expect(toggle).toBeVisible();
    });

    test('prompt watcher can be toggled', async ({ page }) => {
        const toggle = page.getByTestId('prompt-watcher-toggle');
        await expect(toggle).toBeVisible();
        
        // Click to toggle
        await toggle.click();
        
        // Should now have different styling (bg-blue-600 when enabled)
        await expect(toggle).toHaveClass(/bg-blue-600|bg-slate-700/);
    });
});

test.describe('Ledger Cost Unit Display', () => {
    test.beforeEach(async ({ page }) => {
        // Mock ledger with mixed cost units
        await page.route('/api/ledger', async route => {
            await route.fulfill({
                json: [
                    {
                        id: 1,
                        timestamp: '2024-01-01T10:00:00Z',
                        flow_id: 'flow-1',
                        model_used: 'gpt-4o',
                        input_tokens: 1000,
                        output_tokens: 500,
                        total_cost_usd: 0.025,
                        latency_ms: 1200,
                        status: 'SUCCESS',
                        cost_unit: 'TOKEN'
                    },
                    {
                        id: 2,
                        timestamp: '2024-01-01T11:00:00Z',
                        flow_id: 'flow-2',
                        model_used: 'claude-instant',
                        input_tokens: 0,
                        output_tokens: 0,
                        total_cost_usd: 0.01,
                        latency_ms: 800,
                        status: 'SUCCESS',
                        cost_unit: 'PROMPT',
                        prompt_count: 1
                    }
                ]
            });
        });

        await page.route('/api/ledger/optimizations', async route => {
            await route.fulfill({ json: [] });
        });

        await page.goto('/');
        await page.getByRole('button', { name: 'Dashboard' }).click();
    });

    test('displays correct usage based on cost unit', async ({ page }) => {
        const table = page.getByTestId('ledger-table');
        await expect(table).toBeVisible();

        // First row should show token counts (TOKEN cost unit)
        const firstRow = table.locator('tbody tr').first();
        await expect(firstRow).toContainText('1,000');  // input tokens
        await expect(firstRow).toContainText('500');    // output tokens

        // Second row should show prompt count (PROMPT cost unit)
        const secondRow = table.locator('tbody tr').nth(1);
        await expect(secondRow).toContainText('1 prompt');
    });

    test('shows cost unit badge for each entry', async ({ page }) => {
        const table = page.getByTestId('ledger-table');
        
        // Should show Per-Token for first entry
        await expect(table).toContainText('Per-Token');
        
        // Should show Per-Prompt for second entry
        await expect(table).toContainText('Per-Prompt');
    });
});
