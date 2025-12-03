import { test, expect } from '@playwright/test';

// Educational Comment: This test suite verifies the Token Ledger UI functionality.
// We are testing:
// 1. Navigation to the Ledger view.
// 2. Visibility of the Ledger table component.
// 3. Correct rendering of data by mocking the API response.
test.describe('Token Ledger UI', () => {

    // Educational Comment: Before each test, we navigate to the home page to ensure a clean state.
    // We also mock the /api/ledger endpoint to ensure the UI renders correctly without needing a real backend.
    test.beforeEach(async ({ page }) => {
        const mockData = [
            { id: '1', timestamp: '2023-10-27T10:00:00Z', flowId: 'flow-1', model: 'gpt-4', inputTokens: 100, outputTokens: 50, cost: 0.005, status: 'success' },
            { id: '2', timestamp: '2023-10-27T11:00:00Z', flowId: 'flow-2', model: 'gpt-3.5-turbo', inputTokens: 200, outputTokens: 100, cost: 0.002, status: 'pending' },
            { id: '3', timestamp: '2023-10-27T12:00:00Z', flowId: 'flow-3', model: 'gpt-4', inputTokens: 50, outputTokens: 20, cost: 0.003, status: 'failed' },
        ];

        await page.route('/api/ledger', async route => {
            await route.fulfill({ json: mockData });
        });

        await page.goto('/');
    });

    test('should navigate to Ledger view from Sidebar', async ({ page }) => {
        // Educational Comment: We simulate a user clicking the "Dashboard" button in the sidebar.
        // Since the app uses state-based routing, we verify the view change by checking for the heading.
        await page.getByRole('button', { name: 'Dashboard' }).click();
        await expect(page.getByRole('heading', { name: 'Token Ledger' })).toBeVisible();
    });

    test('should display the ledger table', async ({ page }) => {
        // Educational Comment: We navigate to the ledger view via the sidebar and check if the table exists.
        await page.getByRole('button', { name: 'Dashboard' }).click();
        await expect(page.getByTestId('ledger-table')).toBeVisible();
    });

    test('should display correct number of rows from mocked API', async ({ page }) => {
        // Educational Comment: Trigger navigation to load the data
        await page.getByRole('button', { name: 'Dashboard' }).click();

        // Educational Comment: We verify that the table body contains the expected number of rows (3).
        // This confirms that the component correctly renders the data it receives.
        const rows = page.getByTestId('ledger-table').locator('tbody tr');
        await expect(rows).toHaveCount(3);
    });
});
