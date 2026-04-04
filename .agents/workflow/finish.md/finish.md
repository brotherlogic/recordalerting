---
description: Automatically commit and push changes to a new branch
---

Whenever you make modifications to files in this repository during a conversation, you should automatically commit and push your changes to a new branch on GitHub, unless the user explicitly tells you not to.

Follow these steps via terminal commands to push the changes:

// turbo-all
1. Verify the changes using `git status` and `git diff`
2. If we're on the main branch, then create and switch to a new branch with a short, descriptive name: `git checkout -b <descriptive-branch-name>`; otherwise use the existing branch going forwards
3. Stage the modified files: `git add <files>`
4. Commit the changes with a clear message: `git commit -m "<Clear description of the changes>"`
5. Push the new branch to origin: `git push -u origin <descriptive-branch-name>`
6. Inform the user that you have pushed the branch and provide the GitHub Pull Request creation link that appears in the git push output.