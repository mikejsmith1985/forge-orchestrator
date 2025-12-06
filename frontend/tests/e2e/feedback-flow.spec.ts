/**
 * E2E Test: Feedback Flow
 * 
 * Tests the complete user journey for the feedback feature:
 * - Opening the feedback modal
 * - Setup view with PAT generation link
 * - Token saving
 * - Feedback submission with screenshots
 * - Error handling
 */
import { test, expect } from '@playwright/test';

test.describe('Feedback Feature', () => {
    test.beforeEach(async ({ page }) => {
        // Mock the welcome endpoint to prevent modal
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    contentType: 'application/json',
                    body: JSON.stringify({ shown: true, currentVersion: '1.1.1' }),
                });
            } else {
                await route.fulfill({ status: 200, body: '{}' });
            }
        });
        
        // Clear any saved tokens
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.removeItem('forge_github_token');
        });
        await page.waitForLoadState('networkidle');
    });

    test('should open feedback modal from sidebar', async ({ page }) => {
        // Click the feedback button in sidebar
        await page.getByTestId('feedback-button').click();
        
        // Modal should be visible
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Should show the setup view since no token is saved
        await expect(page.getByText('Setup Required')).toBeVisible();
    });

    test('should display setup view with prefilled PAT link', async ({ page }) => {
        // Open feedback modal
        await page.getByTestId('feedback-button').click();
        
        // Check for setup instructions
        await expect(page.getByText(/Personal Access Token \(PAT\)/i)).toBeVisible();
        
        // Check for the generate token link with correct attributes
        const tokenLink = page.getByTestId('generate-token-link');
        await expect(tokenLink).toBeVisible();
        await expect(tokenLink).toHaveText(/Generate Token on GitHub/i);
        
        // Verify the link has correct URL with prefilled parameters
        const href = await tokenLink.getAttribute('href');
        expect(href).toContain('github.com/settings/tokens/new');
        expect(href).toContain('scopes=public_repo');
        expect(href).toContain('description=Forge+Orchestrator+Feedback');
        
        // Verify it opens in new tab
        await expect(tokenLink).toHaveAttribute('target', '_blank');
        
        // Check for scope information
        await expect(page.getByText('public_repo')).toBeVisible();
    });

    test('should show setup instructions and requirements', async ({ page }) => {
        await page.getByTestId('feedback-button').click();
        
        // Should explain what PAT is needed for
        await expect(page.getByText(/Create issues in the forge-orchestrator repository/i)).toBeVisible();
        await expect(page.getByText(/Upload screenshots/i)).toBeVisible();
    });

    test('should save GitHub token', async ({ page }) => {
        await page.getByTestId('feedback-button').click();
        
        // Enter a test token
        const tokenInput = page.getByTestId('github-token-input');
        await tokenInput.fill('ghp_test_token_12345');
        
        // Click save settings
        await page.getByRole('button', { name: /Save Settings/i }).click();
        
        // Should switch to feedback view
        await expect(page.getByText('Issue Description')).toBeVisible();
        await expect(page.getByText('Setup Required')).not.toBeVisible();
        
        // Token should be saved in localStorage
        const savedToken = await page.evaluate(() => {
            return localStorage.getItem('forge_github_token');
        });
        expect(savedToken).toBe('ghp_test_token_12345');
    });

    test('should require token before saving', async ({ page }) => {
        await page.getByTestId('feedback-button').click();
        
        // Try to save without entering token
        await page.getByRole('button', { name: /Save Settings/i }).click();
        
        // Should show error
        await expect(page.getByText(/GitHub Token is required/i)).toBeVisible();
        
        // Should still be in setup view
        await expect(page.getByText('Setup Required')).toBeVisible();
    });

    test('should allow updating token settings', async ({ page }) => {
        // Save a token first
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_old_token');
        });
        
        // Open feedback modal
        await page.getByTestId('feedback-button').click();
        
        // Should show feedback view
        await expect(page.getByText('Issue Description')).toBeVisible();
        
        // Click update settings
        await page.getByRole('button', { name: /Update Settings/i }).click();
        
        // Should go back to setup view
        await expect(page.getByText('Setup Required')).toBeVisible();
        
        // Token input should have the old token
        const tokenInput = page.getByTestId('github-token-input');
        await expect(tokenInput).toHaveValue('ghp_old_token');
    });

    test('should show feedback form when token is configured', async ({ page }) => {
        // Pre-configure token
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        // Open feedback modal
        await page.getByTestId('feedback-button').click();
        
        // Should show feedback form
        await expect(page.getByText('Issue Description')).toBeVisible();
        await expect(page.getByPlaceholder(/Describe the issue/i)).toBeVisible();
        
        // Should have screenshot section
        await expect(page.getByRole('button', { name: /Capture Screen/i })).toBeVisible();
        
        // Should have submit button (disabled when empty)
        const submitButton = page.getByRole('button', { name: /Submit Feedback/i });
        await expect(submitButton).toBeVisible();
        await expect(submitButton).toBeDisabled();
    });

    test('should enable submit button when description is entered', async ({ page }) => {
        // Pre-configure token
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        
        const submitButton = page.getByRole('button', { name: /Submit Feedback/i });
        await expect(submitButton).toBeDisabled();
        
        // Enter description
        await page.getByPlaceholder(/Describe the issue/i).fill('Test feedback issue');
        
        // Submit should now be enabled
        await expect(submitButton).toBeEnabled();
    });

    test('should handle feedback submission with mocked API', async ({ page }) => {
        // Mock the GitHub API endpoints - use correct pattern
        await page.route('https://api.github.com/repos/mikejsmith1985/forge-orchestrator/issues', async (route) => {
            await route.fulfill({
                status: 201,
                contentType: 'application/json',
                body: JSON.stringify({
                    html_url: 'https://github.com/mikejsmith1985/forge-orchestrator/issues/123',
                    number: 123
                })
            });
        });
        
        // Pre-configure token
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        
        // Wait for modal to be visible
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Enter feedback
        await page.getByPlaceholder(/Describe the issue/i).fill('This is a test issue for the e2e test');
        
        // Submit
        await page.getByRole('button', { name: /Submit Feedback/i }).click();
        
        // Should show processing message first
        await expect(page.getByText(/Processing|Creating/i)).toBeVisible({ timeout: 2000 });
        
        // Should show success message
        await expect(page.getByText(/Issue #123 created/i)).toBeVisible({ timeout: 10000 });
        
        // Should show link to issue
        const issueLink = page.getByRole('link', { name: /Issue #123 created/i });
        await expect(issueLink).toBeVisible();
        await expect(issueLink).toHaveAttribute('href', 'https://github.com/mikejsmith1985/forge-orchestrator/issues/123');
    });

    test('should handle API authentication errors gracefully', async ({ page }) => {
        // Mock 401 response
        await page.route('https://api.github.com/repos/mikejsmith1985/forge-orchestrator/issues', async (route) => {
            await route.fulfill({
                status: 401,
                contentType: 'application/json',
                body: JSON.stringify({ message: 'Bad credentials' })
            });
        });
        
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_invalid_token');
        });
        
        await page.getByTestId('feedback-button').click();
        await page.getByPlaceholder(/Describe the issue/i).fill('Test issue');
        await page.getByRole('button', { name: /Submit Feedback/i }).click();
        
        // Should show error message with helpful text
        await expect(page.getByText(/Invalid GitHub token/i)).toBeVisible({ timeout: 5000 });
    });

    test('should handle permission errors', async ({ page }) => {
        // Mock 403 response
        await page.route('https://api.github.com/repos/mikejsmith1985/forge-orchestrator/issues', async (route) => {
            await route.fulfill({
                status: 403,
                contentType: 'application/json',
                body: JSON.stringify({ message: 'Forbidden' })
            });
        });
        
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        await page.getByPlaceholder(/Describe the issue/i).fill('Test issue');
        await page.getByRole('button', { name: /Submit Feedback/i }).click();
        
        // Should show permission error (check for key parts of the message)
        await expect(page.getByText(/lacks permissions/i)).toBeVisible({ timeout: 10000 });
        await expect(page.getByText(/public_repo/i)).toBeVisible();
    });

    test('should close modal on escape key', async ({ page }) => {
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Press escape
        await page.keyboard.press('Escape');
        
        // Modal should close
        await expect(page.getByTestId('feedback-modal')).not.toBeVisible();
    });

    test('should close modal on backdrop click', async ({ page }) => {
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Click backdrop (outside modal content)
        await page.getByTestId('feedback-modal').click({ position: { x: 10, y: 10 } });
        
        // Modal should close
        await expect(page.getByTestId('feedback-modal')).not.toBeVisible();
    });

    test('should close modal on X button click', async ({ page }) => {
        await page.goto('/');
        await page.evaluate(() => {
            localStorage.setItem('forge_github_token', 'ghp_test_token');
        });
        
        await page.getByTestId('feedback-button').click();
        await expect(page.getByTestId('feedback-modal')).toBeVisible();
        
        // Click close button
        await page.getByRole('button', { name: /Close/i }).click();
        
        // Modal should close
        await expect(page.getByTestId('feedback-modal')).not.toBeVisible();
    });
});
