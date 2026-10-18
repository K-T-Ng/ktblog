---
title: What I Learned Building a Static Blog, Part 2
pubDate: 2026-10-18
description: "Notes for backend developers new to frontend: adding Tailwind, tags, dark mode, RSS, SEO, and search to an Astro blog, while keeping Lighthouse at 90+."
tags: ["astro", "cloudflare", "tailwind"]
---

> Who this is for: backend developers with no frontend experience. [Part 1](/posts/building-ktblog-part-1/) covers the basics: npm, Astro components, and deploying to Cloudflare Pages.

## Introduction

I'm a backend developer learning everything else. This blog is both the notes and the project: I'm building it from scratch — frontend, backend, databases, and deployment — one milestone at a time.

In [Part 1](/posts/building-ktblog-part-1/), I got a bare-bones blog live: Astro, hand-written CSS, and Markdown posts, deployed to Cloudflare Pages. It worked, but that was about all it did.

M2 was about polishing the frontend. I added:

- **Tailwind CSS** for styling, plus its typography plugin for rendered Markdown.
- **Tags**, with a page for each tag.
- **Dark mode** that follows the OS setting, plus a toggle that remembers your choice.
- **RSS and SEO**: a feed, meta descriptions, Open Graph tags, a sitemap, and a real 404 page.
- **Search** that runs entirely in the browser.

The goal was a Lighthouse score of 90 or more in every category. The site already scored above 90 before I started, so the real challenge was keeping it there while adding features.

Looking back, one theme kept coming up: the CSS cascade, the rules that decide which style wins when several of them apply to the same element. It surprised me three times in this milestone, and I'll point it out each time.

## Measure first

Before changing anything, I measured a baseline of the M1 site, so I'd have something to compare against.

