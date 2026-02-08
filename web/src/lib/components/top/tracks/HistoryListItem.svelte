<script>
    import * as Avatar from "$lib/components/ui/avatar/index";
    import {Button} from "$lib/components/ui/button/index";
    import {Trash} from "@lucide/svelte";

    let {
        track,
        mode,
    } = $props();
</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-2 items-center">
        <Avatar.Root class="rounded-md h-8 w-8 ">
            <Avatar.Image src={track.picture_url}
                          alt={`${track.title}'s picture`}/>
            <Avatar.Fallback>{track.title}</Avatar.Fallback>
        </Avatar.Root>
        <div class="flex flex-col">
            <div class="flex gap-2">
                <p>{track.title}</p>
                {#if mode === "artists"}
                    {#each track.artists as artist}
                        <a href="/artists/{artist.id}">{artist.name}</a>
                    {/each}
                {:else}
                    <a href="/albums/{track.album.id}">{track.album.name}</a>
                {/if}
                <p class="text-sm text-muted-foreground"></p>
            </div>
            <p class="text-sm text-muted-foreground">{track.scrobbled_at}</p>
        </div>
    </div>
    <Button variant="outline" size="icon">
        <!--        TODO: Handle delete-->
        <Trash/>
    </Button>
</article>
