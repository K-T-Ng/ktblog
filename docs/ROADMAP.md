# Roadmap

Build a personal blog from scratch to learn frontend, backend, databases, and deployment end to end. Each milestone ships something that works and ends with a blog post about it.

## Milestones

- **M1: Static blog live.** Build the site with Astro and hand-written CSS, write posts in Markdown, and publish it on a custom domain with Cloudflare Pages.
- **M2: Frontend polish.** Move to Tailwind, then add tags, dark mode, RSS and SEO, and search. Aim for Lighthouse scores of 90 or more.
- **M3: Go API and first VPS.** Start with the standard `net/http` library, then move to Gin. Add views and likes held in memory with safe concurrency, a React island to show them, and a VPS deployment with systemd and Caddy.
- **M4: PostgreSQL.** Persist views and likes and add comments using pgx, sqlc, and goose. Learn transactions and isolation levels, and restore from a backup at least once.
- **M5: Containers and CI/CD.** Move to Docker Compose, and let GitHub Actions run tests and lint and deploy automatically.
- **M6: Operations and security.** Add structured logs, uptime alerts, error tracking, rate limiting, and spam protection.
- **M7: Electives.** Pick by interest: admin and auth, Redis, full-text search, object storage, WebSockets, Prometheus and Grafana.

## Tech Stack

| Layer | Choice |
|---|---|
| Frontend | Astro + TypeScript; CSS → Tailwind; React islands + shadcn/ui |
| Content | Markdown in git (Astro Content Collections) |
| Backend | Go (standard library → Gin) |
| Database | PostgreSQL with pgx, sqlc, goose |
| Hosting | Cloudflare Pages (web); VPS with Caddy (api) |
| Deploy | systemd → Docker Compose + GitHub Actions |
