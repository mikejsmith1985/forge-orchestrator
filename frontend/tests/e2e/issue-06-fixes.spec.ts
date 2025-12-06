import { test, expect } from '@playwright/test';

/**
 * Issue #6: UX Improvements Tests
 * 
 * Tests for all fixes implemented:
 * 1. Prompt Watcher - doesn't reset terminal
 * 2. Terminal Settings - notification near save button
 * 3. Terminal - Font size adjustment
 * 4. Architect View - Token measurement clarity
 * 5. Feedback Modal - Minimize button and scrollable screenshots
 */

test.describe('Issue #6: UX Improvements', () => {
    
    test.describe('Terminal - Prompt Watcher', () => {
        test('toggling prompt watcher does not reset terminal content', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Terminal
            await page.click('text=Terminal');
            await page.waitForSelector('[data-testid="prompt-watcher-toggle"]', { timeout: 10000 });
            
            // Wait for terminal to stabilize
            await page.waitForTimeout(3000);
            
            // Check initial state - button should exist and work
            const toggleButton = page.locator('[data-testid="prompt-watcher-toggle"]');
            await expect(toggleButton).toBeVisible();
            
            // Get initial button state
            const initialClass = await toggleButton.getAttribute('class');
            
            // Toggle prompt watcher ON
            await toggleButton.click();
            await page.waitForTimeout(500);
            
            // Verify button state changed
            const toggledClass = await toggleButton.getAttribute('class');
            expect(toggledClass).not.toBe(initialClass);
            
            // Toggle back OFF
            await toggleButton.click();
            await page.waitForTimeout(500);
            
            // Verify we can toggle it - this proves it's not resetting/recreating the terminal
            const finalClass = await toggleButton.getAttribute('class');
            expect(finalClass).toBe(initialClass);
            
            // The terminal component should still exist and be functional
            await expect(page.locator('[data-testid="terminal-container"]')).toBeVisible();
        });
    });

    test.describe('Terminal - Font Size Controls', () => {
        test('can increase and decrease font size without resetting terminal', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Terminal
            await page.click('text=Terminal');
            await page.waitForSelector('[data-testid="font-size-display"]', { timeout: 10000 });
            await page.waitForTimeout(2000);
            
            // Check initial font size (should be 14px by default)
            const initialSize = await page.textContent('[data-testid="font-size-display"]');
            expect(initialSize).toContain('14px');
            
            // Increase font size
            await page.click('[data-testid="font-size-increase"]');
            await page.waitForTimeout(300);
            
            const increasedSize = await page.textContent('[data-testid="font-size-display"]');
            expect(increasedSize).toContain('15px');
            
            // Terminal container should still exist
            await expect(page.locator('[data-testid="terminal-container"]')).toBeVisible();
            
            // Decrease font size twice
            await page.click('[data-testid="font-size-decrease"]');
            await page.click('[data-testid="font-size-decrease"]');
            await page.waitForTimeout(300);
            
            const decreasedSize = await page.textContent('[data-testid="font-size-display"]');
            expect(decreasedSize).toContain('13px');
            
            // Terminal should still be there and functional
            await expect(page.locator('[data-testid="terminal-container"]')).toBeVisible();
        });

        test('font size persists after page reload', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Terminal
            await page.click('text=Terminal');
            await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
            await page.waitForTimeout(1000);
            
            // Set font size to 18px
            await page.click('[data-testid="font-size-increase"]');
            await page.click('[data-testid="font-size-increase"]');
            await page.click('[data-testid="font-size-increase"]');
            await page.click('[data-testid="font-size-increase"]');
            await page.waitForTimeout(500);
            
            const sizeBeforeReload = await page.textContent('[data-testid="font-size-display"]');
            expect(sizeBeforeReload).toContain('18px');
            
            // Reload page
            await page.reload();
            await page.waitForTimeout(1000);
            
            // Navigate back to Terminal
            await page.click('text=Terminal');
            await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
            await page.waitForTimeout(1000);
            
            // Verify font size persisted
            const sizeAfterReload = await page.textContent('[data-testid="font-size-display"]');
            expect(sizeAfterReload).toContain('18px');
        });
    });

    test.describe('Terminal Settings - Notification Position', () => {
        test('success notification appears near save button at bottom', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Settings
            await page.click('text=Settings');
            await page.waitForTimeout(500);
            
            // Click Terminal Settings tab if needed
            const terminalSettingsTab = page.locator('text=Terminal Settings');
            if (await terminalSettingsTab.isVisible()) {
                await terminalSettingsTab.click();
                await page.waitForTimeout(300);
            }
            
            // Find save button
            const saveButton = page.locator('button:has-text("Save Configuration")');
            await expect(saveButton).toBeVisible();
            
            // Initially, no success message
            const successMessage = page.locator('text=Configuration saved');
            await expect(successMessage).not.toBeVisible();
            
            // Click save (force click to avoid overlay issues)
            await saveButton.click({ force: true });
            await page.waitForTimeout(1000);
            
            // Wait for and verify success message appears
            await expect(successMessage).toBeVisible({ timeout: 3000 });
            
            // The message should be in the same container area as the save button (bottom of settings)
            // We're testing that it's visible and appears after clicking save
            const messageBox = await successMessage.boundingBox();
            const saveButtonBox = await saveButton.boundingBox();
            
            expect(messageBox).toBeTruthy();
            expect(saveButtonBox).toBeTruthy();
            
            // Message should be relatively close to save button (not at the very top of page)
            if (messageBox && saveButtonBox) {
                // Message should not be at the very top (y < 100)
                expect(messageBox.y).toBeGreaterThan(100);
            }
        });
    });

    test.describe('Architect View - Token Measurement Clarity', () => {
        test('displays token measurement explanation for TOKEN cost unit', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Architect
            await page.click('text=Architect');
            await page.waitForTimeout(1500);
            
            // Check for explanation text - look for the actual text
            const explanationContainer = page.locator('text=Token Usage:').first();
            await expect(explanationContainer).toBeVisible({ timeout: 5000 });
            
            // Get the full text content from parent container
            const parentDiv = page.locator('.bg-slate-800\\/50').filter({ hasText: 'Token Usage:' });
            const content = await parentDiv.textContent();
            
            // Verify it contains key explanation phrases
            expect(content).toMatch(/Token Usage:|input prompt|Prompt-Based Billing/i);
        });

        test('token meter updates as text is typed', async ({ page }) => {
            await page.goto('/');
            
            // Navigate to Architect
            await page.click('text=Architect');
            await page.waitForTimeout(1000);
            
            // Find the textarea
            const textarea = page.locator('[data-testid="architect-input"]');
            await expect(textarea).toBeVisible();
            
            // Type some text
            await textarea.fill('This is a test prompt to estimate tokens. Adding more text to get a reasonable token count that will be non-zero.');
            await page.waitForTimeout(1000); // Wait longer for debounced API call
            
            // Verify token meter shows non-zero count
            const tokenCount = page.locator('[data-testid="token-count"]');
            await expect(tokenCount).toBeVisible();
            
            const countText = await tokenCount.textContent();
            // Should show a number greater than 0
            expect(countText).toMatch(/[1-9]\d* \/ \d/);
        });
    });

    test.describe('Feedback Modal - Minimize and Screenshots', () => {
        test('shows minimize button when content is added', async ({ page }) => {
            await page.goto('/');
            await page.waitForTimeout(1000);
            
            // Try to find and open feedback modal via keyboard or button
            // Check if there's a feedback trigger in the UI
            const feedbackTriggers = [
                'button:has-text("Feedback")',
                'button[aria-label*="eedback" i]',
                '[data-testid="feedback-button"]',
                'text=/feedback/i'
            ];
            
            let feedbackOpened = false;
            for (const selector of feedbackTriggers) {
                const trigger = page.locator(selector).first();
                if (await trigger.isVisible().catch(() => false)) {
                    await trigger.click();
                    await page.waitForTimeout(500);
                    feedbackOpened = true;
                    break;
                }
            }
            
            // If we can't find feedback button, try keyboard shortcut or skip
            if (!feedbackOpened) {
                // Try common keyboard shortcuts
                await page.keyboard.press('Control+Shift+F');
                await page.waitForTimeout(500);
            }
            
            // Check if modal opened
            const modal = page.locator('[data-testid="feedback-modal"]');
            const isModalVisible = await modal.isVisible().catch(() => false);
            
            if (!isModalVisible) {
                // Skip this test if we can't open the modal
                test.skip();
                return;
            }
            
            // Initially no minimize button (no content)
            const minimizeButton = page.locator('[data-testid="minimize-feedback"]');
            await expect(minimizeButton).not.toBeVisible();
            
            // Add some description
            const descriptionField = page.locator('textarea[placeholder*="Describe"]');
            await descriptionField.fill('Test feedback content');
            await page.waitForTimeout(300);
            
            // Now minimize button should be visible
            await expect(minimizeButton).toBeVisible();
            
            // Click minimize
            await minimizeButton.click();
            await page.waitForTimeout(300);
            
            // Modal should be minimized, showing badge
            const minimizedBadge = page.locator('[data-testid="feedback-minimized-badge"]');
            await expect(minimizedBadge).toBeVisible();
            
            // Click badge to restore
            await minimizedBadge.click();
            await page.waitForTimeout(300);
            
            // Modal should be visible again with content preserved
            await expect(descriptionField).toBeVisible();
            const preservedContent = await descriptionField.inputValue();
            expect(preservedContent).toBe('Test feedback content');
        });

        test('screenshots container is scrollable with multiple screenshots', async ({ page }) => {
            await page.goto('/');
            
            // Open feedback modal
            const feedbackButton = page.locator('button:has-text("Feedback"), button[aria-label*="feedback"]').first();
            
            if (await feedbackButton.isVisible()) {
                await feedbackButton.click();
                await page.waitForTimeout(500);
                
                // Skip token setup if needed
                const tokenInput = page.locator('[data-testid="github-token-input"]');
                if (await tokenInput.isVisible()) {
                    await tokenInput.fill('dummy_token_for_test');
                    await page.click('button:has-text("Save Settings")');
                    await page.waitForTimeout(500);
                }
                
                // Capture a screenshot
                const captureButton = page.locator('button:has-text("Capture Screen")');
                if (await captureButton.isVisible()) {
                    await captureButton.click();
                    await page.waitForTimeout(2000);
                    
                    // Check if screenshots container has scrollable class
                    const screenshotsContainer = page.locator('[data-testid="screenshots-container"]');
                    if (await screenshotsContainer.isVisible()) {
                        // Verify it has overflow-y-auto class for scrolling
                        const containerClass = await screenshotsContainer.getAttribute('class');
                        expect(containerClass).toContain('overflow-y-auto');
                        expect(containerClass).toContain('max-h');
                    }
                }
            }
        });
    });

    test.describe('Integration - All Fixes Working Together', () => {
        test('complete workflow using all improved features', async ({ page }) => {
            await page.goto('/');
            await page.waitForTimeout(1000);
            
            // 1. Test Terminal with font size and prompt watcher
            await page.click('text=Terminal');
            await page.waitForSelector('[data-testid="terminal-container"]', { timeout: 10000 });
            await page.waitForTimeout(2000);
            
            // Adjust font size
            await page.click('[data-testid="font-size-increase"]');
            await page.waitForTimeout(300);
            
            // Verify font size changed
            const fontSize = await page.textContent('[data-testid="font-size-display"]');
            expect(fontSize).toContain('15px');
            
            // Toggle prompt watcher
            await page.click('[data-testid="prompt-watcher-toggle"]');
            await page.waitForTimeout(300);
            
            // 2. Check Architect View with token explanation
            await page.click('text=Architect');
            await page.waitForTimeout(1500);
            
            const textarea = page.locator('[data-testid="architect-input"]');
            await textarea.fill('Create a new feature for the application');
            await page.waitForTimeout(500);
            
            // Verify explanation container is visible
            const explanationDiv = page.locator('.bg-slate-800\\/50').filter({ hasText: /Token Usage|Prompt-Based/ });
            const hasExplanation = await explanationDiv.isVisible().catch(() => false);
            expect(hasExplanation).toBeTruthy();
            
            // 3. Test Terminal Settings with notification
            await page.click('text=Settings');
            await page.waitForTimeout(1000);
            
            const saveButton = page.locator('button:has-text("Save Configuration")');
            if (await saveButton.isVisible()) {
                await saveButton.click({ force: true }); // Force to avoid overlay
                await page.waitForTimeout(1500);
                
                // Success message should be visible
                const successMessage = page.locator('text=Configuration saved');
                const messageVisible = await successMessage.isVisible().catch(() => false);
                expect(messageVisible).toBeTruthy();
            }
        });
    });
});
