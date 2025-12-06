import { test, expect } from '@playwright/test';

test.describe('Key Management UI', () => {
    test.beforeEach(async ({ page }) => {
        // Mock the status endpoint
        await page.route('/api/keys/status', async (route) => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                json: {
                    keys: [
                        { provider: 'anthropic', isSet: true },
                        { provider: 'openai', isSet: false },
                    ],
                },
            });
        });

        // Mock the update endpoint
        await page.route('/api/keys', async (route) => {
            const body = JSON.parse(route.request().postData() || '{}');
            if (body.provider === 'openai' && body.key === 'sk-test-123') {
                await route.fulfill({ status: 200 });
            } else {
                await route.fulfill({ status: 400 });
            }
        });

        await page.goto('/');
    });

    test('should navigate to settings page', async ({ page }) => {
        // Click settings in sidebar
        await page.getByRole('button', { name: 'Settings' }).click();

        // Wait for loading to finish
        await expect(page.locator('.animate-spin')).not.toBeVisible();

        // Verify header
        await expect(page.getByRole('heading', { name: 'Key Management' })).toBeVisible();
        await expect(page.getByText('Securely manage your API keys')).toBeVisible();
    });

    /**
     * Task 4.1: Security Assurance Text Test
     * Verifies the security assurance message is displayed prominently
     */
    test('displays security assurance message (Task 4.1)', async ({ page }) => {
        await page.getByRole('button', { name: 'Settings' }).click();
        await expect(page.locator('.animate-spin')).not.toBeVisible();

        // Verify security assurance box is visible
        const securityBox = page.locator('[data-testid="security-assurance"]');
        await expect(securityBox).toBeVisible();

        // Verify security message content
        await expect(securityBox).toContainText('Secure Storage');
        await expect(securityBox).toContainText('encrypted');
        await expect(securityBox).toContainText('operating system');
        await expect(securityBox).toContainText('keyring');
        await expect(securityBox).toContainText('never');
    });

    test('should display correct key status', async ({ page }) => {
        await page.getByRole('button', { name: 'Settings' }).click();
        await expect(page.locator('.animate-spin')).not.toBeVisible();

        // Verify Anthropic is configured (green check)
        const anthropicSection = page.getByTestId('provider-card-anthropic');
        await expect(anthropicSection.getByText('Configured')).toBeVisible();
        await expect(anthropicSection.getByText('Key is currently set')).toBeVisible();

        // Verify OpenAI is not configured (red x)
        const openaiSection = page.getByTestId('provider-card-openai');
        await expect(openaiSection.getByText('Not Configured')).toBeVisible();
        await expect(openaiSection.getByText('Enter your API key')).toBeVisible();
    });

    test('should allow updating a key', async ({ page }) => {
        // Override the status mock for this test to handle the state change
        let isKeySet = false;
        await page.route('/api/keys/status', async (route) => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                json: {
                    keys: [
                        { provider: 'anthropic', isSet: true },
                        { provider: 'openai', isSet: isKeySet },
                    ],
                },
            });
        });

        await page.getByRole('button', { name: 'Settings' }).click();
        await expect(page.locator('.animate-spin')).not.toBeVisible();

        const openaiSection = page.getByTestId('provider-card-openai');
        const input = openaiSection.getByPlaceholder('Enter openai API Key');
        const saveBtn = openaiSection.getByRole('button', { name: 'Save Key' });

        // Verify initial state
        await expect(saveBtn).toBeDisabled();
        await expect(openaiSection.getByText('Not Configured')).toBeVisible();

        // Enter key
        await input.fill('sk-test-123');
        await expect(saveBtn).toBeEnabled();

        // Update our local state when the POST happens
        await page.route('/api/keys', async (route) => {
            const body = JSON.parse(route.request().postData() || '{}');
            if (body.provider === 'openai' && body.key === 'sk-test-123') {
                isKeySet = true; // Update state for next GET
                await route.fulfill({ 
                    status: 200,
                    contentType: 'application/json',
                    json: { status: 'ok', message: 'API key saved successfully' }
                });
            } else {
                await route.fulfill({ status: 400 });
            }
        });

        // Save
        await saveBtn.click();

        // Wait for the status to be refetched (UI will fetch after save)
        await page.waitForTimeout(500);

        // Verify input is cleared after save
        await expect(input).toBeEmpty();

        // Verify status changes to Configured
        // The UI should re-fetch status after save
        await expect(openaiSection.getByText('Configured')).toBeVisible();
        await expect(openaiSection.getByText('Key is currently set')).toBeVisible();
    });

    test('should mask api keys', async ({ page }) => {
        await page.getByRole('button', { name: 'Settings' }).click();

        const input = page.getByPlaceholder('Enter anthropic API Key');
        await expect(input).toHaveAttribute('type', 'password');
    });
});
