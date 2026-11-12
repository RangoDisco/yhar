<script lang="ts">
    import * as Tooltip from "#lib/components/ui/tooltip/index.ts";
    import {Button} from "#lib/components/ui/button/index.ts";
    import {Trash} from "@lucide/svelte";
    import dayjs from "dayjs";
    import relativeTime from "dayjs/plugin/relativeTime";

    import {page} from "$app/state";
    import type {Scrobble} from "#lib/types/content.ts";
    import ContentPicture from "#lib/components/ContentPicture.svelte";

    dayjs.extend(relativeTime);

    type Props = {
        scrobble: Scrobble;
        parentType: "artists" | "albums";
        handleDelete: (id: string) => void;
        shouldShowDeleteButton: boolean;
    };

    let {scrobble, parentType, handleDelete, shouldShowDeleteButton}: Props = $props();
</script>

<article class="flex items-center justify-between gap-4">
    <div class="flex w-full min-w-0 items-center gap-4">
        <ContentPicture
                url={`${scrobble.track.pictureDomain}/images/${scrobble.track.picturePath}`}
                alt={`${scrobble.track.title}`}
                contentType="tracks"
        />
        <div class="flex w-full min-w-0 justify-between">
            <div class="flex flex-col min-w-0">
                <a class="line-clamp-1 w-full min-w-0 text-lg hover:underline" href="/users/{page.params.userID}/top/tracks/{scrobble.track.id}">{scrobble.track.title}</a>
                <div class="flex gap-1">
                    {#if parentType === 'artists'}
                        {#each scrobble.track.artists as artist, i}
                            {#if i !== 0}
                                ·
                            {/if}
                            <a
                                    class="h-6 text-muted-foreground hover:underline"
                                    href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a
                            >
                        {/each}
                    {:else}
                        <a
                                class="text-muted-foreground hover:underline"
                                href="/users/{page.params.userID}/top/albums/{scrobble.track.album.id}"
                        >{scrobble.track.album.title}</a
                        >
                    {/if}
                </div>
            </div>

            <Tooltip.Root>
                <Tooltip.Trigger>
                    <p class="text-sm whitespace-nowrap text-muted-foreground">
                        {dayjs(scrobble.scrobbled_at).fromNow()}
                    </p>
                </Tooltip.Trigger>
                <Tooltip.Content>
                    <p class="text-sm">{dayjs(scrobble.scrobbled_at).format('YYYY-MM-DD HH:mm ')}</p>
                </Tooltip.Content>
            </Tooltip.Root>
        </div>
    </div>
    {#if shouldShowDeleteButton}
        <Button
                variant="outline"
                size="icon"
                onclick={() => handleDelete(scrobble.id)}
                aria-label="Delete scrobble"
        >
            <Trash/>
        </Button>
    {/if}
</article>
