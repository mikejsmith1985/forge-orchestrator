# Forge Orchestrator - User Guide

**Version:** 1.1.0  
**Last Updated:** December 2024

Welcome to Forge Orchestrator! This guide will walk you through every feature in plain, easy-to-understand language.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Welcome Screen](#welcome-screen)
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

Forge Orchestrator is your AI workflow command center. Think of it like a control panel for managing AI interactions, tracking how much you're using, and creating automated workflows.

### How to Launch

1. Double-click the `forge-orchestrator` application file
2. A browser window will automatically open
3. The app runs at `http://localhost:8080` by default

---

## Welcome Screen

### What It Does

- Shows you an overview of all the features when you first open the app
- Appears automatically the first time you use the app
- Also shows up after you update to a new version

### How to Use It

1. **Read the features** - The welcome screen shows 6 main things you can do:
   - **Architect** - Write down your ideas and see token counts
   - **Ledger** - Track how many tokens you've used
   - **Flows** - Create visual workflows connecting AI agents
   - **Commands** - Save frequently-used commands
   - **API Keys** - Set up your AI provider credentials
   - **Optimization** - Get suggestions to save tokens

2. **Dismiss it** - You can close the welcome screen by:
   - Clicking the **"Get Started"** button
   - Pressing the **Escape** key
   - Pressing **Enter** or **Space**
   - Clicking the **X** in the corner

3. **Won't annoy you** - Once you close it, it won't show again until you update to a new version

---

## Navigation (Sidebar)

### What It Does

The sidebar is your main menu on the left side of the screen. It lets you jump between different sections of the app.

### The Menu Items

| Icon | Name | What It Does |
|------|------|--------------|
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

1. **Navigate there** - Click "Architect" in the sidebar (or go to `/architect`)

2. **Type your ideas** - Click in the big text box and start typing whatever you want:
   - Project requirements
   - Task descriptions
   - Questions for AI
   - Anything!

3. **Watch the token counter** - As you type, the bottom of the screen shows:
   - **Token count** - How many tokens your text uses
   - **Method** - How the tokens were counted (tiktoken = accurate, heuristic = estimate)
   - **Provider** - Which AI provider the count is for (like OpenAI)

### Tips

- Tokens are like "word pieces" that AI uses to understand text
- More tokens = more cost when using AI
- The counter updates automatically as you type (with a small delay)

---

## Dashboard / Ledger

### What It Does

The Ledger is your token usage history. It shows every interaction you've had and how many tokens were used.

### How to Use It

1. **Navigate there** - Click "Dashboard" in the sidebar (or go to `/ledger`)

2. **View your history** - You'll see a list of past interactions showing:
   - When it happened
   - What was sent/received
   - How many tokens were used
   - Which AI provider was used

3. **Look for patterns** - Use this to:
   - See which tasks use the most tokens
   - Track your usage over time
   - Find opportunities to save tokens

### Optimization Suggestions

- The system analyzes your usage patterns
- It suggests ways to reduce token usage
- Click "Apply" on any suggestion to implement it

---

## Commands

### What It Does

Commands are like shortcuts. You save commands you use often, then run them with one click.

### How to Use It

1. **Navigate there** - Click "Commands" in the sidebar (or go to `/commands`)

2. **View your commands** - See all your saved commands as cards

3. **Create a new command:**
   - Click the **"Add Command"** or **"+"** button
   - Give it a name (like "List Files")
   - Enter the actual command (like `ls -la`)
   - Optionally assign a keyboard shortcut
   - Click **Save**

4. **Run a command:**
   - Click on any command card
   - Or use the keyboard shortcut you assigned

5. **Edit a command:**
   - Click the edit/pencil icon on any command card
   - Make your changes
   - Click Save

6. **Delete a command:**
   - Click the delete/trash icon on any command card
   - Confirm you want to delete it

### Keyboard Shortcuts

- First 10 commands get shortcuts: `Ctrl+Shift+1` through `Ctrl+Shift+0`
- Additional commands get letter shortcuts: `Ctrl+Shift+A`, `Ctrl+Shift+B`, etc.

---

## Flows

### What It Does

Flows let you create visual workflows that connect different AI agents together. Think of it like drawing a flowchart where each box does something.

### How to Use It

1. **Navigate there** - Click "Flows" in the sidebar (or go to `/flows`)

2. **View your flows** - See a list of all your saved workflows

3. **Create a new flow:**
   - Click **"New Flow"** or **"+"** button
   - You'll enter the flow editor

4. **In the Flow Editor:**
   - **Add nodes** - Click and drag to add new steps
   - **Connect nodes** - Draw lines between nodes to show the order
   - **Configure nodes** - Click on a node to set what it does
   - **Name your flow** - Give it a meaningful name at the top

5. **Node Types:**
   - **Input** - Where data comes in
   - **AI Agent** - Processes with an AI model
   - **Output** - Where results go
   - **Condition** - Makes decisions based on data

6. **Save your flow:**
   - Click the **Save** button
   - Your flow appears in the flows list

7. **Run a flow:**
   - Click the **Run** or **Execute** button
   - Watch as each node processes in order
   - See the results when it finishes

8. **Delete a flow:**
   - Click the delete icon on any flow in the list
   - Confirm deletion

---

## Settings (API Keys)

### What It Does

This is where you set up your AI provider credentials. You need API keys to use AI services like Anthropic (Claude), OpenAI (ChatGPT), or Google.

### How to Use It

1. **Navigate there** - Click "Settings" in the sidebar (or go to `/settings`)

2. **See available providers:**
   - **Anthropic** - For Claude AI
   - **OpenAI** - For ChatGPT/GPT-4
   - **Google** - For Gemini

3. **Check the status:**
   - ✅ **Green checkmark** = Key is configured
   - ❌ **Red X** = Key is not set up

4. **Add or update a key:**
   - Find the provider you want (like "Anthropic")
   - Paste your API key in the password field
   - Click the **"Save Key"** button
   - You'll see a green success message

5. **If it fails:**
   - You'll see a red error message
   - Check that your key is correct
   - Make sure you copied the whole key

### Where to Get API Keys

- **Anthropic:** https://console.anthropic.com/
- **OpenAI:** https://platform.openai.com/api-keys
- **Google:** https://makersuite.google.com/app/apikey

### Security Note

- Your keys are stored securely in your system's keyring
- They are never displayed after you save them
- They never leave your computer (except to authenticate with the AI provider)

---

## Sending Feedback

### What It Does

Found a bug? Have an idea? The feedback feature lets you report issues directly to the developers with screenshots and logs included.

### How to Use It

1. **Open the feedback form:**
   - Click **"Send Feedback"** in the sidebar
   - A popup window appears

2. **First-time setup (one time only):**
   - You need a GitHub Personal Access Token
   - Click the **"Generate Token"** link
   - GitHub opens - follow the prompts to create a token
   - Copy the token (starts with `ghp_`)
   - Paste it in the token field
   - Click **"Save Settings"**

3. **Write your feedback:**
   - Describe the issue or suggestion in the text box
   - Be specific! Include:
     - What you were doing
     - What you expected to happen
     - What actually happened

4. **Add screenshots (optional but helpful):**
   - Click **"Capture Screen"**
   - The app takes a screenshot of what's behind the popup
   - You can capture multiple screenshots
   - Click the **X** on any screenshot to remove it

5. **Submit your feedback:**
   - Click **"Submit Feedback"**
   - Wait for the upload to complete
   - You'll see a success message with a link to your issue

6. **Cancel anytime:**
   - Click **"Cancel"** or the **X** to close without submitting
   - Your screenshots are saved if you reopen the form

### What Gets Sent

- Your description
- Your screenshots (uploaded to GitHub)
- Your browser info (to help debug)
- Recent application logs (to help find the bug)

---

## Updates

### What It Does

The app can check for new versions and help you update.

### How Updates Work

1. **Automatic checking:**
   - The app checks for updates when you open it
   - It also checks every 30 minutes

2. **When an update is available:**
   - A purple **"Update Available"** button appears in the sidebar
   - A small notification may pop up

3. **View update details:**
   - Click the purple button or the notification
   - A popup shows:
     - Current version
     - New version
     - What's changed (release notes)

4. **Download the update:**
   - Click the **"Download"** button
   - The new version downloads
   - Follow any prompts to install

5. **Dismiss for later:**
   - Click **"Later"** or close the popup
   - The reminder won't show again for 24 hours
   - The purple button stays visible

### Checking Manually

- Click on the version number in the sidebar (like "v1.1.0")
- This opens the update popup even if no update is available

---

## Keyboard Shortcuts Summary

| Shortcut | What It Does |
|----------|--------------|
| `Ctrl+Shift+1-9` | Run commands 1-9 |
| `Ctrl+Shift+0` | Run command 10 |
| `Ctrl+Shift+A-Z` | Run commands 11+ |
| `Escape` | Close popups/modals |
| `Enter` / `Space` | Dismiss welcome screen |

---

## Troubleshooting

### App Won't Start

- Make sure no other program is using port 8080
- Try running as administrator
- Check if your antivirus is blocking it

### API Key Won't Save

- Make sure you copied the entire key
- Check your internet connection
- Try a different browser

### Feedback Won't Submit

- Verify your GitHub token is correct
- Make sure the token has `public_repo` permission
- Check your internet connection

### Updates Not Showing

- Click the version number to manually check
- Make sure you're connected to the internet
- Try restarting the app

---

## Getting Help

- **Send Feedback** - Use the built-in feedback feature
- **GitHub Issues** - Visit the project repository
- **Documentation** - Check the Project Documentation folder

---

*Happy orchestrating! 🚀*
