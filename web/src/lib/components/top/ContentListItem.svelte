<script>
    import * as Avatar from "$lib/components/ui/avatar/index";

    import {page} from "$app/state";

    let {
        contentID = $bindable(null),
        index,
        title,
        parentType,
        parents,
        pictureUrl,
        contentType,
        scrobbleCount,
    } = $props();

</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-4 items-center">
        <h3 class="text-lg">{index + 1}</h3>
        <div class="flex gap-2 items-center">
            <Avatar.Root class="{contentType !== 'artists' ? 'rounded-md' : ''} h-12 w-12 ">
                <Avatar.Image src={pictureUrl}
                              alt={`${title}'s picture`}/>
                <Avatar.Fallback>{title}</Avatar.Fallback>
            </Avatar.Root>
            <div class="flex flex-col">
                {#if contentID}
                    <a class="hover:underline" href="/users/{page.params.userID}/top/{contentType}/{contentID}">{title}</a>
                {:else}
                    <p>{title}</p>
                {/if}
                <div class="flex gap-1">
                    {#each parents as parent, i}
                        {#if i < 3}
                            {#if i !== 0}
                                ·
                            {/if}
                            <a class="text-sm text-muted-foreground hover:underline"
                               href="/users/{page.params.userID}/top/{parentType}/{parent.id}">{parent.name ?? parent.title}</a>
                        {/if}
                    {/each}
                </div>
            </div>
        </div>
    </div>
    <p class="text-muted-foreground">{scrobbleCount}</p>
</article>
