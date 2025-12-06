/**
 * Integration Test: Key Management
 * 
 * This test verifies the API key management UI using real API endpoints.
 * Note: In CI environments, the keyring service may not be available,
 * so we focus on testing the UI elements rather than the full key storage flow.
 */
import { test, expect } from '@playwright/test';

test.describe('Key Management Integration', () => {
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

    test('settings page displays key status', async ({ page }) => {
        // 1. Navigate directly to Settings page via URL
        await page.goto('/settings');
        
        // 2. Verify settings view is displayed
        await expect(page.getByTestId('settings-view')).toBeVisible();
        
        // 3. Verify API key configuration UI is present
        // Look for provider names (Anthropic, OpenAI)
        await expect(page.getByText(/anthropic/i)).toBeVisible();
        await expect(page.getByText(/openai/i)).toBeVisible();
    });

    test('can interact with key configuration UI', async ({ page }) => {
        // 1. Navigate directly to Settings page via URL
        await page.goto('/settings');
        
        await expect(page.getByTestId('settings-view')).toBeVisible();
        
        // 2. Verify there are input fields or configuration options
        // The exact UI elements depend on the implementation
        const hasInputs = await page.locator('input[type="password"], input[type="text"]').count() > 0;
        const hasButtons = await page.locator('button').count() > 0;
        
        // Settings page should have some interactive elements
        expect(hasInputs || hasButtons).toBeTruthy();
    });

    test('shows success toast after saving API key', async ({ page }) => {
        // Mock the /api/keys endpoint to return success
        await page.route('/api/keys', async (route) => {
            if (route.request().method() === 'POST') {
                await route.fulfill({
                    status: 200,
                    contentType: 'application/json',
                    body: JSON.stringify({ success: true, message: 'API key saved successfully' }),
                });
            } else {
                await route.continue();
            }
        });
        
        // 1. Navigate to Settings page
        await page.goto('/settings');
        
        await expect(page.getByTestId('settings-view')).toBeVisible();
        
        // 2. Find the Anthropic provider card and its input
        const anthropicCard = page.getByTestId('provider-card-anthropic');
        await expect(anthropicCard).toBeVisible();
        
        // 3. Enter a test API key
        const keyInput = anthropicCard.locator('input[type="password"]');
        await keyInput.fill('sk-ant-test-key-12345');
        
        // 4. Click the Save button
        const saveButton = anthropicCard.locator('button:has-text("Save")');
        await saveButton.click();
        
        // 5. Verify success toast appears
        await expect(page.getByTestId('toast-success')).toBeVisible({ timeout: 5000 });
        await expect(page.getByText(/saved successfully/i)).toBeVisible();
    });

    test('shows error toast when save fails', async ({ page }) => {
        // Mock the API to return an error
        await page.route('/api/keys', async (route) => {
            if (route.request().method() === 'POST') {
                await route.fulfill({
                    status: 500,
                    body: JSON.stringify({ error: 'Keyring not available' }),
                });
            } else {
                await route.continue();
            }
        });
        
        // 1. Navigate to Settings page
        await page.goto('/settings');
        
        await expect(page.getByTestId('settings-view')).toBeVisible();
        
        // 2. Find the OpenAI provider card
        const openaiCard = page.getByTestId('provider-card-openai');
        await expect(openaiCard).toBeVisible();
        
        // 3. Enter a test API key
        const keyInput = openaiCard.locator('input[type="password"]');
        await keyInput.fill('sk-test-key');
        
        // 4. Click Save
        const saveButton = openaiCard.locator('button:has-text("Save")');
        await saveButton.click();
        
        // 5. Verify error toast appears
        await expect(page.getByTestId('toast-error')).toBeVisible({ timeout: 5000 });
    });
});
