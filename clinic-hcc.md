# HCC Clinic

## mattpocock skills

- <https://github.com/mattpocock/skills>

```md

/grill-with-docs 

---
npx skills add mattpocock/skills

/setup-matt-pocock-skills 

/improve-codebase-architecture 

```

## superpowers

- <https://github.com/obra/superpowers/blob/main/docs/README.opencode.md>

```md
use skill tool to load brainstorming
use skill tool to load writing-plan

```

---

## playwright

```md
need a demo simulation of printing receipt flow by `playwright-cli --head` 
- open http://localhost:5173, 
- navigate to receipt printing page, 
- input data fields by faker data for demo, 
- input 2-3 receipt items with different types, 
- click print button to call browser print fuction
- click save button to save the receipt data
- all actions pause by 3 seconds

load skill /playwright-cli
condense all test\demo-receipt-*.spec.ts, tests/DEMO-QUICKSTART.md, fix all problem, errors, handle invalid inputs iteratively, use fast simulation without actions pause to spped up your tests, make the demo simulation of printing receipt flow compelte until you reach the browser print dialogue. 

condense all test\demo-receipt-*.spec.ts into single simulation run, do not over engineering

# ---
npm install -g @playwright/cli@latest
npm init playwright@latest

npm install -D playwright  
npx playwright install chromium
playwright-cli install --skills

```

## skill

```md
# List all installed skills (project and global)
npx skills list

# List only global skills
npx skills ls -g

# Filter by specific agents
npx skills ls -a claude-code -a opencode 

```

## install

```md
# git
git remote add origin https://github.com/hccyip-clinic/clinic-hugo.git
git remote set-url origin https://github.com/hccyip-clinic/clinic-hugo.git

git branch -M main
git push -u origin main

git config --global user.name "Your Name"


# global

npm install -g @github/copilot
npm install -g opencode-ai@latest --allow-scripts=opencode-ai


npm install -D vitest
npm install -D tsx
npm install --save-dev fake-indexeddb

# npm
npm outdated
npm update

npm install @faker-js/faker@latest
npm install @types/node@latest
npm install vite-plugin-svgr@latest

```
