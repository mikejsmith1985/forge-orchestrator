import { test, expect } from '@playwright/test';

test.describe('Forge Architect View', () => {
    test.beforeEach(async ({ page }) => {
        // Mock budget API for consistent test behavior
        await page.route('/api/budget*', async route => {
            await route.fulfill({
                json: {
                    totalBudget: 10.00,
                    spentToday: 0,
                    remainingBudget: 10.00,
                    remainingPrompts: 1000,
                    costUnit: 'TOKEN',
                    model: 'gpt-4o'
                }
            });
        });
        await page.goto('/architect');
    });

    test('should display architect input and token meter', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meter = page.getByTestId('token-meter');
        
        await expect(input).toBeVisible();
        await expect(meter).toBeVisible();
    });

    test('should update token meter width when typing', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter-bar"]');

        // Initial state should be 0%
        await expect(meterBar).toHaveCSS('width', '0px');

        // tiktoken calculates tokens more accurately than /4 heuristic
        // Type enough text to get a measurable width change (>= 1%)
        const text = 'a'.repeat(3200);
        await input.fill(text);

        // Wait for debounced API call and verify width changes from 0
        // Using regex to match any non-zero percentage
        await expect(meterBar).toHaveAttribute('style', /width: [1-9][0-9]?%/, { timeout: 10000 });
    });

    test('should change color to red when pasting large text', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter-bar"]');

        // Initial color should be green
        await expect(meterBar).toHaveClass(/bg-green-500/);

        // Max tokens = 8000. > 90% is Red.
        // tiktoken counts ~8 chars per token for repeated 'a'
        // Need > 7200 tokens for red = > 57600 chars
        const largeText = 'a'.repeat(70000);
        await input.fill(largeText);

        // Wait for debounced API call and verify color class changes to red
        await expect(meterBar).toHaveClass(/bg-red-500/, { timeout: 10000 });
    });
});

// Mobile Responsiveness Tests (Contract #42)
test.describe('Architect View - Mobile (360x640)', () => {
    test.beforeEach(async ({ page }) => {
        // Mock budget API for consistent test behavior
        await page.route('/api/budget*', async route => {
            await route.fulfill({
                json: {
                    totalBudget: 10.00,
                    spentToday: 0,
                    remainingBudget: 10.00,
                    remainingPrompts: 1000,
                    costUnit: 'TOKEN',
                    model: 'gpt-4o'
                }
            });
        });
        await page.setViewportSize({ width: 360, height: 640 });
        await page.goto('/architect');
    });

    test('should display architect input on mobile', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        await expect(input).toBeVisible();
        
        // Verify it has reasonable width for mobile (accounting for layout padding)
        const box = await input.boundingBox();
        expect(box?.width).toBeGreaterThan(200);
    });

    test('should display token meter on mobile', async ({ page }) => {
        const meter = page.getByTestId('token-meter');
        await expect(meter).toBeVisible();
        
        // Should be below the textarea, not side-by-side
        const input = page.getByTestId('architect-input');
        const inputBox = await input.boundingBox();
        const meterBox = await meter.boundingBox();
        
        expect(meterBox?.y).toBeGreaterThan(inputBox!.y + inputBox!.height - 10);
    });

    test('should update meter when typing on mobile', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter-bar"]');
        
        // Click to focus (use click instead of tap for browser compatibility)
        await input.click();
        await input.fill('a'.repeat(3200));
        
        // Wait for debounced API call and verify width changes
        await expect(meterBar).toHaveAttribute('style', /width: [1-9][0-9]?%/, { timeout: 10000 });
    });

    test('should not have horizontal scroll', async ({ page }) => {
        // Check that body width equals viewport width (no horizontal overflow)
        const bodyWidth = await page.evaluate(() => document.body.scrollWidth);
        expect(bodyWidth).toBeLessThanOrEqual(360);
    });

    test('should have readable text (font size check)', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const fontSize = await input.evaluate((el) => {
            return parseFloat(window.getComputedStyle(el).fontSize);
        });
        
        // Text should be at least 14px for readability on mobile
        expect(fontSize).toBeGreaterThanOrEqual(14);
    });

    test('should have adequate touch target size', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const box = await input.boundingBox();
        
        // Touch targets should be at least 44px (Apple HIG guideline)
        expect(box?.height).toBeGreaterThanOrEqual(44);
    });
});

test.describe('Architect View - Large Mobile (414x896 iPhone 11 Pro Max)', () => {
    test.beforeEach(async ({ page }) => {
        // Mock budget API for consistent test behavior
        await page.route('/api/budget*', async route => {
            await route.fulfill({
                json: {
                    totalBudget: 10.00,
                    spentToday: 0,
                    remainingBudget: 10.00,
                    remainingPrompts: 1000,
                    costUnit: 'TOKEN',
                    model: 'gpt-4o'
                }
            });
        });
        await page.setViewportSize({ width: 414, height: 896 });
        await page.goto('/architect');
    });

    test('should work on iPhone 11 Pro Max size', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        await expect(input).toBeVisible();
        
        // Verify input has reasonable width for this viewport
        const box = await input.boundingBox();
        expect(box?.width).toBeGreaterThan(250);
    });

    test('should handle large text paste and show red meter', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        
        // Fill with large text to trigger red color (> 90% of 8000 tokens)
        // tiktoken counts ~8 chars per token for repeated 'a'
        await input.fill('a'.repeat(70000));
        
        const meterBar = page.locator('[data-testid="token-meter-bar"]');
        // Wait for debounced API call
        await expect(meterBar).toHaveClass(/bg-red-500/, { timeout: 10000 });
    });

    test('should not have horizontal scroll on large mobile', async ({ page }) => {
        const bodyWidth = await page.evaluate(() => document.body.scrollWidth);
        expect(bodyWidth).toBeLessThanOrEqual(414);
    });

    test('typing works correctly on mobile viewport', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        
        // Click to focus and type (use click instead of tap for browser compatibility)
        await input.click();
        await input.fill('Hello from mobile viewport test');
        
        await expect(input).toHaveValue('Hello from mobile viewport test');
    });
});
