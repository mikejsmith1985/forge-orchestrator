import { test, expect } from '@playwright/test';

test.describe('Forge Architect View', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/');
    });

    test('should update token meter width when typing', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter"] .bg-gray-800 > div');

        // Initial state should be 0%
        await expect(meterBar).toHaveCSS('width', '0px');

        // Type 800 characters (200 tokens) -> 200/8000 = 2.5%
        // We'll type enough to get a measurable width.
        // 8000 tokens max. 4 chars = 1 token.
        // To get 50% (4000 tokens), we need 16000 characters.
        // Let's try 10% (800 tokens) -> 3200 characters.

        const text = 'a'.repeat(3200);
        await input.fill(text);

        // Verify width is approximately 10%
        // The width is set as a style percentage.
        await expect(meterBar).toHaveAttribute('style', /width: 10%/);
    });

    test('should change color to red when pasting large text', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter"] .bg-gray-800 > div');

        // Initial color should be green
        await expect(meterBar).toHaveClass(/bg-green-500/);

        // Max tokens = 8000.
        // > 90% is Red. 90% of 8000 = 7200 tokens.
        // 7200 tokens * 4 chars/token = 28800 chars.
        // Let's paste 30000 chars to be safe.
        const largeText = 'a'.repeat(30000);

        // Use evaluate to set value quickly for "paste" simulation or just fill
        await input.fill(largeText);

        // Verify color class changes to red
        await expect(meterBar).toHaveClass(/bg-red-500/);
    });
});
