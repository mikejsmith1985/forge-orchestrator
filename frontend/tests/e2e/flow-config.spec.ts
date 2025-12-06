import { test, expect } from '@playwright/test';

/**
 * Flow Node Configuration Tests - Updated for V2.1 Remediation Plan
 * 
 * Task 3.1-3.3: Tests the redesigned Flow Editor with:
 * - Simplified node configuration (command input instead of dropdowns)
 * - Two distinct node types: Shell Command and LLM Prompt
 * - Premium execution confirmation modal for LLM nodes
 */

test.describe('Flow Node Configuration', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to flow editor - create a new flow
        await page.goto('/flows/new');
        // Wait for the new node types to be visible
        await page.waitForSelector('[data-testid="shell-node-drag"]', { timeout: 10000 });
    });

    test('displays both Shell and LLM node types in sidebar', async ({ page }) => {
        // Verify Shell Command node is in sidebar
        await expect(page.locator('[data-testid="shell-node-drag"]')).toBeVisible();
        await expect(page.locator('[data-testid="shell-node-drag"]')).toContainText('Shell Command');
        await expect(page.locator('[data-testid="shell-node-drag"]')).toContainText('Zero-Token');
        
        // Verify LLM Prompt node is in sidebar
        await expect(page.locator('[data-testid="llm-node-drag"]')).toBeVisible();
        await expect(page.locator('[data-testid="llm-node-drag"]')).toContainText('LLM Prompt');
        await expect(page.locator('[data-testid="llm-node-drag"]')).toContainText('Premium');
    });

    test('Add Shell node → Configure → Save → Node shows updated label', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag a Shell node to the canvas
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Fill in the label
        const labelInput = page.locator('[data-testid="config-label-input"]');
        await labelInput.clear();
        await labelInput.fill('Run Tests');

        // Enter a shell command
        const commandTextarea = page.locator('[data-testid="config-command-textarea"]');
        await commandTextarea.fill('npm test');

        // Click Save
        await page.locator('[data-testid="config-save-btn"]').click();

        // Verify panel closes and label is updated
        await expect(configPanel).not.toBeVisible();
        await expect(page.locator('text=Run Tests')).toBeVisible();
    });

    test('Shell node type shows Zero-Token indicator', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await agentNode.click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // Shell type should be selected by default
        const shellTypeBtn = page.locator('[data-testid="node-type-shell"]');
        await expect(shellTypeBtn).toHaveClass(/border-green-500/);

        // Should show zero-token message
        await expect(configPanel).toContainText('no token consumption');
    });

    test('LLM node type shows Premium indicator and token meter', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const llmNodeDrag = page.locator('[data-testid="llm-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await llmNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await agentNode.click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // LLM type should be selected
        const llmTypeBtn = page.locator('[data-testid="node-type-llm"]');
        await expect(llmTypeBtn).toHaveClass(/border-purple-500/);

        // Should show premium message
        await expect(configPanel).toContainText('consumes tokens');

        // Token meter should be visible
        await expect(configPanel.locator('[data-testid="token-meter"]')).toBeVisible();
    });

    test('LLM node shows confirmation modal when saving (Task 3.3)', async ({ page }) => {
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
        await commandTextarea.fill('copilot -p "Refactor the auth module"');

        // Click Save - should show confirmation modal
        await page.locator('[data-testid="config-save-btn"]').click();

        // Confirmation modal should appear
        const confirmModal = page.locator('[data-testid="premium-confirm-modal"]');
        await expect(confirmModal).toBeVisible();
        await expect(confirmModal).toContainText('Premium Resource Confirmation');
        await expect(confirmModal).toContainText('consume tokens');

        // Confirm the premium execution
        await page.locator('[data-testid="confirm-premium-btn"]').click();

        // Modal and panel should close
        await expect(confirmModal).not.toBeVisible();
        await expect(configPanel).not.toBeVisible();
    });

    test('can switch between Shell and LLM node types', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        await page.locator('[data-testid="agent-node"]').first().click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // Initially shell type
        await expect(page.locator('[data-testid="node-type-shell"]')).toHaveClass(/border-green-500/);

        // Switch to LLM
        await page.locator('[data-testid="node-type-llm"]').click();
        await expect(page.locator('[data-testid="node-type-llm"]')).toHaveClass(/border-purple-500/);

        // Switch back to Shell
        await page.locator('[data-testid="node-type-shell"]').click();
        await expect(page.locator('[data-testid="node-type-shell"]')).toHaveClass(/border-green-500/);
    });

    test('Unconfigured node shows warning indicator', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });

        // Warning text should be visible for unconfigured node
        const warningText = page.locator('[data-testid="node-warning"]');
        await expect(warningText).toBeVisible();
        await expect(warningText).toContainText('Click to configure');
    });

    test('Cancel button closes panel without saving', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        await page.locator('[data-testid="agent-node"]').first().click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // Change the label
        const labelInput = page.locator('[data-testid="config-label-input"]');
        await labelInput.clear();
        await labelInput.fill('Should Not Save');

        // Click Cancel
        await page.locator('[data-testid="config-cancel-btn"]').click();

        // Verify panel closes and label was NOT saved
        await expect(configPanel).not.toBeVisible();
        await expect(page.locator('text=Should Not Save')).not.toBeVisible();
    });

    test('Token meter updates as command is typed for LLM nodes', async ({ page }) => {
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

        const tokenMeter = configPanel.locator('[data-testid="token-meter"]');
        await expect(tokenMeter).toBeVisible();

        // Type in the command
        const commandTextarea = page.locator('[data-testid="config-command-textarea"]');
        await commandTextarea.fill('This is a test prompt that should update the token count');

        // Verify token count updates
        await expect(tokenMeter).toContainText(/\d+ \/ 4,000/);
    });

    test('Close button (X) closes config panel', async ({ page }) => {
        const canvas = page.locator('.react-flow');
        const shellNodeDrag = page.locator('[data-testid="shell-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await shellNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        await page.waitForTimeout(500);
        await page.locator('[data-testid="agent-node"]').first().click();

        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible();

        // Click close button
        await page.locator('[data-testid="config-close-btn"]').click();

        await expect(configPanel).not.toBeVisible();
    });
});
