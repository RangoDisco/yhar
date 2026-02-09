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

<section class="flex flex-col gap-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl md:text-3xl">{title}</h1>
        {#if url}
            <a class="text-sm text-muted-foreground md:text-base hover:underline" href="/users/{page.params.userID}/top/{url}">View all</a>
        {/if}
    </div>
    {@render children()}
</section>