The tool is [Lighthouse](https://developer.chrome.com/docs/lighthouse/overview), Google's open-source tool for auditing web pages. It's built into Chrome DevTools (the Lighthouse tab) and scores a page from 0 to 100 in four categories: Performance, Accessibility, Best Practices, and SEO. For a backend developer, think of it as a benchmark and a linter rolled into one.

Here's the baseline, measured on a local production build (`npm run build && npm run preview`):

| Page | Device  | Performance | Accessibility | Best Practices | SEO | LCP   | TBT   | CLS |
|------|---------|-------------|---------------|----------------|-----|-------|-------|-----|
| Home | Desktop | 100         | 91            | 100            | 90  | 0.2 s | 0 ms  | 0   |
| Home | Mobile  | 100         | 91            | 100            | 90  | 0.9 s | 50 ms | 0   |
| Post | Desktop | 100         | 93            | 100            | 90  | 0.2 s | 0 ms  | 0   |
| Post | Mobile  | 100         | 93            | 100            | 90  | 0.9 s | 0 ms  | 0   |

The last three columns are the key performance metrics:

- **LCP** (Largest Contentful Paint): how fast the largest piece of content appears.
- **TBT** (Total Blocking Time): how long JavaScript blocks the main thread.
- **CLS** (Cumulative Layout Shift): how much the layout jumps around while loading.

Every score was already 90 or above, since Astro ships no JavaScript by default. The points were lost on two audits: a missing meta description, and low color contrast in the footer and in code comments. Both get fixed later in this post.

A few traps can make the numbers misleading:

1. **Measuring the dev server.** `astro dev` and `astro preview` both default to port 4321, so the URL alone doesn't tell you which one you're testing. I once measured the dev server by mistake, and mobile Performance dropped to 84. The extra weight came from Vite's hot-reload client and Astro's dev toolbar, none of which ships to production. It's like benchmarking a debug build.
2. **Measuring SEO on a preview URL.** Cloudflare Pages preview deployments send an `x-robots-tag: noindex` header to keep them out of search results. Lighthouse penalizes pages that can't be indexed, so SEO dropped to 45. Measure SEO on production only.
3. **Browser extensions.** Extensions logged errors in the console, which cost 4 points in Best Practices. One error came from an ad blocker blocking Cloudflare's analytics beacon, a script that isn't in my source code at all. Cloudflare injects it into the HTML at the edge, so what's served isn't exactly what was built. Running Lighthouse in an incognito window keeps extensions out of the measurement.

## Tailwind 4

[Tailwind](https://tailwindcss.com/) is a utility-first CSS framework. Instead of naming a class and writing its CSS in a separate file, you compose small, single-purpose classes directly in the HTML.

Before, with hand-written CSS:

```css
/* global.css */
.post-card {
    padding: 1rem;
    border-radius: 0.5rem;
    background-color: white;
}
```

```html
<div class="post-card">
```

After, with Tailwind, the same three declarations as classes:

```html
<div class="p-4 rounded-lg bg-white">
```

States and conditions are handled by **variants**: a prefix that makes a class apply only in that state.

```html
<div class="p-4 rounded-lg bg-white hover:bg-gray-100 dark:bg-gray-800">
```

`hover:` applies when the mouse is over the element, and `dark:` applies in dark mode. In plain CSS, each of those would be another rule: a `.post-card:hover` selector, and a dark-mode rule, which we'll get to in the dark mode section.

What I liked:

- **No naming, no jumping between files.** The styles sit on the element they style.
- **A consistent scale.** `p-4`, `gray-800`, and friends come from a fixed set of spacing and color values, so the design stays consistent without any discipline on my part.
- **Easy deletion.** Delete a component and its styles go with it. No dead CSS left behind.
- **Build-time only.** Tailwind scans the source and generates CSS for just the classes you use. No JavaScript ships to the browser.

What I didn't:

- **Long class lists.** Some elements end up with a dozen classes.
- **A new vocabulary.** You keep the [docs](https://tailwindcss.com/docs) open for a while.
- **It's still CSS.** Tailwind doesn't replace knowing CSS; it's shorthand for it.

Setup was one command, `npx astro add tailwind`. Version 4 is configured in CSS (`@import "tailwindcss";` in a stylesheet) through a Vite plugin, so most v3 tutorials with a `tailwind.config.js` no longer apply.

### Preflight: where did my headings go?

Right after installing, the site looked broken: headings were the same size as body text, links lost their blue underline, and list bullets disappeared. Nothing was actually broken. That was [Preflight](https://tailwindcss.com/docs/preflight), Tailwind's CSS reset:

```css
h1, h2, h3, h4, h5, h6 {
    font-size: inherit;
    font-weight: inherit;
}
```

Every browser ships its own default stylesheet, and the defaults differ between browsers. Preflight wipes them out so every browser starts from the same blank page, and you add back exactly what you want. The idea is "opt in instead of opt out": with utilities, you'd otherwise spend half your time undoing defaults, like removing heading margins.

### The cascade, part 1

My first reaction was "Preflight overrode my CSS". It didn't. My CSS never set heading sizes. Those came from the browser's own stylesheet, and Preflight overrode *that*.

This is the cascade: the CSS rules for deciding which declaration wins when several apply to the same element. I only know the specificity part: an id beats a class, and a class beats a tag. But the cascade checks other things first. For normal declarations, from weakest to strongest:

1. **The browser's default stylesheet** (the "user-agent stylesheet").
2. **Your CSS inside a cascade layer** (`@layer`), in the order the layers are declared. Tailwind declares `theme`, `base`, `components`, and `utilities`. Preflight lives in `base`, and classes like `p-4` live in `utilities`.
3. **Your CSS outside any layer**, such as a hand-written `global.css` or an Astro `<style>` block.

Specificity only breaks ties *within* the same level. A rule from a higher level wins no matter how specific the lower-level rule is.

That has a sharp edge while old CSS and Tailwind coexist:

```css
/* global.css, hand-written, no layer */
a {
    color: black;
}
```

```html
<a class="text-blue-700" href="/">Home</a>
```

The link stays black. `text-blue-700` is a more specific selector (a class beats a tag), but it lives in a layer and the hand-written rule doesn't, so the hand-written rule wins. No warning, no error. So as I migrated each part of the site to Tailwind, I deleted its old CSS right away.

To see this in action, inspect an `h1` in DevTools. The Styles panel shows the Preflight rule under `@layer base`, with the user-agent rule crossed out below it.

### Typography for Markdown

Utility classes need an element to sit on. Post content is different: Astro renders it from Markdown, so I can't put a class on its `<h2>` or `<ul>`, and Preflight had stripped them bare.

The [typography plugin](https://github.com/tailwindlabs/tailwindcss-typography) solves this. Wrap the content in `<article class="prose">`, and it styles every heading, list, quote, and code block inside with sensible defaults. It's the opposite of utility-first: one class that targets a whole tree of elements you don't control. Utilities for my own markup, `prose` for rendered Markdown.

You can tweak it with element modifiers in the form `prose-{element}:{utility}`, like `prose-a:text-blue-700` for link color.

That's also where I hit my first real bug. I gave inline code a gray background with `prose-code:bg-gray-100`. It also matched the `<code>` *inside* syntax-highlighted code blocks, which put light gray under light text, at a contrast of about 1.1:1, almost invisible. Markdown produces `<code>` in two places, inline and inside `<pre>`, and the modifier targets both. The fix was an arbitrary variant that only matches code whose parent isn't a `<pre>`:

```html
<article class="prose [&_:not(pre)>code]:bg-gray-100">
```

`&` is the element itself, `_` stands for a space (a descendant), and `:not(pre)>code` is a `<code>` whose parent isn't `<pre>`. The lesson: when a selector targets an element type, check every place that element can appear.

## Tags

Tags group posts by topic. A post about deployment gets a `cloudflare` tag, and `/tags/cloudflare/` lists every post with that tag. This part was less about CSS and more about TypeScript, and about what it does and doesn't check.

### Validating tags: regex or enum?

Tags live in each post's frontmatter, and the content collection schema validates them at build time. I considered two ways to validate them:

- **A regex** accepts any tag in the right format, for example lowercase kebab-case. That stops `Cloudflare` or `cloudflare pages`, but not typos: `claudeflare` passes, and so do both `cloudflare-page` and `cloudflare-pages`, which silently creates two tag pages for the same topic.
- **An enum** accepts only the tags I've listed. Typos and near-duplicates fail the build. The cost is that a new tag has to be added to the schema before a post can use it.

I went with the enum. On a solo blog, the extra step is small, and it even helps: it pushes me to reuse existing tags instead of inventing new ones. It's the same idea as Pydantic's `Literal[...]` or a Kotlin `enum class`.

Here's the schema in `content.config.ts`:

```typescript
tags: z.array(z.enum(["astro", "cloudflare", "node-js", "tailwind"]))
    .default([])
    .refine(
        (tags) => tags.length === new Set(tags).size, {
            error: "Duplicate tags are not allowed."
        }
    ),
```

- `.default([])` lets a post leave out `tags` entirely.
- `.refine()` adds a custom rule. A `Set` drops duplicates, so if it's smaller than the array, a tag was repeated.

A typo fails the build with an error that points at the exact tag and lists the allowed values:

```text
tags.2: Invalid option: expected one of "astro"|"cloudflare"|"node-js"|"tailwind"
```

One gotcha: build the tag pages from the tags the posts actually use, not from the enum. Otherwise a tag nobody uses yet gets an empty page.

### A successful build doesn't mean the types are right

While extracting the post list into a `PostList` component, I typed its prop wrong:

```typescript
interface Props {
    posts: CollectionEntry<"posts">;  // missing [], should be an array
}
```

The IDE flagged it right away. But `npm run build` passed and the page rendered fine.

Astro builds with Vite, which uses esbuild to compile TypeScript. esbuild *removes* type annotations without *checking* them. At runtime, the value was an array all along, so `posts.map(...)` worked. Only `astro check` (`npm run check`) runs the real type checker, and it reported the error.

I ran into the same thing in [Part 1](/posts/building-ktblog-part-1/). It's like Python type hints: the interpreter ignores them, and you need mypy to check them. That's different from Kotlin, where a type error stops the compile. So my verification step is always `npm run check && npm run build`.

The checker can also be silenced without any warning. My first version of the tag page used `new Map()` with no type parameters, which TypeScript infers as `Map<any, any>`. Everything that came out of it was `any`, so `astro check` reported zero errors because there was nothing left to check. Typing it as `Map<string, CollectionEntry<"posts">[]>` brought the checks back.

### A JavaScript trap: `toSorted()` without a comparator

Called without a comparison function, `toSorted()` (and `sort()`) converts every element to a string and compares the strings:

```javascript
[10, 9, 1].toSorted();  // [1, 10, 9]
```

With objects, it's worse: every post becomes the string `"[object Object]"`, so they all compare equal and the order doesn't change. Python's `sorted()` would raise a `TypeError` on a list of dicts. JavaScript just "works". I didn't notice at first because, with only one post, there was nothing to sort.

Always pass a comparator. Sorting posts newest first:

```typescript
posts.toSorted((a, b) => b.data.pubDate.getTime() - a.data.pubDate.getTime());
```

## Dark mode

I wanted two things from dark mode:

- On a first visit, follow the visitor's OS or browser setting.
- Once they flip the toggle, remember the choice on later visits.

### The `dark:` variant

Tailwind makes the styling part easy. `dark:` is a variant, like `hover:`:

```html
<body class="bg-white text-gray-900 dark:bg-gray-900 dark:text-white">
```

By default, `dark:` compiles to a `prefers-color-scheme: dark` media query, which follows the OS setting. That covers the first requirement with no JavaScript at all.

The second requirement breaks it. JavaScript can't change `prefers-color-scheme`, so a toggle button has nothing to flip. The [Tailwind docs](https://tailwindcss.com/docs/dark-mode#toggling-dark-mode-manually) show the fix: redefine the variant to look for a `dark` class instead.

```css
@custom-variant dark (&:where(.dark, .dark *));
```

Now `dark:` applies when the element itself, or any ancestor, has the `dark` class. JavaScript can add or remove that class on `<html>`, and every `dark:` class on the site follows without a single edit. `:where()` keeps the added selector's specificity at zero, so it doesn't change which rules win.

### Where does the code run?

My first attempt at "follow the OS on a first visit" was to check the setting in the component's frontmatter and pass `isDarkMode` down as a prop. It can't work, and the reason is the most useful thing I learned in this milestone.

An `.astro` file runs in three places:

1. **The frontmatter** (between the `---` fences) runs at **build time**, in Node. There's no `window`, `document`, or `localStorage`. It's like the Python code in a FastAPI route.
2. **The template** is rendered to static HTML, also at build time. It's like a Jinja template.
3. **A `<script>` tag** in the template is shipped to the browser and runs on **each visitor's machine**.

A static build produces one HTML file per page, shared by every visitor. At build time there is no visitor, so there's no OS setting or stored choice to read. Anything that depends on the visitor has to run in the browser.

The browser side uses [`localStorage`](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage), a small key-value store that each site gets in the browser. Values are strings, and they persist across visits:

```javascript
localStorage.setItem("theme", "dark");
localStorage.getItem("theme");  // "dark", or null if never set
```

Access can throw when the visitor blocks site data, so both calls are wrapped in `try`/`catch`. If reading fails, the site falls back to the OS setting.

### FOUC: a race condition

The theme has to be decided before the page is painted. Otherwise a dark-mode visitor sees a white page that turns dark a moment later. That flash is called FOUC, a flash of unstyled (here, wrongly styled) content.

A normal Astro `<script>` is bundled and emitted as `<script type="module">`, and module scripts are deferred: they run after the whole document has been parsed. The browser may paint before that:

```text
Deferred script (type="module"):
  parse <head> → parse <body> → paint (light) → script adds "dark" → repaint (dark)

Inline script in <head> (is:inline):
  parse <head> → script adds "dark" → parse <body> → paint (dark)
```

`<script is:inline>` tells Astro to emit the script as is. It becomes a classic script, which blocks the parser and runs immediately, before `<body>` even exists. `<html>` already exists at that point, so the script can set its class:

```html
<script is:inline>
    {
        let theme;
        try {
            theme = localStorage.getItem("theme");
        } catch {}  // Follow the OS setting when storage is blocked.

        let isDarkMode;
        switch (theme) {
            case "dark":
                isDarkMode = true;
                break;
            case "light":
                isDarkMode = false;
                break;
            default:
                isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
                break;
        }
        document.documentElement.classList.toggle("dark", isDarkMode);
    }
</script>
```

The surprise came when I tried to see the flash with the deferred version. I couldn't, not even on a throttled "Slow 4G" connection. The page is tiny, and the stylesheet download took longer than parsing, so the deferred script happened to run before the first paint anyway. Only after I moved the script to the end of `<body>` did the flash show up.

So the deferred version wasn't correct, just lucky. FOUC is a race between the first paint and the script, and a bigger page, a slower CPU, or a cached stylesheet could flip the result. The inline head script doesn't win the race by being faster. The HTML spec guarantees it runs before the body is parsed, so there's no race at all. It's the frontend version of fixing a race condition with a lock instead of hoping requests arrive in the right order.

Inline scripts have costs: they block parsing, and Astro doesn't process them, so they skip TypeScript entirely. I kept this one to a few lines and put the toggle button's click handler in a normal, type-checked `<script>`, since it only has to be ready before the first click.

### Code blocks: the cascade, part 2

Code blocks were the last thing still light in dark mode, and the trickiest.

Astro highlights code blocks at build time with [Shiki](https://shiki.style/), and Shiki writes the colors straight into the HTML as inline styles:

```html
<pre class="astro-code github-dark" style="background-color:#24292e;color:#e1e4e8">
```

That was already a problem before dark mode. The baseline's remaining accessibility failure was in code blocks: the default theme's comment color had a contrast of 3.04:1, below the 4.5:1 that [WCAG AA](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html) requires for normal text.

Shiki supports [dual themes](https://shiki.style/guide/dual-themes), configured in `astro.config.mjs`:

```javascript
markdown: {
    shikiConfig: {
        themes: {
            light: "github-light-high-contrast",
            dark: "github-dark-default",
        },
    },
},
```

With two themes, the light colors still go in real inline properties, but the dark colors only go in CSS custom properties:

```html
<span style="color:#66707B;--shiki-dark:#8B949E">
```

That's a comment token from one of my posts: `#66707B` in light mode, and `#8B949E` waiting in a variable for dark mode.

Shiki doesn't switch anything for you. You write the CSS that applies `--shiki-dark` in dark mode, and here's where the cascade comes back. An inline style beats every selector, no matter how specific. The only way for a stylesheet to override a normal inline declaration is `!important`:

```css
.astro-code {
    background-color: var(--color-gray-100) !important;
}

.astro-code,
.astro-code span {
    @variant dark {
        color: var(--shiki-dark) !important;
        background-color: var(--shiki-dark-bg) !important;
    }
}
```

`!important` is usually a code smell, but this is the case it exists for: overriding styles you don't control. `@variant dark` is the CSS-side form of `dark:`, so this rule follows the same `dark` class as everything else.

Two gotchas:

- Shiki's docs use the `.shiki` class in their examples, but Astro renders `class="astro-code"`. Copying the docs' CSS matched nothing.
- The gray background in light mode lowers the contrast of the comments. GitHub's regular light theme passes only on pure white (4.55:1) and drops to 4.13:1 on `gray-100`. I checked the bundled light themes, and almost none have comments that pass even on their own background. `github-light-high-contrast` stays at 4.58:1 on `gray-100`, and `github-dark-default` gets 6.15:1.

After this change, the post page's accessibility score went from 93 to 100 in both themes.

## RSS, SEO, and a missing 404

This part was mostly metadata, plus one real bug that had been hiding since M1.

### RSS

An RSS feed is an XML file (`/rss.xml`) that lists recent posts. Readers subscribe in a feed reader and get new posts without visiting the site. Think of it as a small, standardized, read-only public API.

Astro's [`@astrojs/rss`](https://docs.astro.build/en/recipes/rss/) package builds it from the content collection. The feed is an **endpoint**, a file under `src/pages/` that exports a `GET` handler instead of a template:

```typescript
// src/pages/rss.xml.ts (simplified)
export async function GET(context: APIContext) {
    const posts = sortByPubDateDesc(await getCollection("posts"));
    return rss({
        title: SITE_NAME,
        description: SITE_DESCRIPTION,
        site: context.site,
        items: posts.map((post) => ({
            title: post.data.title,
            pubDate: post.data.pubDate,
            description: post.data.description,
            link: `/posts/${post.id}/`,
        })),
    });
}
```

It looks like a FastAPI route, but there's an important difference: in a static build, the handler runs **once, at build time**, and its response is written to `dist/rss.xml`. No server runs it per request.

Feeds need absolute URLs, so `site: "https://ktblog.dev"` went into `astro.config.mjs`. The canonical URLs, Open Graph tags, and sitemap below all reuse it.

### SEO metadata

Static HTML from Astro is already easy for search engines to read, so SEO was mostly about filling in tags in the layout's `<head>`:

- `<meta name="description">`: the summary under the title in search results. This was one of the baseline failures.
- `<link rel="canonical">`: the one official URL for the page, so preview URLs and `?ref=` variants don't count as duplicates.
- [Open Graph](https://ogp.me/) tags (`og:title`, `og:description`, `og:type`...): the preview card when a link is shared.

I made `description` a **required** field in the content schema and in the layout's props. A post without one now fails the build, instead of showing up later as a Lighthouse warning. It's the same move as making a field non-nullable in a Pydantic model or a Kotlin data class.

The type checker can't catch everything, though. My first draft used `og:site`, which doesn't exist (it's `og:site_name`), and `name=` instead of `property=` for Open Graph tags. Both are valid TypeScript and valid HTML. Only comparing the output against the spec found them.

### The soft 404

The real bug showed up when Lighthouse tried to read `robots.txt`. The site didn't have one, yet `https://ktblog.dev/robots.txt` returned **200**, with the homepage's HTML. Lighthouse parsed each line of HTML as a robots rule and reported 22 errors.

Any unknown URL did the same:

```text
$ curl -I https://ktblog.dev/this-page-does-not-exist
HTTP/2 200
content-type: text/html; charset=utf-8
```

When a Cloudflare Pages project has no `404.html`, Pages assumes it's a single-page app and serves `index.html` for every unknown path, with status 200. That's called a soft 404: the page is missing, but the status code says otherwise. Search engines can index the "missing" pages as copies of the homepage. It also made a missing `rss.xml` look like a working URL that returned a web page.

It never showed up locally, because `astro preview` returns a real 404. Only production behaved differently.

The fix was one file, `src/pages/404.astro`. Astro builds it to `404.html`, and Pages then serves it with status 404. That's a different mindset from a backend: in FastAPI or Spring, the code decides the status code. On a static host, the platform decides, and all I control is which files exist.

With real 404s in place, the rest was quick:

- **Sitemap**: [`@astrojs/sitemap`](https://docs.astro.build/en/guides/integrations-guide/sitemap/), installed with `npx astro add sitemap`. It's an *integration*, not a library: instead of my code calling it, it hooks into the build and receives the list of built pages, much like a Spring Boot starter.
- **robots.txt**: another endpoint, so the sitemap URL comes from `site` instead of being hard-coded:

```typescript
// src/pages/robots.txt.ts (simplified)
export async function GET(context: APIContext): Promise<Response> {
    const sitemapUrl = new URL("sitemap-index.xml", context.site);
    return new Response(
        ["User-agent: *", "Allow: /", "", `Sitemap: ${sitemapUrl}`, ""].join("\n"),
    );
}
```

SEO went from 92 to 100 on every page.

## Search

The last feature was search. There's no backend yet (that's M3), so it had to run entirely in the browser.

### Pagefind

I compared two options:

- **A hand-written JSON index**: an endpoint like `rss.xml.ts` that writes every post's title and description to `search.json`, filtered in the browser. The whole file downloads on the first search and grows with every post.
- **[Pagefind](https://pagefind.app/)**: a tool that runs after the build, reads the generated HTML in `dist/`, and writes a full-text index split into small chunks. The browser downloads only the chunks a query needs. It also handles word stemming and ranking, and ships a search UI.

I picked Pagefind. It runs as an npm `postbuild` script, which npm runs automatically after `build`, so Cloudflare Pages runs it too:

```json
"scripts": {
    "build": "astro build",
    "postbuild": "pagefind --site dist"
}
```

The search itself runs in the browser. Pagefind's core is written in Rust and compiled to WebAssembly, and it only loads when someone actually searches, along with the index chunks. Every page still loads the search UI's script and stylesheet, about 47 KB compressed.

The flip side: the dev server never has an index, because `astro dev` doesn't run the build. Search only works after `npm run build && npm run preview`.

### The cascade, part 3

Pagefind 1.5 ships its UI as [web components](https://developer.mozilla.org/en-US/docs/Web/API/Web_components), custom HTML tags like `<pagefind-modal>`. I put the search button next to the dark mode toggle and spaced them with `space-x-10`. The spacing did nothing.

`space-x-10` works by putting a left margin on the children. Pagefind's stylesheet contains this rule, which isolates its button from the host site's CSS:

```css
pagefind-modal-trigger {
    all: initial;
}
```

`all: initial` resets every property, margin included. It's also *unlayered*, and Tailwind's utilities live in a layer. Remember the cascade from the Tailwind section: unlayered CSS beats layered CSS, whatever the specificity. So Pagefind's reset won.

The fix was to stop styling the child and let the parent own the spacing:

```html
<div class="flex items-center gap-3">
    <pagefind-modal-trigger />
    <ThemeToggle />
</div>
```

`gap` is a property of the flex container, so nothing a child does can reset it.

### Dark mode, again

The search modal stayed white in dark mode. Pagefind doesn't follow `prefers-color-scheme` or my `dark` class. It reads its colors from CSS variables (`--pf-background`, `--pf-text`, and so on) set on `:root`.

I overrode them in `global.css`, using the site's own colors, under the same `dark` class as everything else:

```css
:root.dark {
    --pf-text: var(--color-white);
    --pf-background: var(--color-gray-900);
    /* ... */
}
```

Why `:root.dark` and not `@variant dark`? Pagefind's defaults are also unlayered, on `:root`. `@variant dark` compiles to `:where(.dark, ...)`, which has zero specificity, so it would tie and the winner would depend on which stylesheet loads last. `:root.dark` is more specific, so it wins every time.

Navigation-mode Lighthouse runs still scored 100 on accessibility, because the modal is closed when the page loads. A **snapshot** run with the modal open and results showing told a different story: 96, with highlighted words at a contrast of 1.01:1, practically invisible. I'd left `--pf-mark` at its default, a near-black meant for the light theme. One more variable fixed it.

## Results

Lighthouse on a post page, before and after:

|                | Baseline (M1, local build) | After M2 (production) |
|----------------|----------------------------|-----------------------|
| Performance    | 100                        | 92–100                |
| Accessibility  | 93                         | 100                   |
| Best Practices | 100                        | 100                   |
| SEO            | 90                         | 100                   |

The "after" column covers light and dark themes on desktop and mobile. Accessibility and SEO are now perfect in every run.

Performance dipped a little, and the reason is honest: search added a render-blocking stylesheet, which costs about 200 ms of first paint on a throttled mobile run, and the search button changes size when its JavaScript loads, a small layout shift (CLS 0.03, still well under the 0.1 "good" threshold). I decided both were worth it for search.

The comparison isn't perfect. The baseline was measured locally, and the final runs are on production behind Cloudflare. But every score is at 90 or above, which was the goal.

## Lessons

- **Measure the right thing.** Dev server versus production build, preview URL versus production, a clean browser versus extensions. Half my surprising numbers came from measuring the wrong thing.
- **Read the output, not just the screen.** The wrong Open Graph tags looked fine in the browser and passed every check. Reading the generated HTML in `dist/` found them, along with a page `<title>` that printed a template string literally.
- **A green check needs evidence that it ran.** After the highlight fix, the snapshot report passed, but a passing report records no page elements. To be sure the modal had really been open, I checked the screenshot embedded in the report. A test with no assertions passes too.
- **Learn the cascade.** It explained three of the most confusing bugs in this milestone, and it's the one CSS topic I'd tell another backend developer to study first.

## What's next

M3 adds the first backend: a Go API for view counts and likes, running on a VPS with systemd and Caddy. The frontend gets its first interactive island to talk to it.
