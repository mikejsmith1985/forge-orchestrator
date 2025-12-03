// Educational Comment: Import Playwright test utilities.
// We use 'test' to define our test cases and 'expect' for assertions.
import { test, expect, type Page } from '@playwright/test';

// Educational Comment: Group all Token Optimizer tests.
// This suite verifies the display and interaction of optimization suggestions
// in the Ledger view.
test.describe('Token Optimizer', () => {

    // Educational Comment: Helper to mock the ledger API which is required for the view to load.
    // We type the page argument to avoid implicit any errors.
    const mockLedger = async (page: Page) => {
        await page.route('/api/ledger', async route => {
            await route.fulfill({ json: [] });
        });
    };

    // Educational Comment: TEST 1 - Verify Optimizations Display
    // Purpose: Ensure that when the API returns suggestions, they are correctly
    // displayed as cards in the UI.
    test('should display optimization suggestions', async ({ page }) => {
        // Mock the ledger API as it's required for the component to load without error.
        await mockLedger(page);

        // Educational Comment: Mock the API response for optimizations.
        // We must match the Suggestion interface expected by OptimizationCard.
        await page.route('/api/ledger/optimizations', async route => {
            const json = [
                {
                    id: 'opt-1',
                    title: 'Unused Token Cleanup',
                    description: 'Remove 50 unused tokens to save space.',
                    estimatedSavings: 50.0,
                    status: 'pending'
                },
                {
                    id: 'opt-2',
                    title: 'Compress History',
                    description: 'Compress old ledger entries.',
                    estimatedSavings: 120.5,
                    status: 'pending'
                }
            ];
            await route.fulfill({ json });
        });

        // Educational Comment: Navigate to the root URL first.
        await page.goto('/');

        // Educational Comment: Click "Dashboard" in the sidebar to switch to the Ledger view.
        await page.click('text=Dashboard');

        // Educational Comment: Verify that the optimization cards are visible.
        await expect(page.getByText('Unused Token Cleanup')).toBeVisible();
        await expect(page.getByText('Estimated Savings: $50.0000')).toBeVisible();
        await expect(page.getByText('Compress History')).toBeVisible();
        await expect(page.getByText('Estimated Savings: $120.5000')).toBeVisible();
    });

    // Educational Comment: TEST 2 - Verify "Click to Apply" Workflow
    // Purpose: Verify that clicking the "Apply" button updates the UI state
    // to indicate the optimization has been applied.
    test('should apply optimization on click', async ({ page }) => {
        // Mock the ledger API.
        await mockLedger(page);

        // Educational Comment: Mock the initial suggestions.
        await page.route('/api/ledger/optimizations', async route => {
            const json = [
                {
                    id: 'opt-1',
                    title: 'Quick Fix',
                    description: 'A quick optimization.',
                    estimatedSavings: 10.0,
                    status: 'pending'
                }
            ];
            await route.fulfill({ json });
        });

        // Educational Comment: Mock the apply action endpoint.
        await page.route('/api/ledger/optimizations/opt-1/apply', async route => {
            await route.fulfill({ status: 200, json: { success: true } });
        });

        await page.goto('/');
        await page.click('text=Dashboard');

        // Educational Comment: Find the "Apply" button for the specific card.
        const applyButton = page.getByRole('button', { name: 'Apply' }).first();
        await expect(applyButton).toBeVisible();

        // Educational Comment: Click the button.
        await applyButton.click();

        // Educational Comment: Verify the button state changes to "Applied".
        await expect(page.getByRole('button', { name: 'Applied' })).toBeVisible();
        // Ensure the button is now disabled.
        await expect(page.getByRole('button', { name: 'Applied' })).toBeDisabled();
    });

    // Educational Comment: TEST 3 - Verify Empty State
    // Purpose: Ensure that when there are no optimizations, the UI shows
    // an appropriate empty state message instead of a blank space or error.
    test('should display empty state when no suggestions', async ({ page }) => {
        // Mock the ledger API.
        await mockLedger(page);

        // Educational Comment: Mock an empty list of suggestions.
        await page.route('/api/ledger/optimizations', async route => {
            await route.fulfill({ json: [] });
        });

        await page.goto('/');
        await page.click('text=Dashboard');

        // Educational Comment: Verify the empty state message is visible.
        await expect(page.getByText('No optimization suggestions yet')).toBeVisible();
    });
});
