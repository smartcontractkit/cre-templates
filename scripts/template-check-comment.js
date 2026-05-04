/**
 * Used by .github/workflows/template-check-comment.yml.
 * Reads /tmp/check-results from the prior workflow's artifact and posts or
 * removes a PR comment via the GitHub API.
 */
module.exports = async ({ github, context }) => {
  const fs = require('fs');

  const read = (filename) => {
    try {
      return fs.readFileSync(`/tmp/check-results/${filename}`, 'utf8').trim();
    } catch {
      return '';
    }
  };

  // Validate PR number — must be a positive integer.
  const prNumber = parseInt(read('pr-number.txt'), 10);
  if (!Number.isInteger(prNumber) || prNumber <= 0) {
    console.log('Invalid or missing PR number in artifact; skipping comment.');
    return;
  }

  const exitCode = read('exit-code.txt');

  const marker = '<!-- template-check-comment -->';
  const runUrl = `${context.serverUrl}/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.payload.workflow_run.id}`;

  const { data: comments } = await github.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: prNumber,
  });
  const existing = comments.find((c) => c.body.includes(marker));

  // ---------------------------------------------------------------
  // All templates passed (or no templates changed) — clean up.
  // ---------------------------------------------------------------
  if (exitCode === '0') {
    if (existing) {
      await github.rest.issues.deleteComment({
        owner: context.repo.owner,
        repo: context.repo.repo,
        comment_id: existing.id,
      });
      console.log('Deleted stale failure comment.');
    }
    return;
  }

  // ---------------------------------------------------------------
  // Build the failure comment.
  // ---------------------------------------------------------------

  // Parse structured results if available; fall back to raw output.
  let tableRows = '';
  let summaryLine = '';
  const rawOutput = read('output.txt');

  try {
    const resultsJson = read('results.json');
    if (resultsJson) {
      const data = JSON.parse(resultsJson);
      summaryLine = `**${data.failed} of ${data.total} template(s) failed.**`;

      // Build a markdown table row for each failed template.
      // SDK Range shows what package.json / go.mod specifies.
      // Resolved shows what was actually installed/used — so a
      // '^1.5.0' range resolving to '2.0.0' is immediately visible.
      const failedTemplates = data.templates.filter((t) => t.status !== 'pass');
      const rows = failedTemplates.map((t) => {
        // Sanitise all values read from the artifact.
        const sanitise = (s) => String(s || '').replace(/[|`]/g, ' ').slice(0, 80);
        const name = sanitise(t.name);
        const lang = sanitise(t.language);
        const range = sanitise(t.sdk_range);
        const resolved = sanitise(t.sdk_resolved);
        const step = sanitise(t.failure_step);
        // Show a short excerpt of the first real error line
        const output = (t.failure_output || '')
          .split('\n')
          .map((l) => l.trim())
          .filter((l) => l && !l.startsWith('npm warn') && !l.startsWith('>'))
          .slice(0, 3)
          .join(' · ')
          .replace(/[|`]/g, ' ')
          .slice(0, 120);
        return `| \`${name}\` | ${lang} | \`${range}\` | \`${resolved}\` | ${step} | ${output} |`;
      });

      if (rows.length > 0) {
        tableRows = [
          '| Template | Language | SDK Range | Resolved | Step | Error |',
          '|----------|----------|-----------|----------|------|-------|',
          ...rows,
        ].join('\n');
      }
    }
  } catch (e) {
    console.log('Could not parse results.json, falling back to raw output:', e.message);
  }

  // If we couldn't build a table, fall back to a plain code block.
  let failureSection = '';
  if (tableRows) {
    failureSection = tableRows;
  } else {
    const resultsMatch = rawOutput.match(/={8,}\nResults:.*\n={8,}[\s\S]*/);
    const excerpt = resultsMatch ? resultsMatch[0].trim() : rawOutput.trim();
    failureSection = '```\n' + excerpt.slice(0, 3000) + '\n```';
  }

  const body = [
    '## ⚠️ Template Check Failures',
    '',
    summaryLine || 'One or more templates changed in this PR failed validation or compilation.',
    '',
    '> **SDK Range** is what `package.json` / `go.mod` specifies.',
    '> **Resolved** is the exact version that was installed and used.',
    '> A range like `^1.5.0` resolving to `2.0.0` indicates a breaking SDK release.',
    '',
    failureSection,
    '',
    `[View full output →](${runUrl})`,
    '',
    '<details>',
    '<summary>What should I do?</summary>',
    '',
    '- **Fix the template:** Update the template so it compiles against the SDK version it specifies.',
    '- **Pin the SDK version:** Change the range (e.g. `^1.5.0`) to an exact version (e.g. `1.5.2`) if you want to lock against a known-good release.',
    '',
    '</details>',
  ].join('\n');

  const commentBody = `${marker}\n${body}`;

  if (existing) {
    await github.rest.issues.updateComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      comment_id: existing.id,
      body: commentBody,
    });
    console.log('Updated existing failure comment.');
  } else {
    await github.rest.issues.createComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: prNumber,
      body: commentBody,
    });
    console.log('Posted new failure comment.');
  }
};
