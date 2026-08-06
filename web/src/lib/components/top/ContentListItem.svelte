<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index";
    import disc from "$lib/assets/placeholders/disc-3.svg";
    import user from "$lib/assets/placeholders/user.jpg";

    import {page} from "$app/state";
    import type {Album, Artist} from "$lib/types/content";

    type Props = {
        contentID?: string | null
        index: number
        title: string
        parentType: "artists" | "albums"
        parents: Artist[] | Album[] | null
        pictureUrl: string | null
        contentType: string
        scrobbleCount: number | null
    }

    let {
        contentID = $bindable(null),
        index,
        title,
        parentType,
        parents,
        pictureUrl,
        contentType,
        scrobbleCount,
    }: Props = $props();

</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-4 items-center">
        <h3 class="text-lg w-6">{index}</h3>
        <div class="flex gap-4 items-center">
            <Avatar.Root class="{contentType !== 'artists' ? 'rounded-md' : ''} h-12 w-12 md:h-16 md:w-16 aspect-square">
                <Avatar.Image src={pictureUrl}
                              alt={`${title}'s picture`}/>
                <Avatar.Fallback>
                    <img class="rounded-full h-12 w-12" src={contentType === "artists" ? user: disc}
                         alt="Artist placeholder"/>
                </Avatar.Fallback>
            </Avatar.Root>
            <div class="flex flex-col">
                {#if contentID}
                    <a class="text-lg hover:underline"
                       href="/users/{page.params.userID}/top/{contentType}/{contentID}">{title}</a>
                {:else}
                    <p class="text-lg">{title}</p>
                {/if}
                <div class="flex gap-1">
                    {#each parents as parent, i}
                        {#if i < 3}
                            {#if i !== 0}
                                ·
                            {/if}
                            <a class="text-muted-foreground hover:underline"
                               href="/users/{page.params.userID}/top/{parentType}/{parent.id}">{parent.name ?? parent.title}</a>
                        {/if}
                    {/each}
                </div>
            </div>
        </div>
    </div>
    <p class="text-muted-foreground">{scrobbleCount}</p>
</article>
