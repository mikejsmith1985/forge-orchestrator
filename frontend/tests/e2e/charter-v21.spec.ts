import { test, expect } from '@playwright/test';

/**
 * V2.1 Charter Compliance Tests
 * 
 * Tests the new features addressing charter violations:
 * - Dynamic Budget Meter in Forge Architect
 * - Context Navigator for intelligent file selection
 * - Flow Node cost confirmation modal with exact cost
 */

test.describe('V2.1 Budget Meter', () => {
    test.beforeEach(async ({ page }) => {
        // Mock the budget API
        await page.route('/api/budget*', async route => {
            await route.fulfill({
                json: {
                    totalBudget: 10.00,
                    spentToday: 2.50,
                    remainingBudget: 7.50,
                    remainingPrompts: 750,
                    costUnit: 'TOKEN',
                    model: 'gpt-4o'
                }
            });
        });

        await page.goto('/architect');
    });

    test('displays budget meter with remaining prompts', async ({ page }) => {
        const budgetMeter = page.getByTestId('budget-meter');
        await expect(budgetMeter).toBeVisible();
        
        // Should show remaining prompts
        const remainingPrompts = page.getByTestId('remaining-prompts');
        await expect(remainingPrompts).toContainText('750');
    });

    test('shows budget progress bar', async ({ page }) => {
        const budgetBar = page.getByTestId('budget-meter-bar');
        await expect(budgetBar).toBeVisible();
        
        // 75% remaining = 75% width
        await expect(budgetBar).toHaveAttribute('style', /width: 75%/);
    });

    test('displays model name in budget meter', async ({ page }) => {
        const budgetMeter = page.getByTestId('budget-meter');
        await expect(budgetMeter).toContainText('gpt-4o');
    });
});

test.describe('V2.1 Context Navigator', () => {
    test.beforeEach(async ({ page }) => {
        // Mock budget API
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

        // Mock context files API
        await page.route('/api/context/files*', async route => {
            const url = new URL(route.request().url());
            const filter = url.searchParams.get('filter');
            
            let files;
            if (filter === 'uncommitted') {
                files = [
                    { path: 'src/App.tsx', status: 'M', type: 'file' },
                    { path: 'src/components/Test.tsx', status: 'A', type: 'file' },
                ];
            } else {
                files = [
                    { path: 'src/index.ts', type: 'file' },
                ];
            }
            
            await route.fulfill({ json: { files } });
        });

        await page.goto('/architect');
    });

    test('displays context navigator collapsed by default', async ({ page }) => {
        const navigator = page.getByTestId('context-navigator');
        await expect(navigator).toBeVisible();
        
        // Filter buttons should not be visible when collapsed
        await expect(page.getByTestId('filter-uncommitted')).not.toBeVisible();
    });

    test('expands context navigator on click', async ({ page }) => {
        const toggle = page.getByTestId('context-navigator-toggle');
        await toggle.click();
        
        // Filter buttons should now be visible
        await expect(page.getByTestId('filter-uncommitted')).toBeVisible();
        await expect(page.getByTestId('filter-recent')).toBeVisible();
        await expect(page.getByTestId('filter-source')).toBeVisible();
    });

    test('shows uncommitted files with status badges', async ({ page }) => {
        await page.getByTestId('context-navigator-toggle').click();
        
        // Wait for files to load
        await page.waitForTimeout(500);
        
        // Should show modified file
        await expect(page.getByText('src/App.tsx')).toBeVisible();
    });

    test('can select files and shows count', async ({ page }) => {
        await page.getByTestId('context-navigator-toggle').click();
        await page.waitForTimeout(500);
        
        // Click on a file to select it
        await page.getByText('src/App.tsx').click();
        
        // Should show selected count
        await expect(page.getByText('1 selected')).toBeVisible();
    });

    test('can clear selection', async ({ page }) => {
        await page.getByTestId('context-navigator-toggle').click();
        await page.waitForTimeout(500);
        
        await page.getByText('src/App.tsx').click();
        await expect(page.getByText('1 selected')).toBeVisible();
        
        await page.getByTestId('clear-selection').click();
        await expect(page.getByText('1 selected')).not.toBeVisible();
    });
});

test.describe('V2.1 Flow Node Cost Confirmation', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/flows/new');
        await page.waitForSelector('[data-testid="llm-node-drag"]', { timeout: 10000 });
    });

    test('confirmation modal shows estimated cost', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const llmNodeDrag = page.locator('[data-testid="llm-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await llmNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        await page.locator('[data-testid="agent-node"]').first().click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // Enter a prompt
        const commandTextarea = page.locator('[data-testid="config-command-textarea"]');
        await commandTextarea.fill('Explain the authentication module in detail');

        // Click Save to show confirmation
        await page.locator('[data-testid="config-save-btn"]').click();

        // Confirmation modal should show estimated cost
        const confirmModal = page.locator('[data-testid="premium-confirm-modal"]');
        await expect(confirmModal).toBeVisible();
        
        // Should show estimated tokens
        await expect(page.getByTestId('estimated-tokens')).toBeVisible();
        
        // Should show estimated cost in dollars
        await expect(page.getByTestId('estimated-cost')).toBeVisible();
        await expect(page.getByTestId('estimated-cost')).toContainText('$');
    });

    test('modal shows rate information', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const llmNodeDrag = page.locator('[data-testid="llm-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await llmNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        await page.locator('[data-testid="agent-node"]').first().click();
        
        const commandTextarea = page.locator('[data-testid="config-command-textarea"]');
        await commandTextarea.fill('Test prompt');
        
        await page.locator('[data-testid="config-save-btn"]').click();

        const confirmModal = page.locator('[data-testid="premium-confirm-modal"]');
        await expect(confirmModal).toContainText('Rate:');
        await expect(confirmModal).toContainText('1K tokens');
    });
});
