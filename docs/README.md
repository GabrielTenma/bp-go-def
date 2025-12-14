# GitHub Pages Setup Instructions

## Files Created

The GitHub Pages site is located in the `docs` folder:

```
docs/
├── index.md                    # Main landing page (English)
├── _config.yml                 # Jekyll configuration
└── assets/
    └── css/
        └── style.scss          # Custom styling
```

## How to Enable GitHub Pages

### 1. Push Changes to GitHub

```powershell
# Navigate to repository root
cd "c:\Users\Tzadkiel\Documents\Development\Git\Own Repos\GabrielTenma\bp-go-def"

# Add all files
git add docs/

# Commit changes
git commit -m "Add GitHub Pages site"

# Push to GitHub
git push origin main
```

### 2. Enable GitHub Pages in Repository Settings

1. Go to your repository on GitHub: https://github.com/GabrielTenma/bp-go-def
2. Click on **Settings** tab
3. Scroll down to **Pages** section (in the left sidebar under "Code and automation")
4. Under **Source**:
   - Select **Deploy from a branch**
   - Branch: `main` (or your default branch)
   - Folder: `/docs`
   - Click **Save**

### 3. Wait for Deployment

GitHub will automatically build and deploy your site. This usually takes 1-3 minutes.

You can check the deployment status:
- Go to **Actions** tab in your repository
- Look for "pages build and deployment" workflow
- Once it shows a green checkmark, your site is live!

### 4. Access Your Site

Your GitHub Pages site will be available at:
**https://gabrieltenma.github.io/bp-go-def/**

---

## Features Included

### Main Landing Page (`index.md`)
- Beautiful feature showcase in English
- Key advantages highlighted:
  - Modular architecture
  - Beautiful monitoring dashboard
  - Terminal Interface (TUI)
  - Complete infrastructure integrations
  - Standardized API patterns
  - Premium developer experience
  - Security-first approach
- Quick start guide
- Screenshots
- Links to developer documentation in `/docs` folder

### Custom Styling
Modern, professional design with:
- Go-themed color palette (#00ADD8)
- Responsive layout for all devices
- Enhanced typography
- Beautiful code blocks
- Smooth animations
- Card-based layouts

---

## Folder Structure

The project has two separate documentation folders:

### `docs/` - GitHub Pages (Public Site)
- Public-facing landing page
- Feature showcase
- Getting started guide
- Links to developer documentation

### `docs_wiki/` - Developer Documentation
- Integration guides
- Architecture diagrams
- API documentation
- Request/Response structure
- Internal development notes

---

## Testing Locally (Optional)

If you want to preview the site locally before pushing:

### Install Jekyll (One-time setup)
```powershell
# Install Ruby (if not already installed)
# Download from: https://rubyinstaller.org/

# Install Bundler
gem install bundler

# Navigate to docs folder
cd docs

# Create Gemfile
@"
source 'https://rubygems.org'
gem 'github-pages', group: :jekyll_plugins
gem 'webrick'
"@ | Out-File -FilePath Gemfile -Encoding utf8

# Install dependencies
bundle install
```

### Run Local Server
```powershell
cd docs
bundle exec jekyll serve
```

Then visit: http://localhost:4000/bp-go-def/

Press `Ctrl+C` to stop the server.

---

## Notes

- The theme used is **jekyll-theme-cayman** (GitHub's official theme)
- All content is in English
- Developer documentation remains in the `docs_wiki/` folder
- GitHub Pages content is in the `docs/` folder
- The site is fully responsive (mobile, tablet, desktop)

---

## Next Steps

1. Push the changes to GitHub
2. Enable GitHub Pages in repository settings (select `/docs` folder)
3. Wait for deployment (1-3 minutes)
4. Visit your site at https://gabrieltenma.github.io/bp-go-def/
5. Share the link!

---

## Updating Content

To update the site in the future:
1. Edit files in the `docs/` folder
2. Commit and push changes
3. GitHub Pages will automatically rebuild and deploy

---

## Customization

You can customize further by editing:
- `docs/_config.yml` - Site metadata and settings
- `docs/assets/css/style.scss` - Styling and colors
- `docs/index.md` - Main page content
