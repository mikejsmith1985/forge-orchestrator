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
        await page.getByRole('button', { name: 'Settings' }).click();
        await expect(page.locator('.animate-spin')).not.toBeVisible();

        const openaiSection = page.getByTestId('provider-card-openai');
        const input = openaiSection.getByPlaceholder('Enter openai API Key');
        const saveBtn = openaiSection.getByRole('button', { name: 'Save Key' });

        // Verify initial state
        await expect(saveBtn).toBeDisabled();

        // Enter key
        await input.fill('sk-test-123');
        await expect(saveBtn).toBeEnabled();

        // Save
        await saveBtn.click();

        // Verify input is cleared after save (mock implementation in component does this)
        await expect(input).toBeEmpty();
    });

    test('should mask api keys', async ({ page }) => {
        await page.getByRole('button', { name: 'Settings' }).click();

        const input = page.getByPlaceholder('Enter anthropic API Key');
        await expect(input).toHaveAttribute('type', 'password');
    });
});
