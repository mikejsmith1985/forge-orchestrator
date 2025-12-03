import { test, expect } from '@playwright/test';

test.describe('Layout', () => {
    test('should display sidebar on desktop', async ({ page }) => {
        await page.goto('/');
        const sidebar = page.getByTestId('sidebar');
        await expect(sidebar).toBeVisible();
    });

    test('should hide sidebar on mobile', async ({ page }) => {
        await page.setViewportSize({ width: 375, height: 667 });
        await page.goto('/');
        const sidebar = page.getByTestId('sidebar');
        await expect(sidebar).toBeHidden();
    });

    test('should display Forge Orchestrator title', async ({ page }) => {
        await page.goto('/');
        await expect(page.getByText('Forge Orchestrator')).toBeVisible();
    });
});
