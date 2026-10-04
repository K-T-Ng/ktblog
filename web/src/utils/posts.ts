import type { CollectionEntry } from "astro:content";

export function sortByPubDateDesc(posts: CollectionEntry<"posts">[]): CollectionEntry<"posts">[] {
    return posts.toSorted(
        (a, b) => b.data.pubDate.getTime() - a.data.pubDate.getTime()
    );
}
