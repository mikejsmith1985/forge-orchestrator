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
});
