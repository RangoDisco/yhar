<script lang="ts">
    import type {Snippet} from "svelte";

    import {page} from "$app/state";

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
        <h1 class="text-lg md:text-xl text-foreground/90">
            {#if url}
                <a class="text-lg md:text-xl text-foreground/90 hover:underline"
                   href="/users/{page.params.userID}/{url}">
                    {title}
                </a>
            {:else}
                {title}
            {/if}
        </h1>
    </div>
    {@render children()}
</section>
