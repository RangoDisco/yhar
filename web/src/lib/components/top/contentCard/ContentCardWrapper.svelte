<script lang="ts">
    import {page} from "$app/state";
    import type {Snippet} from "svelte";

    type Props = {
        contentType: string;
        contentID?: string | null;
        children: Snippet;
    }

    let {
        contentType,
        contentID = $bindable(null),
        children,
    }: Props = $props();
</script>

<!-- Make card clickable if contentID is provided -->
{#if contentID}
    <a class="flex flex-col items-center justify-center gap-2"
       href="/users/{page.params.userID}/top/{contentType}/{contentID}" disabled>
        {@render children()}
    </a>
{:else}
    {@render children()}
{/if}
