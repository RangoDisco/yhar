<script lang="ts">
    import type {Snippet} from "svelte";

    import {page} from "$app/state"

    type Props = {
        title: string;
        url?: string | null;
        children: Snippet
    }

    let {
        title,
        url = $bindable(null),
        children
    }: Props = $props();
</script>

<section class="flex flex-col gap-6 w-full">
    <div class="flex justify-between items-center">
        <h1 class="text-lg md:text-xl text-foreground/90">{title}</h1>
        {#if url}
            <a class="text-sm text-muted-foreground md:text-base hover:underline" href="/users/{page.params.userID}/{url}">View all</a>
        {/if}
    </div>
    {@render children()}
</section>
