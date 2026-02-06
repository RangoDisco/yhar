<script>
    import * as Avatar from "$lib/components/ui/avatar/index";

    let {
        index,
        title,
        parents,
        pictureUrl,
        contentType,
        scrobbleCount,
    } = $props();
</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-4 items-center">
        <h3 class="text-lg">{index + 1}</h3>
        <div class="flex gap-2">
            <Avatar.Root class="{contentType !== 'artists' ? 'rounded-md' : ''} h-12 w-12 ">
                <Avatar.Image src={pictureUrl}
                              alt={`${title}'s picture`}/>
                <Avatar.Fallback>{title}</Avatar.Fallback>
            </Avatar.Root>
            <div class="flex flex-col">
                <p>{title}</p>
                <div>
                    {#each parents as parent}
                        <span class="text-sm text-muted-foreground">{parent.name}</span>
                    {/each}
                </div>
            </div>
        </div>
    </div>
    <p class="text-muted-foreground">{scrobbleCount}</p>
</article>
