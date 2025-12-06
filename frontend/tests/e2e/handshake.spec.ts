import { test, expect } from '@playwright/test';
import { exec } from 'child_process';
import { promisify } from 'util';
import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';

const execAsync = promisify(exec);

// Get __dirname equivalent in ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Handshake Automation Tests
 * 
 * Tests the automated handshake generation and synchronization system
 */

test.describe.serial('Handshake Automation', () => {
    const projectRoot = path.join(__dirname, '../../..');  // frontend/tests/e2e -> root
    const scriptsDir = path.join(projectRoot, 'scripts');
    const handshakeFile = path.join(projectRoot, 'FORGE_HANDSHAKE.md');

    test('generate-handshake script exists and is executable', async () => {
        const scriptPath = path.join(scriptsDir, 'generate-handshake.sh');
        
        // Check file exists
        expect(fs.existsSync(scriptPath)).toBeTruthy();
        
        // Check it's executable
        const stats = fs.statSync(scriptPath);
        const isExecutable = (stats.mode & fs.constants.X_OK) !== 0;
        expect(isExecutable).toBeTruthy();
    });

    test('validate-handshake script exists and is executable', async () => {
        const scriptPath = path.join(scriptsDir, 'validate-handshake.sh');
        
        expect(fs.existsSync(scriptPath)).toBeTruthy();
        
        const stats = fs.statSync(scriptPath);
        const isExecutable = (stats.mode & fs.constants.X_OK) !== 0;
        expect(isExecutable).toBeTruthy();
    });

    test('watch-releases script exists and is executable', async () => {
        const scriptPath = path.join(scriptsDir, 'watch-releases.sh');
        
        expect(fs.existsSync(scriptPath)).toBeTruthy();
        
        const stats = fs.statSync(scriptPath);
        const isExecutable = (stats.mode & fs.constants.X_OK) !== 0;
        expect(isExecutable).toBeTruthy();
    });

    test('sync-terminal-handshake script exists and is executable', async () => {
        const scriptPath = path.join(projectRoot, 'sync-terminal-handshake.sh');
        
        expect(fs.existsSync(scriptPath)).toBeTruthy();
        
        const stats = fs.statSync(scriptPath);
        const isExecutable = (stats.mode & fs.constants.X_OK) !== 0;
        expect(isExecutable).toBeTruthy();
    });

    test('can generate handshake document', async () => {
        // Run generation script
        const { stdout, stderr } = await execAsync(
            './scripts/generate-handshake.sh',
            { cwd: projectRoot }
        );
        
        // Should complete without error
        expect(stdout).toContain('Generating Forge Orchestrator Handshake');
        expect(stdout).toContain('Backend Version');
        expect(stdout).toContain('✅');
        
        // Handshake file should exist
        expect(fs.existsSync(handshakeFile)).toBeTruthy();
    });

    test('generated handshake has required content', async () => {
        // Generate handshake first
        await execAsync('./scripts/generate-handshake.sh', { cwd: projectRoot });
        
        // Read handshake file
        const content = fs.readFileSync(handshakeFile, 'utf-8');
        
        // Check required sections
        expect(content).toContain('Forge Orchestrator');
        expect(content).toContain('Core Architecture');
        expect(content).toContain('API Endpoints');
        expect(content).toContain('UI Components');
        expect(content).toContain('Feature Requirements');
        expect(content).toContain('Configuration');
        expect(content).toContain('Testing');
        expect(content).toContain('Release Process');
        
        // Check version format
        expect(content).toMatch(/\*\*Version\*\*:.*\d+\.\d+/);
        
        // Check timestamp format
        expect(content).toMatch(/\*\*Last Updated\*\*:.*\d{4}-\d{2}-\d{2}/);
        
        // Check has component count
        expect(content).toMatch(/\d+ React Components/);
        
        // Check has API endpoints
        expect(content).toContain('`/');
    });

    test('can validate handshake document', async () => {
        // Generate first
        await execAsync('./scripts/generate-handshake.sh', { cwd: projectRoot });
        
        // Run validation
        const { stdout } = await execAsync(
            './scripts/validate-handshake.sh',
            { cwd: projectRoot }
        );
        
        // Should pass validation
        expect(stdout).toContain('Validating handshake');
        expect(stdout).toContain('✅ Validation passed');
        expect(stdout).toContain('feature checkboxes');
    });

    test('Makefile has handshake targets', async () => {
        const makefilePath = path.join(projectRoot, 'Makefile');
        
        expect(fs.existsSync(makefilePath)).toBeTruthy();
        
        const content = fs.readFileSync(makefilePath, 'utf-8');
        
        // Check for handshake targets
        expect(content).toContain('handshake:');
        expect(content).toContain('validate-handshake:');
        expect(content).toContain('sync-terminal:');
        expect(content).toContain('watch-terminal:');
    });

    test('documentation files exist', async () => {
        // Check for documentation
        const docsPath = path.join(projectRoot, 'docs/RELEASE_AUTOMATION.md');
        expect(fs.existsSync(docsPath)).toBeTruthy();
        
        const quickRef = path.join(projectRoot, 'handoffs/HANDSHAKE_QUICK_REF.md');
        expect(fs.existsSync(quickRef)).toBeTruthy();
        
        // Check content
        const docsContent = fs.readFileSync(docsPath, 'utf-8');
        expect(docsContent).toContain('Release Automation Guide');
        expect(docsContent).toContain('Handshake Flow');
        
        const quickRefContent = fs.readFileSync(quickRef, 'utf-8');
        expect(quickRefContent).toContain('Quick Reference');
        expect(quickRefContent).toContain('make handshake');
    });

    test('GitHub workflow includes handshake generation', async () => {
        const workflowPath = path.join(projectRoot, '.github/workflows/release.yml');
        
        expect(fs.existsSync(workflowPath)).toBeTruthy();
        
        const content = fs.readFileSync(workflowPath, 'utf-8');
        
        // Check for handshake steps
        expect(content).toContain('Generate Handshake Document');
        expect(content).toContain('Validate Handshake Document');
        expect(content).toContain('generate-handshake.sh');
        expect(content).toContain('validate-handshake.sh');
        expect(content).toContain('FORGE_HANDSHAKE.md');
    });

    test('handshake includes orchestrator-specific features', async () => {
        await execAsync('./scripts/generate-handshake.sh', { cwd: projectRoot });
        
        const content = fs.readFileSync(handshakeFile, 'utf-8');
        
        // Should mention orchestrator features
        expect(content).toContain('Architect');
        expect(content).toContain('Workflow');
        expect(content).toContain('Flow');
        expect(content).toContain('Ledger');
        expect(content).toContain('Command Deck');
    });

    test('handshake includes terminal features', async () => {
        await execAsync('./scripts/generate-handshake.sh', { cwd: projectRoot });
        
        const content = fs.readFileSync(handshakeFile, 'utf-8');
        
        // Should mention terminal features from Terminal project
        expect(content).toContain('Terminal');
        expect(content).toContain('WebSocket');
        expect(content).toContain('PTY');
        expect(content).toContain('WSL');
        expect(content).toContain('Auto-respond');
        expect(content).toContain('Auto-reconnection');
    });
});
