<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index";
    import {Button} from "$lib/components/ui/button/index";
    import {Trash, ListMusic} from "@lucide/svelte";
    import dayjs from "dayjs";
    import relativeTime from "dayjs/plugin/relativeTime";
    import * as Tooltip from "$lib/components/ui/tooltip/index.js";

    import {page} from "$app/state";
    import type {Scrobble} from "$lib/types/content";

    dayjs.extend(relativeTime);

    type Props = {
        scrobble: Scrobble,
        parentType: "artists" | "albums"
        handleDelete: (id: string) => void
    }

    let {
        scrobble,
        parentType,
        handleDelete
    }: Props = $props();
</script>

<article class="flex items-center justify-between">
    <div class="flex gap-4 items-center w-full min-w-0">
        <Avatar.Root class="rounded-md h-8 w-8 lg:h-10 lg:w-10">
            <Avatar.Image src={scrobble.track.picture_url}
                          alt={`${scrobble.track.title}'s picture`}/>
            <Avatar.Fallback class="rounded-md h-8 w-8">
                <ListMusic size={18} class="text-muted-foreground"/>
            </Avatar.Fallback>
        </Avatar.Root>
        <div class="flex flex-col w-full min-w-0">
            <div class="flex gap-2 items-center min-w-0">
                <span class="text-lg line-clamp-1 w-full min-w-0">{scrobble.track.title}</span>
                <Tooltip.Root>
                    <Tooltip.Trigger>
                        <p class="text-sm text-muted-foreground whitespace-nowrap w-40">{dayjs(scrobble.scrobbled_at).fromNow()}</p>
                    </Tooltip.Trigger>
                    <Tooltip.Content>
                        <p class="text-sm">{dayjs(scrobble.scrobbled_at).format("YYYY-MM-DD HH:mm ")}</p>
                    </Tooltip.Content>
                </Tooltip.Root>
            </div>
            <div class="flex gap-1">
                {#if parentType === "artists"}
                    {#each scrobble.track.artists as artist, i}
                        {#if i !== 0}
                            ·
                        {/if}
                        <a class="text-muted-foreground hover:underline h-6"
                           href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a>
                    {/each}
                {:else}
                    <a class="text-muted-foreground hover:underline"
                       href="/users/{page.params.userID}/top/albums/{scrobble.track.album.id}">{scrobble.track.album.title}</a>
                {/if}
            </div>
        </div>
    </div>
    <Button variant="outline" size="icon" onclick={() => handleDelete(scrobble.id)} aria-label="Delete scrobble">
        <Trash/>
    </Button>
</article>
