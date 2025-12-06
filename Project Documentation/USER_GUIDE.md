# Forge Orchestrator - User Guide

**Version:** 1.1.1  
**Last Updated:** December 2024

Welcome to Forge Orchestrator! This guide will walk you through every feature in plain, easy-to-understand language.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Terminal (Home View)](#terminal-home-view)
3. [Navigation (Sidebar)](#navigation-sidebar)
4. [Architect View](#architect-view)
5. [Dashboard / Ledger](#dashboard--ledger)
6. [Commands](#commands)
7. [Flows](#flows)
8. [Settings (API Keys)](#settings-api-keys)
9. [Sending Feedback](#sending-feedback)
10. [Updates](#updates)

---

## Getting Started

### What is Forge Orchestrator?

Forge Orchestrator is your **terminal-first AI workflow command center**. At its heart is a real terminal (PTY) where you can type commands and see actual shell output. Around that terminal, you get:

- Token usage tracking for AI providers
- Visual workflow automation
- Secure API key management
- Command shortcuts

### How to Launch

1. Run the `forge-orchestrator` binary
2. Open your browser to `http://localhost:8080` (or the URL shown in the terminal)
3. You'll land on the **Terminal** view - your primary workspace

---

## Terminal (Home View)

### What It Does

The Terminal is Forge Orchestrator's **main feature**. It's a real shell running on your machine, streamed to your browser. Everything you type executes in your actual environment.

### The Interface

| Element | Description |
|---------|-------------|
| 🔴🟡🟢 Traffic Lights | Decorative header (macOS style) |
| "Terminal" Label | Shows you're in the terminal view |
| Green/Red Dot | Connection status (green = connected) |
| **Prompt Watcher** Toggle | Auto-respond to y/n prompts |
| Terminal Area | Where you type and see output |

### How to Use It

1. **Type commands** - Just like any terminal
   ```bash
   ls -la
   git status
   npm install
   ```

2. **See output** - Real shell output appears in real-time

3. **Connection indicator** - The small dot shows:
   - 🟢 Green = Connected and ready
   - 🔴 Red = Disconnected (will auto-reconnect)

### Prompt Watcher

The **Prompt Watcher** button in the top-right corner enables automatic responses to confirmation prompts:

- **Off (default)**: You manually respond to "Continue? [y/n]" prompts
- **On (blue)**: Automatically sends "y" when it detects prompts like:
  - `[y/n]`
  - `[Y/n]`
  - `Continue?`
  - `Are you sure?`

**Use cases:**
- Unattended installations (`npm install`, `apt-get`)
- Batch operations that need confirmation
- Automated testing scripts

⚠️ **Use with caution** - This will auto-confirm potentially destructive operations!

---

## Navigation (Sidebar)

### What It Does

The sidebar is your main menu on the left side of the screen. It lets you jump between different sections of the app.

### The Menu Items

| Icon | Name | What It Does |
|------|------|--------------|
| 🖥️ | **Terminal** | Your main shell (default view) |
| 🧠 | Architect | Write ideas and count tokens |
| 📊 | Dashboard | View your token usage history |
| ⚡ | Commands | Manage saved commands |
| 🔀 | Flows | Create AI workflows |
| ⚙️ | Settings | Manage API keys |

### Other Buttons in the Sidebar

- **Send Feedback** - Report bugs or suggest features
- **Update Available** - Shows when a new version is ready (purple button)
- **Version Number** - Click to check for updates
- **Green Dot** - Shows the app is running ("Online")

### Mobile Users

- On phones or small screens, the sidebar hides automatically
- Tap the **hamburger menu** (☰) in the top-left to open it
- Tap anywhere outside the menu to close it

---

## Architect View

### What It Does

The Architect is where you "brain dump" your ideas. It's like a notepad that also counts how many AI tokens your text would use.

### How to Use It

1. **Navigate there** - Click "Architect" in the sidebar

2. **Type your ideas** - Click in the big text box and start typing:
   - Project requirements
   - Task descriptions
   - Questions for AI

3. **Watch the token counter** - The bottom shows:
   - 🪙 **Token count** - How many tokens your text uses
   - **Method** - tiktoken (accurate) or heuristic (estimate)
   - **Provider** - Which AI provider the count is for

### Dynamic Budget Meter

The token meter now shows your budget in the correct currency:

- **🪙 tokens** - For providers like OpenAI/Anthropic
- **💬 prompts** - For per-request pricing models

---

## Dashboard / Ledger

### What It Does

The Ledger is your token usage history. It shows every AI interaction and how much it cost.

### How to Use It

1. **Navigate there** - Click "Dashboard" in the sidebar

2. **View your history** - Each row shows:
   - Timestamp
   - Flow ID
   - Model used
   - **Usage** (dynamic based on billing type):
     - Token-based: Shows ↓ input / ↑ output tokens
     - Prompt-based: Shows number of prompts
   - Latency
   - Cost in USD
   - Status

3. **Cost Unit Badge** - Each entry shows its billing type:
   - 🪙 **Per-Token** - Traditional token billing
   - 💬 **Per-Prompt** - Per-request billing

### Optimization Suggestions

- The system analyzes your usage patterns
- It suggests ways to reduce costs
- Click "Apply" on any suggestion to implement it

---

## Commands

### What It Does

Commands are like shortcuts. Save commands you use often, then run them with one click.

### How to Use It

1. **Navigate there** - Click "Commands" in the sidebar

2. **View your commands** - See all saved commands as cards

3. **Create a new command:**
   - Click the **"Add Command"** button
   - Give it a name (like "Run Tests")
   - Enter the command (`npm test`)
   - Click **Save**

4. **Run a command:**
   - Click on any command card
   - The command is sent to the Terminal

5. **Edit/Delete:**
   - Use the edit/delete icons on each card

### Keyboard Shortcuts

- First 10 commands: `Ctrl+Shift+1` through `Ctrl+Shift+0`
- Additional commands: `Ctrl+Shift+A`, `Ctrl+Shift+B`, etc.

---

## Flows

### What It Does

Flows let you create visual workflows that chain together shell commands and AI operations.

### The New Node Types (V2.1)

Flows now have **two distinct node types**:

| Node Type | Badge | Token Cost | Use For |
|-----------|-------|------------|---------|
| **Shell Command** | ⚡ Zero-Token | Free | Local scripts, git, npm |
| **LLM Prompt** | 💎 Premium | Uses budget | AI-powered tasks |

### How to Use It

1. **Navigate there** - Click "Flows" in the sidebar

2. **Create a new flow** - Click **"Create New Flow"**

3. **Drag nodes** from the sidebar:
   - **Shell Command** (green) - For local execution
   - **LLM Prompt** (purple) - For AI operations
   - **Input/Output** nodes - For data flow

4. **Configure a node** - Click on it to open the config panel:
   - **Label** - Give it a meaningful name
   - **Node Type** - Switch between Shell and LLM
   - **Command/Prompt** - The actual command to run

5. **Connect nodes** - Drag from output handles to input handles

6. **Save and Execute** - Use the toolbar buttons

### Premium Confirmation Modal

When you configure an **LLM Prompt node**:

1. Fill in your prompt
2. Click **Save**
3. A **confirmation modal** appears showing:
   - Estimated token cost
   - Warning that this uses premium budget
4. Click **Confirm & Save** to proceed

This prevents accidental token spending!

### Example Flow

```
[Input] → [Shell: git status] → [LLM: Summarize changes] → [Shell: create-pr.sh] → [Output]
```

---

## Settings (API Keys)

### What It Does

This is where you set up your AI provider credentials securely.

### Security Assurance

At the top of the Settings page, you'll see a **security notice**:

> 🔐 **Secure Storage**
> Your API keys are **encrypted** and stored in your operating system's native keyring (macOS Keychain, Windows Credential Manager, or Linux Secret Service). Keys are **never** exposed to the browser or stored in plain text.

### How to Use It

1. **Navigate there** - Click "Settings" in the sidebar

2. **See available providers:**
   - **Anthropic** - For Claude AI
   - **OpenAI** - For ChatGPT/GPT-4

3. **Check the status:**
   - ✅ **Configured** = Key is set
   - ❌ **Not Configured** = Key needed

4. **Add or update a key:**
   - Paste your API key in the password field
   - Click **"Save Key"**
   - See green success message

### Where to Get API Keys

- **Anthropic:** https://console.anthropic.com/
- **OpenAI:** https://platform.openai.com/api-keys

---

## Sending Feedback

### What It Does

Report bugs or suggest features directly to the developers. The app will create a GitHub issue with your feedback, screenshots, and diagnostic logs to help us fix problems faster.

### Initial Setup (One-Time)

The first time you use feedback, you'll need to set up a GitHub Personal Access Token (PAT):

1. **Why you need it:**
   - Allows the app to create issues on your behalf
   - Uploads screenshots to help developers see the problem
   - Completely safe - only has permission to create issues

2. **Click "Generate Token on GitHub":**
   - This opens GitHub with a **prefilled form**
   - Token name: "Forge Orchestrator Feedback"
   - Permission scope: `public_repo` (create issues)
   - Click **"Generate token"** on GitHub

3. **Copy and paste:**
   - GitHub shows your new token (starts with `ghp_`)
   - Copy it (you won't see it again!)
   - Paste into the Forge Orchestrator input field
   - Click **"Save Settings"**

### Sending Feedback

Once setup is complete:

1. **Open the form** - Click "Send Feedback" in the sidebar

2. **Describe the issue:**
   - Be specific about what happened
   - What you were doing when it occurred
   - What you expected vs. what actually happened

3. **Add screenshots (optional but helpful):**
   - Click **"Capture Screen"** to take a screenshot
   - The screenshot is automatically added to your feedback
   - You can add multiple screenshots
   - Click the trash icon to remove unwanted screenshots

4. **Submit:**
   - Click **"Submit Feedback"**
   - The app uploads screenshots and creates a GitHub issue
   - You'll see the issue number and a link when done
   - The modal closes automatically

### What Gets Included

Your feedback submission includes:

- Your description
- Screenshots you captured
- Browser information (user agent)
- Timestamp
- Application logs (helps with debugging)

### Troubleshooting

**"Invalid GitHub token" error:**
- Your token may have expired
- Click "Update Settings" and generate a new token
- Make sure you copied the entire token

**"Token lacks permissions" error:**
- The token needs `public_repo` scope
- Use the "Generate Token on GitHub" link to create a proper token
- Don't manually create tokens - use the prefilled link

**Screenshot upload failed:**
- Your token may be invalid
- Check your internet connection
- The issue will still be created without the screenshot

---

## Updates

### How Updates Work

1. **Automatic checking** - Checks on startup and every 30 minutes

2. **When available:**
   - Purple **"Update Available"** button appears
   - Click to see what's new

3. **Download:**
   - Click **Download** in the update modal
   - Follow installation prompts

4. **Check manually:**
   - Click the version number in the sidebar

---

## Keyboard Shortcuts

| Shortcut | What It Does |
|----------|--------------|
| `Ctrl+Shift+1-9` | Run commands 1-9 |
| `Ctrl+Shift+0` | Run command 10 |
| `Ctrl+Shift+A-Z` | Run commands 11+ |
| `Escape` | Close modals |

---

## Troubleshooting

### Terminal Shows "Disconnected"

- The backend may not be running
- Check that `forge-orchestrator` is running in your terminal
- Refresh the page to reconnect

### API Key Won't Save

- Make sure you copied the entire key
- Check your internet connection
- Verify the key starts with the expected prefix (e.g., `sk-` for OpenAI)

### Flow Nodes Not Saving

- For LLM nodes, you must confirm the premium modal
- Make sure you entered a command/prompt

### Updates Not Showing

- Click the version number to manually check
- Check your internet connection

---

## Getting Help

- **Send Feedback** - Use the built-in feature
- **GitHub Issues** - Visit the project repository
- **Documentation** - Check the `docs/` folder

---

*Happy orchestrating! 🚀*
