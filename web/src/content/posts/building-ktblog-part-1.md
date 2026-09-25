---
title: What I Learned Building a Static Blog, Part 1
pubDate: 2026-09-26
---

> Who this is for: backend developers with no frontend experience.

## Introduction

I'm a backend developer learning everything else. This blog is both the notes and the project: I'm building it from scratch — frontend, backend, databases, and deployment — one milestone at a time.

The project is split into seven milestones:

- **M1: Static blog live.** Astro, hand-written CSS, Markdown posts, deployed to Cloudflare Pages.
- **M2: Frontend polish.** Tailwind, tags, dark mode, RSS, SEO, and search.
- **M3: Go API and first VPS.** Views and likes served by a Go API, deployed with systemd and Caddy.
- **M4: PostgreSQL.** Persist views and likes, and add comments.
- **M5: Containers and CI/CD.** Docker Compose and GitHub Actions.
- **M6: Operations and security.** Logs, alerts, error tracking, rate limiting, and spam protection.
- **M7: Electives.** Whatever I find interesting by then.

This post covers what I learned in M1.

## npm, npx, and nvm

Before this project, my frontend experience was limited to running existing projects and making small changes. I used these tools without really knowing how they fit together. Here they are, mapped to tools I already knew:

| Node                | Python             | Java / Kotlin      |
|---------------------|--------------------|--------------------|
| `node`              | `python`           | JVM                |
| `npm`               | pip / uv           | Maven / Gradle     |
| `package.json`      | `pyproject.toml`   | `build.gradle.kts` |
| `package-lock.json` | `uv.lock`          | Gradle lockfile    |
| `node_modules/`     | `.venv/`           | Gradle cache       |
| `npx`               | `uvx` / `pipx run` | —                  |
| `nvm`               | pyenv              | SDKMAN             |

- **npm** is the package manager. It installs packages into the project's own `node_modules/` by default, so every project is isolated, a bit like having a virtualenv you never need to activate.
- **npx** runs a command provided by a package. It looks in the project's `node_modules/.bin` first, and if the package isn't there, downloads it to a cache and runs it from there. `npm create astro@latest` is shorthand for `npx create-astro@latest`, a project generator in the same spirit as Spring Initializr.
- **nvm** installs and switches between Node versions. Installing a Node version also installs its npm, so every Node version has its own npm.

### `.nvmrc` vs. `engines.node`

Both files mention a Node version, but they do different jobs:

- `.nvmrc` **selects** a version (`24`). nvm reads it, and so does Cloudflare Pages, so my laptop and the build server run the same Node. It's the equivalent of `.python-version`.
- `engines.node` in `package.json` **constrains** the range (`>=22.12.0`). npm only warns when it isn't met. It's the equivalent of `requires-python`.

## Deploying a monorepo to Cloudflare Pages

The repository holds both the frontend and, later, the Go API:

```text
web/     Astro frontend
api/     Go API
```

Setting up Pages:

