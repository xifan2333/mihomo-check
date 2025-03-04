# Gist Saving Method

## Deployment

- Create a Gist at your convenience.

- Configure the gist id in `config.yaml`.

- Configure the gist token in `config.yaml`.

## Worker Reverse Proxy for GitHub API

- Deploy the [worker](./cloudflare/worker.js) to Cloudflare Workers.

- Set `GITHUB_USER` in `Variables and Secrets` to your GitHub username.

- Set `GITHUB_ID` in `Variables and Secrets` to your gist id.

- Set `AUTH_TOKEN` in `Variables and Secrets` to your access token.

- Configure `github-api-mirror` to your worker address.

```
    github-api-mirror: "https://your-worker-url/github"
```

## Subscription Retrieval

> If the Worker is configured, change the `key` accordingly.
> The subscription format is `https://your-worker-url/gist?key=all.yaml&token=AUTH_TOKEN`.

- Subscribe to all

```
https://gist.githubusercontent.com/YOUR_GITHUB_USERNAME/YOUR_GIST_ID/raw/all.yaml
```

- Unlock OpenAI nodes

```
https://gist.githubusercontent.com/YOUR_GITHUB_USERNAME/YOUR_GIST_ID/raw/openai.yaml
```

- Unlock Netflix nodes

```
https://gist.githubusercontent.com/YOUR_GITHUB_USERNAME/YOUR_GIST_ID/raw/netflix.yaml
```

- Unlock Disney nodes

```
https://gist.githubusercontent.com/YOUR_GITHUB_USERNAME/YOUR_GIST_ID/raw/disney.yaml
```

- Unlock YouTube nodes

```
https://gist.githubusercontent.com/YOUR_GITHUB_USERNAME/YOUR_GIST_ID/raw/youtube.yaml
```
