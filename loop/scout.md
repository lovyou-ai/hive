# Scout Report — Iteration 8

## Map (from code + infra)

Read state.md. Checked Fly secrets — DATABASE_URL, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, AUTH_REDIRECT_URL all set. App is running, health checks passing. The "visitors get 503" worry was wrong — infrastructure is wired up.

Reassessed: since the infra gap doesn't exist, the next highest-leverage discoverability improvement is sitemap.xml + robots.txt. The site has 250+ pages but search engines have no way to discover them.

## Gap Type

Missing infrastructure — no sitemap.xml, no robots.txt.

## The Gap

Search engines can't efficiently discover 305 pages (43 blog posts, 14 layers, 201+ primitives, 28 agent primitives, 13 grammars, static pages) without a sitemap. No robots.txt means no pointer to the sitemap.

## Why This Gap

SEO meta tags (iteration 7) give each page good metadata, but search engines still need to find the pages. A sitemap is the standard mechanism. 305 indexed pages with proper meta descriptions = significant organic search potential for long-tail queries.

## Filled Looks Like

GET /sitemap.xml returns valid XML with all public URLs. GET /robots.txt allows all crawlers and points to the sitemap.
