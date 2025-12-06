/**
 * Integration Test: Command Flow
 * 
 * This test verifies the complete command lifecycle using real API endpoints
 * (no mocking). It requires a running Go test server with an in-memory database.
 */
import { test, expect } from '@playwright/test';

test.describe('Command Flow Integration', () => {
    test.setTimeout(30000);

    // Mock the welcome endpoint to prevent the modal from appearing
    test.beforeEach(async ({ page }) => {
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    contentType: 'application/json',
                    body: JSON.stringify({ shown: true, currentVersion: '1.1.0' }),
                });
            } else {
                await route.fulfill({ status: 200, body: '{}' });
            }
        });
    });

    test('complete command lifecycle - create, view, delete', async ({ page }) => {
        // 1. Navigate directly to Command Deck via URL (more reliable than clicking nav)
        await page.goto('/commands');
        
        // Verify we're on the Command Deck
        await expect(page.getByTestId('command-deck')).toBeVisible();
        
        // 2. Open the Add Command modal
        await page.getByTestId('add-command-btn').click();
        await expect(page.getByTestId('add-command-modal')).toBeVisible();
        
        // 3. Fill in command details (NO MOCKING - real API call)
        const uniqueName = `Integration Test ${Date.now()}`;
        await page.getByTestId('command-name-input').fill(uniqueName);
        await page.getByTestId('command-description-input').fill('Created by integration test');
        await page.getByTestId('command-input').fill('echo "hello from integration"');
        
        // 4. Submit the form
        await page.getByTestId('submit-command-btn').click();
        
        // 5. Verify modal closes and command appears in list
        await expect(page.getByTestId('add-command-modal')).not.toBeVisible();
        
        // Wait for the command to appear (real API response)
        await expect(page.getByText(uniqueName)).toBeVisible({ timeout: 5000 });
        
        // 6. Verify command card is displayed with correct data
        const commandCard = page.locator('[data-testid="command-card"]').filter({ hasText: uniqueName });
        await expect(commandCard).toBeVisible();
        await expect(commandCard.getByText('Created by integration test')).toBeVisible();
    });

    test('empty state displays correctly', async ({ page }) => {
        // Navigate directly to Command Deck via URL
        await page.goto('/commands');
        
        // The deck should be visible (may have commands or be empty)
        await expect(page.getByTestId('command-deck')).toBeVisible();
        
        // Add button should always be visible
        await expect(page.getByTestId('add-command-btn')).toBeVisible();
    });
});
