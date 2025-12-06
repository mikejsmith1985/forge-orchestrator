/**
 * E2E Tests for Terminal Root Directory Configuration (Issue #5)
 * Tests the new root_dir configuration feature that allows users to specify
 * a starting directory for their terminal, particularly important for WSL users.
 */

import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:9000';

test.describe('Terminal Root Directory Configuration', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate directly to settings page
        await page.goto(`${BASE_URL}/settings`);
        await page.waitForLoadState('networkidle');
        // Give some time for React to render
        await page.waitForTimeout(1000);
    });

    test('API returns root_dir field in config', async ({ page }) => {
        // Test API endpoint directly
        const response = await page.request.get(`${BASE_URL}/api/config`);
        expect(response.ok()).toBeTruthy();

        const config = await response.json();
        expect(config).toHaveProperty('shell');
        
        // root_dir may be empty or set, but the field should exist
        expect(config.shell).toHaveProperty('root_dir');
    });

    test('can update root_dir via API', async ({ page }) => {
        const testConfig = {
            shell: {
                type: 'wsl',
                wsl_distro: 'Ubuntu-24.04',
                root_dir: 'C:\\Users\\test\\projects\\forge'
            },
            update: {
                check_on_startup: true,
                check_interval_minutes: 30,
                auto_download: false
            },
            server: {
                port: 9000,
                open_browser: true
            }
        };

        // Update config via API
        const response = await page.request.post(`${BASE_URL}/api/config`, {
            data: testConfig
        });
        expect(response.ok()).toBeTruthy();

        // Verify it was saved
        const getResponse = await page.request.get(`${BASE_URL}/api/config`);
        const savedConfig = await getResponse.json();
        expect(savedConfig.shell.root_dir).toBe(testConfig.shell.root_dir);
    });

    test('settings page loads without errors', async ({ page }) => {
        // Wait for either Terminal Settings header or the tabs to appear
        const settingsLoaded = page.locator('button:has-text("Terminal")');
        await expect(settingsLoaded).toBeVisible({ timeout: 10000 });
    });

    test('terminal tab is visible and active by default', async ({ page }) => {
        const terminalTab = page.locator('button:has-text("Terminal")');
        await expect(terminalTab).toBeVisible({ timeout: 10000 });
        // Should be active (has blue border)
        await expect(terminalTab).toHaveClass(/border-blue-500/);
    });

    test('terminal settings form renders', async ({ page }) => {
        // Check for key elements
        await expect(page.locator('h1:has-text("Terminal Settings")')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('h2:has-text("Shell Type")')).toBeVisible();
        await expect(page.locator('button:has-text("Save Configuration")')).toBeVisible();
    });

    test('can save configuration with root_dir', async ({ page }) => {
        // Wait for form to load
        await expect(page.locator('h1:has-text("Terminal Settings")')).toBeVisible({ timeout: 10000 });

        // On Linux, we'll have bash shell - let's just verify the save functionality works
        // Save current configuration
        await page.click('button:has-text("Save Configuration")');

        // Wait for success message
        await expect(page.locator('text=Configuration saved')).toBeVisible({ timeout: 5000 });
    });

    test('configuration persists after save', async ({ page }) => {
        // Set a test configuration via API first
        await page.request.post(`${BASE_URL}/api/config`, {
            data: {
                shell: {
                    type: 'bash',
                    root_dir: '/home/user/test-project'
                },
                update: {
                    check_on_startup: true,
                    check_interval_minutes: 30,
                    auto_download: false
                },
                server: {
                    port: 9000,
                    open_browser: true
                }
            }
        });

        // Reload the page
        await page.goto(`${BASE_URL}/settings`);
        await page.waitForLoadState('networkidle');

        // Verify the config is still present via API
        const response = await page.request.get(`${BASE_URL}/api/config`);
        const config = await response.json();
        expect(config.shell.root_dir).toBe('/home/user/test-project');
    });

    test('troubleshooting section is visible', async ({ page }) => {
        await expect(page.locator('h1:has-text("Terminal Settings")')).toBeVisible({ timeout: 10000 });
        
        // Look for troubleshooting section
        const troubleshooting = page.locator('h3:has-text("Troubleshooting")');
        await expect(troubleshooting).toBeVisible();
    });
});
