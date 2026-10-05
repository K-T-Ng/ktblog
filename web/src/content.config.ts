import { defineCollection } from 'astro:content';
import { glob } from 'astro/loaders';
import { z } from 'astro/zod';

const posts = defineCollection({
    loader: glob({ base: "./src/content/posts", pattern: "*.md" }),
    schema: z.object({
        title: z.string(),
        pubDate: z.coerce.date(),
        description: z.string(),
        tags: z.array(z.enum(["astro", "cloudflare", "node-js"]))
            .default([])
            .refine(
                (tags) => tags.length === new Set(tags).size, {
                    error: "Duplicate tags are not allowed."
                }
            ),
    }),
});

export const collections = { posts };
