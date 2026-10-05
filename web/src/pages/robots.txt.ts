import type { APIContext } from "astro";

export async function GET(context: APIContext): Promise<Response> {
    const site = context.site;
    if (site === undefined) {
        throw new Error("No site configured for robots.txt.");
    }

    const siteMapUrl = new URL("sitemap-index.xml", site);
    const contents = [
        "User-agent: *",
        "Allow: /",
        "",
        `Sitemap: ${ siteMapUrl }`,
        "",
    ];
    return new Response(contents.join("\n"));
}
