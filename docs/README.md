Building the Documentation Locally
When contributing to the documentation, it’s essential to build and preview it locally to ensure that your changes display correctly. This guide provides detailed steps on how to do that using Docker and Yarn.

Prerequisites
Before getting started, make sure you have the following installed:

Docker: Version >= 2.1.0.3
Yarn: Version >= 1.22.4
How to Build the Docs Locally
Follow these steps to build the documentation and preview it in your browser:

1. Start Docker
Ensure Docker is installed and running. You can check if Docker is active by running:
bash
Copy code
docker info
2. Switch to the Docs Directory
Open your terminal and navigate to the docs folder:
bash
Copy code
cd docs
3. Build and Preview the Site
Run the following command to build the docs and launch a local preview:
bash
Copy code
make docs
The site will be available at:
http://localhost:3002/docs/grafana/latest/
Any changes you make in the sources/ directory will automatically refresh the preview.
4. Use Local Static Assets (Optional)
If you have the grafana/website repo checked out in the same directory as the grafana repo, you can run:
bash
Copy code
make docs-local-static
This ensures that local assets like images are used instead of remote ones.
Editing Content Guidelines
Most content resides in the sources/ directory. Some pages, however, are auto-generated from TypeScript files. Please follow these instructions to edit content correctly.

Auto-Generated Content
Markdown Location:
docs/sources/panels-visualizations/query-transform-data/transform-data/index.md

TypeScript Source Files:

scripts/docs/generate-transformations.ts – General content.
public/app/features/transformers/docs/content.ts – Transformation-specific content.
Note: Always use reference-style links in content.ts to prevent issues with the UI.

Using Internal Links with relref
When linking to other pages within the documentation, use Hugo’s relref shortcode to maintain consistency.

Example:

hugo
Copy code
{{< relref "example.md" >}}
If Hugo reports an ambiguous link, provide the full path (e.g., folder/example.md).

Managing Redirects
When moving or removing pages, add redirects to avoid broken links.

If moving a page: Add an aliases entry in the page’s front matter to redirect the old URL to the new one.
If deleting a page: Add an aliases entry in the most relevant existing page.
Tip: If copying a page to create a new one, remove old aliases from the copy to avoid conflicting redirects.

Editing the Side Menu
The side menu is generated from the file structure. Use the weight parameter in the front matter to control the order of pages.

Change Menu Text: Use the menuTitle parameter if the menu text should differ from the page title.
Adding Images
Follow the guidelines for adding images, diagrams, and screenshots in the Image Guidelines.

Deploying Changes to grafana.com
When your pull request (PR) is merged, the changes in the docs/sources directory are automatically synced to the Grafana website via GitHub Actions.

PR to main branch:

Changes sync to content/docs/grafana/next.
Published at: https://grafana.com/docs/grafana/next/
PR to the release branch:

Changes sync to content/docs/grafana/latest.
Published at: https://grafana.com/docs/grafana/latest/
Once synced, the site will automatically update—no further action needed.

Summary
This README provides the necessary steps to build, edit, and deploy the documentation efficiently. It ensures that contributions are formatted correctly and aligned with Grafana’s standards. Follow these instructions to ensure smooth local builds and successful deployments.

Key Highlights of Changes
Clearer Instructions: Step-by-step build and preview process.
Structured Content: Auto-generated content guidelines separated for clarity.
Helpful Tips: Use of relref and aliases to manage links and redirects.
Deployment Clarity: Explained GitHub Action sync process for PRs.
