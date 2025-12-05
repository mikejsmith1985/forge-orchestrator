import { test, expect } from '@playwright/test';

test.describe('Flows Editor Visual Verification', () => {
    test.beforeEach(async ({ page }) => {
        // Mock the flows API
        await page.route('/api/flows', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([])
            });
        });
        await page.route('/api/flows/*', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ id: 1, name: 'Test Flow', data: '{}', status: 'active' })
            });
        });
    });

    test('can navigate to flows and create new flow', async ({ page }) => {
        // Navigate to Flows List
        await page.goto('/flows');
        await expect(page.getByRole('heading', { name: 'Flows' })).toBeVisible();

        // Navigate to Create New Flow
        await page.getByRole('button', { name: 'Create New Flow' }).first().click();
        
        // Verify we're on the flow editor by checking for the flow name input
        await expect(page.locator('input[placeholder="Flow Name"]')).toBeVisible();

        // Wait for React Flow to render
        await page.waitForTimeout(1000);

        // Drag and drop check (visual only, hard to verify logic without complex setup)
        // We just verify the canvas is there
        await expect(page.locator('.react-flow')).toBeVisible();
    });
});
