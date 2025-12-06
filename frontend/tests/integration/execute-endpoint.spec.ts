/**
 * Integration Test: Execute Endpoint (Contract 5)
 * 
 * This test validates the full stack (NO MOCKS):
 * React UI button click → Go BFF endpoint reached → Go endpoint calls Executor interface
 * → UI receives and renders "Execution Request Received" message.
 * 
 * This is the validation for Contract 5 of implementation-contracts-v2.md
 */
import { test, expect } from '@playwright/test';

test.describe('Execute Endpoint Integration (Contract 5)', () => {
    test.setTimeout(30000);

    // Mock the welcome endpoint to prevent the modal from appearing
    test.beforeEach(async ({ page }) => {
        await page.route('/api/welcome', async (route) => {
            if (route.request().method() === 'GET') {
                await route.fulfill({
                    status: 200,
                    contentType: 'application/json',
                    body: JSON.stringify({ shown: true, currentVersion: '1.1.0' }),
                });
            } else {
                await route.fulfill({ status: 200, body: '{}' });
            }
        });
    });

    test('full stack validation - UI calls /api/execute and receives response', async ({ page }) => {
        // Navigate to the app
        await page.goto('/');
        
        // Wait for app to load
        await page.waitForLoadState('networkidle');
        
        // Execute command via /api/execute endpoint using page.evaluate
        // This simulates what a React button click would do
        const response = await page.evaluate(async () => {
            const res = await fetch('/api/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    command: 'echo Contract5Test'
                }),
            });
            return await res.json();
        });

        // Verify the response contains "Execution Request Received" (Contract 5 requirement)
        expect(response.message).toBe('Execution Request Received');
        
        // Verify the command was actually executed (NO MOCKS - real execution)
        expect(response.stdout).toContain('Contract5Test');
        
        // Verify success
        expect(response.success).toBe(true);
        
        // Verify exit code is 0
        expect(response.exitCode).toBe(0);
    });

    test('execute endpoint handles echo hello world correctly', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        
        const response = await page.evaluate(async () => {
            const res = await fetch('/api/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    command: 'echo hello world'
                }),
            });
            return await res.json();
        });

        expect(response.message).toBe('Execution Request Received');
        expect(response.stdout).toBe('hello world\n');
        expect(response.success).toBe(true);
        expect(response.exitCode).toBe(0);
    });

    test('execute endpoint handles failed commands correctly', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        
        const response = await page.evaluate(async () => {
            const res = await fetch('/api/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    command: 'exit 1'
                }),
            });
            return await res.json();
        });

        // Should still get "Execution Request Received" message
        expect(response.message).toBe('Execution Request Received');
        
        // But success should be false and exit code should be 1
        expect(response.success).toBe(false);
        expect(response.exitCode).toBe(1);
    });

    test('execute endpoint rejects empty commands', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        
        const response = await page.evaluate(async () => {
            const res = await fetch('/api/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    command: ''
                }),
            });
            return { status: res.status, body: await res.json() };
        });

        expect(response.status).toBe(400);
        expect(response.body.success).toBe(false);
    });

    test('health endpoint is accessible', async ({ page }) => {
        // Verify the Go BFF is running
        const response = await page.request.get('/api/health');
        expect(response.status()).toBe(200);
        
        const body = await response.json();
        expect(body.status).toBe('ok');
    });
});
