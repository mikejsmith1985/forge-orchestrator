/**
 * E2E Test: Feedback Modal
 * 
 * Tests the feedback/bug reporting functionality.
 */
import { test, expect } from '@playwright/test';

test.describe('Feedback Modal', () => {
    test.setTimeout(30000);

    test('feedback button is visible in sidebar', async ({ page }) => {
        await page.goto('/');
        
        // Look for feedback button
        const feedbackButton = page.getByRole('button', { name: /feedback/i });
        await expect(feedbackButton).toBeVisible();
    });

    test('clicking feedback button opens modal', async ({ page }) => {
        await page.goto('/');
        
        // Click feedback button
        await page.getByRole('button', { name: /feedback/i }).click();
        
        // Modal should be visible
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
    });

    test('feedback modal has required fields', async ({ page }) => {
        await page.goto('/');
        await page.getByRole('button', { name: /feedback/i }).click();
        
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Should have description field
        await expect(page.getByPlaceholder(/describe|issue|feedback/i)).toBeVisible();
        
        // Should have submit button
        await expect(page.getByRole('button', { name: /submit|send/i })).toBeVisible();
        
        // Should have cancel/close button
        await expect(page.getByRole('button', { name: /cancel|close/i })).toBeVisible();
    });

    test('feedback modal can capture screenshot', async ({ page }) => {
        await page.goto('/');
        await page.getByRole('button', { name: /feedback/i }).click();
        
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Should have screenshot capture button
        const captureButton = page.getByRole('button', { name: /capture|screenshot/i });
        await expect(captureButton).toBeVisible();
    });

    test('feedback modal can be closed', async ({ page }) => {
        await page.goto('/');
        await page.getByRole('button', { name: /feedback/i }).click();
        
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Close the modal
        await page.getByRole('button', { name: /cancel|close/i }).first().click();
        
        // Modal should be hidden
        await expect(page.getByTestId('feedback-modal')).not.toBeVisible();
    });

    test('feedback modal requires description before submit', async ({ page }) => {
        await page.goto('/');
        await page.getByRole('button', { name: /feedback/i }).click();
        
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Submit button should be disabled when description is empty
        const submitButton = page.getByRole('button', { name: /submit|send/i });
        await expect(submitButton).toBeDisabled();
        
        // Type a description
        await page.getByPlaceholder(/describe|issue|feedback/i).fill('Test feedback description');
        
        // Submit button should now be enabled
        await expect(submitButton).toBeEnabled();
    });
});
