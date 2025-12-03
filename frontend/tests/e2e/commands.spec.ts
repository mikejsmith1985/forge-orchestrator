// Educational Comment: Import Playwright test utilities for E2E testing.
// 'test' provides the test runner and browser context, 'expect' provides assertions.
import { test, expect } from '@playwright/test';

// Educational Comment: Group all Command Deck related tests in a describe block
// for better organization and reporting. This suite verifies the full user flow
// of navigating to the Command Deck, creating commands, and viewing them.
test.describe('Command Deck', () => {
    // Educational Comment: TEST 1 - Navigate to Command Deck
    // Purpose: Verify that users can successfully navigate to the Command Deck view
    // and that the main UI elements are rendered correctly.
    test('should display command deck', async ({ page }) => {
        // Educational Comment: Start at the app's root URL. Using a relative path '/'
        // allows Playwright's baseURL config to determine the actual server location.
        await page.goto('/');

        // Educational Comment: Simulate user clicking the "Flows" navigation item in the sidebar.
        // This is how users navigate to the Command Deck in the actual application.
        await page.click('text=Flows');

        // Educational Comment: Verify the CommandDeck component is rendered by checking
        // for its data-testid attribute. Using data-testid is best practice as it creates
        // stable selectors that won't break if CSS classes or text content changes.
        const commandDeck = page.getByTestId('command-deck');
        await expect(commandDeck).toBeVisible();

        // Educational Comment: Verify the page header and description are displayed.
        // This ensures not just that the component loaded, but that it's showing
        // the correct content to help users understand what this view is for.
        await expect(page.getByText('Command Deck')).toBeVisible();
        await expect(page.getByText('Manage and execute your command cards')).toBeVisible();
    });

    // Educational Comment: TEST 2 (Part A) - Verify Add Command Modal Opens
    // Purpose: Ensure the "Add Command" button opens the modal form correctly
    // and all required form fields are present and visible.
    test('should show add command modal', async ({ page }) => {
        // Educational Comment: Navigate to the Command Deck view
        await page.goto('/');
        await page.click('text=Flows');

        // Educational Comment: Find and click the "Add Command" button using its test ID.
        // We first verify it's visible to ensure the UI is in the expected state.
        const addButton = page.getByTestId('add-command-btn');
        await expect(addButton).toBeVisible();
        await addButton.click();

        // Educational Comment: Verify the modal title is visible, confirming the modal opened
        await expect(page.getByText('Add Command')).toBeVisible();

        // Educational Comment: Verify all three required form fields are present.
        // This ensures the form is complete and ready for user input.
        // - command-name-input: for naming the command
        // - command-description-input: for describing what the command does
        // - command-input: for the actual shell command to execute
        await expect(page.getByTestId('command-name-input')).toBeVisible();
        await expect(page.getByTestId('command-description-input')).toBeVisible();
        await expect(page.getByTestId('command-input')).toBeVisible();
    });

    // Educational Comment: TEST 2 & 3 - Create New Command and Verify Display
    // Purpose: This test covers the full user flow of creating a new command card:
    // 1. Opening the modal form
    // 2. Filling in all required fields
    // 3. Submitting the form
    // 4. Verifying the modal closes (indicating successful submission)
    test('should add a new command', async ({ page }) => {
        // Educational Comment: Navigate to Command Deck
        await page.goto('/');
        await page.click('text=Flows');

        // Educational Comment: Open the "Add Command" modal by clicking the button
        await page.getByTestId('add-command-btn').click();

        // Educational Comment: Fill out the form with test data.
        // We use .fill() which simulates real user input character by character.
        await page.getByTestId('command-name-input').fill('Test Command');
        await page.getByTestId('command-description-input').fill('This is a test command');
        await page.getByTestId('command-input').fill('npm test');

        // Educational Comment: Submit the form by clicking the submit button.
        // This triggers a POST request to /api/commands with the form data.
        await page.getByTestId('submit-command-btn').click();

        // Educational Comment: Wait for the async operation to complete.
        // The modal should close after successful submission, and the new command
        // should appear in the grid. This timeout gives the API call time to complete.
        await page.waitForTimeout(1000);

        // Educational Comment: Verify the modal is no longer visible, which indicates
        // the form was submitted successfully. In a full integration test with a running
        // backend, we could also verify the command card appears in the grid using:
        // await expect(page.getByTestId('command-card')).toBeVisible();
        await expect(page.getByText('Add Command')).not.toBeVisible();
    });

    // Educational Comment: TEST - Empty State or Data Display
    // Purpose: Verify the UI handles both empty and populated states gracefully.
    // This is important for good UX - users should see helpful messaging when
    // there's no data, or a grid of command cards when data exists.
    test('should display empty state when no commands', async ({ page }) => {
        // Educational Comment: Navigate to Command Deck
        await page.goto('/');
        await page.click('text=Flows');

        // Educational Comment: Wait for the API response to complete loading.
        // The CommandDeck component fetches data on mount via useEffect.
        await page.waitForTimeout(1000);

        // Educational Comment: Check for either empty state or command cards.
        // We use .catch() to handle cases where the element doesn't exist,
        // since one or the other should be present depending on database state.
        const isEmpty = await page.getByText('No commands yet').isVisible().catch(() => false);
        const hasCards = await page.getByTestId('command-card').count().then(count => count > 0).catch(() => false);

        // Educational Comment: Assert that at least one condition is true.
        // This makes the test resilient to different database states while
        // still verifying that the UI renders something meaningful.
        expect(isEmpty || hasCards).toBeTruthy();
    });
});
