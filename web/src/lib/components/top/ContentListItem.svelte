<script lang="ts">
    import {page} from "$app/state";
    import type {Album, Artist, ContentType} from "#lib/types/content.js";
    import ContentPicture from "#lib/components/ContentPicture.svelte";

    type Props = {
        contentID?: string | null;
        index: number;
        title: string;
        parentType: "artists" | "albums";
        parents: Artist[] | Album[] | null;
        pictureUrl: string | null;
        contentType: ContentType;
        scrobbleCount: number | null;
    };

    let {
        contentID = $bindable(null),
        index,
        title,
        parentType,
        parents,
        pictureUrl,
        contentType,
        scrobbleCount
    }: Props = $props();
</script>

<article class="flex items-center justify-between gap-4">
    <div class="flex items-center gap-4">
        <h3 class="w-6 text-lg">{index}</h3>
        <div class="flex items-center gap-4">
            <ContentPicture
                    class="h-12 w-12"
                    url={pictureUrl}
                    contentType={contentType}
                    alt={`${title}'s picture`}
            />
            <div class="flex flex-col">
                {#if contentID}
                    <a
                            class="text-lg hover:underline"
                            href="/users/{page.params.userID}/top/{contentType}/{contentID}">{title}</a
                    >
                {:else}
                    <p class="text-lg">{title}</p>
                {/if}
                <div class="flex gap-1">
                    {#each parents as parent, i}
                        {#if i < 2}
                            {#if i !== 0}
                                ·
                            {/if}
                            <a
                                    class="text-muted-foreground hover:underline text-nowrap"
                                    href="/users/{page.params.userID}/top/{parentType}/{parent.id}"
                            >{parent.name ?? parent.title}</a
                            >
                        {/if}
                    {/each}
                </div>
            </div>
        </div>
    </div>
    <p class="text-muted-foreground">{scrobbleCount}</p>
</article>
