# Deployment Guide

This guide walks you through deploying RSS Watcher to GitHub Actions.

## Prerequisites

- GitHub account
- iOS device with Bark app installed
- RSS feeds you want to monitor

## Step-by-Step Deployment

### 1. Get Your Bark Device Key

1. Install [Bark](https://apps.apple.com/app/bark-customed-notifications/id1403753865) from the App Store
2. Open the app
3. You'll see a URL like: `https://api.day.app/YOUR_DEVICE_KEY/`
4. Copy the `YOUR_DEVICE_KEY` part (you'll need this in step 4)

### 2. Fork or Clone This Repository

**Option A: Fork (Recommended for most users)**
1. Click the "Fork" button at the top right of this repository
2. This creates a copy in your GitHub account

**Option B: Clone and Push to New Repo**
```bash
git clone https://github.com/rsswatcher/rsswatcher.git
cd rsswatcher
# Create a new repository on GitHub, then:
git remote set-url origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
git push -u origin main
```

### 3. Configure Your RSS Feeds

1. Edit the `feeds.yaml` file in your repository
2. Add your RSS feeds following this format:

```yaml
feeds:
  - id: my-tech-blog
    name: Tech Blog
    url: https://techblog.example.com/rss
    notify: true
    dedupe_key: guid
    aggregate: false

  - id: news-feed
    name: Daily News
    url: https://news.example.com/feed
    notify: true
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
```

**Field Explanations:**
- `id`: Unique identifier (no spaces, use hyphens)
- `name`: Display name for notifications
- `url`: Full URL to the RSS/Atom feed
- `notify`: Set to `true` to receive notifications
- `dedupe_key`: How to identify unique items (`guid`, `link`, or `title`)
- `aggregate`: `true` to combine multiple new items into one notification
- `aggregate_window_minutes`: Time window for aggregation (optional)

### 4. Configure GitHub Secrets

1. Go to your repository on GitHub
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secrets:

| Name | Value | Required |
|------|-------|----------|
| `BARK_DEVICE_KEY` | Your Bark device key from step 1 | Yes |
| `BARK_SERVER` | Custom Bark server (if self-hosting) | No |

**To add a secret:**
1. Click "New repository secret"
2. Enter the name (e.g., `BARK_DEVICE_KEY`)
3. Paste the value
4. Click "Add secret"

### 5. Enable GitHub Actions

GitHub Actions should be enabled by default, but to verify:

1. Go to the **Actions** tab in your repository
2. If you see a message about workflows being disabled, click "I understand my workflows, go ahead and enable them"

### 6. Test the Workflow

**Manual Test:**
1. Go to the **Actions** tab
2. Click on "RSS Monitor (Go + Bark)" in the left sidebar
3. Click "Run workflow" button (on the right)
4. Select the branch (usually `main`)
5. Click the green "Run workflow" button
6. Wait for the workflow to complete
7. Check your iOS device for notifications!

### 7. Verify Automatic Execution

The workflow is configured to run every 30 minutes. To verify:

1. Wait for the scheduled time
2. Go to the **Actions** tab
3. You should see automatic runs appearing
4. Click on a run to see detailed logs

## Customizing the Schedule

To change how often the workflow runs, edit `.github/workflows/rss-monitor.yml`:

```yaml
on:
  schedule:
    - cron: '*/30 * * * *'  # Current: Every 30 minutes
```

**Common cron schedules:**
- Every 15 minutes: `*/15 * * * *`
- Every hour: `0 * * * *`
- Every 6 hours: `0 */6 * * *`
- Every day at 9 AM: `0 9 * * *`

**Note:** GitHub Actions may delay scheduled workflows during high-load periods. Expect 5-10 minute delays.

## Troubleshooting

### No Notifications Received

**Check 1: Verify Bark is working**
```bash
curl "https://api.day.app/YOUR_DEVICE_KEY/Test/Hello"
```
You should receive a test notification.

**Check 2: Verify GitHub Actions ran**
1. Go to Actions tab
2. Click on the latest workflow run
3. Look for errors in the logs

**Check 3: Verify secrets are set**
1. Settings → Secrets and variables → Actions
2. Ensure `BARK_DEVICE_KEY` is listed

**Check 4: Check feed configuration**
- Ensure `notify: true` is set for your feeds
- Verify feed URLs are accessible
- Check that feeds have new content

### Workflow Fails to Run

**Check 1: Actions enabled**
- Go to Settings → Actions → General
- Ensure "Allow all actions and reusable workflows" is selected

**Check 2: Workflow syntax**
- Check for YAML syntax errors in workflow files
- Use a YAML validator

**Check 3: Repository permissions**
- Settings → Actions → General → Workflow permissions
- Select "Read and write permissions"
- Check "Allow GitHub Actions to create and approve pull requests"

### State Not Updating

**Check 1: Workflow permissions**
The workflow needs write access to commit state changes:
1. Settings → Actions → General
2. Under "Workflow permissions", select "Read and write permissions"

**Check 2: Branch protection**
If your `main` branch is protected:
1. Settings → Branches
2. Edit branch protection rules
3. Allow github-actions bot to push

### Rate Limiting or Feed Errors

If you see too many errors in logs:

1. **Reduce feed count** - Remove some feeds temporarily
2. **Increase interval** - Change from 30 to 60 minutes
3. **Check feed URLs** - Some feeds may be down or rate-limited
4. **Add delays** - Modify code to add delays between feeds

## Advanced Configuration

### Self-Hosted Bark Server

If you run your own Bark server:

1. Deploy Bark server (see [Bark documentation](https://github.com/Finb/Bark))
2. Add `BARK_SERVER` secret with your server URL
3. Example: `https://bark.yourdomain.com`

### Multiple Devices

To notify multiple devices:

**Option 1: Bark Groups** (Recommended)
Set up group notifications in Bark app, then use group names in feed configuration.

**Option 2: Multiple Secrets**
1. Add multiple device key secrets: `BARK_DEVICE_KEY_1`, `BARK_DEVICE_KEY_2`
2. Modify the notifier code to send to multiple keys

### Custom Notification Format

Edit `internal/notifier/bark.go` to customize:
- Notification title format
- Message content
- Notification sounds
- Icons and badges

### External State Storage

For advanced users who want state stored externally (Redis, S3, etc.):

1. Modify `internal/state/state.go`
2. Implement remote storage backend
3. Update GitHub Actions workflow to not commit state

## Monitoring and Maintenance

### Check Workflow Status

Create a simple dashboard:
1. Star your repository
2. Check https://github.com/YOUR_USERNAME?tab=stars
3. Look for action status badges

### Add Status Badge

Add to your README.md:
```markdown
![RSS Monitor](https://github.com/YOUR_USERNAME/rsswatcher/actions/workflows/rss-monitor.yml/badge.svg)
```

### Review Logs

Periodically check workflow logs for:
- Failed feed fetches
- Notification errors
- Performance issues

### Update Dependencies

Every few months:
```bash
go get -u ./...
go mod tidy
git commit -am "Update dependencies"
git push
```

## Security Best Practices

1. **Never commit secrets** - Always use GitHub Secrets
2. **Review state file** - Ensure no sensitive data is stored
3. **Use HTTPS feeds** - Prefer secure feed URLs
4. **Monitor access** - Regularly review repository access
5. **Keep dependencies updated** - Update Go modules regularly

## Cost Considerations

- **GitHub Actions**: 2,000 minutes/month free
- **Typical usage**: ~5 minutes/day = 150 minutes/month
- **Result**: Completely free for most users!

If you exceed limits:
- Reduce check frequency
- Remove some feeds
- Upgrade to GitHub Pro (more free minutes)

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section above
2. Review [GitHub Actions logs](https://github.com/YOUR_USERNAME/YOUR_REPO/actions)
3. Search [existing issues](https://github.com/rsswatcher/rsswatcher/issues)
4. Create a new issue with:
   - Description of the problem
   - Relevant logs (remove sensitive data!)
   - Your configuration (remove secrets!)

## Next Steps

After successful deployment:

1. Add more RSS feeds to `feeds.yaml`
2. Customize notification formats
3. Adjust check frequency
4. Share your setup with others!

Enjoy your automated RSS notifications! 🎉
