import type { APIContext } from "astro";
import rss, { type RSSOptions } from "@astrojs/rss";
import { getCollection } from "astro:content";
import { sortByPubDateDesc } from "../utils/posts.ts";

export async function GET(context: APIContext) {
    const site = context.site;
    if (site === undefined) {
        throw new Error("No site configured for RSS.");
    }
    const feedUrl = new URL("rss.xml", site);

    const posts = await getCollection("posts");
    const sortedPosts = sortByPubDateDesc(posts);
    const rssOptions: RSSOptions = {
        title: "KT's Blog",
        description: "Notes from building this blog from scratch: the problems I solved and the techniques I want to try next.",
        site: site,
        xmlns: {
            atom: "http://www.w3.org/2005/Atom",
        },
        customData: `<language>en</language><atom:link href="${ feedUrl }" rel="self" type="application/rss+xml" />`,
        items: sortedPosts.map((post) => ({
            title: post.data.title,
            pubDate: post.data.pubDate,
            link: `/posts/${post.id}/`,
            categories: post.data.tags,
        })),
    };
    return rss(rssOptions);
}
