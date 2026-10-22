import { getCollection } from "astro:content";

export async function GET(): Promise<Response> {
    const posts = await getCollection("posts");
    const slugs = posts.map((post) => post.id);
    return new Response(
        JSON.stringify({ slugs }),
        { headers: { "Content-Type": "application/json" } }
    );
}
