import { test, expect } from '@playwright/test';

test.describe('Flow Node Configuration', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to flow editor - create a new flow
        await page.goto('/flows/new');
        await page.waitForSelector('[data-testid="agent-node-drag"]', { timeout: 10000 });
    });

    test('Add node → Configure → Save → Node shows updated label', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            // Drag the agent node to the canvas center
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear and click on it
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Fill in the form fields
        const labelInput = page.locator('[data-testid="config-label-input"]');
        await labelInput.clear();
        await labelInput.fill('My Custom Agent');

        const roleSelect = page.locator('[data-testid="config-role-select"]');
        await roleSelect.selectOption('Architect');

        const promptTextarea = page.locator('[data-testid="config-prompt-textarea"]');
        await promptTextarea.fill('Analyze the code and create a detailed architecture plan');

        // Click Save
        const saveBtn = page.locator('[data-testid="config-save-btn"]');
        await saveBtn.click();

        // Verify config panel closes
        await expect(configPanel).not.toBeVisible();

        // Verify node shows updated label
        await expect(page.locator('text=My Custom Agent')).toBeVisible();
    });

    test('Unconfigured node shows warning indicator', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });

        // Verify warning indicator is visible (yellow dot for unconfigured)
        const warningIndicator = page.locator('[data-testid="node-unconfigured"]');
        await expect(warningIndicator).toBeVisible();

        // Verify warning text is visible
        const warningText = page.locator('[data-testid="node-warning"]');
        await expect(warningText).toBeVisible();
        await expect(warningText).toContainText('Not configured');
    });

    test('Cancel button closes panel without saving', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear and click on it
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Change the label
        const labelInput = page.locator('[data-testid="config-label-input"]');
        await labelInput.clear();
        await labelInput.fill('Should Not Save');

        // Click Cancel
        const cancelBtn = page.locator('[data-testid="config-cancel-btn"]');
        await cancelBtn.click();

        // Verify config panel closes
        await expect(configPanel).not.toBeVisible();

        // Verify the original label is still there (not the new one)
        await expect(page.locator('text=Should Not Save')).not.toBeVisible();
        // The agent node on the canvas should still have the original label
        await expect(page.locator('[data-testid="agent-node"]').locator('text=Agent Node')).toBeVisible();
    });

    test('Token meter updates as prompt is typed', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear and click on it
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Find token meter
        const tokenMeter = configPanel.locator('[data-testid="token-meter"]');
        await expect(tokenMeter).toBeVisible();

        // Type in the prompt
        const promptTextarea = page.locator('[data-testid="config-prompt-textarea"]');
        await promptTextarea.fill('This is a test prompt that should update the token count');

        // Verify token count updates (check that it shows some tokens)
        await expect(tokenMeter).toContainText(/\d+ \/ 4,000 tokens/);
    });

    test('Role dropdown has all required options', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear and click on it
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Check role dropdown options
        const roleSelect = page.locator('[data-testid="config-role-select"]');
        
        // Verify all required role options exist
        await expect(roleSelect.locator('option[value="Architect"]')).toHaveText('Planner / Architect');
        await expect(roleSelect.locator('option[value="Implementation"]')).toHaveText('Developer / Coder');
        await expect(roleSelect.locator('option[value="Test"]')).toHaveText('QA / Tester');
        await expect(roleSelect.locator('option[value="Optimizer"]')).toHaveText('Auditor / Optimizer');
    });

    test('Close button (X) closes config panel', async ({ page }) => {
        // Get the canvas container
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeVisible();

        // Drag an agent node to the canvas
        const agentNodeDrag = page.locator('[data-testid="agent-node-drag"]');
        const canvasBounds = await canvas.boundingBox();
        
        if (canvasBounds) {
            await agentNodeDrag.dragTo(canvas, {
                targetPosition: { x: canvasBounds.width / 2, y: canvasBounds.height / 2 }
            });
        }

        // Wait for the node to appear and click on it
        await page.waitForTimeout(500);
        const agentNode = page.locator('[data-testid="agent-node"]').first();
        await expect(agentNode).toBeVisible({ timeout: 5000 });
        await agentNode.click();

        // Verify config panel opens
        const configPanel = page.locator('[data-testid="node-config-panel"]');
        await expect(configPanel).toBeVisible({ timeout: 5000 });

        // Click close button (X)
        const closeBtn = page.locator('[data-testid="config-close-btn"]');
        await closeBtn.click();

        // Verify config panel closes
        await expect(configPanel).not.toBeVisible();
    });
});
