/**
 * E2E Test: Welcome Modal
 * 
 * Tests the welcome/splash screen functionality for new users and version upgrades.
 */
import { test, expect } from '@playwright/test';

test.describe('Welcome Modal', () => {
    test.setTimeout(30000);

    test('shows welcome modal on first visit', async ({ page }) => {
        // Mock the welcome API to return shown=false (first time user)
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    body: JSON.stringify({
                        shown: false,
                        currentVersion: '1.1.0',
                        lastVersion: ''
                    }),
                });
            } else {
                await route.fulfill({ status: 200, body: JSON.stringify({ status: 'ok' }) });
            }
        });

        await page.goto('/');
        
        // Welcome modal should be visible
        await expect(page.getByTestId('welcome-modal')).toBeVisible({ timeout: 5000 });
        
        // Should show the application name
        await expect(page.getByText(/forge orchestrator/i)).toBeVisible();
        
        // Should show version
        await expect(page.getByText(/1\.1\.0/)).toBeVisible();
    });

    test('welcome modal shows feature overview', async ({ page }) => {
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    body: JSON.stringify({ shown: false, currentVersion: '1.1.0', lastVersion: '' }),
                });
            } else {
                await route.fulfill({ status: 200, body: JSON.stringify({ status: 'ok' }) });
            }
        });

        await page.goto('/');
        await expect(page.getByTestId('welcome-modal')).toBeVisible();

        // Should list main features
        await expect(page.getByText(/architect/i)).toBeVisible();
        await expect(page.getByText(/ledger/i)).toBeVisible();
        await expect(page.getByText(/flows/i)).toBeVisible();
        await expect(page.getByText(/settings|api keys/i)).toBeVisible();
    });

    test('welcome modal can be dismissed', async ({ page }) => {
        let welcomePosted = false;
        
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    body: JSON.stringify({ shown: false, currentVersion: '1.1.0', lastVersion: '' }),
                });
            } else if (route.request().method() === 'POST') {
                welcomePosted = true;
                await route.fulfill({ status: 200, body: JSON.stringify({ status: 'ok' }) });
            }
        });

        await page.goto('/');
        await expect(page.getByTestId('welcome-modal')).toBeVisible();

        // Click the Get Started / Close button
        await page.getByRole('button', { name: /get started|close|continue/i }).click();

        // Modal should be dismissed
        await expect(page.getByTestId('welcome-modal')).not.toBeVisible();

        // Should have called POST /api/welcome
        expect(welcomePosted).toBe(true);
    });

    test('does not show welcome modal if already shown for version', async ({ page }) => {
        await page.route('/api/welcome', async (route) => {
            await route.fulfill({
                status: 200,
                body: JSON.stringify({
                    shown: true,
                    currentVersion: '1.1.0',
                    lastVersion: '1.1.0'
                }),
            });
        });

        await page.goto('/');
        
        // Wait for page to load
        await page.waitForTimeout(1000);
        
        // Welcome modal should NOT be visible
        await expect(page.getByTestId('welcome-modal')).not.toBeVisible();
    });

    test('shows welcome modal on version upgrade', async ({ page }) => {
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    body: JSON.stringify({
                        shown: false, // Not shown for new version
                        currentVersion: '1.2.0',
                        lastVersion: '1.1.0' // User saw 1.1.0 before
                    }),
                });
            } else {
                await route.fulfill({ status: 200, body: JSON.stringify({ status: 'ok' }) });
            }
        });

        await page.goto('/');
        
        // Welcome modal should be visible for version upgrade
        await expect(page.getByTestId('welcome-modal')).toBeVisible({ timeout: 5000 });
    });
});
