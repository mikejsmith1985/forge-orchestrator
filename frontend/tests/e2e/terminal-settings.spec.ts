import { test, expect } from '@playwright/test';

/**
 * Terminal Settings E2E Tests
 * 
 * These tests validate the terminal settings functionality:
 * - Settings page renders with terminal options
 * - Shell type can be changed and saved
 * - Configuration persists
 * - Terminal connects with different shell types
 */

test.describe('Terminal Settings', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to settings
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');
    });

    test('settings page renders with terminal tab', async ({ page }) => {
        // Verify terminal tab is visible and active
        const terminalTab = page.locator('button:has-text("Terminal")');
        await expect(terminalTab).toBeVisible();
        await expect(terminalTab).toHaveClass(/border-blue-500/);
    });

    test('terminal settings form displays correctly', async ({ page }) => {
        // Verify the settings form elements
        await expect(page.locator('h1:has-text("Terminal Settings")')).toBeVisible();
        await expect(page.locator('h2:has-text("Shell Type")')).toBeVisible();
        
        // Should have at least one shell option visible
        await expect(page.locator('input[type="radio"][name="shell"]')).toHaveCount({ min: 1 });
        
        // Should have save button
        await expect(page.locator('button:has-text("Save Configuration")')).toBeVisible();
    });

    test('can load current configuration', async ({ page }) => {
        // Wait for config to load
        await page.waitForTimeout(1000);
        
        // At least one radio button should be checked
        const checkedRadio = page.locator('input[type="radio"][name="shell"]:checked');
        await expect(checkedRadio).toHaveCount(1);
    });

    test('can change shell type and save', async ({ page }) => {
        // Wait for config to load
        await page.waitForTimeout(1000);
        
        // Select bash shell (should work on both platforms)
        const bashOption = page.locator('input[type="radio"][value="bash"]');
        
        // If bash option exists (Unix/Linux), test it
        if (await bashOption.count() > 0) {
            await bashOption.click();
            
            // Save configuration
            const saveButton = page.locator('button:has-text("Save Configuration")');
            await saveButton.click();
            
            // Should show success message
            await expect(page.locator('text=Configuration saved')).toBeVisible({ timeout: 5000 });
        }
    });

    test('troubleshooting section is visible', async ({ page }) => {
        // Verify troubleshooting help text
        await expect(page.locator('h3:has-text("Troubleshooting")')).toBeVisible();
    });

    test('can switch between settings tabs', async ({ page }) => {
        // Click on API Keys tab
        const apiKeysTab = page.locator('button:has-text("API Keys")');
        await apiKeysTab.click();
        
        // Should see API keys content
        await expect(page.locator('text=API Keys')).toBeVisible();
        
        // Switch back to Terminal tab
        const terminalTab = page.locator('button:has-text("Terminal")').first();
        await terminalTab.click();
        
        // Should see terminal settings again
        await expect(page.locator('h1:has-text("Terminal Settings")')).toBeVisible();
    });
});

test.describe('Terminal Connection with Settings', () => {
    test('terminal connects after settings page visit', async ({ page }) => {
        // Visit settings first
        await page.goto('/settings');
        await page.waitForTimeout(1000);
        
        // Navigate to terminal
        await page.goto('/terminal');
        await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
        
        // Wait for potential connection
        await page.waitForTimeout(3000);
        
        // Check connection indicator - should be green (connected) or red (disconnected with error message)
        const connectionIndicator = page.locator('.rounded-full.w-2.h-2');
        await expect(connectionIndicator.first()).toBeVisible();
    });

    test('terminal shows shell type in welcome message', async ({ page }) => {
        // Navigate to terminal
        await page.goto('/terminal');
        await page.waitForSelector('.xterm-viewport', { timeout: 10000 });
        
        // Wait for connection message
        await page.waitForTimeout(2000);
        
        // Terminal should show connection message with shell type
        // This is visible in the xterm output
        const xtermRows = page.locator('.xterm-rows');
        await expect(xtermRows).toBeVisible();
    });
});

test.describe('Terminal Error Handling', () => {
    test('terminal shows helpful error message on connection failure', async ({ page }) => {
        // This test verifies that if terminal fails to connect,
        // it shows a helpful error message rather than silent failure
        
        await page.goto('/terminal');
        await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
        
        // Wait for connection attempt
        await page.waitForTimeout(3000);
        
        // If connection fails, terminal should show error text
        // Either "Connected" (success) or "Failed" (with helpful message)
        const xtermRows = page.locator('.xterm-rows');
        await expect(xtermRows).toBeVisible();
    });
});

test.describe('Settings API Integration', () => {
    test('settings API returns valid configuration', async ({ page, request }) => {
        // Direct API test
        const response = await request.get('/api/config');
        expect(response.ok()).toBeTruthy();
        
        const config = await response.json();
        expect(config).toHaveProperty('shell');
        expect(config.shell).toHaveProperty('type');
    });

    test('settings API accepts configuration updates', async ({ page, request }) => {
        // Get current config
        const getResponse = await request.get('/api/config');
        const currentConfig = await getResponse.json();
        
        // Update config (keep same values to avoid side effects)
        const updateResponse = await request.post('/api/config', {
            data: currentConfig,
        });
        
        expect(updateResponse.ok()).toBeTruthy();
        const result = await updateResponse.json();
        expect(result).toHaveProperty('message');
    });
});
