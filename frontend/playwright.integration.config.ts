import { defineConfig, devices } from '@playwright/test';

/**
 * Integration test configuration for Playwright.
 * 
 * Unlike the E2E tests which mock API responses, integration tests
 * run against a real Go backend server with an in-memory SQLite database.
 * 
 * Usage:
 * 1. Start the Go test server: go test ./tests/integration/... -run TestMain -count=1
 * 2. Set INTEGRATION_SERVER_URL environment variable
 * 3. Run: npx playwright test --config=playwright.integration.config.ts
 */
export default defineConfig({
    testDir: './tests/integration',
    fullyParallel: false, // Run serially to avoid race conditions with shared database
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 2 : 0,
    workers: 1, // Single worker for integration tests
    reporter: [
        ['html', { outputFolder: 'playwright-report-integration' }],
        ['json', { outputFile: 'test-results/integration-results.json' }]
    ],
    use: {
        // Use environment variable for server URL, fallback to localhost:8080
        baseURL: process.env.INTEGRATION_SERVER_URL || 'http://localhost:8080',
        trace: 'on-first-retry',
        screenshot: 'only-on-failure',
    },
    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],
    // No webServer config - the Go test server must be started separately
    webServer: undefined,
});
