import { test, expect } from '@playwright/test';

test('Flows Editor Visual Verification', async ({ page }) => {
    // Navigate to Flows List
    await page.goto('/flows');
    await expect(page.getByRole('heading', { name: 'Flows' })).toBeVisible();

    // Navigate to Create New Flow
    await page.getByText('Create New Flow').click();
    await expect(page.getByText('New Flow')).toBeVisible();

    // Wait for React Flow to render
    await page.waitForTimeout(1000);

    // Drag and drop check (visual only, hard to verify logic without complex setup)
    // We just verify the canvas is there
    await expect(page.locator('.react-flow')).toBeVisible();
});
