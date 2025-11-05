# Quick Start Guide

Get RSS Watcher running in under 5 minutes!

## 1. Get Bark Device Key (1 minute)

1. Install [Bark](https://apps.apple.com/app/bark-customed-notifications/id1403753865) on your iPhone
2. Open the app
3. Copy your device key (the long string in the URL shown)

## 2. Fork This Repository (30 seconds)

Click the "Fork" button at the top right of this page.

## 3. Add Your Bark Key (30 seconds)

In your forked repository:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `BARK_DEVICE_KEY`
4. Value: [paste your device key]
5. Click **Add secret**

## 4. Add Your RSS Feeds (2 minutes)

1. Click on `feeds.yaml` in your repository
2. Click the pencil icon (Edit)
3. Replace the content with:

```yaml
feeds:
  - id: my-first-feed
    name: My Blog
    url: https://example.com/rss.xml  # Replace with your RSS feed URL
    notify: true
    dedupe_key: guid
    aggregate: false
```

4. Click **Commit changes**

## 5. Test It! (1 minute)

1. Go to the **Actions** tab
2. Click **RSS Monitor (Go + Bark)** on the left
3. Click **Run workflow** button (top right)
4. Click the green **Run workflow** button
5. Wait ~30 seconds
6. Check your iPhone for a notification! 🎉

## Done!

Your RSS monitor will now run automatically every 30 minutes.

## Next Steps

- Add more feeds to `feeds.yaml`
- Read the full [README.md](README.md) for all features
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for advanced configuration

## Troubleshooting

**No notification?**

1. Test Bark: `curl "https://api.day.app/YOUR_KEY/Test/Hello"`
2. Check Actions logs for errors
3. Verify your feed URL works in a browser

**Need help?** [Open an issue](https://github.com/rsswatcher/rsswatcher/issues)
