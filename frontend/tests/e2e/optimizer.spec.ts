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
        // Educational Comment: Mock the API responses BEFORE navigation.
        // The route() method intercepts network requests matching the pattern.
        await page.route(/\/api\/ledger$/, async route => {
            await route.fulfill({ 
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([])
            });
        });
        
        await page.route(/\/api\/ledger\/optimizations$/, async route => {
            await route.fulfill({ 
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([
                    {
                        id: 1,
                        type: 'model_switch',
                        title: 'Unused Token Cleanup',
                        description: 'Remove 50 unused tokens to save space.',
                        estimated_savings: 50.0,
                        savings_unit: 'USD',
                        target_flow_id: 'flow-1',
                        apply_action: '{"action":"switch_model"}',
                        status: 'pending'
                    },
                    {
                        id: 2,
                        type: 'prompt_optimization',
                        title: 'Compress History',
                        description: 'Compress old ledger entries.',
                        estimated_savings: 120.5,
                        savings_unit: 'USD',
                        target_flow_id: 'flow-2',
                        apply_action: '{"action":"optimize_prompt"}',
                        status: 'pending'
                    }
                ])
            });
        });

        // Educational Comment: Navigate to the root URL first.
        await page.goto('/');

        // Educational Comment: Click "Dashboard" in the sidebar to switch to the Ledger view.
        await page.click('text=Dashboard');

        // Wait for data loading to complete
        await page.waitForLoadState('networkidle');

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
        // Mock the ledger API BEFORE navigation using regex pattern.
        await page.route(/\/api\/ledger$/, async route => {
            await route.fulfill({ 
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([])
            });
        });

        // Educational Comment: Mock the initial suggestions with correct data shape.
        await page.route(/\/api\/ledger\/optimizations$/, async route => {
            await route.fulfill({ 
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([
                    {
                        id: 1,
                        type: 'model_switch',
                        title: 'Quick Fix',
                        description: 'A quick optimization.',
                        estimated_savings: 10.0,
                        savings_unit: 'USD',
                        target_flow_id: 'flow-1',
                        apply_action: '{"action":"switch_model"}',
                        status: 'pending'
                    }
                ])
            });
        });

        // Educational Comment: Mock the apply action endpoint.
        await page.route(/\/api\/ledger\/optimizations\/1\/apply/, async route => {
            await route.fulfill({ 
                status: 200, 
                contentType: 'application/json',
                body: JSON.stringify({ success: true }) 
            });
        });

        await page.goto('/');
        await page.click('text=Dashboard');
        
        // Wait for data loading to complete
        await page.waitForLoadState('networkidle');

        // Educational Comment: Find the "Apply" button for the specific card.
        const applyButton = page.getByRole('button', { name: 'Apply' }).first();
        await expect(applyButton).toBeVisible();

        // Educational Comment: Click the button to open confirmation modal.
        await applyButton.click();
        
        // Click "Apply Changes" in the confirmation modal
        await page.getByRole('button', { name: 'Apply Changes' }).click();

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
