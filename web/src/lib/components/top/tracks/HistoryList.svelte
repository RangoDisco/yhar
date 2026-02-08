<script lang="ts">
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";

    type Props = {
        tracks: [],
        mode: "artists" | "albums"
    }

    let {
        tracks,
        mode
    }: Props = $props();
</script>

<section class="flex flex-col gap-2 md:w-1/2">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl">History</h1>
        <a class="text-sm text-muted-foreground" href="history">View all</a>
    </div>
    <div class="flex flex-col gap-2">
        {#each tracks as track}
            <HistoryListItem track={track} mode={mode}/>
        {/each}
    </div>
</section>
