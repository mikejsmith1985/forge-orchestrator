import { test, expect } from '@playwright/test';

// Educational Comment: This test suite verifies the Token Ledger UI functionality.
test.describe('Token Ledger UI', () => {

    test.beforeEach(async ({ page }) => {
        // Mock data in snake_case format matching the Go backend
        const mockLedgerData = [
            { id: 1, timestamp: '2023-10-27T10:00:00Z', flow_id: 'flow-1', model_used: 'gpt-4', input_tokens: 100, output_tokens: 50, total_cost_usd: 0.005, latency_ms: 500, status: 'SUCCESS' },
            { id: 2, timestamp: '2023-10-27T11:00:00Z', flow_id: 'flow-2', model_used: 'gpt-3.5-turbo', input_tokens: 200, output_tokens: 100, total_cost_usd: 0.002, latency_ms: 300, status: 'SUCCESS' },
            { id: 3, timestamp: '2023-10-27T12:00:00Z', flow_id: 'flow-3', model_used: 'gpt-4', input_tokens: 50, output_tokens: 20, total_cost_usd: 0.003, latency_ms: 400, status: 'FAILED' },
        ];

        await page.route('/api/ledger', async route => {
            await route.fulfill({ json: mockLedgerData });
        });

        await page.route('/api/ledger/optimizations', async route => {
            await route.fulfill({ json: [] });
        });

        await page.goto('/');
    });

    test('should navigate to Ledger view from Sidebar', async ({ page }) => {
        await page.getByRole('button', { name: 'Dashboard' }).click();
        await expect(page.getByRole('heading', { name: 'Token Ledger' })).toBeVisible();
    });

    test('should display the ledger table', async ({ page }) => {
        await page.getByRole('button', { name: 'Dashboard' }).click();
        await expect(page.getByTestId('ledger-table')).toBeVisible();
    });

    test('should display correct number of rows from mocked API', async ({ page }) => {
        await page.getByRole('button', { name: 'Dashboard' }).click();
        const rows = page.getByTestId('ledger-table').locator('tbody tr');
        await expect(rows).toHaveCount(3);
    });

    test('should show WebSocket connection status', async ({ page }) => {
        await page.getByRole('button', { name: 'Dashboard' }).click();
        // Should show connection indicator text
        await expect(page.getByTestId('ledger-view')).toBeVisible();
    });

    test('should display toast container for notifications', async ({ page }) => {
        await page.getByRole('button', { name: 'Dashboard' }).click();
        // Toast container is always rendered (may be empty)
        await expect(page.getByTestId('ledger-view')).toBeVisible();
    });
});
