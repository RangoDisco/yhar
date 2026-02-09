<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index";

    import {page} from "$app/state";

    let {
        title,
        contentID = $bindable(null),
        pictureUrl,
        contentType
    } = $props();

</script>

<article class="w-[30%] md:w-[15%] flex flex-col items-center justify-center gap-2">
    {#if contentID}
        <a class="flex flex-col items-center justify-center gap-2"
           href="/users/{page.params.userID}/top/{contentType}/{contentID}" disabled>
            <img class="{contentType !== 'artists' ? 'rounded-md' : 'rounded-full'} w-full aspect-square"
                 src={pictureUrl}
                 alt={`${title}'s picture`}/>
            <p class="text-center mb-auto">{title}</p>
        </a>
    {:else}
        <img class="{contentType !== 'artists' ? 'rounded-md' : 'rounded-full'} w-full aspect-square" src={pictureUrl}
             alt={`${title}'s picture`}/>
        <p class="text-center mb-auto">{title}</p>
    {/if}
</article>
