<script>
    import * as Avatar from "$lib/components/ui/avatar/index";
    import {Button} from "$lib/components/ui/button/index";
    import {Trash} from "@lucide/svelte";
    import dayjs from "dayjs";
    import relativeTime from "dayjs/plugin/relativeTime";

    import {page} from "$app/state"

    dayjs.extend(relativeTime);

    let {
        track,
        parentType
    } = $props();
</script>

<article class="flex items-center justify-between">
    <div class="flex gap-2 items-center w-full">
        <Avatar.Root class="rounded-md h-8 w-8">
            <Avatar.Image src={track.picture_url}
                          alt={`${track.title}'s picture`}/>
            <Avatar.Fallback>{track.title}</Avatar.Fallback>
        </Avatar.Root>
        <div class="flex flex-col w-full">
            <div class="flex gap-2 items-center">
                <span class="line-clamp-1 w-full min-w-0">{track.title}</span>
                <p class="text-sm text-muted-foreground whitespace-nowrap w-40">{dayjs(track.scrobbled_at).fromNow()}</p>
            </div>
            <div class="flex gap-1">
                {#if parentType === "artists"}
                    {#each track.artists as artist, i}
                        {#if i !== 0}
                            ·
                        {/if}
                        <a class="text-sm text-muted-foreground hover:underline h-6"
                           href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a>
                    {/each}
                {:else}
                    <a class="text-sm text-muted-foreground hover:underline"
                       href="/users/{page.params.userID}/top/albums/{track.album.id}">{track.album.title}</a>
                {/if}
            </div>
        </div>
    </div>
    <Button variant="outline" size="icon">
        <!--        TODO: Handle delete-->
        <Trash/>
    </Button>
</article>