1. Create an account on [Cloudflare](https://www.cloudflare.com/).
2. Go to **Workers & Pages** and select **Create application**.
3. The dashboard now defaults to Workers. Pages is labeled "legacy", and you get to it through **Continue to Pages**.
4. Select **Import an existing Git repository**, then choose your GitHub account and repository.
5. Fill in the build settings:

| Setting                | Value           |
|------------------------|-----------------|
| Framework preset       | Astro           |
| Build command          | `npm run build` |
| Build output directory | `dist`          |
| Root directory         | `web`           |

Two details are easy to get wrong. The root directory must be `web`, or Cloudflare looks for `package.json` in the repository root and the build fails. And the output directory is relative to the root directory, so it's `dist`, not `web/dist`.

After that, every push to `main` deploys automatically, and every pull request gets its own preview URL. HTTPS, HTTP/2, and a CDN come for free.

Why Pages if it's labeled legacy? It still works, my roadmap was already built around it, and moving a static site later is cheap: the build output is just files in `dist/`.

## Layouts, slots, and scoped CSS

### Component structure

An Astro component has two parts:

```astro
---
// Component script (TypeScript), runs at build time
---
<!-- Component template (HTML + expressions + CSS) -->
```

### Component props

Props let you pass values into a component, like function arguments:

```astro
---
interface Props {
    title: string;
}

const { title } = Astro.props;
---
<h1>Title: {title}</h1>
```

Then you use it like this:

```astro
---
import Foo from '../components/Foo.astro';
---

<Foo title="Hello World" />
```

The `Props` interface isn't just documentation. If you forget `title` or pass the wrong type, the editor and `astro check` flag it, much like a Pydantic model.

### Slot

`<slot />` is a placeholder for content passed in from the outside. If you've used Jinja's `{% block %}` or Thymeleaf layouts, it's the same idea:

```astro
---
interface Props {
    title: string;
}

const { title } = Astro.props;
---
<h1>Title: {title}</h1>
<div><slot /></div>
```

Anything between the component's tags goes into the slot:

```astro
<Foo title="Hello World">
    <p>This is the content</p>
</Foo>
```

This is how my `BaseLayout` works: the header and footer live in the layout, and each page only provides its own content.

### Scoped CSS

A `<style>` block inside an Astro component only applies to that component. To see what that means, I added this to my layout:

```css
h1 { color: red; }
```

The `<h1>` on my homepage didn't turn red, even though it ends up inside the layout's `<main>`. It's written in the page file, passed in through the slot, so it doesn't belong to the layout. Changing the selector to `header a` did work, because that link is written in the layout itself.

The browser's dev tools show how this works. Astro adds an attribute such as `data-astro-cid-z4jru4n3` to every element a component writes, and rewrites its selectors to match only elements with that attribute:

```css
h1[data-astro-cid-z4jru4n3] { color: red; }
```

So site-wide styles, like fonts and the text inside Markdown posts, belong in a global stylesheet imported by the layout. Styles that only belong to one component stay scoped.

## Content Collections

### Zod

Zod is a runtime validator and compile-time type helper that ensures data matches a specific structure. If you use Python, think of it as Pydantic for TypeScript.

Every post needs at least a title and a publish date, so the site can display them. They go in the Markdown frontmatter:

```markdown
---
title: Hello World
pubDate: 2026-09-25
---
Your content goes here.
```

And `src/content.config.ts` uses Zod to make sure every Markdown file has them:

```ts
import { defineCollection } from 'astro:content';
import { glob } from 'astro/loaders';
import { z } from 'astro/zod';

const posts = defineCollection({
    loader: glob({ base: './src/content/posts', pattern: '*.md' }),
    schema: z.object({
        title: z.string(),
        pubDate: z.coerce.date(),
    }),
});

export const collections = { posts };
```

If a post is missing its title or has an invalid date, the build fails instead of publishing a broken page.

### getCollection

`getCollection()` is how you query a collection. Astro reads the `collections` object exported from `src/content.config.ts` and gives you a query function for each key. The string you pass in is that key.

```html
---
import { getCollection } from 'astro:content';

const posts = await getCollection('posts');
posts.sort((a, b) => b.data.pubDate.getTime() - a.data.pubDate.getTime());
---

<ul>
    {posts.map((post) => (
        <li><a href={`/posts/${post.id}/`}>{post.data.title}</a></li>
    ))}
</ul>
```

If you come from Spring Data, this should feel familiar. You declare `interface PostRepository extends JpaRepository<Post, Long>` and get `findAll()` for free. Here you declare a collection and get `getCollection()`.

It also takes an optional filter function, such as `getCollection('posts', (post) => ...)`, which is useful for things like hiding drafts.

The big difference from a repository is **when** it runs. There is no database and no server: `getCollection()` reads the Markdown files once, at build time, and the result is baked into plain HTML.

### getStaticPaths

Each post gets its own page from a single file: `src/pages/posts/[id].astro`. The square brackets make `id` a route parameter, much like `@app.get("/posts/{id}")` in FastAPI.

But a FastAPI route handles requests as they arrive. A static site has no server waiting for requests, so every page must already exist as an HTML file before anyone visits. Astro needs to know every possible `id` up front, and `getStaticPaths()` is where you tell it:

```astro
---
import { getCollection, render } from 'astro:content';
import BaseLayout from '../../layouts/BaseLayout.astro';

export async function getStaticPaths() {
    const posts = await getCollection('posts');
    return posts.map((post) => ({
        params: { id: post.id },  // fills [id] in the URL
        props: { post },          // passed to the page as Astro.props
    }));
}

const { post } = Astro.props;
const { Content } = await render(post);  // Markdown -> component
---

<BaseLayout title={post.data.title}>
    <h1>{post.data.title}</h1>
    <Content />
</BaseLayout>
```

For every object in the returned array, Astro renders the page once and writes it to disk. `hello-world.md` becomes `dist/posts/hello-world/index.html`, which Cloudflare Pages serves at `/posts/hello-world/`.

This also means adding a post is a build step, not a database insert. Push a new Markdown file, and the page only exists after the next build.

## Gotchas

### Time zones

A post's date is `2026-09-25`, and the post page displayed it like this:

```astro
<time>{post.data.pubDate.toLocaleDateString()}</time>
```

It looked right on my machine. When I ran the same build with the time zone set to Los Angeles, the page said September 24.

The reason is that this code runs **at build time**. YAML parses `2026-09-25` as midnight UTC, and `toLocaleDateString()` converts it to the time zone of the machine running the build. Every reader sees the same HTML, so what varies is not the reader's time zone but the build server's. The locale (`9/25/2026` vs. `25/09/2026`) depends on the build machine too.

The fix is to set both explicitly:

```ts
date.toLocaleDateString('en-GB', { timeZone: 'UTC' })
```

### Data structure from getCollection

My first version of the sort on the homepage looked like this:

```ts
posts.sort((a, b) => b.pubDate.getTime() - a.pubDate.getTime());
```

It's wrong, and the build still passed. To see why it's wrong, here is what one entry from `getCollection('posts')` actually looks like (dumped as JSON, trimmed):

```json
{
  "id": "hello-world",
  "collection": "posts",
  "data": {
    "title": "Hello World",
    "pubDate": "2026-09-25T00:00:00.000Z"
  },
  "body": "This is a testing post\n\nThis is a list\n* Item\n...",
  "filePath": "src/content/posts/hello-world.md",
  "digest": "ece2a0ef7d55f40b",
  "rendered": {
    "html": "<p>This is a testing post</p>\n<p>This is a list</p>\n<ul>...",
    "metadata": { "headings": [], "frontmatter": { ... }, ... }
  }
}
```

| Field        | What it is                                                              |
|--------------|-------------------------------------------------------------------------|
| `id`         | Generated from the file name. It becomes the URL slug.                  |
| `collection` | Which collection the entry belongs to.                                  |
| `data`       | Your frontmatter, validated and converted by your Zod schema.           |
| `body`       | The raw Markdown, without the frontmatter.                              |
| `filePath`   | Where the file lives, relative to the project root.                     |
| `digest`     | A hash of the content, used by Astro for caching.                       |
| `rendered`   | The Markdown already converted to HTML, plus metadata such as headings. |

Your fields live under `data`, and Astro's own fields sit next to it. This keeps them from colliding: if a post's frontmatter had its own `id`, it would be `post.data.id` and would never overwrite Astro's `post.id`. It works like `ResponseEntity` in Spring, where headers and status sit outside and your payload is in `getBody()`.

`pubDate` looks like a string above only because JSON has no date type. In code it is a real `Date`, because the schema uses `z.coerce.date()`. That's why `.getTime()` works on it.

So the fix is `b.data.pubDate`. The more interesting question is why the broken version built fine:

1. **There was only one post.** With a single element, `sort()` has nothing to compare, so the comparison function never ran and never threw. A second post would have broken the build. Code that never runs is code that was never tested.
2. **`astro build` does not type-check.** It compiles TypeScript but doesn't verify types. `astro check` (from the `@astrojs/check` package) reports the error right away: `Property 'pubDate' does not exist`.

### Installing into the wrong directory

In a monorepo, npm acts on whichever directory you run it in. I ran `npm install @astrojs/check typescript` from the repository root instead of `web/`. There was no `package.json` there, so npm created a new one and installed the packages into the root. The frontend project never got them. Now I run `pwd` before any npm command.

## What's next

M2 is frontend polish: Tailwind, dark mode, RSS, SEO, and search, with a Lighthouse score of 90 or more as the goal.
